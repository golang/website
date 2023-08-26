FROM golang:1.20-alpine AS base
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download 
COPY . .

FROM base as build-server
RUN go build -o /bin/server ./cmd/golangorg

FROM scratch as prod 
COPY --from=build-server /bin/server /bin/
EXPOSE 8080
CMD [ "/bin/server", "-http=0.0.0.0:8080" ]
