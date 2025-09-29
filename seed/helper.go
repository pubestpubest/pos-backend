package seed

func ptr[T any](v T) *T     { return &v }
func ptrBool(v bool) *bool  { return &v }
func ptrInt(v int) *int     { return &v }
func ptrI64(v int64) *int64 { return &v }
