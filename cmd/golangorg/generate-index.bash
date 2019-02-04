#!/usr/bin/env bash

# Copyright 2011 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# This script creates a .zip file representing the $GOROOT file system
# and computes the corresponding search index files.
#
# These are used in production (see app.prod.yaml)

set -e -u -x

ZIPFILE=golangorg.zip
INDEXFILE=golangorg.index
SPLITFILES=index.split.

error() {
	echo "error: $1"
	exit 2
}

install() {
	go install
}

getArgs() {
	if [ ! -v GOLANGORG_DOCSET ]; then
		GOLANGORG_DOCSET="$(go env GOROOT)"
		echo "GOLANGORG_DOCSET not set explicitly, using GOROOT instead"
	fi

	# safety checks
	if [ ! -d "$GOLANGORG_DOCSET" ]; then
		error "$GOLANGORG_DOCSET is not a directory"
	fi

	# reporting
	echo "GOLANGORG_DOCSET = $GOLANGORG_DOCSET"
}

makeZipfile() {
	echo "*** make $ZIPFILE"
	rm -f $ZIPFILE goroot
	ln -s "$GOLANGORG_DOCSET" goroot
	zip -q -r $ZIPFILE goroot/* # glob to ignore dotfiles (like .git)
	rm goroot
}

makeIndexfile() {
	echo "*** make $INDEXFILE"
	golangorg=$(go env GOPATH)/bin/golangorg
	# NOTE: run golangorg without GOPATH set. Otherwise third-party packages will end up in the index.
	GOPATH= $golangorg -write_index -goroot goroot -index_files=$INDEXFILE -zip=$ZIPFILE
}

splitIndexfile() {
	echo "*** split $INDEXFILE"
	rm -f $SPLITFILES*
	split -b8m $INDEXFILE $SPLITFILES
}

cd $(dirname $0)

install
getArgs "$@"
makeZipfile
makeIndexfile
splitIndexfile
rm $INDEXFILE

echo "*** setup complete"
