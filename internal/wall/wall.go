// Package moat provides the first layer of defence for applications, by checking for attempted 'jailbreaks', and
// to scan a given text prompt for Personally Identifiable Information (PII) using an external PII scanner.
//
// Jailbreaking is the process of attempting to subvert the application's security defences, and is a commonly used
// to bypass ChatGPTs filters and generate offensive content. In this context, we are looking to see if a user is inputting
// user data to try and force the application or the LLM to generate a response that does something it shouldn't.
//
// For more information on PII scanning, see the pii and pii_aws packages documentation.
package wall

import (
	"fmt"
	"github.com/safetorun/PromptDefender/cache"
	"log"

	"github.com/safetorun/PromptDefender/badwords"
	"github.com/safetorun/PromptDefender/pii"
	"github.com/safetorun/PromptDefender/tracer"
)

type Wall struct {
	PiiScanner         pii.Scanner
	BadWordsCheck      *badwords.BadWords
	XmlEscapingScanner XmlEscapingScanner
	RemoteApiCaller    RemoteApiCaller
	logger             *log.Logger
	Cache              *cache.Cache
}

type PromptToCheck struct {
	Prompt           string
	ScanPii          *bool
	XmlTagToCheckFor *string
	CheckForBadWords *bool
	FastCheck        *bool
}

type PiiDetectionResult struct {
	ContainsPii bool
}

type CheckResult struct {
	PiiResult         *PiiDetectionResult
	ContainsBadWords  bool
	XmlScannerResult  *XmlEscapingDetectionResult
	InjectionDetected bool
}

type WallOpts func(*Wall) error

func New(opts ...WallOpts) (*Wall, error) {
	m := &Wall{}

	m.logger = log.Default()

	for _, opt := range opts {
		err := opt(m)
		if err != nil {
			return nil, err
		}
	}

	return m, nil
}

func (m *Wall) CheckWall(check PromptToCheck, t tracer.Tracer) (*CheckResult, error) {

	var piiResult *PiiDetectionResult = nil
	var xmlResult *XmlEscapingDetectionResult = nil
	var injectionDetected = false

	if m.Cache != nil {
		cached, cachedResult, err := checkCache(m.Cache, check)
		if err != nil {
			m.logger.Println(fmt.Sprintf("Error checking cache: %v", err))
			return nil, err
		}

		if cached {
			m.logger.Println("Cache hit")
			return cachedResult, nil
		} else {
			m.logger.Println("Cache miss")
		}
	}

	var containsBadWords = false
	if check.CheckForBadWords == nil {
		bw, err := m.checkPromptContainsBadwords(check, t)

		if err != nil {
			return nil, err
		}

		containsBadWords = *bw
	}

	if check.ScanPii == nil {
		piiResult = nil
	} else if *check.ScanPii {
		piiRe, err := m.checkPromptForPii(check, t)
		piiResult = piiRe

		if err != nil {
			return nil, err
		}
	} else {
		piiResult = &PiiDetectionResult{ContainsPii: false}
	}

	if check.XmlTagToCheckFor != nil {
		xmlResultInner, err := m.checkForXmlEscaping(check, t)
		if err != nil {
			return nil, err
		}

		xmlResult = xmlResultInner
	}

	if m.RemoteApiCaller != nil && (check.FastCheck == nil || !*check.FastCheck) {
		detected, err := m.checkForInjectionDetected(check, t)
		if err != nil {
			m.logger.Println("Error checking for injection ", err)
			return nil, err
		}
		injectionDetected = *detected
	}

	checkResult := CheckResult{
		PiiResult:         piiResult,
		ContainsBadWords:  containsBadWords,
		XmlScannerResult:  xmlResult,
		InjectionDetected: injectionDetected,
	}

	err := storeCache(m.Cache, check, &checkResult)

	if err != nil {
		m.logger.Printf("Error storing cache: %v\n", err)
	}

	return &checkResult, nil
}

func (m *Wall) checkForInjectionDetected(check PromptToCheck, t tracer.Tracer) (*bool, error) {
	wrappedMethod := tracer.TracerGenericsWrapper[string, MatchLevel](m.RemoteApiCaller.CallRemoteApi)
	injectionResponse, err := t.TraceDecorator(wrappedMethod, "checking_remote_api")(check.Prompt)
	injectionDetected := injectionResponse == ExactMatch || injectionResponse == VeryClose

	if err != nil {
		return nil, err
	}

	return &injectionDetected, nil
}

func (m *Wall) checkForXmlEscaping(check PromptToCheck, t tracer.Tracer) (*XmlEscapingDetectionResult, error) {
	wrappedMethod := tracer.TracerGenericsWrapper2[string, string, *XmlEscapingDetectionResult](m.XmlEscapingScanner.Scan)
	xmlResultInner, err := t.TraceDecorator(wrappedMethod, "xml_check")(check.Prompt, *check.XmlTagToCheckFor)
	return xmlResultInner.(*XmlEscapingDetectionResult), err
}

func (m *Wall) checkPromptForPii(check PromptToCheck, t tracer.Tracer) (*PiiDetectionResult, error) {
	wrappedMethod := tracer.TracerGenericsWrapper[string, *PiiDetectionResult](m.retrievePiiScore)
	piiRe, err := t.TraceDecorator(wrappedMethod, "scanning_pii")(check.Prompt)
	return piiRe.(*PiiDetectionResult), err
}

func (m *Wall) checkPromptContainsBadwords(check PromptToCheck, t tracer.Tracer) (*bool, error) {
	m.logger.Printf("Checking for bad words in prompt: %+v\n", check.Prompt)

	wrappedMethod := tracer.TracerGenericsWrapper[string, *bool](m.BadWordsCheck.CheckPromptContainsBadWords)

	// Execute m.BadWordsCheck.CheckPromptContainsBadWords with the wrapped method
	containsBadWords, err := t.TraceDecorator(wrappedMethod, "scanning_bad_words")(check.Prompt)

	if err != nil || containsBadWords == nil {
		return nil, fmt.Errorf("error checking bad words: %v", err)
	}

	return containsBadWords.(*bool), nil
}

func (m *Wall) retrievePiiScore(basePrompt string) (*PiiDetectionResult, error) {
	piiResult, err := m.PiiScanner.Scan(basePrompt)

	if err != nil {
		fmt.Printf("Error scanning for PII %v\n", err)
		return nil, err
	}

	return &PiiDetectionResult{ContainsPii: piiResult.ContainingPii}, nil
}
