package util

import (
	"strings"
	"slices"
)

func StandardizeBindingsArray(arr []string) []string {
	new := []string{}
	for _, s := range arr {
		new = append(new, StandardizeBinding(s))
	}
	return new
}

func StandardizeBinding(binding string) string {
	arr := strings.Split(binding, "+")

	for i, s := range arr {
		lower := strings.ToLower(s)
		arr[i] = lower
		
		if strings.HasPrefix(lower, "rune[") && strings.HasSuffix(lower, "]") {
			arr[i] = lower[5 : len(s)-1]
		}
	}

	slices.Sort(arr)
	return strings.ToLower(strings.Join(arr, "+"))
}

func NonNilLen[T any](slice []T) int {
	count := 0
	for _, item := range slice {
		if any(item) != nil {
			count++
		}
	}
	return count
}