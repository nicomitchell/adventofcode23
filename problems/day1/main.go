package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"strings"
)

var problemPart = 1

func main() {
	if len(os.Args) < 2 {
		log.Fatal("missing argument: input file")
	}
	inputFilepath := os.Args[1]
	if len(os.Args) == 3 {
		problemPartS := os.Args[2]
		if problemPartS == "2" {
			problemPart = 2
		} else if problemPartS != "1" {
			fmt.Println("unrecognized problem stage; assuming stage 1")
		}
	}
	f, err := os.Open(inputFilepath)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	reader := bufio.NewReader(f)
	var sum int
	var readErr error
	for readErr == nil {
		var line []byte
		line, readErr = reader.ReadBytes('\n')
		if readErr != nil && readErr != io.EOF {
			log.Fatal(readErr)
		}
		if len(line) == 0 {
			continue
		}
		num, err := parseLine(line)
		if err != nil {
			fmt.Printf("invalid number from line %s", string(line))
			continue
		}
		fmt.Println("got number from line: ", num)
		sum += num
	}
	fmt.Println("Got sum: ", sum)
	return
}

func parseLine(line []byte) (int, error) {
	var f, l byte
	var fIndexSS, lIndexSS int
	if problemPart == 2 {
		f, l, fIndexSS, lIndexSS = parseDigitsAsSubStrings(string(line))
	}

	// 'a' allows us to easily confirm whether we've found a digit yet
	var firstDigit, lastDigit byte = 'a', 'a'
	firstIndexDigit, lastIndexDigit := len(line), -1
	fInd := bytes.IndexAny(line, "0123456789")
	if fInd != -1 {
		firstIndexDigit = fInd
		firstDigit = line[fInd]
	}
	lInd := bytes.LastIndexAny(line, "0123456789")
	if lInd != -1 {
		lastIndexDigit = lInd
		lastDigit = line[lInd]
	}
	if problemPart == 2 {
		if f != 'a' && firstIndexDigit > fIndexSS {
			firstIndexDigit = fIndexSS
			firstDigit = f
		}
		if l != 'a' && lastIndexDigit < lIndexSS {
			lastIndexDigit = lIndexSS
			lastDigit = l
		}
	}
	// original manual iter solution
	// for i := range line {
	// 	if firstDigit != 'a' && lastDigit != 'a' {
	// 		break
	// 	}
	// 	if firstDigit == 'a' {
	// 		if line[i] >= '0' && line[i] <= '9' {
	// 			firstDigit = line[i]
	// 			firstIndexDigit = i
	// 		}
	// 	}
	// 	if lastDigit == 'a' {
	// 		if line[len(line)-1-i] >= '0' && line[len(line)-1-i] <= '9' {
	// 			lastDigit = line[len(line)-1-i]
	// 			lastIndexDigit = len(line) - 1 - i
	// 		}
	// 	}
	// }
	numString := string([]byte{firstDigit, lastDigit})
	return strconv.Atoi(numString)
}

// returns first/last numbers parsed from substrings as well as the index they are found in
// numbers returned as bytes for easy type comparison with parseLine output
var validSubstringsForDigits = map[string]byte{
	"one":   '1',
	"two":   '2',
	"three": '3',
	"four":  '4',
	"five":  '5',
	"six":   '6',
	"seven": '7',
	"eight": '8',
	"nine":  '9',
	"zero":  '0',
}

func parseDigitsAsSubStrings(line string) (f byte, l byte, fIndex int, lIndex int) {
	f, l = 'a', 'a'
	fIndex, lIndex = len(line), -1
	for substr, digit := range validSubstringsForDigits {
		ind := strings.Index(line, substr)
		if ind != -1 {
			if ind < fIndex {
				fIndex = ind
				f = digit
			}
		}
		lInd := strings.LastIndex(line, substr)
		if lInd != -1 {
			if lInd > lIndex {
				lIndex = lInd //+ len(substr) - 1
				l = digit
			}
		}
	}
	return
}
