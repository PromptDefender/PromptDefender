package aiprompt

type RemoteAIChecker interface {
	CheckAI(prompt string) (*string, error)
}
