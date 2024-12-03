package main

import (
	"fmt"
	"testing"

	"github.com/eagledb14/paperlink/engagement"
)


func cleanupEngagements(dbList ...string) {
	engagements := engagement.LoadEngagements()
	fmt.Println(len(engagements))

	for _, e := range engagements {
		for _, db := range dbList {
			if e.Name == db {
				e.Delete()
			}
		}
	}

}

func cleanupTemplates(dbList ...string) {
	engagements := engagement.LoadTemplates()

	for _, e := range engagements {
		for _, db := range dbList {
			if e.Name == db {
				e.Delete()
			}
		}
	}

}

func contains(egs []engagement.Engagement, name string) bool {
	for _, e := range egs {
		if e.Name == name {
			return true
		}
	}
	return false
}

func assertEngagements(egs []engagement.Engagement, name ...string) {

}

func TestEngagementInsert(t *testing.T) {
	defer cleanupEngagements("test1", "test2", "test3")

	engagement.NewEngagement("test1", "test1", "test1@email")
	engagement.NewEngagement("test2", "test2", "test2@email")
	engagement.NewEngagement("test3", "test3", "test3@email")

	engagements := engagement.LoadEngagements()

	if len(engagements) < 3 {
		panic("not enough engagements")
	}

	if !contains(engagements, "test1") {
		panic("missing test1")
	}
	if !contains(engagements, "test2") {
		panic("missing test2")
	}
	if !contains(engagements, "test3") {
		panic("missing test3")
	}
}

func TestTemplateInsert(t *testing.T) {
	defer cleanupTemplates("test1", "test2", "test3")

	engagement.NewTemplate("test1")
	engagement.NewTemplate("test2")
	engagement.NewTemplate("test3")

	engagements := engagement.LoadTemplates()

	if len(engagements) < 3 {
		panic("not enough engagements")
	}

	if !contains(engagements, "test1") {
		panic("missing test1")
	}
	if !contains(engagements, "test2") {
		panic("missing test2")
	}
	if !contains(engagements, "test3") {
		panic("missing test3")
	}
}

func TestEngagementFromTemplate(t *testing.T) {
	defer cleanupTemplates("test1", "test2", "test3")
	defer cleanupEngagements("test1", "test2", "test3")

	engagement.NewTemplate("test1")
	engagement.NewTemplate("test2")
	engagement.NewTemplate("test3")

	engagement.NewEngagementFromTemplate("test1", "test1", "test1", "test1@email")
	engagement.NewEngagementFromTemplate("test2", "test2", "test2", "test2@email")
	engagement.NewEngagementFromTemplate("test3", "test3", "test3", "test3@email")

	engagements := engagement.LoadEngagements()

	if len(engagements) < 3 {
		panic("not enough engagements")
	}

	if !contains(engagements, "test1") {
		panic("missing test1")
	}
	if !contains(engagements, "test2") {
		panic("missing test2")
	}
	if !contains(engagements, "test3") {
		panic("missing test3")
	}
}
