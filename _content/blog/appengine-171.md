---
title: Go updates in App Engine 1.7.1
date: 2012-08-22
by:
- Andrew Gerrand
tags:
- appengine
- release
summary: App Engine SDK 1.7.1 adds memcache and other functionality for Go.
---


This week we released version 1.7.1 of the App Engine SDK.
It includes some significant updates specific to the App Engine runtime for Go.

The [memcache package](https://developers.google.com/appengine/docs/go/memcache/reference) has
had some additions to its [Codec](https://developers.google.com/appengine/docs/go/memcache/reference#Codec) convenience type.
The SetMulti, AddMulti, CompareAndSwap, and CompareAndSwapMulti methods
make it easier to store and update encoded data in the [Memcache Service](https://developers.google.com/appengine/docs/go/memcache/overview).

The [bulkloader tool](https://developers.google.com/appengine/docs/go/tools/uploadingdata) can
now be used with Go apps,
allowing users to upload and download datastore records in bulk.
This is useful for backups and offline processing,
and a great help when migrating Python or Java apps to the Go runtime.

The [Images Service](https://developers.google.com/appengine/docs/go/images/overview) is
now available to Go users.
The new [appengine/image package](https://developers.google.com/appengine/docs/go/images/reference) supports serving
images directly from Blobstore and resizing or cropping those images on the fly.
Note that this is not the full image service as provided by the Python and Java SDKs,
as much of the equivalent functionality is available in the [standard Go image package](/pkg/image/) and
external packages such as [graphics-go](http://code.google.com/p/graphics-go/).

The new [runtime.RunInBackground](https://developers.google.com/appengine/docs/go/backends/runtime#RunInBackground) function
allows backend requests to spawn a new request independent of the initial request.
These can run in the background as long as the backend stays alive.

Finally, we have filled in some missing functionality:
the [xmpp package](https://developers.google.com/appengine/docs/go/xmpp/reference) now
supports sending presence updates and chat invitations and retrieving the
presence state of another user,
and the [user package](https://developers.google.com/appengine/docs/go/users/reference) supports
authenticating clients with OAuth.

You can grab the new SDK from the [App Engine downloads page](https://developers.google.com/appengine/downloads#Google_App_Engine_SDK_for_Go) and
browse the [updated documentation](https://developers.google.com/appengine/docs/go).
