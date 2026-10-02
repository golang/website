#!/bin/bash

# Copyright 2011 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

# This code parses mapping.old and finds a correspondence from the old
# urls (e.g. #42) to the corresponding path (e.g. /concurrency/3).

find_url() {
    title="$1"
    file=$(grep -l "* $title$" ./*.article)
    if [[ -z $file ]]; then
        echo "undefined" >&2
        return 1
    fi
    titles=$(grep "^* " "$file" | awk '{print NR, $0}')
    page=$(echo "$titles" | grep "* $title$" | awk '{print $1}')
    if [[ $(echo "$page" | wc -l) -gt "1" ]]; then
        echo "multiple matches found for $title; find 'CHOOSE BETWEEN' in the output" >&2
        page="CHOOSE BETWEEN $page"
    fi
    page=$(echo "$page")
    lesson=$(echo "$file" | sed -E 's|^(.*)/(.*).article$|\2|')
    echo "'/$lesson/$page'"
    return 0
}

mapping=$(cat mapping.old)

cd ../content || exit
while read -r page; do
    num=$(echo "$page" | awk '{print $1}')
    title=$(echo "$page" | cut -d' ' -f2-)
    url=$(find_url "$title")
    echo "    '#$num': $url, // $title"
done <<< "$mapping"
cd - > /dev/null || exit
