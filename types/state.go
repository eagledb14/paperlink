package types

import (
	"errors"

	"github.com/eagledb14/paperlink/engagement"
)


type State struct {
	Engagements []engagement.Engagement
	Templates []engagement.Engagement
	EngagementMap map[string]int
	TemplateMap map[string]int
}

func NewState() *State {
	newState := State{}
	newState.Engagements = engagement.LoadEngagements()
	newState.Templates = engagement.LoadTemplates()

	newState.EngagementMap = make(map[string]int)
	newState.TemplateMap = make(map[string]int)

	for i, e := range newState.Engagements {
		newState.EngagementMap[e.Name] = i
	}
	for i, t := range newState.Templates {
		newState.TemplateMap[t.Name] = i
	}

	return &newState
}

func (s *State) GetEngagement(name string) (*engagement.Engagement, error) {
	index, exists := s.EngagementMap[name]
	if !exists {
		return nil, errors.New("Missing Enagement")

	}
	e := s.Engagements[index]
	
	return &e, nil
}

func (s *State) GetTemplate(name string) (*engagement.Engagement, error) {
	index, exists := s.TemplateMap[name]
	if !exists {
		return nil, errors.New("Missing Template")

	}
	e := s.Templates[index]
	
	return &e, nil
}

func (s *State) AddEngagement(newEngagement engagement.Engagement) {
	s.Engagements = append(s.Engagements, newEngagement)
	s.EngagementMap[newEngagement.Name] = len(s.Engagements) - 1
}

func (s *State) AddTemplate(newTemplate engagement.Engagement) {
	s.Templates = append(s.Templates, newTemplate)
	s.TemplateMap[newTemplate.Name] = len(s.Templates) - 1
}

func (s *State) DeleteEnagement(name string) {
	index := s.EngagementMap[name]
	s.Engagements[index].Delete()
	s.Engagements = append(s.Engagements[:index], s.Engagements[index+1:]...)

	for i, e := range s.Engagements {
		s.EngagementMap[e.Name] = i
	}
	for i, t := range s.Templates {
		s.TemplateMap[t.Name] = i
	}
}

func (s *State) DeleteTemplate(name string) {
	index := s.TemplateMap[name]
	s.Templates[index].Delete()
	s.Templates = append(s.Templates[:index], s.Templates[index+1:]...)


	for i, e := range s.Engagements {
		s.EngagementMap[e.Name] = i
	}
	for i, t := range s.Templates {
		s.TemplateMap[t.Name] = i
	}
}
