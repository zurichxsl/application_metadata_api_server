package server

import (
	"application_metadata_api_server/server/api"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"
	"testing"
)

func TestAppValidator_ValidatePut(t *testing.T) {
	validator := newAppValidator()

	testCases := []struct {
		name                    string
		app                     *api.App
		expectedValidationError bool
	}{
		{
			name: "missing required field values, expect validation fail",
			app: &api.App{
				Title: "11",
			},
			expectedValidationError: true,
		},
		{
			name: "valid input, validation pass",
			app: &api.App{
				Title:   "11",
				Version: "v",
				Maintainers: []api.Maintainer{
					{
						Name:  "aiden",
						Email: "aiden@gmail.com",
					},
				},
				Company: "c",
				Website: "http:sss",
				Source:  "ddd",
				License: "ddd",
			},
			expectedValidationError: false,
		},
		{
			name: "invalid email input, validation fail",
			app: &api.App{
				Title:   "11",
				Version: "v",
				Maintainers: []api.Maintainer{
					{
						Name:  "aiden",
						Email: "aidefffc.com",
					},
				},
				Company: "c",
				Website: "http:sss",
				Source:  "ddd",
				License: "ddd",
			},
			expectedValidationError: true,
		},
		{
			name: "nested field missing, validation fail",
			app: &api.App{
				Title:   "11",
				Version: "v",
				Maintainers: []api.Maintainer{
					{
						Email: "aide@fffc.com",
					},
				},
				Company: "c",
				Website: "http:sss",
				Source:  "ddd",
				License: "ddd",
			},
			expectedValidationError: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			yamlData, err := yaml.Marshal(test.app)
			assert.Nil(t, err)
			_, err = validator.ValidatePut(yamlData)
			if test.expectedValidationError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}

}

func TestAppValidator_ValidateSearch(t *testing.T) {
	validator := newAppValidator()

	testCases := []struct {
		name                    string
		app                     interface{}
		expectedValidationError bool
	}{
		{
			name: "valid search query, validation pass",
			app: api.App{
				Title: "11",
			},
			expectedValidationError: false,
		},
		{
			name:                    "invalid search query, validation fail",
			app:                     "abcd",
			expectedValidationError: true,
		},
	}
	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			yamlData, err := yaml.Marshal(test.app)
			assert.Nil(t, err)
			_, err = validator.ValidateSearch(yamlData)
			if test.expectedValidationError {
				assert.NotNil(t, err)
			} else {
				assert.Nil(t, err)
			}
		})
	}

}
