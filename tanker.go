package main

type Tanker interface {
	GetFuel() float64
	// TODO: rename
	ReduceFuel(val float64)
}
