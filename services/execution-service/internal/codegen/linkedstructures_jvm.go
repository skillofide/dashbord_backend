package codegen

import "strings"

// javaLinkedClasses are appended after the submission, not before it.
//
// Java requires imports to precede every type declaration, so prepending the
// node classes would push a submission's own `import java.util.*;` after a
// class and stop the file compiling. Extra package-private top-level classes
// are legal anywhere in the file, so the end is the safe place.
const javaLinkedClasses = `

// Definition for a binary tree node.
class TreeNode {
    int val;
    TreeNode left;
    TreeNode right;
    TreeNode() {}
    TreeNode(int val) { this.val = val; }
    TreeNode(int val, TreeNode left, TreeNode right) { this.val = val; this.left = left; this.right = right; }
}

// Definition for singly-linked list node.
class ListNode {
    int val;
    ListNode next;
    ListNode() {}
    ListNode(int val) { this.val = val; }
    ListNode(int val, ListNode next) { this.val = val; this.next = next; }
}
`

// javaLinkedHelpers build and serialise the structures. They live inside
// Solution alongside the rest of the generated driver.
const javaLinkedHelpers = `
    private static java.util.List<String> __cgElems(String s) {
        String b = __cgTrim(s);
        if (b.isEmpty()) return new java.util.ArrayList<>();
        return __cgSplit(b);
    }
    private static TreeNode __cgBuildTree(String s) {
        java.util.List<String> a = __cgElems(s);
        if (a.isEmpty() || __cgIsNull(a.get(0))) return null;
        TreeNode root = new TreeNode(Integer.parseInt(a.get(0).trim()));
        java.util.Deque<TreeNode> q = new java.util.ArrayDeque<>();
        q.add(root);
        int i = 1;
        while (!q.isEmpty() && i < a.size()) {
            TreeNode node = q.poll();
            if (i < a.size()) {
                String v = a.get(i++);
                if (!__cgIsNull(v)) { node.left = new TreeNode(Integer.parseInt(v.trim())); q.add(node.left); }
            }
            if (i < a.size()) {
                String v = a.get(i++);
                if (!__cgIsNull(v)) { node.right = new TreeNode(Integer.parseInt(v.trim())); q.add(node.right); }
            }
        }
        return root;
    }
    private static boolean __cgIsNull(String s) { return s.trim().equals("null"); }
    private static String __cgDumpTree(TreeNode root) {
        if (root == null) return "[]";
        java.util.List<String> out = new java.util.ArrayList<>();
        // LinkedList, not ArrayDeque: the queue deliberately carries nulls to
        // mark absent children, and ArrayDeque throws on a null element.
        java.util.LinkedList<TreeNode> q = new java.util.LinkedList<>();
        q.add(root);
        while (!q.isEmpty()) {
            TreeNode node = q.poll();
            if (node == null) { out.add("null"); continue; }
            out.add(String.valueOf(node.val));
            q.add(node.left);
            q.add(node.right);
        }
        // Trailing nulls carry no information.
        while (!out.isEmpty() && out.get(out.size() - 1).equals("null")) out.remove(out.size() - 1);
        return "[" + String.join(",", out) + "]";
    }
    private static ListNode __cgBuildList(String s) {
        java.util.List<String> a = __cgElems(s);
        ListNode head = null, tail = null;
        for (String v : a) {
            ListNode node = new ListNode(Integer.parseInt(v.trim()));
            if (tail == null) { head = node; tail = node; }
            else { tail.next = node; tail = node; }
        }
        return head;
    }
    private static String __cgDumpList(ListNode head) {
        java.util.List<String> out = new java.util.ArrayList<>();
        int seen = 0;
        while (head != null) {
            out.add(String.valueOf(head.val));
            head = head.next;
            // A cycle would otherwise run to the time limit, which reads as
            // "too slow" rather than "wrong".
            if (++seen > 100000) break;
        }
        return "[" + String.join(",", out) + "]";
    }
`

// javaLinkedType spells a linked type in Java.
func javaLinkedType(t Type) (string, bool) {
	switch t {
	case TTreeNode:
		return "TreeNode", true
	case TListNode:
		return "ListNode", true
	}
	return "", false
}

func javaBuildExpr(t Type, src string) string {
	switch t {
	case TTreeNode:
		return "__cgBuildTree(" + src + ")"
	case TListNode:
		return "__cgBuildList(" + src + ")"
	}
	return ""
}

func javaDumpExpr(t Type, src string) string {
	switch t {
	case TTreeNode:
		return "__cgDumpTree(" + src + ")"
	case TListNode:
		return "__cgDumpList(" + src + ")"
	}
	return ""
}

// appendJavaLinkedClasses adds the node definitions after the submission.
func appendJavaLinkedClasses(s Signature, code string) string {
	if !needsLinkedPrelude(s) {
		return code
	}
	return strings.TrimRight(code, "\n") + "\n" + javaLinkedClasses
}
