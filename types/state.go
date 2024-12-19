package types

import (
	"errors"
	"sort"
	"sync"

	"github.com/eagledb14/paperlink/auth"
	"github.com/eagledb14/paperlink/dictionary"
	"github.com/eagledb14/paperlink/engagement"
)


type State struct {
	Engagements []engagement.Engagement
	Templates []engagement.Engagement
	lock sync.RWMutex

	Dictionary dictionary.Dictionary
	Auth *auth.Auth
	Clients sync.Map
}

func NewState() *State {
	newState := State{}
	newState.Engagements = engagement.LoadEngagements()
	sortEngagements(newState.Engagements)

	newState.Templates = engagement.LoadTemplates()
	sortEngagements(newState.Templates)

	newState.Dictionary = dictionary.LoadDictionary()
	newState.Auth = auth.NewAuth()

	return &newState
}

func (s *State) GetEngagement(name string) (*engagement.Engagement, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, e := range s.Engagements {
		if e.Name == name {
			return &e, nil
		}
	}
	
	return nil, errors.New("Engagement Not Found")
}

func (s *State) GetTemplate(name string) (*engagement.Engagement, error) {
	s.lock.RLock()
	defer s.lock.RUnlock()
	for _, e := range s.Templates {
		if e.Name == name {
			return &e, nil
		}
	}
	
	return nil, errors.New("Template Not Found")
}

func (s *State) AddEngagement(newEngagement engagement.Engagement) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Engagements = append([]engagement.Engagement{newEngagement}, s.Engagements...)
}

func (s *State) AddTemplate(newTemplate engagement.Engagement) {
	s.lock.Lock()
	defer s.lock.Unlock()
	s.Templates = append([]engagement.Engagement{newTemplate}, s.Templates...)
}

func (s *State) DeleteEnagement(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for i, e := range s.Engagements {
		if e.Name == name {
			s.Engagements = append(s.Engagements[:i], s.Engagements[i+1:]...)
			return
		}
	}
}

func (s *State) DeleteTemplate(name string) {
	s.lock.Lock()
	defer s.lock.Unlock()

	for i, e := range s.Templates {
		if e.Name == name {
			s.Templates = append(s.Templates[:i], s.Templates[i+1:]...)
			return
		}
	}
}

func sortEngagements(e []engagement.Engagement) {
	sort.Slice(e, func(i, j int) bool {
		return e[j].TimeStamp.Before(e[i].TimeStamp)
	})
}
