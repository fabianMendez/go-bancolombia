package main

import (
	"fmt"
	"log"
	"os"

	// "os"

	"github.com/joho/godotenv"
)

func main() {
	_ = godotenv.Load()

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
		log.Fatal("could not login: ", err)
	}

	defer cl.Logout()

	balance, err := cl.GetDepositsBalance()
	if err != nil {
		log.Fatal("could not get deposits balance:", err)
	}
	fmt.Printf("%s - %s: (%s) %s\n", balance.ProductName, balance.Number, balance.Currency, balance.AvailableBalance)

	for i := 1; i <= 6; i++ {
		detail, err := cl.GetSavingsDetail(i % 3)
		if err != nil {
			log.Fatal("could not get savings detail:", i, err)
		}
		for _, d := range detail {
			fmt.Println(d.BranchID, d.Amount, d.Date, d.Description, d.OptionalRef)
		}
	}
}
