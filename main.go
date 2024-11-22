package main

import (
	// "github.com/gofiber/fiber/v2"
	// "fmt"

	"github.com/eagledb14/paperlink/engagement"
	// "github.com/eagledb14/paperlink/net"
)

func main() {
	// net.Run()

	edb := engagement.NewEngagementDb()
	// edb.DeleteEngagement("test1")
	// _ = edb
	// err := edb.CreateBlankEngagementDb("test1", "test person", "test@test.com")
	// err := e
	// fmt.Println(err)
	tdb := engagement.NewTemplateDb()
	tdb.CreateEngagement("test1", "1", "1")
	tdb.CreateEngagement("test2", "2", "2")

	edb.CreateEngagementFromTemplate("monkey", "test2", "2", "2")
	// tdb.DeleteEngagement("test2")
	// err := tdb.CreateTemplate("test2")
	// fmt.Println(err)

	// err = edb.CreateEngagementFromTemplate("test2", "", "")
	// fmt.Println(err)
}
