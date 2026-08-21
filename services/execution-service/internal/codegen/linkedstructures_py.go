package codegen

import "strings"

// pyLinkedPrelude mirrors the JavaScript one: node definitions the learner can
// refer to, plus the conversions the generated driver uses.
const pyLinkedPrelude = `# Definition for a binary tree node.
class TreeNode:
    def __init__(self, val=0, left=None, right=None):
        self.val = val
        self.left = left
        self.right = right


# Definition for singly-linked list node.
class ListNode:
    def __init__(self, val=0, next=None):
        self.val = val
        self.next = next


# ─── generated: linked-structure serialisation ───
def __cg_build_tree(a):
    if not a or a[0] is None:
        return None
    root = TreeNode(a[0])
    q = [root]
    i = 1
    while q and i < len(a):
        node = q.pop(0)
        if i < len(a):
            v = a[i]
            i += 1
            if v is not None:
                node.left = TreeNode(v)
                q.append(node.left)
        if i < len(a):
            v = a[i]
            i += 1
            if v is not None:
                node.right = TreeNode(v)
                q.append(node.right)
    return root


def __cg_dump_tree(root):
    if root is None:
        return []
    out = []
    q = [root]
    while q:
        node = q.pop(0)
        if node is None:
            out.append(None)
            continue
        out.append(node.val)
        q.append(node.left)
        q.append(node.right)
    # Trailing nulls carry no information.
    while out and out[-1] is None:
        out.pop()
    return out


def __cg_build_list(a):
    head = tail = None
    for v in (a or []):
        node = ListNode(v)
        if tail is None:
            head = tail = node
        else:
            tail.next = node
            tail = node
    return head


def __cg_dump_list(head):
    out = []
    seen = 0
    while head is not None:
        out.append(head.val)
        head = head.next
        # A cycle would otherwise run until the time limit, which reads as
        # "too slow" rather than "wrong".
        seen += 1
        if seen > 100000:
            break
    return out


`

func pyBuildExpr(t Type, src string) string {
	switch t {
	case TTreeNode:
		return "__cg_build_tree(" + src + ")"
	case TListNode:
		return "__cg_build_list(" + src + ")"
	}
	return src
}

func pyDumpExpr(t Type, src string) string {
	switch t {
	case TTreeNode:
		return "__cg_dump_tree(" + src + ")"
	case TListNode:
		return "__cg_dump_list(" + src + ")"
	}
	return src
}

func pyLinkedType(t Type) (string, bool) {
	switch t {
	case TTreeNode:
		return "TreeNode", true
	case TListNode:
		return "ListNode", true
	}
	return "", false
}

func prependPyLinkedPrelude(s Signature, code string) string {
	if !needsLinkedPrelude(s) {
		return code
	}
	return pyLinkedPrelude + strings.TrimLeft(code, "\n")
}
