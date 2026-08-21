package codegen

import (
	"fmt"
	"strings"
)

// Method is one callable on a class-mode problem's type.
type Method struct {
	Name       string  `json:"name"`
	Params     []Param `json:"params"`
	ReturnType Type    `json:"returnType"`
}

// ClassSignature describes a problem built around a type rather than a
// function: a constructor plus methods, graded by driving a sequence of calls
// against one instance.
//
// LRUCache, MinStack, Trie and TicTacToe are all this shape. State carried
// between calls is the point of each of them — `put` then `get` cannot be
// expressed as a single function call — so they cannot be squeezed into a
// function signature without changing what they teach.
type ClassSignature struct {
	// ClassName is the type the learner defines, e.g. "LRUCache".
	ClassName string   `json:"className"`
	Ctor      []Param  `json:"ctor"`
	Methods   []Method `json:"methods"`
}

// Validate rejects a class signature the generators cannot emit faithfully.
func (c ClassSignature) Validate() error {
	if strings.TrimSpace(c.ClassName) == "" {
		return fmt.Errorf("codegen: class has no name")
	}
	if len(c.Methods) == 0 {
		return fmt.Errorf("codegen: class %s declares no methods", c.ClassName)
	}
	for _, p := range c.Ctor {
		if !p.Type.Valid() || p.Type == TVoid {
			return fmt.Errorf("codegen: constructor parameter %q has unsupported type %q", p.Name, p.Type)
		}
	}
	seen := map[string]bool{}
	for _, m := range c.Methods {
		if strings.TrimSpace(m.Name) == "" {
			return fmt.Errorf("codegen: class %s has an unnamed method", c.ClassName)
		}
		if seen[m.Name] {
			return fmt.Errorf("codegen: class %s declares %q twice", c.ClassName, m.Name)
		}
		seen[m.Name] = true
		if !m.ReturnType.Valid() {
			return fmt.Errorf("codegen: %s.%s has unsupported return type %q", c.ClassName, m.Name, m.ReturnType)
		}
		for _, p := range m.Params {
			if !p.Type.Valid() || p.Type == TVoid {
				return fmt.Errorf("codegen: %s.%s parameter %q has unsupported type %q", c.ClassName, m.Name, p.Name, p.Type)
			}
		}
	}
	return nil
}

// StarterClassJavaScript renders the skeleton for a class-mode problem.
func StarterClassJavaScript(c ClassSignature) string {
	var b strings.Builder
	ctorNames := paramNames(c.Ctor)

	b.WriteString("/**\n")
	for _, p := range c.Ctor {
		fmt.Fprintf(&b, " * @param {%s} %s\n", jsDocType(p.Type), p.Name)
	}
	b.WriteString(" */\n")
	fmt.Fprintf(&b, "function %s(%s) {\n", c.ClassName, strings.Join(ctorNames, ", "))
	b.WriteString("    // Set up your state here\n")
	b.WriteString("}\n")

	for _, m := range c.Methods {
		names := paramNames(m.Params)
		b.WriteString("\n/**\n")
		for _, p := range m.Params {
			fmt.Fprintf(&b, " * @param {%s} %s\n", jsDocType(p.Type), p.Name)
		}
		if m.ReturnType != TVoid {
			fmt.Fprintf(&b, " * @return {%s}\n", jsDocType(m.ReturnType))
		}
		b.WriteString(" */\n")
		fmt.Fprintf(&b, "%s.prototype.%s = function(%s) {\n", c.ClassName, m.Name, strings.Join(names, ", "))
		b.WriteString("    // Write your code here\n")
		if m.ReturnType != TVoid {
			fmt.Fprintf(&b, "    %s\n", jsZeroReturn(m.ReturnType))
		}
		b.WriteString("};\n")
	}
	return b.String()
}

// DriverClassJavaScript renders the harness for a class-mode problem.
//
// Two lines of input: the call sequence and the arguments for each call. The
// first entry names the class and carries the constructor's arguments; every
// later entry is a method call on that one instance. Output is the list of
// return values, null wherever a method returns nothing — which is what makes
// the sequence checkable rather than just the final state.
func DriverClassJavaScript(c ClassSignature) string {
	var b strings.Builder
	b.WriteString("\n// ─── generated driver ───\n")
	b.WriteString("(function () {\n")
	b.WriteString("  const __lines = require('fs').readFileSync(0, 'utf-8').split('\\n');\n")
	b.WriteString("  const __calls = JSON.parse(__lines[0]);\n")
	b.WriteString("  const __args = JSON.parse(__lines[1]);\n")
	fmt.Fprintf(&b, "  const __instance = new %s(...(__args[0] || []));\n", c.ClassName)
	b.WriteString("  const __out = [null];\n")
	b.WriteString("  for (let i = 1; i < __calls.length; i++) {\n")
	b.WriteString("    const __r = __instance[__calls[i]](...(__args[i] || []));\n")
	b.WriteString("    __out.push(__r === undefined ? null : __r);\n")
	b.WriteString("  }\n")
	b.WriteString("  process.stdout.write(JSON.stringify(__out));\n")
	b.WriteString("})();\n")
	return b.String()
}

// StarterClassPython renders the Python skeleton.
func StarterClassPython(c ClassSignature) string {
	var b strings.Builder
	fmt.Fprintf(&b, "class %s:\n", c.ClassName)

	ctorParts := []string{"self"}
	for _, p := range c.Ctor {
		ctorParts = append(ctorParts, fmt.Sprintf("%s: %s", p.Name, pyType(p.Type)))
	}
	fmt.Fprintf(&b, "    def __init__(%s):\n", strings.Join(ctorParts, ", "))
	b.WriteString("        # Set up your state here\n")
	b.WriteString("        pass\n")

	for _, m := range c.Methods {
		parts := []string{"self"}
		for _, p := range m.Params {
			parts = append(parts, fmt.Sprintf("%s: %s", p.Name, pyType(p.Type)))
		}
		fmt.Fprintf(&b, "\n    def %s(%s) -> %s:\n", m.Name, strings.Join(parts, ", "), pyType(m.ReturnType))
		b.WriteString("        # Write your code here\n")
		if m.ReturnType == TVoid {
			b.WriteString("        pass\n")
		} else {
			fmt.Fprintf(&b, "        %s\n", pyZeroReturn(m.ReturnType))
		}
	}
	return b.String()
}

// DriverClassPython renders the Python harness.
func DriverClassPython(c ClassSignature) string {
	var b strings.Builder
	b.WriteString("\n# ─── generated driver ───\n")
	b.WriteString("if __name__ == \"__main__\":\n")
	b.WriteString("    import sys as _sys, json as _json\n")
	b.WriteString("    _lines = _sys.stdin.read().split(\"\\n\")\n")
	b.WriteString("    _calls = _json.loads(_lines[0])\n")
	b.WriteString("    _args = _json.loads(_lines[1])\n")
	fmt.Fprintf(&b, "    _instance = %s(*(_args[0] or []))\n", c.ClassName)
	b.WriteString("    _out = [None]\n")
	b.WriteString("    for _i in range(1, len(_calls)):\n")
	b.WriteString("        _r = getattr(_instance, _calls[_i])(*(_args[_i] or []))\n")
	b.WriteString("        _out.append(_r)\n")
	b.WriteString("    _sys.stdout.write(_json.dumps(_out, separators=(\",\", \":\")))\n")
	return b.String()
}

// WrapClass assembles a runnable program for a class-mode submission.
//
// Only JavaScript and Python for now: both dispatch a method by name at
// runtime, which is what the call sequence needs. Java, C++ and Go would each
// need a generated switch over the method names with per-method argument
// decoding, and emitting a driver that silently skipped an unknown call would
// grade wrongly rather than fail loudly.
func WrapClass(lang string, c ClassSignature, code string) (string, error) {
	if err := c.Validate(); err != nil {
		return "", err
	}
	switch strings.ToLower(lang) {
	case LangJavaScript:
		return code + "\n" + DriverClassJavaScript(c), nil
	case LangPython:
		return code + "\n" + DriverClassPython(c), nil
	}
	return "", fmt.Errorf("codegen: %s does not support class-design problems yet", lang)
}

// StarterClass renders the skeleton for a class-mode problem.
func StarterClass(lang string, c ClassSignature) (string, error) {
	if err := c.Validate(); err != nil {
		return "", err
	}
	switch strings.ToLower(lang) {
	case LangJavaScript:
		return StarterClassJavaScript(c), nil
	case LangPython:
		return StarterClassPython(c), nil
	}
	return "", fmt.Errorf("codegen: %s does not support class-design problems yet", lang)
}

func paramNames(ps []Param) []string {
	out := make([]string, len(ps))
	for i, p := range ps {
		out[i] = p.Name
	}
	return out
}
