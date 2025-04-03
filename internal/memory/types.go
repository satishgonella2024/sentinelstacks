// Package memory provides memory storage implementations
package memory

import "github.com/satishgonella2024/sentinelstacks/pkg/types"

// SimilarityMatch represents a similarity search result for backwards compatibility
// This should be replaced with types.SimilarityResult in the future
type SimilarityMatch struct {
	Key      string
	Score    float32
	Metadata map[string]interface{}
}

// ToSimilarityResult converts a SimilarityMatch to a types.SimilarityResult
func (m SimilarityMatch) ToSimilarityResult() types.SimilarityResult {
	return types.SimilarityResult{
		ID:       m.Key,
		Score:    m.Score,
		Metadata: m.Metadata,
	}
}

// SimilarityResultsToMatches converts a slice of types.SimilarityResult to a slice of SimilarityMatch
func SimilarityResultsToMatches(results []types.SimilarityResult) []SimilarityMatch {
	matches := make([]SimilarityMatch, len(results))
	for i, result := range results {
		matches[i] = SimilarityMatch{
			Key:      result.ID,
			Score:    result.Score,
			Metadata: result.Metadata,
		}
	}
	return matches
}

// MatchesToSimilarityResults converts a slice of SimilarityMatch to a slice of types.SimilarityResult
func MatchesToSimilarityResults(matches []SimilarityMatch) []types.SimilarityResult {
	results := make([]types.SimilarityResult, len(matches))
	for i, match := range matches {
		results[i] = match.ToSimilarityResult()
	}
	return results
}
