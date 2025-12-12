package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/safetorun/PromptDefender/keep"
	"github.com/safetorun/PromptDefender/tracer"
	"github.com/safetorun/PromptDefender/user_repository"
	"github.com/safetorun/PromptDefender/wall"
	"go.opentelemetry.io/otel"
)

type Server struct {
	KeepInstance *keep.Keep
	WallInstance *wall.Wall
	UserRepo     user_repository.UserRepository
	APIKey       string // In a real scenario, this might come from middleware context
	// For now we assume a single API key or extract it from request context
}

// Helper to get API key from context or request header.
// The generated middleware `ApiKeyAuth` puts scopes in context, but maybe we need to extract the key itself.
// The OpenAPI spec says `x-api-key` header.
func getAPIKey(r *http.Request) string {
	return r.Header.Get("x-api-key")
}

func (s *Server) BuildKeep(w http.ResponseWriter, r *http.Request) {
	var req KeepRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	randomiseXmlTag := false
	if req.RandomiseXmlTag != nil {
		randomiseXmlTag = *req.RandomiseXmlTag
	}

	answer, err := s.KeepInstance.BuildKeep(keep.StartingPrompt{
		Prompt:       req.Prompt,
		RandomiseTag: randomiseXmlTag,
	})

	if err != nil {
		if keep.IsPromptRequiredError(err) {
			http.Error(w, "prompt cannot be empty", http.StatusBadRequest)
			return
		}
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	resp := KeepResponse{
		ShieldedPrompt: answer.NewPrompt,
		XmlTag:         answer.Tag,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) ListUsers(w http.ResponseWriter, r *http.Request) {
	apiKey := getAPIKey(r)
	users, err := s.UserRepo.GetUsers(apiKey)
	if err != nil {
		log.Printf("Error getting users: %v", err)
		http.Error(w, "Error retrieving users", http.StatusInternalServerError)
		return
	}

	var result []User
	for _, u := range users {
		uid := u.UserOrSessionId
		result = append(result, User{UserId: &uid})
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(result)
}

func (s *Server) AddUser(w http.ResponseWriter, r *http.Request) {
	var req User
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.UserId == nil {
		http.Error(w, "userId is required", http.StatusBadRequest)
		return
	}

	apiKey := getAPIKey(r)
	err := s.UserRepo.CreateUser(user_repository.UserCore{
		UserOrSessionId: *req.UserId,
		ApiKeyId:        apiKey,
	})

	if err != nil {
		log.Printf("Error creating user: %v", err)
		http.Error(w, "Error creating user", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (s *Server) RemoveUser(w http.ResponseWriter, r *http.Request, userId string) {
	apiKey := getAPIKey(r)
	err := s.UserRepo.DeleteUser(userId, apiKey)
	if err != nil {
		log.Printf("Error deleting user: %v", err)
		http.Error(w, "Error deleting user", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (s *Server) GetUser(w http.ResponseWriter, r *http.Request, userId string) {
	apiKey := getAPIKey(r)
	user, err := s.UserRepo.GetUserByID(userId, apiKey)
	if err != nil {
		if err == user_repository.ErrUserIDNotFound {
			http.Error(w, "User not found", http.StatusNotFound)
			return
		}
		log.Printf("Error retrieving user: %v", err)
		http.Error(w, "Error retrieving user", http.StatusInternalServerError)
		return
	}

	resp := User{UserId: &user.UserOrSessionId}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) BuildWall(w http.ResponseWriter, r *http.Request) {
	var req WallRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	ctx := r.Context()
	wallTracer := otel.Tracer("wall")
	t := tracer.NewTracer(ctx, wallTracer)

	answer, err := s.WallInstance.CheckWall(wall.PromptToCheck{
		Prompt:           req.Prompt,
		ScanPii:          req.ScanPii,
		XmlTagToCheckFor: req.XmlTag,
		CheckForBadWords: req.CheckBadwords,
		FastCheck:        req.FastCheck,
	}, t)

	if err != nil {
		log.Printf("Error checking wall: %v", err)
		http.Error(w, "Error checking wall", http.StatusInternalServerError)
		return
	}

	// Check suspicious user and session using remote call (recursive call to itself if we used the same structure,
	// but here we might need to assume the same server handles it or skip if we consolidate).
	// The original code called `CheckSuspiciousUser`. That made an HTTP request to another lambda.
	// Since we are now a single server, we can check directly against the Repo!
	// However, the original `WallLambda` made a remote call.
	// Let's implement `CheckSuspiciousUser` to check against `s.UserRepo` directly.

	apiKey := getAPIKey(r)

	var suspiciousUser *bool
	if req.UserId != nil {
		isSuspicious, err := s.checkSuspiciousUserDirect(*req.UserId, apiKey)
		if err != nil {
			log.Printf("Error checking suspicious user: %v", err)
			// We continue? Original code returned error.
			http.Error(w, "Error checking suspicious user", http.StatusInternalServerError)
			return
		}
		suspiciousUser = isSuspicious
	}

	var suspiciousSession *bool
	if req.SessionId != nil {
		isSuspicious, err := s.checkSuspiciousUserDirect(*req.SessionId, apiKey)
		if err != nil {
			log.Printf("Error checking suspicious session: %v", err)
			http.Error(w, "Error checking suspicious session", http.StatusInternalServerError)
			return
		}
		suspiciousSession = isSuspicious
	}

	var piiDetected *bool
	if answer.PiiResult != nil {
		piiDetected = &answer.PiiResult.ContainsPii
	}

	var xmlEscaping *bool
	if answer.XmlScannerResult != nil {
		xmlEscaping = &answer.XmlScannerResult.ContainsXmlEscaping
	}

	jb := false
	potentialJailbreak := &jb
	if answer.InjectionDetected || answer.ContainsBadWords {
		jb = true
		potentialJailbreak = &jb
	}

	resp := WallResponse{
		ContainsPii:          piiDetected,
		PotentialJailbreak:   potentialJailbreak,
		PotentialXmlEscaping: xmlEscaping,
		SuspiciousUser:       suspiciousUser,
		SuspiciousSession:    suspiciousSession,
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

func (s *Server) checkSuspiciousUserDirect(userId string, apiKey string) (*bool, error) {
	// Check against UserRepo directly instead of HTTP call
	_, err := s.UserRepo.GetUserByID(userId, apiKey)
	if err != nil {
		if err == user_repository.ErrUserIDNotFound {
			b := false
			return &b, nil
		}
		return nil, err
	}
	b := true
	return &b, nil
}
