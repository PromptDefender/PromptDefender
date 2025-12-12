module github.com/safetorun/PromptDefender/deployments/aws/lambda_keep

go 1.20

require (
	github.com/aws/aws-lambda-go v1.46.0
	github.com/safetorun/PromptDefender/aiprompt v0.0.0-20231210112259-15b98f65cf67
	github.com/safetorun/PromptDefender/cache v0.0.0-00010101000000-000000000000
	github.com/safetorun/PromptDefender/internal/base_aws v0.0.0-20240130072532-9a4bc83dc7d2
	github.com/safetorun/PromptDefender/keep v0.0.0-20231210112259-15b98f65cf67
	go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda v0.48.0
	go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig v0.48.0
	go.opentelemetry.io/contrib/propagators/aws v1.23.0
	go.opentelemetry.io/otel v1.23.1
)

replace (
	github.com/safetorun/PromptDefender/aiprompt => ../../internal/aiprompt
	github.com/safetorun/PromptDefender/cache => ../../internal/cache
	github.com/safetorun/PromptDefender/internal/base_aws => ../../internal/base_aws
	github.com/safetorun/PromptDefender/keep => ../../internal/keep
	github.com/safetorun/PromptDefender/request_log => ../../pkg/request_log
	github.com/safetorun/PromptDefender/utils => ../../internal/utils
)

require (
	github.com/aws/aws-sdk-go v1.55.8 // indirect
	github.com/cenkalti/backoff/v4 v4.2.1 // indirect
	github.com/go-logr/logr v1.4.1 // indirect
	github.com/go-logr/stdr v1.2.2 // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/grpc-ecosystem/grpc-gateway/v2 v2.19.0 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/safetorun/PromptDefender/request_log v0.0.0-00010101000000-000000000000 // indirect
	github.com/safetorun/PromptDefender/utils v0.0.0-00010101000000-000000000000 // indirect
	go.opentelemetry.io/contrib/detectors/aws/lambda v0.48.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace v1.23.0 // indirect
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc v1.23.0 // indirect
	go.opentelemetry.io/otel/metric v1.23.1 // indirect
	go.opentelemetry.io/otel/sdk v1.23.0 // indirect
	go.opentelemetry.io/otel/trace v1.23.1 // indirect
	go.opentelemetry.io/proto/otlp v1.1.0 // indirect
	golang.org/x/net v0.20.0 // indirect
	golang.org/x/sys v0.16.0 // indirect
	golang.org/x/text v0.14.0 // indirect
	google.golang.org/genproto/googleapis/api v0.0.0-20240102182953-50ed04b92917 // indirect
	google.golang.org/genproto/googleapis/rpc v0.0.0-20240102182953-50ed04b92917 // indirect
	google.golang.org/grpc v1.61.0 // indirect
	google.golang.org/protobuf v1.33.0 // indirect
)
