package utils

import "time"

// StringPtr returns a pointer to string
func StringPtr(s string) *string {
	return &s
}

// TimePtr returns a pointer to time.Time
func TimePtr(t time.Time) *time.Time {
	return &t
}

// Int64Ptr returns a pointer to int64
func Int64Ptr(i int64) *int64 {
	return &i
}

// Float64Ptr returns a pointer to float64
func Float64Ptr(f float64) *float64 {
	return &f
}

// BoolPtr returns a pointer to bool
func BoolPtr(b bool) *bool {
	return &b
}
