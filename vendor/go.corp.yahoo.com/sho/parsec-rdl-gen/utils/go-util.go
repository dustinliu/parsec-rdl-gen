package utils

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/ardielle/ardielle-go/rdl"
	"os"
	"path/filepath"
	"strings"
)

func optionalAnyToString(any interface{}) string {
	if any == nil {
		return "null"
	}
	switch v := any.(type) {
	case *bool:
		return fmt.Sprintf("%v", *v)
	case *int8:
		return fmt.Sprintf("%d", *v)
	case *int16:
		return fmt.Sprintf("%d", *v)
	case *int32:
		return fmt.Sprintf("%d", *v)
	case *int64:
		return fmt.Sprintf("%d", *v)
	case *float32:
		return fmt.Sprintf("%g", *v)
	case *float64:
		return fmt.Sprintf("%g", *v)
	case *string:
		return *v
	case bool:
		return fmt.Sprintf("%v", v)
	case int8:
		return fmt.Sprintf("%d", v)
	case int16:
		return fmt.Sprintf("%d", v)
	case int32:
		return fmt.Sprintf("%d", v)
	case int64:
		return fmt.Sprintf("%d", v)
	case float32:
		return fmt.Sprintf("%g", v)
	case float64:
		return fmt.Sprintf("%g", v)
	case string:
		return fmt.Sprintf("%v", v)
	default:
		panic("optionalAnyToString")
	}
}

func formatBlock(s string, leftCol int, rightCol int, prefix string) string {
	if s == "" {
		return ""
	}
	tab := spaces(leftCol)
	var buf bytes.Buffer
	max := 80
	col := leftCol
	lines := 1
	tokens := strings.Split(s, " ")
	for _, tok := range tokens {
		toklen := len(tok)
		if col+toklen >= max {
			buf.WriteString("\n")
			lines++
			buf.WriteString(tab)
			buf.WriteString(prefix)
			buf.WriteString(tok)
			col = leftCol + 3 + toklen
		} else {
			if col == leftCol {
				col += len(prefix)
				buf.WriteString(prefix)
			} else {
				buf.WriteString(" ")
			}
			buf.WriteString(tok)
			col += toklen + 1
		}
	}
	buf.WriteString("\n")
	emptyPrefix := strings.Trim(prefix, " ")
	pad := tab + emptyPrefix + "\n"
	return pad + buf.String() + pad
}

func FormatComment(s string, leftCol int, rightCol int) string {
	return formatBlock(s, leftCol, rightCol, "// ")
}

func spaces(count int) string {
	return stringOfChar(count, ' ')
}

func stringOfChar(count int, b byte) string {
	buf := make([]byte, 0, count)
	for i := 0; i < count; i++ {
		buf = append(buf, b)
	}
	return string(buf)
}

func addFields(reg rdl.TypeRegistry, dst []*rdl.StructFieldDef, t *rdl.Type) []*rdl.StructFieldDef {
	switch t.Variant {
	case rdl.TypeVariantStructTypeDef:
		st := t.StructTypeDef
		if st.Type != "Struct" {
			dst = addFields(reg, dst, reg.FindType(st.Type))
		}
		for _, f := range st.Fields {
			dst = append(dst, f)
		}
	}
	return dst
}

func FlattenedFields(reg rdl.TypeRegistry, t *rdl.Type) []*rdl.StructFieldDef {
	return addFields(reg, make([]*rdl.StructFieldDef, 0), t)
}

func Capitalize(text string) string {
	return strings.ToUpper(text[0:1]) + text[1:]
}

func Uncapitalize(text string) string {
	return strings.ToLower(text[0:1]) + text[1:]
}

func LeftJustified(text string, width int) string {
	return text + spaces(width-len(text))
}

func OutputWriter(outdir string, name string, ext string) (*bufio.Writer, *os.File, string, error) {
	sname, path := GetOutputPathInfo(outdir, name, ext)
	if path == "" {
		return bufio.NewWriter(os.Stdout), nil, sname, nil
	}
	f, err := os.Create(path)
	if err != nil {
		return nil, nil, "", err
	}
	writer := bufio.NewWriter(f)
	return writer, f, sname, nil
}

func GetOutputPathInfo(outdir string, name string, ext string) (string, string) {
	sname := "anonymous"
	if strings.HasSuffix(outdir, ext) {
		name = filepath.Base(outdir)
		sname = name[:len(name)-len(ext)]
		outdir = filepath.Dir(outdir)
	}
	if name != "" {
		sname = name
	}
	if outdir == "" {
		return sname, ""
	}
	outfile := sname
	if !strings.HasSuffix(outfile, ext) {
		outfile += ext
	}
	path := filepath.Join(outdir, outfile)
	return sname, path
}

func generationHeader(banner string) string {
	return fmt.Sprintf("//\n// This file generated by %s\n//", banner)
}

func generationPackage(schema *rdl.Schema) string {
	pkg := "main"
	if schema.Name != "" {
		pkg = strings.ToLower(string(schema.Name))
	}
	return pkg
}

/*
// generate common runtime code for client and server
func generateUtil(schema *rdl.Schema, writer io.Writer) error {
	basenameFunc := func(s string) string {
		i := strings.LastIndex(s, ".")
		if i >= 0 {
			s = s[i+1:]
		}
		return s
	}
	funcMap := template.FuncMap{
		"basename": basenameFunc,
	}
	t := template.Must(template.New("util").Funcs(funcMap).Parse(utilTemplate))
	return t.Execute(writer, schema)
}

const utilTemplate = `
`
*/

func Split(str string, delim rune) []string {
	pieces := make([]string, 0)
	escaped := false
	offset := 0
	length := len(str)

	for pos, rune := range str {
		if pos + 1 == length && rune != delim {
			pieces = append(pieces, str[offset:])
		} else {
			if pos != 0 && rune == '"' && str[pos - 1] != '\\' {
				escaped = !escaped
			} else if rune == delim && !escaped {
				pieces = append(pieces, str[offset:pos])
				offset = pos + 1
			}
		}
	}
	return pieces
}
