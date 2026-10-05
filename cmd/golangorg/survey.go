// Copyright 2026 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	_ "embed"
	"time"
)

type surveyConfig struct {
	ID, Name, URL string
	Start, End    time.Time
}

//go:embed survey.json
var surveyJSON []byte

func surveyBanner() *surveyConfig {
	if disableTemporaryBanner {
		return nil
	}
	now := time.Now()
	if now.Before(currentSurvey.Start) || now.After(currentSurvey.End) {
		return nil
	}
	return &currentSurvey
}
