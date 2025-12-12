package main

import (
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
	"github.com/safetorun/PromptDefender/keep"
	"github.com/safetorun/PromptDefender/pii_aws"
	sagemaker_jailbreak_model "github.com/safetorun/PromptDefender/remote_sagemaker_call"
	"github.com/safetorun/PromptDefender/user_repository_ddb"
	"github.com/safetorun/PromptDefender/wall"
)

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	openAIKey := os.Getenv("open_ai_api_key")
	if openAIKey == "" {
		log.Fatal("open_ai_api_key environment variable is required")
	}

	cacheTableName := os.Getenv("CACHE_TABLE_NAME")
	// It's possible cache table name is optional/empty in some configs, but lambda code seemingly expected it.
	// We'll proceed.

	// Initialize Keep
	var keepOpts keep.KeepOption = func(c *keep.Keep) {
		if cacheTableName != "" {
			ddbCache := cache.New(cacheTableName)
			c.Cache = &ddbCache
		}
	}
	keepInstance := keep.New(aiprompt.NewOpenAI(openAIKey), keepOpts)

	// Initialize Wall
	wallOpts := func(c *wall.Wall) error {
		c.PiiScanner = pii_aws.New() // Accesses AWS Comprehend via pii_aws
		c.BadWordsCheck = badwords.New(badwords_embeddings.New(embeddings.New(openAIKey)))
		c.XmlEscapingScanner = wall.NewBasicXmlEscapingScaner()

		sagemakerEndpoint := os.Getenv("SAGEMAKER_ENDPOINT_JAILBREAK")
		if sagemakerEndpoint != "" {
			apiCaller := sagemaker_jailbreak_model.New(sagemakerEndpoint)
			c.RemoteApiCaller = &apiCaller
		} else {
			log.Println("SAGEMAKER_ENDPOINT_JAILBREAK not set, injection detection via Sagemaker disabled")
		}

		if cacheTableName != "" {
			ddbCache := cache.New(cacheTableName)
			c.Cache = &ddbCache
		}
		return nil
	}

	wallInstance, err := wall.New(wallOpts)
	if err != nil {
		log.Fatalf("Error creating wall instance: %v", err)
	}

	// Initialize User Repo
	userRepo := user_repository_ddb.New() // Uses AWS config from env

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
	// This creates the handler that routes to ServerInterface methods
	HandlerFromMux(server, r)

	log.Printf("Starting server on port %s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
