package render

import (
	"fmt"
	"math"
	"strings"

	"github.com/yardnsm/gohever"
)

const (
	BarLength = 33

	BarFill  = "="
	BarBlank = "·"
)

func CardBalance(status gohever.CardStatus) string {
	barParts := int(math.Floor(float64(BarLength) / float64(len(status.Factors))))

	numBarFill := int(math.Ceil(status.MonthlyUsage * BarLength))
	numBarBlank := BarLength - numBarFill

	var (
		render         []string
		topTimeline    []string
		bottomTimeline []string
		amount         int
	)

	for i := range status.Factors {
		factor := status.Factors[i]
		amount += int(factor.Amount)
		percentage := int(math.Round(100 - 100*factor.Factor))

		// I thought that maybe in the future I would like to calculate the length based on the
		// relative space the factors takes, instead of it being fixed.
		nodeLength := barParts

		var (
			node   []string
			filler []string
		)

		//	   1000₪
		//         │
		//     30% │

		node = append(node, padLeft(fmt.Sprintf("%d₪", amount), nodeLength))
		node = append(node, padLeft("│", nodeLength))
		node = append(node, padLeft(fmt.Sprintf("%d%% │", percentage), nodeLength))

		// Filler is basically an empty node
		filler = append(filler, node...)
		for i := range filler {
			filler[i] = strings.Repeat(" ", nodeLength)
		}

		if i%2 == 0 {
			// Top timeline
			topTimeline = append(topTimeline, strings.Join(node, "\n"))
			bottomTimeline = append(bottomTimeline, strings.Join(filler, "\n"))
		} else {
			// Bottom timeline
			topTimeline = append(topTimeline, strings.Join(filler, "\n"))
			bottomTimeline = append(bottomTimeline, strings.Join(node, "\n"))
		}
	}

	render = append(render, concatLines(topTimeline))
	render = append(render, strings.Repeat(BarFill, numBarFill)+strings.Repeat(BarBlank, numBarBlank))
	render = append(render, reverseString(concatLines(bottomTimeline)))

	return strings.Join(render, "\n")
}
