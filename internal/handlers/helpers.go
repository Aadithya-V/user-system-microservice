package handlers

import "errors"

var NULLISLAND [2]float64 = [2]float64{0.0, 0.0}

func validateCoordinates(lat, lon float64) error {
	if lat < -85.05112878 || lat > 85.05112878 || lon < -180.0 || lon > 180.0 {
		return errors.New("Invalid Coordinates")
	} else {
		return nil
	}
}
