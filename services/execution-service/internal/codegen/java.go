package codegen

import (
	"fmt"
	"strings"
)

// javaType spells a Type as Java source.
func javaType(t Type) string {
	if name, ok := javaLinkedType(t); ok {
		return name
	}
	switch t {
	case TInt:
		return "int"
	case TLong:
		return "long"
	case TDouble:
		return "double"
	case TBool:
		return "boolean"
	case TString:
		return "String"
	case TChar:
		return "char"
	case TIntArr:
		return "int[]"
	case TLongArr:
		return "long[]"
	case TDblArr:
		return "double[]"
	case TBoolArr:
		return "boolean[]"
	case TStrArr:
		return "String[]"
	case TInt2D:
		return "int[][]"
	case TStr2D:
		return "String[][]"
	case TVoid:
		return "void"
	}
	return "Object"
}

// StarterJava renders the Java skeleton.
func StarterJava(s Signature) string {
	parts := make([]string, len(s.Params))
	for i, p := range s.Params {
		parts[i] = fmt.Sprintf("%s %s", javaType(p.Type), p.Name)
	}

	var b strings.Builder
	b.WriteString("public class Solution {\n")
	fmt.Fprintf(&b, "    public %s %s(%s) {\n", javaType(s.ReturnType), s.EntryPoint, strings.Join(parts, ", "))
	b.WriteString("        // Write your code here\n")
	if s.ReturnType != TVoid {
		fmt.Fprintf(&b, "        %s\n", javaZeroReturn(s.ReturnType))
	}
	b.WriteString("    }\n}\n")
	return appendJavaLinkedClasses(s, b.String())
}

func javaZeroReturn(t Type) string {
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
		return "return null;"
	}
	return "return null;"
}

// javaParseExpr returns an expression producing a value of the given type from
// a String holding one JSON line.
func javaParseExpr(t Type, src string) string {
	if expr := javaBuildExpr(t, src); expr != "" {
		return expr
	}
	switch t {
	case TInt:
		return fmt.Sprintf("Integer.parseInt(%s.trim())", src)
	case TLong:
		return fmt.Sprintf("Long.parseLong(%s.trim())", src)
	case TDouble:
		return fmt.Sprintf("Double.parseDouble(%s.trim())", src)
	case TBool:
		return fmt.Sprintf("Boolean.parseBoolean(%s.trim())", src)
	case TChar:
		return fmt.Sprintf("__cgStr(%s).charAt(0)", src)
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
	return "null"
}

// WrapJava injects a generated main plus JSON helpers into the Solution class.
//
// The JDK ships no JSON parser, so rather than embedding a general one the
// driver parses each line against the declared type — which it knows exactly.
// The old reflection driver had to guess the type from a Method object and
// split arguments with a hand-rolled bracket counter on a single line, so any
// problem taking more than one line of input was unreachable.
func WrapJava(s Signature, code string) string {
	var m strings.Builder
	m.WriteString("\n    // ─── generated driver ───\n")
	m.WriteString("    public static void main(String[] __args) throws Exception {\n")
	m.WriteString("        java.io.BufferedReader __rd = new java.io.BufferedReader(new java.io.InputStreamReader(System.in));\n")
	m.WriteString("        Solution __sol = new Solution();\n")

	names := make([]string, len(s.Params))
	for i, p := range s.Params {
		fmt.Fprintf(&m, "        String __l%d = __rd.readLine(); if (__l%d == null) __l%d = \"\";\n", i, i, i)
		fmt.Fprintf(&m, "        %s %s = %s;\n", javaType(p.Type), p.Name, javaParseExpr(p.Type, fmt.Sprintf("__l%d", i)))
		names[i] = p.Name
	}

	call := fmt.Sprintf("__sol.%s(%s)", s.EntryPoint, strings.Join(names, ", "))
	if s.ReturnType == TVoid {
		fmt.Fprintf(&m, "        %s;\n", call)
		if len(s.Params) > 0 {
			fmt.Fprintf(&m, "        System.out.print(__cgJson(%s));\n", s.Params[0].Name)
		}
	} else {
		fmt.Fprintf(&m, "        %s __result = %s;\n", javaType(s.ReturnType), call)
		if dump := javaDumpExpr(s.ReturnType, "__result"); dump != "" {
			fmt.Fprintf(&m, "        System.out.print(%s);\n", dump)
		} else {
			m.WriteString("        System.out.print(__cgJson(__result));\n")
		}
	}
	m.WriteString("    }\n")
	m.WriteString(javaHelpers)
	if needsLinkedPrelude(s) {
		m.WriteString(javaLinkedHelpers)
	}

	// Inject before the class's final closing brace.
	last := strings.LastIndex(code, "}")
	if last == -1 {
		return appendJavaLinkedClasses(s, code+"\n"+m.String())
	}
	return appendJavaLinkedClasses(s, code[:last]+m.String()+"\n}\n")
}

// javaHelpers are the type-directed parse and serialise routines the driver
// calls. Names are __cg-prefixed so they cannot shadow anything a learner
// writes in the same class.
const javaHelpers = `
    private static String __cgTrim(String s) {
        s = s.trim();
        if (s.startsWith("[") && s.endsWith("]")) s = s.substring(1, s.length() - 1);
        return s.trim();
    }
    private static java.util.List<String> __cgSplit(String s) {
        java.util.List<String> out = new java.util.ArrayList<>();
        int depth = 0; boolean q = false; StringBuilder cur = new StringBuilder();
        for (int i = 0; i < s.length(); i++) {
            char c = s.charAt(i);
            if (c == '"' && (i == 0 || s.charAt(i - 1) != '\\')) q = !q;
            if (!q && (c == '[' || c == '{')) depth++;
            if (!q && (c == ']' || c == '}')) depth--;
            if (c == ',' && depth == 0 && !q) { out.add(cur.toString()); cur = new StringBuilder(); }
            else cur.append(c);
        }
        if (cur.length() > 0) out.add(cur.toString());
        return out;
    }
    private static String __cgStr(String s) {
        s = s.trim();
        if (s.startsWith("\"") && s.endsWith("\"") && s.length() >= 2) s = s.substring(1, s.length() - 1);
        return s.replace("\\\"", "\"").replace("\\n", "\n").replace("\\t", "\t").replace("\\\\", "\\");
    }
    private static int[] __cgIntArr(String s) {
        String b = __cgTrim(s);
        if (b.isEmpty()) return new int[0];
        java.util.List<String> parts = __cgSplit(b);
        int[] a = new int[parts.size()];
        for (int i = 0; i < a.length; i++) a[i] = Integer.parseInt(parts.get(i).trim());
        return a;
    }
    private static long[] __cgLongArr(String s) {
        String b = __cgTrim(s);
        if (b.isEmpty()) return new long[0];
        java.util.List<String> parts = __cgSplit(b);
        long[] a = new long[parts.size()];
        for (int i = 0; i < a.length; i++) a[i] = Long.parseLong(parts.get(i).trim());
        return a;
    }
    private static double[] __cgDblArr(String s) {
        String b = __cgTrim(s);
        if (b.isEmpty()) return new double[0];
        java.util.List<String> parts = __cgSplit(b);
        double[] a = new double[parts.size()];
        for (int i = 0; i < a.length; i++) a[i] = Double.parseDouble(parts.get(i).trim());
        return a;
    }
    private static boolean[] __cgBoolArr(String s) {
        String b = __cgTrim(s);
        if (b.isEmpty()) return new boolean[0];
        java.util.List<String> parts = __cgSplit(b);
        boolean[] a = new boolean[parts.size()];
        for (int i = 0; i < a.length; i++) a[i] = Boolean.parseBoolean(parts.get(i).trim());
        return a;
    }
    private static String[] __cgStrArr(String s) {
        String b = __cgTrim(s);
        if (b.isEmpty()) return new String[0];
        java.util.List<String> parts = __cgSplit(b);
        String[] a = new String[parts.size()];
        for (int i = 0; i < a.length; i++) a[i] = __cgStr(parts.get(i));
        return a;
    }
    private static int[][] __cgInt2D(String s) {
        String b = __cgTrim(s);
        if (b.isEmpty()) return new int[0][];
        java.util.List<String> rows = __cgSplit(b);
        int[][] a = new int[rows.size()][];
        for (int i = 0; i < a.length; i++) a[i] = __cgIntArr(rows.get(i));
        return a;
    }
    private static String[][] __cgStr2D(String s) {
        String b = __cgTrim(s);
        if (b.isEmpty()) return new String[0][];
        java.util.List<String> rows = __cgSplit(b);
        String[][] a = new String[rows.size()][];
        for (int i = 0; i < a.length; i++) a[i] = __cgStrArr(rows.get(i));
        return a;
    }
    private static String __cgEsc(String s) {
        return s.replace("\\", "\\\\").replace("\"", "\\\"").replace("\n", "\\n").replace("\t", "\\t");
    }
    private static String __cgJson(Object v) {
        if (v == null) return "null";
        if (v instanceof String) return "\"" + __cgEsc((String) v) + "\"";
        if (v instanceof Character) return "\"" + __cgEsc(String.valueOf(v)) + "\"";
        if (v instanceof Boolean || v instanceof Integer || v instanceof Long) return String.valueOf(v);
        if (v instanceof Double || v instanceof Float) {
            double d = ((Number) v).doubleValue();
            if (d == Math.rint(d) && !Double.isInfinite(d)) return String.valueOf((long) d);
            return String.valueOf(d);
        }
        if (v instanceof int[]) {
            StringBuilder sb = new StringBuilder("["); int[] a = (int[]) v;
            for (int i = 0; i < a.length; i++) { if (i > 0) sb.append(","); sb.append(a[i]); }
            return sb.append("]").toString();
        }
        if (v instanceof long[]) {
            StringBuilder sb = new StringBuilder("["); long[] a = (long[]) v;
            for (int i = 0; i < a.length; i++) { if (i > 0) sb.append(","); sb.append(a[i]); }
            return sb.append("]").toString();
        }
        if (v instanceof double[]) {
            StringBuilder sb = new StringBuilder("["); double[] a = (double[]) v;
            for (int i = 0; i < a.length; i++) { if (i > 0) sb.append(","); sb.append(__cgJson(a[i])); }
            return sb.append("]").toString();
        }
        if (v instanceof boolean[]) {
            StringBuilder sb = new StringBuilder("["); boolean[] a = (boolean[]) v;
            for (int i = 0; i < a.length; i++) { if (i > 0) sb.append(","); sb.append(a[i]); }
            return sb.append("]").toString();
        }
        if (v instanceof Object[]) {
            StringBuilder sb = new StringBuilder("["); Object[] a = (Object[]) v;
            for (int i = 0; i < a.length; i++) { if (i > 0) sb.append(","); sb.append(__cgJson(a[i])); }
            return sb.append("]").toString();
        }
        return String.valueOf(v);
    }
`
