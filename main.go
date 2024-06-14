package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	amqp "github.com/rabbitmq/amqp091-go"
)

type InfoConnection struct {
	ID   string `json:"id"`
	Addr string `json:"address"`
}

type ClientMessage struct {
	ID          string   `json:"id"`
	ChatID      string   `json:"chatId"`
	SenderID    string   `json:"senderId"`
	Content     string   `json:"content"`
	SentAt      string   `json:"sentAt"`
	Attachments []string `json:"attachments"`
}

type Server struct {
	conns    map[string][]*websocket.Conn // Map untuk menyimpan koneksi berdasarkan ID pengguna
	connsMu  sync.Mutex
	nextID   int
	upgrader websocket.Upgrader
}

func NewServer() *Server {
	return &Server{
		conns:    make(map[string][]*websocket.Conn),
		nextID:   1,
		upgrader: websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }},
	}
}

func (s *Server) handleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := s.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	// Baca pesan pertama dari klien untuk mendapatkan ID pengguna
	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading user ID:", err)
		return
	}

	var userID struct {
		ID string `json:"id"`
	}
	if err := json.Unmarshal(msg, &userID); err != nil {
		log.Println("Error decoding user ID:", err)
		return
	}

	// Tambahkan koneksi ke dalam map conns dengan ID pengguna yang unik
	s.connsMu.Lock()
	connID := fmt.Sprintf("%s-%d", userID.ID, s.nextID)
	s.nextID++
	if _, exists := s.conns[userID.ID]; !exists {
		s.conns[userID.ID] = make([]*websocket.Conn, 0)
	}
	s.conns[userID.ID] = append(s.conns[userID.ID], conn)
	s.connsMu.Unlock()

	// Kirim ID koneksi ke klien
	conn.WriteMessage(websocket.TextMessage, []byte(connID))

	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			// Hapus koneksi dari map saat koneksi tertutup atau ada error
			s.connsMu.Lock()
			for i, c := range s.conns[userID.ID] {
				if c == conn {
					s.conns[userID.ID] = append(s.conns[userID.ID][:i], s.conns[userID.ID][i+1:]...)
					break
				}
			}
			if len(s.conns[userID.ID]) == 0 {
				delete(s.conns, userID.ID)
			}
			s.connsMu.Unlock()
			return
		}

		var clientMsg ClientMessage
		if err := json.Unmarshal(msg, &clientMsg); err != nil {
			log.Println("Error decoding message:", err)
			continue
		}

		// Kirim pesan hanya ke koneksi tujuan
		s.connsMu.Lock()
		found := false
		for _, destConn := range s.conns[clientMsg.ChatID] {
			// Kirim pesan sesuai dengan struktur yang diberikan
			destConn.WriteJSON(clientMsg)
			// Kirim respons ke klien bahwa pesan berhasil dikirim
			conn.WriteMessage(websocket.TextMessage, []byte("Pesan berhasil dikirim"))
			found = true
		}
		s.connsMu.Unlock()

		if !found {
			conn.WriteMessage(websocket.TextMessage, []byte("Penerima pesan tidak ditemukan"))
		}
	}
}

func (s *Server) getConnectionInfo() []InfoConnection {
	s.connsMu.Lock()
	defer s.connsMu.Unlock()

	var connections []InfoConnection

	for connID, connList := range s.conns {
		for _, conn := range connList {
			info := InfoConnection{
				ID:   connID,
				Addr: conn.RemoteAddr().String(),
			}
			connections = append(connections, info)
		}
	}

	return connections
}

func connectionsHandler(w http.ResponseWriter, _ *http.Request, server *Server) {
	connections := server.getConnectionInfo()

	// Mengembalikan informasi koneksi dalam format JSON
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(connections)
}

func (s *Server) consume() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"hello", // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack (ubah ke false)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	var forever chan struct{}

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			var clientMsg ClientMessage
			if err := json.Unmarshal(d.Body, &clientMsg); err != nil {
				log.Println("Error decoding message:", err)
				d.Nack(false, true) // Nack message if decoding fails
				continue
			}

			s.connsMu.Lock()
			for _, destConn := range s.conns[clientMsg.ChatID] {
				// Kirim pesan sesuai dengan struktur yang diberikan
				destConn.WriteJSON(clientMsg)
			}
			s.connsMu.Unlock()

			// Konfirmasi bahwa pesan telah berhasil di-handle
			if err := d.Ack(false); err != nil {
				log.Printf("Error acknowledging message: %s", err)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	<-forever
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func main() {
	server := NewServer()

	go server.consume()

	http.HandleFunc("/ws", server.handleWebSocket)
	http.HandleFunc("/connections", func(w http.ResponseWriter, r *http.Request) {
		connectionsHandler(w, r, server)
	})

	fmt.Println("Server is listening on port 8181")
	log.Fatal(http.ListenAndServe(":8181", nil))
}
