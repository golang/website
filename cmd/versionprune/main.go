// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"sort"
	"time"

	"google.golang.org/api/appengine/v1"
)

var (
	dryRun       = flag.Bool("dry_run", true, "print but do not run changes")
	keepDuration = flag.Duration("keep_duration", 24*time.Hour, "keep versions with age < `t`")
	keepNumber   = flag.Int("keep_number", 5, "keep at least `n` versions")
	project      = flag.String("project", "", "GCP project `name` (required)")
	service      = flag.String("service", "", "AppEngine service `name` (required)")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: versionprune -project=name -service=name [options]\n")
	flag.PrintDefaults()
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	log.SetPrefix("versionprune: ")
	log.SetFlags(0)

	if *project == "" || *service == "" {
		usage()
	}
	if *keepDuration < 0 {
		log.Fatalf("-keep_duration=%v must be >= 0", *keepDuration)
	}
	if *keepNumber < 0 {
		log.Fatalf("-keep_number=%d must be >= 0", *keepNumber)
	}

	if err := run(context.Background()); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(2)
	}
}

// version is an intermediate representation of an AppEngine version.
type version struct {
	AEVersion *appengine.Version
	// Message is a human-readable message of why a version was bucketed.
	Message string
	// Created is a parsed time of the AppEngine Version CreateTime.
	Created time.Time
}

// run fetches, processes, and (if DryRun is false) deletes stale AppEngine versions for the specified Service.
func run(ctx context.Context) error {
	aes, err := appengine.NewService(ctx)
	if err != nil {
		return fmt.Errorf("creating appengine client: %w", err)
	}
	ass := appengine.NewAppsServicesService(aes)
	asvs := appengine.NewAppsServicesVersionsService(aes)

	s, err := ass.Get(*project, *service).Do()
	if err != nil {
		return fmt.Errorf("fetching service: %w", err)
	}

	vs, err := getVersions(ctx, asvs)
	if err != nil {
		return fmt.Errorf("fetching versions: %w", err)
	}

	bs, err := bucket(s.Split.Allocations, vs, *keepNumber, *keepDuration)
	if err != nil {
		return fmt.Errorf("bucketing versions: %w", err)
	}
	printIntent(bs)
	if err := act(asvs, bs, *dryRun); err != nil {
		return fmt.Errorf("executing: %w", err)
	}
	return nil
}

func getVersions(ctx context.Context, asvs *appengine.AppsServicesVersionsService) ([]*appengine.Version, error) {
	var versions []*appengine.Version
	err := asvs.List(*project, *service).Pages(ctx, func(r *appengine.ListVersionsResponse) error {
		versions = append(versions, r.Versions...)
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Slice(versions, func(i, j int) bool {
		// Sort by create time, descending.
		return versions[i].CreateTime > versions[j].CreateTime
	})
	return versions, nil
}

type buckets struct {
	keep   []version
	delete []version
}

// bucket splits c.versions into intended actions.
func bucket(allocs map[string]float64, versions []*appengine.Version, keepNumber int, keepDuration time.Duration) (buckets, error) {
	var bs buckets
	for _, av := range versions {
		v := version{AEVersion: av}
		created, err := time.Parse(time.RFC3339, av.CreateTime)
		if err != nil {
			return bs, fmt.Errorf("failed to parse time %q for version %s: %v", av.CreateTime, av.Id, err)
		}
		v.Created = created
		if s, ok := allocs[av.Id]; ok {
			v.Message = fmt.Sprintf("version is serving traffic. split: %v%%", s*100)
			bs.keep = append(bs.keep, v)
			continue
		}
		if len(bs.keep) < keepNumber {
			v.Message = fmt.Sprintf("keeping the latest %d versions (%d)", keepNumber, len(bs.keep))
			bs.keep = append(bs.keep, v)
			continue
		}
		if dur := time.Since(v.Created); dur < keepDuration {
			v.Message = fmt.Sprintf("keeping recent versions (%s)", dur)
			bs.keep = append(bs.keep, v)
			continue
		}
		v.Message = "bye"
		bs.delete = append(bs.delete, v)
	}
	return bs, nil
}

func printIntent(bs buckets) {
	fmt.Printf("target project:\t[%v]\n", *project)
	fmt.Printf("target service:\t[%v]\n", *service)
	fmt.Printf("versions:\t(%v)\n", len(bs.delete)+len(bs.keep))
	fmt.Println()

	fmt.Printf("versions to keep (%v): [\n", len(bs.keep))
	for _, v := range bs.keep {
		fmt.Printf("\t%v: %v\n", v.AEVersion.Id, v.Message)
	}
	fmt.Println("]")

	fmt.Printf("versions to delete (%v): [\n", len(bs.delete))
	for _, v := range bs.delete {
		fmt.Printf("\t%v: %v\n", v.AEVersion.Id, v.Message)
	}
	fmt.Println("]")
}

// act performs delete requests for AppEngine services to be deleted. No deletions are performed if DryRun is true.
func act(asvs *appengine.AppsServicesVersionsService, bs buckets, dryRun bool) error {
	for _, v := range bs.delete {
		if dryRun {
			fmt.Printf("dry-run: skipping delete %v: %v\n", v.AEVersion.Id, v.Message)
			continue
		}
		fmt.Printf("deleting %v/%v/%v\n", *project, *service, v.AEVersion.Id)
		if _, err := asvs.Delete(*project, *service, v.AEVersion.Id).Do(); err != nil {
			return fmt.Errorf("error deleting %v/%v/%v: %w", *project, *service, v.AEVersion.Id, err)
		}
	}
	return nil
}
