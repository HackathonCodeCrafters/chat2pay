package service

// stringPtr converts a string to *string, returns nil if string is empty
func stringPtr(s string) *string {
	if s == "" {
		return nil
	}
	return &s
}

// uint64Ptr converts a uint64 to *uint64, returns nil if value is 0
func uint64Ptr(v uint64) *uint64 {
	if v == 0 {
		return nil
	}
	return &v
}
