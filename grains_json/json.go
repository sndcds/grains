package grains_json

import (
	"encoding/json"
	"fmt"
)

// IntArray is a reusable type for a JSON array of integers.
// It supports both formats:
//
//	[{"category": 4}, {"category": 5}]  -> []int{4,5}
//	[4,5]                                -> []int{4,5}
type IntArray []int

// UnmarshalJSON implements json.Unmarshaler
func (ia *IntArray) UnmarshalJSON(data []byte) error {
	// Try simple array of ints first
	var simple []int
	if err := json.Unmarshal(data, &simple); err == nil {
		*ia = simple
		return nil
	}

	// Fallback: array of objects with "category" key
	var objs []map[string]interface{}
	if err := json.Unmarshal(data, &objs); err != nil {
		return fmt.Errorf("IntArray: failed to unmarshal JSON: %w", err)
	}

	var out []int
	for _, obj := range objs {
		val, ok := obj["category"]
		if !ok {
			return fmt.Errorf("IntArray: key 'category' not found in object %v", obj)
		}
		switch v := val.(type) {
		case float64: // JSON numbers are float64
			out = append(out, int(v))
		default:
			return fmt.Errorf("IntArray: key 'category' is not a number: %v", val)
		}
	}

	*ia = out
	return nil
}

// MarshalJSON implements json.Marshaler
func (ia IntArray) MarshalJSON() ([]byte, error) {
	return json.Marshal([]int(ia))
}
