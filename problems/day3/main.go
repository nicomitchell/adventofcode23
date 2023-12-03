package main

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
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
	raw, err := io.ReadAll(f)
	if err != nil {
		log.Fatal(err)
	}
	var sum int
	lines := bytes.Split(raw, []byte{'\n'})
	if problemPart == 1 {
		symbolCoords := getSymbolCoordinates(lines)
		coordsToParse := make(map[coordinate]bool)
		for _, symbolCoord := range symbolCoords {
			adjacentDigitCoords := getAdjacentDigits(symbolCoord, lines)
			for _, coord := range adjacentDigitCoords {
				//use a map to maybe avoid collisions
				coordsToParse[coord] = true
			}
		}
		for coord := range coordsToParse {
			// we need to mark all the coords of the digit as false when we check them
			// to avoid double counting numbers
			if coordsToParse[coord] {
				num, checked := parseNumberAtCoord(coord, lines)
				sum += num
				for _, checkedCoord := range checked {
					if _, ok := coordsToParse[checkedCoord]; ok {
						coordsToParse[checkedCoord] = false
					}
				}
			}
		}
	} else if problemPart == 2 {
		gearCoords := getGearCoordinates(lines)
		for _, coord := range gearCoords {
			coordsToParse := make(map[coordinate]bool)
			adjacentDigits := getAdjacentDigits(coord, lines)
			for _, adjCoord := range adjacentDigits {
				//use a map to maybe avoid collisions
				coordsToParse[adjCoord] = true
			}
			// we only count it if there are 2 adjacent numbers
			numberOfAdjacentNumbers := 0
			gearRatio := 1
			for digitCoord := range coordsToParse {
				if coordsToParse[digitCoord] {
					number, checked := parseNumberAtCoord(digitCoord, lines)
					for _, checkedCoord := range checked {
						if _, ok := coordsToParse[checkedCoord]; ok {
							coordsToParse[checkedCoord] = false
						}
					}
					numberOfAdjacentNumbers++
					gearRatio *= number
				}
			}
			if numberOfAdjacentNumbers == 2 {
				sum += gearRatio
			}
		}
	}
	fmt.Println("Got sum: ", sum)
	return
}

type coordinate struct {
	x, y int
}

// returns coordinates of all symbols
func getSymbolCoordinates(lines [][]byte) []coordinate {
	out := []coordinate{}
	for y, line := range lines {
		for x, symbol := range line {
			//assume any char that's not a digit or period is a symbol
			if symbol != '.' && (symbol < '0' || symbol > '9') {
				out = append(out, coordinate{x: x, y: y})
			}
		}
	}
	return out
}

func getAdjacentDigits(coord coordinate, lines [][]byte) []coordinate {
	digitCoords := []coordinate{}
	for yMod := -1; yMod <= 1; yMod++ {
		for xMod := -1; xMod <= 1; xMod++ {
			checked := coordinate{x: coord.x + xMod, y: coord.y + yMod}
			val := lines[checked.y][checked.x]
			if val >= '0' && val <= '9' {
				digitCoords = append(digitCoords, checked)
			}
		}
	}
	return digitCoords
}

// returns a list of the coords that make up the number
func parseNumberAtCoord(coord coordinate, lines [][]byte) (int, []coordinate) {
	parsed := []coordinate{coord}
	minX, maxX := coord.x, coord.x
	for x := coord.x - 1; x >= 0; x-- {
		if lines[coord.y][x] >= '0' && lines[coord.y][x] <= '9' {
			minX = x
			parsed = append(parsed, coordinate{x: x, y: coord.y})
		} else {
			break
		}
	}
	for x := coord.x + 1; x < len(lines[coord.y]); x++ {
		if lines[coord.y][x] >= '0' && lines[coord.y][x] <= '9' {
			maxX = x
			parsed = append(parsed, coordinate{x: x, y: coord.y})
		} else {
			break
		}
	}
	number, err := strconv.Atoi(string(lines[coord.y][minX : maxX+1]))
	if err != nil {
		log.Fatal(err)
	}
	return number, parsed
}

// returns coordinates of all gears (*)
func getGearCoordinates(lines [][]byte) []coordinate {
	out := []coordinate{}
	for y, line := range lines {
		for x, symbol := range line {
			//assume any char that's not a digit or period is a symbol
			if symbol == '*' {
				out = append(out, coordinate{x: x, y: y})
			}
		}
	}
	return out
}
