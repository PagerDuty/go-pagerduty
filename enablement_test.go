package pagerduty

import (
	"context"
	"net/http"
	"testing"
)

// Test data for enablements
var (
	testEnablementFeatureAIOps = "aiops"
	testServiceID              = "PABCDEF"
	testEventOrchestrationID   = "1b49abe7-26db-4439-a715-c6d883acfb3e"
)

// Mock response data based on API examples
const (
	mockServiceEnablementsSuccessResponse = `{
		"enablements": [
			{
				"feature": "aiops",
				"enabled": true
			}
		]
	}`

	mockServiceEnablementsWithWarningResponse = `{
		"enablements": [
			{
				"feature": "aiops",
				"enabled": true
			}
		],
		"warnings": ["Your account is not entitled to use AIOps features for this Service."]
	}`

	mockServiceEnablementsDisabledResponse = `{
		"enablements": [
			{
				"feature": "aiops",
				"enabled": false
			}
		]
	}`

	mockEventOrchestrationEnablementsSuccessResponse = `{
		"enablements": [
			{
				"feature": "aiops",
				"enabled": true
			}
		]
	}`

	mockEventOrchestrationEnablementsWithWarningResponse = `{
		"enablements": [
			{
				"feature": "aiops",
				"enabled": true
			}
		],
		"warnings": ["You can't use AIOps functionality with this Orchestration because your account hasn't purchased AIOps"]
	}`

	// Service enablement update response (uses list format)
	mockServiceUpdateEnablementSuccessResponse = `{
		"enablements": [
			{
				"feature": "aiops",
				"enabled": false
			}
		]
	}`

	mockServiceUpdateEnablementEnabledResponse = `{
		"enablements": [
			{
				"feature": "aiops",
				"enabled": true
			}
		]
	}`

	// Event orchestration enablement update response (uses single format)
	mockEventOrchestrationUpdateEnablementSuccessResponse = `{
		"enablement": {
			"feature": "aiops",
			"enabled": false
		}
	}`

	mockEventOrchestrationUpdateEnablementEnabledResponse = `{
		"enablement": {
			"feature": "aiops",
			"enabled": true
		}
	}`

	mockErrorResponse400 = `{
		"error": {
			"code": 2001,
			"message": "Invalid Input Provided",
			"errors": ["The feature name 'invalid' is not supported"]
		}
	}`

	mockErrorResponse403 = `{
		"error": {
			"code": 2003,
			"message": "Forbidden",
			"errors": ["You are not authorized to access this resource"]
		}
	}`

	mockErrorResponse404 = `{
		"error": {
			"code": 2100,
			"message": "Not Found",
			"errors": ["The specified service does not exist"]
		}
	}`

	mockErrorResponse500 = `{
		"error": {
			"code": 3001,
			"message": "Internal Server Error",
			"errors": ["An unexpected error occurred"]
		}
	}`
)

// ListServiceEnablements - Success Response
func TestService_ListEnablements(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/"+testServiceID+"/enablements", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockServiceEnablementsSuccessResponse))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.ListServiceEnablementsWithContext(context.Background(), testServiceID)

	want := &EnablementsWithWarnings{
		Enablements: []Enablement{
			{
				Feature: testEnablementFeatureAIOps,
				Enabled: true,
			},
		},
		Warnings: nil,
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// ListServiceEnablements - Success Response with Warning (warnings are now returned in result)
func TestService_ListEnablementsWithWarning(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/"+testServiceID+"/enablements", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockServiceEnablementsWithWarningResponse))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.ListServiceEnablementsWithContext(context.Background(), testServiceID)

	want := &EnablementsWithWarnings{
		Enablements: []Enablement{
			{
				Feature: testEnablementFeatureAIOps,
				Enabled: true,
			},
		},
		Warnings: []string{"Your account is not entitled to use AIOps features for this Service."},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// ListServiceEnablements - Disabled Response
func TestService_ListEnablementsDisabled(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/"+testServiceID+"/enablements", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockServiceEnablementsDisabledResponse))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.ListServiceEnablementsWithContext(context.Background(), testServiceID)

	want := &EnablementsWithWarnings{
		Enablements: []Enablement{
			{
				Feature: testEnablementFeatureAIOps,
				Enabled: false,
			},
		},
		Warnings: nil,
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// ListServiceEnablements - 403 Forbidden Error
func TestService_ListEnablements403Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/"+testServiceID+"/enablements", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(mockErrorResponse403))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.ListServiceEnablementsWithContext(context.Background(), testServiceID)

	if !testErrCheck(t, "ListServiceEnablements", "access forbidden", err) {
		return
	}
}

// ListServiceEnablements - 404 Not Found Error
func TestService_ListEnablements404Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/"+testServiceID+"/enablements", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(mockErrorResponse404))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.ListServiceEnablementsWithContext(context.Background(), testServiceID)

	if !testErrCheck(t, "ListServiceEnablements", "not found", err) {
		return
	}
}

// ListServiceEnablements - 500 Internal Server Error
func TestService_ListEnablements500Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/"+testServiceID+"/enablements", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(mockErrorResponse500))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.ListServiceEnablementsWithContext(context.Background(), testServiceID)

	if !testErrCheck(t, "ListServiceEnablements", "server error", err) {
		return
	}
}

// UpdateServiceEnablement - Success Response (Disable)
func TestService_UpdateEnablement(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/"+testServiceID+"/enablements/"+testEnablementFeatureAIOps, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockServiceUpdateEnablementSuccessResponse))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.UpdateServiceEnablementWithContext(context.Background(), testServiceID, testEnablementFeatureAIOps, false)

	want := &EnablementWithWarnings{
		Enablement: &Enablement{
			Feature: testEnablementFeatureAIOps,
			Enabled: false,
		},
		Warnings: nil,
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// UpdateServiceEnablement - Success Response (Enable)
func TestService_UpdateEnablementEnable(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/"+testServiceID+"/enablements/"+testEnablementFeatureAIOps, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockServiceUpdateEnablementEnabledResponse))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.UpdateServiceEnablementWithContext(context.Background(), testServiceID, testEnablementFeatureAIOps, true)

	want := &EnablementWithWarnings{
		Enablement: &Enablement{
			Feature: testEnablementFeatureAIOps,
			Enabled: true,
		},
		Warnings: nil,
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// UpdateServiceEnablement - Invalid Feature Error
func TestService_UpdateEnablementInvalidFeature(t *testing.T) {
	setup()
	defer teardown()

	client := defaultTestClient(server.URL, "foo")
	_, err := client.UpdateServiceEnablementWithContext(context.Background(), testServiceID, "invalid_feature", false)

	if !testErrCheck(t, "UpdateServiceEnablement", "unsupported feature", err) {
		return
	}
}

// UpdateServiceEnablement - Empty Service ID Error
func TestService_UpdateEnablementEmptyServiceID(t *testing.T) {
	setup()
	defer teardown()

	client := defaultTestClient(server.URL, "foo")
	_, err := client.UpdateServiceEnablementWithContext(context.Background(), "", testEnablementFeatureAIOps, false)

	if !testErrCheck(t, "UpdateServiceEnablement", "entity ID cannot be empty", err) {
		return
	}
}

// UpdateServiceEnablement - 403 Forbidden Error
func TestService_UpdateEnablement403Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/"+testServiceID+"/enablements/"+testEnablementFeatureAIOps, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(mockErrorResponse403))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.UpdateServiceEnablementWithContext(context.Background(), testServiceID, testEnablementFeatureAIOps, false)

	if !testErrCheck(t, "UpdateServiceEnablement", "access forbidden", err) {
		return
	}
}

// UpdateServiceEnablement - 404 Not Found Error
func TestService_UpdateEnablement404Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/"+testServiceID+"/enablements/"+testEnablementFeatureAIOps, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(mockErrorResponse404))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.UpdateServiceEnablementWithContext(context.Background(), testServiceID, testEnablementFeatureAIOps, false)

	if !testErrCheck(t, "UpdateServiceEnablement", "not found", err) {
		return
	}
}

// UpdateServiceEnablement - 500 Internal Server Error
func TestService_UpdateEnablement500Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/"+testServiceID+"/enablements/"+testEnablementFeatureAIOps, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(mockErrorResponse500))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.UpdateServiceEnablementWithContext(context.Background(), testServiceID, testEnablementFeatureAIOps, false)

	if !testErrCheck(t, "UpdateServiceEnablement", "server error", err) {
		return
	}
}

// ListEventOrchestrationEnablements - Success Response
func TestEventOrchestration_ListEnablements(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/"+testEventOrchestrationID+"/enablements", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockEventOrchestrationEnablementsSuccessResponse))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.ListEventOrchestrationEnablementsWithContext(context.Background(), testEventOrchestrationID)

	want := &EnablementsWithWarnings{
		Enablements: []Enablement{
			{
				Feature: testEnablementFeatureAIOps,
				Enabled: true,
			},
		},
		Warnings: nil,
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// ListEventOrchestrationEnablements - Success Response with Warning
func TestEventOrchestration_ListEnablementsWithWarning(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/"+testEventOrchestrationID+"/enablements", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockEventOrchestrationEnablementsWithWarningResponse))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.ListEventOrchestrationEnablementsWithContext(context.Background(), testEventOrchestrationID)

	want := &EnablementsWithWarnings{
		Enablements: []Enablement{
			{
				Feature: testEnablementFeatureAIOps,
				Enabled: true,
			},
		},
		Warnings: []string{"You can't use AIOps functionality with this Orchestration because your account hasn't purchased AIOps"},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// ListEventOrchestrationEnablements - 403 Forbidden Error
func TestEventOrchestration_ListEnablements403Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/"+testEventOrchestrationID+"/enablements", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(mockErrorResponse403))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.ListEventOrchestrationEnablementsWithContext(context.Background(), testEventOrchestrationID)

	if !testErrCheck(t, "ListEventOrchestrationEnablements", "access forbidden", err) {
		return
	}
}

// ListEventOrchestrationEnablements - 404 Not Found Error
func TestEventOrchestration_ListEnablements404Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/"+testEventOrchestrationID+"/enablements", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(mockErrorResponse404))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.ListEventOrchestrationEnablementsWithContext(context.Background(), testEventOrchestrationID)

	if !testErrCheck(t, "ListEventOrchestrationEnablements", "not found", err) {
		return
	}
}

// ListEventOrchestrationEnablements - 500 Internal Server Error
func TestEventOrchestration_ListEnablements500Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/"+testEventOrchestrationID+"/enablements", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(mockErrorResponse500))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.ListEventOrchestrationEnablementsWithContext(context.Background(), testEventOrchestrationID)

	if !testErrCheck(t, "ListEventOrchestrationEnablements", "server error", err) {
		return
	}
}

// UpdateEventOrchestrationEnablement - Success Response
func TestEventOrchestration_UpdateEnablement(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/"+testEventOrchestrationID+"/enablements/"+testEnablementFeatureAIOps, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockEventOrchestrationUpdateEnablementSuccessResponse))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.UpdateEventOrchestrationEnablementWithContext(context.Background(), testEventOrchestrationID, testEnablementFeatureAIOps, false)

	want := &EnablementWithWarnings{
		Enablement: &Enablement{
			Feature: testEnablementFeatureAIOps,
			Enabled: false,
		},
		Warnings: []string{},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// UpdateEventOrchestrationEnablement - Success Response (Enable)
func TestEventOrchestration_UpdateEnablementEnable(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/"+testEventOrchestrationID+"/enablements/"+testEnablementFeatureAIOps, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockEventOrchestrationUpdateEnablementEnabledResponse))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.UpdateEventOrchestrationEnablementWithContext(context.Background(), testEventOrchestrationID, testEnablementFeatureAIOps, true)

	want := &EnablementWithWarnings{
		Enablement: &Enablement{
			Feature: testEnablementFeatureAIOps,
			Enabled: true,
		},
		Warnings: []string{},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

// UpdateEventOrchestrationEnablement - Invalid Feature Error
func TestEventOrchestration_UpdateEnablementInvalidFeature(t *testing.T) {
	setup()
	defer teardown()

	client := defaultTestClient(server.URL, "foo")
	_, err := client.UpdateEventOrchestrationEnablementWithContext(context.Background(), testEventOrchestrationID, "invalid_feature", false)

	if !testErrCheck(t, "UpdateEventOrchestrationEnablement", "unsupported feature", err) {
		return
	}
}

// UpdateEventOrchestrationEnablement - Empty Orchestration ID Error
func TestEventOrchestration_UpdateEnablementEmptyOrchestrationID(t *testing.T) {
	setup()
	defer teardown()

	client := defaultTestClient(server.URL, "foo")
	_, err := client.UpdateEventOrchestrationEnablementWithContext(context.Background(), "", testEnablementFeatureAIOps, false)

	if !testErrCheck(t, "UpdateEventOrchestrationEnablement", "entity ID cannot be empty", err) {
		return
	}
}

// UpdateEventOrchestrationEnablement - 403 Forbidden Error
func TestEventOrchestration_UpdateEnablement403Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/"+testEventOrchestrationID+"/enablements/"+testEnablementFeatureAIOps, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusForbidden)
		_, _ = w.Write([]byte(mockErrorResponse403))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.UpdateEventOrchestrationEnablementWithContext(context.Background(), testEventOrchestrationID, testEnablementFeatureAIOps, false)

	if !testErrCheck(t, "UpdateEventOrchestrationEnablement", "access forbidden", err) {
		return
	}
}

// UpdateEventOrchestrationEnablement - 404 Not Found Error
func TestEventOrchestration_UpdateEnablement404Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/"+testEventOrchestrationID+"/enablements/"+testEnablementFeatureAIOps, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		_, _ = w.Write([]byte(mockErrorResponse404))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.UpdateEventOrchestrationEnablementWithContext(context.Background(), testEventOrchestrationID, testEnablementFeatureAIOps, false)

	if !testErrCheck(t, "UpdateEventOrchestrationEnablement", "not found", err) {
		return
	}
}

// UpdateEventOrchestrationEnablement - 500 Internal Server Error
func TestEventOrchestration_UpdateEnablement500Error(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/"+testEventOrchestrationID+"/enablements/"+testEnablementFeatureAIOps, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(mockErrorResponse500))
	})

	client := defaultTestClient(server.URL, "foo")
	_, err := client.UpdateEventOrchestrationEnablementWithContext(context.Background(), testEventOrchestrationID, testEnablementFeatureAIOps, false)

	if !testErrCheck(t, "UpdateEventOrchestrationEnablement", "server error", err) {
		return
	}
}

// Context-based tests
func TestService_ListEnablementsWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/"+testServiceID+"/enablements", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockServiceEnablementsSuccessResponse))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.ListServiceEnablementsWithContext(context.Background(), testServiceID)

	want := &EnablementsWithWarnings{
		Enablements: []Enablement{
			{
				Feature: testEnablementFeatureAIOps,
				Enabled: true,
			},
		},
		Warnings: nil,
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestService_UpdateEnablementWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/services/"+testServiceID+"/enablements/"+testEnablementFeatureAIOps, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockServiceUpdateEnablementSuccessResponse))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.UpdateServiceEnablementWithContext(context.Background(), testServiceID, testEnablementFeatureAIOps, false)

	want := &EnablementWithWarnings{
		Enablement: &Enablement{
			Feature: testEnablementFeatureAIOps,
			Enabled: false,
		},
		Warnings: nil,
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestEventOrchestration_ListEnablementsWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/"+testEventOrchestrationID+"/enablements", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockEventOrchestrationEnablementsSuccessResponse))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.ListEventOrchestrationEnablementsWithContext(context.Background(), testEventOrchestrationID)

	want := &EnablementsWithWarnings{
		Enablements: []Enablement{
			{
				Feature: testEnablementFeatureAIOps,
				Enabled: true,
			},
		},
		Warnings: nil,
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}

func TestEventOrchestration_UpdateEnablementWithContext(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/event_orchestrations/"+testEventOrchestrationID+"/enablements/"+testEnablementFeatureAIOps, func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "PUT")
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(mockEventOrchestrationUpdateEnablementSuccessResponse))
	})

	client := defaultTestClient(server.URL, "foo")
	res, err := client.UpdateEventOrchestrationEnablementWithContext(context.Background(), testEventOrchestrationID, testEnablementFeatureAIOps, false)

	want := &EnablementWithWarnings{
		Enablement: &Enablement{
			Feature: testEnablementFeatureAIOps,
			Enabled: false,
		},
		Warnings: []string{},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}
