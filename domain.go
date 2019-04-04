package domain

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os"
)

// GoogleAuth represents the JSON credentials file provided by Google for Service Accounts
type GoogleAuth struct {
	Type                    string `json:"type"`
	ProjectID               string `json:"project_id"`
	PrivateKeyID            string `json:"private_key_id"`
	PrivateKey              string `json:"private_key"`
	ClientEmail             string `json:"client_email"`
	ClientID                string `json:"client_id"`
	AuthURI                 string `json:"auth_uri"`
	TokenURI                string `json:"token_uri"`
	AuthProviderX509CertURL string `json:"auth_provider_x509_cert_url"`
	ClientX509CertURL       string `json:"client_x509_cert_url"`
}

// MemberSourceApiConfig holds API information for fetching members to be sycned to a Google Group
//    BaseURL is the API base without trailing slash. It is assumed all groups to be sycned will be sub-paths from here
//    User is username for basic https auth
//    Pass is password for basic https auth
type MemberSourceApiConfig struct {
	BaseURL string
	User    string
	Pass    string
}

// GroupMap holds the source API path to the group list report along with the Google Group address it should be synced to
type GroupMap struct {
	SourcePath         string
	GoogleGroupAddress string
}

// AppConfig holds the full application configuration
type AppConfig struct {
	GoogleAuth            GoogleAuth
	GoogleDelegatedAdmin  string
	MemberSourceApiConfig MemberSourceApiConfig
	GroupMaps             []GroupMap
}

// GroupDiff holds the information needed by the sync process in relation to
//   the source group name,
//   the target group name (e.g. the google group),
//   the email addresses of the members of the source group
//   the email address of the current members of the target group
//   the email address of the source group that should be added to the target group
//   the email address that should be delete from the target group, since they
//     don't appear in the source group
type GroupDiff struct {
	SourceGroup     string
	TargetGroup     string
	SourceMembers   []string
	TargetMembers   []string
	MembersToAdd    []string
	MembersToDelete []string
}

// MemberSourceResponse represents the structure of the API response with list of members
type MemberSourceResponse struct {
	ReportEntry []struct {
		Email string `json:"Email"`
	} `json:"Report_Entry"`
}

// IsStringInStringSlice checks whether there is a match for a string
//  in a slice of strings
func IsStringInStringSlice(needle string, haystack []string) bool {
	for _, candidate := range haystack {
		if needle == candidate {
			return true
		}
	}
	return false
}

// GetEnv returns the value of the requested environment variable
//   or the given default value, if the environment variable's value is an
//   empty string.
func GetEnv(name, defaultValue string) string {
	value := os.Getenv(name)
	if value == "" {
		value = defaultValue
	}

	return value
}

// LoadAppConfig reads in the configuration.json file into an AppConfig struct
func LoadAppConfig(filename string) (AppConfig, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		log.Printf("unable to application config file %s, error: %s\n", filename, err.Error())
		return AppConfig{}, err
	}

	appConfig := AppConfig{}
	err = json.Unmarshal(data, &appConfig)
	if err != nil {
		log.Printf("unable to unmarshal application configuration file data, error: %s\n", err.Error())
		return AppConfig{}, err
	}

	return appConfig, nil
}
