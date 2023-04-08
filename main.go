package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"syscall"

	"github.com/dustin/go-humanize"
	"github.com/joho/godotenv"
	"github.com/olekukonko/tablewriter"
	"golang.org/x/term"
)

func fmtMoney(n float64) string {
	return "$ " + humanize.FormatFloat("#,###.##", n)
}

func main() {
	_ = godotenv.Load()

	cl, err := NewClient()
	if err != nil {
		log.Fatal(err)
	}

	user, password, err := readCredentials()
 if err != nil {
  log.Fatal("couldn't read credentials: ", err)
 }

	err = cl.Login(user, password)
	if err != nil {
		log.Fatal("could not login: ", err)
	}
	log.Println("logged in")

	defer cl.Logout()

	/*err = cl.AccountEnrroll()
	if err != nil {
		log.Fatal("could not enroll account: ", err)
	} else {
		log.Println("account enrolled")
		os.Exit(0)
	}*/

	balance, err := cl.GetDepositsBalance()
	if err != nil {
		log.Fatal("could not get deposits balance:", err)
	}

	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader([]string{"Tipo", "Número", "Saldo disponible"})
	table.Append([]string{balance.Description, balance.Number, "$ " + balance.AvailableBalance})
	table.Render()

	// fmt.Printf("%#v\n", balance)

	tableB := tablewriter.NewWriter(os.Stdout)
	tableB.SetHeader([]string{"Fecha", "Oficina", "Descripción", "Referencia", "Monto"})
	tableB.SetColumnAlignment([]int{tablewriter.ALIGN_DEFAULT, tablewriter.ALIGN_DEFAULT,
		tablewriter.ALIGN_LEFT, tablewriter.ALIGN_LEFT, tablewriter.ALIGN_RIGHT})

	for i := 1; i <= 3; i++ {
		detail, err := cl.GetSavingsDetail(balance.ID, i%3)
		if err != nil {
			log.Fatal("could not get savings detail:", i, err)
		}
		for _, d := range detail {
			row := []string{d.Date.Format("2006/01/02"), d.BranchID, d.Description, d.OptionalRef, fmtMoney(d.Amount)}
			fgcolor := tablewriter.FgGreenColor
			if d.Amount < 0 {
				fgcolor = tablewriter.FgRedColor
			} else if d.Amount < 50000 {
				fgcolor = tablewriter.FgYellowColor
			}
			// tableB.Append(row)
			tableB.Rich(row, []tablewriter.Colors{{}, {}, {}, {}, {fgcolor}})
		}
	}

	tableB.Render()
}

func readCredentials() (string, string, error) {
 var err error
 rdr := bufio.NewReader(os.Stdin)

	user, found := os.LookupEnv("AUTH_USER")
	if !found {
		fmt.Print("User: ")
		user, err = rdr.ReadString('\n')
		if err != nil {
			return "", "", fmt.Errorf("User required")
		}
		user = strings.TrimSpace(user)
	}

	password, found := os.LookupEnv("AUTH_PASS")
	if !found {
		fmt.Print("Password: ")
		bytePassword, err := term.ReadPassword(int(syscall.Stdin))
		if err != nil {
			return "", "", fmt.Errorf("Password required")
		}
		fmt.Println()
		password = string(bytePassword)
	}
 return user, password, nil
}
