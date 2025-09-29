package utils

import "github.com/google/uuid"

// Helper functions for safe pointer dereferencing with sensible defaults

// DerefString safely dereferences a string pointer, returning empty string if nil
func DerefString(p *string) string {
	if p == nil {
		return ""
	}
	return *p
}

// DerefInt64 safely dereferences an int64 pointer, returning 0 if nil
func DerefInt64(p *int64) int64 {
	if p == nil {
		return 0
	}
	return *p
}

// DerefBool safely dereferences a bool pointer, returning false if nil
func DerefBool(p *bool) bool {
	if p == nil {
		return false
	}
	return *p
}

// DerefUUID safely dereferences a UUID pointer, returning uuid.Nil if nil
func DerefUUID(p *uuid.UUID) uuid.UUID {
	if p == nil {
		return uuid.Nil
	}
	return *p
}

// DerefInt safely dereferences an int pointer, returning 0 if nil
func DerefInt(p *int) int {
	if p == nil {
		return 0
	}
	return *p
}
