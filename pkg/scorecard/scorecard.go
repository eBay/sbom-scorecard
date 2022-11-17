package scorecard

type SbomReport interface {
	Report() string
}

func PrettyPercent(num, denom int) int {
	if denom == 0 {
		return 0
	}
	return 100 * (1.0 * num) / denom
}
