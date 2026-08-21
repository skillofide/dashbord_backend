package codegen

import (
	"os"
	"testing"
)

const goTwoSumSolution = `package main

func twoSum(nums []int, target int) []int {
	seen := make(map[int]int)
	for i, n := range nums {
		if j, ok := seen[target-n]; ok {
			return []int{j, i}
		}
		seen[n] = i
	}
	return nil
}
`

// A learner who imports a package themselves must not collide with the imports
// the driver injects.
const goTwoSumWithUserImport = `package main

import "os"

func twoSum(nums []int, target int) []int {
	_ = os.Getenv
	seen := make(map[int]int)
	for i, n := range nums {
		if j, ok := seen[target-n]; ok {
			return []int{j, i}
		}
		seen[n] = i
	}
	return nil
}
`

const javaTwoSumSolution = `import java.util.*;

public class Solution {
    public int[] twoSum(int[] nums, int target) {
        Map<Integer, Integer> seen = new HashMap<>();
        for (int i = 0; i < nums.length; i++) {
            if (seen.containsKey(target - nums[i])) return new int[]{seen.get(target - nums[i]), i};
            seen.put(nums[i], i);
        }
        return new int[0];
    }
}
`

const cppTwoSumSolution = `#include <bits/stdc++.h>
using namespace std;

class Solution {
public:
    vector<int> twoSum(vector<int>& nums, int target) {
        unordered_map<int,int> seen;
        for (int i = 0; i < (int)nums.size(); i++) {
            auto it = seen.find(target - nums[i]);
            if (it != seen.end()) return {it->second, i};
            seen[nums[i]] = i;
        }
        return {};
    }
};
`

// TestEmitSample writes the generated artefacts for Two Sum to a directory so
// they can be executed inside the real runner images. Unit tests confirm the
// generated text; only running it in the container confirms it works. Guarded
// by an env var so an ordinary `go test` run does not write files.
func TestEmitSample(t *testing.T) {
	dir := os.Getenv("CODEGEN_EMIT_DIR")
	if dir == "" {
		t.Skip("set CODEGEN_EMIT_DIR to emit sample artefacts")
	}
	sig := twoSumSig()

	files := map[string]string{
		"two-sum.starter.js":     StarterJavaScript(sig),
		"two-sum.driver.js":      DriverJavaScript(sig),
		"two-sum.starter.py":     StarterPython(sig),
		"two-sum.driver.py":      DriverPython(sig),
		"two-sum.wrapped.go":     WrapGo(sig, goTwoSumSolution),
		"two-sum.userimport.go":  WrapGo(sig, goTwoSumWithUserImport),
		"two-sum.starterwrap.go": WrapGo(sig, StarterGo(sig)),
		"Solution.java":          WrapJava(sig, javaTwoSumSolution),
		"SolutionStarter.java":   WrapJava(sig, StarterJava(sig)),
		"two-sum.wrapped.cpp":    WrapCpp(sig, cppTwoSumSolution),
		"two-sum.starter.cpp":    WrapCpp(sig, StarterCpp(sig)),
	}
	for name, body := range files {
		if err := os.WriteFile(dir+"/"+name, []byte(body), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
}

// TestEmitLinkedSamples writes wrapped linked-structure programs so they can be
// run inside the real runner images. Unit tests confirm the generated text;
// only a compiler confirms it works.
func TestEmitLinkedSamples(t *testing.T) {
	dir := os.Getenv("CODEGEN_EMIT_DIR")
	if dir == "" {
		t.Skip("set CODEGEN_EMIT_DIR to emit sample artefacts")
	}
	tree := Signature{
		EntryPoint: "isValidBST",
		Params:     []Param{{Name: "root", Type: TTreeNode}},
		ReturnType: TBool,
	}
	javaTree := `import java.util.*;

public class Solution {
    public boolean isValidBST(TreeNode root) {
        return check(root, null, null);
    }
    private boolean check(TreeNode n, Integer lo, Integer hi) {
        if (n == null) return true;
        if (lo != null && n.val <= lo) return false;
        if (hi != null && n.val >= hi) return false;
        return check(n.left, lo, n.val) && check(n.right, n.val, hi);
    }
}
`
	list := Signature{
		EntryPoint: "swapPairs",
		Params:     []Param{{Name: "head", Type: TListNode}},
		ReturnType: TListNode,
	}
	javaList := `public class Solution {
    public ListNode swapPairs(ListNode head) {
        ListNode dummy = new ListNode(0);
        dummy.next = head;
        ListNode prev = dummy;
        while (prev.next != null && prev.next.next != null) {
            ListNode a = prev.next, b = a.next;
            a.next = b.next;
            b.next = a;
            prev.next = b;
            prev = a;
        }
        return dummy.next;
    }
}
`
	cppList := `class Solution {
public:
    ListNode* swapPairs(ListNode* head) {
        ListNode dummy(0);
        dummy.next = head;
        ListNode* prev = &dummy;
        while (prev->next && prev->next->next) {
            ListNode *a = prev->next, *b = a->next;
            a->next = b->next;
            b->next = a;
            prev->next = b;
            prev = a;
        }
        return dummy.next;
    }
};
`
	cppTree := `class Solution {
public:
    bool isValidBST(TreeNode* root) { return check(root, nullptr, nullptr); }
private:
    bool check(TreeNode* n, TreeNode* lo, TreeNode* hi) {
        if (!n) return true;
        if (lo && n->val <= lo->val) return false;
        if (hi && n->val >= hi->val) return false;
        return check(n->left, lo, n) && check(n->right, n, hi);
    }
};
`
	goList := `package main

func swapPairs(head *ListNode) *ListNode {
	dummy := &ListNode{}
	dummy.Next = head
	prev := dummy
	for prev.Next != nil && prev.Next.Next != nil {
		a, b := prev.Next, prev.Next.Next
		a.Next = b.Next
		b.Next = a
		prev.Next = b
		prev = a
	}
	return dummy.Next
}
`
	goTree := `package main

func isValidBST(root *TreeNode) bool {
	var check func(n *TreeNode, lo, hi *int) bool
	check = func(n *TreeNode, lo, hi *int) bool {
		if n == nil {
			return true
		}
		if lo != nil && n.Val <= *lo {
			return false
		}
		if hi != nil && n.Val >= *hi {
			return false
		}
		v := n.Val
		return check(n.Left, lo, &v) && check(n.Right, &v, hi)
	}
	return check(root, nil, nil)
}
`

	files := map[string]string{
		"TreeSolution.java": WrapJava(tree, javaTree),
		"ListSolution.java": WrapJava(list, javaList),
		"tree.cpp":          WrapCpp(tree, cppTree),
		"list.cpp":          WrapCpp(list, cppList),
		"tree_linked.go":    WrapGo(tree, goTree),
		"list_linked.go":    WrapGo(list, goList),
	}
	for name, body := range files {
		if err := os.WriteFile(dir+"/"+name, []byte(body), 0o644); err != nil {
			t.Fatalf("write %s: %v", name, err)
		}
	}
}
