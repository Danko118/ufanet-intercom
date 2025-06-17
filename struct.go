package main

type Intercom struct {
	Address string  `json:"address"`
	Aparts  []int64 `json:"aparts"`
	Vendor  string  `json:"vendor"`

	MAC string `json:"-"`
}
