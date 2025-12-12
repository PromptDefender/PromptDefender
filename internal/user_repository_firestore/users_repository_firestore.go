package user_repository_firestore

import (
	"context"

	"cloud.google.com/go/firestore"
	"github.com/safetorun/PromptDefender/user_repository"
	"google.golang.org/api/iterator"
)

const IndexKey = "ApiKeyId"

type UserRepositoryFirestore struct {
	Client         *firestore.Client
	CollectionName string
}

func New(projectID string, collectionName string) (*UserRepositoryFirestore, error) {
	ctx := context.Background()
	client, err := firestore.NewClient(ctx, projectID)
	if err != nil {
		return nil, err
	}
	return &UserRepositoryFirestore{
		Client:         client,
		CollectionName: collectionName,
	}, nil
}

func (r *UserRepositoryFirestore) GetUserByID(id string, apiKeyId string) (*user_repository.UserCore, error) {
	ctx := context.Background()
	// Assuming ID is unique, direct lookup
	// Note: DynamoDB implementation uses composite key or specific schema.
	// If UserOrSessionId is the document ID, we can do direct lookup.
	// If not, we might need a query.
	// DdB implementation: Key: UserOrSessionId, ApiKeyId.
	// In Firestore, if we use UserOrSessionId as doc ID, we can verify ApiKeyId in fields.
	// Or we can query where UserOrSessionId == id and ApiKeyId == apiKeyId.

	// Let's implement as Query for safety unless we decide on schema.
	// Querying is safer to match DdB "GetItem" with specific keys.

	iter := r.Client.Collection(r.CollectionName).
		Where("UserOrSessionId", "==", id).
		Where("ApiKeyId", "==", apiKeyId).
		Documents(ctx)

	doc, err := iter.Next()
	if err == iterator.Done {
		return nil, user_repository.ErrUserIDNotFound
	}
	if err != nil {
		return nil, err
	}

	var user user_repository.UserCore
	if err := doc.DataTo(&user); err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *UserRepositoryFirestore) GetUsers(apiKeyId string) ([]user_repository.UserCore, error) {
	ctx := context.Background()
	iter := r.Client.Collection(r.CollectionName).Where(IndexKey, "==", apiKeyId).Documents(ctx)

	var users []user_repository.UserCore
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return nil, err
		}
		var user user_repository.UserCore
		if err := doc.DataTo(&user); err != nil {
			return nil, err
		}
		users = append(users, user)
	}
	return users, nil
}

func (r *UserRepositoryFirestore) CreateUser(user user_repository.UserCore) error {
	ctx := context.Background()
	// We need a document ID. DdB uses UserOrSessionId usually, or we can let Firestore auto-generate.
	// But uniqueness of UserOrSessionId + ApiKeyId might be important.
	// If we use UserOrSessionId as doc ID, it enforces one user per ID global (or per collection).
	// Let's use auto-id or UserOrSessionId if it's unique enough.
	// Given DdB usage, UserOrSessionId seems to be the primary key part.
	// Let's assume UserOrSessionId is unique enough to use as DocRef if we want upsert, OR just Add().

	// DdB PutItem overwrites.
	// If we want to overwrite based on UserOrSessionId:
	// _, err := r.Client.Collection(r.CollectionName).Doc(user.UserOrSessionId).Set(ctx, user)
	// But DdB has composite key (UserOrSessionId, ApiKeyId). Firestore doc IDs are single strings.
	// We can make a composite ID "UserOrSessionId:ApiKeyId" or just Add() and rely on queries.
	// However, Update/Delete needs to locate it. DeleteUser takes both IDs.

	// Let's use Query to find and update/delete if we don't enforce a specific DocID scheme.
	// ACTUALLY, simpler: just create a new document if it doesn't exist?
	// DdB PutItem overwrites.

	// Best practice for Firestore replacement of DdB composite key:
	// Use a deterministic DocID = "UserOrSessionId:ApiKeyId" (sanitized) to enable O(1) Get/Delete.
	// Let's try that.

	docID := user.UserOrSessionId + ":" + user.ApiKeyId
	_, err := r.Client.Collection(r.CollectionName).Doc(docID).Set(ctx, user)
	return err
}

func (r *UserRepositoryFirestore) DeleteUser(userOrSessionId string, apiKeyId string) error {
	ctx := context.Background()
	docID := userOrSessionId + ":" + apiKeyId
	_, err := r.Client.Collection(r.CollectionName).Doc(docID).Delete(ctx)
	return err
}
