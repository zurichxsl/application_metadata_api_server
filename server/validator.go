package server

import (
	"application_metadata_api_server/server/api"
	"fmt"
	"reflect"
	"regexp"
	"sigs.k8s.io/yaml"
	"strings"
)

// Validator is the interface to validate App schema field who has "validate" tag
type Validator interface {
	// ValidatePut validates any field of App who has "validate" tag, be it a struct or single field, also automatically validates nested structs
	ValidatePut(req []byte) (api.App, ValidationError)
	// ValidateSearch validates if input is a App struct, but does not validate around "validate" tag
	ValidateSearch(req []byte) (api.App, ValidationError)
}

// appValidator is an implementation of Validator
type appValidator struct {
	validators map[string]func(name string, obj interface{}) (bool, error)
}

func newAppValidator() Validator {
	return &appValidator{
		validators: map[string]func(name string, obj interface{}) (bool, error){
			"email":    isEmailValid,
			"required": isRequired,
		},
	}
}

func (v *appValidator) ValidatePut(req []byte) (api.App, ValidationError) {
	app := &api.App{}
	if err := yaml.Unmarshal(req, app); err != nil {
		return *app, NewInvalidSpec(err)
	}
	value := reflect.ValueOf(app)
	validationErr := v.traverseField(value)
	return *app, validationErr
}

func (v *appValidator) ValidateSearch(req []byte) (api.App, ValidationError) {
	app := &api.App{}
	if err := yaml.Unmarshal(req, app); err != nil {
		return *app, NewInvalidSpec(err)
	}
	return *app, nil
}

// validateStruct check if any struct level validations, after all field validations already checked.
func (v *appValidator) validateStruct(cur reflect.Value) ValidationError {
	for i := 0; i < cur.NumField(); i++ {
		field := cur.Type().Field(i)
		valueField := cur.Field(i)
		vTags := getValidateTags(field)
		for _, vTag := range vTags {
			handlerFn, ok := v.validators[vTag]
			if ok {
				_, err := handlerFn(field.Name, valueField.Interface())
				if err != nil {
					return NewInvalidSpec(err)
				}
			}
		}
		if len(vTags) > 0 {
			err := v.traverseField(valueField)
			if err != nil {
				return NewInvalidSpec(err)
			}
		}
	}
	return nil
}

// traverseField validates any field, be it a struct or single field,
// also automatically validates nested structs
func (v *appValidator) traverseField(cur reflect.Value) ValidationError {
	k := cur.Kind()
	switch k {
	case reflect.Pointer:
		err := v.validateStruct(cur.Elem())
		if err != nil {
			return NewInvalidSpec(err)
		}
	case reflect.Slice:
		for j := 0; j < cur.Len(); j++ {
			err := v.traverseField(cur.Index(j))
			if err != nil {
				return NewInvalidSpec(err)
			}
		}
	case reflect.Struct:
		err := v.validateStruct(cur)
		if err != nil {
			return NewInvalidSpec(err)
		}
	}
	return nil
}

func isEmailValid(name string, obj interface{}) (bool, error) {
	e := obj.(string)
	emailRegex := regexp.MustCompile(`^[a-z0-9._%+\-]+@[a-z0-9.\-]+\.[a-z]{2,4}$`)
	v := emailRegex.MatchString(e)
	if !v {
		return v, fmt.Errorf("%s email is invalid", name)
	}
	return v, nil
}

func getStructTag(f reflect.StructField, tagName string) string {
	return f.Tag.Get(tagName)
}

func getValidateTags(v reflect.StructField) []string {
	validateTag := getStructTag(v, "validate")
	if len(validateTag) == 0 {
		return []string{}
	}
	return strings.Split(validateTag, ",")
}

func isRequired(name string, obj interface{}) (bool, error) {
	v := notEmpty(obj)
	if !v {
		return v, fmt.Errorf("%s is required", name)
	}
	return v, nil
}

func notEmpty(obj interface{}) bool {
	value := reflect.ValueOf(obj)
	kind := value.Type().Kind()
	switch kind {
	case reflect.String:
		return len(strings.TrimSpace(value.String())) > 0
	case reflect.Map, reflect.Slice:
		return value.Len() > 0
	case reflect.Ptr, reflect.Interface, reflect.Func:
		return !value.IsNil()
	default:
		return value.IsValid() && value.Interface() != reflect.Zero(value.Type()).Interface()
	}
}
