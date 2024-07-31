// Copyright 2024 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
)

type Benchmark struct {
	Name string
	Data []float64
}

var perf = []Benchmark{
	{Name: "Go 1: Uint64", Data: []float64{2.29656e+00, 3.23495e+00, 2.67690e+00, 3.84810e+00, 2.51320e+00, 1.87559e+00, 4.83222e+00}},
	{Name: "PCG: Uint64", Data: []float64{1.52985e+00, 8.03210e+00, 2.46820e+00, 1.13633e+01, 4.19290e+00, 2.21153e+00, 6.86715e+00}},
	{Name: "ChaCha8: Uint64", Data: []float64{3.10785e+00, 5.99320e+00, 4.22015e+00, 8.54389e+00, 4.64550e+00, 3.54895e+00, 7.48717e+00}},
	{Name: "Go 1: N(1000)", Data: []float64{3.04079e+00, 1.43600e+01, 4.84235e+00, 2.50285e+01, 3.05638e+00, 2.41394e+00, 1.27050e+01}},
	{Name: "PCG: N(1000)", Data: []float64{2.47375e+00, 1.03700e+01, 3.88155e+00, 1.49215e+01, 4.09644e+00, 2.31600e+00, 9.82755e+00}},
	{Name: "ChaCha8: N(1000)", Data: []float64{4.03063e+00, 8.37230e+00, 5.74820e+00, 1.17947e+01, 4.85670e+00, 3.99968e+00, 1.01725e+01}},
}

var columns = []string{
	"amd",
	"amd32",
	"intel",
	"intel32",
	"m1",
	"m3",
	"taut2a",
}

var descs = []string{
	"AMD Ryzen 9 7950X",
	"AMD Ryzen 9 7950X running 32-bit code",
	"11th Gen Intel Core i7-1185G7",
	"11th Gen Intel Core i7-1185G7 running 32-bit code",
	"Apple M1",
	"Apple M3",
	"Google Cloud Tau T2A (Ampere Altra)",
}

var svghdr = `<?xml version="1.0" standalone="no"?>
<!DOCTYPE svg PUBLIC "-//W3C//DTD SVG 1.1//EN"
  "http://www.w3.org/Graphics/SVG/1.1/DTD/svg11.dtd">
<svg height="170" width="400" version="1.1"
     xmlns="http://www.w3.org/2000/svg">
  <defs>
    <style type="text/css"><![CDATA[
      text {
        font-family: sans-serif, Arial;
        font-size: 12px;
      }
      text.head {
        font-weight: bold;
        font-size: 14px;
      }
    ]]></style>
  </defs>
`

func writeSVG(col int) {
	var buf bytes.Buffer
	buf.WriteString(svghdr)

	y := 5
	height := 20
	skip := 20
	scale := 50.
	for _, bench := range perf {
		if true || bench.Data[col] > 10 {
			scale = 25
			break
		}
	}
	fills := []string{
		"#ffaaaa",
		"#ccccff",
		"#ffffaa",
		"#ffaaaa",
		"#ccccff",
		"#ffffaa",
	}
	y += skip
	fmt.Fprintf(&buf, "<text x='%d' y='%d' class='head'><tspan dx='5' dy='-0.5em'>%s</tspan></text>\n",
		0, y, descs[col])
	for i, bench := range perf {
		val := bench.Data[col]
		fill := fills[i]
		barx := int(val * scale)
		fmt.Fprintf(&buf, "<rect x='5' y='%d' height='%d' width='%d' fill='%s' stroke='black' />\n", y+(skip-height)/2, height, barx, fill)
		labelx := 5 + barx
		if labelx > 395 {
			labelx = 395
		}
		fmt.Fprintf(&buf, "<text x='%d' y='%d' text-anchor='end'><tspan dy='-0.5em'>%.1fns</tspan></text>\n",
			labelx-3, y+skip-(skip-height)/2, val)
		textx := 10
		if labelx < 130 {
			textx = labelx + 5
		}
		fmt.Fprintf(&buf, "<text x='%d' y='%d'><tspan dy='-0.5em'>%s</tspan></text>\n",
			textx, y+skip-(skip-height)/2, bench.Name)
		y += skip
		if i == 2 {
			y += skip / 2
		}
	}
	y += 5

	buf.WriteString("</svg>\n")
	if err := os.WriteFile(columns[col]+".svg", buf.Bytes(), 0666); err != nil {
		log.Fatal(err)
	}
}

func main() {
	for i := range perf[0].Data {
		writeSVG(i)
	}
}
