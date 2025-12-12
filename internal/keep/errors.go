package keep

type PromptRequiredError struct {
}

func (e *PromptRequiredError) Error() string {
	return "prompt cannot be empty"
}

func NewPromptRequiredError() *PromptRequiredError {
	return &PromptRequiredError{}
}

func IsPromptRequiredError(err error) bool {
	_, ok := err.(*PromptRequiredError)
	return ok
}
