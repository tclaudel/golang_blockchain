package values

import (
	"fmt"
	"strings"
)

type MiningDifficulty struct {
	difficulty int
	zeros      string
}

const (
	easiest = iota + 2
	easy
	medium
	hard
)

var (
	Easiest = MiningDifficulty{difficulty: easiest, zeros: strings.Repeat("0", easiest)}
	Easy    = MiningDifficulty{difficulty: easy, zeros: strings.Repeat("0", easy)}
	Medium  = MiningDifficulty{difficulty: medium, zeros: strings.Repeat("0", medium)}
	Hard    = MiningDifficulty{difficulty: hard, zeros: strings.Repeat("0", hard)}
)

func MiningDifficultyFromInt(difficulty int) (MiningDifficulty, error) {
	switch difficulty {
	case easiest:
		return Easiest, nil
	case easy:
		return Easy, nil
	case medium:
		return Medium, nil
	case hard:
		return Hard, nil
	default:
		return MiningDifficulty{}, fmt.Errorf("invalid difficulty: %d", difficulty)
	}
}

func (d MiningDifficulty) Difficulty() int {
	return d.difficulty
}

func (d MiningDifficulty) Zeros() string {
	return strings.Repeat("0", d.difficulty)
}
