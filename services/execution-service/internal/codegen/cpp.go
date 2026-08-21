package codegen

import (
	"fmt"
	"strings"
)

// cppType spells a Type as C++ source.
func cppType(t Type) string {
	if name, ok := cppLinkedType(t); ok {
		return name
	}
	switch t {
	case TInt:
		return "int"
	case TLong:
		return "long long"
	case TDouble:
		return "double"
	case TBool:
		return "bool"
	case TString:
		return "string"
	case TChar:
		return "char"
	case TIntArr:
		return "vector<int>"
	case TLongArr:
		return "vector<long long>"
	case TDblArr:
		return "vector<double>"
	case TBoolArr:
		return "vector<bool>"
	case TStrArr:
		return "vector<string>"
	case TInt2D:
		return "vector<vector<int>>"
	case TStr2D:
		return "vector<vector<string>>"
	case TVoid:
		return "void"
	}
	return "auto"
}

// StarterCpp renders the C++ skeleton.
func StarterCpp(s Signature) string {
	parts := make([]string, len(s.Params))
	for i, p := range s.Params {
		t := cppType(p.Type)
		// Containers by reference, scalars by value — the convention learners
		// will have seen everywhere else.
		// Pointers to linked structures are passed as-is; only containers and
		// strings take a reference.
		if _, linked := cppLinkedType(p.Type); !linked && (strings.HasPrefix(t, "vector") || t == "string") {
			t += "&"
		}
		parts[i] = fmt.Sprintf("%s %s", t, p.Name)
	}

	var b strings.Builder
	b.WriteString("#include <bits/stdc++.h>\nusing namespace std;\n\n")
	b.WriteString(prependCppLinkedTypes(s, ""))
	b.WriteString("class Solution {\npublic:\n")
	fmt.Fprintf(&b, "    %s %s(%s) {\n", cppType(s.ReturnType), s.EntryPoint, strings.Join(parts, ", "))
	b.WriteString("        // Write your code here\n")
	if s.ReturnType != TVoid {
		fmt.Fprintf(&b, "        %s\n", cppZeroReturn(s.ReturnType))
	}
	b.WriteString("    }\n};\n")
	return b.String()
}

func cppZeroReturn(t Type) string {
	switch t {
	case TInt, TLong, TChar:
		return "return 0;"
	case TDouble:
		return "return 0.0;"
	case TBool:
		return "return false;"
	case TString:
		return "return \"\";"
	case TIntArr, TLongArr, TDblArr, TBoolArr, TStrArr, TInt2D, TStr2D:
		return "return {};"
	}
	return "return {};"
}

func cppParseExpr(t Type, src string) string {
	if expr := cppBuildExpr(t, src); expr != "" {
		return expr
	}
	switch t {
	case TInt:
		return fmt.Sprintf("stoi(%s)", src)
	case TLong:
		return fmt.Sprintf("stoll(%s)", src)
	case TDouble:
		return fmt.Sprintf("stod(%s)", src)
	case TBool:
		return fmt.Sprintf("(__cgTrimWs(%s) == \"true\")", src)
	case TChar:
		return fmt.Sprintf("__cgStr(%s)[0]", src)
	case TString:
		return fmt.Sprintf("__cgStr(%s)", src)
	case TIntArr:
		return fmt.Sprintf("__cgIntArr(%s)", src)
	case TLongArr:
		return fmt.Sprintf("__cgLongArr(%s)", src)
	case TDblArr:
		return fmt.Sprintf("__cgDblArr(%s)", src)
	case TBoolArr:
		return fmt.Sprintf("__cgBoolArr(%s)", src)
	case TStrArr:
		return fmt.Sprintf("__cgStrArr(%s)", src)
	case TInt2D:
		return fmt.Sprintf("__cgInt2D(%s)", src)
	case TStr2D:
		return fmt.Sprintf("__cgStr2D(%s)", src)
	}
	return "{}"
}

// WrapCpp appends the helper block and a generated main to the submission.
func WrapCpp(s Signature, code string) string {
	var b strings.Builder
	// The helpers below use string and vector, so the driver cannot depend on
	// the submission having included them. Prepending is safe: <bits/stdc++.h>
	// carries an include guard, and a repeated `using namespace std` is legal.
	b.WriteString("#include <bits/stdc++.h>\nusing namespace std;\n")
	b.WriteString(prependCppLinkedTypes(s, code))
	b.WriteString("\n// ─── generated driver ───\n")
	b.WriteString(cppHelpers)
	if needsLinkedPrelude(s) {
		b.WriteString(cppLinkedHelpers)
	}
	b.WriteString("int main() {\n")
	b.WriteString("    ios::sync_with_stdio(false);\n")

	names := make([]string, len(s.Params))
	for i, p := range s.Params {
		fmt.Fprintf(&b, "    string __l%d; getline(cin, __l%d);\n", i, i)
		fmt.Fprintf(&b, "    %s %s = %s;\n", cppType(p.Type), p.Name, cppParseExpr(p.Type, fmt.Sprintf("__l%d", i)))
		names[i] = p.Name
	}

	call := fmt.Sprintf("__sol.%s(%s)", s.EntryPoint, strings.Join(names, ", "))
	b.WriteString("    Solution __sol;\n")
	if s.ReturnType == TVoid {
		fmt.Fprintf(&b, "    %s;\n", call)
		if len(s.Params) > 0 {
			fmt.Fprintf(&b, "    cout << __cgJson(%s);\n", s.Params[0].Name)
		}
	} else {
		fmt.Fprintf(&b, "    auto __result = %s;\n", call)
		if dump := cppDumpExpr(s.ReturnType, "__result"); dump != "" {
			fmt.Fprintf(&b, "    cout << %s;\n", dump)
		} else {
			b.WriteString("    cout << __cgJson(__result);\n")
		}
	}
	b.WriteString("    return 0;\n}\n")
	return b.String()
}

// cppHelpers mirror the Java ones: parse against the declared type rather than
// embedding a general JSON parser, and serialise by hand so output is byte
// identical to what the other four languages emit.
const cppHelpers = `
static string __cgTrimWs(const string& s) {
    size_t a = s.find_first_not_of(" \t\r\n");
    if (a == string::npos) return "";
    size_t b = s.find_last_not_of(" \t\r\n");
    return s.substr(a, b - a + 1);
}
static string __cgInner(const string& s) {
    string t = __cgTrimWs(s);
    if (t.size() >= 2 && t.front() == '[' && t.back() == ']') t = t.substr(1, t.size() - 2);
    return __cgTrimWs(t);
}
static vector<string> __cgSplit(const string& s) {
    vector<string> out; int depth = 0; bool q = false; string cur;
    for (size_t i = 0; i < s.size(); i++) {
        char c = s[i];
        if (c == '"' && (i == 0 || s[i-1] != '\\')) q = !q;
        if (!q && (c == '[' || c == '{')) depth++;
        if (!q && (c == ']' || c == '}')) depth--;
        if (c == ',' && depth == 0 && !q) { out.push_back(cur); cur.clear(); }
        else cur += c;
    }
    if (!cur.empty()) out.push_back(cur);
    return out;
}
static string __cgStr(const string& s) {
    string t = __cgTrimWs(s);
    if (t.size() >= 2 && t.front() == '"' && t.back() == '"') t = t.substr(1, t.size() - 2);
    string o; 
    for (size_t i = 0; i < t.size(); i++) {
        if (t[i] == '\\' && i + 1 < t.size()) {
            char n = t[++i];
            if (n == 'n') o += '\n'; else if (n == 't') o += '\t'; else o += n;
        } else o += t[i];
    }
    return o;
}
static vector<int> __cgIntArr(const string& s) {
    string b = __cgInner(s); vector<int> v; if (b.empty()) return v;
    for (auto& p : __cgSplit(b)) v.push_back(stoi(__cgTrimWs(p)));
    return v;
}
static vector<long long> __cgLongArr(const string& s) {
    string b = __cgInner(s); vector<long long> v; if (b.empty()) return v;
    for (auto& p : __cgSplit(b)) v.push_back(stoll(__cgTrimWs(p)));
    return v;
}
static vector<double> __cgDblArr(const string& s) {
    string b = __cgInner(s); vector<double> v; if (b.empty()) return v;
    for (auto& p : __cgSplit(b)) v.push_back(stod(__cgTrimWs(p)));
    return v;
}
static vector<bool> __cgBoolArr(const string& s) {
    string b = __cgInner(s); vector<bool> v; if (b.empty()) return v;
    for (auto& p : __cgSplit(b)) v.push_back(__cgTrimWs(p) == "true");
    return v;
}
static vector<string> __cgStrArr(const string& s) {
    string b = __cgInner(s); vector<string> v; if (b.empty()) return v;
    for (auto& p : __cgSplit(b)) v.push_back(__cgStr(p));
    return v;
}
static vector<vector<int>> __cgInt2D(const string& s) {
    string b = __cgInner(s); vector<vector<int>> v; if (b.empty()) return v;
    for (auto& p : __cgSplit(b)) v.push_back(__cgIntArr(p));
    return v;
}
static vector<vector<string>> __cgStr2D(const string& s) {
    string b = __cgInner(s); vector<vector<string>> v; if (b.empty()) return v;
    for (auto& p : __cgSplit(b)) v.push_back(__cgStrArr(p));
    return v;
}
static string __cgEsc(const string& s) {
    string o;
    for (char c : s) {
        if (c == '"') o += "\\\""; else if (c == '\\') o += "\\\\";
        else if (c == '\n') o += "\\n"; else if (c == '\t') o += "\\t";
        else o += c;
    }
    return o;
}
static string __cgJson(bool v) { return v ? "true" : "false"; }
static string __cgJson(int v) { return to_string(v); }
static string __cgJson(long long v) { return to_string(v); }
static string __cgJson(double v) {
    if (v == (long long)v) return to_string((long long)v);
    ostringstream o; o << v; return o.str();
}
static string __cgJson(char v) { return string("\"") + __cgEsc(string(1, v)) + "\""; }
static string __cgJson(const string& v) { return "\"" + __cgEsc(v) + "\""; }
template <typename T>
static string __cgJson(const vector<T>& v) {
    string o = "[";
    for (size_t i = 0; i < v.size(); i++) { if (i) o += ","; o += __cgJson(v[i]); }
    return o + "]";
}
`
