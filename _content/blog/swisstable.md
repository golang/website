---
title: Faster Go maps with Swiss Tables
date: 2025-02-26
by:
- Michael Pratt
summary: Go 1.24 improves map performance with a brand new map implementation
---

The hash table is a central data structure in computer science, and it provides the implementation for the map type in many languages, including Go.

The concept of a hash table was [first described](https://spectrum.ieee.org/hans-peter-luhn-and-the-birth-of-the-hashing-algorithm) by Hans Peter Luhn in 1953 in an internal IBM memo that suggested speeding up search by placing items into "buckets" and using a linked list for overflow when buckets already contain an item.
Today we would call this a [hash table using chaining](https://en.wikipedia.org/wiki/Hash_table#Separate_chaining).

In 1954, Gene M. Amdahl, Elaine M. McGraw, and Arthur L. Samuel first used an "open addressing" scheme when programming the IBM 701.
When a bucket already contains an item, the new item is placed in the next empty bucket.
This idea was formalized and published in 1957 by W. Wesley Peterson in ["Addressing for Random-Access Storage"](https://ieeexplore.ieee.org/document/5392733).
Today we would call this a [hash table using open addressing with linear probing](https://en.wikipedia.org/wiki/Hash_table#Open_addressing).

With data structures that have been around this long, it's easy to think that they must be "done"; that we know everything there is to know about them and they can't be improved anymore.
That's not true!
Computer science research continues to make advancements in fundamental algorithms, both in terms of algorithmic complexity and taking advantage of modern CPU hardware.
For example, Go 1.19 [switched the `sort` package](/doc/go1.19#sortpkgsort) from a traditional quicksort, to [pattern-defeating quicksort](https://arxiv.org/pdf/2106.05123.pdf), a novel sorting algorithm from Orson R. L. Peters, first described in 2015.

Like sorting algorithms, hash table data structures continue to see improvements.
In 2017, Sam Benzaquen, Alkis Evlogimenos, Matt Kulukundis, and Roman Perepelitsa at Google presented [a new C++ hash table design](https://www.youtube.com/watch?v=ncHmEUmJZf4), dubbed "Swiss Tables".
In 2018, their implementation was [open sourced in the Abseil C++ library](https://abseil.io/blog/20180927-swisstables).

Go 1.24 includes a completely new implementation of the built-in map type, based on the Swiss Table design.
In this blog post we'll look at how Swiss Tables improve upon traditional hash tables, and at some of the unique challenges in bringing the Swiss Table design to Go's maps.

## Open-addressed hash table

Swiss Tables are a form of open-addressed hash table, so let's do a quick overview of how a basic open-addressed hash table works.

In an open-addressed hash table, all items are stored in a single backing array.
We'll call each location in the array a *slot*.
The slot to which a key belongs is primarily determined by the *hash function*, `hash(key)`.
The hash function maps each key to an integer, where the same key always maps to the same integer, and different keys ideally follow a uniform random distribution of integers.
The defining feature of open-addressed hash tables is that they resolve collisions by storing the key elsewhere in the backing array.
So, if the slot is already full (a *collision*), then a *probe sequence* is used to consider other slots until an empty slot is found.
Let's take a look at a sample hash table to see how this works.

### Example

Below you can see a 16-slot backing array for a hash table, and the key (if any) stored in each slot.
The values are not shown, as they are not relevant to this example.

<style>
/*
go.dev .Article max-width is 55em. Only enable horizontal scrolling if the
screen is narrow enough to require scrolling (narrower than article width)
because otherwise some platforms (e.g., Chrome on macOS) display a scrollbar
even when the screen is wide enough.
*/
@media screen and (max-width: 55em) {
    .swisstable-table-container {
        /* Scroll horizontally on overflow (likely on mobile) */
        overflow: scroll;
    }
}

.swisstable-table {
    /* Combine table inner borders (1px total rather than 2px, one for cell above and one for cell below. */
    border-collapse: collapse;
    /* All column widths equal. */
    table-layout: fixed;
    /* Center table within container div */
    margin: 0 auto;
}

.swisstable-table-cell {
    /* Black border between cells. */
    border: 1px solid;
    /* Add visual spacing around contents. */
    padding: 0.5em 1em 0.5em 1em;
    /* Center within cell. */
    text-align: center;
}
</style>

<div class="swisstable-table-container">
    <table class="swisstable-table">
        <thead>
            <tr>
                <th class="swisstable-table-cell">Slot</th>
                <th class="swisstable-table-cell">0</th>
                <th class="swisstable-table-cell">1</th>
                <th class="swisstable-table-cell">2</th>
                <th class="swisstable-table-cell">3</th>
                <th class="swisstable-table-cell">4</th>
                <th class="swisstable-table-cell">5</th>
                <th class="swisstable-table-cell">6</th>
                <th class="swisstable-table-cell">7</th>
                <th class="swisstable-table-cell">8</th>
                <th class="swisstable-table-cell">9</th>
                <th class="swisstable-table-cell">10</th>
                <th class="swisstable-table-cell">11</th>
                <th class="swisstable-table-cell">12</th>
                <th class="swisstable-table-cell">13</th>
                <th class="swisstable-table-cell">14</th>
                <th class="swisstable-table-cell">15</th>
            </tr>
        </thead>
        <tbody>
            <tr>
                <td class="swisstable-table-cell">Key</td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell">56</td>
                <td class="swisstable-table-cell">32</td>
                <td class="swisstable-table-cell">21</td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell">78</td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
            </tr>
        </tbody>
    </table>
</div>

To insert a new key, we use the hash function to select a slot.
Since there are only 16 slots, we need to restrict to this range, so we'll use `hash(key) % 16` as the target slot.
Suppose we want to insert key `98` and `hash(98) % 16 = 7`.
Slot 7 is empty, so we simply insert 98 there.
On the other hand, suppose we want to insert key `25` and `hash(25) % 16 = 3`.
Slot 3 is a collision because it already contains key 56.
Thus we cannot insert here.

We use a probe sequence to find another slot.
There are a variety of well-known probe sequences.
The original and most straightforward probe sequence is *linear probing*, which simply tries successive slots in order.

So, in the `hash(25) % 16 = 3` example, since slot 3 is in use, we would consider slot 4 next, which is also in use.
So too is slot 5.
Finally, we'd get to empty slot 6, where we'd store key 25.

Lookup follows the same approach.
A lookup of key 25 would start at slot 3, check whether it contains key 25 (it does not), and then continue linear probing until it finds key 25 in slot 6.

This example uses a backing array with 16 slots.
What happens if we insert more than 16 elements?
If the hash table runs out of space, it will grow, usually by doubling the size of the backing array.
All existing entries are reinserted into the new backing array.

Open-addressed hash tables don't actually wait until the backing array is completely full to grow because as the array gets more full, the average length of each probe sequence increases.
In the example above using key 25, we must probe 4 different slots to find an empty slot.
If the array had only one empty slot, the worst case probe length would be O(n).
That is, you may need to scan the entire array.
The proportion of used slots is called the *load factor*, and most hash tables define a *maximum load factor* (typically 70-90%) at which point they will grow to avoid the extremely long probe sequences of very full hash tables.

## Swiss Table

The Swiss Table [design](https://abseil.io/about/design/swisstables) is also a form of open-addressed hash table.
Let's see how it improves over a traditional open-addressed hash table.
We still have a single backing array for storage, but we will break the array into logical *groups* of 8 slots each.
(Larger group sizes are possible as well. More on that below.)

In addition, each group has a 64-bit *control word* for metadata.
Each of the 8 bytes in the control word corresponds to one of the slots in the group.
The value of each byte denotes whether that slot is empty, deleted, or in use.
If it is in use, the byte contains the lower 7 bits of the hash for that slot's key (called `h2`).

<!-- Group table followed by control word table. Both are in the same container so they scroll together on mobile. -->
<div class="swisstable-table-container">
    <table class="swisstable-table">
        <thead>
            <tr>
                <th class="swisstable-table-cell"></th>
                <th class="swisstable-table-cell" colspan="8">Group 0</th>
                <th class="swisstable-table-cell" colspan="8">Group 1</th>
            </tr>
            <tr>
                <th class="swisstable-table-cell">Slot</th>
                <th class="swisstable-table-cell">0</th>
                <th class="swisstable-table-cell">1</th>
                <th class="swisstable-table-cell">2</th>
                <th class="swisstable-table-cell">3</th>
                <th class="swisstable-table-cell">4</th>
                <th class="swisstable-table-cell">5</th>
                <th class="swisstable-table-cell">6</th>
                <th class="swisstable-table-cell">7</th>
                <th class="swisstable-table-cell">0</th>
                <th class="swisstable-table-cell">1</th>
                <th class="swisstable-table-cell">2</th>
                <th class="swisstable-table-cell">3</th>
                <th class="swisstable-table-cell">4</th>
                <th class="swisstable-table-cell">5</th>
                <th class="swisstable-table-cell">6</th>
                <th class="swisstable-table-cell">7</th>
            </tr>
        </thead>
        <tbody>
            <tr>
                <td class="swisstable-table-cell">Key</td>
                <td class="swisstable-table-cell">56</td>
                <td class="swisstable-table-cell">32</td>
                <td class="swisstable-table-cell">21</td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell">78</td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
            </tr>
        </tbody>
    </table>
    <br/> <!-- Visual space between the tables -->
    <table class="swisstable-table">
        <thead>
            <tr>
                <th class="swisstable-table-cell"></th>
                <th class="swisstable-table-cell" colspan="8">64-bit control word 0</th>
                <th class="swisstable-table-cell" colspan="8">64-bit control word 1</th>
            </tr>
            <tr>
                <th class="swisstable-table-cell">Slot</th>
                <th class="swisstable-table-cell">0</th>
                <th class="swisstable-table-cell">1</th>
                <th class="swisstable-table-cell">2</th>
                <th class="swisstable-table-cell">3</th>
                <th class="swisstable-table-cell">4</th>
                <th class="swisstable-table-cell">5</th>
                <th class="swisstable-table-cell">6</th>
                <th class="swisstable-table-cell">7</th>
                <th class="swisstable-table-cell">0</th>
                <th class="swisstable-table-cell">1</th>
                <th class="swisstable-table-cell">2</th>
                <th class="swisstable-table-cell">3</th>
                <th class="swisstable-table-cell">4</th>
                <th class="swisstable-table-cell">5</th>
                <th class="swisstable-table-cell">6</th>
                <th class="swisstable-table-cell">7</th>
            </tr>
        </thead>
        <tbody>
            <tr>
                <td class="swisstable-table-cell">h2</td>
                <td class="swisstable-table-cell">23</td>
                <td class="swisstable-table-cell">89</td>
                <td class="swisstable-table-cell">50</td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell">47</td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
                <td class="swisstable-table-cell"></td>
            </tr>
        </tbody>
    </table>
</div>

Insertion works as follows:

1. Compute `hash(key)` and break the hash into two parts: the upper 57-bits (called `h1`) and the lower 7 bits (called `h2`).
2. The upper bits (`h1`) are used to select the first group to consider: `h1 % 2` in this case, as there are only 2 groups.
3. Within a group, all slots are equally eligible to hold the key. We must first determine whether any slot already contains this key, in which case this is an update rather than a new insertion.
4. If no slot contains the key, then we look for an empty slot to place this key.
5. If no slot is empty, then we continue the probe sequence by searching the next group.

Lookup follows the same basic process.
If we find an empty slot in step 4, then we know an insertion would have used this slot and can stop the search.

Step 3 is where the Swiss Table magic happens.
We need to check whether any slot in a group contains the desired key.
Naively, we could just do a linear scan and compare all 8 keys.
However, the control word lets us do this more efficiently.
Each byte contains the lower 7 bits of the hash (`h2`) for that slot.
If we determine which bytes of the control word contain the `h2` we are looking for, we'll have a set of candidate matches.

In other words, we want to do a byte-by-byte equality comparison within the control word.
For example, if we are looking for key 32, where `h2 = 89`, the operation we want looks like so.

<!-- Visualization of SIMD comparison -->
<div class="swisstable-table-container">
    <table class="swisstable-table">
        <tbody>
            <tr>
                <td class="swisstable-table-cell"><strong>Test word</strong></td>
                <td class="swisstable-table-cell">89</td>
                <td class="swisstable-table-cell">89</td>
                <td class="swisstable-table-cell">89</td>
                <td class="swisstable-table-cell">89</td>
                <td class="swisstable-table-cell">89</td>
                <td class="swisstable-table-cell">89</td>
                <td class="swisstable-table-cell">89</td>
                <td class="swisstable-table-cell">89</td>
            </tr>
            <tr>
                <td class="swisstable-table-cell"><strong>Comparison</strong></td>
                <td class="swisstable-table-cell">==</td>
                <td class="swisstable-table-cell">==</td>
                <td class="swisstable-table-cell">==</td>
                <td class="swisstable-table-cell">==</td>
                <td class="swisstable-table-cell">==</td>
                <td class="swisstable-table-cell">==</td>
                <td class="swisstable-table-cell">==</td>
                <td class="swisstable-table-cell">==</td>
            </tr>
            <tr>
                <td class="swisstable-table-cell"><strong>Control word</strong></td>
                <td class="swisstable-table-cell">23</td>
                <td class="swisstable-table-cell">89</td>
                <td class="swisstable-table-cell">50</td>
                <td class="swisstable-table-cell">-</td>
                <td class="swisstable-table-cell">-</td>
                <td class="swisstable-table-cell">-</td>
                <td class="swisstable-table-cell">-</td>
                <td class="swisstable-table-cell">-</td>
            </tr>
            <tr>
                <td class="swisstable-table-cell"><strong>Result</strong></td>
                <td class="swisstable-table-cell">0</td>
                <td class="swisstable-table-cell">1</td>
                <td class="swisstable-table-cell">0</td>
                <td class="swisstable-table-cell">0</td>
                <td class="swisstable-table-cell">0</td>
                <td class="swisstable-table-cell">0</td>
                <td class="swisstable-table-cell">0</td>
                <td class="swisstable-table-cell">0</td>
            </tr>
        </tbody>
    </table>
</div>

This is an operation supported by [SIMD](https://en.wikipedia.org/wiki/Single_instruction,_multiple_data) hardware, where a single instruction performs parallel operations on independent values within a larger value (*vector*). In this case, we [can implement this operation](https://cs.opensource.google/go/go/+/master:src/internal/runtime/maps/group.go;drc=a08984bc8f2acacebeeadf7445ecfb67b7e7d7b1;l=155?ss=go) using a set of standard arithmetic and bitwise operations when special SIMD hardware is not available.

The result is a set of candidate slots.
Slots where `h2` does not match do not have a matching key, so they can be skipped.
Slots where `h2` does match are potential matches, but we must still check the entire key, as there is potential for collisions (1/128 probability of collision with a 7-bit hash, so still quite low).

This operation is very powerful, as we have effectively performed 8 steps of a probe sequence at once, in parallel.
This speeds up lookup and insertion by reducing the average number of comparisons we need to perform.
This improvement to probing behavior allowed both the Abseil and Go implementations to increase the maximum load factor of Swiss Table maps compared to prior maps, which lowers the average memory footprint.

## Go challenges

Go's built-in map type has some unusual properties that pose additional challenges to adopting a new map design.
Two were particularly tricky to deal with.

### Incremental growth

When a hash table reaches its maximum load factor, it needs to grow the backing array.
Typically this means the next insertion doubles the size of the array, and copies all entries to the new array.
Imagine inserting into a map with 1GB of entries.
Most insertions are very fast, but the one insertion that needs to grow the map from 1GB to 2GB will need to copy 1GB of entries, which will take a long time.

Go is frequently used for latency-sensitive servers, so we don't want operations on built-in types that can have arbitrarily large impact on tail latency.
Instead, Go maps grow incrementally, so that each insertion has an upper bound on the amount of growth work it must do.
This bounds the latency impact of a single map insertion.

Unfortunately, the Abseil (C++) Swiss Table design assumes all at once growth, and the probe sequence depends on the total group count, making it difficult to break up.

Go's built-in map addresses this with another layer of indirection by splitting each map into multiple Swiss Tables.
Rather than a single Swiss Table implementing the entire map, each map consists of one or more independent tables that cover a subset of the key space.
An individual table stores a maximum of 1024 entries.
A variable number of upper bits in the hash are used to select which table a key belongs to.
This is a form of [*extendible hashing*](https://en.wikipedia.org/wiki/Extendible_hashing), where the number of bits used increases as needed to differentiate the total number of tables.

During insertion, if an individual table needs to grow, it will do so all at once, but other tables are unaffected.
Thus the upper bound for a single insertion is the latency of growing a 1024-entry table into two 1024-entry tables, copying 1024 entries.

### Modification during iteration

Many hash table designs, including Abseil's Swiss Tables, forbid modifying the map during iteration.
The Go language spec [explicitly allows](/ref/spec#For_statements:~:text=The%20iteration%20order,iterations%20is%200.) modifications during iteration, with the following semantics:

* If an entry is deleted before it is reached, it will not be produced.
* If an entry is updated before it is reached, the updated value will be produced.
* If a new entry is added, it may or may not be produced.

A typical approach to hash table iteration is to simply walk through the backing array and produce values in the order they are laid out in memory.
This approach runs afoul of the above semantics, most notably because insertions may grow the map, which would shuffle the memory layout.

We can avoid the impact of shuffle during growth by having the iterator keep a reference to the table it is currently iterating over.
If that table grows during iteration, we keep using the old version of the table and thus continue to deliver keys in the order of the old memory layout.

Does this work with the above semantics?
New entries added after growth will be missed entirely, as they are only added to the grown table, not the old table.
That's fine, as the semantics allow new entries not to be produced.
Updates and deletions are a problem, though: using the old table could produce stale or deleted entries.

This edge case is addressed by using the old table only to determine the iteration order.
Before actually returning the entry, we consult the grown table to determine whether the entry still exists, and to retrieve the latest value.

This covers all the core semantics, though there are even more small edge cases not covered here.
Ultimately, the permissiveness of Go maps with iteration results in iteration being the most complex part of Go's map implementation.

## Future work

In [microbenchmarks](/issue/54766#issuecomment-2542444404), map operations are up to 60% faster than in Go 1.23.
Exact performance improvement varies quite a bit due to the wide variety of operations and uses of maps, and some edge cases do regress compared to Go 1.23.
Overall, in full application benchmarks, we found a geometric mean CPU time improvement of around 1.5%.

There are more map improvements we want to investigate for future Go releases.
For example, we may be able to [increase the locality of](/issue/70835) operations on maps that are not in the CPU cache.

We could also further improve the control word comparisons.
As described above, we have a portable implementation using standard arithmetic and bitwise operations.
However, some architectures have SIMD instructions that perform this sort of comparison directly.
Go 1.24 already uses 8-byte SIMD instructions for amd64, but we could extend support to other architectures.
More importantly, while standard instructions operate on up to 8-byte words, SIMD instructions nearly always support at least 16-byte words.
This means we could increase the group size to 16 slots, and perform 16 hash comparisons in parallel instead of 8.
This would further decrease the average number of probes required for lookups.

## Acknowledgements

A Swiss Table-based Go map implementation has been a long time coming, and involved many contributors.
I want to thank YunHao Zhang ([@zhangyunhao116](https://github.com/zhangyunhao116)), PJ Malloy ([@thepudds](https://github.com/thepudds)), and [@andy-wm-arthur](https://github.com/andy-wm-arthur) for building initial versions of a Go Swiss Table implementation.
Peter Mattis ([@petermattis](https://github.com/petermattis)) combined these ideas with solutions to the Go challenges above to build [`github.com/cockroachdb/swiss`](https://pkg.go.dev/github.com/cockroachdb/swiss), a Go-spec compliant Swiss Table implementation.
The Go 1.24 built-in map implementation is heavily based on Peter's work.
Thank you to everyone in the community that contributed!
