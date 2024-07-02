package token

import "math/rand"

func formatToken(a string) string {
	var b = len(a)
	var d string
	for b > 0 {
		b--
		c := a[b : b+1]
		if c == "-" {
			d = "-" + d
		} else {
			d = c + d
		}
	}
	return d
}

func encodeCharacter(a string) int {
	c := "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	d := "abcdefghijklmnopqrstuvwxyz"
	e := "0123456789"

	// Search for the character in the uppercase alphabet
	for i := 0; i < 26; i++ {
		if c[i:i+1] == a {
			return i
		}
		if d[i:i+1] == a {
			return i
		}
	}

	// Search for the character in the digits
	for i := 0; i < 10; i++ {
		if e[i:i+1] == a {
			return i
		}
	}

	// If character not found, return 0
	return 0
}

func getRandomInt(min, max int) int {
	return rand.Intn(max-min+1) + min
}
