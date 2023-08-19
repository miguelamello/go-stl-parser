package model

import (
	"os"
)

type Model interface {
	FacetCounter(filePointer *os.File) ([]facet, int, error)
	SurfaceArea(facets []facet) float64
}