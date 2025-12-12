# PromptDefender API Server

This is the single API server implementation for PromptDefender, replacing the individual Lambdas.

## Prerequisites

- Go 1.23+
- AWS Credentials configured (for User Repository DynamoDB and PII Scanner)

## Environment Variables

| Variable | Description | Required |
| --- | --- | --- |
| `open_ai_api_key` | OpenAI API Key for embeddings and keep generations | Yes |
| `SAGEMAKER_ENDPOINT_JAILBREAK` | AWS Sagemaker Endpoint for Jailbreak detection | No (Jailbreak detection disabled if missing) |
| `CACHE_TABLE_NAME` | DynamoDB Table name for caching results | No (Caching disabled if missing) |
| `PORT` | HTTP Port to listen on (default: 8080) | No |
| `AWS_REGION` | AWS Region (usually required by SDK if not in profile) | Yes |

## Running

```bash
cd cmd/server
go build
export open_ai_api_key="sk-..."
./server
```

## API

The server implements the API defined in `api/openapi.yml`.
