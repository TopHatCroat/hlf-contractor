package price

const (
	TypeName = "ChargeTransactionPrice"
)

type Entity struct {
	Price int `json:"price,omitempty"`
}

func (c Entity) Key() ([]string, error) {
	return []string{TypeName}, nil
}
