---
title: New unique package
date: 2024-08-27
by:
- Michael Knyszek
tags:
- interning
- unique
summary: New package for interning in Go 1.23.
---

The standard library of Go 1.23 now includes the [new `unique` package](https://pkg.go.dev/unique).
The purpose behind this package is to enable the canonicalization of
comparable values.
In other words, this package lets you deduplicate values so that they point
to a single, canonical, unique copy, while efficiently managing the canonical
copies under the hood.
You might be familiar with this concept already, called
["interning"](https://en.wikipedia.org/wiki/Interning_(computer_science)).
Let's dive in to see how it works, and why it's useful.

## A simple implementation of interning

At a high level, interning is very simple.
Take the code sample below, which deduplicates strings using just a regular
map.

```
var internPool map[string]string

// Intern returns a string that is equal to s but that may share storage with
// a string previously passed to Intern.
func Intern(s string) string {
	pooled, ok := internPool[s]
	if !ok {
		// Clone the string in case it's part of some much bigger string.
		// This should be rare, if interning is being used well.
		pooled = strings.Clone(s)
		internPool[pooled] = pooled
	}
	return pooled
}
```

This is useful for when you're constructing a lot of strings that are likely to
be duplicates, like when parsing a text format.

This implementation is super simple and works well enough for some cases, but it
has a few problems:

* It never removes strings from the pool.
* It cannot be safely used by multiple goroutines concurrently.
* It only works with strings, even though the idea is quite general.

There's also a missed opportunity in this implementation, and it's subtle.
Under the hood, [strings are immutable structures consisting of a pointer
and a length](/blog/slices).
When comparing two strings, if the pointers are not equal, then we must
compare their contents to determine equality.
But if we know that two strings are canonicalized, then it *is* sufficient
to just check their pointers.

## Enter the `unique` package

The new `unique` package introduces a function similar to `Intern` called
[`Make`](https://pkg.go.dev/unique#Make).

It works about the same way as `Intern`.
Internally there's also a global map ([a fast generic concurrent
map](https://pkg.go.dev/internal/concurrent@go1.23.0)) and `Make` looks up the
provided value in that map.
But it also differs from `Intern` in two important ways.
Firstly, it accepts values of any comparable type.
And secondly, it returns a wrapper value, a
[`Handle[T]`](https://pkg.go.dev/unique#Handle), from which the canonical value
can be retrieved.

This `Handle[T]` is key to the design.
A `Handle[T]` has the property that two `Handle[T]` values are equal if and
only if the values used to create them are equal.
What's more, the comparison of two `Handle[T]` values is cheap: it comes down
to a pointer comparison.
Compared to comparing two long strings, that's an order of magnitude cheaper!

So far, this is nothing you can't do in ordinary Go code.

But `Handle[T]` also has a second purpose: as long as a `Handle[T]` exists for
a value, the map will retain the canonical copy of the value.
Once all `Handle[T]` values that map to a specific value are gone, the
package marks that internal map entry as deletable, to be reclaimed in the near
future.
This sets a clear policy for when to remove entries from the map: when the
canonical entries are no longer being used, then the garbage collector is free
to clean them up.

If you've used Lisp before, this may all sound quite familiar to you.
Lisp [symbols](https://en.wikipedia.org/wiki/Symbol_(programming)) are interned
strings, but not strings themselves, and all symbols' string values are
guaranteed to be in the same pool.
This relationship between symbols and strings parallels the relationship
between `Handle[string]` and `string`.

## A real-world example

So, how might one use `unique.Make`?
Look no further than the `net/netip` package in the standard library, which
interns values of type `addrDetail`, part of the
[`netip.Addr`](https://pkg.go.dev/net/netip#Addr) structure.

Below is an abridged version of the actual code from `net/netip` that uses
`unique`.

```
// Addr represents an IPv4 or IPv6 address (with or without a scoped
// addressing zone), similar to net.IP or net.IPAddr.
type Addr struct {
	// Other irrelevant unexported fields...

	// Details about the address, wrapped up together and canonicalized.
	z unique.Handle[addrDetail]
}

// addrDetail indicates whether the address is IPv4 or IPv6, and if IPv6,
// specifies the zone name for the address.
type addrDetail struct {
	isV6   bool   // IPv4 is false, IPv6 is true.
	zoneV6 string // May be != "" if IsV6 is true.
}

var z6noz = unique.Make(addrDetail{isV6: true})

// WithZone returns an IP that's the same as ip but with the provided
// zone. If zone is empty, the zone is removed. If ip is an IPv4
// address, WithZone is a no-op and returns ip unchanged.
func (ip Addr) WithZone(zone string) Addr {
	if !ip.Is6() {
		return ip
	}
	if zone == "" {
		ip.z = z6noz
		return ip
	}
	ip.z = unique.Make(addrDetail{isV6: true, zoneV6: zone})
	return ip
}
```

Since many IP addresses are likely to use the same zone and this zone is part
of their identity, it makes a lot of sense to canonicalize them.
The deduplication of zones reduces the average memory footprint of each
`netip.Addr`, while the fact that they're canonicalized means `netip.Addr`
values are more efficient to compare, since comparing zone names becomes a
simple pointer comparison.

## A footnote about interning strings

While the `unique` package is useful, `Make` is admittedly not quite like
`Intern` for strings, since the `Handle[T]` is required to keep a string from
being deleted from the internal map.
This means you need to modify your code to retain handles as well as strings.

But strings are special in that, although they behave like values, they
actually contain pointers under the hood, as we mentioned earlier.
This means that we could potentially canonicalize just the underlying storage
of the string, hiding the details of a `Handle[T]` inside the string itself.
So, there is still a place in the future for what I'll call _transparent string
interning_, in which strings can be interned without the `Handle[T]` type,
similar to the `Intern` function but with semantics more closely resembling
`Make`.

In the meantime, `unique.Make("my string").Value()` is one possible workaround.
Even though failing to retain the handle will allow the string to be deleted
from `unique`'s internal map, map entries are not deleted immediately.
In practice, entries will not be deleted until at least the next garbage
collection completes, so this workaround still allows for some degree of
deduplication in the periods between collections.

## Some history, and looking toward the future

The truth is that the `net/netip` package actually interned zone strings since
it was first introduced.
The interning package it used was an internal copy of the
[go4.org/intern](https://pkg.go.dev/go4.org/intern) package.
Like the `unique` package, it has a `Value` type (which looks a lot like a
`Handle[T]`, pre-generics), has the notable property that entries in the
internal map are removed once their handles are no longer referenced.

But to achieve this behavior, it has to do some unsafe things.
In particular, it makes some assumptions about the garbage collector's behavior
to implement [_weak pointers_](https://en.wikipedia.org/wiki/Weak_reference)
outside the runtime.
A weak pointer is a pointer that doesn't prevent the garbage collector from
reclaiming a variable; when this happens, the pointer automatically becomes
nil.
As it happens, weak pointers are _also_ the core abstraction underlying the
`unique` package.

That's right: while implementing the `unique` package, we added proper weak
pointer support to the garbage collector.
And after stepping through the minefield of regrettable design decisions that
accompany weak pointers (like, should weak pointers track [object
resurrection](https://en.wikipedia.org/wiki/Object_resurrection)? No!), we were
astonished by how simple and straightforward all of it turned out to be.
Astonished enough that weak pointers are now a [public
proposal](/issue/67552).

This work also led us to reexamine finalizers, resulting in another proposal
for an easier-to-use and more efficient [replacement for
finalizers](/issue/67535).
With [a hash function for comparable values](/issue/54670) on the way as well,
the future of [building memory-efficient
caches](/issue/67552#issuecomment-2200755798) in Go is bright!
