package helper

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/milvus-io/milvus-sdk-go/v2/client"
	"github.com/milvus-io/milvus-sdk-go/v2/entity"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func GenerateTimeUUID() string {
	now := time.Now()
	timeComponent := fmt.Sprintf("%04d%02d%02d%02d%02d%02d%09d",
		now.Year(), now.Month(), now.Day(),
		now.Hour(), now.Minute(), now.Second(),
		now.Nanosecond())
	if len(timeComponent) > 32 {
		timeComponent = timeComponent[:32]
	}
	for len(timeComponent) < 32 {
		timeComponent += "0"
	}
	return fmt.Sprintf("%s-%s-%s-%s-%s",
		timeComponent[0:8],
		timeComponent[8:12],
		timeComponent[12:16],
		timeComponent[16:20],
		timeComponent[20:32])
}

func NewDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return db, nil
}

func CreateMilvusClientAndCollection(collectionName, MilvusAddress, dim string) (client.Client, error) {
	// 1. Create a context with a timeout.
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// 2. Connect to Milvus.
	log.Printf("Connecting to Milvus at %s...", MilvusAddress)
	milvusClient, err := client.NewClient(
		ctx,
		client.Config{
			Address: MilvusAddress,
		},
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Milvus: %w", err)
	}
	log.Println("Successfully connected to Milvus.")

	// 3. Define the schema for the collection with the corrected field definitions.
	schema := &entity.Schema{
		CollectionName: collectionName,
		Description:    "Face embeddings collection",
		AutoID:         false, // We will provide our own primary keys
		Fields: []*entity.Field{
			{
				Name:       "object_id",
				DataType:   entity.FieldTypeVarChar,
				PrimaryKey: true, // CORRECT: Field name is PrimaryKey
				AutoID:     false,
				// CORRECT: MaxLength is a type parameter for VarChar
				TypeParams: map[string]string{
					entity.TypeParamMaxLength: "36", // UUID length, value must be a string
				},
			},
			{
				Name:     "vector",
				DataType: entity.FieldTypeFloatVector,
				TypeParams: map[string]string{
					entity.TypeParamDim: dim, // The dimension of your vectors
				},
			},
		},
	}

	// 4. Check if the collection already exists.
	has, err := milvusClient.HasCollection(ctx, collectionName)
	if err != nil {
		return nil, fmt.Errorf("failed to check for collection: %w", err)
	}

	if has {
		log.Printf("Collection '%s' already exists.", collectionName)
	} else {
		// 5. Create the collection if it does not exist.
		log.Printf("Collection '%s' does not exist. Creating...", collectionName)
		err = milvusClient.CreateCollection(ctx, schema, 1) // 1 shard
		if err != nil {
			return nil, fmt.Errorf("failed to create collection: %w", err)
		}
		log.Printf("Successfully created collection '%s'", collectionName)
	}

	// 6. Load the collection into memory for searching.
	log.Printf("Loading collection '%s' into memory...", collectionName)
	err = milvusClient.LoadCollection(ctx, collectionName, false) // false means synchronous loading
	if err != nil {
		return nil, fmt.Errorf("failed to load collection: %w", err)
	}
	log.Printf("Collection '%s' loaded successfully.", collectionName)

	// 7. Return the ready-to-use client.
	return milvusClient, nil
}
