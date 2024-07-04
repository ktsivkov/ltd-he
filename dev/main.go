package main

import (
	"fmt"
	"log"
	"time"

	"github.com/ktsivkov/ltd-he/pkg/token"
)

func main() {
	tokenService := token.NewService()
	now, _ := time.Parse("2022-Jan-02", "2013-Feb-03")
	// Define a slice of integers
	t, err := tokenService.Token("Rammtein#2934", 454, 277, 1300, 5, 0, 53, 96, now, true)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Printf("The generated token is: %v\n", t)

	isValid, err := tokenService.ValidateToken("Rammtein#2934", t)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("The token is: %v\n", isValid)

	presetToken := "FGZ{d42brL1f?RHUxIgk=6Rbbe1@3ISilu2ux#zmN~UN|[p>@ov=%K}R"
	isValid5, err := tokenService.ValidateToken("Rammtein#2934", presetToken)
	if err != nil {
		log.Println(err)
	}
	fmt.Printf("The token is: %v\n", isValid5)
}
