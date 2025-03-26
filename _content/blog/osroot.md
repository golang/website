---
title: Traversal-resistant file APIs
date: 2025-03-12
by:
- Damien Neil
tags:
- file
- os
summary: New file access APIs in Go 1.24.
---

A *path traversal vulnerability* arises when an attacker can trick a program
into opening a file other than the one it intended.
This post explains this class of vulnerability,
some existing defenses against it, and describes how the new
[`os.Root`](/pkg/os#Root) API added in Go 1.24 provides
a simple and robust defense against unintentional path traversal.

## Path traversal attacks

"Path traversal" covers a number of related attacks following a common pattern:
A program attempts to open a file in some known location, but an attacker causes
it to open a file in a different location.

If the attacker controls part of the filename, they may be able to use relative
directory components ("..") to escape the intended location:

    f, err := os.Open(filepath.Join(trustedLocation, "../../../../etc/passwd"))

On Windows systems, some names have special meaning:

    // f will print to the console.
    f, err := os.Create(filepath.Join(trustedLocation, "CONOUT$"))

If the attacker controls part of the local filesystem, they may be able to use
symbolic links to cause a program to access the wrong file:

    // Attacker links /home/user/.config to /home/otheruser/.config:
    err := os.WriteFile("/home/user/.config/foo", config, 0o666)

If the program defends against symlink traversal by first verifying that the intended file
does not contain any symlinks, it may still be vulnerable to
[time-of-check/time-of-use (TOCTOU) races](https://en.wikipedia.org/wiki/Time-of-check_to_time-of-use),
where the attacker creates a symlink after the program's check:

    // Validate the path before use.
    cleaned, err := filepath.EvalSymlinks(unsafePath)
    if err != nil {
      return err
    }
    if !filepath.IsLocal(cleaned) {
      return errors.New("unsafe path")
    }

    // Attacker replaces part of the path with a symlink.
    // The Open call follows the symlink:
    f, err := os.Open(cleaned)

Another variety of TOCTOU race involves moving a directory that forms part of a path
mid-traversal. For example, the attacker provides a path such as "a/b/c/../../etc/passwd",
and renames "a/b/c" to "a/b" while the open operation is in progress.

## Path sanitization

Before we tackle path traversal attacks in general, let's start with path sanitization.
When a program's threat model does not include attackers with access to the local file system,
it can be sufficient to validate untrusted input paths before use.

Unfortunately, sanitizing paths can be surprisingly tricky,
especially for portable programs that must handle both Unix and Windows paths.
For example, on Windows ```filepath.IsAbs(`\foo`)``` reports `false`,
because the path "\foo" is relative to the current drive.

In Go 1.20, we added the [`path/filepath.IsLocal`](/pkg/path/filepath#IsLocal)
function, which reports whether a path is "local". A "local" path is one which:

  - does not escape the directory in which it is evaluated ("../etc/passwd" is not allowed);
  - is not an absolute path ("/etc/passwd" is not allowed);
  - is not empty ("" is not allowed);
  - on Windows, is not a reserved name ("COM1" is not allowed).

In Go 1.23, we added the [`path/filepath.Localize`](/pkg/path/filepath#Localize)
function, which converts a /-separated path into a local operating system path.

Programs that accept and operate on potentially attacker-controlled paths should almost
always use `filepath.IsLocal` or `filepath.Localize` to validate or sanitize those paths.

## Beyond sanitization

Path sanitization is not sufficient when attackers may have access to part of
the local filesystem.

Multi-user systems are uncommon these days, but attacker access to the filesystem
can still occur in a variety of ways.
An unarchiving utility that extracts a tar or zip file may be induced
to extract a symbolic link and then extract a file name that traverses that link.
A container runtime may give untrusted code access to a portion of the local filesystem.

Programs may defend against unintended symlink traversal by using the
[`path/filepath.EvalSymlinks`](/pkg/path/filepath#EvalSymlinks)
function to resolve links in untrusted names before validation, but as described
above this two-step process is vulnerable to TOCTOU races.

Before Go 1.24, the safer option was to use a package such as
[github.com/google/safeopen](/pkg/github.com/google/safeopen),
that provides path traversal-resistant functions for opening a potentially-untrusted
filename within a specific directory.

## Introducing `os.Root`

In Go 1.24, we are introducing new APIs in the `os` package to safely open
a file in a location in a traversal-resistant fashion.

The new [`os.Root`](/pkg/os#Root) type represents a directory somewhere
in the local filesystem. Open a root with the [`os.OpenRoot`](/pkg/os#OpenRoot)
function:

    root, err := os.OpenRoot("/some/root/directory")
    if err != nil {
      return err
    }
    defer root.Close()

`Root` provides methods to operate on files within the root.
These methods all accept filenames relative to the root,
and disallow any operations that would escape from the root either
using relative path components ("..") or symlinks.

    f, err := root.Open("path/to/file")

`Root` permits relative path components and symlinks that do not escape the root.
For example, `root.Open("a/../b")` is permitted. Filenames are resolved using the
semantics of the local platform: On Unix systems, this will follow
any symlink in "a" (so long as that link does not escape the root);
while on Windows systems this will open "b" (even if "a" does not exist).

`Root` currently provides the following set of operations:

    func (*Root) Create(string) (*File, error)
    func (*Root) Lstat(string) (fs.FileInfo, error)
    func (*Root) Mkdir(string, fs.FileMode) error
    func (*Root) Open(string) (*File, error)
    func (*Root) OpenFile(string, int, fs.FileMode) (*File, error)
    func (*Root) OpenRoot(string) (*Root, error)
    func (*Root) Remove(string) error
    func (*Root) Stat(string) (fs.FileInfo, error)

In addition to the `Root` type, the new
[`os.OpenInRoot`](/pkg/os#OpenInRoot) function
provides a simple way to open a potentially-untrusted filename within a
specific directory:

    f, err := os.OpenInRoot("/some/root/directory", untrustedFilename)

The `Root` type provides a simple, safe, portable API for operating with untrusted filenames.

## Caveats and considerations

### Unix

On Unix systems, `Root` is implemented using the `openat` family of system calls.
A `Root` contains a file descriptor referencing its root directory and will track that
directory across renames or deletion.

`Root` defends against symlink traversal but does not limit traversal
of mount points. For example, `Root` does not prevent traversal of
Linux bind mounts. Our threat model is that `Root` defends against
filesystem constructs that may be created by ordinary users (such
as symlinks), but does not handle ones that require root privileges
to create (such as bind mounts).

### Windows

On Windows, `Root` opens a handle referencing its root directory.
The open handle prevents that directory from being renamed or deleted until the `Root` is closed.

`Root` prevents access to reserved Windows device names such as `NUL` and `COM1`.

### WASI

On WASI, the `os` package uses the WASI preview 1 filesystem API,
which are intended to provide traversal-resistant filesystem access.
Not all WASI implementations fully support filesystem sandboxing,
however, and `Root`'s defense against traversal is limited to that provided
by the WASI implementation.

### GOOS=js

When GOOS=js, the `os` package uses the Node.js file system API.
This API does not include the openat family of functions,
and so `os.Root` is vulnerable to TOCTOU (time-of-check-time-of-use) races in symlink
validation on this platform.

When GOOS=js, a `Root` references a directory name rather than a file descriptor,
and does not track directories across renames.

### Plan 9

Plan 9 does not have symlinks.
On Plan 9, a `Root` references a directory name and performs lexical sanitization of
filenames.

### Performance

`Root` operations on filenames containing many directory components can be much more expensive
than the equivalent non-`Root` operation. Resolving ".." components can also be expensive.
Programs that want to limit the cost of filesystem operations can use `filepath.Clean` to
remove ".." components from input filenames, and may want to limit the number of
directory components.

## Who should use os.Root?

You should use `os.Root` or `os.OpenInRoot` if:

  - you are opening a file in a directory; AND
  - the operation should not access a file outside that directory.

For example, an archive extractor writing files to an output directory should use
`os.Root`, because the filenames are potentially untrusted and it would be incorrect
to write a file outside the output directory.

However, a command-line program that writes output to a user-specified location
should not use `os.Root`, because the filename is not untrusted and may
refer to anywhere on the filesystem.

As a good rule of thumb, code which calls `filepath.Join` to combine a fixed directory
and an externally-provided filename should probably use `os.Root` instead.

    // This might open a file not located in baseDirectory.
    f, err := os.Open(filepath.Join(baseDirectory, filename))

    // This will only open files under baseDirectory.
    f, err := os.OpenInRoot(baseDirectory, filename)

## Future work

The `os.Root` API is new in Go 1.24.
We expect to make additions and refinements to it in future releases.

The current implementation prioritizes correctness and safety over performance.
Future versions will take advantage of platform-specific APIs, such as
Linux's `openat2`, to improve performance where possible.

There are a number of filesystem operations which `Root` does not support yet, such as
creating symbolic links and renaming files. Where possible, we will add support for these
operations. A list of additional functions in progress is in
[go.dev/issue/67002](/issue/67002).
