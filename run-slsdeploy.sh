#!/usr/bin/env bash

mkdir -p ~/.aws

cp /go/src/github.com/silinternational/serverless-google-groups-sync/aws.credentials ~/.aws/credentials

sls deploy -v --force