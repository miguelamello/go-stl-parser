package main

import (
	"fmt"
	"os"

	"github.com/miguelamello/stl-parser/model"
)

func outputDetails(model model.Model, file *os.File) {
	facets, count, err := model.FacetCounter(file)
	if err != nil {
		fmt.Println("Error reading the STL file:", err)
		return
	}
	area := model.SurfaceArea(facets)
	fmt.Println("Number of facets:", count)
	fmt.Println("Surface area:", area)
}

func main() {

	// Expecting the filepath as a command-line argument.
	if len(os.Args) != 2 {
		fmt.Println("Please provide the filename to be analyzed")
		return
	}

	// Provided filepath must be a valid file
	filepath := os.Args[1]
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("No such file or directory")
		fmt.Println("Please provide a valid filepath")
		return
	}

	defer file.Close()

	// A valid file pointer must be provided
	if file == nil {
		fmt.Println("Error opening file")
		return
	}

	// Check if the file is ASCII or binary
	isAscii, err := model.IsAscii(filepath)
	if err != nil {
		fmt.Println("Error reading the STL file:", err)
		return
	}

	// Decide which model to use
	if isAscii {
		ascii := model.Ascii{}
		outputDetails(ascii, file)
	} else {
		binary := model.Binary{}
		outputDetails(binary, file)
	}

}
