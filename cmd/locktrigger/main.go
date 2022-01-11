// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// Locktrigger “locks” a given build trigger, making sure that
// the currently running build is the only trigger running.
//
// Usage:
//
//	locktrigger -project=$PROJECT_ID -build=$BUILD_ID
//
// The $PROJECT_ID and $BUILD_ID are typically written literally in cloudbuild.yaml
// and then substituted by Cloud Build.
//
// When a project uses “continuous deployment powered by Cloud Build”,
// the deployment is a little bit too continuous: when multiple commits
// land in a short time window, Cloud Build will run all the triggered
// build jobs in parallel. If each job does “gcloud app deploy”, there
// is no guarantee which will win: perhaps an older commit will complete
// last, resulting in the newest commit not actually being the final
// deployed version of the site. This should probably be fixed in
// “continuous deployment powered by Cloud Build”, but until then,
// locktrigger works around the problem.
//
// All triggered builds must run locktrigger to guarantee mutual exclusion.
// When there is contention—that is, when multiple builds are running and
// they all run locktrigger—the build corresponding to the newest commit
// is permitted to continue running, and older builds are canceled.
//
// When locktrigger exits successfully, then, at that moment, the current
// build is (or recently was) the only running build for its trigger.
// Of course, another build may start immediately after locktrigger exits.
// As long as that build also runs locktrigger, then either it will cancel
// itself (if it is older than we are), or it will cancel us before proceeding
// (if we are older than it is).
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	cloudbuild "cloud.google.com/go/cloudbuild/apiv1/v2"
	"google.golang.org/api/iterator"
	cloudbuildpb "google.golang.org/genproto/googleapis/devtools/cloudbuild/v1"
)

var (
	project = flag.String("project", "", "GCP project `name` (required)")
	build   = flag.String("build", "", "GCP build `id` (required)")
	repo    = flag.String("repo", "", "`URL` of repository (required)")
)

func usage() {
	fmt.Fprintf(os.Stderr, "usage: locktrigger -project=name -build=id -repo=URL\n")
	os.Exit(2)
}

func main() {
	flag.Usage = usage
	flag.Parse()
	log.SetPrefix("locktrigger: ")
	log.SetFlags(0)

	if *project == "" || *build == "" || *repo == "" {
		usage()
	}

	ctx := context.Background()
	c, err := cloudbuild.NewClient(ctx)
	if err != nil {
		log.Fatal(err)
	}
	defer c.Close()

	// Find commit hash of local Git
	myHash := run("git", "rev-parse", "HEAD")
	log.Printf("my hash: %v", myHash)

	// Find build object for current build, check that it matches.
	self := getBuild(c, ctx, *build)
	if hash := self.Substitutions["COMMIT_SHA"]; hash != myHash {
		log.Fatalf("build COMMIT_SHA does not match local hash: %v != %v", hash, myHash)
	}
	log.Printf("my build: %v", self.Id)
	if self.BuildTriggerId == "" {
		log.Fatalf("build has no trigger ID")
	}
	log.Printf("my trigger: %v", self.BuildTriggerId)

	// List all builds for our trigger that are still running.
	req := &cloudbuildpb.ListBuildsRequest{
		ProjectId: *project,
		// Note: Really want "status=WORKING buildTriggerId="+self.BuildTriggerId,
		// but that fails with an InvalidArgument error for unknown reasons.
		// status=WORKING will narrow the list down to something reasonable,
		// and we filter the unrelated triggers below.
		Filter: "status=WORKING",
	}
	it := c.ListBuilds(ctx, req)
	foundSelf := false
	shallow := false
	if _, err := os.Stat(run("git", "rev-parse", "--git-dir") + "/shallow"); err == nil {
		shallow = true
	}
	for {
		b, err := it.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			log.Fatalf("reading builds: %v (%q)", err, req.Filter)
		}
		if b.BuildTriggerId != self.BuildTriggerId {
			continue
		}

		// Check whether this build is an older or newer commit.
		// If this build is older, cancel it.
		// If this build is newer, cancel ourselves.
		if b.Id == self.Id {
			foundSelf = true
			continue
		}
		hash := b.Substitutions["COMMIT_SHA"]
		if hash == "" {
			log.Fatalf("cannot find COMMIT_SHA for build %v", b.Id)
		}
		if hash == myHash {
			log.Fatalf("found another build %v at same commit %v", b.Id, hash)
		}

		// Fetch the full Git repo so we can answer the history questions.
		// This is delayed until now to avoid the expense of fetching the full repo
		// if we are the only build that is running.
		if shallow {
			log.Printf("git fetch --unshallow")
			run("git", "fetch", "--unshallow", *repo)
			shallow = false
		}

		// Contention.
		// Find the common ancestor between us and that build,
		// to tell whether we're older, it's older, or we're unrelated.
		log.Printf("checking %v", hash)
		switch run("git", "merge-base", myHash, hash) {
		default:
			log.Fatalf("unexpected build for unrelated commit %v", hash)

		case myHash:
			// myHash is older than b's hash. Cancel self.
			log.Printf("canceling self, for build %v commit %v", b.Id, hash)
			cancel(c, ctx, self.Id)

		case hash:
			// b's hash is older than myHash. Cancel b.
			log.Printf("canceling build %v commit %v", b.Id, hash)
			cancel(c, ctx, b.Id)
		}
	}

	// If we listed all the in-progress builds, we should have seen ourselves.
	if !foundSelf {
		log.Fatalf("reading builds: didn't find self")
	}
}

// getBuild returns the build info for the build with the given id.
func getBuild(c *cloudbuild.Client, ctx context.Context, id string) *cloudbuildpb.Build {
	req := &cloudbuildpb.GetBuildRequest{
		ProjectId: *project,
		Id:        id,
	}
	b, err := c.GetBuild(ctx, req)
	if err != nil {
		log.Fatalf("getbuild %v: %v", id, err)
	}
	return b
}

// cancel cancels the build with the given id.
func cancel(c *cloudbuild.Client, ctx context.Context, id string) {
	req := &cloudbuildpb.CancelBuildRequest{
		ProjectId: *project,
		Id:        id,
	}
	_, err := c.CancelBuild(ctx, req)
	if err != nil {
		// Not Fatal: maybe cancel failed because the build exited.
		// Waiting for it to stop running below will take care of that case.
		log.Printf("cancel %v: %v", id, err)
	}

	// Wait for build to report being stopped,
	// in case cancel only queues the cancellation and doesn't actually wait,
	// or in case cancel failed.
	// Willing to wait a few minutes.
	now := time.Now()
	for time.Since(now) < 3*time.Minute {
		b := getBuild(c, ctx, id)
		if b.Status != cloudbuildpb.Build_WORKING {
			log.Printf("canceled %v: now %v", id, b.Status)
			return
		}
		time.Sleep(10 * time.Second)
	}
	log.Fatalf("cancel %v: did not stop", id)
}

// run runs the given command line and returns the standard output, with spaces trimmed.
func run(args ...string) string {
	var stdout, stderr bytes.Buffer
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		log.Fatalf("exec %v: %v\n%s%s", args, err, stdout.String(), stderr.String())
	}
	return strings.TrimSpace(stdout.String())
}
