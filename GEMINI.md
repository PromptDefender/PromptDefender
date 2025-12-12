# PromptDefender

PromptDefender (also known as PromptShield or SafeToRun) is a multi-layer defense system designed to protect LLM applications against prompt injection attacks, data leakage (PII), and other security threats.

## Architecture

The system consists of several core components:

*   **Wall**: The first line of defense. It sanitizes input by checking for:
    *   Jailbreak attempts (using heuristics and ML models).
    *   PII (Personally Identifiable Information).
    *   Bad words/profanity.
    *   XML escaping attacks.
    *   Suspicious user activity.
*   **Keep**: Wraps the user's prompt with "instruction defense" to prevent the LLM from being manipulated. It manages the construction of safe prompts.
*   **User Repository**: Manages a list of suspicious users or sessions to block attackers.
*   **Drawbridge (Planned)**: Validates LLM responses for canaries and unsafe content.

## Tech Stack

*   **Language**: Go (Golang)
*   **API Framework**: `go-chi` (HTTP router)
*   **Infrastructure**: Terraform (AWS Lambda, API Gateway, DynamoDB)
*   **API Spec**: OpenAPI 3.0 (`api/openapi.yml`)
*   **Code Generation**: `oapi-codegen` for generating Go server and client code from OpenAPI specs.

## Development Workflow

### Prerequisites
*   Go 1.20+
*   Terraform
*   `make`
*   Python 3 (for jailbreak data processing scripts)

### Key Commands

The project uses a `Makefile` to manage common tasks:

*   **Generate Code**: `make generate`
    *   Regenerates the Go API server code from `api/openapi.yml`.
    *   Run this after modifying the OpenAPI spec.
*   **Build**: `make build`
    *   Compiles the server binary in `cmd/server`.
*   **Test**: `make test`
    *   Runs unit tests across all modules.
*   **Integration Test**: `make integration_test`
    *   Runs integration tests against a deployed environment. Requires `URL` and `DEFENDER_API_KEY` env vars.
*   **Deploy**: `make deploy`
    *   Deploys the infrastructure using Terraform.

### Running Locally

To run the server locally, you need to set specific environment variables and run the entry point in `cmd/server`.

1.  **Set Environment Variables**:
    *   `PORT`: Server port (default: 8080)
    *   `open_ai_api_key`: Required for embeddings and bad word checks.
    *   `CACHE_TABLE_NAME`: (Optional) DynamoDB table name for caching.
    *   `SAGEMAKER_ENDPOINT_JAILBREAK`: (Optional) Endpoint for the SageMaker jailbreak detection model.

2.  **Run**:
    ```bash
    go run cmd/server/main.go
    ```

## Key Directories

*   `cmd/server/`: Main application entry point (`main.go`) and generated API server code.
*   `internal/`: Core application logic.
    *   `wall/`: Logic for the "Wall" defense layer.
    *   `keep/`: Logic for the "Keep" defense layer.
    *   `aiprompt/`: Interactions with OpenAI.
    *   `badwords/`: Bad word detection.
    *   `pii/` & `pii_aws/`: PII scanning logic.
*   `api/`: OpenAPI specifications (`openapi.yml`).
*   `terraform/`: Infrastructure as Code definitions for AWS deployment.
*   `test/`: Integration tests and load tests (`k6`).
*   `builder/`: Python scripts for processing jailbreak datasets and building models.

## Deployment

Deployment is handled via Terraform.
*   `terraform-base-infrastructure/`: Contains base resources (e.g., HuggingFace datasets).
*   `terraform/`: Contains the main application infrastructure (Lambda functions, API Gateway, etc.).
