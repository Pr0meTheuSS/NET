package server

import (
	"errors"
	"fmt"
	"log"
	datarepresenter "main/DataRepresenter"
	"main/connection"
	"main/speedmeasure"
	"net"
	"runtime"
	"strconv"
)

type Server struct {
	ip                     string
	port                   int
	listenTimeoutInSeconds int
	listener               net.Listener
}

func NewServer(ip string, port, listenTimeoutInSeconds int) *Server {
	return &Server{
		ip:                     ip,
		port:                   port,
		listenTimeoutInSeconds: listenTimeoutInSeconds,
	}
}

func (s Server) ListenAndServe() error {
	var err error
	s.listener, err = net.Listen("tcp", s.ip+":"+strconv.Itoa(s.port))
	defer s.listener.Close()

	if err != nil {
		return errors.New(fmt.Sprintf("Cannot listen port. Error:%s", err))
	}

	return s.serveConnections()
}

const (
	defaultUploadsDir = "uploads"
)

func (s Server) serveConnections() error {
	dataRepresenter := datarepresenter.NewDataRepresenterCli()
	for {
		fmt.Println("Goroutine amount:", runtime.NumGoroutine())
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("Accept connection error:", err.Error())
			continue
		}

		log.Println("Accept new connection from client")

		go func() {
			c := connection.NewServerSideConnection(conn, &FileConsumer{})
			ch := dataRepresenter.Register()

			sm := speedmeasure.NewSpeedMeasurer(c, ch)
			sm.MeasureSpeed()
			c.ServerServe()
			dataRepresenter.Unregister(ch)
			sm.Close()
		}()
		fmt.Println("Goroutine amount:", runtime.NumGoroutine())
	}
}
