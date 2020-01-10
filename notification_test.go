package pagerduty

import (
	"net/http"
	"testing"
)

// ListNotifications
func TestNotification_List(t *testing.T) {
	setup()
	defer teardown()

	mux.HandleFunc("/notifications", func(w http.ResponseWriter, r *http.Request) {
		testMethod(t, r, "GET")
		w.Write([]byte(`{"notifications": [{"id": "1"}]}`))
	})

	var listObj = APIListObject{Limit: 0, Offset: 0, More: false, Total: 0}
	var client = &Client{apiEndpoint: server.URL, authToken: "foo", HTTPClient: defaultHTTPClient}
	var opts = ListNotificationOptions{
		APIListObject: listObj,
		Includes:      []string{},
		Filter:        "foo",
		Since:         "bar",
		Until:         "baz",
	}
	res, err := client.ListNotifications(opts)

	want := &ListNotificationsResponse{
		APIListObject: listObj,
		Notifications: []Notification{
			{
				ID: "1",
			},
		},
	}

	if err != nil {
		t.Fatal(err)
	}
	testEqual(t, want, res)
}
