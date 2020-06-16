package internal

import "time"

// GetTimePointer returns a pointer to a time.Time instance for test cases
func GetTimePointer(t time.Time) *time.Time {
	return &t
}
