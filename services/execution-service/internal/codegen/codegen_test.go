package codegen

import "testing"

func twoSumSig() Signature {
	return Signature{
		EntryPoint: "twoSum",
		Params:     []Param{{Name: "nums", Type: TIntArr}, {Name: "target", Type: TInt}},
		ReturnType: TIntArr,
		Compare:    CompareExact,
	}
}

func TestValidate(t *testing.T) {
	if err := twoSumSig().Validate(); err != nil {
		t.Fatalf("valid signature rejected: %v", err)
	}

	bad := []struct {
		name string
		sig  Signature
	}{
		{"no entry point", Signature{ReturnType: TInt}},
		{"bad return", Signature{EntryPoint: "f", ReturnType: "GraphNode"}},
		{"bad param type", Signature{EntryPoint: "f", ReturnType: TInt, Params: []Param{{Name: "a", Type: "widget"}}}},
		{"unnamed param", Signature{EntryPoint: "f", ReturnType: TInt, Params: []Param{{Name: " ", Type: TInt}}}},
		{"duplicate param", Signature{EntryPoint: "f", ReturnType: TInt, Params: []Param{{Name: "a", Type: TInt}, {Name: "a", Type: TInt}}}},
		{"void param", Signature{EntryPoint: "f", ReturnType: TInt, Params: []Param{{Name: "a", Type: TVoid}}}},
		{"unknown compare", Signature{EntryPoint: "f", ReturnType: TInt, Compare: "vibes"}},
	}
	for _, tc := range bad {
		t.Run(tc.name, func(t *testing.T) {
			if err := tc.sig.Validate(); err == nil {
				t.Errorf("expected %s to be rejected", tc.name)
			}
		})
	}
}

// The starter must not contain a working body. Starters seeded with the
// reference solution are why most of the Easy set could be cleared by pressing
// Submit without typing anything.
func TestStarterHasNoAnswer(t *testing.T) {
	js := StarterJavaScript(twoSumSig())
	for _, giveaway := range []string{"new Map(", "for (", "map.has"} {
		if contains(js, giveaway) {
			t.Errorf("JS starter leaks an implementation: contains %q\n%s", giveaway, js)
		}
	}
	if !contains(js, "function twoSum(nums, target)") {
		t.Errorf("JS starter missing the declared signature:\n%s", js)
	}

	py := StarterPython(twoSumSig())
	if !contains(py, "def twoSum(nums: list[int], target: int) -> list[int]:") {
		t.Errorf("Python starter missing the declared signature:\n%s", py)
	}
	if contains(py, "seen") {
		t.Errorf("Python starter leaks an implementation:\n%s", py)
	}
}

// The driver must read one JSON value per line. The old one decoded the whole
// input as a single array and fell back to passing the raw string, so a
// two-line input silently reached the solution as one argument.
func TestDriverReadsOneLinePerParam(t *testing.T) {
	js := DriverJavaScript(twoSumSig())
	for _, want := range []string{
		"const nums = JSON.parse(__lines[0]);",
		"const target = JSON.parse(__lines[1]);",
		"twoSum(nums, target)",
	} {
		if !contains(js, want) {
			t.Errorf("JS driver missing %q\n%s", want, js)
		}
	}

	py := DriverPython(twoSumSig())
	for _, want := range []string{
		"nums = _json.loads(_lines[0])",
		"target = _json.loads(_lines[1])",
		"_json.dumps(_result, separators=(\",\", \":\"))",
	} {
		if !contains(py, want) {
			t.Errorf("Python driver missing %q\n%s", want, py)
		}
	}
	// .trim() is a JavaScript method; three hand-written Python drivers called
	// it on a str and raised AttributeError on every submission.
	if contains(py, ".trim()") {
		t.Errorf("Python driver calls .trim(), which does not exist on str:\n%s", py)
	}
	// print() on a list emits Python repr, not JSON.
	if contains(py, "print(") {
		t.Errorf("Python driver uses print() instead of json.dumps:\n%s", py)
	}
}

// A void entry point mutates its first argument; that argument is the answer.
func TestVoidDriverPrintsFirstArg(t *testing.T) {
	sig := Signature{
		EntryPoint: "sortColors",
		Params:     []Param{{Name: "nums", Type: TIntArr}},
		ReturnType: TVoid,
	}
	if !contains(DriverJavaScript(sig), "JSON.stringify(nums)") {
		t.Errorf("void JS driver should print the mutated argument")
	}
	if !contains(DriverPython(sig), "_json.dumps(nums, separators=") {
		t.Errorf("void Python driver should print the mutated argument")
	}
}

func contains(haystack, needle string) bool {
	return len(needle) > 0 && len(haystack) >= len(needle) && indexOf(haystack, needle) >= 0
}

func indexOf(h, n string) int {
	for i := 0; i+len(n) <= len(h); i++ {
		if h[i:i+len(n)] == n {
			return i
		}
	}
	return -1
}

// A zero-parameter entry point must not declare a reader it never uses. Go
// rejects an unused variable, so getImageTag()-shaped problems failed to
// compile until the declaration became conditional.
func TestGoDriverNoUnusedReader(t *testing.T) {
	sig := Signature{EntryPoint: "getImageTag", ReturnType: TString}
	out := WrapGo(sig, StarterGo(sig))
	if contains(out, "__rd :=") {
		t.Errorf("zero-parameter Go driver declares an unused reader:\n%s", out)
	}
	// The imports must still all be referenced.
	for _, want := range []string{"_ = __cgBufio.NewReader", "_ = __cgJSON.Marshal", "_ = __cgOS.Stdout"} {
		if !contains(out, want) {
			t.Errorf("Go driver missing %q, import would be unused:\n%s", want, out)
		}
	}

	// With parameters, the reader is required.
	sig2 := Signature{EntryPoint: "f", Params: []Param{{Name: "n", Type: TInt}}, ReturnType: TInt}
	if !contains(WrapGo(sig2, StarterGo(sig2)), "__rd :=") {
		t.Error("Go driver with parameters must declare the reader")
	}
}

// Linked structures must be built before the entry point sees them, and
// serialised back afterwards. Decoding the test data straight into the
// parameter would hand a tree problem a plain array.
func TestLinkedStructureDriver(t *testing.T) {
	tree := Signature{
		EntryPoint: "isValidBST",
		Params:     []Param{{Name: "root", Type: TTreeNode}},
		ReturnType: TBool,
	}
	drv := DriverJavaScript(tree)
	if !contains(drv, "__cgBuildTree(JSON.parse(__lines[0]))") {
		t.Errorf("tree parameter is not built from its serialised form:\n%s", drv)
	}
	starter := StarterJavaScript(tree)
	if !contains(starter, "function TreeNode(val, left, right)") {
		t.Errorf("starter does not define TreeNode:\n%s", starter)
	}

	list := Signature{
		EntryPoint: "reverseList",
		Params:     []Param{{Name: "head", Type: TListNode}},
		ReturnType: TListNode,
	}
	wrapped, err := Wrap(LangJavaScript, list, StarterJavaScript(list))
	if err != nil {
		t.Fatalf("wrap: %v", err)
	}
	if !contains(wrapped, "__cgBuildList(") || !contains(wrapped, "__cgDumpList(") {
		t.Errorf("list problem must build on the way in and serialise on the way out:\n%s", wrapped)
	}
	// The prelude has to come first so the submission can refer to ListNode.
	if indexOf(wrapped, "function ListNode(") > indexOf(wrapped, "function reverseList(") {
		t.Error("node definitions must precede the submission")
	}

	// Every language must now define the node type and build it from the
	// serialised array. A driver that decoded the array straight into the
	// parameter would hand a tree problem a plain list of integers.
	nodeMarkers := map[string][]string{
		LangJava: {"class ListNode {", "__cgBuildList("},
		LangCpp:  {"struct ListNode {", "__cgBuildList("},
		LangGo:   {"type ListNode struct {", "__cgBuildList("},
	}
	for lang, markers := range nodeMarkers {
		out, err := Wrap(lang, list, starterFor(t, lang, list))
		if err != nil {
			t.Errorf("%s refused a linked-structure problem: %v", lang, err)
			continue
		}
		for _, m := range markers {
			if !contains(out, m) {
				t.Errorf("%s driver missing %q:\n%s", lang, m, out)
			}
		}
	}
	// Non-linked signatures are unaffected.
	if _, err := Wrap(LangGo, twoSumSig(), StarterGo(twoSumSig())); err != nil {
		t.Errorf("plain signature wrongly refused: %v", err)
	}
}

func starterFor(t *testing.T, lang string, sig Signature) string {
	t.Helper()
	out, err := Starter(lang, sig)
	if err != nil {
		t.Fatalf("starter %s: %v", lang, err)
	}
	return out
}
