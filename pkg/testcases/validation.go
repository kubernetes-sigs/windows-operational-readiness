package testcases

import (
	"fmt"
	"github.com/go-playground/validator/v10"
)

// validateYAML validate the input YAML and returns the error
func (o *OpTestConfig) validateYAML() error {
	validate := validator.New()
	validate.RegisterStructValidation(OpTestConfigValidation, OpTestConfig{})
	if err := validate.Struct(o); err != nil {
		return err
	}
	return nil
}

// OpTestConfigValidation set the required fields and is used by the validator function
func OpTestConfigValidation(sl validator.StructLevel) {
	opTestConfig := sl.Current().Interface().(OpTestConfig)
	for _, opTestCase := range opTestConfig.OpTestCases {
		if opTestCase.Category == "" {
			fmt.Println("Category Required")
			sl.ReportError(opTestCase.Category, "category", "Category", "categoryRequired", "")
		}
		if opTestCase.Description == "" {
			fmt.Println("Description Required")
			sl.ReportError(opTestCase.Description, "description", "Description", "descriptionRequired", "")
		}
		if len(opTestCase.Focus) == 0 || (len(opTestCase.Focus) == 1 && len(opTestCase.Focus[0]) == 0) {
			fmt.Println("Focus Required")
			sl.ReportError(opTestCase.Focus, "focus", "Focus", "focusRequired", "")
		}
	}
}
