---
title: Robust generic functions on slices
date: 2024-02-22
by:
- Valentin Deleplace
summary: Avoiding memory leaks in the slices package.
---

The [slices](/pkg/slices) package provides functions that work for slices of any type.
In this blog post we'll discuss how you can use these functions more effectively by understanding how slices are represented in memory and how that affects the garbage collector, and we'll cover how we recently adjusted these functions to make them less surprising.

With [Type parameters](/blog/deconstructing-type-parameters) we can write functions like [slices.Index](/pkg/slices#Index) once for all types of slices of comparable elements:

```
// Index returns the index of the first occurrence of v in s,
// or -1 if not present.
func Index[S ~[]E, E comparable](s S, v E) int {
	for i := range s {
		if v == s[i] {
			return i
		}
	}
	return -1
}
```

It is no longer necessary to implement `Index` again for each different type of element.

The [slices](/pkg/slices) package contains many such helpers to perform common operations on slices:

```
	s := []string{"Bat", "Fox", "Owl", "Fox"}
	s2 := slices.Clone(s)
	slices.Sort(s2)
	fmt.Println(s2) // [Bat Fox Fox Owl]
	s2 = slices.Compact(s2)
	fmt.Println(s2)                  // [Bat Fox Owl]
	fmt.Println(slices.Equal(s, s2)) // false
```

Several new functions (`Insert`, `Replace`, `Delete`, etc.) modify the slice. To understand how they work, and how to properly use them, we need to examine the underlying structure of slices.

A slice is a view of a portion of an array. [Internally](/blog/slices-intro), a slice contains a pointer, a length, and a capacity. Two slices can have the same underlying array, and can view overlapping portions.

For example, this slice `s` is a view on 4 elements of an array of size 6:

{{image "generic-slice-functions/1_sample_slice_4_6.svg" 450}}

If a function changes the length of a slice passed as a parameter, then it needs to return a new slice to the caller. The underlying array may remain the same if it doesn't have to grow. This explains why [append](/blog/slices) and `slices.Compact` return a value, but `slices.Sort`, which merely reorders the elements, does not.

Consider the task of deleting a portion of a slice. Prior to generics, the standard way to delete the portion `s[2:5]` from the slice `s` was to call the [append](/ref/spec#Appending_and_copying_slices) function to copy the end portion over the middle portion:

```
s = append(s[:2], s[5:]...)
```

The syntax was complex and error-prone, involving subslices and a variadic parameter. We added [slices.Delete](/pkg/slices#Delete) to make it easier to delete elements:

```
func Delete[S ~[]E, E any](s S, i, j int) S {
       return append(s[:i], s[j:]...)
}
```

The one-line function `Delete` more clearly expresses the programmer's intent. Let’s consider a slice `s` of length 6 and capacity 8, containing pointers:

{{image "generic-slice-functions/2_sample_slice_6_8.svg" 600}}

This call deletes the elements at `s[2]`, `s[3]`, `s[4]`  from the slice `s`:

```
s = slices.Delete(s, 2, 5)
```

{{image "generic-slice-functions/3_delete_s_2_5.svg" 600}}

The gap at the indices 2, 3, 4 is filled by shifting the element `s[5]` to the left, and setting the new length to `3`.

`Delete` need not allocate a new array, as it shifts the elements in place. Like `append`, it returns a new slice. Many other functions in the `slices` package follow this pattern, including `Compact`, `CompactFunc`, `DeleteFunc`, `Grow`, `Insert`, and `Replace`.

When calling these functions we must consider the original slice invalid, because the underlying array has been modified. It would be a mistake to call the function but ignore the return value:

```
	slices.Delete(s, 2, 5) // incorrect!
	// s still has the same length, but modified contents
```

## A problem of unwanted liveness

Before Go 1.22, `slices.Delete` didn't modify the elements between the new and original lengths of the slice. While the returned slice wouldn't include these elements, the "gap" created at the end of the original, now-invalidated slice continued to hold onto them. These elements could contain pointers to large objects (a 20MB image), and the garbage collector would not release the memory associated with these objects. This resulted in a memory leak that could lead to significant performance issues.

In this above example, we’re successfully deleting the pointers `p2`, `p3`, `p4` from `s[2:5]`, by shifting one element to the left. But `p3` and `p4` are still present in the underlying array, beyond the new length of `s`. The garbage collector won’t reclaim them. Less obviously, `p5` is not one of the deleted elements, but its memory may still leak because of the `p5` pointer kept in the gray part of the array.

This could be confusing for developers, if they were not aware that "invisible" elements were still using memory.

So we had two options:

* Either keep the efficient implementation of `Delete`. Let users set obsolete pointers to `nil` themselves, if they want to make sure the values pointed to can be freed.
* Or change `Delete` to always set the obsolete elements to zero. This is extra work, making `Delete` slightly less efficient. Zeroing pointers (setting them to `nil`) enables the garbage collection of the objects, when they become otherwise unreachable.

It was not obvious which option was best. The first one provided performance by default, and the second one provided memory frugality by default.

## The fix

A key observation is that "setting the obsolete pointers to `nil`" is not as easy as it seems. In fact, this task is so error-prone that we should not put the burden on the user to write it. Out of pragmatism, we chose to modify the implementation of the five functions `Compact`, `CompactFunc`, `Delete`, `DeleteFunc`, `Replace` to "clear the tail". As a nice side effect, the cognitive load is reduced and users now don’t need to worry about these memory leaks.

In Go 1.22, this is what the memory looks like after calling Delete:

{{image "generic-slice-functions/4_delete_s_2_5_nil.svg" 600}}

The code changed in the five functions uses the new built-in function [clear](/pkg/builtin#clear) (Go 1.21) to set the obsolete elements to the zero value of the element type of `s`:

{{image "generic-slice-functions/5_Delete_diff.png" 800}}

The zero value of `E` is `nil` when `E` is a type of pointer, slice, map, chan, or interface.

## Tests failing

This change has led to some tests that passed in Go 1.21 now failing in Go 1.22, when the slices functions are used incorrectly. This is good news. When you have a bug, tests should let you know.

If you ignore the return value of `Delete`:

```
slices.Delete(s, 2, 3)  // !! INCORRECT !!
```

then you may incorrectly assume that `s` does not contain any nil pointer. [Example in the Go Playground](/play/p/NDHuO8vINHv).

If you ignore the return value of `Compact`:

```
slices.Sort(s) // correct
slices.Compact(s) // !! INCORRECT !!
```

then you may incorrectly assume that `s` is properly sorted and compacted. [Example](/play/p/eFQIekiwlnu).

If you assign the return value of `Delete` to another variable, and keep using the original slice:

```
u := slices.Delete(s, 2, 3)  // !! INCORRECT, if you keep using s !!
```

then you may incorrectly assume that `s` does not contain any nil pointer. [Example](/play/p/rDxWmJpLOVO).

If you accidentally shadow the slice variable, and keep using the original slice:

```
s := slices.Delete(s, 2, 3)  // !! INCORRECT, using := instead of = !!
```

then you may incorrectly assume that `s` does not contain any nil pointer. [Example](/play/p/KSpVpkX8sOi).


## Conclusion

The API of the `slices` package is a net improvement over the traditional pre-generics syntax to delete or insert elements.

We encourage developers to use the new functions, while avoiding the "gotchas" listed above.

Thanks to the recent changes in the implementation, a class of memory leaks is automatically avoided, without any change to the API, and with no extra work for the developers.


## Further reading

The signature of the functions in the `slices` package is heavily influenced by the specifics of the representation of slices in memory. We recommend reading

*   [Go Slices: usage and internals](/blog/slices-intro)
*   [Arrays, slices: The mechanics of 'append'](/blog/slices)
*   The [dynamic array](https://en.wikipedia.org/wiki/Dynamic_array) data structure
*   The [documentation](/pkg/slices) of the package slices

The [original proposal](/issue/63393) about zeroing obsolete elements contains many details and comments.
