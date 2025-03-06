package utils

import (
	"os"

	"golang.org/x/exp/rand"
)

func Getenv(key string, fallback string) string {
	value, exists := os.LookupEnv(key)
	if exists {
		return value
	} else {
		return fallback
	}
}

func RemoveSliceElement[T comparable](slice *([]T), target T) bool {
	for i, el := range *slice {
		if el == target {
			*slice = append((*slice)[:i], (*slice)[i+1:]...)
			return true
		}
	}
	return false
}

func ShuffleSlice[T any](slice *([]T)) {
	length := len(*slice)
	for i := 0; i < length; i++ {
		j := rand.Intn(i + 1)
		(*slice)[i], (*slice)[j] = (*slice)[j], (*slice)[i]
	}
}

func Mod(a, b int) int {
	return (a%b + b) % b
}
