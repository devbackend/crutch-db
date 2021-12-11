package server

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"

	"github.com/devbackend/crutch-db/internal/operation"
	"github.com/devbackend/crutch-db/internal/parser"
)

type Server interface {
	Start() error
	Stop(ctx context.Context) error
}

func New(options ...Option) Server {
	srv := new(tcpServer)

	for _, option := range options {
		option(srv)
	}

	return srv
}

type tcpServer struct {
	host      string
	port      int
	cmdParser parser.Parser
	cmdRunner *operation.Runner
}

func (s *tcpServer) Start() error {
	addr := net.JoinHostPort(s.host, strconv.Itoa(s.port))

	ln, err := net.Listen("tcp", addr)
	defer func() {
		err := ln.Close()
		if err != nil {
			log.Println("err", err)
		}
	}()

	if err != nil {
		return err
	}

	fmt.Println("started on", addr)

	for {
		conn, err := ln.Accept()
		if err != nil {
			return err
		}

		fmt.Println("connection established")

		go s.handleConn(conn)
	}
}

func (s *tcpServer) Stop(ctx context.Context) error {
	return nil
}

func (s *tcpServer) handleConn(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err == io.EOF {
			fmt.Println("connection closed")
			return
		}

		if err != nil {
			s.logError(err)
			continue
		}

		cmd, err := s.cmdParser.Parse(message)
		if err != nil {
			_, err = conn.Write([]byte("parse error: " + err.Error() + "\n"))
			if err != nil {
				log.Println("err", err)
			}

			continue
		}

		res, err := s.cmdRunner.Run(cmd)
		if err != nil {
			_, err = conn.Write([]byte("ERR " + err.Error() + "\n"))
			if err != nil {
				log.Println("err", err)
			}

			continue
		}

		_, err = conn.Write([]byte("RES " + res + "\n"))
		if err != nil {
			log.Println("err", err)
		}
	}
}

func (s *tcpServer) logError(err error) {
	log.Println(err)
}
