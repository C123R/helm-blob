# helm-blob [![Build Status](https://travis-ci.com/C123R/helm-blob.svg?token=9FxxpxQR16mxxejVyCbA&branch=master)](https://travis-ci.com/C123R/helm-blob)

`helm-blob` plugin allows you to manage helm repositories on the blob storage like Azure Blob, GCS, S3, etc.

This plugin supports operations like uploading or deletion of charts from remote Helm Repository hosted on Blob Storage. It could be used to initialize the new Helm Repository.

`helm-blob` was inspired by [Alex Khaerov's](https://github.com/hayorov) [helm-gcs](https://github.com/hayorov/helm-gcs) plugin with extending support for Azure Blob storage and S3, which makes helm-blob to support Azure Blob, GCS, S3 storage.

This plugin uses Go Cloud's [Blob](https://gocloud.dev/howto/blob/) package.

## Installation

```sh
helm plugin install https://github.com/C123R/helm-blob.git
```

To install specific version of:

```sh
helm plugin install https://github.com/C123R/helm-blob.git --version 0.2.0
```

### If you are still using Helm Below Version 3:

```sh
helm plugin install https://github.com/C123R/helm-blob.git --version 0.1.0
```

## Usage

**Note:** This plugin will not provide new blob storage, You must first create blob storage container/bucket that will be used as a remote chart repository.

- ### Initialize a new chart repository

  ```sh
  helm blob init azblob://helmrepo

  OR

  helm blob init gs://helmrepo/charts
  ```

- ### Add your repository to Helm

  ```sh
  helm repo add azurehelm azblob://helmrepo
  ```

- ### Push a new chart to your repository

  ```sh
  helm blob push mychart.tar.gz azurehelm
  ```

- ### Updating Helm cache (Required after pushing new chart)

  ```sh
  helm repo update
  ```

- ### Fetch the chart

  ```sh
  helm fetch azurehelm/mychart
  ```

- ### Delete a chart

  ```sh
  helm blob delete mychart azurehelm
  ```

  Note: This will delete all chart versions from remote repository. To delete a specific chart:

  ```sh
  helm blob delete mychart -v 0.3.0 azurehelm
  ```

## Authentication

Helm blob's plugin authentication varies depending upon the blob provider as mentioned below:

- ### S3

  S3 provider support AWS [default credential provider](https://docs.aws.amazon.com/sdk-for-go/v1/developer-guide/configuring-sdk.html#specifying-credentials) chain in the following order:

  - Environment variables.

  - Shared credentials file.

  - If your application is running on an Amazon EC2 instance, IAM role for Amazon EC2.

- ### Azure Blob

  Currently it supports authentication only with [environment variables](https://docs.microsoft.com/en-us/azure/storage/common/storage-azure-cli#set-default-azure-storage-account-environment-variables):

  - AZURE_STORAGE_ACCOUNT
  - AZURE_STORAGE_KEY or AZURE_STORAGE_SAS_TOKEN

- ### [GCS](https://cloud.google.com/docs/authentication/production)

  GCS provider uses [Application Default Credentials](https://cloud.google.com/docs/authentication/production) in the following order:

  - Environment Variable (GOOGLE_APPLICATION_CREDENTIALS)
  - Default Service Account from the compute instance(Compute Engine, Kubernetes Engine, Cloud function etc).

  To authenticate against GCS you can:

  - Use the [application default credentials](https://cloud.google.com/sdk/gcloud/reference/auth/application-default/)

  - Use a service account via the global flag `--service-account`

  See [the GCP documentation](https://cloud.google.com/docs/authentication/production#providing_credentials_to_your_application) for more information.
