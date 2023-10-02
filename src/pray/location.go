package pray

type Location struct {
	Latitude  float64
	Longitude float64
}

func GetLocation() Location {
	return Location{Latitude: 0.0, Longitude: 0.0}
}