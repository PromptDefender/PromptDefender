package embeddings

type EmbeddingValue struct {
	EmbeddingValue []float64
}

type Embeddings interface {
	CreateEmbeddings(string) (*EmbeddingValue, error)
	RetrieveBadwordsEmbeddings() (*[]EmbeddingValue, error)
}
