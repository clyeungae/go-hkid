package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const checkDigit = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func main() {
	genTime := getUserInput()

	if genTime <= 0 {
		return
	}

	for i := 0; i < genTime; i++ {
		fmt.Println(genHKID())
	}
}

func getUserInput() int {
	var genTime int

	fmt.Println("Enter how many HKID you want")
	fmt.Scan(&genTime)

	return genTime
}

func genHKID() string {
	// introduce time to get rid of deterministic behavior of rand
	source := rand.NewSource(time.Now().UnixNano())
	randomGen := rand.New(source)

	var runeArr []string

	// random generator config
	maxChar := 9
	minChar := 8

	totalChar := randomGen.Intn(maxChar-minChar+1) + minChar
	checkSum := 0

	totalLength := totalChar - 1
	leadingAlphabetLength := totalChar - 7

	for currentCharIndex := 0; currentCharIndex < totalLength; currentCharIndex++ {
		var char string
		if currentCharIndex < leadingAlphabetLength {
			char = string(alphabet[randomGen.Intn(len([]rune(alphabet)))])
		} else {
			char = fmt.Sprint(randomGen.Intn(10))
		}
		runeArr = append(runeArr, char)

		if currentCharIndex == 0 {
			if totalLength == 8 {
				checkSum += 9 * strings.Index(checkDigit, string(char))
			} else {
				checkSum += 9 * 36
				checkSum += (totalChar - currentCharIndex) * strings.Index(checkDigit, string(char))
			}
		} else {
			checkSum += (totalChar - currentCharIndex) * strings.Index(checkDigit, string(char))
		}
	}

	checkSum %= 11

	lastDigit := "0"
	if checkSum != 0 {
		lastDigit = strings.Split(checkDigit, "")[11-checkSum]
	}

	runeArr = append(runeArr, "(", fmt.Sprint(lastDigit), ")")
	return joinString((runeArr))
}

func joinString(strs []string) string {
	var strBuilder strings.Builder
	for _, r := range strs {
		strBuilder.WriteString(r)
	}
	return strBuilder.String()
}
