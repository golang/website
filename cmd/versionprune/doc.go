// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/*
Versionprune prunes stale AppEngine versions for a specified service.

The command by default will:
- keep the latest 5 versions
- keep any version that is serving traffic
- keep any version that is younger than 24h
- ignore versions with invalid dates (doesn't seem to be a real thing)

	Sample output:

	target project:	[go-discovery]
	target service:	[go-dev]
	versions:	(18)

	versions to keep (11): [
		20191101t013408: version is serving traffic. split: 100%
		20191031t211924: keeping the latest 5 versions (2)
		20191031t211903: keeping the latest 5 versions (3)
		20191031t205920: keeping the latest 5 versions (4)
		20191031t232247: keeping the latest 5 versions (5)
		20191031t232028: keeping recent versions (2h56m4.73591512s)
		20191031t220312: keeping recent versions (4h13m20.735921508s)
		20191031t211935: keeping recent versions (4h56m55.73592447s)
		20191031t211824: keeping recent versions (4h58m10.735928067s)
		20191031t200353: keeping recent versions (6h12m38.735932792s)
		20191031t150644: keeping recent versions (11h9m44.735935312s)
	]
	versions to delete (7): [
		20191030t225128: bye
		20191030t214823: bye
		20191030t214355: bye
		20191030t204338: bye
		20191030t202841: bye
		20191030t195403: bye
		20191030t192250: bye
	]
	deleting go-discovery/go-dev/20191030t225128
	...
*/
package main
