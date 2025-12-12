package pii

type ScanResult struct {
	ContainingPii bool
}

type Scanner interface {
	Scan(pii string) (*ScanResult, error)
}
