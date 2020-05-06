package types

import "strings"

// QueryResResolve Queries Result Payload for a resolve query
type QueryResResolve struct {
	Value string `json:"value"`
}

// implement fmt.Stringer
func (r QueryResResolve) String() string {
	return r.Value
}

// QueryResNames Queries Result Payload for a names query
type QueryResNames []string

// implement fmt.Stringer
func (n QueryResNames) String() string {
	return strings.Join(n[:], "\n")
}

// QueryResDescription Queries Result Payload for a description query
type QueryResDescription struct {
	Description string `json:"description"`
}

// implement fmt.Stringer
func (r QueryResDescription) String() string {
	return r.Description
}
