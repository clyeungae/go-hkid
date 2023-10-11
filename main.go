package main

import (
	"bytes"
	"fmt"
	"hkid/lib/crypto"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"golang.org/x/exp/slices"
)

const alphabet = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const checkDigit = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZ"

func main() {
	genTime, prefix, needAES := getUserInput()

	var key []byte
	var iv uint
	var blockSize uint
	if needAES {
		key, iv, blockSize = handleAESStructure()
	}

	for i := uint(0); i < genTime; i++ {
		plainText := genHKID(prefix)

		hashedText := ""
		if needAES {
			cipherText := crypto.AES256EncryptWithByteKey(plainText, key, iv, blockSize)
			hashedText = crypto.SHA256HashingWithByte(cipherText)
		}
		fmt.Printf("%v %v\n", plainText, hashedText)
	}
}

func getUserInput() (uint, string, bool) {
	var genTime uint
	var prefix string
	var needAES bool

	fmt.Println("Enter how many HKID you want")
	fmt.Scanln(&genTime)

	if genTime <= 0 {
		genTime = 10
	}

	var iPrefix string
	fmt.Println("Please specific a prefix (1 - 2 char)")
	fmt.Scanln(&iPrefix)

	// format checking
	re := regexp.MustCompile(`^[a-zA-Z]+$`)
	if !re.MatchString(iPrefix) {
		prefix = ""
	} else if len(prefix) > 2 {
		prefix = prefix[:2]
	} else {
		prefix = iPrefix
	}

	var iAES string
	fmt.Println("Do you need generate AES? (y/n)")
	fmt.Scanln(&iAES)

	yes := []string{"y", "yes"}
	needAES = slices.Contains(yes, iAES)

	return genTime, strings.ToUpper(prefix), needAES
}

func genHKID(prefix string) string {
	// introduce time to get rid of deterministic behavior of rand
	source := rand.NewSource(time.Now().UnixNano())
	randomGen := rand.New(source)

	var runeArr []string

	// random generator config
	maxChar := 9
	minChar := 8

	var totalChar int
	switch len(prefix) {
	case 2:
		totalChar = maxChar
	case 1:
		totalChar = minChar
	default:
		totalChar = randomGen.Intn(maxChar-minChar+1) + minChar
	}

	checkSum := 0

	totalLength := totalChar - 1
	leadingAlphabetLength := totalChar - 7

	for currentCharIndex := 0; currentCharIndex < totalLength; currentCharIndex++ {
		var char string
		if currentCharIndex < leadingAlphabetLength {
			switch len(prefix) {
			case 2:
				char = prefix[:1]
				prefix = prefix[1:2]
			case 1:
				char = prefix
				prefix = ""
			default:
				char = string(alphabet[randomGen.Intn(len([]rune(alphabet)))])
			}

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

func hexToByte(hex string) []byte {
	var byteBuilder bytes.Buffer

	for i := 0; i < len(hex)/2; i++ {
		byteBuilder.Write([]byte(hex[i*2 : i*2+1]))
	}
	return byteBuilder.Bytes()
}

func handleAESStructure() ([]byte, uint, uint) {
	// only for demo only, no hard code key and iv in prod
	key := "12345678901234567890123456789012"

	var blockSize uint
	var iv uint
	var iKey string
	var bKey []byte

	fmt.Println("Enter AES Key")
	fmt.Scanln(&iKey)

	if len(iKey) > 0 {
		key = iKey
	}

	switch len(key) {
	case 16, 24, 32:
		bKey = []byte(key)
		break
	default:
		bKey = hexToByte(key)
	}

	fmt.Println("Enter AES Block Size")
	fmt.Scanln(&blockSize)

	fmt.Println("Enter AES IV")
	fmt.Scanln(&iv)

	return bKey, iv, blockSize
}
