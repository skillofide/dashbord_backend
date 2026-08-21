package codegen

import (
	"fmt"
	"strings"
)

// Language names accepted by Wrap and Starter.
const (
	LangJavaScript = "javascript"
	LangPython     = "python"
	LangJava       = "java"
	LangCpp        = "cpp"
	LangGo         = "go"
)

// Wrap assembles a runnable program from a function-mode submission.
//
// This is the single place that decides how a submission becomes a program.
// It replaced a per-language pile of hardcoded drivers keyed on five problem
// slugs plus a reflection fallback that guessed the entry point and the
// argument list — the fallback silently mis-graded any problem with more than
// one line of input, a string return, or an entry point that was not the first
// function in the file.
func Wrap(lang string, s Signature, code string) (string, error) {
	if err := s.Validate(); err != nil {
		return "", err
	}
	switch strings.ToLower(lang) {
	case LangJavaScript:
		return jsLinkedWrap(s, code), nil
	case LangPython:
		return prependPyLinkedPrelude(s, code) + "\n" + DriverPython(s), nil
	case LangGo:
		return WrapGo(s, code), nil
	case LangJava:
		return WrapJava(s, code), nil
	case LangCpp:
		return WrapCpp(s, code), nil
	}
	return "", fmt.Errorf("codegen: no driver for language %q", lang)
}

// Starter renders the skeleton a learner begins from.
//
// Starters must be generated rather than stored per language by hand. Hand
// authoring is how 361 of the 442 Easy problems ended up shipping a working
// solution in the editor — solvable by pressing Submit on arrival — while 66
// others shipped `console.log("Hello World")` and could not be solved at all.
func Starter(lang string, s Signature) (string, error) {
	if err := s.Validate(); err != nil {
		return "", err
	}
	// Linked structures are only modelled for JavaScript and Python so far. The
	// other three would need a node type, a builder and a serialiser each, and
	// emitting a driver that treats a tree as an array would grade wrongly
	// rather than fail loudly.

	switch strings.ToLower(lang) {
	case LangJavaScript:
		return StarterJavaScript(s), nil
	case LangPython:
		return StarterPython(s), nil
	case LangGo:
		return StarterGo(s), nil
	case LangJava:
		return StarterJava(s), nil
	case LangCpp:
		return StarterCpp(s), nil
	}
	return "", fmt.Errorf("codegen: no starter for language %q", lang)
}
