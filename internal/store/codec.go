package store

import "encoding/json"

func unmarshal(data []byte, target any) error { return json.Unmarshal(data, target) }
