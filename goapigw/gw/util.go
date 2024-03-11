package main

import (
	"encoding/json"
	"fmt"
)

func Map2Struct[T any, K comparable, V any](t *T, m map[K]V) (*T, error) {
	marshal, err := json.Marshal(m)
	if err != nil {
		return nil, fmt.Errorf("marshal failed: %v", err.Error())
	}

	err = json.Unmarshal(marshal, t)
	if err != nil {
		return nil, fmt.Errorf("unmarshal failed: %v", err.Error())
	}

	return t, nil
}
