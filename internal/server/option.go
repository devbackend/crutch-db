package server

import (
	"github.com/devbackend/crutch-db/internal/operation"
	"github.com/devbackend/crutch-db/internal/parser"
)

type Option func(srv *tcpServer)

// WithPort add host to tcpServer
func WithPort(port int) Option {
	return func(srv *tcpServer) {
		srv.port = port
	}
}

// WithParser add command parser to tcpServer
func WithParser(parser parser.Parser) Option {
	return func(srv *tcpServer) {
		srv.cmdParser = parser
	}
}

// WithOperationRunner add command runner to tcpServer
func WithOperationRunner(runner *operation.Runner) Option {
	return func(srv *tcpServer) {
		srv.cmdRunner = runner
	}
}
