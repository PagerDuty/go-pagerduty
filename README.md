# go-pgerduty

go-pagerduty is a CLI and [go](https://golang.org/) client library for [PagerDuty v2 API](https://v2.developer.pagerduty.com/v2/page/api-reference).
[godoc](http://godoc.org/github.com/PagerDuty/go-pagerduty)

## Installation

```
go get github.com/PagerDuty/go-pagerduty
```

## Usage

### CLI

### From golang libraries

```go
package main

import (
	"fmt"
	"github.com/PagerDuty/go-pagerduty"
)

var	authtoken = "" // Set your auth token here

func main() {
	var opts pagerduty.ListEscalationPoliciesOptions
	client := pagerduty.NewClient(authtoken)
	if eps, err := client.ListEscalationPolicies(opts); err != nil {
		panic(err)
	} else {
		for _, p := range eps.EscalationPolicies {
			fmt.Println(p.Name)
		}
	}
}
```

## License
[Apache 2](http://www.apache.org/licenses/LICENSE-2.0)

## Contributing

1. Fork it ( https://github.com/PagerDuty/go-pagerduty/fork )
2. Create your feature branch (`git checkout -b my-new-feature`)
3. Commit your changes (`git commit -am 'Add some feature'`)
4. Push to the branch (`git push origin my-new-feature`)
5. Create a new Pull Request
