package main

import (
	// "github.com/eagledb14/paperlink/engagement"

	"github.com/eagledb14/paperlink/net"
	// "fmt"
	// "github.com/eagledb14/paperlink/types"
)

func main() {
	// e[0].db.Query(`SELECT name, contact, email, timeStamp FROM engagements`)
	// s := types.NewState()
	// fmt.Println(s.TemplateMap)
	// engagement.NewEngagementFromTemplate("test1", "test1", "thing", "other thing")
	// engagement.NewEngagementFromTemplate("test1", "test2", "thing", "other thing")
	// engagement.NewEngagementFromTemplate("test1", "test with spaces 1", "thing", "other thing")
	// engagement.NewEngagementFromTemplate("test1", "test with spaces 2", "thing", "other thing")
	net.Run()
	// e := engagement.LoadEngagements()
	// fmt.Println(e)
	// fmt.Println(e[0].UpdateEngagement("cheese", e[0].Contact, e[0].Email))
	// fmt.Println(e)
	// e := engagement.NewEngagementFromTemplate("test1", "monkey", "1", "2")
	// e.UpdateEngagement("cheese", e.Contact, e.Email)
	// fmt.Println(e)
	//
	// t := engagement.LoadTemplates()
	// fmt.Println(t)
	// t := engagement.NewTemplate("test1")
	// _ = t
	// e := engagement.NewEngagement("monkey", "1", "2")
	// t := engagement.NewTemplate("test1")
	// for i := range 10 {
	// 	istr := strconv.Itoa(i)
	// 	t.InsertSection(i, "title", "body")
	// 	t.InsertAsset(istr, "name", "type")
	// 	t.InsertFinding(i, "title", time.Now(), "summary", "desc")
	// }
	// e := engagement.NewEngagementFromTemplate("test1", "monkey", "1", "2")
	// engagement.LoadEngagements()

	// e.UpdateEngagement("")
	// e.InsertSection(1, )
	// _ = e
	// e.Close()
	// t.Close()
	// net.Run()
	// edb := engagement.NewEngagementDb()
	// tdb := engagement.NewTemplateDb()
	// tdb.CreateEngagement("test1", "", "")
	//
	// var wg sync.WaitGroup
	// for i := range 1_000 {
	// 	wg.Add(1)
	// 	go func(i int, wg *sync.WaitGroup, edb *engagement.EngagementDb) {
	// 		defer wg.Done()
	// 		istr := strconv.Itoa(i)
	// 		err := edb.CreateEngagement(istr, istr, istr)
	// 		_ = err
	// 		// edb.DeleteEngagement(istr)
	// 		// fmt.Println(i, err)
	// 	}(i, &wg, &edb)
	// }
	// wg.Wait()
	// edb.CreateEngagementFromTemplate("monkeyman", "test1", "1", "2")

	// engagement.NewEngagement("")
	// wg := sync.WaitGroup{}
	// for i := range 1_000 {
	// 	wg.Add(1)
	// 	go func(i int, wg *sync.WaitGroup) {
	// 		defer wg.Done()
	// 		istr := strconv.Itoa(i)
	// 		e := engagement.NewEngagement(istr, istr, istr)
	// 		_ = e
	// 		// err := e.CreateEngagement(istr, istr, istr)
	// 		fmt.Println(i)
	// 	}(i, &wg)
	// }
	// wg.Wait()
	// edb.DeleteEngagement("test1")
	// _ = edb
	// err := edb.CreateBlankEngagementDb("test1", "test person", "test@test.com")
	// err := e
	// fmt.Println(err)
	// tdb := engagement.NewTemplateDb()
	// tdb.CreateEngagement("test1", "1", "1")
	// tdb.CreateEngagement("test2", "2", "2")
	//
	// edb.CreateEngagementFromTemplate("monkey", "test2", "2", "2")
	// tdb.DeleteEngagement("test2")
	// err := tdb.CreateTemplate("test2")
	// fmt.Println(err)

	// err = edb.CreateEngagementFromTemplate("test2", "", "")
	// fmt.Println(err)
}
