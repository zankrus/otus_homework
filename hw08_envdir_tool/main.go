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
	log.Printf("dir: %v\n", dirPath)

	env, err := ReadDir(dirPath)
	if err != nil {
		log.Fatal(err)
	}

	command := os.Args[2:]
	log.Printf("command: %v\n", command)

	RunCmd(command, env)
}
