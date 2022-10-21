package main

type Flower struct {
	Params []float64 `json:"params"`
	Name   string    `json:"name"`
}

type Flowers struct {
	Fl []Flower
}
