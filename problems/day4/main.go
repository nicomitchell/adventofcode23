package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
)

var problemPart = "1"

func main() {
	if len(os.Args) < 2 {
		log.Fatal("missing argument: input file")
	}
	inputFilepath := os.Args[1]
	if len(os.Args) == 3 {
		problemPart = os.Args[2]
	}
	f, err := os.Open(inputFilepath)
	defer f.Close()
	if err != nil {
		log.Fatal(err)
	}
	var sum int
	if problemPart == "1" {
		reader := bufio.NewReader(f)
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
			fmt.Print(string(line))
			nums := bytes.Split(line, []byte{':'})[1]
			parts := bytes.Split(nums, []byte{'|'})
			winningNumbers := getNumbers(bytes.Trim(parts[0], " \n"))
			myNumbers := getNumbers(bytes.Trim(parts[1], " \n"))
			lineSum := 0
			for n, _ := range myNumbers {
				if winningNumbers[n] {
					if lineSum == 0 {
						lineSum = 1
					} else {
						lineSum *= 2
					}
				}
			}
			sum += lineSum
		}
	} else if problemPart == "2" {
		raw, err := io.ReadAll(f)
		if err != nil {
			log.Fatal(err)
		}
		lines := bytes.Split(raw, []byte{'\n'})
		cardCopies := make([]int, len(lines))
		// 1 copy of each card
		for i := range cardCopies {
			cardCopies[i] = 1
		}
		for i, line := range lines {
			sum += cardCopies[i]
			nums := bytes.Split(line, []byte{':'})[1]
			parts := bytes.Split(nums, []byte{'|'})
			winningNumbers := getNumbers(bytes.Trim(parts[0], " \n"))
			myNumbers := getNumbers(bytes.Trim(parts[1], " \n"))
			var matches int
			for n, _ := range myNumbers {
				if winningNumbers[n] {
					matches++
				}
			}
			for copyIndex := i + 1; copyIndex <= i+matches && copyIndex < len(cardCopies); copyIndex++ {
				cardCopies[copyIndex] += cardCopies[i]
			}
		}
	}
	fmt.Println("Got sum: ", sum)
	return
}

// using a map to make comparing them easier
func getNumbers(line []byte) map[int]bool {
	out := make(map[int]bool)
	nums := bytes.Split(line, []byte{' '})
	for _, unparsed := range nums {
		if len(unparsed) == 0 {
			continue
		}
		n, err := strconv.Atoi(string(unparsed))
		if err != nil {
			log.Fatal(err)
		}
		out[n] = true
	}
	return out
}
