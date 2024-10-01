package utils

import "fmt"

func FindCommonElement(slice1, slice2 []uint) (uint, error) {
	elements := make(map[uint]struct{})

	for _, elem := range slice1 {
		elements[elem] = struct{}{}
	}

	var commonElem uint
	found := false

	for _, elem := range slice2 {
		if _, exists := elements[elem]; exists {
			if found {
				// More than one common element found
				return 0, fmt.Errorf("more than one common element found")
			}
			commonElem = elem
			found = true
		}
	}

	if !found {
		// No common element found
		return 0, fmt.Errorf("no common element found")
	}

	return commonElem, nil
}

func SliceDiff(slice1, slice2 []string) []string {
	elementMap := make(map[string]struct{})
	for _, v := range slice2 {
		elementMap[v] = struct{}{}
	}

	result := []string{}
	for _, v := range slice1 {
		if _, found := elementMap[v]; !found {
			result = append(result, v)
		}
	}

	return result
}

func RemoveValue(slice []string, valueToRemove string) []string {
	var result []string

	for _, v := range slice {
		if v != valueToRemove {
			result = append(result, v)
		}
	}

	return result
}
