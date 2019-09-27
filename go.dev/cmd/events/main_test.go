package main

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
)

var (
	noEventGroup       = &Group{Name: "no event group"}
	upcomingEventGroup = &Group{
		Name: "Upcoming Event Group",
		Timezone: "Europe/Oslo",
		NextEvent: &Event{
			ID:            "12345",
			Name:          "Upcoming Event",
			Time:          1262976000000,
		},
	}
	fakeGroups = map[string]*Group{
		"noEvent": noEventGroup,
		"ueg":     upcomingEventGroup,
	}
)

type fakeClient struct{}

func (f fakeClient) getGroupsSummary() (*GroupsSummary, error) {
	return &GroupsSummary{Chapters: []*Chapter{
		{URLName: "noEvent"},
		{
			URLName: "ueg",
			Description: "We host our own events\n",
		},
	}}, nil
}

func (f fakeClient) getGroup(urlName string) (*Group, error) {
	g, ok := fakeGroups[urlName]
	if !ok {
		return nil, fmt.Errorf("no group %q", urlName)
	}
	return g, nil
}

func TestGetUpcomingEvents(t *testing.T) {
	want := &UpcomingEvents{All: []EventData{
		{
			Name:        "Upcoming Event",
			ID:          "12345",
			Description: "We host our own events<br/>\n",
			LocalDate:   "Jan 8, 2010",
			LocalTime:   "2010-01-08T19:40:00+01:00",
			URL:         "https://www.meetup.com/ueg/events/12345",
		},
	}}
	f := fakeClient{}
	got, err := getUpcomingEvents(f)
	if err != nil {
		t.Fatalf("getUpcomingEvents(%v) error = %v, wanted no error", f, err)
	}
	if diff := cmp.Diff(want, got); diff != "" {
		t.Errorf("getUpcomingEvents(%v) mismatch (-want +got):\n%s", f, diff)
	}
}
