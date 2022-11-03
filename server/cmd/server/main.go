package main

import (
	"log"

	"github.com/gin-gonic/gin"

	"intelops.io/server/pkg/api"
)

func main() {
	r := gin.Default()
	api.Setup(r)
	if err := r.Run(); err != nil {
		log.Fatalf("failed to start server : %s", err.Error())
	}
}
