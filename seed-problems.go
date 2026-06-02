package main

import (
	"context"
	"fmt"
	"os"

	"github.com/jackc/pgx/v5"
)

type PracticeSet struct {
	ID         string
	Title      string
	Level      string
	LevelColor string
	BgColor    string
}

type Example struct {
	Input       string
	Output      string
	Explanation string
}

type Hint struct {
	Title string
	Body  string
}

type TestCase struct {
	Input    string
	Expected string
	IsHidden bool
}

type Problem struct {
	ID           string
	Slug         string
	Title        string
	Difficulty   string
	Topic        string
	XP           int
	Statement    string
	SetID        string
	Tags         []string
	Examples     []Example
	Hints        []Hint
	JavascriptSC string
	PythonSC     string
	JavaSC       string
	CppSC        string
	GoSC         string
	TestCases    []TestCase
}

func main() {
	dsn := os.Getenv("POSTGRES_DSN")
	if dsn == "" {
		dsn = "postgres://skillofide:password@localhost:5432/skillofide?sslmode=disable"
	}

	ctx := context.Background()
	conn, err := pgx.Connect(ctx, dsn)
	if err != nil {
		fmt.Printf("Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer conn.Close(ctx)

	fmt.Println("Connected to PostgreSQL database for seeding...")

	// Clear existing tables
	tables := []string{
		"starter_codes",
		"test_cases",
		"examples",
		"hints",
		"problem_constraints",
		"problem_tags",
		"problem_user_status",
		"problems",
		"practice_sets",
	}

	for _, table := range tables {
		_, err := conn.Exec(ctx, fmt.Sprintf("TRUNCATE TABLE %s CASCADE", table))
		if err != nil {
			fmt.Printf("Failed to truncate table %s: %v\n", table, err)
		}
	}
	fmt.Println("Cleared existing problem tables.")

	// 1. Insert practice sets
	sets := []PracticeSet{
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4ca001",
			Title:      "Foundational Basics",
			Level:      "Beginner",
			LevelColor: "#22c55e",
			BgColor:    "#1e1b4b",
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4ca002",
			Title:      "Path to Proficiency",
			Level:      "Intermediate",
			LevelColor: "#3b82f6",
			BgColor:    "#1e1b4b",
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Title:      "Masters of Algorithms",
			Level:      "Advanced",
			LevelColor: "#9b5cf6",
			BgColor:    "#1e1b4b",
		},
	}

	for _, s := range sets {
		_, err := conn.Exec(ctx, `
			INSERT INTO practice_sets (id, title, level, level_color, bg_color, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, now(), now())
		`, s.ID, s.Title, s.Level, s.LevelColor, s.BgColor)
		if err != nil {
			fmt.Printf("Failed to insert practice set %s: %v\n", s.Title, err)
			os.Exit(1)
		}
	}
	fmt.Println("Seeded practice sets.")

	// 2. Define problems
	problems := []Problem{
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4c0001",
			Slug:       "op1",
			Title:      "Arithmetic Operators Basics",
			Difficulty: "Easy",
			Topic:      "Operators",
			XP:         50,
			Statement:  "Write a program that takes two integer inputs `a` and `b` and returns their sum, difference, product, and quotient (integer division) in an array format: `[sum, difference, product, quotient]`.\n\nMake sure to handle standard arithmetic rules.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca001",
			Tags:       []string{"Basic Math", "Operators"},
			Examples: []Example{
				{
					Input:       "a = 10\nb = 2",
					Output:      "[12, 8, 20, 5]",
					Explanation: "10+2 = 12, 10-2 = 8, 10*2 = 20, 10/2 = 5.",
				},
			},
			Hints: []Hint{
				{
					Title: "Division operator",
					Body:  "Use integer division. In JavaScript, use Math.floor(a / b). In Python, use a // b.",
				},
			},
			JavascriptSC: `/**
 * @param {number} a
 * @param {number} b
 * @return {number[]}
 */
function arithmeticOperations(a, b) {
    // Write your code here
    return [
        a + b,
        a - b,
        a * b,
        Math.floor(a / b)
    ];
}`,
			PythonSC: `def arithmeticOperations(a: int, b: int) -> list[int]:
    # Write your code here
    return [
        a + b,
        a - b,
        a * b,
        a // b
    ]`,
			JavaSC: `public class Solution {
    public int[] arithmeticOperations(int a, int b) {
        // Write your code here
        return new int[] {
            a + b,
            a - b,
            a * b,
            a / b
        };
    }
}`,
			CppSC: `#include <vector>
using namespace std;

class Solution {
public:
    vector<int> arithmeticOperations(int a, int b) {
        // Write your code here
        return { a + b, a - b, a * b, a / b };
    }
};`,
			GoSC: `package main

import "fmt"

func arithmeticOperations(a, b int) []int {
    return []int{
        a + b,
        a - b,
        a * b,
        a / b,
    }
}`,
			TestCases: []TestCase{
				{Input: "10\n2", Expected: "[12,8,20,5]", IsHidden: false},
				{Input: "15\n3", Expected: "[18,12,45,5]", IsHidden: false},
				{Input: "-4\n2", Expected: "[-2,-6,-8,-2]", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4c0002",
			Slug:       "arr5",
			Title:      "Two Sum Problem",
			Difficulty: "Medium",
			Topic:      "Arrays",
			XP:         100,
			Statement:  "Given an array of integers `nums` and an integer `target`, return *indices of the two numbers such that they add up to `target`*.\n\nYou may assume that each input would have ***exactly* one solution**, and you may not use the *same* element twice.\n\nYou can return the answer in any order.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca002",
			Tags:       []string{"Arrays", "Hash Table", "Search"},
			Examples: []Example{
				{
					Input:       "nums = [2, 7, 11, 15], target = 9",
					Output:      "[0, 1]",
					Explanation: "Because nums[0] + nums[1] == 9, we return [0, 1].",
				},
				{
					Input:       "nums = [3, 2, 4], target = 6",
					Output:      "[1, 2]",
					Explanation: "Because nums[1] + nums[2] == 6, we return [1, 2].",
				},
			},
			Hints: []Hint{
				{
					Title: "Hash Map Approach",
					Body:  "Use a hash map to store the index of each element. For each element, check if its complement (target - num) exists in the map.",
				},
			},
			JavascriptSC: `/**
 * @param {number[]} nums
 * @param {number} target
 * @return {number[]}
 */
function twoSum(nums, target) {
    // Write your code here
    const map = new Map();
    for (let i = 0; i < nums.length; i++) {
        const complement = target - nums[i];
        if (map.has(complement)) {
            return [map.get(complement), i];
        }
        map.set(nums[i], i);
    }
    return [];
}`,
			PythonSC: `def twoSum(nums: list[int], target: int) -> list[int]:
    # Write your code here
    seen = {}
    for i, num in enumerate(nums):
        complement = target - num
        if complement in seen:
            return [seen[complement], i]
        seen[num] = i
    return []`,
			JavaSC: `import java.util.*;

public class Solution {
    public int[] twoSum(int[] nums, int target) {
        // Write your code here
        Map<Integer, Integer> map = new HashMap<>();
        for (int i = 0; i < nums.length; i++) {
            int complement = target - nums[i];
            if (map.containsKey(complement)) {
                return new int[] { map.get(complement), i };
            }
            map.put(nums[i], i);
        }
        return new int[0];
    }
}`,
			CppSC: `#include <vector>
#include <unordered_map>
using namespace std;

class Solution {
public:
    vector<int> twoSum(vector<int>& nums, int target) {
        // Write your code here
        unordered_map<int, int> seen;
        for (int i = 0; i < nums.size(); ++i) {
            int complement = target - nums[i];
            if (seen.count(complement)) {
                return {seen[complement], i};
            }
            seen[nums[i]] = i;
        }
        return {};
    }
};`,
			GoSC: `package main

func twoSum(nums []int, target int) []int {
    seen := make(map[int]int)
    for i, num := range nums {
        complement := target - num
        if idx, found := seen[complement]; found {
            return []int{idx, i}
        }
        seen[num] = i
    }
    return nil
}`,
			TestCases: []TestCase{
				{Input: "[2,7,11,15]\n9", Expected: "[0,1]", IsHidden: false},
				{Input: "[3,2,4]\n6", Expected: "[1,2]", IsHidden: false},
				{Input: "[3,3]\n6", Expected: "[0,1]", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4c0003",
			Slug:       "cond1",
			Title:      "Basic If-Else Statement",
			Difficulty: "Easy",
			Topic:      "Conditionals",
			XP:         50,
			Statement:  "Given an integer `n`, write a function to check if the number is **Even** or **Odd**.\nReturn the string `\"Even\"` or `\"Odd\"` accordingly.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca001",
			Tags:       []string{"Conditionals", "Logic"},
			Examples: []Example{
				{
					Input:       "n = 4",
					Output:      "\"Even\"",
					Explanation: "4 is divisible by 2, so it is even.",
				},
				{
					Input:       "n = 7",
					Output:      "\"Odd\"",
					Explanation: "7 is not divisible by 2, so it is odd.",
				},
			},
			Hints: []Hint{
				{
					Title: "Modulo operator",
					Body:  "Use the modulo operator % to check the remainder when divided by 2.",
				},
			},
			JavascriptSC: `/**
 * @param {number} n
 * @return {string}
 */
function checkEvenOdd(n) {
    // Write your code here
    return n % 2 === 0 ? "Even" : "Odd";
}`,
			PythonSC: `def checkEvenOdd(n: int) -> str:
    # Write your code here
    return "Even" if n % 2 == 0 else "Odd"`,
			JavaSC: `public class Solution {
    public String checkEvenOdd(int n) {
        // Write your code here
        return n % 2 == 0 ? "Even" : "Odd";
    }
}`,
			CppSC: `#include <string>
using namespace std;

class Solution {
public:
    string checkEvenOdd(int n) {
        // Write your code here
        return n % 2 == 0 ? "Even" : "Odd";
    }
};`,
			GoSC: `package main

func checkEvenOdd(n int) string {
    if n%2 == 0 {
        return "Even"
    }
    return "Odd"
}`,
			TestCases: []TestCase{
				{Input: "4", Expected: "\"Even\"", IsHidden: false},
				{Input: "7", Expected: "\"Odd\"", IsHidden: false},
				{Input: "0", Expected: "\"Even\"", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4c0004",
			Slug:       "loop1",
			Title:      "Sum of First N Numbers (For)",
			Difficulty: "Easy",
			Topic:      "Loops",
			XP:         50,
			Statement:  "Given a positive integer `N`, calculate the sum of all natural numbers from `1` to `N` inclusive using a loop.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca001",
			Tags:       []string{"Loops", "Basic Math"},
			Examples: []Example{
				{
					Input:       "N = 5",
					Output:      "15",
					Explanation: "1 + 2 + 3 + 4 + 5 = 15.",
				},
			},
			Hints: []Hint{
				{
					Title: "Looping",
					Body:  "Use a for loop running from 1 to N, adding the loop counter to a running sum variable.",
				},
			},
			JavascriptSC: `/**
 * @param {number} n
 * @return {number}
 */
function sumOfN(n) {
    // Write your code here
    let sum = 0;
    for (let i = 1; i <= n; i++) {
        sum += i;
    }
    return sum;
}`,
			PythonSC: `def sumOfN(n: int) -> int:
    # Write your code here
    total = 0
    for i in range(1, n + 1):
        total += i
    return total`,
			JavaSC: `public class Solution {
    public int sumOfN(int n) {
        // Write your code here
        int sum = 0;
        for (int i = 1; i <= n; i++) {
            sum += i;
        }
        return sum;
    }
}`,
			CppSC: `class Solution {
public:
    int sumOfN(int n) {
        // Write your code here
        int sum = 0;
        for (int i = 1; i <= n; ++i) {
            sum += i;
        }
        return sum;
    }
};`,
			GoSC: `package main

func sumOfN(n int) int {
    sum := 0
    for i := 1; i <= n; i++ {
        sum += i
    }
    return sum
}`,
			TestCases: []TestCase{
				{Input: "5", Expected: "15", IsHidden: false},
				{Input: "10", Expected: "55", IsHidden: false},
				{Input: "100", Expected: "5050", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4c0005",
			Slug:       "str2",
			Title:      "Palindrome String Check",
			Difficulty: "Easy",
			Topic:      "Strings",
			XP:         50,
			Statement:  "Check if a given string `s` is a palindrome, considering only alphanumeric characters and ignoring cases.\nReturn `true` if it is, and `false` otherwise.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca001",
			Tags:       []string{"Strings", "Two Pointers"},
			Examples: []Example{
				{
					Input:       "s = \"A man, a plan, a canal: Panama\"",
					Output:      "true",
					Explanation: "\"amanaplanacanalpanama\" is a palindrome.",
				},
				{
					Input:       "s = \"race a car\"",
					Output:      "false",
					Explanation: "\"raceacar\" is not a palindrome.",
				},
			},
			Hints: []Hint{
				{
					Title: "Cleaning the string",
					Body:  "First clean the string by keeping only alphanumeric characters and converting everything to lowercase. Then check if it equals its reverse.",
				},
			},
			JavascriptSC: `/**
 * @param {string} s
 * @return {boolean}
 */
function isPalindrome(s) {
    // Write your code here
    const clean = s.toLowerCase().replace(/[^a-z0-9]/g, '');
    return clean === clean.split('').reverse().join('');
}`,
			PythonSC: `def isPalindrome(s: str) -> bool:
    # Write your code here
    clean = "".join(c.lower() for c in s if c.isalnum())
    return clean == clean[::-1]`,
			JavaSC: `public class Solution {
    public boolean isPalindrome(String s) {
        // Write your code here
        String clean = s.toLowerCase().replaceAll("[^a-z0-9]", "");
        int left = 0, right = clean.length() - 1;
        while (left < right) {
            if (clean.charAt(left) != clean.charAt(right)) return false;
            left++;
            right--;
        }
        return true;
    }
}`,
			CppSC: `#include <string>
#include <cctype>
using namespace std;

class Solution {
public:
    bool isPalindrome(string s) {
        // Write your code here
        string clean = "";
        for (char c : s) {
            if (isalnum(c)) {
                clean += tolower(c);
            }
        }
        int left = 0, right = clean.length() - 1;
        while (left < right) {
            if (clean[left] != clean[right]) return false;
            left++;
            right--;
        }
        return true;
    }
};`,
			GoSC: `package main

import "unicode"

func isPalindrome(s string) bool {
    var clean []rune
    for _, r := range s {
        if unicode.IsLetter(r) || unicode.IsDigit(r) {
            clean = append(clean, unicode.ToLower(r))
        }
    }
    left, right := 0, len(clean)-1
    for left < right {
        if clean[left] != clean[right] {
            return false
        }
        left++
        right--
    }
    return true
}`,
			TestCases: []TestCase{
				{Input: "\"A man, a plan, a canal: Panama\"", Expected: "true", IsHidden: false},
				{Input: "\"race a car\"", Expected: "false", IsHidden: false},
				{Input: "\" \"", Expected: "true", IsHidden: true},
			},
		},
	}

	// Dynamically append the other 37 placeholder problems
	type problemMeta struct {
		Slug       string
		Title      string
		Difficulty string
		Topic      string
	}

	metaList := []problemMeta{
		{"op2", "Relational & Logical Operators", "Easy", "Operators"},
		{"op3", "Bitwise Left & Right Shift", "Medium", "Operators"},
		{"op4", "Ternary Operator Practice", "Easy", "Operators"},
		{"op5", "Operator Precedence Rules", "Medium", "Operators"},
		{"op6", "Compound Assignment Challenges", "Hard", "Operators"},
		{"cond2", "Switch-Case Grade Checker", "Easy", "Conditionals"},
		{"cond3", "Nested If Conditions", "Medium", "Conditionals"},
		{"cond4", "Leap Year Calculator", "Medium", "Conditionals"},
		{"cond5", "Find Maximum of Three Numbers", "Easy", "Conditionals"},
		{"cond6", "Tax Bracket Calculator", "Hard", "Conditionals"},
		{"loop2", "Factorial Calculator (While)", "Medium", "Loops"},
		{"loop3", "Fibonacci Series Generator", "Medium", "Loops"},
		{"loop4", "Prime Number Verification", "Medium", "Loops"},
		{"loop5", "Pattern Printing - Pyramid", "Hard", "Loops"},
		{"loop6", "Do-While Menu Driven Program", "Easy", "Loops"},
		{"func1", "Write Your First Function", "Easy", "Functions"},
		{"func2", "Area of Circle Calculator", "Easy", "Functions"},
		{"func3", "Recursive Factorial Function", "Medium", "Functions"},
		{"func4", "Function Overloading Challenge", "Medium", "Functions"},
		{"func5", "Pass by Value vs Reference", "Hard", "Functions"},
		{"func6", "Closure and Scope Practice", "Hard", "Functions"},
		{"arr1", "Find Min & Max in Array", "Easy", "Arrays"},
		{"arr2", "Reverse an Array in Place", "Easy", "Arrays"},
		{"arr3", "Search Element (Linear vs Binary)", "Medium", "Arrays"},
		{"arr4", "Merge Two Sorted Arrays", "Medium", "Arrays"},
		{"arr6", "Rotate 2D Matrix (Rotate Image)", "Hard", "Arrays"},
		{"str1", "Count Vowels and Consonants", "Easy", "Strings"},
		{"str3", "Reverse Words in a Sentence", "Medium", "Strings"},
		{"str4", "Anagram Validation", "Medium", "Strings"},
		{"str5", "String Substring Search (KMP)", "Hard", "Strings"},
		{"str6", "Longest Substring Without Repeating", "Hard", "Strings"},
		{"obj1", "Create Simple Student Object", "Easy", "Objects"},
		{"obj2", "Accessors & Mutators (Get/Set)", "Easy", "Objects"},
		{"obj3", "Object Serialization to JSON", "Medium", "Objects"},
		{"obj4", "Inheritance and Prototype Chain", "Medium", "Objects"},
		{"obj5", "Deep Copy vs Shallow Copy", "Hard", "Objects"},
		{"obj6", "Design Patterns: Singleton Object", "Hard", "Objects"},
	}

	for idx, m := range metaList {
		var setID string
		var xp int
		switch m.Difficulty {
		case "Easy":
			setID = "54574a34-9a68-4e65-ab9a-af05db4ca001"
			xp = 50
		case "Medium":
			setID = "54574a34-9a68-4e65-ab9a-af05db4ca002"
			xp = 100
		case "Hard":
			setID = "54574a34-9a68-4e65-ab9a-af05db4ca003"
			xp = 150
		}

		problems = append(problems, Problem{
			ID:         fmt.Sprintf("54574a34-9a68-4e65-ab9a-af05db4c%04d", 100+idx),
			Slug:       m.Slug,
			Title:      m.Title,
			Difficulty: m.Difficulty,
			Topic:      m.Topic,
			XP:         xp,
			Statement:  fmt.Sprintf("In this challenge, you are asked to solve the algorithmic problem: **%s**.\n\nSince this is a placeholder challenge, you simply need to return the input value directly.\n\nMake sure your function has the correct signature as shown in the starter code.", m.Title),
			SetID:      setID,
			Tags:       []string{m.Topic, m.Difficulty},
			Examples: []Example{
				{
					Input:       "1",
					Output:      "1",
					Explanation: "The function returns the input value directly.",
				},
			},
			Hints: []Hint{
				{
					Title: "Problem Solving",
					Body:  "Follow standard logic and constraints defined for the topic.",
				},
			},
			JavascriptSC: `/**
 * Solve the ` + m.Title + ` problem.
 * @param {any} input
 * @return {any}
 */
function solveChallenge(input) {
    // Write your code here
    return input;
}`,
			PythonSC: `def solveChallenge(input_val):
    # Write your code here
    return input_val`,
			JavaSC: `public class Solution {
    public Object solveChallenge(Object inputVal) {
        // Write your code here
        return inputVal;
    }
}`,
			CppSC: `#include <iostream>
using namespace std;

class Solution {
public:
    void solveChallenge() {
        // Write your code here
    }
};`,
			GoSC: `package main

func solveChallenge(input string) string {
    return input
}`,
			TestCases: []TestCase{
				{Input: "1", Expected: "1", IsHidden: false},
				{Input: "2", Expected: "2", IsHidden: true},
			},
		})
	}

	for _, p := range problems {
		// Insert problem
		_, err := conn.Exec(ctx, `
			INSERT INTO problems (id, slug, title, difficulty, topic, xp, statement, set_id, created_at, updated_at)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, now(), now())
		`, p.ID, p.Slug, p.Title, p.Difficulty, p.Topic, p.XP, p.Statement, p.SetID)
		if err != nil {
			fmt.Printf("Failed to insert problem %s: %v\n", p.Title, err)
			os.Exit(1)
		}

		// Insert tags
		for _, tag := range p.Tags {
			_, err = conn.Exec(ctx, `
				INSERT INTO problem_tags (problem_id, tag)
				VALUES ($1, $2)
			`, p.ID, tag)
			if err != nil {
				fmt.Printf("Failed to insert tag %s for problem %s: %v\n", tag, p.Title, err)
			}
		}

		// Insert examples
		for idx, ex := range p.Examples {
			_, err = conn.Exec(ctx, `
				INSERT INTO examples (problem_id, input, output, explanation, order_index)
				VALUES ($1, $2, $3, $4, $5)
			`, p.ID, ex.Input, ex.Output, ex.Explanation, idx)
			if err != nil {
				fmt.Printf("Failed to insert example for problem %s: %v\n", p.Title, err)
			}
		}

		// Insert hints
		for idx, h := range p.Hints {
			_, err = conn.Exec(ctx, `
				INSERT INTO hints (problem_id, order_index, title, body)
				VALUES ($1, $2, $3, $4)
			`, p.ID, idx, h.Title, h.Body)
			if err != nil {
				fmt.Printf("Failed to insert hint for problem %s: %v\n", p.Title, err)
			}
		}

		// Insert starter codes
		_, err = conn.Exec(ctx, `
			INSERT INTO starter_codes (problem_id, javascript, python, java, cpp, go)
			VALUES ($1, $2, $3, $4, $5, $6)
		`, p.ID, p.JavascriptSC, p.PythonSC, p.JavaSC, p.CppSC, p.GoSC)
		if err != nil {
			fmt.Printf("Failed to insert starter codes for problem %s: %v\n", p.Title, err)
		}

		// Insert test cases
		for idx, tc := range p.TestCases {
			_, err = conn.Exec(ctx, `
				INSERT INTO test_cases (problem_id, input, expected_output, is_hidden, time_limit_ms, memory_limit_mb, order_index)
				VALUES ($1, $2, $3, $4, 2000, 256, $5)
			`, p.ID, tc.Input, tc.Expected, tc.IsHidden, idx)
			if err != nil {
				fmt.Printf("Failed to insert test case for problem %s: %v\n", p.Title, err)
			}
		}
	}

	fmt.Println("Seeded all problems, starter codes, examples, hints, and test cases successfully!")
}
