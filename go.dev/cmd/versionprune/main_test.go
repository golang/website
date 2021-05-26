package main

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"google.golang.org/api/appengine/v1"
)

func TestDecide(t *testing.T) {
	recentTime := time.Now()
	recentTime = recentTime.Truncate(time.Minute)
	allocs := map[string]float64{"currentlyServingID": 1.0}

	tests := []struct {
		desc       string
		versions   []*appengine.Version
		keepNumber int
		wantKeep   []version
		wantDelete []version
		wantErr    bool
	}{
		{
			desc: "no versions",
		},
		{
			desc:     "invalid Version time",
			versions: []*appengine.Version{{Id: "invalid time", CreateTime: "abc123"}},
			wantErr:  true,
		},
		{
			desc:     "old versions",
			versions: []*appengine.Version{{Id: "old one", CreateTime: "2018-01-02T15:04:05Z"}},
			wantDelete: []version{{
				AEVersion: &appengine.Version{Id: "old one", CreateTime: "2018-01-02T15:04:05Z"},
				Created:   time.Date(2018, 1, 2, 15, 4, 5, 0, time.UTC),
			}},
		},
		{
			desc:     "versions serving",
			versions: []*appengine.Version{{Id: "currentlyServingID", CreateTime: "2018-01-02T15:04:05Z"}},
			wantKeep: []version{{
				AEVersion: &appengine.Version{Id: "currentlyServingID", CreateTime: "2018-01-02T15:04:05Z"},
				Created:   time.Date(2018, 1, 2, 15, 4, 5, 0, time.UTC),
			}},
		},
		{
			desc:     "within 24h",
			versions: []*appengine.Version{{Id: "some id", CreateTime: recentTime.Format(time.RFC3339)}},
			wantKeep: []version{{
				AEVersion: &appengine.Version{Id: "some id", CreateTime: recentTime.Format(time.RFC3339)},
				Created:   recentTime,
			}},
		},
		{
			desc: "keeps KeepNumber versions",
			versions: []*appengine.Version{
				{Id: "some id", CreateTime: recentTime.Format(time.RFC3339)},
				{Id: "currentlyServingID", CreateTime: "2018-01-02T15:04:05Z"},
				{Id: "not serving", CreateTime: "2019-01-02T15:04:05Z"},
				{Id: "this one should be deleted", CreateTime: "2019-01-02T15:04:05Z"},
			},
			keepNumber: 3,
			wantKeep: []version{
				{
					AEVersion: &appengine.Version{Id: "some id", CreateTime: recentTime.Format(time.RFC3339)},
					Created:   recentTime,
				},
				{
					AEVersion: &appengine.Version{Id: "currentlyServingID", CreateTime: "2018-01-02T15:04:05Z"},
					Created:   time.Date(2018, 1, 2, 15, 4, 5, 0, time.UTC),
				},
				{
					AEVersion: &appengine.Version{Id: "not serving", CreateTime: "2019-01-02T15:04:05Z"},
					Created:   time.Date(2019, 1, 2, 15, 4, 5, 0, time.UTC),
				},
			},
			wantDelete: []version{{
				AEVersion: &appengine.Version{Id: "this one should be deleted", CreateTime: "2019-01-02T15:04:05Z"},
				Created:   time.Date(2019, 1, 2, 15, 4, 5, 0, time.UTC),
			}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.desc, func(t *testing.T) {
			bs, err := bucket(allocs, tt.versions, tt.keepNumber, 24*time.Hour)
			if (err != nil) != tt.wantErr {
				t.Errorf("bucket(%v, %v, %v, %v) = %v, %v, wantErr %v", allocs, tt.versions, tt.keepNumber, 24*time.Hour, bs, err, tt.wantErr)
				return
			}
			ignoreFields := cmpopts.IgnoreFields(version{}, "Message")
			if diff := cmp.Diff(tt.wantKeep, bs.keep, ignoreFields); diff != "" {
				t.Errorf("c.Keep mismatch (-want +got):\n%s", diff)
			}
			if diff := cmp.Diff(tt.wantDelete, bs.delete, ignoreFields); diff != "" {
				t.Errorf("c.Delete mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAct(t *testing.T) {
	bs := buckets{delete: []version{{AEVersion: &appengine.Version{Id: "test ID"}}}}
	if err := act(nil, bs, true); err != nil {
		t.Errorf("c.act() = %v, wanted no error", err)
	}
	defer func(t *testing.T) {
		t.Helper()
		if recover() == nil {
			// c.act() should panic with no asvs set, showing that the DryRun flag worked.
			// faking out the appengine admin client is hard.
			t.Errorf("recover() = nil, wanted panic")
		}
	}(t)
	act(nil, bs, false)
}
