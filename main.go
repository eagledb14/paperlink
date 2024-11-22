package main

import (
	"strconv"
	"time"

	"github.com/eagledb14/paperlink/engagement"
)

func main() {
	// t := engagement.NewTemplate("test1")
	// _ = t
	e := engagement.NewEngagement("monkey", "1", "2")
	for i := range 10 {
		istr := strconv.Itoa(i)
		e.InsertSection(i, "title", "body")
		e.InsertAsset(istr, "name", "type")
		e.InsertFinding(i, "title", time.Now(), "summary", "desc")
	}
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
