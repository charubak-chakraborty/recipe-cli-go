package parser

//Parser ... Interface Definition for Parser
type Parser interface {
	Parse() error
	GenerateResult() error
	ValidateInput(file *string, searchByRecipeTag *string, searchByPostCode *string, searchByDeliveryTime *string) error
	CanBeDelivered(timerange string, searchtime string) bool
}
