package domain

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestIsStringInStringSlice_True(t *testing.T) {
	type TestCase struct {
		needle   string
		haystack []string
	}

	testData := []TestCase{
		{"a", []string{"a"}},
		{"a", []string{"b", "c", "a"}},
	}

	for index, data := range testData {
		results := IsStringInStringSlice(data.needle, data.haystack)

		if !results {
			errStart := fmt.Sprintf("At testData index %d, the string was not found in the slice but should have been.\n", index)
			t.Errorf("%s   String: %s.\n   Slice: %v", errStart, data.needle, data.haystack)
			return
		}
	}
}

func TestIsStringInStringSlice_False(t *testing.T) {
	type TestCase struct {
		needle   string
		haystack []string
	}

	testData := []TestCase{
		{"a", []string{}},
		{"a", []string{"b", "c"}},
	}

	for index, data := range testData {
		results := IsStringInStringSlice(data.needle, data.haystack)

		if results {
			errStart := fmt.Sprintf("At testData index %d, the string was found in the slice but should NOT have been.\n", index)
			t.Errorf("%s   String: %s.\n   Slice: %v", errStart, data.needle, data.haystack)
			return
		}
	}
}

func TestGetGroupMembersFromSource(t *testing.T) {
	group1Data := []string{
		"user1@domain.com", "user2@domain.com", "user3@domain.com",
	}
	group1ResponseBody, _ := json.Marshal(&group1Data)

	group2Data := []string{
		"user1@domain.com", "user2@domain.com",
	}
	group2ResponseBody, _ := json.Marshal(&group2Data)

	mux := http.NewServeMux()
	server := httptest.NewServer(mux)

	mux.HandleFunc("/group1", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("content-type", "application/json")
		fmt.Fprintf(w, string(group1ResponseBody))
	})

	mux.HandleFunc("/group2", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
		w.Header().Set("content-type", "application/json")
		fmt.Fprintf(w, string(group2ResponseBody))
	})

	apiConfig := MemberSourceApiConfig{
		BaseURL: server.URL,
		User:    "test",
		Pass:    "test",
	}

	group1members, err := GetGroupMembersFromSource(apiConfig, "/group1")
	if err != nil {
		t.Errorf("failed to get group1 members, error: %s", err.Error())
	}

	if len(group1members) != len(group1Data) {
		t.Errorf("group1 member response does not mach number of results in test data, got: %v", group1members)
	}

	group2members, err := GetGroupMembersFromSource(apiConfig, "/group2")
	if err != nil {
		t.Errorf("failed to get group2 members, error: %s", err.Error())
	}

	if len(group2members) != len(group2Data) {
		t.Errorf("group2 member response does not mach number of results in test data, got: %v", group2members)
	}
}
