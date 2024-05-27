package websocket

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	socketio "github.com/googollee/go-socket.io"
	amqp "github.com/rabbitmq/amqp091-go"
)

type MessageData struct {
	ID      string `json:"id"`
	Message string `json:"message"`
}

type InfoConnection struct {
	ID   string `json:"id"`
	Addr string `json:"address"`
}

type ClientMessage struct {
	DestinationID string `json:"destination_id"`
	Content       string `json:"content"`
}

type Server struct {
	conns   map[string]socketio.Conn // Map to store connections by ID
	connsMu sync.Mutex
	nextID  int
}

func NewServer() *Server {
	return &Server{
		conns:  make(map[string]socketio.Conn),
		nextID: 1,
	}
}

func (s *Server) handleSocketIO(server *socketio.Server) {
	server.OnConnect("/", func(conn socketio.Conn) error {
		s.connsMu.Lock()
		connID := fmt.Sprintf("%d", s.nextID)
		s.nextID++
		s.conns[connID] = conn
		s.connsMu.Unlock()

		conn.SetContext(connID)
		conn.Emit("connectionID", connID)
		log.Printf("Connected: %s", connID)
		return nil
	})

	server.OnEvent("/", "message", func(conn socketio.Conn, msg string) {
		var clientMsg ClientMessage
		if err := json.Unmarshal([]byte(msg), &clientMsg); err != nil {
			log.Println("Error decoding message:", err)
			return
		}

		s.connsMu.Lock()
		destConn, found := s.conns[clientMsg.DestinationID]
		s.connsMu.Unlock()

		if found {
			destConn.Emit("message", clientMsg.Content)
			conn.Emit("status", "Message sent successfully")
		} else {
			conn.Emit("status", "Destination not found")
		}
	})

	server.OnDisconnect("/", func(conn socketio.Conn, reason string) {
		connID := conn.Context().(string)
		s.connsMu.Lock()
		delete(s.conns, connID)
		s.connsMu.Unlock()
		log.Printf("Disconnected: %s", connID)
	})
}

func (s *Server) getConnectionInfo() []InfoConnection {
	s.connsMu.Lock()
	defer s.connsMu.Unlock()

	var connections []InfoConnection
	for connID, conn := range s.conns {
		info := InfoConnection{
			ID:   connID,
			Addr: conn.RemoteAddr().String(),
		}
		connections = append(connections, info)
	}
	return connections
}

func connectionsHandler(w http.ResponseWriter, _ *http.Request, server *Server) {
	connections := server.getConnectionInfo()

	// Return connection information in JSON format
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(connections); err != nil {
		log.Println("JSON Encode error:", err)
	}
}

func (s *Server) consume() {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	failOnError(err, "Failed to connect to RabbitMQ")
	defer conn.Close()

	ch, err := conn.Channel()
	failOnError(err, "Failed to open a channel")
	defer ch.Close()

	q, err := ch.QueueDeclare(
		"ClientQueue", // name
		false,         // durable
		false,         // delete when unused
		false,         // exclusive
		false,         // no-wait
		nil,           // arguments
	)
	failOnError(err, "Failed to declare a queue")

	msgs, err := ch.Consume(
		q.Name, // queue
		"",     // consumer
		false,  // auto-ack (change to false)
		false,  // exclusive
		false,  // no-local
		false,  // no-wait
		nil,    // args
	)
	failOnError(err, "Failed to register a consumer")

	go func() {
		for d := range msgs {
			log.Printf("Received a message: %s", d.Body)

			// Unmarshal the message body to determine the destination
			var clientMsg ClientMessage
			if err := json.Unmarshal(d.Body, &clientMsg); err != nil {
				log.Println("Error decoding message:", err)
				d.Nack(false, true)
				continue
			}

			// Forward the message to the appropriate WebSocket client
			s.connsMu.Lock()
			destConn, found := s.conns[clientMsg.DestinationID]
			s.connsMu.Unlock()

			if found {
				destConn.Emit("message", clientMsg.Content)
				log.Printf("Message forwarded to client %s", clientMsg.DestinationID)
				d.Ack(false)
			} else {
				log.Printf("Destination client %s not found", clientMsg.DestinationID)
				d.Nack(false, true)
			}
		}
	}()

	log.Printf(" [*] Waiting for messages. To exit press CTRL+C")
	select {}
}

func failOnError(err error, msg string) {
	if err != nil {
		log.Panicf("%s: %s", msg, err)
	}
}

func StartWS() {
	server := NewServer()

	socketIOServer := socketio.NewServer(nil)
	server.handleSocketIO(socketIOServer)

	http.Handle("/socket.io/", socketIOServer)
	http.HandleFunc("/connections", func(w http.ResponseWriter, r *http.Request) {
		connectionsHandler(w, r, server)
	})

	go func() {
		if err := socketIOServer.Serve(); err != nil {
			log.Fatalf("SocketIO server error: %s", err)
		}
	}()
	defer socketIOServer.Close()

	// Start the RabbitMQ consumer in a separate goroutine
	go server.consume()

	fmt.Println("WebSocket server is listening on port 8181")
	log.Fatal(http.ListenAndServe(":8181", nil))
}
