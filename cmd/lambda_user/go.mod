module github.com/safetorun/PromptDefender/deployments/aws/lambda_users

go 1.24.0

require (
	github.com/aws/aws-lambda-go v1.46.0
	github.com/safetorun/PromptDefender/internal/base_aws v0.0.0-20240130072532-9a4bc83dc7d2
	github.com/stretchr/testify v1.11.1
	github.com/thoas/go-funk v0.9.3
	go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda v0.48.0
	go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig v0.48.0
	go.opentelemetry.io/contrib/propagators/aws v1.23.0
	go.opentelemetry.io/otel v1.40.0
)

require (
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/go-logr/logr v1.4.3 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.0 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	go.opentelemetry.io/auto/sdk v1.2.1 // indirect
	go.opentelemetry.io/contrib/detectors/aws/lambda v0.48.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.23.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.23.0 // indirect
	go.opentelemetry.io/otel/metric v1.40.0 // indirect
	go.opentelemetry.io/otel/sdk v1.40.0 // indirect
	go.opentelemetry.io/otel/trace v1.40.0 // indirect
	go.opentelemetry.io/proto/otlp v1.1.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.40.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240102182953-50ed04b92917 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240102182953-50ed04b92917 // indirect
	google.golang.org/grpc v1.61.0 // indirect
	google.golang.org/protobuf v1.32.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
