package model

type Parameters struct {
	Name   string
	Values []string
}

type Param struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}
