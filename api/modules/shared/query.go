package shared

import (
	"encoding/json"
	"github.com/pkg/errors"
)

var (
	defaultRange = []int{0, 20}
)

type QueryParams struct {
	Filter map[string]string
	Range  []int
	Sort   []string
}

func ParseQueryParams(query string) (QueryParams, error) {
	params := QueryParams{}
	err := json.Unmarshal([]byte(query), &params)
	if err != nil {
		return params, errors.Wrap(err, "failed to parse query params")
	}

	return params, nil
}

func (qp *QueryParams) GetFilter(key string) string {
	res, ok := qp.Filter[key]
	if !ok {
		return ""
	}

	return res
}

func (qp *QueryParams) GetRange() []int {
	if qp.Range == nil {
		return defaultRange
	}

	return qp.Range
}

func (qp *QueryParams) GetSort() []string {
	if qp.Sort == nil {
		return []string{}
	}

	return qp.Sort
}
