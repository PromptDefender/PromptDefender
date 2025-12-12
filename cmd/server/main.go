package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/safetorun/PromptDefender/aiprompt"
	"github.com/safetorun/PromptDefender/badwords"
	"github.com/safetorun/PromptDefender/badwords_embeddings"
	"github.com/safetorun/PromptDefender/cache"
	"github.com/safetorun/PromptDefender/embeddings"
	"github.com/safetorun/PromptDefender/huggingface_jailbreak_model"
	"github.com/safetorun/PromptDefender/keep"
	"github.com/safetorun/PromptDefender/pii_google"
	"github.com/safetorun/PromptDefender/user_repository_firestore"
	"github.com/safetorun/PromptDefender/wall"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	// GCP Configuration
	projectID := os.Getenv("GOOGLE_PROJECT_ID")
	if projectID == "" {
		log.Fatal("GOOGLE_PROJECT_ID environment variable is required")
	}
	location := os.Getenv("GOOGLE_LOCATION")
	if location == "" {
		location = "us-central1"
	}

	// AI Configuration
	vertexModel := os.Getenv("VERTEX_AI_MODEL") // e.g., gemini-1.5-flash
	if vertexModel == "" {
		vertexModel = "gemini-1.5-flash"
	}

	// Cache Configuration
	cacheCollectionName := os.Getenv("CACHE_COLLECTION_NAME")
	if cacheCollectionName == "" {
		cacheCollectionName = "cache"
	}

	// Initialize Keep
	var keepOpts keep.KeepOption = func(c *keep.Keep) {
		// Firestore Cache
		fc, err := cache.NewFirestore(projectID, cacheCollectionName)
		if err != nil {
			log.Printf("Warning: Failed to initialize Firestore cache: %v", err)
		} else {
			c.Cache = &fc
		}
	}

	vertexAI, err := aiprompt.NewVertexAI(projectID, location, vertexModel)
	if err != nil {
		log.Fatalf("Error converting Vertex AI client: %v", err)
	}
	defer vertexAI.Close()

	keepInstance := keep.New(vertexAI, keepOpts)

	// Initialize Wall
	wallOpts := func(c *wall.Wall) error {
		// PII Scanner (Google DLP)
		piiScanner, err := pii_google.New(projectID)
		if err != nil {
			return fmt.Errorf("failed to init PII scanner: %w", err)
		}
		c.PiiScanner = piiScanner

		// Badwords (Vertex Embeddings)
		// Note: Badwords embeddings check relies on creating embeddings for input.
		embeddingModel := "text-embedding-004"
		vertexEmbeddings, err := embeddings.NewVertex(projectID, location, embeddingModel)
		if err != nil {
			return fmt.Errorf("failed to init Vertex embeddings: %w", err)
		}

		c.BadWordsCheck = badwords.New(badwords_embeddings.New(vertexEmbeddings))

		// XML Escaping
		c.XmlEscapingScanner = wall.NewBasicXmlEscapingScaner()

		// Remote Jailbreak Check (HuggingFace)
		huggingFaceToken := os.Getenv("HUGGINGFACE_TOKEN")
		if huggingFaceToken != "" {
			hfCaller := huggingface_jailbreak_model.NewRemoteApiCaller(huggingFaceToken)
			c.RemoteApiCaller = &hfCaller
		} else {
			log.Println("Warning: HUGGINGFACE_TOKEN not set, skipping remote jailbreak check")
		}

		if cacheCollectionName != "" {
			fc, err := cache.NewFirestore(projectID, cacheCollectionName)
			if err == nil {
				c.Cache = &fc
			}
		}
		return nil
	}

	wallInstance, err := wall.New(wallOpts)
	if err != nil {
		log.Fatalf("Error creating wall instance: %v", err)
	}

	// Initialize User Repo (Firestore)
	userCollection := os.Getenv("USERS_COLLECTION")
	if userCollection == "" {
		userCollection = "users"
	}
	userRepo, err := user_repository_firestore.New(projectID, userCollection)
	if err != nil {
		log.Fatalf("Error creating user repository: %v", err)
	}

	// Initialize Server
	server := &Server{
		KeepInstance: keepInstance,
		WallInstance: wallInstance,
		UserRepo:     userRepo,
	}

	// Router setup
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// We use HandlerFromMux to register the routes from api.gen.go to our router
	HandlerFromMux(server, r)

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
