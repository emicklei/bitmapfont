package bitmapfont

import (
	"bytes"
	"fmt"
	"io"
)

type Font struct {
	info     Info
	Common   Common
	page     Page
	chars    map[uint8]Char
	kernings []Kerning
}

func NewFont() *Font {
	return &Font{chars: map[uint8]Char{}}
}

func (f *Font) addChar(c Char) {
	f.chars[c.Id] = c
}

func (f *Font) addKerning(k Kerning) {
	f.kernings = append(f.kernings, k)
}

func (f *Font) String() string {
	buf := new(bytes.Buffer)
	io.WriteString(buf, "Font(")
	fmt.Fprintf(buf, "face=%s, size=%d, bold=%t, #chars=%d, #kernings=%d", f.info.Face, f.info.Size, f.info.Bld, len(f.chars), len(f.kernings))
	io.WriteString(buf, ")")
	return buf.String()
}

func (f *Font) CharAt(c uint8) Char {
	ch, ok := f.chars[c]
	if !ok {
		return f.chars[" "[0]]
	}
	return ch
}

func (f *Font) Scales() (float32, float32) {
	return f.Common.ScaleW, f.Common.ScaleH
}

func (f *Font) LineHeight() float32 {
	return f.Common.LineHeight
}

func (f *Font) Base() float32 {
	return float32(f.Common.Base)
}

func (f *Font) AmountBetween(left, right uint8) float32 {
	// linear search for now
	for _, each := range f.kernings {
		if each.First == left && each.Second == right {
			return each.Amount
		}
	}
	return 0
}
