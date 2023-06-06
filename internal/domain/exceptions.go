package domain

import "fmt"

type WeatherNotFoundError struct {
	City string
}

func (e WeatherNotFoundError) Error() string {
	return fmt.Sprintf("weather with city %s not found", e.City)
}
