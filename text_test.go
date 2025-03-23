package bitmapfont

import (
	_ "image/png"
	"testing"
)

func TestNewTextNoFont(t *testing.T) {
	defer func() {
		if err := recover(); err == nil {
			t.Error(err)
		}
	}()
	NewText("", 0, 0, 100, 100, nil)
}
