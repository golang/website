#!/bin/bash
# Copyright 2021 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# This script is meant to be run from Cloud Build as a substitute
# for "gcloud app deploy", as in:
#
#	go-app-deploy.sh [--project=name] app.yaml
#
# It should not be run by hand and is therefore not marked executable.
#
# It customizes the usual "gcloud app deploy" in two ways.
#
# First, it deploys with --no-promote and only promotes after
# the deployed app has passed an in-prod readiness test.
#
# Second, it chooses an app version like 2021-06-02-204309-2c120970
# giving the date, time, and commit hash of the commit being deployed.
# This handles the case where multiple commits are being run through
# Cloud Build at once and would otherwise end up with timestamps in
# arbitrary order depending on the order in which Cloud Build happened
# to reach each's gcloud app deploy command. With our choice, the
# versions are ordered in git time order.

set -e

project=golang-org
case "$1" in
--project=*)
	project=$(echo $1 | sed 's/--project=//')
	shift
esac

yaml=app.yaml
case "$1" in
*.yaml)
	yaml=$1
	shift
esac

if [ $# != 0 ]; then
	echo 'usage: go-app-deploy.sh [--project=name] path/to/app.yaml' >&2
	exit 2
fi

version=$(git log -n1 --date='format:%Y-%m-%d-%H%M%S' --pretty='format:%cd-%h')

service=$(awk '$1=="service:" {print $2}' $yaml)

servicedot="-$service-dot"
if [ "$service" = default ]; then
	servicedot=""
fi
host="$version-dot$servicedot-$project.appspot.com"

echo "### deploying to https://$host"
gcloud -q --project=$project app deploy -v $version --no-promote $yaml

curl --version

for i in 1 2 3 4 5; do
	if curl -s --fail --show-error "https://$host/_readycheck"; then
		echo '### site is up!'

		# We used to have to worry about go-app-deploy.sh running for
		# other commits simultaneously, so we checked that we weren't
		# going to stomp a newer version of the web site.
		# The use of locktrigger in cloudbuild.yaml should make this impossible,
		# but keep checking anyway.
		serving=$(gcloud app services describe --project=$project $service | grep ': 1.0')
		if [ "$serving" '>' "$version" ]; then
			echo "### serving version $serving is newer than our $version; not promoting"
			exit 1
		fi

		echo '### promoting'
		gcloud -q --project=$project app services set-traffic $service --splits=$version=1
		exit 0
	fi
	echo '### not healthy'
	curl "https://$host/_readycheck" # show response body
done

echo "### failed to become healthy; giving up"
exit 1
