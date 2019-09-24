package main

type State string

const (
	TypeName = "User"

	UserStateActive  State = "active"
	UserStateBlocked State = "blocked"
)

type Entity struct {
	MSP     string `json:"msp,omitempty"`
	Email   string `json:"email,omitempty"`
	Balance int    `json:"balance,omitempty"`
	State   State  `json:"state,omitempty"`
}

func (c Entity) Key() ([]string, error) {
	return []string{TypeName, c.MSP, c.Email}, nil
}

type RestrictedResponse struct {
	MSP   string `json:"msp,omitempty"`
	Email string `json:"email,omitempty"`
	State State  `json:"state,omitempty"`
}
