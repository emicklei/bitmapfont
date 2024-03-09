package bitmapfont

import "github.com/go-gl/gl/v2.1/gl"

type OpenGLFont struct {
	*Font
	texture uint32
}

// NewOpenGLFont createa a new font by reading the FNT and texture image file (PNG)
// This creates an OpenGL resource ; do not forget to call Delete() when no longer used.
func NewOpenGLFont(fntFilename, imageFilename string) (*OpenGLFont, error) {
	r := NewFontReader()
	font, err := r.ReadFile(fntFilename)
	if err != nil {
		return nil, err
	}
	b := NewTextureBuilder()
	txt, err := b.BuildFromFile(imageFilename)
	if err != nil {
		return nil, err
	}
	return &OpenGLFont{
		Font:    font,
		texture: txt,
	}, nil
}

func (f *OpenGLFont) Delete() {
	gl.DeleteTextures(1, &f.texture)
}
