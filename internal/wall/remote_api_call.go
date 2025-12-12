package wall

type MatchLevel int

const (
	ExactMatch MatchLevel = iota
	VeryClose
	Medium
	NoMatch
	TotallyDifferent
)

type RemoteApiCaller interface {
	// Call the remote API and return the response which will contain
	// an injection score
	CallRemoteApi(string) (MatchLevel, error)
}
