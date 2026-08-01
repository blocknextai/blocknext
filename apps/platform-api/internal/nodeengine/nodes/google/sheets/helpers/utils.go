package helpers

import (
	"strings"

	"github.com/blocknextai/go-packages/apperror"
	"github.com/blocknextai/go-packages/cast"
)

var (
	errInvalidColumnLetter = apperror.Internal("invalid column letter")
	errColumnNotFound      = apperror.Internal("column not found")
)

func ColumnLetterToIndex(col string) (int, error) {
	col = strings.ToUpper(col)
	result := 0
	for _, char := range col {
		if char < 'A' || char > 'Z' {
			return 0, errInvalidColumnLetter
		}
		result = result*26 + int(char-'A'+1)
	}
	return result, nil
}

func IndexToColumnLetter(index int) string {
	result := ""
	for index >= 0 {
		result = string(rune('A'+index%26)) + result
		index = index/26 - 1
	}
	return result
}

func FindColumnIndex(columnIdentifier string, headerRow []string) (int, error) {
	for i, header := range headerRow {
		if strings.EqualFold(strings.TrimSpace(header), strings.TrimSpace(columnIdentifier)) {
			return i, nil
		}
	}

	if index, err := ColumnLetterToIndex(columnIdentifier); err == nil {
		if index < len(headerRow) {
			return index, nil
		}
	}

	return -1, errColumnNotFound
}

func FormatSpreadsheetData(values [][]string, headerRow int, filterColumn string, filterValue string, limit int) []map[string]any {
	if len(values) == 0 {
		return []map[string]any{}
	}

	hasHeader := headerRow > 0 && headerRow <= len(values)
	headerRowIndex := headerRow - 1

	var headers []string
	var dataRows [][]string

	if hasHeader {
		headers = values[headerRowIndex]

		for i, row := range values {
			if i != headerRowIndex {
				dataRows = append(dataRows, row)
			}
		}
	} else {
		if len(values) > 0 {
			for i := 0; i < len(values[0]); i++ {
				headers = append(headers, "column_"+cast.ToString(i+1))
			}
		}
		dataRows = values
	}

	formattedData := []map[string]any{}

	for rowIndex, row := range dataRows {
		obj := make(map[string]any)

		actualRowNumber := rowIndex
		if hasHeader {
			actualRowNumber = rowIndex + headerRow + 1
		}

		for i, header := range headers {
			obj["row_number"] = actualRowNumber
			if i < len(row) {
				obj[header] = row[i]
			} else {
				obj[header] = ""
			}
		}

		if filterColumn != "" && filterValue != "" {
			if value, exists := obj[filterColumn]; exists {
				valueStr, _ := value.(string)

				if filterValue == "false" && valueStr == "" {
				} else if valueStr != filterValue {
					continue
				}
			} else {
				if filterValue == "false" || filterValue == "" {
				} else {
					continue
				}
			}
		}

		formattedData = append(formattedData, obj)

		if limit > 0 && len(formattedData) >= limit {
			break
		}
	}

	if formattedData == nil {
		return []map[string]any{}
	}

	return formattedData
}
