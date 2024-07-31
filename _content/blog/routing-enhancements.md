---
title: Routing Enhancements for Go 1.22
date: 2024-02-13
by:
- Jonathan Amsterdam, on behalf of the Go team
summary: Go 1.22's additions to patterns for HTTP routes.
---

Go 1.22 brings two enhancements to the `net/http` package's router: method
matching and wildcards. These features let you express common routes as
patterns instead of Go code. Although they are simple to explain and use,
it was a challenge to come up with the right rules for selecting the winning
pattern when several match a request.

We made these changes as part of our continuing effort to make Go a great
language for building production systems. We studied many third-party web
frameworks, extracted what we felt were the most used features, and integrated
them into `net/http`. Then we validated our choices and improved our design by
collaborating with the community in a [GitHub discussion](
https://github.com/golang/go/discussions/60227) and a [proposal issue](/issue/61410).
Adding these features to the standard library means one fewer dependency for
many projects. But third-party web frameworks remain a fine choice for current
users or programs with advanced routing needs.

## Enhancements

The new routing features almost exclusively affect the pattern string passed
to the two `net/http.ServeMux` methods `Handle` and `HandleFunc`, and the
corresponding top-level functions `http.Handle` and `http.HandleFunc`. The only
API changes are two new methods on `net/http.Request` for working with wildcard
matches.

We'll illustrate the changes with a hypothetical blog server in which every post
has an integer identifier. A request like `GET /posts/234` retrieves the post with
ID 234. Before Go 1.22, the code for handling those requests would start with a
line like this:

    http.HandleFunc("/posts/", handlePost)

The trailing slash routes all requests beginning `/posts/` to the `handlePost`
function, which would have to check that the HTTP method was GET, extract
the identifier, and retrieve the post. Since the method check isn't strictly
necessary to satisfy the request, it would be a natural mistake to omit it. That
would mean that a request like `DELETE /posts/234` would fetch the post, which
is surprising at the least.

In Go 1.22, the existing code will continue to work, or you could instead write this:

    http.HandleFunc("GET /posts/{id}", handlePost2)

This pattern matches a GET request whose path begins "/posts/" and has two
segments. (As a special case, GET also matches HEAD; all the other methods match
exactly.) The `handlePost2` function no longer needs to check the method, and
extracting the identifier string can be written using the new `PathValue` method
on `Request`:

    idString := req.PathValue("id")

The rest of `handlePost2` would behave like `handlePost`, converting the string
identifier to an integer and fetching the post.

Requests like `DELETE /posts/234` will fail if no other matching pattern is
registered. In accordance with [HTTP semantics](
https://httpwg.org/specs/rfc9110.html#status.405), a `net/http` server will reply
to such a request with a `405 Method Not Allowed` error that lists the available methods
in an `Allow` header.

A wildcard can match an entire segment, like `{id}` in the example above, or if
it ends in `...` it can match all the remaining segments of the path, as in the
pattern `/files/{pathname...}`.

There is one last bit of syntax. As we showed above, patterns ending in a slash,
like `/posts/`, match all paths beginning with that string. To match only the
path with the trailing slash, you can write `/posts/{$}`. That will match
`/posts/` but not `/posts` or `/posts/234`.

And there is one last bit of API: `net/http.Request` has a `SetPathValue` method
so that routers outside the standard library can make the results of their own
path parsing available via `Request.PathValue`.

## Precedence

Every HTTP router must deal with overlapping patterns, like `/posts/{id}` and
`/posts/latest`. Both of these patterns match the path "posts/latest", but at most
one can serve the request. Which pattern takes precedence?

Some routers disallow overlaps; others use the pattern that was registered last.
Go has always allowed overlaps, and has chosen the longer pattern regardless
of registration order. Preserving order-independence was important to us (and
necessary for backwards compatibility), but we needed a better rule than
"longest wins." That rule would select `/posts/latest` over `/posts/{id}`, but
would choose `/posts/{identifier}` over both. That seems wrong: the wildcard
name shouldn't matter. It feels like `/posts/latest` should always win this
competition, because it matches a single path instead of many.

Our quest for a good precedence rule led us to consider many properties of
patterns. For example, we considered preferring the pattern with the longest
literal (non-wildcard) prefix. That would choose `/posts/latest` over `/posts/
{id}`. But it wouldn't distinguish between `/users/{u}/posts/latest` and
`/users/{u}/posts/{id}`, and it seems like the former should take precedence.

We eventually chose a rule based on what the patterns mean instead of how they
look. Every valid pattern matches a set of requests. For example,
`/posts/latest` matches requests with the path `/posts/latest`, while `/posts/{id}`
matches requests with any two-segment path whose first segment is "posts". We
say that one pattern is _more specific_ than another if it matches a strict subset
of requests. The pattern `/posts/latest` is more specific than `/posts/{id}`
because the latter matches every request that the former does, and more.

The precedence rule is simple: the most specific pattern wins. This rule
matches our intuition that `posts/latests` should be preferred to `posts/{id}`,
and `/users/{u}/posts/latest` should be preferred to `/users/{u}/posts/{id}`.
It also makes sense for methods. For example, `GET /posts/{id}` takes
precedence over `/posts/{id}` because the first only matches GET and HEAD
requests, while the second matches requests with any method.

The "most specific wins" rule generalizes the original "longest wins" rule for
the path parts of original patterns, those without wildcards or `{$}`. Such
patterns only overlap when one is a prefix of the other, and the longer is the
more specific.

What if two patterns overlap but neither is more specific? For example, `/posts/{id}`
and `/{resource}/latest` both match `/posts/latest`. There is no obvious answer to
which takes precedence, so we consider these patterns to conflict with each other.
Registering both of them (in either order!) will panic.

The precedence rule works exactly as above for methods and paths, but we had to
make one exception for hosts to preserve compatibility: if two patterns would
otherwise conflict and one has a host while the other does not, then the pattern
with the host takes precedence.

Students of computer science may recall the beautiful theory of regular
expressions and regular languages. Each regular expression picks out a regular
language, the set of strings matched by the expression. Some questions are
easier to pose and answer by talking about languages rather than expressions.
Our precedence rule was inspired by this theory. Indeed, each routing pattern
corresponds to a regular expression, and sets of matching requests play the role of
regular languages.

Defining precedence by languages instead of expressions makes it easy to state
and understand. But there is a downside to having a rule based on potentially
infinite sets: it isn't clear how to implement it efficiently. It turns out we
can determine whether two patterns conflict by walking them segment by segment.
Roughly speaking, if one pattern has a literal segment wherever the other has a
wildcard, it is more specific; but if literals align with wildcards in both
directions, the patterns conflict.

As new patterns are registered on a `ServeMux`, it checks for conflicts with previously
registered patterns. But checking every pair of patterns would take quadratic
time. We use an index to skip patterns that cannot conflict with a new pattern;
in practice, it works quite well. In any case, this check happens when
patterns are registered, usually at server startup. The time to match incoming
requests in Go 1.22 hasn't changed much from previous versions.

## Compatibility

We made every effort to keep the new functionality compatible with older
versions of Go. The new pattern syntax is a superset of the old, and the new
precedence rule generalizes the old one. But there are a few edge cases. For
example, previous versions of Go accepted patterns with braces and treated
them literally, but Go 1.22 uses braces for wildcards. The GODEBUG setting
`httpmuxgo121` restores the old behavior.

For more details about these routing enhancements, see the [`net/http.ServeMux`
documentation](/pkg/net/http#ServeMux).


