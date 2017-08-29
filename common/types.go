package common

import (
	"encoding/json"
	"sort"
	"time"
)

type ContactSuggestion struct {
	Start     string `json:"start"`
	BeginTime time.Time
	End       string `json:"end"`
	EndTime   time.Time
	Users     []string `json:"users"`
}

type ContactSuggestions []ContactSuggestion

func (slice ContactSuggestions) Len() int {
	return len(slice)
}

func (slice ContactSuggestions) Less(x, y int) bool {
	return slice[x].BeginTime.Before(slice[y].BeginTime)
}

func (slice ContactSuggestions) Swap(x, y int) {
	slice[x], slice[y] = slice[y], slice[x]
}

type contact ContactSuggestion

func (this *ContactSuggestion) UnmarshalJSON(b []byte) error {
	var suggestion contact
	err := json.Unmarshal(b, &suggestion)

	startTime, err := time.Parse(time.RFC3339, suggestion.Start)
	if err != nil {
		return err
	}
	suggestion.BeginTime = startTime
	endTime, err := time.Parse(time.RFC3339, suggestion.End)
	if err != nil {
		return err
	}
	suggestion.EndTime = endTime
	*this = ContactSuggestion(suggestion)
	return err
}

func NewContactSuggestion(start time.Time, end time.Time, users []string) *ContactSuggestion {
	return &ContactSuggestion{BeginTime: start, EndTime: end, Users: users}
}

type Suggestion struct {
	Text string `json:"text"`
	URL  string `json:"url"`
}

func NewSuggestion(text string, url string) *Suggestion {
	return &Suggestion{Text: text, URL: url}
}

type Config struct {
	Token            string `json:"token"`
	User             string `json:"user"`
	SurpressMessages int    `json:"surpressMessageDuration"`
	Message          string `json:"message`
	Start            string `json:"start"`
	InitalizedTime   time.Time
	StartTime        time.Time
	EndTime          time.Time
	End              string             `json:"end"`
	Suggestions      []Suggestion       `json:"suggestions"`
	Contacts         ContactSuggestions `json:"contacts"`
}

type config Config

func (this *Config) UnmarshalJSON(b []byte) error {
	var config config
	err := json.Unmarshal(b, &config)
	start, err := time.Parse(time.RFC3339, config.Start)
	if err != nil {
		return err
	}

	config.StartTime = start
	end, err := time.Parse(time.RFC3339, config.End)
	if err != nil {
		return err
	}
	config.EndTime = end
	sort.Sort(config.Contacts)
	*this = Config(config)
	return err
}
