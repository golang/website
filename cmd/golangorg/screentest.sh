#!/bin/bash
# Copyright 2024 The Go Authors. All rights reserved.
# Use of this source code is governed by a BSD-style
# license that can be found in the LICENSE file.

go run ./cmd/screentest \
	-test http://localhost:6060/go.dev \
	-want https://go.dev \
	'./cmd/golangorg/testdata/screentest/*'
