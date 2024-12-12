package main

import (
	"flag"
	"fmt"

	"github.com/eagledb14/paperlink/auth"
)

func main() {
	admin := flag.String("admin", "", "Creates new admin account with temporary password")
	flag.Parse()
	if *admin == "" {
		Run()
	} else {
		a := auth.NewAuth()
		tempPass, _ := a.GenerateCookie()

		fmt.Println(tempPass)
		a.NewUser(*admin, tempPass, true)
	}
}
