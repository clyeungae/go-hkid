package main

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphablet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
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
	// introduce time to get rid of deterministic behaviour of rand
	source := rand.NewSource(time.Now().UnixNano())
	randomGen := rand.New(source)

	var runeArr []string
	countDownNum := 8
	checkSum := 0

	for i := 0; i < 7; i++ {
		var char string
		if i == 0 {
			char = string(alphablet[randomGen.Intn(len([]rune(alphablet)))])
		} else {
			char = fmt.Sprint(randomGen.Intn(10))
		}
		runeArr = append(runeArr, char)

		checkSum += countDownNum * strings.Index(checkDigit, string(char))
		countDownNum--
	}

	checkSum = (checkSum + 324) % 11

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
