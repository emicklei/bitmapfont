package bitmapfont

import (
	"bufio"
	"log"
	"os"
	"strconv"
	"strings"
)

type FontReader struct {
	font *Font
}

func NewFontReader() *FontReader {
	return new(FontReader)
}

func (r *FontReader) ReadFile(filename string) (*Font, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r.font = NewFont()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		spaceSeparated := strings.Split(scanner.Text(), " ")
		tag := spaceSeparated[0]
		props := spaceSeparated[1:]
		kvs, err := parseKeyValues(props)
		if err != nil {
			return nil, err
		}
		r.buildFont(tag, kvs)
	}
	return r.font, nil
}

func (r *FontReader) buildFont(tag string, kvs map[string]value) {
	switch tag {
	case "info":
		r.font.info = buildInfo(kvs)
	case "common":
		r.font.Common = buildCommon(kvs)
	case "page":
		r.font.page = buildPage(kvs)
	case "char":
		r.font.addChar(buildChar(kvs))
	case "kerning":
		r.font.addKerning(buildKerning(kvs))
	}
}

type value struct {
	intValue    int
	stringValue string
	intArray    []int
}

func (v value) f32() float32 {
	return float32(v.intValue)
}

func parseKeyValues(props []string) (map[string]value, error) {
	kvs := map[string]value{}
	for _, each := range props {
		if len(each) > 0 {
			kv := strings.SplitN(each, "=", 2) // value may contain "="
			if len(kv) != 2 {
				log.Printf("unexpected non-key-value:[%q]\n", each)
				continue
			}
			v := value{}
			if strings.HasPrefix(kv[1], "\"") {
				v.stringValue = strings.Trim(kv[1], "\"")
			} else if strings.Contains(kv[1], ",") {
				nums := strings.Split(kv[1], ",")
				ints := []int{}
				for _, other := range nums {
					i, err := strconv.Atoi(other)
					if err != nil {
						return kvs, err
					}
					ints = append(ints, i)
				}
				v.intArray = ints
			} else {
				i, err := strconv.Atoi(kv[1])
				if err != nil {
					return kvs, err
				}
				v.intValue = i
			}
			kvs[kv[0]] = v
		}
	}
	return kvs, nil
}
