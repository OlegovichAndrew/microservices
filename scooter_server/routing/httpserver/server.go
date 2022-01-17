package httpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"net"
	"net/http"
	"scooter_micro/proto"
	"time"
)

const (
	defaultReadTimeout     = 5 * time.Second
	defaultWriteTimeout    = 5 * time.Second
	defaultIdleTimeout     = 30 * time.Second
	defaultShutdownTimeout = 3 * time.Second
	defaultAddr            = ":8085"
)

//Client is a client's struct who connects to the "scooter-run" page.
type Client struct {
	w    io.Writer
	done chan struct{}
}

//Server is a struct of the http-server which has a channel for gRPC connection.
type Server struct {
	server          *http.Server
	notify          chan error
	shutdownTimeout time.Duration
	client          map[int]*Client
	taken           map[int]bool
	codes           map[int]int
	in              chan *proto.ClientMessage
	Structure       chan *proto.ScooterClient
	ScooterIdMap    map[uint64]proto.ScooterService_RegisterServer
	*proto.UnimplementedScooterServiceServer
}

type Option func(*Server)

//New creates and starts the http-server
func New(handler http.Handler, opts ...Option) *Server {
	httpServer := &http.Server{
		Handler:      handler,
		ReadTimeout:  defaultReadTimeout,
		WriteTimeout: defaultWriteTimeout,
		IdleTimeout:  defaultIdleTimeout,
		Addr:         defaultAddr,
	}

	server := &Server{
		server:          httpServer,
		notify:          make(chan error, 1),
		shutdownTimeout: defaultShutdownTimeout,
		client:          make(map[int]*Client),
		taken:           make(map[int]bool),
		codes:           make(map[int]int),
		in:              make(chan *proto.ClientMessage),
		Structure:       make(chan *proto.ScooterClient),
		ScooterIdMap:    make(map[uint64]proto.ScooterService_RegisterServer),
	}

	for _, opt := range opts {
		opt(server)
	}

	server.run()

	return server
}

func (s *Server) Notify() <-chan error {
	return s.notify
}

func (s *Server) Shutdown() error {
	ctx, cancel := context.WithTimeout(context.Background(), s.shutdownTimeout)
	defer cancel()

	return s.server.Shutdown(ctx)
}

func Port(port string) Option {
	return func(s *Server) {
		s.server.Addr = net.JoinHostPort("", port)
	}
}

//ScooterHandler is a special handler which adds a new stream client to the server.
func (s *Server) ScooterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new client connected")
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	client := &Client{
		w:    w,
		done: make(chan struct{}),
	}
	s.AddClient(client)

	<-client.done
	fmt.Println("connection closed")
}

//AddClient is a Server's function for adding attached Client.
func (s *Server) AddClient(c *Client) {
	s.client[1] = c
}

func (s *Server) MatchStreamToScooterId(ctx context.Context, stream proto.ScooterService_RegisterServer) {
	for k, v := range s.ScooterIdMap {
		if v == nil {
			s.ScooterIdMap[k] = stream
			break
		}
	}
}

//Register is a function for implementing gRPC-service.
func (s *Server) Register(stream proto.ScooterService_RegisterServer) error {
	s.MatchStreamToScooterId(context.Background(), stream)
	fmt.Println(s.ScooterIdMap)
	go func() {
		for {
			msg, err := stream.Recv()
			if err != nil {
				//fmt.Printf("Error: %v", err)
				err = status.Errorf(codes.Internal, "unexpected error %v", err)
			}

			if msg.GetId() != 0 {
				s.in <- msg
			}
			continue
		}
	}()

	go func() {
		select {
		case data := <- s.Structure:
			err := stream.Send(data)
			if err != nil {
				log.Printf("send error %v", err)
			}
		}

	}()

	return nil
}

//Receive is the function which receive a message from the gRPC stream and direct it to the Server's 'in' channel.
func (s *Server) Receive(stream proto.ScooterService_ReceiveServer) error {
	var err error

	for {
		msg, err := stream.Recv()
		if err != nil {
			fmt.Println(err)
			err = status.Errorf(codes.Internal, "unexpected error %v", err)
			break
		}

		s.in <- msg

	}

	return err
}

//run runs the Server and wait for messages into the channel. Then encode them and print to the console.
func (s *Server) run() {
	go func() {
		for {
			select {
			case msg := <-s.in:
				var buf bytes.Buffer
				json.NewEncoder(&buf).Encode(msg)

				for _, client := range s.client {

					go func(c *Client) {
						if _, err := fmt.Fprintf(c.w, "data: %v\n\n", buf.String()); err != nil {
							fmt.Println(err)
							delete(s.client, 1)
							close(c.done)
							return
						}

						if f, ok := c.w.(http.Flusher); ok {
							f.Flush()
						}
						fmt.Printf("data: %v\n", buf.String())
					}(client)
				}
			}
		}
	}()
}
