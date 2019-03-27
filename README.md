# serverless-google-groups-sync

This application is intended to retrieve a list of email addresses that should be members of a Google Group 
and perform synchronization logic to add/remove members from the Google Group based on the list received from API call. 
It is written to take advantage of the Serverless framework and AWS Lambda so that it can be scheduled to run as often 
as you desire based on volatility of source list. 

## Prerequisites
Before you can use this script you need a few things including: 

1. AWS credentials for the Serverless framework to provision and publish the Lambda function as well as an S3 bucket for 
storing application configuration.
2. Google Service Account credentials for integrating with Google Groups
3. A GSuite admin account with permissions to the Google Groups APIs that this script can act on behalf of (delegated user)
4. A source API that returns an array of email addresses per Google Group you want synced

## Application Configuration
Since this application needs some environmental configuration that is too complex for simple environment variables it 
has been developed to retrieve a JSON file from a secure S3 bucket. The structure of the `config.json` file is:

```json
{
  "GoogleAuth": {
    "type": "service_account",
    "project_id": "abc-theme-123456",
    "private_key_id": "abc123",
    "private_key": "-----BEGIN PRIVATE KEY-----\nMIIabc...\nabc...\n...xyz\n-----END PRIVATE KEY-----\n",
    "client_email": "my-sync-bot@abc-theme-123456.iam.gserviceaccount.com",
    "client_id": "123456789012345678901",
    "auth_uri": "https://accounts.google.com/o/oauth2/auth",
    "token_uri": "https://oauth2.googleapis.com/token",
    "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
    "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/my-sync-bot%40abc-theme-123456.iam.gserviceaccount.com"
  },
  "GoogleDelegatedAdmin": "gsuite-admin-with-groups-api-access@domain.com",
  "MemberSourceApiConfig": {
    "BaseURL": "https://api.domain.com",
    "User": "user",
    "Pass": "pass"
  },
  "GroupMaps": [
    {
      "SourcePath": "/group1",
      "GoogleGroupAddress": "group1@groups.domain.com"
    },
    {
      "SourcePath": "/group2",
      "GoogleGroupAddress": "group2@groups.domain.com"
    }
  ]
}
```

Below are further instructions on how to get the Google credentials shown above.

### AWS Credentials
See https://stackoverflow.com/a/46133337/856070 and
save AWS credentials to ./aws.credentials

### Google Service Account Configuration

(see https://stackoverflow.com/questions/53808710/authenticate-to-google-admin-directory-api#answer-53808774 and
 https://developers.google.com/admin-sdk/reports/v1/guides/delegation)

In the google developer console ...
* Create a new Service Account and a corresponding JSON credential file.
* Delegate Domain-Wide Authority to the Service Account.
* The email address for this user should be stored in the `config.json` as the `GoogleDelegatedAdmin` value

The JSON credential file should contain something like this ...

```json
{
  "type": "service_account",
  "project_id": "abc-theme-123456",
  "private_key_id": "abc123",
  "private_key": "-----BEGIN PRIVATE KEY-----\nMIIabc...\nabc...\n...xyz\n-----END PRIVATE KEY-----\n",
  "client_email": "my-sync-bot@abc-theme-123456.iam.gserviceaccount.com",
  "client_id": "123456789012345678901",
  "auth_uri": "https://accounts.google.com/o/oauth2/auth",
  "token_uri": "https://oauth2.googleapis.com/token",
  "auth_provider_x509_cert_url": "https://www.googleapis.com/oauth2/v1/certs",
  "client_x509_cert_url": "https://www.googleapis.com/robot/v1/metadata/x509/my-sync-bot%40abc-theme-123456.iam.gserviceaccount.com"
}
```

These contents will need to be copied into the `config.json` file as the value of the `GoogleAuth` key.

The sync job will need to use the Service Account credentials to impersonate another user that has
domain superadmin privilege and who has logged in at least once into G Suite and
accepted the terms and conditions.

### Source API for users to be synced into groups
To keep things simple this application expects the response from the API to be a JSON array of email addresses, ex:

```json
["user1@domain.com","user2@domain.com","user3@domain.com"]
```

Configuration for the source API is stored in `config.json` as `MemberSourceApiConfig` with the structure:

```json
{
    "BaseURL": "https://api.domain.com",
    "User": "user",
    "Pass": "pass"
}
```

### Map of API source paths to Google Group addresses
In the main application `config.json` file you will create a map of source paths to Google Group addresses as the value
for `GroupMaps`

Format should be:

```json
[
    {
      "SourcePath": "/group1",
      "GoogleGroupAddress": "group1@groups.domain.com"
    },
    {
      "SourcePath": "/group2",
      "GoogleGroupAddress": "group2@groups.domain.com"
    }
]
```

## Deployment
After gathering all the prerequisite data and creating your `config.json` file, you should be ready for deployment.

Rather than requiring you to install all the build and deployment depencies locally you can use the included Dockerfile to 
perform the work. 

1. Copy `.env.example` to `.env` and update values in it
2. Run `make build` - this will use Docker Compose to build the included Docker image and build the Go binary 
and place it in the `./bin/` directory
3. Run `make deploy` - this will run the `serveless deploy` command inside the Docker image to deploy to AWS

If you want to remove the deployed resources, run `make remove`

