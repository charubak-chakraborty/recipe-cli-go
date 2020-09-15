package parser

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidateInput(t *testing.T) {

	var (
		filepath             string = "path/to/file"
		searchbyrecipetag    string = "Potato"
		searchbypostcode     string = "10210"
		searchbydeliverytime string = "1AM - 10PM"
	)

	parser := NewParser()
	err := parser.ValidateInput(&filepath, &searchbyrecipetag, &searchbypostcode, &searchbydeliverytime)
	assert.Nil(t, err)

}

func TestValidateInputError(t *testing.T) {

	var (
		filepath             string = ""
		searchbyrecipetag    string = ""
		searchbypostcode     string = ""
		searchbydeliverytime string = ""
	)

	parser := NewParser()
	err := parser.ValidateInput(&filepath, &searchbyrecipetag, &searchbypostcode, &searchbydeliverytime)
	assert.Error(t, err)

}

func TestCanBeDelivered(t *testing.T) {

	parser := NewParser()
	canbedeliveredTrue_1 := parser.CanBeDelivered("Thursday 6PM - 12AM", "11PM - 12AM")
	canbedeliveredTrue_2 := parser.CanBeDelivered("Thursday 6PM - 12AM", "11PM - 12AM")

	canbedeliveredFalse_1 := parser.CanBeDelivered("Thursday 6PM - 12AM", "11PM - 1PM")
	canbedeliveredFalse_2 := parser.CanBeDelivered("Thursday 6PM - 12AM", "12AM - 12AM")

	assert.True(t, canbedeliveredTrue_1)
	assert.True(t, canbedeliveredTrue_2)

	assert.False(t, canbedeliveredFalse_1)
	assert.False(t, canbedeliveredFalse_2)

}

func TestParseAndGenerate(t *testing.T) {

	var (
		filepath             string = "../input.json"
		searchbyrecipetag    string = "Potato"
		searchbypostcode     string = "10210"
		searchbydeliverytime string = "1AM - 11PM"
	)

	parser := NewParser()
	err := parser.ValidateInput(&filepath, &searchbyrecipetag, &searchbypostcode, &searchbydeliverytime)
	assert.Nil(t, err)
	err = parser.Parse()
	assert.Nil(t, err)
	err = parser.GenerateResult()
	assert.Nil(t, err)

}
