package main

import (
	"math"
)

// AQIResult contains the calculated AQI value and associated information
type AQIResult struct {
	AQI            int
	Category       string
	Color          string
	SensitiveGroup string
}

// PM25Breakpoint represents a breakpoint for PM2.5 AQI calculation
type PM25Breakpoint struct {
	ConcLo float32
	ConcHi float32
	AQILo  int
	AQIHi  int
}

// PM10Breakpoint represents a breakpoint for PM10 AQI calculation
type PM10Breakpoint struct {
	ConcLo float32
	ConcHi float32
	AQILo  int
	AQIHi  int
}

// PM2.5 breakpoints based on EPA standards
var pm25Breakpoints = []PM25Breakpoint{
	{0.0, 12.0, 0, 50},       // Good
	{12.1, 35.4, 51, 100},    // Moderate
	{35.5, 55.4, 101, 150},   // Unhealthy for Sensitive Groups
	{55.5, 150.4, 151, 200},  // Unhealthy
	{150.5, 250.4, 201, 300}, // Very Unhealthy
	{250.5, 350.4, 301, 400}, // Hazardous
	{350.5, 500.4, 401, 500}, // Hazardous
}

// PM10 breakpoints based on EPA standards
var pm10Breakpoints = []PM10Breakpoint{
	{0, 54, 0, 50},       // Good
	{55, 154, 51, 100},   // Moderate
	{155, 254, 101, 150}, // Unhealthy for Sensitive Groups
	{255, 354, 151, 200}, // Unhealthy
	{355, 424, 201, 300}, // Very Unhealthy
	{425, 504, 301, 400}, // Hazardous
	{505, 604, 401, 500}, // Hazardous
}

// AQI categories
var aqiCategories = []struct {
	Threshold int
	Category  string
	Color     string
}{
	{50, "Good", "rgb(0,228,0)"},
	{100, "Moderate", "rgb(255,255,0)"},
	{150, "Unhealthy for Sensitive Groups", "rgb(255,126,0)"},
	{200, "Unhealthy", "rgb(255,0,0)"},
	{300, "Very Unhealthy", "rgb(143,63,151)"},
	{500, "Hazardous", "rgb(126,0,35)"},
}

// CalculateAQI calculates the AQI value from a concentration using the provided breakpoints
func calculateAQI(concentration float32, breakpoints interface{}) int {
	var concLo, concHi float32
	var aqiLo, aqiHi int
	found := false

	switch bp := breakpoints.(type) {
	case []PM25Breakpoint:
		// Truncate PM2.5 to 1 decimal place
		concentration = float32(math.Floor(float64(concentration)*10) / 10)
		for _, b := range bp {
			if concentration >= b.ConcLo && concentration <= b.ConcHi {
				concLo = b.ConcLo
				concHi = b.ConcHi
				aqiLo = b.AQILo
				aqiHi = b.AQIHi
				found = true
				break
			}
		}
	case []PM10Breakpoint:
		// Truncate PM10 to integer
		concentration = float32(int(concentration))
		for _, b := range bp {
			if concentration >= b.ConcLo && concentration <= b.ConcHi {
				concLo = b.ConcLo
				concHi = b.ConcHi
				aqiLo = b.AQILo
				aqiHi = b.AQIHi
				found = true
				break
			}
		}
	}

	if !found {
		// Beyond AQI - use highest breakpoint and linear extrapolation
		switch bp := breakpoints.(type) {
		case []PM25Breakpoint:
			last := bp[len(bp)-1]
			concLo = last.ConcLo
			concHi = last.ConcHi
			aqiLo = last.AQILo
			aqiHi = last.AQIHi
		case []PM10Breakpoint:
			last := bp[len(bp)-1]
			concLo = last.ConcLo
			concHi = last.ConcHi
			aqiLo = last.AQILo
			aqiHi = last.AQIHi
		}
	}

	// AQI equation: I = ((IHi - ILo) / (BPHi - BPLo)) * (Cp - BPLo) + ILo
	aqi := ((float32(aqiHi-aqiLo) / (concHi - concLo)) * (concentration - concLo)) + float32(aqiLo)

	// Round to nearest integer
	return int(math.Round(float64(aqi)))
}

// CalculatePM25AQI calculates the AQI from PM2.5 concentration (μg/m³)
func CalculatePM25AQI(concentration float32) AQIResult {
	aqi := calculateAQI(concentration, pm25Breakpoints)

	category := ""
	color := ""

	// Determine category and color
	for _, info := range aqiCategories {
		if aqi <= info.Threshold {
			category = info.Category
			color = info.Color
			break
		}
	}

	// If AQI > 500, it's still Hazardous
	if aqi > 500 {
		category = aqiCategories[len(aqiCategories)-1].Category
		color = aqiCategories[len(aqiCategories)-1].Color
	}

	// Sensitive groups for PM2.5
	sensitiveGroup := ""
	if aqi > 100 {
		sensitiveGroup = "People with heart or lung disease, older adults, children, and people of lower socioeconomic status"
	}

	return AQIResult{
		AQI:            aqi,
		Category:       category,
		Color:          color,
		SensitiveGroup: sensitiveGroup,
	}
}

// CalculatePM10AQI calculates the AQI from PM10 concentration (μg/m³)
func CalculatePM10AQI(concentration float32) AQIResult {
	aqi := calculateAQI(concentration, pm10Breakpoints)

	category := ""
	color := ""

	// Determine category and color
	for _, info := range aqiCategories {
		if aqi <= info.Threshold {
			category = info.Category
			color = info.Color
			break
		}
	}

	// If AQI > 500, it's still Hazardous
	if aqi > 500 {
		category = "Hazardous"
		color = "rgb(126,0,35)"
	}

	// Sensitive groups for PM10
	sensitiveGroup := ""
	if aqi > 100 {
		sensitiveGroup = "People with heart or lung disease, older adults, children, and people of lower socioeconomic status"
	}

	return AQIResult{
		AQI:            aqi,
		Category:       category,
		Color:          color,
		SensitiveGroup: sensitiveGroup,
	}
}

// CalculateOverallAQI calculates the overall AQI (highest of PM2.5 and PM10)
func CalculateOverallAQI(pm25 float32, pm10 float32) AQIResult {
	pm25Result := CalculatePM25AQI(pm25)
	pm10Result := CalculatePM10AQI(pm10)

	// Return the result with the higher AQI
	if pm25Result.AQI >= pm10Result.AQI {
		return pm25Result
	}
	return pm10Result
}
