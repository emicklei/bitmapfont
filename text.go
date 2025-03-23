package bitmapfont

import (
	"strings"

	"github.com/go-gl/gl/v2.1/gl"
)

type Text struct {
	multiline   string
	X           float32 // left-top x-coordinate
	Y           float32 // left-top y-coordinate
	widthZ      float32 // text must fit into this width, or zero
	height      float32 // text must fit into this height
	font        *OpenGLFont
	vertices    [][]TextureVertex
	actualWidth float32
}

// Centered returns a new text centered on its X and Y coordinate.
func (t Text) Centered() Text {
	t.X -= t.actualWidth / 2
	t.Y -= t.height / 2
	t.vertices, t.actualWidth = t.computeVertices()
	return t
}

// NewText returns a new Text value for rendering a (multiline) string using a Font inside a 2D box.
// If the width is set to zero then the text will have a unscaled width computed from the text.
func NewText(text string, leftTopX, leftTopY, widthOrZero, height float32, font *OpenGLFont) Text {
	if font == nil {
		panic("font required")
	}
	t := Text{multiline: text, X: leftTopX, Y: leftTopY, widthZ: widthOrZero, height: height, font: font}
	t.vertices, t.actualWidth = t.computeVertices()
	return t
}

func (t Text) Width() float32 {
	return t.widthZ
}

func (t Text) ActualWidth() float32 {
	return t.actualWidth
}

func (t Text) Height() float32 {
	return t.height
}

func (t Text) computeVertices() (all [][]TextureVertex, actualWidth float32) {
	left, top := t.X, t.Y
	sw, sh := t.font.Scales()
	uw, uh := t.unscaledDimension()
	sx := t.widthZ / uw
	sy := t.height / uh
	if t.widthZ == 0 {
		sx = sy
	}
	actualWidth = uw * sx

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
		for i := range len(each) {
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
