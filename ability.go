package pagerduty

// ListAbilityResponse is the response when calling the ListAbility API endpoint.
type ListAbilityResponse struct {
	Abilities []string `json:"abilities"`
}

// ListAbilities lists all abilities on your account.
func (c *PagerdutyClient) ListAbilities() (*ListAbilityResponse, error) {
	resp, err := c.Get("/abilities")
	if err != nil {
		return nil, err
	}
	var result ListAbilityResponse
	return &result, DecodeJSON(resp, &result)
}

// TestAbility Check if your account has the given ability.
func (c *PagerdutyClient) TestAbility(ability string) error {
	_, err := c.Get("/abilities/" + ability)
	return err
}
