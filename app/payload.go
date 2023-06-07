package app

import (
	"encoding/json"
	"fmt"
	"marble/marble-backend/models"
	"marble/marble-backend/pure_utils"
	"strings"
	"time"

	"github.com/go-playground/validator"
	dynamicstruct "github.com/ompluscator/dynamic-struct"
)

func buildDynamicStruct(fields map[models.FieldName]models.Field) dynamicstruct.DynamicStruct {
	custom_type := dynamicstruct.NewStruct()

	var stringPointerType *string
	var intPointerType *int
	var floatPointerType *float64
	var boolPointerType *bool
	var timePointerType *time.Time

	// those fields are mandatory for all tables
	custom_type.AddField("Object_id", stringPointerType, `validate:"required"`)
	custom_type.AddField("Updated_at", timePointerType, `validate:"required"`)

	for fieldName, field := range fields {
		name := string(fieldName)
		switch strings.ToLower(name) {
		case "object_id", "updated_at":
			// already added above, with a different validation tag
			break
		default:
			switch field.DataType {
			case models.Bool:
				custom_type.AddField(pure_utils.Capitalize(name), boolPointerType, "")
			case models.Int:
				custom_type.AddField(pure_utils.Capitalize(name), intPointerType, "")
			case models.Float:
				custom_type.AddField(pure_utils.Capitalize(name), floatPointerType, "")
			case models.String:
				custom_type.AddField(pure_utils.Capitalize(name), stringPointerType, "")
			case models.Timestamp:
				custom_type.AddField(pure_utils.Capitalize(name), timePointerType, "")
			}
		}
	}
	return custom_type.Build()
}

func validateParsedJson(instance interface{}) error {
	validate := validator.New()
	err := validate.Struct(instance)
	if err != nil {

		// This error should happen in the dynamic struct is badly formatted, or if the tags
		// contain bad values. If this returns an error, it's a 500 error.
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}

		// Otherwise it's a 400 error and we can access the reasons from here
		count := 0
		for _, err := range err.(validator.ValidationErrors) {
			fmt.Printf("The input object is not valid: key %v, validation tag: '%v', receive value %v", err.Field(), err.Tag(), err.Param())
			count++
		}
		if count > 0 {
			return models.FormatValidationError
		}
	}
	return nil
}

func ParseToDataModelObject(table models.Table, jsonBody []byte) (models.Payload, error) {
	fields := table.Fields

	custom_type := buildDynamicStruct(fields)

	dynamicStructInstance := custom_type.New()
	dynamicStructReader := dynamicstruct.NewReader(dynamicStructInstance)

	// This is where errors can happen while parson the json. We could for instance have badly formatted
	// json, or timestamps.
	// We could also have more serious errors, like a non-capitalized field in the dynamic struct that
	// causes a panic. We should manage the errors accordingly.
	decoder := json.NewDecoder(strings.NewReader(string(jsonBody)))
	// Reject fields that are not present in the data model/the dynamic struct
	decoder.DisallowUnknownFields()

	if err := decoder.Decode(&dynamicStructInstance); err != nil {
		return models.Payload{}, fmt.Errorf("%w: %w", models.FormatValidationError, err)
	}

	// If the data has been successfully parsed, we can validate it
	// This is done using the validate tags on the dynamic struct
	// There are two possible cases of error
	err := validateParsedJson(dynamicStructInstance)
	if err != nil {
		return models.Payload{}, err
	}

	return models.Payload{Reader: dynamicStructReader, Table: table}, nil
}
