package model

import (
	"os"
	"bufio"
	"strings"
	"strconv"
	"runtime"
	"sync"
	"errors"
)

type Ascii struct {}

// Return a array of facets found in the STL file 
// and the number of facets read successfully
func (ascii Ascii) FacetCounter(filePointer *os.File) ([]facet, int, error) {

	var counter int
	var facets []facet
	scanner := bufio.NewScanner(filePointer)
	buf := make([]byte, 1024*256) // 256KB buffer (adjustable)
	scanner.Buffer(buf, bufio.MaxScanTokenSize)

	// Start reading the facets
	for scanner.Scan() {

		line := strings.TrimSpace(scanner.Text())

		// Skip the header (the line starting with "solid")
		if strings.HasPrefix(line, "solid") {
			continue
		}

		// Check if reached the end of the file
		if strings.HasPrefix(line, "endsolid") {
			break
		}

		// Check if this line defines a facet
		if strings.HasPrefix(line, "facet normal") {

			counter++
			var facet facet

			// Parse the normal vector
			fields := strings.Fields(line)
			
			if len(fields) != 5 {
				return nil, 0, errors.New("malformed facet normal found")
			}

			// Convert the string values to float64
			for i := 0; i < 3; i++ {
				facet.Normal[i], _ = strconv.ParseFloat(fields[i+2], 64)
			}
			
			// Read the vertices of the facet
			for i := 0; i < 4; i++ {
				if scanner.Scan() {
					vertexLine := strings.TrimSpace(scanner.Text())
					if strings.HasPrefix(vertexLine, "vertex") {
						vertexFields := strings.Fields(vertexLine)
						if len(vertexFields) != 4 {
							return nil, 0, errors.New("malformed vertex found")
						}
						var vertex [3]float64
						for j := 0; j < 3; j++ {
							// Convert the string values to float64
							vertex[j], _ = strconv.ParseFloat(vertexFields[j+1], 64) 
						}
						facet.Vertices = append(facet.Vertices, vertex)
					}
				}
			}

			if len(facet.Vertices) != 3 {
				return nil, 0, errors.New("missing vertex found")
			}

			facets = append(facets, facet)

		}
	}

	if err := scanner.Err(); err != nil {
		return nil, 0, err
	}

	return facets, counter, nil

}

// Using go routines and channels to calculate the surface area of the model
// Spliting the facets into chunks and calculating the area of each chunk in parallel
func (ascii Ascii) SurfaceArea(facets []facet) float64 {

	var totalArea float64
	var wg sync.WaitGroup // WaitGroup to wait for all workers to finish
	areaCh := make(chan float64) // Channel to receive the area of each chunk
	numWorkers := runtime.NumCPU() // Number of CPU cores
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
		totalArea += area
	}

	return totalArea

}

