package lodbc

import (
	"encoding/json"
	"fmt"
	"regexp"
)

// Defines the supported types of query options
type QueryOptionKey int

const (
	//Result set number to read when queries return multiple result sets
	ResultSetNum QueryOptionKey = iota
)

// Identifier to add in SQL query to indiciate start/end of options
const optionIdentifier = "@!!@"

var regexOptions = regexp.MustCompile("^@!!@.*@!!@")
var regexOptionIdentifier = regexp.MustCompile("@!!@")

// Structure defining query option key and value
type QueryOption struct {
	Key   QueryOptionKey
	Value interface{}
}

func NewQueryOption(key QueryOptionKey, value interface{}) QueryOption {
	return QueryOption{Key: key, Value: value}
}

/*
 * Adds a query option to a SQL statement
 */
func AddQueryOption(sqlQuery string, option QueryOption) (string, error) {
	return AddQueryOptions(sqlQuery, []QueryOption{option})
}

/*
 * Adds query options to a SQL statement
 */
func AddQueryOptions(sqlQuery string, options []QueryOption) (string, error) {
	//Return the SQL query if there are no options specified
	if len(options) == 0 {
		return sqlQuery, nil
	}

	//Use JSON encoding to encode a string with the options key/values
	encodedOptions, err := json.Marshal(options)
	if err != nil {
		return "", err
	}

	//Format a new string that begins/ends with optionIdentifier and contains encodedOptions
	return fmt.Sprintf("%v%v%v%v", optionIdentifier, string(encodedOptions), optionIdentifier, sqlQuery), nil
}

/*
 * Parses options from a SQL statement
 */
func parseQueryOptions(sqlQuery string) ([]QueryOption, error) {

	//Parse matchOptions
	matchOptions := regexOptions.FindString(sqlQuery)

	//If there are no options, return an empty slice
	if matchOptions == "" {
		return make([]QueryOption, 0), nil
	}

	//Remove the beginning and ending identifiers
	matchOptions = regexOptionIdentifier.ReplaceAllString(matchOptions, "")

	//JSON decode the matchOptions
	options := make([]QueryOption, 0)
	err := json.Unmarshal([]byte(matchOptions), &options)
	if err != nil {
		return nil, err
	}

	return options, nil
}

/*
 * Return the value for a specified option
 */
func getOptionValue(options []QueryOption, key QueryOptionKey) (interface{}, bool) {
	//Iterate through the options to try to find the key
	for _, option := range options {
		if option.Key == key {
			return option.Value, true
		}
	}

	//Option was not found
	return "", false
}

/*
 * Removes all options from the SQL statement
 */
func removeOptions(sqlQuery string) string {
	return regexOptions.ReplaceAllString(sqlQuery, "")
}
