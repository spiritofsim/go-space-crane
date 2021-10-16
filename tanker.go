package main

type Tank interface {
	GetFuel() float64
	// TODO: rename
	ReduceFuel(val float64)
}
