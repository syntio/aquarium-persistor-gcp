# Aquarium GCP Persistor


Persistor is a Google Cloud Platform (GCP) component used for permanent data storage.

Its primary purpose is to serve as a back-up component that works independently of the rest of the system.

It is used to simplify the subsequent processing of data in the case of logical errors or data mapping. Also, it is used as a safety step.

Based on our previous experience, it is very useful to have data stored in cheap storage such as GCP Storage used in Persistor. Data stored that way can be easily made available to different processing needs (depending on the use case, e.g. BigQuery is not always the best option to access the data).


**Technology stack:**
 - Programing language: Go
 - Big data: GCP Storage
 - Compute: GCP Cloud Functions


**More details about Persistor can be found on**
- [Wiki page](../../wiki)

 
## Repository structure

 Persistor is divided into five modules. A module is a collection of related Go packages that are released together. A file named go.mod in each folder declares the module path: the import path prefix for all packages within the module.

 
The following structure allows us not to build/deploy the entire project, but only the modules we are planning to use.

 

```bash
+---invoker
|       go.mod
|       invoker.go
|
+---lib
|       getEnvVariable.go
|       go.mod
|       invokerInfo.go
|       puller.go
|       pullerInfo.go
|       storage.go
|       storageInfo.go
|
+---pull
|       go.mod
|       pull.go
|
+---push
|       go.mod
|       push.go
|
+---streamingPull
|       go.mod
|       streamingPull.go
|   .gitignore
|   CODE_OF_CONDUCT.md
|   CONTRIBUTING.md
|   LICENSE.md
|   README.md
\   SUPPORT_INFORMATION.md
```

## Dependencies

**GCP**

- Existing or new project 
- Enabled billing for your Cloud project
- Enabled the Cloud Functions, Cloud Pub/Sub and Cloud Build APIs.
- Pub/Sub Topic
- Storage Bucket

**Program language:** Go 1.13

 
## Features


Persistor supports three different subscriber mechanisms:

  - **Push** for minor activity of data publishers, topic and subscription need to be in the same project
  - **Synchronous Pull** for periodically active data publishers
  - **Streaming Pull** for continuous flow of messages/events


### Storing messages

Subscriber component reads messages from Pub/Sub subscription, and store them on GCP Storage (GCS) in **raw** format. 

Files on GCS are stored in folder structure based on ingestion date. Format of a folder: 

`[Bucket Name]/[YYYY]/[MM]/[DD]/[HH]`

Each message is stored as separate object on GCS, with content identical to payload of Pub/Sub message.

## Developing


Instructions on how to set up a local development environment and how to integrate it into an existing environment

To start developing the project further:

```shell
git clone https://github.com/syntio/aquarium-persistor-gcp.git
cd aquarium-persistor-gcp
```

In order to create a new folder/package in the existing solution and to be able to use the library, it is necessary to create and configure the *go.mod* by adding the following line:


```go
replace github.com/syntio/aquarium-persistor-gcp/lib => ../lib
```

The library import path requires the following form:

```go
github.com/syntio/aquarium-persistor-gcp/lib
```

## Deployment and Configuration


Instructions on how to establish GCP persistor connection between Pub/Sub and GCS storage using Cloud shell can be found [here](../../../wiki/Deployment-via-gcloud-shell).

## Limitations

Pull does not always store the exact number of messages that was determined by the Number of messages parameter. There could be 1 - 3 more messages pulled and stored.
Timeout parameter of a Cloud Function that pulls and stores messages synchronously needs to be adjusted so that the function terminates shortly after the last message is stored. Otherwise the function will stay active until timeout.


## Links
 
Issue tracker: https://github.com/syntio/aquarium-persistor-gcp/issues
- *In case of sensitive bugs like security vulnerabilities, please contact  support@syntio.net directly instead of using issue tracker. We value your effort to improve the security and privacy of this project!*


## Contributing
Please refer to [CONTRIBUTING.md](./CONTRIBUTING.md)

## Developed by

The repository is developed and maintained by Syntio Labs.

## Licence 
Licensed under the Apache License, Version 2.0
