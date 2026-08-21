package codegen

import "strings"

// jsLinkedPrelude defines the node types and the conversions between a node
// structure and its serialised array form.
//
// It is prepended to the submission, not appended, so a learner can refer to
// TreeNode and ListNode in their own code the way the problem statement assumes
// they exist. The __cg helpers are only called by the generated driver.
const jsLinkedPrelude = `// Definition for a binary tree node.
function TreeNode(val, left, right) {
    this.val = (val === undefined ? 0 : val);
    this.left = (left === undefined ? null : left);
    this.right = (right === undefined ? null : right);
}

// Definition for singly-linked list node.
function ListNode(val, next) {
    this.val = (val === undefined ? 0 : val);
    this.next = (next === undefined ? null : next);
}

// ─── generated: linked-structure serialisation ───
function __cgBuildTree(a) {
    if (!Array.isArray(a) || a.length === 0 || a[0] === null) return null;
    const root = new TreeNode(a[0]);
    const q = [root];
    let i = 1;
    while (q.length > 0 && i < a.length) {
        const node = q.shift();
        if (i < a.length) {
            const v = a[i++];
            if (v !== null) { node.left = new TreeNode(v); q.push(node.left); }
        }
        if (i < a.length) {
            const v = a[i++];
            if (v !== null) { node.right = new TreeNode(v); q.push(node.right); }
        }
    }
    return root;
}
function __cgDumpTree(root) {
    if (!root) return [];
    const out = [];
    const q = [root];
    while (q.length > 0) {
        const node = q.shift();
        if (node === null) { out.push(null); continue; }
        out.push(node.val);
        q.push(node.left === undefined ? null : node.left);
        q.push(node.right === undefined ? null : node.right);
    }
    // Trailing nulls carry no information; LeetCode's own serialisation omits them.
    while (out.length > 0 && out[out.length - 1] === null) out.pop();
    return out;
}
function __cgBuildList(a) {
    if (!Array.isArray(a)) return null;
    let head = null, tail = null;
    for (const v of a) {
        const node = new ListNode(v);
        if (tail === null) { head = node; tail = node; }
        else { tail.next = node; tail = node; }
    }
    return head;
}
function __cgDumpList(head) {
    const out = [];
    let seen = 0;
    while (head) {
        out.push(head.val);
        head = head.next;
        // A cycle would otherwise hang the container until the time limit,
        // which reads to the learner as "too slow" rather than "wrong".
        if (++seen > 100000) break;
    }
    return out;
}

`

// needsLinkedPrelude reports whether any part of the signature is a linked
// structure.
func needsLinkedPrelude(s Signature) bool {
	if s.ReturnType.Linked() {
		return true
	}
	for _, p := range s.Params {
		if p.Type.Linked() {
			return true
		}
	}
	return false
}

// jsBuildExpr wraps a decoded JSON value into the structure the entry point
// expects.
func jsBuildExpr(t Type, src string) string {
	switch t {
	case TTreeNode:
		return "__cgBuildTree(" + src + ")"
	case TListNode:
		return "__cgBuildList(" + src + ")"
	}
	return src
}

// jsDumpExpr converts a returned structure back to its serialised form.
func jsDumpExpr(t Type, src string) string {
	switch t {
	case TTreeNode:
		return "__cgDumpTree(" + src + ")"
	case TListNode:
		return "__cgDumpList(" + src + ")"
	}
	return src
}

func jsLinkedDocType(t Type) (string, bool) {
	switch t {
	case TTreeNode:
		return "TreeNode", true
	case TListNode:
		return "ListNode", true
	}
	return "", false
}

// prependLinkedPrelude puts the node definitions above the submission.
func prependLinkedPrelude(s Signature, code string) string {
	if !needsLinkedPrelude(s) {
		return code
	}
	return jsLinkedPrelude + strings.TrimLeft(code, "\n")
}
