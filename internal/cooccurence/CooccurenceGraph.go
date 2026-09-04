package cooccurence

import (
	"math"
	"sort"

	"github.com/AnneeAvakyan/litanalyzer/internal/domain/entities"
)

type EdgeWeight map[[2]int]int

func BuildGraph(mentions []entities.Mention, windowSize int) EdgeWeight {
	weight := make(EdgeWeight)
	floatWeight := make(map[[2]int]float64)

	sort.Slice(mentions, func(i, j int) bool {
		return mentions[i].GlobalSentenceIndex < mentions[j].GlobalSentenceIndex
	})

	for i := 0; i < len(mentions); i++ {
		for j := i + 1; j < len(mentions); j++ {
			if mentions[j].GlobalSentenceIndex-mentions[i].GlobalSentenceIndex > windowSize {
				break
			}

			if mentions[i].CharacterID == mentions[j].CharacterID {
				continue
			}

			key := edgeKey(mentions[i].CharacterID, mentions[j].CharacterID)
			distance := mentions[j].GlobalSentenceIndex - mentions[i].GlobalSentenceIndex
			floatWeight[key] += 1.0 / float64(distance+1)
		}
	}

	for k, v := range floatWeight {
		weight[k] += int(math.Round(v))
	}

	return weight
}

func edgeKey(a, b int) [2]int {
	if a < b {
		return [2]int{a, b}
	}
	return [2]int{b, a}
}
