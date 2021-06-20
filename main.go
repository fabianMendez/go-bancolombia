package main

import (
	"log"
	"os"

	// "os"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	cl, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	user, found := os.LookupEnv("AUTH_USER")
	if !found {
		log.Fatal("User required")
	}

	password, found := os.LookupEnv("AUTH_PASS")
	if !found {
		log.Fatal("Password required")
	}

	err = cl.Login(user, password)
	if err != nil {
		log.Fatal("could not login:", err)
	}

	defer cl.Logout()

	err = cl.GetDepositsBalance()
	if err != nil {
		log.Fatal("could not get deposits balance:", err)
	}
}
