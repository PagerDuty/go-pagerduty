package pagerduty

import (
	"bytes"
	"fmt"
	"net/http"
	"testing"
	"text/template"
)

// ListTeams
func TestTeam_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"teams": [{"id": "1"}]}`))
	})

	listObj := APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	client := defaultTestClient(server.URL, "foo")
	opts := ListTeamOptions{
		APIListObject: listObj,
		Query:         "foo",
	}
	res, err := client.ListTeams(opts)

	want := &ListTeamResponse{
		APIListObject: listObj,
		Teams: []Team{
			{
				APIObject: APIObject{
					ID: "1",
				},
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Create Team
func TestTeam_Create(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "POST")
		_, _ = w.Write([]byte(`{"team": {"id": "1","name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	input := &Team{
		Name: "foo",
	}
	res, err := client.CreateTeam(input)

	want := &Team{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Delete Team
func TestTeam_Delete(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	client := defaultTestClient(server.URL, "foo")
	id := "1"
	err := client.DeleteTeam(id)
	if err != nil {
		t.Fatal(err)
	}
}

// Get Team
func TestTeam_Get(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{"team": {"id": "1","name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")
	id := "1"
	res, err := client.GetTeam(id)

	want := &Team{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Update Team
func TestTeam_Update(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		_, _ = w.Write([]byte(`{"team": {"id": "1","name":"foo"}}`))
	})

	client := defaultTestClient(server.URL, "foo")

	input := &Team{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}
	id := "1"
	res, err := client.UpdateTeam(id, input)

	want := &Team{
		APIObject: APIObject{
			ID: "1",
		},
		Name: "foo",
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// Remove Escalation Policy from Team
func TestTeam_RemoveEscalationPolicyFromTeam(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1/escalation_policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	client := defaultTestClient(server.URL, "foo")
	teamID := "1"
	epID := "1"

	err := client.RemoveEscalationPolicyFromTeam(teamID, epID)
	if err != nil {
		t.Fatal(err)
	}
}

// Add Escalation Policy to Team
func TestTeam_AddEscalationPolicyToTeam(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1/escalation_policies/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	client := defaultTestClient(server.URL, "foo")
	teamID := "1"
	epID := "1"

	err := client.AddEscalationPolicyToTeam(teamID, epID)
	if err != nil {
		t.Fatal(err)
	}
}

// Remove User from Team
func TestTeam_RemoveUserFromTeam(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1/users/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "DELETE")
	})

	client := defaultTestClient(server.URL, "foo")
	teamID := "1"
	userID := "1"

	err := client.RemoveUserFromTeam(teamID, userID)
	if err != nil {
		t.Fatal(err)
	}
}

// Add User to Team
func TestTeam_AddUserToTeam(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/teams/1/users/1", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
	})

	client := defaultTestClient(server.URL, "foo")
	teamID := "1"
	userID := "1"

	err := client.AddUserToTeam(teamID, userID)
	if err != nil {
		t.Fatal(err)
	}
}

func userID(offset, index int) int {
	return offset + index
}

func lastIndex(length, index int) bool {
	return length-1 == index
}

const membersResponseTemplate = `
{
    {{template "pageInfo" . }},
    "members": [
        {{- $length := len .roles -}}
        {{- $offset := .offset -}}
        {{- range $index, $role := .roles -}}
        {
            "user": {
                "id": "ID{{userID $offset $index}}"
            },
            "role": "{{ $role }}"
        }
        {{- if not (lastIndex $length $index) }},
        {{end -}}
        {{- end }}
    ]
}
`

var memberPageTemplate = template.Must(pageTemplate.New("membersResponse").
	Funcs(templateUtilityFuncs).
	Parse(membersResponseTemplate))

const (
	testValidTeamID = "MYTEAM"
	testAPIKey      = "MYKEY"
	testBadURL      = "A-FAKE-URL"
	testMaxPageSize = 3
)

var templateUtilityFuncs = template.FuncMap{
	"lastIndex": lastIndex,
	"userID":    userID,
}

var pageTemplate = template.Must(template.New("pageInfo").Parse(`
    "more": {{- .more -}},
    "limit": {{- .limit -}},
    "offset": {{- .offset -}}
`))

type pageDetails struct {
	lowNumber, highNumber, limit, offset int
	more                                 bool
}

func genMembersRespPage(details pageDetails, t *testing.T) string {
	if details.limit == 0 {
		details.limit = 25 // Default to 25, PD's API default.
	}

	possibleRoles := []string{"manager", "responder", "observer"}
	roles := make([]string, 0)
	for ; details.lowNumber <= details.highNumber; details.lowNumber++ {
		roles = append(roles, possibleRoles[details.lowNumber%len(possibleRoles)])
	}

	buffer := bytes.NewBufferString("")
	err := memberPageTemplate.Execute(buffer, map[string]interface{}{
		"roles":  roles,
		"more":   details.more,
		"limit":  details.limit,
		"offset": details.offset,
	})
	if err != nil {
		t.Fatalf("Failed to apply values to template: %v", err)
	}

	return string(buffer.String())
}

func genRespPages(amount,
	maxPageSize int,
	pageGenerator func(pageDetails, *testing.T) string,
	t *testing.T) []string {
	pages := make([]string, 0)

	lowNumber := 1
	offset := 0
	more := true

	for {
		tempHighNumber := amount

		if lowNumber+(maxPageSize-1) < amount {
			// Still more pages to come, this page doesn't hit upper.
			tempHighNumber = lowNumber + (maxPageSize - 1)
		} else {
			// Last page, with at least 1 user.
			more = false
		}

		// Generate page using current lower and upper.
		page := pageGenerator(pageDetails{
			lowNumber:  lowNumber,
			highNumber: tempHighNumber,
			limit:      maxPageSize,
			more:       more,
			offset:     offset,
		}, t)

		pages = append(pages, page)

		if !more {
			// Hit the last page, stop.
			return pages
		}
		// Move the offset and lower up to prepare for next page.
		offset += maxPageSize
		lowNumber += maxPageSize
	}
}

func TestListMembersSuccess(t *testing.T) {
	setup()
	defer teardown()

	expectedNumResults := testMaxPageSize - 1
	page := genRespPages(expectedNumResults, testMaxPageSize, genMembersRespPage, t)[0]

	mux.HandleFunc("/teams/"+testValidTeamID+"/members", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, page)
	})

	api := defaultTestClient(server.URL, testAPIKey)
	members, err := api.ListMembers(testValidTeamID, ListMembersOptions{})
	if err != nil {
		t.Fatalf("Failed to get members: %v", err)
	}

	if len(members.Members) != expectedNumResults {
		t.Fatalf("Expected %d team members, got: %v", expectedNumResults, err)
	}
}

func TestListMembersError(t *testing.T) {
	api := defaultTestClient(testBadURL, testAPIKey)
	members, err := api.ListMembers(testValidTeamID, ListMembersOptions{})
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	if members != nil {
		t.Fatalf("Expected nil members response, got: %v", members)
	}
}

func TestListAllMembersSuccessMultiplePages(t *testing.T) {
	setup()
	defer teardown()

	expectedNumResults := testMaxPageSize*3 + 1
	currentPage := 0
	pages := genRespPages(expectedNumResults, testMaxPageSize, genMembersRespPage, t)

	mux.HandleFunc("/teams/"+testValidTeamID+"/members", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, pages[currentPage])
		currentPage++
	})

	api := defaultTestClient(server.URL, testAPIKey)

	members, err := api.ListAllMembers(testValidTeamID)
	if err != nil {
		t.Fatalf("Failed to get members: %v", err)
	}

	if len(members) != expectedNumResults {
		t.Fatalf("Expected %d team members, got: %v", expectedNumResults, err)
	}
}

func TestListAllMembersError(t *testing.T) {
	api := defaultTestClient(testBadURL, testAPIKey)
	members, err := api.ListAllMembers(testValidTeamID)
	if err == nil {
		t.Fatalf("Expected error, got nil")
	}
	if len(members) != 0 {
		t.Fatalf("Expected 0 members, got: %v", members)
	}
}
