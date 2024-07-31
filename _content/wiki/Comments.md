---
title: Comments
---

<!--
This is just a placeholder page for enabling a test.
In the deployed site it is overwritten with the content of go.googlesource.com/wiki.
-->

Every package should have a package comment. It should immediately precede the ` package ` statement in one of the files in the package. (It only needs to appear in one file.) It should begin with a single sentence that begins "Package _packagename_" and give a concise summary of the package functionality. This introductory sentence will be used in godoc's list of all packages.

Subsequent sentences and/or paragraphs can give more details. Sentences should be properly punctuated.

```go
// Package superman implements methods for saving the world.
//
// Experience has shown that a small number of procedures can prove
// helpful when attempting to save the world.
package superman
```

Nearly every top-level type, const, var and func should have a comment. A comment for bar should be in the form "_bar_ floats on high o'er vales and hills.". The first letter of _bar_ should not be capitalized unless it's capitalized in the code.

```go
// enterOrbit causes Superman to fly into low Earth orbit, a position
// that presents several possibilities for planet salvation.
func enterOrbit() os.Error {
  ...
}
```

All text that you indent inside a comment, godoc will render as a pre-formatted block. This facilitates code samples.

```go
// fight can be used on any enemy and returns whether Superman won.
//
// Examples:
//
//  fight("a random potato")
//  fight(LexLuthor{})
//
func fight(enemy interface{}) bool {
	// This is testing proper escaping in the wiki.
	for i := 0; i < 10; i++ {
		println("fight!")
	}
}
```


