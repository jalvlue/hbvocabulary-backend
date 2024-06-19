package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
)

func main() {
	fmt.Println("hello world")
	err := godotenv.Load("etc/.env")
	if err != nil {
		log.Fatal("error loading .env file")
	}

	fmt.Println(os.Getenv("DB_HOST"))
}
