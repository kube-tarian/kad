package server

import (
	"github.com/gin-gonic/gin"
	"github.com/kube-tarian/kad/server/pkg/client"
	"github.com/kube-tarian/kad/server/pkg/logging"
)

type Server struct {
	log    logging.Logger
	client *client.Agent
}

func NewServer(log logging.Logger) (*Server, error) {
	client, err := client.NewAgent()
	if err != nil {
		log.Errorf("failed to connect agent internal error", err)
		return nil, err
	}

	return &Server{
		log:    log,
		client: client,
	}, nil
}

func (s *Server) Close(c *gin.Context) {
	s.client.Close()
}
