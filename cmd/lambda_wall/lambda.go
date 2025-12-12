package main

import (
	"context"
	"fmt"
	"github.com/safetorun/PromptDefender/cache"
	sagemaker_jailbreak_model "github.com/safetorun/PromptDefender/remote_sagemaker_call"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/safetorun/PromptDefender/badwords"
	"github.com/safetorun/PromptDefender/badwords_embeddings"
	"github.com/safetorun/PromptDefender/embeddings"
	"github.com/safetorun/PromptDefender/internal/base_aws"
	"github.com/safetorun/PromptDefender/pii_aws"
	"github.com/safetorun/PromptDefender/tracer"
	"github.com/safetorun/PromptDefender/wall"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda"
	"go.opentelemetry.io/contrib/instrumentation/github.com/aws/aws-lambda-go/otellambda/xrayconfig"
	"go.opentelemetry.io/contrib/propagators/aws/xray"
	"go.opentelemetry.io/otel"
)

var (
	wallTracer = otel.Tracer("wall")
)

type WallLambda struct {
	wallInstance *wall.Wall
	context      context.Context
	apiKey       string
	url          string
}

func (m *WallLambda) Handle(wallRequest WallRequest) (*WallResponse, error) {

	t := tracer.NewTracer(m.context, wallTracer)

	answer, err := m.wallInstance.CheckWall(wall.PromptToCheck{
		Prompt:           wallRequest.Prompt,
		ScanPii:          wallRequest.ScanPii,
		XmlTagToCheckFor: wallRequest.XmlTag,
		CheckForBadWords: wallRequest.CheckBadwords,
		FastCheck:        wallRequest.FastCheck,
	},
		t,
	)

	_, span := wallTracer.Start(m.context, "check_suspicious_user")
	suspiciousUser, err := CheckSuspiciousUser(m.context, m.url, wallRequest.UserId, m.apiKey)

	if err != nil {
		return nil, err
	}

	span.End()

	_, span = wallTracer.Start(m.context, "check_suspicious_session")
	suspiciousSession, err := CheckSuspiciousUser(m.context, m.url, wallRequest.SessionId, m.apiKey)

	if err != nil {
		return nil, err
	}

	span.End()

	if err != nil {
		return nil, err
	}

	var piiDetected *bool = nil

	if answer.PiiResult != nil {
		piiDetected = &answer.PiiResult.ContainsPii
	}

	var xmlEscaping *bool = nil

	if answer.XmlScannerResult != nil {
		xmlEscaping = &answer.XmlScannerResult.ContainsXmlEscaping
	}

	jb := false
	var potentialJailbreak = &jb

	if answer.InjectionDetected || answer.ContainsBadWords {
		jb := true
		potentialJailbreak = &jb
	}

	return &WallResponse{
		ContainsPii:          piiDetected,
		PotentialJailbreak:   potentialJailbreak,
		PotentialXmlEscaping: xmlEscaping,
		SuspiciousUser:       suspiciousUser,
		SuspiciousSession:    suspiciousSession,
	}, nil
}

func CheckSuspiciousUser(ctx context.Context, url string, userId *string, apiKey string) (*bool, error) {
	if userId == nil {
		return nil, nil
	}

	_, span := wallTracer.Start(ctx, "check_suspicious_user")
	defer span.End()

	suspiciousUser := NewRemote(url, apiKey)

	log.Default().Println("Checking suspicious user with id: ", *userId)
	isSuspicous, err := suspiciousUser.CheckSuspiciousUser(ctx, *userId)

	if err != nil {
		log.Default().Println("Error checking suspicious user: ", err)
		return nil, err
	}

	return isSuspicous, nil
}

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	ctx, span := wallTracer.Start(ctx, "wall_setup")

	openAIKey, exists := os.LookupEnv("open_ai_api_key")

	log.Default().Println("Received request for wall lambda")

	if !exists {
		return events.APIGatewayProxyResponse{StatusCode: 400}, fmt.Errorf("error with configuration")
	}

	addAllConfigurations := func(c *wall.Wall) error {
		c.PiiScanner = pii_aws.New()
		c.BadWordsCheck = badwords.New(badwords_embeddings.New(embeddings.New(openAIKey)))
		c.XmlEscapingScanner = wall.NewBasicXmlEscapingScaner()
		var apiCaller = sagemaker_jailbreak_model.New(os.Getenv("SAGEMAKER_ENDPOINT_JAILBREAK"))
		c.RemoteApiCaller = &apiCaller
		ddbCache := cache.New(os.Getenv("CACHE_TABLE_NAME"))
		c.Cache = &ddbCache
		return nil
	}

	wallInstance, err := wall.New(addAllConfigurations)

	url := retrieveUrl(request)
	log.Default().Println("Received request for wall lambda with url: ", url)

	wallLambda := WallLambda{
		wallInstance: wallInstance,
		context:      ctx,
		apiKey:       request.RequestContext.Identity.APIKey,
		url:          url,
	}

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

	span.End()
	ctx, span = wallTracer.Start(ctx, "wall_handler_exec")
	defer span.End()

	response, err := base_aws.BaseHandler[WallRequest, WallResponse](request, &wallLambda)

	if err != nil {
		return events.APIGatewayProxyResponse{StatusCode: 400}, err
	}

	return response, nil
}

func retrieveUrl(request events.APIGatewayProxyRequest) string {
	scheme := request.Headers["X-Forwarded-Proto"]
	host := request.Headers["Host"]

	if !strings.Contains("safetorun.com", host) {
		host = host + "/" + request.RequestContext.Stage
	}

	url := fmt.Sprintf("%s://%s", scheme, host)
	return url
}

func main() {

	ctx := context.Background()

	tp, err := xrayconfig.NewTracerProvider(ctx)
	if err != nil {
		fmt.Printf("error creating tracer provider: %v", err)
	}

	defer func(ctx context.Context) {
		err := tp.Shutdown(ctx)
		if err != nil {
			fmt.Printf("error shutting down tracer provider: %v", err)
		}
	}(ctx)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(xray.Propagator{})

	lambda.Start(otellambda.InstrumentHandler(Handler, xrayconfig.WithRecommendedOptions(tp)...))

}
