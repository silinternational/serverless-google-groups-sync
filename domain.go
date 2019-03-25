package domain

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

// Functions that populates the members of a certain group
type GroupMembersGetter func(*GroupDiff) error

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

// MemberSourceApiConfig holds API information for fetching members to be sycned to a Google Group
//    BaseURL is the API base without trailing slash. It is assumed all groups to be sycned will be sub-paths from here
//    User is username for basic https auth
//    Pass is password for basic https auth
type MemberSourceApiConfig struct {
	BaseURL string
	User    string
	Pass    string
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

// GetGroupMembersFromSource calls HTTPS API to get members for given group
func GetGroupMembersFromSource(apiConfig MemberSourceApiConfig, groupPath string) ([]string, error) {

	client := &http.Client{}
	req, err := http.NewRequest("GET", apiConfig.BaseURL+groupPath, nil)
	req.SetBasicAuth(apiConfig.User, apiConfig.Pass)
	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return []string{}, err
	}

	bodyText, err := ioutil.ReadAll(resp.Body)

	members := []string{}
	err = json.Unmarshal(bodyText, &members)
	if err != nil {
		log.Println(err)
		return []string{}, err
	}

	return members, nil
}
