package bitmapfont

import (
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
)

type Text struct {
	multiline string
	X         float32 // center x-coordinate
	Y         float32 // center y-coordinate
	width     float32 // text must fit into this width
	height    float32 // text must fit into this height
	font      *OpenGLFont
	vertices  [][]TextureVertex
}

// NewText return a new Text value for rendering a (multiline) string using a Font inside a 2d box.
func NewText(text string, x, y, width, height float32, font *OpenGLFont) Text {
	t := Text{multiline: text, X: x, Y: y, width: width, height: height, font: font}
	t.vertices = t.computeVertices()
	return t
}

func (t Text) Width() float32 {
	return t.width
}

func (t Text) Height() float32 {
	return t.height
}

func (t Text) computeVertices() (all [][]TextureVertex) {
	left, top := t.X, t.Y
	sw, sh := t.font.Scales()
	uw, uh := t.unscaledDimension()
	sx := t.width / uw
	sy := t.height / uh

	// split multiline text
	for _, each := range strings.Split(t.multiline, "\n") {
		var lastId uint8 = 0
		// each char
		for i := 0; i < len(each); i++ {
			char := t.font.CharAt(each[i])
			if lastId != 0 {
				// lookup space in between chars
				left += t.font.AmountBetween(lastId, char.Id) * sx
				lastId = char.Id
			}
			charTop := top + char.Yoffset*sy
			charBottom := charTop + char.Height*sy
			// all quad points
			vertices := []TextureVertex{
				{char.X / sw, char.Y / sh, left, charTop},
				{(char.X + char.Width) / sw, char.Y / sh, left + char.Width*sx, charTop},
				{(char.X + char.Width) / sw, (char.Y + char.Height) / sh, left + char.Width*sx, charBottom},
				{char.X / sw, (char.Y + char.Height) / sh, left, charBottom},
			}
			all = append(all, vertices)
			left += char.Xadvance * sx
		}
		top += t.font.LineHeight() * sy
		left = t.X
	}
	return
}

func (t Text) unscaledDimension() (width float32, height float32) {
	for _, each := range strings.Split(t.multiline, "\n") {
		lineWidth := float32(0)
		var lastId uint8 = 0
		for i := 0; i < len(each); i++ {
			char := t.font.CharAt(each[i])
			if lastId != 0 {
				// lookup space in between chars
				lineWidth += t.font.AmountBetween(lastId, char.Id)
				lastId = char.Id
			}
			lineWidth += char.Xadvance
		}
		if lineWidth > width {
			width = lineWidth
		}
		height += t.font.LineHeight()
	}
	return
}

// TextureVertex captures one corner of a character to render.
type TextureVertex struct {
	S, T float32 // texture coordinates
	X, Y float32 // position coordinates
}

// Render calls a function with 4 TextureVertex values per character.
// http://www.glprogramming.com/red/chapter09.html
// http://www.angelcode.com/products/bmfont/doc/render_text.html
func (t Text) Render() {
	gl.BindTexture(gl.TEXTURE_2D, t.font.texture)
	gl.Enable(gl.TEXTURE_2D)
	for _, each := range t.vertices {
		gl.Begin(gl.QUADS)
		for _, other := range each {
			gl.TexCoord2f(other.S, other.T)
			gl.Vertex2f(other.X, other.Y)
		}
		gl.End()
	}
	gl.Disable(gl.TEXTURE_2D)
}
