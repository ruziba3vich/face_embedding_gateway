package handler

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/ruziba3vich/face_embedding_gateway/genprotos/face_recognition_service"
	"github.com/ruziba3vich/face_embedding_gateway/internal/models"
	"github.com/ruziba3vich/face_embedding_gateway/internal/service"
	lgg "github.com/ruziba3vich/prodonik_lgger"
)

type FaceRecognitionHandler struct {
	faceEmbedderClient face_recognition_service.FaceEmbedderClient
	logger             *lgg.Logger
}

func NewFaceRecognitionHandler(faceEmbedderClient face_recognition_service.FaceEmbedderClient, logger *lgg.Logger) *FaceRecognitionHandler {
	return &FaceRecognitionHandler{
		faceEmbedderClient: faceEmbedderClient,
		logger:             logger,
	}
}

func (h *FaceRecognitionHandler) HandleImageEmbedding(c *gin.Context) {
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read image"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := h.faceEmbedderClient.GetEmbedding(ctx, &face_recognition_service.ImageRequest{Image: imageData})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("gRPC error: %v", err)})
		return
	}

	if resp.Error != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.Error})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"embedding_length": len(resp.Embedding),
		"embedding":        resp.Embedding,
	})
}

type UserHandler struct {
	service            *service.Service
	faceEmbedderClient face_recognition_service.FaceEmbedderClient
	logger             *lgg.Logger
}

func (h *UserHandler) HandleImageEmbedding(c *gin.Context) {
	file, _, err := c.Request.FormFile("image")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "image file is required"})
		return
	}
	defer file.Close()

	imageData, err := io.ReadAll(file)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to read image"})
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	objectID := c.Query("object_id")
	if len(objectID) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "object_id is not provided"})
		return
	}

	resp, err := h.faceEmbedderClient.GetEmbedding(ctx, &face_recognition_service.ImageRequest{Image: imageData})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": fmt.Sprintf("gRPC error: %v", err)})
		return
	}

	if resp.Error != "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": resp.Error})
		return
	}

	if err := h.service.StoreVector(ctx, objectID, resp.Embedding); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"response": "successfully stored"})
}

func NewUserHandler(service *service.Service, faceEmbedderClient face_recognition_service.FaceEmbedderClient, logger *lgg.Logger) *UserHandler {
	return &UserHandler{
		service:            service,
		faceEmbedderClient: faceEmbedderClient,
		logger:             logger,
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var user models.User
	if err := c.ShouldBindJSON(&user); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON: " + err.Error()})
		return
	}

	if err := h.service.Create(&user); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create user"})
		return
	}

	c.JSON(http.StatusCreated, user)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	id := c.Param("id")
	user, err := h.service.GetByID(id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}
	if user == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "user not found"})
		return
	}
	c.JSON(http.StatusOK, user)
}

func (h *UserHandler) DeleteUserByID(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.DeleteByID(id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "user deleted"})
}

/*

	conn, err := grpc.Dial("localhost:7178", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "gRPC connection failed"})
		return
	}
	defer conn.Close()

	client := face_recognition_service.NewFaceEmbedderClient(conn)
*/
