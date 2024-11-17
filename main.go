package main

import (
	"fmt"
	"gemini-gen-ai/infrastructure"
)

func main() {
	defer fmt.Println("Server stopped...")
	fmt.Println("Server working....")
	infrastructure.Env()
}
