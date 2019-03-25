# serverless-google-groups-sync

### AWS Credentials
See https://stackoverflow.com/a/46133337/856070 and
save AWS credentials to ./aws.credentials

### Google Service Account Configuration

(see https://stackoverflow.com/questions/53808710/authenticate-to-google-admin-directory-api#answer-53808774 and
 https://developers.google.com/admin-sdk/reports/v1/guides/delegation)

In the google developer console ...
* Create a new Service Account and a corresponding JSON credential file.
* Delegate Domain-Wide Authority to the Service Account.

The JSON credential file should contain something like this ...

```
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

The sync job will need to use the Service Account credentials to impersonate another user that has
domain superadmin privilege and who has logged in at least once into G Suite and
accepted the terms and conditions.

### List of Corresponding Group Names

The serverless build process will create an S3 bucket if it doesn't already exist with a name
that is based on ????? and ?????.

This bucket will not contain the creds.json file created during the Google Service Account Configuration (above)
as well as a json file containing a list of the group name pairs. This file should be named groups-map.json
and should look something like this (note: do not include a comma after the last entry) ...

```
[
 {source_group: "abc", google_group: "ABC"},
 {source_group: "def", google_group: "DEF"}
]
```