package model

import (
	"io"
	"os"
	"runtime"
	"sync"
	"encoding/binary"
)

type Binary struct{}

// Return a array of facets found in the STL file
// and the number of facets read successfully
func (_binary Binary) FacetCounter(filePointer *os.File) ([]facet, int, error) {

	var facets []facet

	// Skip the header (80 bytes), which is ignored for binary STL files
	filePointer.Seek(80, io.SeekStart)

	for {
		var facet facet

		// Read the normal vector (12 bytes, 3x float64 for each component)
		err := binary.Read(filePointer, binary.LittleEndian, &facet.Normal)
		if err != nil {
			// Check for EOF (end of file) error, which indicates the end of facets
			if err.Error() == "unexpected EOF" || err == io.EOF {
				break // Exit loop if EOF is encountered
			}
			return nil, 0, err // Return other errors as they are unexpected
		}

		// Read the vertices (36 bytes, 3x float64 for each vertex)
		var vertex [3]float64
		for j := 0; j < 3; j++ {
			err = binary.Read(filePointer, binary.LittleEndian, &vertex)
			if err != nil {
				// Check for EOF (end of file) error, which indicates the end of facets
				if err.Error() == "unexpected EOF" || err == io.EOF {
					break // Exit loop if EOF is encountered
				}
				return nil, 0, err // Return other errors as they are unexpected
			}
			facet.Vertices = append(facet.Vertices, vertex)
		}

		// Skip 2 bytes (attribute byte count) for each facet
		var attributeCount uint16
		err = binary.Read(filePointer, binary.LittleEndian, &attributeCount)
		if err != nil {
			// Check for EOF (end of file) error, which indicates the end of facets
			if err.Error() == "unexpected EOF" || err == io.EOF {
				break // Exit loop if EOF is encountered
			}
			return nil, 0, err // Return other errors as they are unexpected
		}

		facets = append(facets, facet)
	}

	numFacets := len(facets)
	return facets, numFacets, nil

}

// Using go routines and channels to calculate the surface area of the model
// Spliting the facets into chunks and calculating the area of each chunk in parallel
func (_binary Binary) SurfaceArea(facets []facet) float64 {

	var totalArea float64
	var wg sync.WaitGroup                                    // WaitGroup to wait for all workers to finish
	areaCh := make(chan float64)                             // Channel to receive the area of each chunk
	numWorkers := runtime.NumCPU()                           // Number of CPU cores
	chunkSize := (len(facets) + numWorkers - 1) / numWorkers // Even distribution of workers

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		startIdx := i * chunkSize
		endIdx := startIdx + chunkSize

		if i == numWorkers-1 { // Last worker takes any remaining facets
			endIdx = len(facets)
		}

		// Function to calculate the chunk area of a set of facets
		go chunkArea(facets[startIdx:endIdx], &wg, areaCh)

		if endIdx == len(facets) {
			break
		}
	}

	go func() {
		wg.Wait()
		close(areaCh)
	}()

	for area := range areaCh {
		totalArea += float64(area)
	}

	return totalArea

}

