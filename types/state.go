package types

import (
	"errors"
	"sort"

	"github.com/eagledb14/paperlink/auth"
	"github.com/eagledb14/paperlink/dictionary"
	"github.com/eagledb14/paperlink/engagement"
)


type State struct {
	Engagements []engagement.Engagement
	Templates []engagement.Engagement
	EngagementMap map[string]int
	TemplateMap map[string]int
	Dictionary dictionary.Dictionary
	Auth *auth.Auth
}

func NewState() *State {
	newState := State{}
	newState.Engagements = engagement.LoadEngagements()
	sortEngagements(newState.Engagements)

	newState.Templates = engagement.LoadTemplates()
	sortEngagements(newState.Templates)

	newState.EngagementMap = make(map[string]int)
	newState.TemplateMap = make(map[string]int)

	updateMap(newState.Engagements, newState.EngagementMap)
	updateMap(newState.Templates, newState.TemplateMap)

	newState.Dictionary = dictionary.LoadDictionary()
	newState.Auth = auth.NewAuth()

	return &newState
}


func (s *State) GetEngagement(name string) (*engagement.Engagement, error) {
	index, exists := s.EngagementMap[name]
	if !exists || index >= len(s.Engagements) {
		return nil, errors.New("Missing Engagement")
	}

	e := s.Engagements[index]
	
	return &e, nil
}

func (s *State) GetTemplate(name string) (*engagement.Engagement, error) {
	index, exists := s.TemplateMap[name]
	if !exists || index >= len(s.Templates) {
		return nil, errors.New("Missing Template")
	}

	e := s.Templates[index]
	
	return &e, nil
}

func (s *State) AddEngagement(newEngagement engagement.Engagement) {
	s.Engagements = append([]engagement.Engagement{newEngagement}, s.Engagements...)
	updateMap(s.Engagements, s.EngagementMap)
}

func (s *State) AddTemplate(newTemplate engagement.Engagement) {
	s.Templates = append([]engagement.Engagement{newTemplate}, s.Templates...)
	updateMap(s.Templates, s.TemplateMap)
}

func (s *State) DeleteEnagement(name string) {
	index := s.EngagementMap[name]
	s.Engagements[index].Delete()
	s.Engagements = append(s.Engagements[:index], s.Engagements[index+1:]...)

	updateMap(s.Engagements, s.EngagementMap)
}

func (s *State) DeleteTemplate(name string) {
	index := s.TemplateMap[name]
	s.Templates[index].Delete()
	s.Templates = append(s.Templates[:index], s.Templates[index+1:]...)

	updateMap(s.Templates, s.TemplateMap)
}

func sortEngagements(e []engagement.Engagement) {
	sort.Slice(e, func(i, j int) bool {
		return e[j].TimeStamp.Before(e[i].TimeStamp)
	})
}

func updateMap(engagements []engagement.Engagement, m map[string]int) {
	for i, e := range engagements {
		m[e.Name] = i
	}
}
