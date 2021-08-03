---
title: What's new in the Go Cloud Development Kit
date: 2019-03-04
by:
- The Go Cloud Development Kit team at Google
summary: Recent changes to the Go Cloud Development Kit (Go CDK).
---

## Introduction

Last July, we [introduced](https://blog.golang.org/go-cloud) the [Go Cloud Development Kit](https://gocloud.dev)
(previously referred to as simply "Go Cloud"),
an open source project building libraries and tools to improve the experience
of developing for the cloud with Go.
We've made a lot of progress since then -- thank you to early contributors!
We look forward to growing the Go CDK community of users and contributors,
and are excited to work closely with early adopters.

## Portable APIs

Our first initiative is a set of portable APIs for common cloud services.
You write your application using these APIs,
and then deploy it on any combination of providers,
including AWS, GCP, Azure, on-premise, or on a single developer machine for testing.
Additional providers can be added by implementing an interface.

These portable APIs are a great fit if any of the following are true:

  - You develop cloud applications locally.
  - You have on-premise applications that you want to run in the cloud (permanently, or as part of a migration).
  - You want portability across multiple clouds.
  - You are creating a new Go application that will use cloud services.

Unlike traditional approaches where you would need to write new application
code for each cloud provider,
with the Go CDK you write your application code once using our portable
APIs to access the set of services listed below.
Then, you can run your application on any supported cloud with minimal config changes.

Our current set of APIs includes:

  - [blob](https://godoc.org/gocloud.dev/blob),
    for persistence of blob data.
    Supported providers include: AWS S3, Google Cloud Storage (GCS),
    Azure Storage, the filesystem, and in-memory.
  - [pubsub](https://godoc.org/gocloud.dev/pubsub) for publishing/subscribing
    of messages to a topic.
    Supported providers include: Amazon SNS/SQS,
    Google Pub/Sub, Azure Service Bus, RabbitMQ, and in-memory.
  - [runtimevar](https://godoc.org/gocloud.dev/runtimevar),
    for watching external configuration variables.
    Supported providers include AWS Parameter Store,
    Google Runtime Configurator, etcd, and the filesystem.
  - [secrets](https://godoc.org/gocloud.dev/secrets),
    for encryption/decryption.
    Supported providers include AWS KMS, GCP KMS,
    Hashicorp Vault, and local symmetric keys.
  - Helpers for connecting to cloud SQL providers. Supported providers include AWS RDS and Google Cloud SQL.
  - We are also working on a document storage API (e.g. MongoDB, DynamoDB, Firestore).

## Feedback

We hope you're as excited about the Go CDK as we are -- check out our [godoc](https://godoc.org/gocloud.dev),
walk through our [tutorial](https://github.com/google/go-cloud/tree/master/samples/tutorial),
and use the Go CDK in your application(s).
We'd love to hear your ideas for other APIs and API providers you'd like to see.

If you're digging into Go CDK please share your experiences with us:

  - What went well?
  - Were there any pain points using the APIs?
  - Are there any features missing in the API you used?
  - Suggestions for documentation improvements.

To send feedback, you can:

  - Submit issues to our public [GitHub repository](https://github.com/google/go-cloud/issues/new/choose).
  - Email [go-cdk-feedback@google.com](mailto:go-cdk-feedback@google.com).
  - Post to our [public Google group](https://groups.google.com/forum/#!forum/go-cloud).

Thanks!
