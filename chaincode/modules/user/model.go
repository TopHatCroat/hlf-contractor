package user

type State string

const (
	TypeName = "User"

	UserStateRegistered State = "registered"
	UserStateBlocked    State = "blocked"
)

type Entity struct {
	Email   string `json:"email,omitempty"`
	Balance int64  `json:"balance,omitempty"`
	State   State  `json:"state,omitempty"`
}

type RestrictedResponse struct {
	Email string `json:"email,omitempty"`
	State State  `json:"state,omitempty"`
}
