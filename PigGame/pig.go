package main

import (
	"fmt"
	"math/rand"
)

const (
	winCondition   = 500
	gamesPerSeries = 10
)

type score struct {
	player, opponent, thisTurn int
}

//list of function with (score) type input and (score, bool) output
type action func(current score) (result score, turnIsOver bool)

//Both roll and stay are "action"-type

//roll return score of thisTurn and check if turnIsOver
func roll(s score) (score, bool) {
	outcome := rand.Intn(6) + 1
	if outcome == 1 {
		return score{s.opponent, s.player, 0}, true
	}
	return score{s.player, s.opponent, s.thisTurn + outcome}, false
}

func stay(s score) (score, bool) {
	return score{s.opponent, s.player + s.thisTurn, 0}, true
}

//if current score is "score", then we do "action"
type strategy func(current score) action

func stayAtK(k int) strategy {
	return func(s score) action {
		if s.thisTurn >= k {
			return stay
		}
		return roll
	}
}

func play(strategy0 strategy, strategy1 strategy) int {
	strategies := []strategy{strategy0, strategy1}
	var s score
	var turnIsOver bool
	currentPlayer := rand.Intn(2)
	for s.player+s.thisTurn < winCondition {
		action := strategies[currentPlayer](s)
		s, turnIsOver = action(s)
		if turnIsOver {
			currentPlayer = 1 - currentPlayer
		}
	}
	return currentPlayer
}

func roundRobin(strategies []strategy) ([]int, int) {
	wins := make([]int, len(strategies))
	for i := 0; i < len(strategies)-1; i++ {
		for j := i + 1; j < len(strategies); j++ {
			winner := play(strategies[i], strategies[j])
			if winner == 0 {
				wins[i]++
			} else {
				wins[j]++
			}
		}
	}
	gamesPerStrategy := gamesPerSeries * (len(strategies) - 1) // no self play
	return wins, gamesPerStrategy
}

func ratioString(vals ...int) string {
	total := 0
	for _, val := range vals {
		total += val
	}

	s := ""
	for _, val := range vals {
		if s != "" {
			s += ", "
		}
		pct := 100 * float64(val) / float64(total)
		s += fmt.Sprintf("%d/%d (%0.1f%%)", val, total, pct)
	}
	return s
}

func main() {
	strategies := make([]strategy, winCondition)
	for k := range strategies {
		strategies[k] = stayAtK(k + 1) //this is a function btw
	}
	wins, games := roundRobin(strategies)

	for k := range strategies {
		fmt.Printf("Wins, losses staying at k =% 4d: %s\n",
			k+1, ratioString(wins[k], games-wins[k]))
	}
}
