package main

import (
	"bufio"
	"log"
	"os"
)

var Words = make([]string, 0)

func LoadWords() {
	file, err := os.Open("words.txt")

	if err != nil {
		log.Fatal("Couldn't load words from words.txt", err)
	}

	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		Words = append(Words, fileScanner.Text())
	}

	file.Close()
}
