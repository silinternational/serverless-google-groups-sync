package googleclient

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/silinternational/serverless-google-groups-sync"
	"golang.org/x/net/context"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/admin/directory/v1"
)

func getServiceForScopes(delegatedAdmin string, googleAuth domain.GoogleAuth, oauthScope string) (*admin.Service, error) {
	service := &admin.Service{}

	googleAuthJson, err := json.Marshal(googleAuth)
	if err != nil {
		log.Printf("unable to marshal google auth struct into json, error: %s\n", err.Error())
		return service, err
	}

	config, err := google.JWTConfigFromJSON(googleAuthJson, oauthScope) // e.g. admin.AdminDirectoryGroupMemberScope
	if err != nil {
		return service, fmt.Errorf("unable to parse client secret file to config: %s", err)
	}

	config.Subject = delegatedAdmin
	client := config.Client(context.Background())

	service, err = admin.New(client)
	if err != nil {
		return service, fmt.Errorf("unable to retrieve directory Service: %s", err)
	}

	return service, nil
}

// GetGoogleAdminService authenticates with the Google API and returns an admin.Service
//  that has the scopes for Group and GroupMember
//  Authentication requires an email address that matches an actual GMail user (e.g. a machine account)
func GetGoogleAdminService(
	delegatedAdmin string,
	googleAuth domain.GoogleAuth,
) (*admin.Service, error) {

	adminService := &admin.Service{}

	googleAuthJson, err := json.Marshal(googleAuth)
	if err != nil {
		log.Printf("unable to marshal google auth struct into json, error: %s\n", err.Error())
		return adminService, err
	}

	config, err := google.JWTConfigFromJSON(googleAuthJson, admin.AdminDirectoryGroupScope, admin.AdminDirectoryGroupMemberScope)
	if err != nil {
		return adminService, fmt.Errorf("unable to parse client secret file to config: %s", err)
	}

	config.Subject = delegatedAdmin
	client := config.Client(context.Background())

	adminService, err = admin.New(client)
	if err != nil {
		return adminService, fmt.Errorf("unable to retrieve directory Service: %s", err)
	}

	return adminService, nil
}

// GetMembersForGroup populates the TargetMembers attribute of a GroupDiff with the email addresses of the
// members of the corresponding Google Group.
func GetMembersForGroup(groupDiff *domain.GroupDiff, adminService *admin.Service) error {
	group := groupDiff.TargetGroup
	membersHolder, err := adminService.Members.List(group).Do()
	if err != nil {
		return fmt.Errorf("unable to get members of group %s: %s", group, err)
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
			Role:  "MEMBER",
			Email: memberEmail,
		}

		_, err := memberAdmin.Insert(groupName, &newMember).Do()
		if err != nil {
			return fmt.Errorf("unable to insert %s in Google group %s: %s", memberEmail, groupName, err)
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
			return fmt.Errorf("unable to delete %s from Google group %s: %s", memberEmail, group, err)
		}
	}

	return nil
}
