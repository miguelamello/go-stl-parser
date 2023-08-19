package model

import (
	"math"
	"os"
	"strings"
	"sync"
)

// Helper function to check if the STL file is binary or ASCII
func IsAscii(filepath string) (bool, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return false, err
	}

	// Read the first 1024 bytes (header) from the STL file
	buf := make([]byte, 1024)
	b, err := file.Read(buf)
	if err != nil {
		return false, err
	}

	// Check for the presence of the keyword "facet" in the buffer
	// If the keyword is found, it's likely a binary STL file
	if strings.Contains(string(buf[:b]), "facet") {
		return true, nil
	}

	return false, nil // No "facet" keyword found, likely an Binary STL file
}

// Function to calculate the chunk area of a set of facets
func chunkArea(facets []facet, wg *sync.WaitGroup, areaCh chan<- float64) {

	defer wg.Done()
	var chunkArea float64

	for _, facet := range facets {
		chunkArea += facetArea(facet)
	}

	areaCh <- chunkArea

}

// Function to calculate the area of a facet
// Formula: Area = 0.5 * √(cp[0]^2 + cp[1]^2 + cp[2]^2)
func facetArea(facet facet) float64 {

	AB := [3]float64{
		facet.Vertices[1][0] - facet.Vertices[0][0],
		facet.Vertices[1][1] - facet.Vertices[0][1],
		facet.Vertices[1][2] - facet.Vertices[0][2],
	}

	AC := [3]float64{
		facet.Vertices[2][0] - facet.Vertices[0][0],
		facet.Vertices[2][1] - facet.Vertices[0][1],
		facet.Vertices[2][2] - facet.Vertices[0][2],
	}

	cp, isValid := crossProduct(AB, AC)

	if !isValid {
		return 0
	}

	// Normalize the cross product vector
	vector := cp[0]*cp[0] + cp[1]*cp[1] + cp[2]*cp[2]
	if math.IsInf(vector, 0) {
		return 0
	}
	return math.Sqrt(vector)

}

// Function to calculate the cross product of two vectors
// Formula: cp = AB x AC
func crossProduct(v1, v2 [3]float64) ([3]float64, bool) {

	cp := [3]float64{
		v1[1]*v2[2] - v1[2]*v2[1],
		v1[2]*v2[0] - v1[0]*v2[2],
		v1[0]*v2[1] - v1[1]*v2[0],
	}

	// Check for NaN (Not a Number) or Inf (Positive Infinity or Negative Infinity)
	for _, value := range cp {
		if math.IsNaN(value) || math.IsInf(value, 0) {
			return [3]float64{0, 0, 0}, false
		}
	}

	return cp, true

}
