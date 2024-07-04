package token

import "math/rand"

func getRandomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}
