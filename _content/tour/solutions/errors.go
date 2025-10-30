//go:build OMIT

// Copyright 2012 The Go Authors.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"fmt"
	"math"
)

type ErrNegativeSqrt float64

func (e ErrNegativeSqrt) Error() string {
	return fmt.Sprintf("Sqrt: negative number %g", e)
}

const delta = 1e-10

func Sqrt(f float64) (float64, error) {
	if f < 0 {
		return 0, ErrNegativeSqrt(f)
	}
	z := f
	for {
		n := z - (z*z-f)/(2*z)
		if math.Abs(n-z) < delta {
			break
		}
		z = n
	}
	return z, nil
}

func main() {
	for _, number := range []float64{2, -2} {
		if root, err := Sqrt(number); err == nil {
			fmt.Println(root)
		} else {
			fmt.Println(err)
		}
	}
}
