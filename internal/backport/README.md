This directory contains backports of upcoming packages from the
standard library, beyond the Go version supplied by the latest Google
App Engine.

The procedure for adding a backport is fairly manual: copy the files
to a new directory and then global search and replace to modify import
paths to point to the backport.

As new versions of Go land on Google App Engine, this directory should
be pruned as much as possible.
