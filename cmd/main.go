package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ruziba3vich/face_embedding_gateway/genprotos/face_recognition_service"
	handler "github.com/ruziba3vich/face_embedding_gateway/internal/http"
	lgg "github.com/ruziba3vich/prodonik_lgger"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	logger, err := lgg.NewLogger("app.log")
	if err != nil {
		log.Fatal(err)
	}

	clientConn, err := grpc.NewClient("localhost:7178", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		logger.Fatal(err)
	}
	defer clientConn.Close()

	client := face_recognition_service.NewFaceEmbedderClient(clientConn)

	handler := handler.NewFaceRecognitionHandler(client, logger)

	router := gin.Default()
	router.POST("/embedd", handler.HandleImageEmbedding)

	if err := router.Run(":7777"); err != nil {
		logger.Fatal(err)
	}
}
