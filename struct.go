package main

type Intercom struct {
	Address string  `json:"address,omitempty"`
	Aparts  []int64 `json:"aparts,omitempty"`
	Vendor  string  `json:"vendor,omitempty"`
	Status  string  `json:"status,omitempty"`
	MAC     string  `json:"mac"`
}

type Event struct {
	Name string `json:"event,omitempty"`
	Args int    `json:"arg,omitempty"`
	Desc string `json:"desc,omitempty"`
	MAC  string `json:"mac,omitempty"`
}
