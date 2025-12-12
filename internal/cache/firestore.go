package cache

import (
	"context"

	"cloud.google.com/go/firestore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type FirestoreCache struct {
	Client         *firestore.Client
	CollectionName string
}

func NewFirestore(projectID string, collectionName string) (Cache, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &FirestoreCache{
		Client:         client,
		CollectionName: collectionName,
	}, nil
}

func (c *FirestoreCache) Get(key string) (*string, error) {
	ctx := context.Background()
	doc, err := c.Client.Collection(c.CollectionName).Doc(key).Get(ctx)
	if err != nil {
		if status.Code(err) == codes.NotFound {
			return nil, nil // Treat not found as cache miss (nil, nil) matching Ddb behavior
		}
		return nil, err
	}

	if !doc.Exists() {
		return nil, nil
	}

	// Use dataTo map
	data := doc.Data()
	if val, ok := data["Value"].(string); ok {
		return &val, nil
	}

	return nil, nil
}

func (c *FirestoreCache) Set(key string, value string) error {
	ctx := context.Background()
	_, err := c.Client.Collection(c.CollectionName).Doc(key).Set(ctx, map[string]interface{}{
		"Value": value,
	})
	return err
}
