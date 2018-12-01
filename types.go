package bitmapfont

// http://www.angelcode.com/products/bmfont/doc/file_format.html

// This tag holds information on how the font was generated.
type Info struct {
	Face     string // This is the name of the true type font.
	Size     int    // The size of the true type font.
	Bld      bool   // The font is bold.
	Italic   bool   // The font is italic.
	Charset  string // The name of the OEM charset used (when not unicode).
	Unicode  bool   // Set to 1 if it is the unicode charset.
	StretchH int    // The font height stretch in percentage. 100% means no stretch.
	Smooth   bool   // Set to 1 if smoothing was turned on.
	AA       int    // The supersampling level used. 1 means no supersampling was used.
	Padding  []int  // The padding for each character (up, right, down, left).
	Spacing  []int  // The spacing for each character (horizontal, vertical).
	Outline  int    // The outline thickness for the characters.
}

func BuildInfo(kvs map[string]value) Info {
	i := Info{}
	for k, v := range kvs {
		switch k {
		case "face":
			i.Face = v.stringValue
		case "size":
			i.Size = v.intValue
		case "bld":
			i.Bld = v.intValue > 0
		case "italic":
			i.Italic = v.intValue > 0
		case "charset":
			i.Charset = v.stringValue
		case "unicode":
			i.Unicode = v.intValue == 1
		case "stretchH":
			i.StretchH = v.intValue
		case "smooth":
			i.Smooth = v.intValue == 1
		case "aa":
			i.AA = v.intValue
		case "padding":
			i.Padding = v.intArray
		case "spacing":
			i.Spacing = v.intArray
		case "outline":
			i.Outline = v.intValue
		}
	}
	return i
}

// This tag holds information common to all characters.
type Common struct {
	LineHeight float32 // This is the distance in pixels between each line of text.
	Base       int     // The number of pixels from the absolute top of the line to the base of the characters.
	ScaleW     float32 // The width of the texture, normally used to scale the x pos of the character image.
	ScaleH     float32 // The height of the texture, normally used to scale the y pos of the character image.
	Pages      int     // The number of texture pages included in the font.
	Packed     int     // Set to 1 if the monochrome characters have been packed into each of the texture channels. In this case alphaChnl describes what is stored in each channel.
	AlphaChnl  int     // Set to 0 if the channel holds the glyph data, 1 if it holds the outline, 2 if it holds the glyph and the outline, 3 if its set to zero, and 4 if its set to one.
	RedChnl    int     // Set to 0 if the channel holds the glyph data, 1 if it holds the outline, 2 if it holds the glyph and the outline, 3 if its set to zero, and 4 if its set to one.
	GreenChnl  int     // Set to 0 if the channel holds the glyph data, 1 if it holds the outline, 2 if it holds the glyph and the outline, 3 if its set to zero, and 4 if its set to one.
	BlueChnl   int     // Set to 0 if the channel holds the glyph data, 1 if it holds the outline, 2 if it holds the glyph and the outline, 3 if its set to zero, and 4 if its set to one.
}

func BuildCommon(kvs map[string]value) Common {
	c := Common{}
	for k, v := range kvs {
		switch k {
		case "lineHeight":
			c.LineHeight = v.f32()
		case "base":
			c.Base = v.intValue
		case "scaleW":
			c.ScaleW = v.f32()
		case "scaleH":
			c.ScaleH = v.f32()
		case "pages":
			c.Pages = v.intValue
		case "packed":
			c.Packed = v.intValue
		}
	}
	return c
}

// This tag gives the name of a texture file. There is one for each page in the font.
type Page struct {
	Id   int    // The page id.
	File string // The texture file name.
}

func BuildPage(kvs map[string]value) Page {
	p := Page{}
	for k, v := range kvs {
		switch k {
		case "id":
			p.Id = v.intValue
		case "file":
			p.File = v.stringValue
		}
	}
	return p
}

// This tag describes on character in the font. There is one for each included character in the font.
type Char struct {
	Id       uint8   // The character id.
	X        float32 // The left position of the character image in the texture.
	Y        float32 // The top position of the character image in the texture.
	Width    float32 // The width of the character image in the texture.
	Height   float32 // The height of the character image in the texture.
	Xoffset  float32 // How much the current position should be offset when copying the image from the texture to the screen.
	Yoffset  float32 // How much the current position should be offset when copying the image from the texture to the screen.
	Xadvance float32 // How much the current position should be advanced after drawing the character.
	Page     int     //	The texture page where the character image is found.
	Channel  int     //	The texture channel where the character image is found (1 = blue, 2 = green, 4 = red, 8 = alpha, 15 = all channels).
	Letter   string  // The letter is represents
}

func BuildChar(kvs map[string]value) Char {
	c := Char{}
	for k, v := range kvs {
		switch k {
		case "id":
			c.Id = uint8(v.intValue)
		case "x":
			c.X = v.f32()
		case "y":
			c.Y = v.f32()
		case "width":
			c.Width = v.f32()
		case "height":
			c.Height = v.f32()
		case "xoffset":
			c.Xoffset = v.f32()
		case "yoffset":
			c.Yoffset = v.f32()
		case "xadvance":
			c.Xadvance = v.f32()
		case "page":
			c.Page = v.intValue
		case "chnl":
			c.Channel = v.intValue
		case "letter":
			c.Letter = v.stringValue
		}
	}
	return c
}

// The kerning information is used to adjust the distance between certain characters, e.g. some characters should be placed closer to each other than others.
type Kerning struct {
	First  uint8   // The first character id.
	Second uint8   // The second character id.
	Amount float32 // How much the x position should be adjusted when drawing the second character immediately following the first.
}

func BuildKerning(kvs map[string]value) Kerning {
	g := Kerning{}
	for k, v := range kvs {
		switch k {
		case "first":
			g.First = uint8(v.intValue)
		case "second":
			g.Second = uint8(v.intValue)
		case "amount":
			g.Amount = v.f32()
		}
	}
	return g
}
