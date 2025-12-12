package wall

import (
	"fmt"
	"log"
	"strings"
)

type BasicXmlScanner struct {
	logger *log.Logger
}

func NewBasicXmlEscapingScaner() XmlEscapingScanner {
	return &BasicXmlScanner{
		logger: log.Default(),
	}
}

func (b BasicXmlScanner) Scan(textToScan string, tagContentToScanFor string) (*XmlEscapingDetectionResult, error) {

	b.logger.Println(fmt.Sprintf("Scanning text for tag: %s in input: %s\n", tagContentToScanFor, textToScan))

	if tagContentToScanFor == "" {
		return &XmlEscapingDetectionResult{
			ContainsXmlEscaping: false,
		}, nil
	}

	tagToScanFor := "<" + tagContentToScanFor + ">"
	otherTagToScanFor := "</" + tagContentToScanFor + ">"

	re := XmlEscapingDetectionResult{
		ContainsXmlEscaping: strings.Contains(textToScan, tagToScanFor) ||
			strings.Contains(textToScan, otherTagToScanFor),
	}

	return &re, nil
}
