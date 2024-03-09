package main

import (
	"runtime"

	_ "image/png"

	"github.com/emicklei/bitmapfont"

	"github.com/go-gl/gl/v2.1/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const (
	Width  = 800
	Height = 400
)

var openglfont *bitmapfont.OpenGLFont

func main() {
	if err := glfw.Init(); err != nil {
		panic("Can't init glfw!" + err.Error())
	}
	defer glfw.Terminate()

	glfw.WindowHint(glfw.Resizable, glfw.True)
	glfw.WindowHint(glfw.ContextVersionMajor, 2)
	glfw.WindowHint(glfw.ContextVersionMinor, 1)

	window, err := glfw.CreateWindow(Width, Height, "Bitmapfont Demo", nil, nil)
	if err != nil {
		panic(err)
	}
	defer window.Destroy()

	window.MakeContextCurrent()
	glfw.SwapInterval(1)
	gl.Init()

	initScene()
	initFontAndText()
	defer openglfont.Delete()

	runtime.LockOSThread()
	for !window.ShouldClose() {
		drawScene()
		window.SwapBuffers()
		glfw.PollEvents()
	}

}

var txt bitmapfont.Text
var footer bitmapfont.Text

func initFontAndText() {
	// make opengl font
	f, err := bitmapfont.NewOpenGLFont("test_ubuntu.fnt", "test_ubuntu.png")
	if err != nil {
		panic(err)
	}
	openglfont = f
	// create text
	var x, y, w, h float32 = 10, 10, 300, 100
	var multitext = `Bitmapfont
easy OpenGL font rendering
for Go`
	txt = bitmapfont.NewText(multitext, x, y, w, h, openglfont)

	footer = bitmapfont.NewText("fixed height and unspecified (=0) width", x, 200, 0, 30, openglfont)
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

func renderText() {
	txt.Render()
	footer.Render()

	// render bounding box
	var x, y, w, h float32 = txt.X, txt.Y, txt.Width(), txt.Height()
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
