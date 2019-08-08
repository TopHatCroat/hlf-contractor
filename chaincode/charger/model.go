package main

import "time"

type State string

const (
	TypeName = "ChargeTransaction"

	ChargeStateStarted   State = "started"
	ChargeStateStopped   State = "stopped"
	ChargeStateCompleted State = "completed"
	ChargeStateRejected  State = "rejected"
)

type Entity struct {
	Contractor string    `json:"contractor,omitempty"`
	ChargeId   string    `json:"charge_id,omitempty"`
	UserEmail  string    `json:"user_email,omitempty"`
	Price      string    `json:"price,omitempty"`
	StartTime  time.Time `json:"start_date,omitempty"`
	EndTime    time.Time `json:"end_date,omitempty"`
	State      State     `json:"state,omitempty"`
}

type StartTransaction struct {
	Contractor string `json:"contractor,omitempty"`
}

type StopTransaction struct {
	Contractor string `json:"contractor,omitempty"`
	ChargeId   string `json:"charge_id,omitempty"`
}

type CompleteTransaction StopTransaction
