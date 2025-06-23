package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/ruziba3vich/face_embedding_gateway/genprotos/face_recognition_service"
	handler "github.com/ruziba3vich/face_embedding_gateway/internal/http"
	"github.com/ruziba3vich/face_embedding_gateway/internal/service"
	"github.com/ruziba3vich/face_embedding_gateway/pkg/helper"
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

	faceRecognitionHandler := handler.NewFaceRecognitionHandler(client, logger)

	db, err := helper.NewDB("")

	milvusClient, err := helper.CreateMilvusClientAndCollection("face_embeddings", "", "512")
	if err != nil {
		logger.Fatal(err)
	}

	userService := service.NewService(db, milvusClient)
	userHandler := handler.NewUserHandler(userService, client, logger)

	router := gin.Default()
	router.POST("/embedd", faceRecognitionHandler.HandleImageEmbedding)
	router.POST("/ceate-user", userHandler.CreateUser)
	router.GET("/get-user", userHandler.GetUserByID)
	router.DELETE("/delete-user", userHandler.DeleteUserByID)
	router.POST("/create-pic", userHandler.HandleImageEmbedding)

	if err := router.Run(":7777"); err != nil {
		logger.Fatal(err)
	}
}
