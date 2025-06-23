package service

import (
	"context"
	"errors"
	"fmt"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"github.com/ruziba3vich/face_embedding_gateway/internal/models"
	"gorm.io/gorm"
)

type (
	Service struct {
		milvusClient client.Client
		db           *gorm.DB
	}
)

func NewService(db *gorm.DB, milvusClient client.Client) *Service {
	return &Service{
		db:           db,
		milvusClient: milvusClient,
	}
}

func (s *Service) Create(user *models.User) error {
	return s.db.Create(user).Error
}

func (s *Service) GetByID(id string) (*models.User, error) {
	var user models.User
	err := s.db.First(&user, "id = ?", id).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, nil
	}
	return &user, err
}

func (s *Service) DeleteByID(id string) error {
	return s.db.Delete(&models.User{}, "id = ?", id).Error
}

func (h *Service) StoreVector(ctx context.Context, objectID string, vector []float32) error {
	if objectID == "" || len(vector) == 0 {
		return errors.New("invalid input: objectID or vector is empty")
	}

	vectorColumn := entity.NewColumnFloatVector("vector", len(vector), [][]float32{vector})
	idColumn := entity.NewColumnVarChar("object_id", []string{objectID})

	_, err := h.milvusClient.Insert(ctx, "face_embeddings", "", idColumn, vectorColumn)
	if err != nil {
		return fmt.Errorf("failed to store vector: %v", err)
	}

	return nil
}
