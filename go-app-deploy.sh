#!/bin/bash
# Copyright 2021 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# This script is meant to be run from Cloud Build as a substitute
# for "gcloud app deploy", as in:
#
#	go-app-deploy.sh app.yaml
#
# It should not be run by hand and is therefore not marked executable.
#
# It customizes the usual "gcloud app deploy" in two ways.
#
# First, it sets --no-promote or --promote according to whether
# the commit had a Website-Publish vote.
#
# Second, it chooses an app version like 2021-06-02-204309-2c120970
# giving the date, time, and commit hash of the commit being deployed.
# This handles the case where multiple commits are being run through
# Cloud Build at once and would otherwise end up with timestamps in
# arbitrary order depending on the order in which Cloud Build happened
# to reach each's gcloud app deploy command. With our choice, the
# versions are ordered in git order.

set -e

promote=$(
	git cat-file -p 'HEAD' |
	awk '
		BEGIN { flag = "--no-promote" }
		/^Reviewed-on:/ { flag = "--no-promote" }
		/^Website-Publish:/ { flag = "--promote" }
		END {print flag}
	'
)

version=$(
	git log -n1 --date='format:%Y-%m-%d-%H%M%S' --pretty='format:%cd-%h'
)

gcloud app deploy $promote -v $version "$@"
