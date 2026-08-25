package flow031

import "encoding/json"

func marshalDocument(value any) ([]byte, error) { return json.MarshalIndent(value, "", "  ") }
