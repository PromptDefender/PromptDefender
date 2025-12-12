module github.com/safetorun/PromptDefender/cmd/server

go 1.25.3

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
	github.com/safetorun/PromptDefender/pii_google v0.0.0
	github.com/safetorun/PromptDefender/tracer v0.0.0-00010101000000-000000000000
	github.com/safetorun/PromptDefender/user_repository v0.0.0-00010101000000-000000000000
	github.com/safetorun/PromptDefender/user_repository_firestore v0.0.0
	github.com/safetorun/PromptDefender/wall v0.0.0-00010101000000-000000000000
	go.opentelemetry.io/otel v1.38.0
)

require (
	cloud.google.com/go v0.121.6 // indirect
	cloud.google.com/go/aiplatform v1.90.0 // indirect
	cloud.google.com/go/auth v0.16.4 // indirect
	cloud.google.com/go/auth/oauth2adapt v0.2.8 // indirect
	cloud.google.com/go/compute/metadata v0.9.0 // indirect
	cloud.google.com/go/dlp v1.27.0 // indirect
	cloud.google.com/go/firestore v1.20.0 // indirect
	cloud.google.com/go/iam v1.5.2 // indirect
	cloud.google.com/go/longrunning v0.6.7 // indirect
	cloud.google.com/go/vertexai v0.15.0 // indirect
	github.com/apapsch/go-jsonmerge/v2 v2.0.0 // indirect
	github.com/aws/aws-sdk-go v1.51.0 // indirect
	github.com/drewlanenga/govector v0.0.0-20220726163947-b958ac08bc93 // indirect
	github.com/felixge/httpsnoop v1.0.4 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.5 // indirect
	github.com/go-openapi/swag v0.19.5 // indirect
	github.com/google/s2a-go v0.1.9 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/googleapis/enterprise-certificate-proxy v0.3.6 // indirect
	github.com/googleapis/gax-go/v2 v2.15.0 // indirect
	github.com/invopop/yaml v0.1.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/mohae/deepcopy v0.0.0-20170929034955-c48cc78d4826 // indirect
	github.com/perimeterx/marshmallow v1.1.4 // indirect
	github.com/safetorun/PromptDefender/pii v0.0.0 // indirect
	github.com/safetorun/PromptDefender/utils v0.0.0-00010101000000-000000000000 // indirect
	github.com/sashabaranov/go-openai v1.17.9 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.61.0 // indirect
	go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp v0.61.0 // indirect
	go.opentelemetry.io/otel/metric v1.38.0 // indirect
	go.opentelemetry.io/otel/trace v1.38.0 // indirect
	golang.org/x/crypto v0.43.0 // indirect
	golang.org/x/net v0.46.1-0.20251013234738-63d1a5100f82 // indirect
	golang.org/x/oauth2 v0.32.0 // indirect
	golang.org/x/sync v0.17.0 // indirect
	golang.org/x/sys v0.37.0 // indirect
	golang.org/x/text v0.30.0 // indirect
	golang.org/x/time v0.12.0 // indirect
	google.golang.org/api v0.247.0 // indirect
	google.golang.org/genproto v0.0.0-20250603155806-513f23925822 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20251022142026-3a174f9686a8 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20251022142026-3a174f9686a8 // indirect
	google.golang.org/grpc v1.77.0 // indirect
	google.golang.org/protobuf v1.36.10 // indirect
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
	github.com/safetorun/PromptDefender/pii_google => ../../internal/pii_google
	github.com/safetorun/PromptDefender/request_log => ../../pkg/request_log
	github.com/safetorun/PromptDefender/tracer => ../../internal/tracer
	github.com/safetorun/PromptDefender/user_repository => ../../internal/user_repository
	github.com/safetorun/PromptDefender/user_repository_firestore => ../../internal/user_repository_firestore
	github.com/safetorun/PromptDefender/utils => ../../internal/utils
	github.com/safetorun/PromptDefender/wall => ../../internal/wall
)
