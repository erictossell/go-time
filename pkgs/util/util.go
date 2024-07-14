package util

import (
	"github.com/charmbracelet/huh"
	"time"
)

func Map[T any, U any](slice []T, f func(T) U) []U {
	result := make([]U, len(slice))
	for i, v := range slice {
		result[i] = f(v)
	}
	return result
}

func IndexOf(slice []string, value string) int {
	for i, item := range slice {
		if item == value {
			return i
		}
	}
	return -1
}

func CreateTagOptions(tags []string) []huh.Option[string] {
	options := make([]huh.Option[string], len(tags))
	for i, tag := range tags {
		options[i] = huh.NewOption(tag, tag)
	}
	return options
}

func TimePtrToStringPtr(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format("2006-01-02 15:04:05")
	return &s
}
