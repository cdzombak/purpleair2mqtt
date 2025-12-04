package main

import (
	"fmt"
	"testing"
)

func TestCalculatePM25AQI(t *testing.T) {
	tests := []struct {
		name          string
		concentration float32
		expectedAQI   int
		expectedCat   string
		expectedColor string
	}{
		{"Good - Low", 5.0, 21, "Good", "Green"},
		{"Good - High", 12.0, 50, "Good", "Green"},
		{"Moderate - Low", 12.1, 51, "Moderate", "Yellow"},
		{"Moderate - Mid", 23.75, 75, "Moderate", "Yellow"},
		{"Moderate - High", 35.4, 100, "Moderate", "Yellow"},
		{"USG - Low", 35.5, 101, "Unhealthy for Sensitive Groups", "Orange"},
		{"USG - High", 55.4, 150, "Unhealthy for Sensitive Groups", "Orange"},
		{"Unhealthy - Low", 55.5, 151, "Unhealthy", "Red"},
		{"Unhealthy - Mid", 100.0, 174, "Unhealthy", "Red"},
		{"Unhealthy - High", 150.4, 200, "Unhealthy", "Red"},
		{"Very Unhealthy - Low", 150.5, 201, "Very Unhealthy", "Purple"},
		{"Very Unhealthy - High", 250.4, 300, "Very Unhealthy", "Purple"},
		{"Hazardous - Low", 250.5, 301, "Hazardous", "Maroon"},
		{"Hazardous - High", 500.4, 500, "Hazardous", "Maroon"},
		{"Beyond AQI", 600.0, 566, "Hazardous", "Maroon"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePM25AQI(tt.concentration)
			if result.AQI != tt.expectedAQI {
				t.Errorf("CalculatePM25AQI(%f) AQI = %d, want %d", tt.concentration, result.AQI, tt.expectedAQI)
			}
			if result.Category != tt.expectedCat {
				t.Errorf("CalculatePM25AQI(%f) Category = %s, want %s", tt.concentration, result.Category, tt.expectedCat)
			}
			if result.Color != tt.expectedColor {
				t.Errorf("CalculatePM25AQI(%f) Color = %s, want %s", tt.concentration, result.Color, tt.expectedColor)
			}
		})
	}
}

func TestCalculatePM10AQI(t *testing.T) {
	tests := []struct {
		name          string
		concentration float32
		expectedAQI   int
		expectedCat   string
		expectedColor string
	}{
		{"Good - Low", 25.0, 23, "Good", "Green"},
		{"Good - High", 54.0, 50, "Good", "Green"},
		{"Moderate - Low", 55.0, 51, "Moderate", "Yellow"},
		{"Moderate - Mid", 100.0, 73, "Moderate", "Yellow"},
		{"Moderate - High", 154.0, 100, "Moderate", "Yellow"},
		{"USG - Low", 155.0, 101, "Unhealthy for Sensitive Groups", "Orange"},
		{"USG - High", 254.0, 150, "Unhealthy for Sensitive Groups", "Orange"},
		{"Unhealthy - Low", 255.0, 151, "Unhealthy", "Red"},
		{"Unhealthy - High", 354.0, 200, "Unhealthy", "Red"},
		{"Very Unhealthy - Low", 355.0, 201, "Very Unhealthy", "Purple"},
		{"Hazardous", 425.0, 301, "Hazardous", "Maroon"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculatePM10AQI(tt.concentration)
			if result.AQI != tt.expectedAQI {
				t.Errorf("CalculatePM10AQI(%f) AQI = %d, want %d", tt.concentration, result.AQI, tt.expectedAQI)
			}
			if result.Category != tt.expectedCat {
				t.Errorf("CalculatePM10AQI(%f) Category = %s, want %s", tt.concentration, result.Category, tt.expectedCat)
			}
			if result.Color != tt.expectedColor {
				t.Errorf("CalculatePM10AQI(%f) Color = %s, want %s", tt.concentration, result.Color, tt.expectedColor)
			}
		})
	}
}

func TestCalculateOverallAQI(t *testing.T) {
	tests := []struct {
		name        string
		pm25        float32
		pm10        float32
		expectedAQI int
		expectedCat string
	}{
		{"Both Good", 10.0, 40.0, 42, "Good"},
		{"PM2.5 Higher", 35.5, 100.0, 101, "Unhealthy for Sensitive Groups"},
		{"PM10 Higher", 20.0, 200.0, 123, "Unhealthy for Sensitive Groups"},
		{"Both Unhealthy", 100.0, 300.0, 174, "Unhealthy"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := CalculateOverallAQI(tt.pm25, tt.pm10)
			if result.AQI != tt.expectedAQI {
				t.Errorf("CalculateOverallAQI(%f, %f) AQI = %d, want %d", tt.pm25, tt.pm10, result.AQI, tt.expectedAQI)
			}
			if result.Category != tt.expectedCat {
				t.Errorf("CalculateOverallAQI(%f, %f) Category = %s, want %s", tt.pm25, tt.pm10, result.Category, tt.expectedCat)
			}
		})
	}
}

// Example function to demonstrate usage
func ExampleCalculatePM25AQI() {
	result := CalculatePM25AQI(35.5)
	fmt.Printf("PM2.5 concentration: 35.5 μg/m³\n")
	fmt.Printf("AQI: %d\n", result.AQI)
	fmt.Printf("Category: %s\n", result.Category)
	fmt.Printf("Color: %s\n", result.Color)
	fmt.Printf("ColorRGB: %s\n", result.ColorRGB)
	// Output:
	// PM2.5 concentration: 35.5 μg/m³
	// AQI: 101
	// Category: Unhealthy for Sensitive Groups
	// Color: Orange
	// ColorRGB: rgb(255,126,0)
}