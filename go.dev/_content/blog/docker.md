---
title: Deploying Go servers with Docker
date: 2014-09-26
by:
- Andrew Gerrand
summary: How to use Docker's new official base images for Go.
---

## Introduction

This week Docker [announced](https://blog.docker.com/2014/09/docker-hub-official-repos-announcing-language-stacks/)
official base images for Go and other major languages,
giving programmers a trusted and easy way to build containers for their Go programs.

In this article we'll walk through a recipe for creating a Docker container for
a simple Go web application and deploying that container to Google Compute Engine.
If you're not familiar with Docker, you should read
[Understanding Docker](https://docs.docker.com/engine/understanding-docker/)
before reading on.

## The demo app

For our demonstration we will use the
[outyet](https://pkg.go.dev/golang.org/x/example/outyet) program from the
[Go examples repository](https://cs.opensource.google/go/x/example),
a simple web server that reports whether the next version of Go has been released
(designed to power sites like [isgo1point4.outyet.org](http://isgo1point4.outyet.org/)).
It has no dependencies outside the standard library and requires no additional
data files at run time; for a web server, it's about as simple as it gets.

Use "go get" to fetch and install outyet in your
[workspace](/doc/code.html#Workspaces):

	$ go get golang.org/x/example/outyet

## Write a Dockerfile

Replace a file named `Dockerfile` in the `outyet` directory with the following contents:

	# Start from a Debian image with the latest version of Go installed
	# and a workspace (GOPATH) configured at /go.
	FROM golang

	# Copy the local package files to the container's workspace.
	ADD . /go/src/golang.org/x/example/outyet

	# Build the outyet command inside the container.
	# (You may fetch or manage dependencies here,
	# either manually or with a tool like "godep".)
	RUN go install golang.org/x/example/outyet

	# Run the outyet command by default when the container starts.
	ENTRYPOINT /go/bin/outyet

	# Document that the service listens on port 8080.
	EXPOSE 8080

This `Dockerfile` specifies how to construct a container that runs `outyet`,
starting with the basic dependencies (a Debian system with Go installed;
the [official `golang` docker image](https://registry.hub.docker.com/_/golang/)),
adding the `outyet` package source, building it, and then finally running it.

The `ADD`, `RUN`, and `ENTRYPOINT` steps are common tasks for any Go project.
To simplify this, there is an
[`onbuild` variant](https://github.com/docker-library/golang/blob/9ff2ccca569f9525b023080540f1bb55f6b59d7f/1.3.1/onbuild/Dockerfile)
of the `golang` image that automatically copies the package source, fetches the
application dependencies, builds the program, and configures it to run on
startup.

With the `onbuild` variant, the `Dockerfile` is much simpler:

	FROM golang:onbuild
	EXPOSE 8080

## Build and run the image

Invoke Docker from the `outyet` package directory to build an image using the `Dockerfile`:

	$ docker build -t outyet .

This will fetch the `golang` base image from Docker Hub, copy the package source
to it, build the package inside it, and tag the resulting image as `outyet`.

To run a container from the resulting image:

	$ docker run --publish 6060:8080 --name test --rm outyet

The `--publish` flag tells docker to publish the container's port `8080` on the
external port `6060`.

The `--name` flag gives our container a predictable name to make it easier to work with.

The `--rm` flag tells docker to remove the container image when the outyet server exits.

With the container running, open `http://localhost:6060/` in a web browser and
you should see something like this:

{{image "docker/outyet.png"}}

(If your docker daemon is running on another machine (or in a virtual machine),
you should replace `localhost` with the address of that machine. If you're
using [boot2docker](http://boot2docker.io/) on OS X or Windows you can find
that address with `boot2docker ip`.)

Now that we've verified that the image works, shut down the running container
from another terminal window:

	$ docker stop test

## Create a repository on Docker Hub

[Docker Hub](https://hub.docker.com/), the container registry from which we
pulled the `golang` image earlier, offers a feature called
[Automated Builds](http://docs.docker.com/docker-hub/builds/) that builds
images from a GitHub or BitBucket repository.

By committing [the Dockerfile](https://go.googlesource.com/example/+/refs/heads/master/outyet/)
to the repository and creating an
[automated build](https://registry.hub.docker.com/u/adg1/outyet/)
for it, anyone with Docker installed can download and run our image with a
single command. (We will see the utility of this in the next section.)

To set up an Automated Build, commit the Dockerfile to your repo on
[GitHub](https://github.com/) or [BitBucket](https://bitbucket.org/),
create an account on Docker Hub, and follow the instructions for
[creating an Automated Build](http://docs.docker.com/docker-hub/builds/).

When you're done, you can run your container using the name of the automated build:

	$ docker run goexample/outyet

(Replace `goexample/outyet` with the name of the automated build you created.)

## Deploy the container to Google Compute Engine

Google provides
[container-optimized Google Compute Engine images](https://developers.google.com/compute/docs/containers/container_vms)
that make it easy to spin up a virtual machine running an arbitrary Docker container.
On startup, a program running on the instance reads a configuration file that
specifies which container to run, fetches the container image, and runs it.

Create a [containers.yaml](https://cloud.google.com/compute/docs/containers/container_vms#container_manifest)
file that specifies the docker image to run and the ports to expose:

	version: v1beta2
	containers:
	- name: outyet
	  image: goexample/outyet
	  ports:
	  - name: http
	    hostPort: 80
	    containerPort: 8080

(Note that we're publishing the container's port `8080` as external port `80`,
the default port for serving HTTP traffic. And, again, you should replace
`goexample/outyet` with the name of your Automated Build.)

Use the [gcloud tool](https://cloud.google.com/sdk/#Quick_Start)
to create a VM instance running the container:

	$ gcloud compute instances create outyet \
		--image container-vm-v20140925 \
		--image-project google-containers \
		--metadata-from-file google-container-manifest=containers.yaml \
		--tags http-server \
		--zone us-central1-a \
		--machine-type f1-micro

The first argument (`outyet`) specifies the instance name, a convenient label
for administrative purposes.

The `--image` and `--image-project` flags specify the special
container-optimized system image to use (copy these flags verbatim).

The `--metadata-from-file` flag supplies your `containers.yaml` file to the VM.

The `--tags` flag tags your VM instance as an HTTP server, adjusting the
firewall to expose port 80 on the public network interface.

The `--zone` and `--machine-type` flags specify the zone in which to run the VM
and the type of machine to run. (To see a list of machine types and the zones,
run `gcloud compute machine-types list`.)

Once this has completed, the gcloud command should print some information about
the instance. In the output, locate the `networkInterfaces` section to find the
instance's external IP address. Within a couple of minutes you should be able
to access that IP with your web browser and see the "Has Go 1.4 been released
yet?" page.

(To see what's happening on the new VM instance you can ssh into it with
`gcloud compute ssh outyet`. From there, try `sudo docker ps` to see which
Docker containers are running.)

## Learn more

This is just the tip of the iceberg—there's a lot more you can do with Go, Docker, and Google Compute Engine.

To learn more about Docker, see their [extensive documentation](https://docs.docker.com/).

To learn more about Docker and Go, see the [official `golang` Docker Hub repository](https://registry.hub.docker.com/_/golang/) and Kelsey Hightower's [Optimizing Docker Images for Static Go Binaries](https://medium.com/@kelseyhightower/optimizing-docker-images-for-static-binaries-b5696e26eb07).

To learn more about Docker and [Google Compute Engine](http://cloud.google.com/compute),
see the [Container-optimized VMs page](https://cloud.google.com/compute/docs/containers/container_vms)
and the [google/docker-registry Docker Hub repository](https://registry.hub.docker.com/u/google/docker-registry/).
