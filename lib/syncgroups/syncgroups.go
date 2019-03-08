package syncgroups

import (
	"google.golang.org/api/admin/directory/v1"
	"github.com/silinternational/serverless-google-groups-sync"
	"github.com/silinternational/serverless-google-groups-sync/lib/googleclient"
	"fmt"
	"log"
)


// DiffGroup populates the slice of members to add to the target group,
//    because they are in the source group but not the target group
//        -- and --
//  populates the slice of members to delete from the target group,
//    because they are in the target group but not in the source group
func DiffGroup(groupDiff *domain.GroupDiff) {
	groupDiff.MembersToAdd = []string{}
	groupDiff.MembersToDelete = []string{}

	for _, mSource := range groupDiff.SourceMembers {
		if !domain.IsStringInStringSlice(mSource, groupDiff.TargetMembers) {
			groupDiff.MembersToAdd = append(groupDiff.MembersToAdd, mSource)
		}
	}

	for _, mTarget := range groupDiff.TargetMembers {
		if !domain.IsStringInStringSlice(mTarget, groupDiff.SourceMembers) {
			groupDiff.MembersToDelete = append(groupDiff.MembersToDelete, mTarget)
		}
	}
}


// DiffAllGroups populates each group's slices of members to add and members to delete
func DiffAllGroups(allGroups []*domain.GroupDiff) []*domain.GroupDiff {
	modifiedGroups := []*domain.GroupDiff{}

	for _, nextDiff := range allGroups {
		DiffGroup(nextDiff)
		modifiedGroups = append(modifiedGroups, nextDiff)
	}

	return modifiedGroups
}

func GetSourceMembersForAllGroups(groupDiffs []*domain.GroupDiff, getter domain.GroupMembersGetter) ([]*domain.GroupDiff, error) {
	for _, nextDiff := range groupDiffs {
		err := getter(nextDiff)

		if err != nil {
			return groupDiffs, err
		}
	}

	return groupDiffs, nil
}


// InitAllGroupDiffs gets all the groups from the source system and then creates
//  a slice of pointers to matching GroupDiffs that have all their attributes populated,
//  including the matching Target groups and their members
func InitAllGroupDiffs(
	correspondingGroups [][2]string,
	googleAdminService *admin.Service,
	sourceMemberGetter domain.GroupMembersGetter,
) ([]*domain.GroupDiff, error) {

	groupDiffs := []*domain.GroupDiff{}
	for _, groupPair := range correspondingGroups {
		newGroupDiff := domain.GroupDiff{
			SourceGroup: groupPair[0],
			TargetGroup: groupPair[1],
		}
		groupDiffs = append(groupDiffs, &newGroupDiff)
	}

	groupDiffs, err := GetSourceMembersForAllGroups(groupDiffs, sourceMemberGetter)
	if err != nil {
		return groupDiffs, err
	}

	groupDiffs, err = googleclient.GetMembersForAllGroups(groupDiffs, googleAdminService)
	if err != nil {
		return groupDiffs, err
	}

	groupDiffs = DiffAllGroups(groupDiffs)

	return groupDiffs, nil
}



func SyncGroups(
	correspondingGroups [][2]string,
	googleAuthUserEmail string,
	credentialsFilePath string,
	sourceMemberGetter domain.GroupMembersGetter,
) error {

	adminService, err := googleclient.GetGoogleAdminService(googleAuthUserEmail, credentialsFilePath)
	if err != nil {
		return err
	}

	groupDiffs, err := InitAllGroupDiffs(correspondingGroups, adminService, sourceMemberGetter)
	if err != nil {
		return fmt.Errorf("Unable to initialize group data:  %s", err)
	}


	for _, nextDiff := range groupDiffs {
		err := googleclient.AddMembersToGroup(nextDiff.TargetGroup, nextDiff.MembersToAdd, adminService)
		if err != nil {
			return fmt.Errorf("Unable to add members to Google group %s:  %s", nextDiff.TargetGroup, err)
		}

		err = googleclient.DeleteMembersFromGroup(nextDiff.TargetGroup, nextDiff.MembersToDelete, adminService)
		if err != nil {
			return fmt.Errorf("Unable to delete members from Google group %s:  %s", nextDiff.TargetGroup, err)
		}
	}

	for _, nextDiff := range groupDiffs {
		log.Printf("\n  Source Group: %s.   Target Group: %s.   Number Added: %d.   Number Deleted: %d  \n",
			nextDiff.SourceGroup, nextDiff.TargetGroup, len(nextDiff.MembersToAdd), len(nextDiff.MembersToDelete))
	}

	return nil
}