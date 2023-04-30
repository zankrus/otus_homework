package main

import (
	"log"
	"os"
)

func main() {

	argLength := len(os.Args[1:])

	if argLength < 3 {
		log.Fatalf("Not enought arguments got %v, but  need min 3", argLength)
	}

	dirPath := os.Args[1]

	env, err := ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	command := os.Args[2:]

	os.Exit(RunCmd(command, env))
}
