package main

import (
	"bufio"
	"log"
	"math/rand"
	"os"
)

var words = make([]string, 0)

func LoadWords() {
	file, err := os.Open("words.txt")

	if err != nil {
		log.Fatal("Couldn't load words from words.txt", err)
	}

	fileScanner := bufio.NewScanner(file)

	for fileScanner.Scan() {
		words = append(words, fileScanner.Text())
	}

	file.Close()
}

func RandomWord() string {
	return words[rand.Intn(len(words))]
}
