package keep

import (
	"encoding/json"
	"fmt"
	"github.com/safetorun/PromptDefender/aiprompt"
	"github.com/safetorun/PromptDefender/cache"
	"github.com/safetorun/PromptDefender/utils"
	"log"
	"math/rand"
	"time"
)

type KeepOption func(*Keep)

type Keep struct {
	openAi aiprompt.RemoteAIChecker
	Cache  *cache.Cache
	Logger *log.Logger
}

type StartingPrompt struct {
	Prompt       string
	RandomiseTag bool
}

type NewPrompt struct {
	NewPrompt string
	Tag       string
}

func New(aiPrompt aiprompt.RemoteAIChecker, options ...KeepOption) *Keep {
	k := &Keep{
		openAi: aiPrompt,
		Logger: log.Default(),
	}

	for _, opt := range options {
		opt(k)
	}

	return k
}

func generateRandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var seededRand *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	b := make([]byte, length)
	for i := range b {
		b[i] = charset[seededRand.Intn(len(charset))]
	}
	return string(b)
}

func (k *Keep) BuildKeep(startingPrompt StartingPrompt) (*NewPrompt, error) {

	tag := "user_input"

	if startingPrompt.Prompt == "" {
		return nil, NewPromptRequiredError()
	}

	if k.Cache != nil {
		cached, cachedResult, err := checkCache(k.Cache, startingPrompt)
		if err != nil {
			k.Logger.Println(fmt.Sprintf("Error checking cache: %v", err))
			return nil, err
		}

		if cached {
			k.Logger.Println("Cache hit")
			return cachedResult, nil
		} else {
			k.Logger.Println("Cache miss")
		}
	}

	if startingPrompt.RandomiseTag {
		tag = generateRandomString(10)
	}

	builtPrompt := HardenedPrompt(SmartPromptRequest{BasePrompt: startingPrompt.Prompt, XmlTagName: tag})

	response, err := k.openAi.CheckAI(builtPrompt)

	if err != nil {
		return nil, err
	}

	newPrompt := NewPrompt{NewPrompt: *response, Tag: tag}

	err = storeCache(k.Cache, startingPrompt, &newPrompt)

	if err != nil {
		k.Logger.Printf("Error storing cache: %v\n", err)
	}

	return &newPrompt, nil
}

func checkCache(cache *cache.Cache, prompt StartingPrompt) (bool, *NewPrompt, error) {
	if cache == nil {
		println("Cache is nil")
		return false, nil, nil
	}

	b, err := json.Marshal(prompt)

	if err != nil {
		println("Error marshalling cache: ", err)
		return false, nil, err
	}

	cachedResult, err := (*cache).Get(utils.HashString(string(b)))

	if err != nil {
		println(fmt.Sprintf("Error getting cache: %v", err))
		return false, nil, err
	}

	if cachedResult != nil {
		var cachedResultReturn *NewPrompt
		err := json.Unmarshal([]byte(*cachedResult), &cachedResultReturn)

		if err != nil {
			println("Error unmarshalling cache: ", err)
			return false, nil, err
		}

		return true, cachedResultReturn, nil
	}

	return false, nil, nil
}

func storeCache(cache *cache.Cache, prompt StartingPrompt, result *NewPrompt) error {
	if cache == nil {
		return nil
	}

	b, err := json.Marshal(prompt)

	if err != nil {
		return err
	}

	bResult, err := json.Marshal(result)

	if err != nil {
		return err
	}

	return (*cache).Set(utils.HashString(string(b)), string(bResult))
}
