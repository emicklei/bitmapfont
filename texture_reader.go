package bitmapfont

import (
	"errors"
	"image"
	"image/draw"
	"io"
	"os"

	"github.com/go-gl/gl/v2.1/gl"
)

type TextureReader struct{}

func NewTextureReader() TextureReader {
	return TextureReader{}
}

func (r TextureReader) ReadFile(filename string) (uint32, error) {
	// read font texture
	f, err := os.Open(filename)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return newTexture(f)
}

// newTexture creates OpenGL texture from image file

// from https://github.com/go-gl/examples/blob/master/glfw31-gl21-cube/cube.go
func newTexture(imgFile io.Reader) (uint32, error) {
	img, _, err := image.Decode(imgFile)
	if err != nil {
		return 0, err
	}

	rgba := image.NewRGBA(img.Bounds())
	if rgba.Stride != rgba.Rect.Size().X*4 {
		return 0, errors.New("unsupported stride")
	}
	draw.Draw(rgba, rgba.Bounds(), img, image.Point{0, 0}, draw.Src)

	var texture uint32
	gl.Enable(gl.TEXTURE_2D)
	gl.GenTextures(1, &texture)
	gl.BindTexture(gl.TEXTURE_2D, texture)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MIN_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_MAG_FILTER, gl.LINEAR)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_S, gl.CLAMP_TO_EDGE)
	gl.TexParameteri(gl.TEXTURE_2D, gl.TEXTURE_WRAP_T, gl.CLAMP_TO_EDGE)
	gl.TexImage2D(
		gl.TEXTURE_2D,
		0,
		gl.RGBA,
		int32(rgba.Rect.Size().X),
		int32(rgba.Rect.Size().Y),
		0,
		gl.RGBA,
		gl.UNSIGNED_BYTE,
		gl.Ptr(rgba.Pix))

	return texture, nil
}
