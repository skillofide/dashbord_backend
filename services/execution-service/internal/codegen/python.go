package codegen

import (
	"fmt"
	"strings"
)

// pyType spells a Type as a Python annotation.
func pyType(t Type) string {
	if name, ok := pyLinkedType(t); ok {
		return name
	}
	switch t {
	case TInt, TLong:
		return "int"
	case TDouble:
		return "float"
	case TBool:
		return "bool"
	case TString, TChar:
		return "str"
	case TIntArr, TLongArr:
		return "list[int]"
	case TDblArr:
		return "list[float]"
	case TBoolArr:
		return "list[bool]"
	case TStrArr:
		return "list[str]"
	case TInt2D:
		return "list[list[int]]"
	case TStr2D:
		return "list[list[str]]"
	case TVoid:
		return "None"
	}
	return "object"
}

// StarterPython renders the Python skeleton.
func StarterPython(s Signature) string {
	parts := make([]string, len(s.Params))
	for i, p := range s.Params {
		parts[i] = fmt.Sprintf("%s: %s", p.Name, pyType(p.Type))
	}

	var b strings.Builder
	fmt.Fprintf(&b, "def %s(%s) -> %s:\n", s.EntryPoint, strings.Join(parts, ", "), pyType(s.ReturnType))
	b.WriteString("    # Write your code here\n")
	if s.ReturnType == TVoid {
		b.WriteString("    pass\n")
	} else {
		fmt.Fprintf(&b, "    %s\n", pyZeroReturn(s.ReturnType))
	}
	return prependPyLinkedPrelude(s, b.String())
}

func pyZeroReturn(t Type) string {
	switch t {
	case TInt, TLong:
		return "return 0"
	case TDouble:
		return "return 0.0"
	case TBool:
		return "return False"
	case TString, TChar:
		return `return ""`
	case TIntArr, TLongArr, TDblArr, TBoolArr, TStrArr, TInt2D, TStr2D:
		return "return []"
	}
	return "return None"
}

// DriverPython renders the Python harness.
//
// json.dumps, never print() on the raw value: Python renders a list as
// "[12, 8, 20, 5]" and a bool as "True", neither of which is JSON. The old
// arithmetic driver printed the list directly and every correct Python
// submission to that problem was marked wrong on the spaces alone.
//
// separators=(",", ":") because json.dumps defaults to ", " and would emit
// "[0, 1]" where JavaScript emits "[0,1]". The typed comparator accepts both,
// but the console shows actual next to expected — a learner who passes should
// not be looking at two strings that differ.
//
// The read is sys.stdin.read().strip(), not .trim(). Three of the hand-written
// Python drivers called .trim(), which is a JavaScript method — those problems
// raised AttributeError on every submission, correct or not.
func DriverPython(s Signature) string {
	var b strings.Builder
	b.WriteString("\n# ─── generated driver ───\n")
	b.WriteString("if __name__ == \"__main__\":\n")
	b.WriteString("    import sys as _sys, json as _json\n")
	b.WriteString("    _lines = _sys.stdin.read().split(\"\\n\")\n")

	for i, p := range s.Params {
		if p.Type.Linked() {
			fmt.Fprintf(&b, "    %s = %s\n", p.Name, pyBuildExpr(p.Type, fmt.Sprintf("_json.loads(_lines[%d])", i)))
			continue
		}
		fmt.Fprintf(&b, "    %s = _json.loads(_lines[%d])\n", p.Name, i)
	}

	names := make([]string, len(s.Params))
	for i, p := range s.Params {
		names[i] = p.Name
	}
	call := fmt.Sprintf("%s(%s)", s.EntryPoint, strings.Join(names, ", "))

	if s.ReturnType == TVoid {
		fmt.Fprintf(&b, "    %s\n", call)
		if len(s.Params) > 0 {
			fmt.Fprintf(&b, "    _sys.stdout.write(_json.dumps(%s, separators=(\",\", \":\")))\n", s.Params[0].Name)
		}
	} else {
		fmt.Fprintf(&b, "    _result = %s\n", call)
		if s.ReturnType.Linked() {
			fmt.Fprintf(&b, "    _result = %s\n", pyDumpExpr(s.ReturnType, "_result"))
		}
		b.WriteString("    _sys.stdout.write(_json.dumps(_result, separators=(\",\", \":\")))\n")
	}
	return b.String()
}
