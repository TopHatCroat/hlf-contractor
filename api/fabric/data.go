package fabric

import "time"

type User struct {
	MSP     string `json:"msp,omitempty"`
	Email   string `json:"email,omitempty"`
	Balance int    `json:"balance,omitempty"`
	State   string `json:"state,omitempty"`
}

type ChargeTransaction struct {
	Contractor string    `json:"contractor,omitempty"`
	ChargeId   string    `json:"charge_id,omitempty"`
	User       string    `json:"user_email,omitempty"`
	Price      int       `json:"price,omitempty"`
	StartTime  time.Time `json:"start_date,omitempty"`
	EndTime    time.Time `json:"end_date,omitempty"`
	State      string    `json:"state,omitempty"`
}

type ChargePrice struct {
	Price int `json:"price,omitempty"`
}
