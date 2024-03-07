package bitmapfont

import "testing"

func TestRead(t *testing.T) {
	r := NewFontReader()
	f, err := r.ReadFile("example/test_ubuntu.fnt")
	if err != nil {
		t.Error(err)
	}
	t.Logf("%v", f)
}
