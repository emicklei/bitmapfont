package bitmapfont

import "testing"

func TestRead(t *testing.T) {
	r := Reader{}
	err := r.Read("example/test_ubuntu.fnt")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", r.font)
}
