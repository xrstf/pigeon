package vm

import (
	"bytes"
	"fmt"
	"strings"
	"text/template"
	"unicode"

	"github.com/PuerkitoBio/pigeon/ast"
)

const codeTpl = `// Code generated by pigeon (https://github.com/PuerkitoBio/pigeon)
// on {{.Now}}
{{if .Init}}
{{.Init}}
{{end}}{{range .As}}
func ({{$.RecvrNm}} *current) on{{.RuleNm}}{{.ExprIx}}({{range $index, $elem := .Parms}}{{if $index}}, {{end}}{{$elem}}{{end}}{{if .Parms}} interface{}{{end}}) (interface{}, error) { {{.Code}} }

func (v *ϡvm) callOn{{.RuleNm}}{{.ExprIx}}() (interface{}, error) {
{{if .Parms}}stack := v.a.peek()
{{end}}return v.cur.on{{.RuleNm}}{{.ExprIx}}({{range $index, $elem := .Parms}}{{if $index}}, {{end}}stack[{{printf "%q" $elem}}]{{end}})
}
{{end}}{{range .Bs}}
func ({{$.RecvrNm}} *current) on{{.RuleNm}}{{.ExprIx}}({{range $index, $elem := .Parms}}{{if $index}}, {{end}}{{$elem}}{{end}}{{if .Parms}} interface{}{{end}}) (bool, error) { {{.Code}} }

func (v *ϡvm) callOn{{.RuleNm}}{{.ExprIx}}() (bool, error) {
{{if .Parms}}stack := v.a.peek()
{{end}}return v.cur.on{{.RuleNm}}{{.ExprIx}}({{range $index, $elem := .Parms}}{{if $index}}, {{end}}stack[{{printf "%q" $elem}}]{{end}})
}
{{end}}
var ϡtheProgram = &ϡprogram{
instrs: []ϡinstr{
{{range $index, $elem := .Instrs}}{{if (not (mod $index 3))}}{{printf "\n"}}{{end}}{{$elem}}, {{end}}
},
instrToRule: []int{
{{range $index, $elem := .InstrToRule}}{{if (not (mod $index 10))}}{{printf "\n"}}{{end}}{{$elem}}, {{end}}
},
ss: []string{
{{range $index, $elem := .Ss}}{{if (not (mod $index 3))}}{{printf "\n"}}{{end}}{{printf "%q" $elem}}, {{end}}
},
ms: []ϡmatcher{
{{range .Ms}}{{matcher .}},
{{end}}},
as: []func(*ϡvm) (interface{}, error){
{{range .As}}(*ϡvm).callOn{{.RuleNm}}{{.ExprIx}},
{{end}}},
bs: []func(*ϡvm) (bool, error){
{{range .Bs}}(*ϡvm).callOn{{.RuleNm}}{{.ExprIx}},
{{end}}},
}
`

var funcMap = template.FuncMap{
	"mod": func(ix, div int) int {
		if ix == 0 {
			return 1
		}
		return ix % div
	},
	"matcher": func(m ast.Expression) string {
		switch m := m.(type) {
		case *ast.AnyMatcher:
			return "ϡanyMatcher{}"
		case *ast.LitMatcher:
			if m.IgnoreCase {
				m.Val = strings.ToLower(m.Val)
			}
			return fmt.Sprintf("ϡstringMatcher{\nignoreCase: %t,\nvalue: %q,\n}",
				m.IgnoreCase, m.Val)
		case *ast.CharClassMatcher:
			if m.IgnoreCase {
				for j, rn := range m.Chars {
					m.Chars[j] = unicode.ToLower(rn)
				}
				for j, rn := range m.Ranges {
					m.Ranges[j] = unicode.ToLower(rn)
				}
			}
			var buf bytes.Buffer
			buf.WriteString(fmt.Sprintf("ϡcharClassMatcher{\nignoreCase: %t,\ninverted: %t,\nchars: []rune{", m.IgnoreCase, m.Inverted))
			for i, rn := range m.Chars {
				if i > 0 {
					buf.WriteString(", ")
				}
				buf.WriteString(fmt.Sprintf("%q", rn))
			}
			buf.WriteString("},\nranges: []rune{")
			for i, rn := range m.Ranges {
				if i > 0 {
					buf.WriteString(", ")
				}
				buf.WriteString(fmt.Sprintf("%q", rn))
			}
			buf.WriteString("},\nclasses: []*unicode.RangeTable{")
			for i, cl := range m.UnicodeClasses {
				if i > 0 {
					buf.WriteString(", ")
				}
				buf.WriteString(fmt.Sprintf("ϡrangeTable(%q)", cl))
			}
			buf.WriteString("},\n}")
			return buf.String()
		default:
			panic(fmt.Sprintf("unknown matcher type %T", m))
		}
	},
}

var tpl = template.New("gen")

func init() {
	template.Must(tpl.Funcs(funcMap).Parse(codeTpl))
}
