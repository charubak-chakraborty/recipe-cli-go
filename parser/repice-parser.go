package parser

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/bcicen/jstream"
)

type parser struct{}

var (
	//hashmap to track recipe and postcode counts
	recipeMap   map[string]int
	postcodeMap map[string]int

	//input filepath, search parameters
	filePath             *string
	searchByRecipeTag    *string
	searchByPostCode     *string
	searchByDeliveryTime *string

	//count to track feasible deliveries in the searched postcode
	feasibleDeliveryCount int
)

func NewParser() Parser {

	//init hashmaps
	recipeMap = make(map[string]int)
	postcodeMap = make(map[string]int)
	return &parser{}
}

//ValidateInput ... basic validations
func (p *parser) ValidateInput(filepath *string, searchbyrecipetag *string, searchbypostcode *string, searchbydeliverytime *string) error {
	if *filepath == "" || *searchbyrecipetag == "" || *searchbypostcode == "" || *searchbydeliverytime == "" {
		return errors.New("missing parameters. must supply 'searchby' , 'postcode' and 'deliverytime' flags.")
	}
	//delivery time validations
	deliverytimetokens := strings.Split(*searchbydeliverytime, " ")
	if len(deliverytimetokens) != 3 {
		return errors.New("invalid deliverytime")
	}

	fromdeliverytime := deliverytimetokens[0]
	seperator := deliverytimetokens[1]
	todeliverytime := deliverytimetokens[2]

	if seperator != "-" {
		return errors.New("invalid deliverytime")
	}
	if len(fromdeliverytime) < 3 || len(todeliverytime) < 3 {
		return errors.New("invalid deliverytime")
	}
	if fromdeliverytime[len(fromdeliverytime)-2:len(fromdeliverytime)] != "AM" && fromdeliverytime[len(fromdeliverytime)-2:len(fromdeliverytime)] != "PM" {
		return errors.New("invalid deliverytime")
	}
	if todeliverytime[len(todeliverytime)-2:len(todeliverytime)] != "AM" && todeliverytime[len(todeliverytime)-2:len(todeliverytime)] != "PM" {
		return errors.New("invalid deliverytime")
	}
	fromdeliverytimeHour, _ := strconv.Atoi(fromdeliverytime[0 : len(fromdeliverytime)-2])
	todeliverytimeHour, _ := strconv.Atoi(fromdeliverytime[0 : len(fromdeliverytime)-2])
	if (fromdeliverytimeHour < 1 && fromdeliverytimeHour > 12) && (todeliverytimeHour < 1 && todeliverytimeHour > 12) {
		return errors.New("invalid deliverytime")
	}

	filePath = filepath
	searchByRecipeTag = searchbyrecipetag
	searchByPostCode = searchbypostcode
	searchByDeliveryTime = searchbydeliverytime
	return nil
}

//Parse ... parse the recipe json dump to create hashmap
func (p *parser) Parse() error {
	//open the file
	file, err := os.Open(*filePath)
	if err != nil {
		return err
	}

	//stream the file content and parse
	streamDecoder := jstream.NewDecoder(file, 1)
	for eachObj := range streamDecoder.Stream() {
		recipeObject := eachObj.Value.(map[string]interface{})
		var recipe, postcode, delivery string
		for k, v := range recipeObject {
			switch k {
			case "recipe":
				recipe = v.(string)
			case "postcode":
				postcode = v.(string)
			case "delivery":
				delivery = v.(string)
			}
		}

		//track counts in hashmap
		recipeMap[recipe] += 1
		postcodeMap[postcode] += 1
		if postcode == *searchByPostCode {
			canBeDelivered := p.CanBeDelivered(delivery, *searchByDeliveryTime)
			if canBeDelivered {
				feasibleDeliveryCount++
			}
		}

	}
	return nil
}

//GenerateResult ... Generate the output json and print it
func (p *parser) GenerateResult() error {
	searchTerms := strings.Split(*searchByRecipeTag, ",")
	countperrecipe := []CountPerRecipe{}

	//struct for desired JSON output format
	output := Output{
		UniqueRecipeCount: len(recipeMap),
		CountPerRecipe:    countperrecipe,
		BusiestPostcode: BusiestPostCode{
			Postcode:      "",
			DeliveryCount: 0,
		},
		MatchByName: []string{},
		CountPerPostCodeAndTime: CountPerPostCodeAndTime{
			Postcode: *searchByPostCode,
			From:     strings.Split(*searchByDeliveryTime, " ")[0],
			To:       strings.Split(*searchByDeliveryTime, " ")[2],
		},
	}

	//calculate unique recipes and their counts
	for k, v := range recipeMap {
		for _, searchTerm := range searchTerms {
			if strings.Contains(k, searchTerm) {
				output.MatchByName = append(output.MatchByName, k)
			}
		}
		output.CountPerRecipe = append(output.CountPerRecipe, CountPerRecipe{
			Recipe: k,
			Count:  v,
		})
	}

	//calculate busiest postcode
	for k, v := range postcodeMap {
		if v > output.BusiestPostcode.DeliveryCount {
			output.BusiestPostcode.Postcode = k
			output.BusiestPostcode.DeliveryCount = v
		}
	}

	//assign feasible delivery count for searched postcode
	output.CountPerPostCodeAndTime.DeliveryCount = feasibleDeliveryCount

	//print out the desired output into stderr
	json, err := json.MarshalIndent(output, "", "	")
	if err != nil {
		return err
	}
	fmt.Println(string(json))
	return nil
}

//CanBeDelivered ... Utility function to check feasibility of recipe delivery within specific time window
func (p *parser) CanBeDelivered(timerange string, searchtime string) bool {
	timeRangeFrom := strings.Split(timerange, " ")[1]
	timeRangeTo := strings.Split(timerange, " ")[3]
	searchTimeFrom := strings.Split(searchtime, " ")[0]
	searchTimeTo := strings.Split(searchtime, " ")[2]

	//12AM check -- convert to 0AM for easy conversion
	if timeRangeFrom == "12AM" {
		timeRangeFrom = "0AM"
	}
	if timeRangeTo == "12AM" {
		timeRangeTo = "0AM"
	}
	if searchTimeFrom == "12AM" {
		searchTimeFrom = "0AM"
	}
	if searchTimeTo == "12AM" {
		searchTimeTo = "0AM"
	}

	//int conversion for time
	timeRangeFromHour, _ := strconv.Atoi(timeRangeFrom[0 : len(timeRangeFrom)-2])
	timeRangeToHour, _ := strconv.Atoi(timeRangeTo[0 : len(timeRangeTo)-2])
	searchTimeFromHour, _ := strconv.Atoi(searchTimeFrom[0 : len(searchTimeFrom)-2])
	searchTimeToHour, _ := strconv.Atoi(searchTimeTo[0 : len(searchTimeTo)-2])

	//calculate zulu time
	if timeRangeFrom[len(timeRangeFrom)-2:len(timeRangeFrom)] == "PM" {
		timeRangeFromHour += 12
	}
	if timeRangeTo[len(timeRangeTo)-2:len(timeRangeTo)] == "PM" {
		timeRangeToHour += 12
	}
	if searchTimeFrom[len(searchTimeFrom)-2:len(searchTimeFrom)] == "PM" {
		searchTimeFromHour += 12
	}
	if searchTimeTo[len(searchTimeTo)-2:len(searchTimeTo)] == "PM" {
		searchTimeToHour += 12
	}
	if searchTimeFromHour >= timeRangeFromHour && searchTimeToHour <= timeRangeToHour {
		return true
	}

	return false
}
