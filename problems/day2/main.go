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
		fmt.Print(string(line))
		if problemPart == 1 {
			id, bux, err := getIDAndBucketsFromLine(line)
			if err != nil {
				log.Fatal(err)
				continue
			}
			if gamesPossible(12, 13, 14, bux) {
				sum += id
			}
		} else if problemPart == 2 {
			mins, err := getMinsForGames(bytes.Split(line, []byte{':'})[1])
			if err != nil {
				log.Fatal(err)
				continue
			}
			fmt.Println(mins)
			sum += minCubesPower(mins)
		}
	}
	fmt.Println("Got sum: ", sum)
	return
}

// each arr in buckets should be an array of size 3 with the r,g,b values inside it
func gamesPossible(maxR, maxG, maxB int, buckets map[string][]int) bool {
	RGBMaxes := map[string]int{"red": maxR, "green": maxG, "blue": maxB}
	for color, max := range RGBMaxes {
		for _, count := range buckets[color] {
			if count > max {
				return false
			}
		}
	}
	return true
}

// gets the power required of each number of cubes
func minCubesPower(mins map[string]int) int {
	product := 1
	for _, v := range mins {
		product *= v
	}
	return product
}

func getIDAndBucketsFromLine(line []byte) (int, map[string][]int, error) {
	//line format: [Game {#}: {#} {color}, {#} color; {#} color, {#} color, {#} color]
	parts := bytes.Split(bytes.Trim(line, " \n"), []byte{':'})
	id, err := getIdFromPrefix(parts[0])
	if err != nil {
		return 0, nil, err
	}
	buckets, err := getBuckets(parts[1])
	if err != nil {
		return 0, nil, err
	}
	return id, buckets, nil

}

func getIdFromPrefix(prefix []byte) (int, error) {
	//prefix fmt: [Game {id}]
	parts := bytes.Split(prefix, []byte{' '})
	idBytes := parts[1]
	return strconv.Atoi(string(idBytes))
}

func getBuckets(line []byte) (map[string][]int, error) {
	// line format: [{#} {color}, {#} {color}; {#} {color}, {#} {color}, {#} {color}]
	buckets := make(map[string][]int)
	bucketLines := bytes.Split(line, []byte{';'})
	for _, bucketLine := range bucketLines {
		colors := bytes.Split(bucketLine, []byte{','})
		for _, color := range colors {
			trimmed := bytes.Trim(color, " ")
			parts := bytes.Split(trimmed, []byte{' '})
			count, err := strconv.Atoi(string(parts[0]))
			if err != nil {
				return nil, err
			}
			color := string(parts[1])
			if counts, ok := buckets[color]; ok {
				buckets[color] = append(counts, count)
			} else {
				buckets[color] = []int{count}
			}
		}
	}
	return buckets, nil
}

func getIDAndMinsFromLine(line []byte) (int, map[string]int, error) {
	//line format: [Game {#}: {#} {color}, {#} color; {#} color, {#} color, {#} color]
	parts := bytes.Split(bytes.Trim(line, " \n"), []byte{':'})
	id, err := getIdFromPrefix(parts[0])
	if err != nil {
		return 0, nil, err
	}
	mins, err := getMinsForGames(parts[1])
	if err != nil {
		return 0, nil, err
	}
	return id, mins, nil

}

func getMinsForGames(line []byte) (map[string]int, error) {
	// line format: [{#} {color}, {#} {color}; {#} {color}, {#} {color}, {#} {color}]
	mins := make(map[string]int)
	bucketLines := bytes.Split(bytes.Trim(line, " \n"), []byte{';'})
	for _, bucketLine := range bucketLines {
		colors := bytes.Split(bucketLine, []byte{','})
		for _, color := range colors {
			trimmed := bytes.Trim(color, " ")
			parts := bytes.Split(trimmed, []byte{' '})
			count, err := strconv.Atoi(string(parts[0]))
			if err != nil {
				return nil, err
			}
			color := string(parts[1])
			if counts, ok := mins[color]; ok {
				if count > counts {
					mins[color] = count
				}
			} else {
				mins[color] = count
			}
		}
	}
	return mins, nil
}
