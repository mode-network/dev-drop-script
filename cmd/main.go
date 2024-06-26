package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/mode-network/dev-drop-script/scripts"
	"log"
	"os"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file", err)
	}

	if len(os.Args) != 3 {
		log.Println(fmt.Sprintf("Usage: create-dev-drop input.csv output.json"))
		os.Exit(1)
	}

	inFile := os.Args[1]
	outFile := os.Args[2]
	scripts.GenerateDevDropSafeFile(inFile, outFile, scripts.GetConfig())
}
