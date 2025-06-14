package main

type Apart struct {
	ID int64 `json:"id"`
}

type Intercom struct {
	Address string  `json:"address"`
	Aparts  []Apart `json:"aparts"`
}
