package lodbc

import (
	"database/sql"
	"github.com/LukeMauldin/lodbc/odbc"
	"reflect"
)

// Utility function to quickly return rows from the database
func FetchRows(db *sql.DB, query string, tmplt interface{}) ([]interface{}, error) {
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	result := make([]interface{}, 0)
	for rows.Next() {
		s := reflect.ValueOf(tmplt).Elem()
		row := make([]interface{}, s.NumField())
		for i := 0; i < s.NumField(); i++ {
			row[i] = s.Field(i).Addr().Interface()
		}
		err := rows.Scan(row...)
		if err != nil {
			return nil, err
		}
		result = append(result, reflect.Indirect(reflect.ValueOf(tmplt)).Interface())
	}
	return result, nil
}

//Converts SQL_NUMERIC_STRUCT to float
func numericToFloat(inputValue odbc.SQL_NUMERIC_STRUCT) float64 {
	//Convert numeric data to float
	theData := make([]byte, 0, 16)
	for _, v := range inputValue.Val {
		theData = append(theData, byte(v))
	}
	outputVal := byteToHextOval(theData)

	//Take into account the scale - converts to a decimal if necessary
	divisor := float64(1)
	if inputValue.Scale > 0 {
		for i := odbc.SQLCHAR(0); i < inputValue.Scale; i++ {
			divisor = divisor * 10
		}
	}

	finalVal := float64(outputVal) / divisor

	//Take into account the sign - if it is 0, convert to a negative
	if inputValue.Sign == 0 {
		finalVal = finalVal * -1
	}

	return finalVal
}

//Helper function for numericToFloat
func byteToHextOval(inputVal []byte) int64 {
	var value int64
	last := int64(1)
	var current int64
	var a, b int64

	for i := 0; i <= 15; i++ {
		current = int64(inputVal[i])
		a = current % 16
		b = current / 16

		value += last * a
		last = last * 16
		value += last * b
		last = last * 16
	}
	return value
}

//Checks the type v for nil
func isNil(v interface{}) bool {
	if v == nil {
		return true
	}
	val := reflect.ValueOf(v)
	switch val.Kind() {
	case reflect.Chan, reflect.Func, reflect.Interface, reflect.Map,
		reflect.Ptr, reflect.Slice, reflect.UnsafePointer:
		return val.IsNil()
	}
	return false
}
