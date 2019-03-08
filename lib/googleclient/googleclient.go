package googleclient

import (
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"fmt"
	"google.golang.org/api/admin/directory/v1"
	"io/ioutil"
	"github.com/silinternational/serverless-google-groups-sync"
)

func getServiceForScopes(googleAuthUserEmail string, credBytes []byte, scope string) (*admin.Service, error) {
	service := &admin.Service{}

	config, err := google.JWTConfigFromJSON(credBytes, scope) // e.g. admin.AdminDirectoryGroupMemberScope
	if err != nil {
		return service, fmt.Errorf("Unable to parse client secret file to config: %s", err)
	}

	config.Subject = googleAuthUserEmail
	client := config.Client(context.Background())


	service, err = admin.New(client)
	if err != nil {
		return service, fmt.Errorf("Unable to retrieve directory Service: %s", err)
	}

	return service, nil
}

// GetGoogleAdminService authenticates with the Google API and returns an admin.Service
//  that has the scopes for Group and GroupMember
//  Authentication requires an email address that matches an actual GMail user (e.g. a machine account)
func GetGoogleAdminService(
	googleAuthUserEmail string,
	credentialsFilePath string,
	)(*admin.Service, error) {

	adminService := &admin.Service{}

	b, err := ioutil.ReadFile(credentialsFilePath) // e.g. "../../credentials.json"
	if err != nil {
		return adminService, fmt.Errorf("Unable to read client secret file: %s", err)
	}

	config, err := google.JWTConfigFromJSON(b, admin.AdminDirectoryGroupScope, admin.AdminDirectoryGroupMemberScope)
	if err != nil {
		return adminService, fmt.Errorf("Unable to parse client secret file to config: %s", err)
	}

	config.Subject = googleAuthUserEmail
	client := config.Client(context.Background())


	adminService, err = admin.New(client)
	if err != nil {
		return adminService, fmt.Errorf("Unable to retrieve directory Service: %s", err)
	}

	return adminService, nil
}


// GetMembersForGroup populates the TargetMembers attribute of a GroupDiff with the email addresses of the
// members of the corresponding Google Group.
func GetMembersForGroup(groupDiff *domain.GroupDiff, adminService *admin.Service) error {
	group := groupDiff.TargetGroup
	membersHolder, err := adminService.Members.List(group).Do()
	if err != nil {
		return fmt.Errorf("Unable to get members of group %s: %s", group, err)
	}

	membersList := membersHolder.Members
	members := []string{}

	for _, nextMember := range membersList {
		members = append(members, nextMember.Email)
	}

	groupDiff.TargetMembers = members

	return nil
}

func GetMembersForAllGroups(
	groupDiffs []*domain.GroupDiff,
	adminService *admin.Service,
) ([]*domain.GroupDiff, error) {

	for _, nextDiff := range groupDiffs {
		err := GetMembersForGroup(nextDiff, adminService)

		if err != nil {
			return groupDiffs, err
		}
	}

	return groupDiffs, nil
}

// AddMembersToGroup inserts new Gmail users into a Google Group
func AddMembersToGroup(groupName string, members []string, adminService *admin.Service) error {
	if len(members) < 1 {
		return nil
	}

	memberAdmin := adminService.Members

	for _, memberEmail := range members {
		newMember := admin.Member{
			Role: "MEMBER",
			Email: memberEmail,
		}

		_, err := memberAdmin.Insert(groupName, &newMember).Do()
		if err != nil {
			return fmt.Errorf("Unable to insert %s in Google group %s: %s", memberEmail, groupName, err)
		}
	}

	return nil
}


// DeleteMembersFromGroup deletes matching Gmail users from a Google Group
func DeleteMembersFromGroup(group string, members []string, adminService *admin.Service) error {
	if len(members) < 1 {
		return nil
	}

	memberAdmin := adminService.Members

	for _, memberEmail := range members {
		err := memberAdmin.Delete(group, memberEmail).Do()
		if err != nil {
			return fmt.Errorf("Unable to delete %s from Google group %s: %s", memberEmail, group, err)
		}
	}

	return nil
}

