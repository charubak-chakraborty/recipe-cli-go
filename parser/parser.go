package parser

//Parser ... interface definition for parser
type Parser interface {
	Parse() error
	GenerateResult() error
	ValidateInput(file *string, searchByRecipeTag *string, searchByPostCode *string, searchByDeliveryTime *string) error
	CanBeDelivered(timerange string, searchtime string) bool
}
