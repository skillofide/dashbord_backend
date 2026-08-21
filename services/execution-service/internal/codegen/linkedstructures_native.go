package codegen

import "strings"

// ─── C++ ─────────────────────────────────────────────────────────────────────

// cppLinkedTypes are prepended, before the submission, because C++ requires a
// type to be declared before it is used and the submission refers to them.
// That is safe here in a way it is not in Java or Go: the wrapper already
// prepends its includes, and C++ has no rule about declarations preceding
// includes.
const cppLinkedTypes = `
// Definition for a binary tree node.
struct TreeNode {
    int val;
    TreeNode *left;
    TreeNode *right;
    TreeNode() : val(0), left(nullptr), right(nullptr) {}
    TreeNode(int x) : val(x), left(nullptr), right(nullptr) {}
    TreeNode(int x, TreeNode *l, TreeNode *r) : val(x), left(l), right(r) {}
};

// Definition for singly-linked list node.
struct ListNode {
    int val;
    ListNode *next;
    ListNode() : val(0), next(nullptr) {}
    ListNode(int x) : val(x), next(nullptr) {}
    ListNode(int x, ListNode *n) : val(x), next(n) {}
};
`

// cppLinkedHelpers build and serialise the structures. Declared after the
// submission with the rest of the driver helpers.
const cppLinkedHelpers = `
static TreeNode* __cgBuildTree(const string& s) {
    string b = __cgInner(s);
    if (b.empty()) return nullptr;
    vector<string> a = __cgSplit(b);
    if (a.empty() || __cgTrimWs(a[0]) == "null") return nullptr;
    TreeNode* root = new TreeNode(stoi(__cgTrimWs(a[0])));
    // deque of raw pointers; nullptr is a legal element here, unlike Java's
    // ArrayDeque, so absent children need no special casing on the way in.
    deque<TreeNode*> q;
    q.push_back(root);
    size_t i = 1;
    while (!q.empty() && i < a.size()) {
        TreeNode* node = q.front(); q.pop_front();
        if (i < a.size()) {
            string v = __cgTrimWs(a[i++]);
            if (v != "null") { node->left = new TreeNode(stoi(v)); q.push_back(node->left); }
        }
        if (i < a.size()) {
            string v = __cgTrimWs(a[i++]);
            if (v != "null") { node->right = new TreeNode(stoi(v)); q.push_back(node->right); }
        }
    }
    return root;
}
static string __cgDumpTree(TreeNode* root) {
    if (!root) return "[]";
    vector<string> out;
    deque<TreeNode*> q;
    q.push_back(root);
    while (!q.empty()) {
        TreeNode* node = q.front(); q.pop_front();
        if (!node) { out.push_back("null"); continue; }
        out.push_back(to_string(node->val));
        q.push_back(node->left);
        q.push_back(node->right);
    }
    while (!out.empty() && out.back() == "null") out.pop_back();
    string o = "[";
    for (size_t i = 0; i < out.size(); i++) { if (i) o += ","; o += out[i]; }
    return o + "]";
}
static ListNode* __cgBuildList(const string& s) {
    string b = __cgInner(s);
    ListNode *head = nullptr, *tail = nullptr;
    if (b.empty()) return nullptr;
    for (auto& p : __cgSplit(b)) {
        ListNode* node = new ListNode(stoi(__cgTrimWs(p)));
        if (!tail) { head = node; tail = node; }
        else { tail->next = node; tail = node; }
    }
    return head;
}
static string __cgDumpList(ListNode* head) {
    vector<string> out;
    int seen = 0;
    while (head) {
        out.push_back(to_string(head->val));
        head = head->next;
        // A cycle would otherwise run to the time limit, which reads as
        // "too slow" rather than "wrong".
        if (++seen > 100000) break;
    }
    string o = "[";
    for (size_t i = 0; i < out.size(); i++) { if (i) o += ","; o += out[i]; }
    return o + "]";
}
`

func cppLinkedType(t Type) (string, bool) {
	switch t {
	case TTreeNode:
		return "TreeNode*", true
	case TListNode:
		return "ListNode*", true
	}
	return "", false
}

func cppBuildExpr(t Type, src string) string {
	switch t {
	case TTreeNode:
		return "__cgBuildTree(" + src + ")"
	case TListNode:
		return "__cgBuildList(" + src + ")"
	}
	return ""
}

func cppDumpExpr(t Type, src string) string {
	switch t {
	case TTreeNode:
		return "__cgDumpTree(" + src + ")"
	case TListNode:
		return "__cgDumpList(" + src + ")"
	}
	return ""
}

// ─── Go ──────────────────────────────────────────────────────────────────────

// goLinkedTypes are appended after the submission. Go allows package-level
// declarations in any order, and appending keeps them clear of the injected
// import block — putting a declaration between those imports and the
// submission's own import line is what previously made every Go submission
// with an import fail to compile.
const goLinkedTypes = `
// Definition for a binary tree node.
type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

// Definition for singly-linked list node.
type ListNode struct {
	Val  int
	Next *ListNode
}

func __cgBuildTree(raw []byte) *TreeNode {
	var vals []*int
	if err := __cgJSON.Unmarshal(raw, &vals); err != nil || len(vals) == 0 || vals[0] == nil {
		return nil
	}
	root := &TreeNode{Val: *vals[0]}
	queue := []*TreeNode{root}
	i := 1
	for len(queue) > 0 && i < len(vals) {
		node := queue[0]
		queue = queue[1:]
		if i < len(vals) {
			if v := vals[i]; v != nil {
				node.Left = &TreeNode{Val: *v}
				queue = append(queue, node.Left)
			}
			i++
		}
		if i < len(vals) {
			if v := vals[i]; v != nil {
				node.Right = &TreeNode{Val: *v}
				queue = append(queue, node.Right)
			}
			i++
		}
	}
	return root
}

func __cgDumpTree(root *TreeNode) []byte {
	if root == nil {
		return []byte("[]")
	}
	var out []*int
	queue := []*TreeNode{root}
	for len(queue) > 0 {
		node := queue[0]
		queue = queue[1:]
		if node == nil {
			out = append(out, nil)
			continue
		}
		v := node.Val
		out = append(out, &v)
		queue = append(queue, node.Left, node.Right)
	}
	for len(out) > 0 && out[len(out)-1] == nil {
		out = out[:len(out)-1]
	}
	b, _ := __cgJSON.Marshal(out)
	return b
}

func __cgBuildList(raw []byte) *ListNode {
	var vals []int
	if err := __cgJSON.Unmarshal(raw, &vals); err != nil {
		return nil
	}
	var head, tail *ListNode
	for _, v := range vals {
		node := &ListNode{Val: v}
		if tail == nil {
			head, tail = node, node
		} else {
			tail.Next = node
			tail = node
		}
	}
	return head
}

func __cgDumpList(head *ListNode) []byte {
	out := []int{}
	seen := 0
	for head != nil {
		out = append(out, head.Val)
		head = head.Next
		// A cycle would otherwise run to the time limit, which reads as
		// "too slow" rather than "wrong".
		seen++
		if seen > 100000 {
			break
		}
	}
	b, _ := __cgJSON.Marshal(out)
	return b
}
`

func goLinkedType(t Type) (string, bool) {
	switch t {
	case TTreeNode:
		return "*TreeNode", true
	case TListNode:
		return "*ListNode", true
	}
	return "", false
}

// appendGoLinkedTypes adds the node definitions and helpers after the driver.
func appendGoLinkedTypes(s Signature, code string) string {
	if !needsLinkedPrelude(s) {
		return code
	}
	return strings.TrimRight(code, "\n") + "\n" + goLinkedTypes
}

// prependCppLinkedTypes puts the struct definitions above the submission.
func prependCppLinkedTypes(s Signature, code string) string {
	if !needsLinkedPrelude(s) {
		return code
	}
	return cppLinkedTypes + "\n" + strings.TrimLeft(code, "\n")
}
