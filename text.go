package bitmapfont

import (
	"strings"
)

type Text struct {
	Text   string
	X      float32 // center x-coordinate
	Y      float32 // center y-coordinate
	Width  float32 // text must fit into this width
	Height float32 // text must fit into this height
	Font   *Font
}

// NewText return a new Text value for rendering a (multiline) string using a Font inside a 2d box.
func NewText(text string, x, y, width, height float32, font *Font) Text {
	return Text{Text: text, X: x, Y: y, Width: width, Height: height, Font: font}
}

// Render calls a function with 4 TextureVertex values per character.
// http://www.glprogramming.com/red/chapter09.html
// http://www.angelcode.com/products/bmfont/doc/render_text.html
func (t Text) Render(callback func([]TextureVertex)) {
	// compute top+left
	left, top := t.X-(t.Width/2), t.Y-(t.Height/2)
	sw, sh := t.Font.Scales()
	// split multiline text
	for _, each := range strings.Split(t.Text, "\n") {
		var lastId uint8 = 0
		// each char
		for i := 0; i < len(each); i++ {
			char := t.Font.CharAt(each[i])
			if lastId != 0 {
				// lookup space in between chars
				left += t.Font.AmountBetween(lastId, char.Id)
				lastId = char.Id
			}
			charTop := top + char.Yoffset
			charBottom := charTop + char.Height
			// all quad points
			vertices := []TextureVertex{
				{char.X / sw, char.Y / sh, left, charTop},
				{(char.X + char.Width) / sw, char.Y / sh, left + char.Width, charTop},
				{(char.X + char.Width) / sw, (char.Y + char.Height) / sh, left + char.Width, charBottom},
				{char.X / sw, (char.Y + char.Height) / sh, left, charBottom},
			}
			callback(vertices)
			left += char.Xadvance
		}
		top += t.Font.LineHeight()
		left = t.X - (t.Width / 2)
	}
}

// TextureVertex captures one corner of a character to render.
type TextureVertex struct {
	S, T float32 // texture coordinates
	X, Y float32 // position coordinates
}
