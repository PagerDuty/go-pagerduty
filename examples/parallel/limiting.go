package main

import (
	"context"
	"fmt"
	"time"

	"github.com/PagerDuty/go-pagerduty"
)

func main() {
	clientId := "client_id_goes_here"
	clientSecret := "secret_goes_here"
	scopes := []string{"as_account-us.us.companysubdomain", "users.read"}

	// Server-side software would use the default token source
	opts := pagerduty.WithScopedOAuth(clientId, clientSecret, scopes)

	// Terraform provider, CLIs, etc. could use a file token source
	// opts := pagerduty.WithScopedOAuthTokenSource(
	// 	pagerduty.NewFileTokenSource(context.Background(), clientId, clientSecret, scopes, "token.json"),
	// )

	client := pagerduty.NewClient("", opts)

	done := make(chan struct{})
	defer close(done)

	r1, e1 := manyRequests(done, client)
	r2, e2 := manyRequests(done, client)
	r3, e3 := manyRequests(done, client)

	for {
		select {
		case resp := <-r1:
			// Log the response message
			fmt.Printf("[1] Items Returned: %d\n", len(resp.Users))
		case err := <-e1:
			// Log the error message
			fmt.Printf("[1] ERROR: %+v\n", err)
		case resp := <-r2:
			// Log the response message
			fmt.Printf("[2] Items Returned: %d\n", len(resp.Users))
		case err := <-e2:
			// Log the error message
			fmt.Printf("[2] ERROR: %+v\n", err)
		case resp := <-r3:
			// Log the response message
			fmt.Printf("[3] Items Returned: %d\n", len(resp.Users))
		case err := <-e3:
			// Log the error message
			fmt.Printf("[3] ERROR: %+v\n", err)
		}
	}
}

func loggingCloser(channel chan *pagerduty.ListUsersResponse) {
	fmt.Printf("Closing channel \n")
	close(channel)
}

func manyRequests(done <-chan struct{}, client *pagerduty.Client) (<-chan *pagerduty.ListUsersResponse, <-chan error) {
	resp := make(chan *pagerduty.ListUsersResponse)
	errc := make(chan error, 1)

	go func() {
		ctx := context.Background()
		var opts pagerduty.ListUsersOptions

		defer loggingCloser(resp)
		for i := 0; i < 3000; i++ {
			time.Sleep(5 * time.Second)
			users, err := client.ListUsersWithContext(ctx, opts)

			if err != nil {
				errc <- err
				fmt.Printf("Sent an error message\n")
				break
			} else {
				resp <- users
				fmt.Printf("Sent a users message\n")
			}
		}
	}()

	return resp, errc
}
