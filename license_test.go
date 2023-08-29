package pagerduty

import (
	"context"
	"net/http"
	"testing"
)

func TestLicenses_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/licenses", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{ "licenses": [ { "id": "PIP248G", "name": "Business (Full User)", "description": "Event Intelligence", "current_value": 234, "allocations_available": 4766, "valid_roles": [ "owner", "admin", "user", "limited_user", "observer", "restricted_access" ], "role_group": "FullUser", "summary": "Business (Full User)", "type": "license", "self": null, "html_url": null } ] }`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.ListLicensesWithContext(context.Background())

	want := &ListLicensesResponse{
		Licenses: []License{
			{

				APIObject: APIObject{
					ID:   "PIP248G",
					Type: "license",
				},
				Name:                 "Business (Full User)",
				Description:          "Event Intelligence",
				CurrentValue:         234,
				AllocationsAvailable: 4766,
				ValidRoles: []string{
					"owner",
					"admin",
					"user",
					"limited_user",
					"observer",
					"restricted_access",
				},
				RoleGroup: "FullUser",
				Summary:   "Business (Full User)",
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestLicenseAllocations_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/license_allocations", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		_, _ = w.Write([]byte(`{
			"license_allocations": [
			  {
				"license": {
				  "id": "PIP248G",
				  "name": "Business (Full User)",
				  "description": "Event Intelligence",
				  "valid_roles": [
					"owner",
					"admin",
					"user",
					"limited_user",
					"observer",
					"restricted_access"
				  ],
				  "role_group": "FullUser",
				  "summary": "Business (Full User)",
				  "type": "license",
				  "self": null,
				  "html_url": null
				},
				"user": {
				  "id": "PIP248H",
				  "type": "user_reference"
				},
				"allocated_at": "2021-06-01T00:00:00-05:00"
			  }
			],
			"offset": 0,
			"limit": 1,
			"more": false,
			"total": 2
		  }`))
	})

	client := defaultTestClient(server.URL, "foo")

	res, err := client.ListLicenseAllocationsWithContext(context.Background(), ListLicenseAllocationsOptions{})

	want := &ListLicenseAllocationsResponse{
		APIListObject: APIListObject{
			Offset: 0,
			Limit:  1,
			More:   false,
			Total:  2,
		},
		LicenseAllocations: []LicenseAllocation{
			{
				License: LicenseAllocated{
					APIObject: APIObject{
						ID:   "PIP248G",
						Type: "license",
					},
					Name:        "Business (Full User)",
					Description: "Event Intelligence",
					ValidRoles: []string{
						"owner",
						"admin",
						"user",
						"limited_user",
						"observer",
						"restricted_access",
					},
					RoleGroup: "FullUser",
					Summary:   "Business (Full User)",
				},
				AllocatedAt: "2021-06-01T00:00:00-05:00",
				User: APIObject{
					ID:   "PIP248H",
					Type: "user_reference",
				},
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}
