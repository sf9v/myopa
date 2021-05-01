package myopa

import "strings"

// M is an alias to map of interfaces
type M map[string]interface{}

func processTerm(query string) []string {
	splitQ := strings.Split(query, ".")
	var result []string
	for _, term := range splitQ {
		result = append(result, removeOpenBrace(term))
	}

	if result == nil {
		return nil
	}

	if len(result) == 1 {
		return []string{result[0]}
	}

	indexName := result[1]
	fieldName := result[2]
	if len(result) > 2 {
		fieldName = strings.Join(result[2:], ".")
	}

	return []string{indexName, fieldName}
}

func removeOpenBrace(input string) string {
	return strings.Split(input, "[")[0]
}
