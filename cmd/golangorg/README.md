# golangorg

## Local Development

For local development, simply build and run. It serves on localhost:6060.

	go run .

## Local Production Mode

To run in production mode locally, you need:

  * the Google Cloud SDK; see https://cloud.google.com/sdk/
  * Redis
  * Go sources under $GOROOT
  * Godoc sources inside $GOPATH
    (`go get -d golang.org/x/website/cmd/golangorg`)

Build with the `golangorg` tag and run:

	go build -tags golangorg
	./golangorg

In production mode it serves on localhost:8080 (not 6060).
The port is controlled by $PORT, as in:

	PORT=8081 ./golangorg

## Local Production Mode using Docker

To run in production mode locally using Docker, build the app's Docker container:

	make docker-build

Make sure redis is running on port 6379:

	$ echo PING | nc localhost 6379
	+PONG
	^C

Run the datastore emulator:

	gcloud beta emulators datastore start --project golang-org

In another terminal window, run the container:

	$(gcloud beta emulators datastore env-init)

	docker run --rm \
		--net host \
		--env GODOC_REDIS_ADDR=localhost:6379 \
		--env DATASTORE_EMULATOR_HOST=$DATASTORE_EMULATOR_HOST \
		--env DATASTORE_PROJECT_ID=$DATASTORE_PROJECT_ID \
		gcr.io/golang-org/godoc

It serves on localhost:8080.

## Deploying to golang.org

1.	Run `make cloud-build deploy` to build the image, push it to gcr.io,
	and deploy to Flex (but not yet update golang.org to point to it).

2.	Check that the new version, mentioned on "target url:" line, looks OK.

3.	If all is well, run `make publish` to publish the new version to golang.org.
	It will run regression tests and then point the load balancer to the newly
	deployed version.

4.	Stop and/or delete any very old versions. (Stopped versions can be re-started.)
	Keep at least one older verson to roll back to, just in case.

	You can view, stop/delete, or migrate traffic between versions via the
	[GCP Console UI](https://console.cloud.google.com/appengine/versions?project=golang-org&serviceId=default&pageState=(%22versionsTable%22:(%22f%22:%22%255B%257B_22k_22_3A_22Environment_22_2C_22t_22_3A10_2C_22v_22_3A_22_5C_22Flexible_5C_22_22_2C_22s_22_3Atrue_2C_22i_22_3A_22env_22%257D%255D%22))).

5.	You're done.

## Troubleshooting

Ensure the Cloud SDK is on your PATH and you have the app-engine-go component
installed (`gcloud components install app-engine-go`) and your components are
up-to-date (`gcloud components update`).

For deployment, make sure you're signed in to gcloud:

	gcloud auth login
