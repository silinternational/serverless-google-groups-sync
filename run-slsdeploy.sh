#!/usr/bin/env bash

export S3_BUCKET_FOR_GROUPS_MAP="${DEV_S3_BUCKET_FOR_GROUPS_MAP}"

export VPC_SG_ID="${DEV_VPC_SG_ID}"
export VPC_SUBNET1="${DEV_VPC_SUBNET1}"
export VPC_SUBNET2="${DEV_VPC_SUBNET2}"
export VPC_SUBNET3="${DEV_VPC_SUBNET3}"

mkdir -p ~/.aws

cp /go/src/github.com/silinternational/serverless-google-groups-sync/aws.credentials ~/.aws/credentials

sls deploy -v --force