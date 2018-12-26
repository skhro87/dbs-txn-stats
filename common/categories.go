package common

import "fmt"

var Categories = []string{"FOOD", "DRINKS/PARTY", "TRANSPORT", "SALARY", "RENT", "UTILITIES", "CLOTHES", "ELECTRONICS", "OTHER"}

func CategoryTranslation(category int) (string, error) {
	category = category - 1
	if category < 0 || category+1 > len(Categories) {
		return "", fmt.Errorf("invalid category : %v", category+1)
	}

	return Categories[category], nil
}
