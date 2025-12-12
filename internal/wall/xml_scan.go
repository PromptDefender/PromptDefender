package wall

type XmlEscapingDetectionResult struct {
	ContainsXmlEscaping bool
}

type XmlEscapingScanner interface {
	Scan(textToScan string, tagToScanFor string) (*XmlEscapingDetectionResult, error)
}
