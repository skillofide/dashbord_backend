package codegen

import (
	"fmt"
	"strings"
)

// goType spells a Type as Go source.
func goType(t Type) string {
	if name, ok := goLinkedType(t); ok {
		return name
	}
	switch t {
	case TInt:
		return "int"
	case TLong:
		return "int64"
	case TDouble:
		return "float64"
	case TBool:
		return "bool"
	case TString:
		return "string"
	case TChar:
		return "byte"
	case TIntArr:
		return "[]int"
	case TLongArr:
		return "[]int64"
	case TDblArr:
		return "[]float64"
	case TBoolArr:
		return "[]bool"
	case TStrArr:
		return "[]string"
	case TInt2D:
		return "[][]int"
	case TStr2D:
		return "[][]string"
	case TVoid:
		return ""
	}
	return "interface{}"
}

// StarterGo renders the Go skeleton for a function-mode problem.
func StarterGo(s Signature) string {
	parts := make([]string, len(s.Params))
	for i, p := range s.Params {
		parts[i] = fmt.Sprintf("%s %s", p.Name, goType(p.Type))
	}
	ret := goType(s.ReturnType)
	if ret != "" {
		ret = " " + ret
	}

	var b strings.Builder
	b.WriteString("package main\n\n")
	fmt.Fprintf(&b, "func %s(%s)%s {\n", s.EntryPoint, strings.Join(parts, ", "), ret)
	b.WriteString("\t// Write your code here\n")
	if s.ReturnType != TVoid {
		fmt.Fprintf(&b, "\t%s\n", goZeroReturn(s.ReturnType))
	}
	b.WriteString("}\n")
	return appendGoLinkedTypes(s, b.String())
}

func goZeroReturn(t Type) string {
	switch t {
	case TInt, TLong, TChar:
		return "return 0"
	case TDouble:
		return "return 0.0"
	case TBool:
		return "return false"
	case TString:
		return `return ""`
	case TIntArr, TLongArr, TDblArr, TBoolArr, TStrArr, TInt2D, TStr2D:
		return "return nil"
	}
	return "return nil"
}

// WrapGo assembles a complete Go program from a function-mode submission.
//
// Imports cannot simply be appended: Go requires every import declaration to
// precede the first top-level declaration, so a block added after the user's
// function is a syntax error. They are injected directly after "package main"
// instead.
//
// Each import is aliased with a __cg prefix. Go permits the same package to be
// imported more than once under different names, so a learner who writes
// `import "os"` themselves does not collide with the driver's own os import —
// and the prefix keeps the driver's names out of the namespace they are
// working in.
//
// Only used for io_mode = 'function'. A submission that already declares main()
// is a complete stdio program and is run verbatim; appending this would give it
// a second main, which is exactly what used to break the whole go-* family.
func WrapGo(s Signature, code string) string {
	imports := "\nimport (\n\t__cgBufio \"bufio\"\n\t__cgJSON \"encoding/json\"\n\t__cgOS \"os\"\n)\n"

	body := code
	if idx := strings.Index(body, "package main"); idx >= 0 {
		cut := idx + len("package main")
		body = body[:cut] + "\n" + imports + body[cut:]
	} else {
		body = "package main\n" + imports + "\n" + body
	}

	return appendGoLinkedTypes(s, body+"\n"+driverGoMain(s))
}

func driverGoMain(s Signature) string {
	var b strings.Builder
	b.WriteString("// ─── generated driver ───\n")
	b.WriteString("func main() {\n")
	// Go rejects an unused import, and a zero-parameter entry point never
	// touches the reader. These blank references keep every injected import
	// used whatever the signature is.
	//
	// They live inside main deliberately. As a package-level `var` block they
	// were a declaration sitting between the injected imports and the
	// submission's own import line, which made the submission's import illegal:
	// "imports must appear before other declarations".
	b.WriteString("\t_ = __cgBufio.NewReader\n\t_ = __cgJSON.Marshal\n\t_ = __cgOS.Stdout\n")
	// Only declare the reader when something is going to read from it. Go
	// rejects an unused *variable* as firmly as an unused import, so a
	// zero-parameter entry point — getImageTag(), getFlexCenterCSS() — failed
	// to compile with "__rd declared and not used".
	if len(s.Params) > 0 {
		b.WriteString("\t__rd := __cgBufio.NewReader(__cgOS.Stdin)\n")
	}

	for i, p := range s.Params {
		fmt.Fprintf(&b, "\t__l%d, _ := __rd.ReadString('\\n')\n", i)
		if _, linked := goLinkedType(p.Type); linked {
			builder := "__cgBuildTree"
			if p.Type == TListNode {
				builder = "__cgBuildList"
			}
			fmt.Fprintf(&b, "\t%s := %s([]byte(__l%d))\n", p.Name, builder, i)
			continue
		}
		fmt.Fprintf(&b, "\tvar %s %s\n", p.Name, goType(p.Type))
		fmt.Fprintf(&b, "\t__cgJSON.Unmarshal([]byte(__l%d), &%s)\n", i, p.Name)
	}

	names := make([]string, len(s.Params))
	for i, p := range s.Params {
		names[i] = p.Name
	}
	call := fmt.Sprintf("%s(%s)", s.EntryPoint, strings.Join(names, ", "))

	if s.ReturnType == TVoid {
		fmt.Fprintf(&b, "\t%s\n", call)
		if len(s.Params) > 0 {
			fmt.Fprintf(&b, "\t__out, _ := __cgJSON.Marshal(%s)\n", s.Params[0].Name)
			b.WriteString("\t__cgOS.Stdout.Write(__out)\n")
		}
	} else {
		fmt.Fprintf(&b, "\t__result := %s\n", call)
		if _, linked := goLinkedType(s.ReturnType); linked {
			dumper := "__cgDumpTree"
			if s.ReturnType == TListNode {
				dumper = "__cgDumpList"
			}
			fmt.Fprintf(&b, "\t__cgOS.Stdout.Write(%s(__result))\n", dumper)
		} else {
			b.WriteString("\t__out, _ := __cgJSON.Marshal(__result)\n")
			b.WriteString("\t__cgOS.Stdout.Write(__out)\n")
		}
	}
	b.WriteString("}\n")
	return b.String()
}
