# Cognito Exporter

A utility to export Amazon Cognito users to CSV file

## Motivation

Amazon Cognito does not natively have functionality to export users

## Prerequisites

You should obtain the following information from Cognito first in order to use this utility:
- *Region*: The region in which the user pool resides
- *User pool ID*: ID of the user pool used for storing users (format: `{region}_{unique_key}`)
- *App ID*: ID of the application (e.g. Client ID in Cognito app client details)
- *Access Key & Secret Key*: This pair of information can be obtained from IAM Console for a user with `Power User` role

## Usage

1. Set environment variables in a secured place, using the `config.env.sample` file as your guidance
2. Refer to `example/main.go` file in this repository to know how to run the export functionality

## Things to do
This utility is working just fine for my personal purposes, but not without limitations. Refer to [this link](https://aws.github.io/aws-sdk-go-v2/docs/getting-started/) to find a way and enhance it to your heart content.
