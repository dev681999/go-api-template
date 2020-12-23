package util

// BoolPtr returns a bool ptr
func BoolPtr(b bool) *bool {
	newB := b
	return &newB
}
