package importer

import "testing"

func TestParseTSV(t *testing.T) {
	rows := ParseTSV("id\tname\taddress\tphone\temail\ttags\na\tA\tLane\t1\ta@x\twood|calm")
	if len(rows) != 1 || rows[0].Name != "A" {
		t.Fatalf("rows=%v", rows)
	}
}
