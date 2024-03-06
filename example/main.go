package main

import (
	"errors"
	"image"
	"image/draw"
	"io"
	"os"
	"runtime"

	_ "image/png"

	"github.com/emicklei/bitmapfont"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	Title  = "Bitmapfont demo"
	Width  = 800
	Height = 400
)

var font *bitmapfont.Font
var fontTexture uint32

func main() {
	if err := glfw.Init(); err != nil {
		panic("Can't init glfw!" + err.Error())
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	window, err := glfw.CreateWindow(Width, Height, Title, nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.MakeContextCurrent()
	glfw.SwapInterval(1)
	gl.Init()

	initScene()
	initFontAndText()
	defer gl.DeleteTextures(1, &fontTexture)

	runtime.LockOSThread()
	for !window.ShouldClose() {
		drawScene()
		window.SwapBuffers()
		glfw.PollEvents()
	}

}

var txt bitmapfont.Text

func initFontAndText() {
	// read font
	r := bitmapfont.NewReader()
	thefont, err := r.Read("test_ubuntu.fnt")
	if err != nil {
		panic(err)
	}
	font = thefont

	// read font texture
	f, err := os.Open("test_ubuntu.png")
	if err != nil {
		panic(err)
	}
	defer f.Close()
	fontTexture, err = newTexture(f)
	if err != nil {
		panic(err)
	}

	// create text
	var x, y, w, h float32 = 10, 10, 300, 100
	var multitext = `Ubanita
	together, we play`
	txt = bitmapfont.NewText(multitext, x, y, w, h, font)
}

func initScene() {
	gl.Disable(gl.DEPTH_TEST)
	gl.Viewport(0, 0, Width, Height)
	gl.MatrixMode(gl.PROJECTION)
	gl.LoadIdentity()
	gl.Ortho(0, Width, Height, 0, 0, 1)
	gl.MatrixMode(gl.MODELVIEW)
	gl.LoadIdentity()
	// Displacement trick for exact pixelization
	gl.Translatef(0.375, 0.375, 0)
	gl.Enable(gl.TEXTURE_2D)
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

func renderText() {

	gl.BindTexture(gl.TEXTURE_2D, fontTexture)
	gl.Enable(gl.TEXTURE_2D)

	txt.Render(func(vertices []bitmapfont.TextureVertex) {
		gl.Begin(gl.QUADS)
		for _, each := range vertices {
			gl.TexCoord2f(each.S, each.T)
			gl.Vertex2f(each.X, each.Y)
		}
		gl.End()
	})
	gl.Disable(gl.TEXTURE_2D)

	// render bounding box
	var x, y, w, h float32 = txt.X, txt.Y, txt.Width, txt.Height
	gl.Begin(gl.LINE_LOOP)
	gl.Vertex2f(x, y)
	gl.Vertex2f(x+w, y)
	gl.Vertex2f(x+w, y+h)
	gl.Vertex2f(x, y+h)
	gl.End()

}

func drawScene() {
	gl.ClearColor(0.0, 0.0, 0.0, 0)
	gl.Clear(gl.COLOR_BUFFER_BIT)
	renderText()
}
