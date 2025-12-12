module github.com/safetorun/PromptDefender/cmd/server

go 1.23

require (
	github.com/getkin/kin-openapi v0.118.0
	github.com/go-chi/chi/v5 v5.0.10
	github.com/oapi-codegen/runtime v1.1.1
	github.com/safetorun/PromptDefender/aiprompt v0.0.0-20231204191015-931fb010b80b
	github.com/safetorun/PromptDefender/badwords v0.0.0-20231204191015-931fb010b80b
	github.com/safetorun/PromptDefender/badwords_embeddings v0.0.0-00010101000000-000000000000
	github.com/safetorun/PromptDefender/cache v0.0.0-00010101000000-000000000000
	github.com/safetorun/PromptDefender/embeddings v0.0.0-20231204191015-931fb010b80b
	github.com/safetorun/PromptDefender/keep v0.0.0-00010101000000-000000000000
	github.com/safetorun/PromptDefender/pii_aws v0.0.0-00010101000000-000000000000
	github.com/safetorun/PromptDefender/remote_sagemaker_call v0.0.0-00010101000000-000000000000
	github.com/safetorun/PromptDefender/tracer v0.0.0-00010101000000-000000000000
	github.com/safetorun/PromptDefender/user_repository v0.0.0-00010101000000-000000000000
	github.com/safetorun/PromptDefender/user_repository_ddb v0.0.0-00010101000000-000000000000
	github.com/safetorun/PromptDefender/wall v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.23.0
)

require (
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/aws/aws-sdk-go v1.51.0 // indirect
	github.com/drewlanenga/govector v0.0.0-20220726163947-b958ac08bc93 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.5 // indirect
	github.com/google/uuid v1.5.0 // indirect
	github.com/invopop/yaml v0.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/perimeterx/marshmallow v1.1.4 // indirect
	github.com/rogpeppe/go-internal v1.14.1 // indirect
	github.com/safetorun/PromptDefender/pii v0.0.0-20231204191015-931fb010b80b // indirect
	github.com/safetorun/PromptDefender/utils v0.0.0-00010101000000-000000000000 // indirect
	github.com/sashabaranov/go-openai v1.17.9 // indirect
	go.opentelemetry.io/otel/metric v1.23.0 // indirect
	go.opentelemetry.io/otel/trace v1.23.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)

replace (
	github.com/safetorun/PromptDefender/aiprompt => ../../internal/aiprompt
	github.com/safetorun/PromptDefender/badwords => ../../internal/badwords
	github.com/safetorun/PromptDefender/badwords_embeddings => ../../internal/badwords_embeddings
	github.com/safetorun/PromptDefender/cache => ../../internal/cache
	github.com/safetorun/PromptDefender/embeddings => ../../internal/embeddings
	github.com/safetorun/PromptDefender/internal/base_aws => ../../internal/base_aws
	github.com/safetorun/PromptDefender/keep => ../../internal/keep
	github.com/safetorun/PromptDefender/pii => ../../internal/pii
	github.com/safetorun/PromptDefender/pii_aws => ../../internal/pii_aws
	github.com/safetorun/PromptDefender/remote_sagemaker_call => ../../internal/sagemaker_jailbreak_model
	github.com/safetorun/PromptDefender/request_log => ../../pkg/request_log
	github.com/safetorun/PromptDefender/tracer => ../../internal/tracer
	github.com/safetorun/PromptDefender/user_repository => ../../internal/user_repository
	github.com/safetorun/PromptDefender/user_repository_ddb => ../../internal/user_repository_ddb
	github.com/safetorun/PromptDefender/utils => ../../internal/utils
	github.com/safetorun/PromptDefender/wall => ../../internal/wall
)
