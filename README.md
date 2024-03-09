## bitmapfont reader and opengl texture builder

## example

![Bitmpafont demo](demo.png)

## format

http://www.angelcode.com/products/bmfont/doc/file_format.html
http://www.angelcode.com/products/bmfont/
http://www.angelcode.com/products/bmfont/doc/render_text.html


## usage

```
	f, _ := bitmapfont.NewOpenGLFont("test_ubuntu.fnt", "test_ubuntu.png")
    defer f.Delete()
	
	var multitext = `Bitmapfont
easy OpenGL font rendering
for Go`
	txt = bitmapfont.NewText(multitext, 10, 10, 300, 100, f)

    txt.Render()
```

## TODO
- multi-byte chars
- optimize amout lookup

(c) 2016, http://ernestmicklei.com. MIT License	