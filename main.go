package main

import (
	"fmt"
	"os"

	parser "./parser"
)

func main() {

	//take input from ENV
	searchBy := os.Getenv("searchby")
	searchByPostCode := os.Getenv("postcode")
	searchByDeliveryTime := os.Getenv("deliverytime")

	//mounted recipe data
	file := "./input.json"

	// create parser, validate input, parse by streaming content and generate the results into stderr
	recipeparser := parser.NewParser()
	err := recipeparser.ValidateInput(&file, &searchBy, &searchByPostCode, &searchByDeliveryTime)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	} else {
		err := recipeparser.Parse()
		if err != nil {
			fmt.Fprintln(os.Stderr, "parsing error")
		}
		err = recipeparser.GenerateResult()
		if err != nil {
			fmt.Fprintln(os.Stderr, "error while generating results")
		}
	}
}
