package syncgroups

import (
	"testing"

	"github.com/silinternational/serverless-google-groups-sync"
	"github.com/silinternational/serverless-google-groups-sync/lib/testutils"
)

func TestDiffGroup(t *testing.T) {

	type TestCase struct {
		source           []string
		target           []string
		expectedToAdd    []string
		expectedToDelete []string
	}

	testData := []TestCase{
		{ // No changes
			source:           []string{"mm", "aa", "gg"},
			target:           []string{"mm", "aa", "gg"},
			expectedToAdd:    []string{},
			expectedToDelete: []string{},
		},
		{ // Everything is different
			source:           []string{"mm", "aa", "gg"},
			target:           []string{"zz", "xx", "yy"},
			expectedToAdd:    []string{"mm", "aa", "gg"},
			expectedToDelete: []string{"zz", "xx", "yy"},
		},
		{ // some of each
			source:           []string{"mm", "aa", "gg"},
			target:           []string{"mm", "aa", "yy"},
			expectedToAdd:    []string{"gg"},
			expectedToDelete: []string{"yy"},
		},
	}

	for dataIndex, data := range testData {
		nextSet := domain.GroupDiff{
			SourceMembers: data.source,
			TargetMembers: data.target,
		}
		DiffGroup(&nextSet)
		areEqual, errMsg := testutils.AreStringSlicesEqual(data.expectedToAdd, nextSet.MembersToAdd)

		if !areEqual {
			t.Errorf("Error with members-to-add for testData %d ... %s", dataIndex, errMsg)
			//fmt.Printf("\n%v\n%v", data.expectedToAdd, nextSet.MembersToAdd)
			return
		}

		areEqual, errMsg = testutils.AreStringSlicesEqual(data.expectedToDelete, nextSet.MembersToDelete)

		if !areEqual {
			t.Errorf("Error with members-to-delete for testData %d ... %s", dataIndex, errMsg)
		}
	}
}

func TestDiffAllGroups(t *testing.T) {

	type TestCase struct {
		groupDiff        *domain.GroupDiff
		expectedToAdd    []string
		expectedToDelete []string
	}

	testData := []TestCase{
		{ // No changes
			groupDiff: &domain.GroupDiff{
				SourceGroup:   "TestSource",
				TargetGroup:   "TestTarget",
				SourceMembers: []string{"mm", "aa", "gg"},
				TargetMembers: []string{"mm", "aa", "gg"},
			},
			expectedToAdd:    []string{},
			expectedToDelete: []string{},
		},
		{ // Everything is different
			groupDiff: &domain.GroupDiff{
				SourceGroup:   "TestSource",
				TargetGroup:   "TestTarget",
				SourceMembers: []string{"mm", "aa", "gg"},
				TargetMembers: []string{"zz", "xx", "yy"},
			},
			expectedToAdd:    []string{"mm", "aa", "gg"},
			expectedToDelete: []string{"zz", "xx", "yy"},
		},
		{ // some of each
			groupDiff: &domain.GroupDiff{
				SourceGroup:   "TestSource",
				TargetGroup:   "TestTarget",
				SourceMembers: []string{"mm", "aa", "gg"},
				TargetMembers: []string{"mm", "aa", "yy"},
			},
			expectedToAdd:    []string{"gg"},
			expectedToDelete: []string{"yy"},
		},
	}

	allGroups := []*domain.GroupDiff{}

	for _, data := range testData {
		allGroups = append(allGroups, data.groupDiff)
	}

	allGroups = DiffAllGroups(allGroups)

	for dataIndex, data := range testData {
		results := allGroups[dataIndex]
		areEqual, errMsg := testutils.AreStringSlicesEqual(data.expectedToAdd, results.MembersToAdd)

		if !areEqual {
			t.Errorf("Error with members-to-add for testData %d ... %s", dataIndex, errMsg)
			//fmt.Printf("\n%v\n%v\n", data.expectedToAdd, results.MembersToAdd)
			return
		}

		areEqual, errMsg = testutils.AreStringSlicesEqual(data.expectedToDelete, results.MembersToDelete)

		if !areEqual {
			t.Errorf("Error with members-to-delete for testData %d ... %s", dataIndex, errMsg)
		}
	}
}

func TestLoadGroupMapsFromConfig(t *testing.T) {
	appConfig, err := domain.LoadAppConfig("configuration-test.json")
	if err != nil {
		t.Errorf("Unable to load app config, error: %s", err.Error())
	}
	expected := 2
	if len(appConfig.GroupMaps) != expected {
		t.Errorf("did not get expected number of group maps from config, expected %v, got %v", expected, len(appConfig.GroupMaps))
	}
}

// This test performs integration work against Google and is therefore commented out to not run on a regular basis
// func TestSyncGroups(t *testing.T) {
// 	appConfig, err := domain.LoadAppConfig("../../config.json")
// 	if err != nil {
// 		t.Errorf("unable to load app config, error: %s", err.Error())
// 	}
//
// 	group1Data := []string{
//
// 	}
// 	group1ResponseBody, _ := json.Marshal(&group1Data)
//
// 	mux := http.NewServeMux()
// 	server := httptest.NewServer(mux)
//
// 	mux.HandleFunc("/group1", func(w http.ResponseWriter, req *http.Request) {
// 		w.WriteHeader(200)
// 		w.Header().Set("content-type", "application/json")
// 		fmt.Fprintf(w, string(group1ResponseBody))
// 	})
//
// 	appConfig.MemberSourceApiConfig.BaseURL = server.URL
//
// 	err = SyncGroups(appConfig)
// 	if err != nil {
// 		t.Errorf("unable to sync groups, error: %s", err.Error())
// 	}
// }
