package codegen

import (
	"fmt"
	"strings"
)

// jsDocType spells a Type the way JSDoc does, so the generated starter carries
// the same type information the driver enforces.
func jsDocType(t Type) string {
	if name, ok := jsLinkedDocType(t); ok {
		return name
	}
	switch t {
	case TInt, TLong, TDouble:
		return "number"
	case TBool:
		return "boolean"
	case TString, TChar:
		return "string"
	case TIntArr, TLongArr, TDblArr:
		return "number[]"
	case TBoolArr:
		return "boolean[]"
	case TStrArr:
		return "string[]"
	case TInt2D:
		return "number[][]"
	case TStr2D:
		return "string[][]"
	case TVoid:
		return "void"
	}
	return "*"
}

// StarterJavaScript renders the skeleton the learner starts from. It contains
// the signature and nothing else — no working body. Starters used to be seeded
// with the reference solution, so most problems could be "solved" by pressing
// Submit on arrival.
func StarterJavaScript(s Signature) string {
	var b strings.Builder
	b.WriteString("/**\n")
	for _, p := range s.Params {
		fmt.Fprintf(&b, " * @param {%s} %s\n", jsDocType(p.Type), p.Name)
	}
	if s.ReturnType != TVoid {
		fmt.Fprintf(&b, " * @return {%s}\n", jsDocType(s.ReturnType))
	}
	b.WriteString(" */\n")

	names := make([]string, len(s.Params))
	for i, p := range s.Params {
		names[i] = p.Name
	}
	fmt.Fprintf(&b, "function %s(%s) {\n", s.EntryPoint, strings.Join(names, ", "))
	b.WriteString("    // Write your code here\n")
	if s.ReturnType != TVoid {
		fmt.Fprintf(&b, "    %s\n", jsZeroReturn(s.ReturnType))
	}
	b.WriteString("}\n")
	return prependLinkedPrelude(s, b.String())
}

func jsZeroReturn(t Type) string {
	switch t {
	case TInt, TLong, TDouble:
		return "return 0;"
	case TBool:
		return "return false;"
	case TString, TChar:
		return "return \"\";"
	case TIntArr, TLongArr, TDblArr, TBoolArr, TStrArr, TInt2D, TStr2D:
		return "return [];"
	}
	return "return null;"
}

// jsLinkedWrap is the driver's assembly for a submission that takes or returns
// a linked structure. The prelude has to sit above the submission, so the
// caller cannot simply append.
func jsLinkedWrap(s Signature, code string) string {
	return prependLinkedPrelude(s, code) + "\n" + DriverJavaScript(s)
}

// DriverJavaScript renders the harness appended to the submission.
//
// One JSON value per line, one line per parameter, in declaration order. This
// is the format the test cases were already stored in; the old driver tried to
// JSON-decode the entire input as a single array, which meant a two-line input
// like "[2,7,11,15]\n9" fell into a catch block and got passed to the solution
// as one raw string. Two Sum's own reference solution failed on Two Sum.
func DriverJavaScript(s Signature) string {
	var b strings.Builder
	b.WriteString("\n// ─── generated driver ───\n")
	b.WriteString("(function () {\n")
	b.WriteString("  const __src = require('fs').readFileSync(0, 'utf-8');\n")
	b.WriteString("  const __lines = __src.split('\\n');\n")

	for i, p := range s.Params {
		if p.Type.Linked() {
			// Decode the array, then build the structure the entry point expects.
			fmt.Fprintf(&b, "  const %s = %s;\n", p.Name, jsBuildExpr(p.Type, fmt.Sprintf("JSON.parse(__lines[%d])", i)))
			continue
		}
		fmt.Fprintf(&b, "  const %s = JSON.parse(__lines[%d]);\n", p.Name, i)
	}

	names := make([]string, len(s.Params))
	for i, p := range s.Params {
		names[i] = p.Name
	}
	call := fmt.Sprintf("%s(%s)", s.EntryPoint, strings.Join(names, ", "))

	if s.ReturnType == TVoid {
		// A void entry point mutates its first argument in place; that argument
		// is the answer.
		b.WriteString("  " + call + ";\n")
		if len(s.Params) > 0 {
			fmt.Fprintf(&b, "  process.stdout.write(JSON.stringify(%s));\n", jsDumpExpr(s.Params[0].Type, s.Params[0].Name))
		}
	} else {
		fmt.Fprintf(&b, "  const __result = %s;\n", call)
		if s.ReturnType.Linked() {
			fmt.Fprintf(&b, "  const __serialised = %s;\n", jsDumpExpr(s.ReturnType, "__result"))
			b.WriteString("  process.stdout.write(JSON.stringify(__serialised));\n")
			b.WriteString("})();\n")
			return b.String()
		}
		// JSON.stringify, not String(): the judge decodes both sides as JSON so
		// a string result must arrive quoted. Printing it bare is what made
		// correct answers to the string problems read as Wrong Answer.
		b.WriteString("  process.stdout.write(JSON.stringify(__result === undefined ? null : __result));\n")
	}
	b.WriteString("})();\n")
	return b.String()
}
