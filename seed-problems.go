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
					Title: "Loops",
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
					Title: "Strings",
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

		// ─── 10 Real LeetCode Array Problems for Masters of Algorithms (Advanced) ───
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0011",
			Slug:       "leetcode-two-sum",
			Title:      "Two Sum",
			Difficulty: "Easy",
			Topic:      "Array",
			XP:         50,
			Statement:  "Given an array of integers `nums` and an integer `target`, return *indices of the two numbers such that they add up to `target`*.\n\nYou may assume that each input would have ***exactly* one solution**, and you may not use the *same* element twice.\n\nYou can return the answer in any order.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca001",
			Tags:       []string{"Array", "Hash Table"},
			Examples: []Example{
				{
					Input:       "nums = [2,7,11,15], target = 9",
					Output:      "[0,1]",
					Explanation: "Because nums[0] + nums[1] == 9, we return [0, 1].",
				},
			},
			Hints: []Hint{
				{
					Title: "Brute Force",
					Body:  "Loop through each element and find if there is another value that adds up to target.",
				},
				{
					Title: "Hash Map",
					Body:  "Use a hash map to store elements and their indices for O(1) complement lookup.",
				},
			},
			JavascriptSC: `/**
 * @param {number[]} nums
 * @param {number} target
 * @return {number[]}
 */
function twoSum(nums, target) {
    const map = new Map();
    for (let i = 0; i < nums.length; i++) {
        const diff = target - nums[i];
        if (map.has(diff)) {
            return [map.get(diff), i];
        }
        map.set(nums[i], i);
    }
    return [];
}`,
			PythonSC: `def twoSum(nums: list[int], target: int) -> list[int]:
    seen = {}
    for i, num in enumerate(nums):
        diff = target - num
        if diff in seen:
            return [seen[diff], i]
        seen[num] = i
    return []`,
			JavaSC: `import java.util.*;

public class Solution {
    public int[] twoSum(int[] nums, int target) {
        Map<Integer, Integer> map = new HashMap<>();
        for (int i = 0; i < nums.length; i++) {
            int diff = target - nums[i];
            if (map.containsKey(diff)) {
                return new int[] { map.get(diff), i };
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
        unordered_map<int, int> seen;
        for (int i = 0; i < nums.size(); i++) {
            int diff = target - nums[i];
            if (seen.count(diff)) {
                return {seen[diff], i};
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
        diff := target - num
        if idx, found := seen[diff]; found {
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
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0012",
			Slug:       "best-time-to-buy-and-sell-stock",
			Title:      "Best Time to Buy and Sell Stock",
			Difficulty: "Easy",
			Topic:      "Array",
			XP:         50,
			Statement:  "You are given an array `prices` where `prices[i]` is the price of a given stock on the `i-th` day.\n\nYou want to maximize your profit by choosing a **single day** to buy one stock and choosing a **different day in the future** to sell that stock.\n\nReturn *the maximum profit you can achieve from this transaction*. If you cannot achieve any profit, return `0`.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca001",
			Tags:       []string{"Array", "Dynamic Programming"},
			Examples: []Example{
				{
					Input:       "prices = [7,1,5,3,6,4]",
					Output:      "5",
					Explanation: "Buy on day 2 (price = 1) and sell on day 5 (price = 6), profit = 6-1 = 5.",
				},
			},
			Hints: []Hint{
				{
					Title: "Tracking Min Price",
					Body:  "Track the minimum price encountered so far and calculate the potential profit at each step.",
				},
			},
			JavascriptSC: `/**
 * @param {number[]} prices
 * @return {number}
 */
function maxProfit(prices) {
    let minPrice = Infinity;
    let maxProfit = 0;
    for (let i = 0; i < prices.length; i++) {
        if (prices[i] < minPrice) {
            minPrice = prices[i];
        } else if (prices[i] - minPrice > maxProfit) {
            maxProfit = prices[i] - minPrice;
        }
    }
    return maxProfit;
}`,
			PythonSC: `def maxProfit(prices: list[int]) -> int:
    min_price = float('inf')
    max_profit = 0
    for price in prices:
        if price < min_price:
            min_price = price
        elif price - min_price > max_profit:
            max_profit = price - min_price
    return max_profit`,
			JavaSC: `public class Solution {
    public int maxProfit(int[] prices) {
        int minPrice = Integer.MAX_VALUE;
        int maxProfit = 0;
        for (int price : prices) {
            if (price < minPrice) {
                minPrice = price;
            } else if (price - minPrice > maxProfit) {
                maxProfit = price - minPrice;
            }
        }
        return maxProfit;
    }
}`,
			CppSC: `#include <vector>
#include <algorithm>
using namespace std;

class Solution {
public:
    int maxProfit(vector<int>& prices) {
        int minPrice = 1e9;
        int maxProfit = 0;
        for (int price : prices) {
            if (price < minPrice) {
                minPrice = price;
            } else {
                maxProfit = max(maxProfit, price - minPrice);
            }
        }
        return maxProfit;
    }
};`,
			GoSC: `package main

func maxProfit(prices []int) int {
    minPrice := 1000000000
    maxProfit := 0
    for _, price := range prices {
        if price < minPrice {
            minPrice = price
        } else if price - minPrice > maxProfit {
            maxProfit = price - minPrice
        }
    }
    return maxProfit
}`,
			TestCases: []TestCase{
				{Input: "[7,1,5,3,6,4]", Expected: "5", IsHidden: false},
				{Input: "[7,6,4,3,1]", Expected: "0", IsHidden: false},
				{Input: "[2,4,1]", Expected: "2", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0013",
			Slug:       "contains-duplicate",
			Title:      "Contains Duplicate",
			Difficulty: "Easy",
			Topic:      "Array",
			XP:         50,
			Statement:  "Given an integer array `nums`, return `true` if any value appears **at least twice** in the array, and return `false` if every element is distinct.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca001",
			Tags:       []string{"Array", "Hash Table"},
			Examples: []Example{
				{
					Input:       "nums = [1,2,3,1]",
					Output:      "true",
					Explanation: "The value 1 appears at indices 0 and 3.",
				},
			},
			Hints: []Hint{
				{
					Title: "Hash Set",
					Body:  "Keep track of visited numbers using a hash set. If you encounter a number already in the set, a duplicate exists.",
				},
			},
			JavascriptSC: `/**
 * @param {number[]} nums
 * @return {boolean}
 */
function containsDuplicate(nums) {
    const set = new Set(nums);
    return set.size !== nums.length;
}`,
			PythonSC: `def containsDuplicate(nums: list[int]) -> bool:
    return len(set(nums)) != len(nums)`,
			JavaSC: `import java.util.*;

public class Solution {
    public boolean containsDuplicate(int[] nums) {
        Set<Integer> set = new HashSet<>();
        for (int num : nums) {
            if (!set.add(num)) return true;
        }
        return false;
    }
}`,
			CppSC: `#include <vector>
#include <unordered_set>
using namespace std;

class Solution {
public:
    bool containsDuplicate(vector<int>& nums) {
        unordered_set<int> set;
        for (int num : nums) {
            if (set.count(num)) return true;
            set.insert(num);
        }
        return false;
    }
};`,
			GoSC: `package main

func containsDuplicate(nums []int) bool {
    seen := make(map[int]bool)
    for _, num := range nums {
        if seen[num] {
            return true
        }
        seen[num] = true
    }
    return false
}`,
			TestCases: []TestCase{
				{Input: "[1,2,3,1]", Expected: "true", IsHidden: false},
				{Input: "[1,2,3,4]", Expected: "false", IsHidden: false},
				{Input: "[1,1,1,3,3,4,3,2,4,2]", Expected: "true", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0014",
			Slug:       "product-of-array-except-self",
			Title:      "Product of Array Except Self",
			Difficulty: "Medium",
			Topic:      "Array",
			XP:         100,
			Statement:  "Given an integer array `nums`, return *an array* `answer` *such that* `answer[i]` *is equal to the product of all the elements of* `nums` *except* `nums[i]`.\n\nThe product of any prefix or suffix of `nums` is **guaranteed** to fit in a **32-bit integer**.\n\nYou must write an algorithm that runs in $O(n)$ time and without using the division operation.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca002",
			Tags:       []string{"Array", "Prefix Product"},
			Examples: []Example{
				{
					Input:       "nums = [1,2,3,4]",
					Output:      "[24,12,8,6]",
					Explanation: "answer[0] = 2*3*4 = 24, answer[1] = 1*3*4 = 12, etc.",
				},
			},
			Hints: []Hint{
				{
					Title: "Prefix & Suffix Products",
					Body:  "Compute prefix products from left to right, and suffix products from right to left, storing them to construct the output.",
				},
			},
			JavascriptSC: `/**
 * @param {number[]} nums
 * @return {number[]}
 */
function productExceptSelf(nums) {
    const n = nums.length;
    const res = new Array(n).fill(1);
    let left = 1;
    for (let i = 0; i < n; i++) {
        res[i] *= left;
        left *= nums[i];
    }
    let right = 1;
    for (let i = n - 1; i >= 0; i--) {
        res[i] *= right;
        right *= nums[i];
    }
    return res;
}`,
			PythonSC: `def productExceptSelf(nums: list[int]) -> list[int]:
    n = len(nums)
    res = [1] * n
    left = 1
    for i in range(n):
        res[i] *= left
        left *= nums[i]
    right = 1
    for i in range(n - 1, -1, -1):
        res[i] *= right
        right *= nums[i]
    return res`,
			JavaSC: `public class Solution {
    public int[] productExceptSelf(int[] nums) {
        int n = nums.length;
        int[] res = new int[n];
        for (int i = 0; i < n; i++) res[i] = 1;
        int left = 1;
        for (int i = 0; i < n; i++) {
            res[i] *= left;
            left *= nums[i];
        }
        int right = 1;
        for (int i = n - 1; i >= 0; i--) {
            res[i] *= right;
            right *= nums[i];
        }
        return res;
    }
}`,
			CppSC: `#include <vector>
using namespace std;

class Solution {
public:
    vector<int> productExceptSelf(vector<int>& nums) {
        int n = nums.size();
        vector<int> res(n, 1);
        int left = 1;
        for (int i = 0; i < n; i++) {
            res[i] *= left;
            left *= nums[i];
        }
        int right = 1;
        for (int i = n - 1; i >= 0; i--) {
            res[i] *= right;
            right *= nums[i];
        }
        return res;
    }
};`,
			GoSC: `package main

func productExceptSelf(nums []int) []int {
    n := len(nums)
    res := make([]int, n)
    for i := range res {
        res[i] = 1
    }
    left := 1
    for i := 0; i < n; i++ {
        res[i] *= left
        left *= nums[i]
    }
    right := 1
    for i := n - 1; i >= 0; i-- {
        res[i] *= right
        right *= nums[i]
    }
    return res
}`,
			TestCases: []TestCase{
				{Input: "[1,2,3,4]", Expected: "[24,12,8,6]", IsHidden: false},
				{Input: "[-1,1,0,-3,3]", Expected: "[0,0,9,0,0]", IsHidden: false},
				{Input: "[2,3]", Expected: "[3,2]", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0015",
			Slug:       "maximum-subarray",
			Title:      "Maximum Subarray",
			Difficulty: "Medium",
			Topic:      "Array",
			XP:         100,
			Statement:  "Given an integer array `nums`, find the subarray with the largest sum, and return its sum.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca002",
			Tags:       []string{"Array", "Kadane's Algorithm"},
			Examples: []Example{
				{
					Input:       "nums = [-2,1,-3,4,-1,2,1,-5,4]",
					Output:      "6",
					Explanation: "The subarray [4,-1,2,1] has the largest sum = 6.",
				},
			},
			Hints: []Hint{
				{
					Title: "Kadane's Algorithm",
					Body:  "Track the current subarray sum and update the maximum sum whenever the current sum exceeds it.",
				},
			},
			JavascriptSC: `/**
 * @param {number[]} nums
 * @return {number}
 */
function maxSubArray(nums) {
    let maxSum = nums[0];
    let currentSum = nums[0];
    for (let i = 1; i < nums.length; i++) {
        currentSum = Math.max(nums[i], currentSum + nums[i]);
        maxSum = Math.max(maxSum, currentSum);
    }
    return maxSum;
}`,
			PythonSC: `def maxSubArray(nums: list[int]) -> int:
    max_sum = nums[0]
    current_sum = nums[0]
    for num in nums[1:]:
        current_sum = max(num, current_sum + num)
        max_sum = max(max_sum, current_sum)
    return max_sum`,
			JavaSC: `public class Solution {
    public int maxSubArray(int[] nums) {
        int maxSum = nums[0];
        int currentSum = nums[0];
        for (int i = 1; i < nums.length; i++) {
            currentSum = Math.max(nums[i], currentSum + nums[i]);
            maxSum = Math.max(maxSum, currentSum);
        }
        return maxSum;
    }
}`,
			CppSC: `#include <vector>
#include <algorithm>
using namespace std;

class Solution {
public:
    int maxSubArray(vector<int>& nums) {
        int maxSum = nums[0];
        int currentSum = nums[0];
        for (int i = 1; i < nums.size(); i++) {
            currentSum = max(nums[i], currentSum + nums[i]);
            maxSum = max(maxSum, currentSum);
        }
        return maxSum;
    }
};`,
			GoSC: `package main

func maxSubArray(nums []int) int {
    maxSum := nums[0]
    currentSum := nums[0]
    for i := 1; i < len(nums); i++ {
        if nums[i] > currentSum + nums[i] {
            currentSum = nums[i]
        } else {
            currentSum = currentSum + nums[i]
        }
        if currentSum > maxSum {
            maxSum = currentSum
        }
    }
    return maxSum
}`,
			TestCases: []TestCase{
				{Input: "[-2,1,-3,4,-1,2,1,-5,4]", Expected: "6", IsHidden: false},
				{Input: "[1]", Expected: "1", IsHidden: false},
				{Input: "[5,4,-1,7,8]", Expected: "23", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0016",
			Slug:       "maximum-product-subarray",
			Title:      "Maximum Product Subarray",
			Difficulty: "Medium",
			Topic:      "Array",
			XP:         100,
			Statement:  "Given an integer array `nums`, find a subarray that has the largest product, and return the product.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca002",
			Tags:       []string{"Array", "Dynamic Programming"},
			Examples: []Example{
				{
					Input:       "nums = [2,3,-2,4]",
					Output:      "6",
					Explanation: "[2,3] has the largest product = 6.",
				},
			},
			Hints: []Hint{
				{
					Title: "DP States",
					Body:  "Since a negative number multiplied by a negative number becomes positive, track both the maximum product and the minimum product at each step.",
				},
			},
			JavascriptSC: `/**
 * @param {number[]} nums
 * @return {number}
 */
function maxProduct(nums) {
    let maxProd = nums[0];
    let minProd = nums[0];
    let res = nums[0];
    for (let i = 1; i < nums.length; i++) {
        const temp = Math.max(nums[i], maxProd * nums[i], minProd * nums[i]);
        minProd = Math.min(nums[i], maxProd * nums[i], minProd * nums[i]);
        maxProd = temp;
        res = Math.max(res, maxProd);
    }
    return res;
}`,
			PythonSC: `def maxProduct(nums: list[int]) -> int:
    max_prod = nums[0]
    min_prod = nums[0]
    res = nums[0]
    for num in nums[1:]:
        temp = max(num, max_prod * num, min_prod * num)
        min_prod = min(num, max_prod * num, min_prod * num)
        max_prod = temp
        res = max(res, max_prod)
    return res`,
			JavaSC: `public class Solution {
    public int maxProduct(int[] nums) {
        int maxProd = nums[0];
        int minProd = nums[0];
        int res = nums[0];
        for (int i = 1; i < nums.length; i++) {
            int temp = Math.max(nums[i], Math.max(maxProd * nums[i], minProd * nums[i]));
            minProd = Math.min(nums[i], Math.min(maxProd * nums[i], minProd * nums[i]));
            maxProd = temp;
            res = Math.max(res, maxProd);
        }
        return res;
    }
}`,
			CppSC: `#include <vector>
#include <algorithm>
using namespace std;

class Solution {
public:
    int maxProduct(vector<int>& nums) {
        int maxProd = nums[0];
        int minProd = nums[0];
        int res = nums[0];
        for (int i = 1; i < nums.size(); i++) {
            int temp = max({nums[i], maxProd * nums[i], minProd * nums[i]});
            minProd = min({nums[i], maxProd * nums[i], minProd * nums[i]});
            maxProd = temp;
            res = max(res, maxProd);
        }
        return res;
    }
};`,
			GoSC: `package main

func maxProduct(nums []int) int {
    maxProd := nums[0]
    minProd := nums[0]
    res := nums[0]
    for i := 1; i < len(nums); i++ {
        num := nums[i]
        t1 := maxProd * num
        t2 := minProd * num
        tempMax := num
        if t1 > tempMax { tempMax = t1 }
        if t2 > tempMax { tempMax = t2 }
        
        tempMin := num
        if t1 < tempMin { tempMin = t1 }
        if t2 < tempMin { tempMin = t2 }
        
        maxProd = tempMax
        minProd = tempMin
        if maxProd > res {
            res = maxProd
        }
    }
    return res
}`,
			TestCases: []TestCase{
				{Input: "[2,3,-2,4]", Expected: "6", IsHidden: false},
				{Input: "[-2,0,-1]", Expected: "0", IsHidden: false},
				{Input: "[-3,-1,-1]", Expected: "3", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0017",
			Slug:       "find-minimum-in-rotated-sorted-array",
			Title:      "Find Minimum in Rotated Sorted Array",
			Difficulty: "Medium",
			Topic:      "Array",
			XP:         100,
			Statement:  "Suppose an array of length `n` sorted in ascending order is rotated between `1` and `n` times.\n\nGiven the sorted rotated array `nums` of unique elements, return *the minimum element of this array*.\n\nYou must write an algorithm that runs in $O(\\log n)$ time.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca002",
			Tags:       []string{"Array", "Binary Search"},
			Examples: []Example{
				{
					Input:       "nums = [3,4,5,1,2]",
					Output:      "1",
					Explanation: "The original array was [1,2,3,4,5] rotated 3 times.",
				},
			},
			Hints: []Hint{
				{
					Title: "Binary Search",
					Body:  "Use binary search. Compare nums[mid] to nums[right] to determine which half is unsorted and contains the pivot.",
				},
			},
			JavascriptSC: `/**
 * @param {number[]} nums
 * @return {number}
 */
function findMin(nums) {
    let left = 0;
    let right = nums.length - 1;
    while (left < right) {
        const mid = Math.floor((left + right) / 2);
        if (nums[mid] > nums[right]) {
            left = mid + 1;
        } else {
            right = mid;
        }
    }
    return nums[left];
}`,
			PythonSC: `def findMin(nums: list[int]) -> int:
    left = 0
    right = len(nums) - 1
    while left < right:
        mid = (left + right) // 2
        if nums[mid] > nums[right]:
            left = mid + 1
        else:
            right = mid
    return nums[left]`,
			JavaSC: `public class Solution {
    public int findMin(int[] nums) {
        int left = 0;
        int right = nums.length - 1;
        while (left < right) {
            int mid = left + (right - left) / 2;
            if (nums[mid] > nums[right]) {
                left = mid + 1;
            } else {
                right = mid;
            }
        }
        return nums[left];
    }
}`,
			CppSC: `#include <vector>
using namespace std;

class Solution {
public:
    int findMin(vector<int>& nums) {
        int left = 0;
        int right = nums.size() - 1;
        while (left < right) {
            int mid = left + (right - left) / 2;
            if (nums[mid] > nums[right]) {
                left = mid + 1;
            } else {
                right = mid;
            }
        }
        return nums[left];
    }
};`,
			GoSC: `package main

func findMin(nums []int) int {
    left := 0
    right := len(nums) - 1
    for left < right {
        mid := (left + right) / 2
        if nums[mid] > nums[right] {
            left = mid + 1
        } else {
            right = mid
        }
    }
    return nums[left]
}`,
			TestCases: []TestCase{
				{Input: "[3,4,5,1,2]", Expected: "1", IsHidden: false},
				{Input: "[4,5,6,7,0,1,2]", Expected: "0", IsHidden: false},
				{Input: "[11,13,15,17]", Expected: "11", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0018",
			Slug:       "search-in-rotated-sorted-array",
			Title:      "Search in Rotated Sorted Array",
			Difficulty: "Medium",
			Topic:      "Array",
			XP:         100,
			Statement:  "There is an integer array `nums` sorted in ascending order (with distinct values).\n\nPrior to being passed to your function, `nums` is possibly rotated at an unknown pivot index `k` (`1 <= k < nums.length`).\n\nGiven the array `nums` after the possible rotation and an integer `target`, return *the index of* `target` *if it is in* `nums`, *or* `-1` *if it is not in* `nums`.\n\nYou must write an algorithm with $O(\\log n)$ runtime complexity.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca002",
			Tags:       []string{"Array", "Binary Search"},
			Examples: []Example{
				{
					Input:       "nums = [4,5,6,7,0,1,2], target = 0",
					Output:      "4",
					Explanation: "0 is at index 4.",
				},
			},
			Hints: []Hint{
				{
					Title: "Two Halves",
					Body:  "During binary search, one of the halves (left to mid or mid to right) will always be normally sorted. Use that to check target bounds.",
				},
			},
			JavascriptSC: `/**
 * @param {number[]} nums
 * @param {number} target
 * @return {number}
 */
function search(nums, target) {
    let left = 0, right = nums.length - 1;
    while (left <= right) {
        const mid = Math.floor((left + right) / 2);
        if (nums[mid] === target) return mid;
        if (nums[left] <= nums[mid]) {
            if (target >= nums[left] && target < nums[mid]) {
                right = mid - 1;
            } else {
                left = mid + 1;
            }
        } else {
            if (target > nums[mid] && target <= nums[right]) {
                left = mid + 1;
            } else {
                right = mid - 1;
            }
        }
    }
    return -1;
}`,
			PythonSC: `def search(nums: list[int], target: int) -> int:
    left, right = 0, len(nums) - 1
    while left <= right:
        mid = (left + right) // 2
        if nums[mid] == target:
            return mid
        if nums[left] <= nums[mid]:
            if nums[left] <= target < nums[mid]:
                right = mid - 1
            else:
                left = mid + 1
        else:
            if nums[mid] < target <= nums[right]:
                left = mid + 1
            else:
                right = mid - 1
    return -1`,
			JavaSC: `public class Solution {
    public int search(int[] nums, int target) {
        int left = 0, right = nums.length - 1;
        while (left <= right) {
            int mid = left + (right - left) / 2;
            if (nums[mid] == target) return mid;
            if (nums[left] <= nums[mid]) {
                if (target >= nums[left] && target < nums[mid]) {
                    right = mid - 1;
                } else {
                    left = mid + 1;
                }
            } else {
                if (target > nums[mid] && target <= nums[right]) {
                    left = mid + 1;
                } else {
                    right = mid - 1;
                }
            }
        }
        return -1;
    }
}`,
			CppSC: `#include <vector>
using namespace std;

class Solution {
public:
    int search(vector<int>& nums, int target) {
        int left = 0, right = nums.size() - 1;
        while (left <= right) {
            int mid = left + (right - left) / 2;
            if (nums[mid] == target) return mid;
            if (nums[left] <= nums[mid]) {
                if (target >= nums[left] && target < nums[mid]) {
                    right = mid - 1;
                } else {
                    left = mid + 1;
                }
            } else {
                if (target > nums[mid] && target <= nums[right]) {
                    left = mid + 1;
                } else {
                    right = mid - 1;
                }
            }
        }
        return -1;
    }
};`,
			GoSC: `package main

func search(nums []int, target int) int {
    left, right := 0, len(nums)-1
    for left <= right {
        mid := (left + right) / 2
        if nums[mid] == target {
            return mid
        }
        if nums[left] <= nums[mid] {
            if target >= nums[left] && target < nums[mid] {
                right = mid - 1
            } else {
                left = mid + 1
            }
        } else {
            if target > nums[mid] && target <= nums[right] {
                left = mid + 1
            } else {
                right = mid - 1
            }
        }
    }
    return -1
}`,
			TestCases: []TestCase{
				{Input: "[4,5,6,7,0,1,2]\n0", Expected: "4", IsHidden: false},
				{Input: "[4,5,6,7,0,1,2]\n3", Expected: "-1", IsHidden: false},
				{Input: "[1]\n0", Expected: "-1", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0019",
			Slug:       "3sum",
			Title:      "3Sum",
			Difficulty: "Medium",
			Topic:      "Array",
			XP:         100,
			Statement:  "Given an integer array `nums`, return all the triplets `[nums[i], nums[j], nums[k]]` such that `i != j`, `i != k`, and `j != k`, and `nums[i] + nums[j] + nums[k] == 0`.\n\nNotice that the solution set must not contain duplicate triplets.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca002",
			Tags:       []string{"Array", "Two Pointers"},
			Examples: []Example{
				{
					Input:       "nums = [-1,0,1,2,-1,-4]",
					Output:      "[[-1,-1,2],[-1,0,1]]",
					Explanation: "nums[0] + nums[1] + nums[3] = (-1) + 0 + 1 = 0. Distinct triplets are [[-1,-1,2],[-1,0,1]].",
				},
			},
			Hints: []Hint{
				{
					Title: "Sorting & Two Pointers",
					Body:  "Sort the array. Iterate with a fixed element, and use a two-pointer scan (left and right) on the remaining elements.",
				},
			},
			JavascriptSC: `/**
 * @param {number[]} nums
 * @return {number[][]}
 */
function threeSum(nums) {
    nums.sort((a, b) => a - b);
    const res = [];
    for (let i = 0; i < nums.length - 2; i++) {
        if (i > 0 && nums[i] === nums[i - 1]) continue;
        let left = i + 1, right = nums.length - 1;
        while (left < right) {
            const sum = nums[i] + nums[left] + nums[right];
            if (sum === 0) {
                res.push([nums[i], nums[left], nums[right]]);
                while (left < right && nums[left] === nums[left + 1]) left++;
                while (left < right && nums[right] === nums[right - 1]) right--;
                left++;
                right--;
            } else if (sum < 0) {
                left++;
            } else {
                right--;
            }
        }
    }
    return res;
}`,
			PythonSC: `def threeSum(nums: list[int]) -> list[list[int]]:
    nums.sort()
    res = []
    for i in range(len(nums) - 2):
        if i > 0 and nums[i] == nums[i - 1]:
            continue
        left, right = i + 1, len(nums) - 1
        while left < right:
            s = nums[i] + nums[left] + nums[right]
            if s == 0:
                res.append([nums[i], nums[left], nums[right]])
                while left < right and nums[left] == nums[left + 1]:
                    left += 1
                while left < right and nums[right] == nums[right - 1]:
                    right -= 1
                left += 1
                right -= 1
            elif s < 0:
                left += 1
            else:
                right -= 1
    return res`,
			JavaSC: `import java.util.*;

public class Solution {
    public List<List<Integer>> threeSum(int[] nums) {
        Arrays.sort(nums);
        List<List<Integer>> res = new ArrayList<>();
        for (int i = 0; i < nums.length - 2; i++) {
            if (i > 0 && nums[i] == nums[i - 1]) continue;
            int left = i + 1, right = nums.length - 1;
            while (left < right) {
                int sum = nums[i] + nums[left] + nums[right];
                if (sum == 0) {
                    res.add(Arrays.asList(nums[i], nums[left], nums[right]));
                    while (left < right && nums[left] == nums[left + 1]) left++;
                    while (left < right && nums[right] == nums[right - 1]) right--;
                    left++;
                    right--;
                } else if (sum < 0) {
                    left++;
                } else {
                    right--;
                }
            }
        }
        return res;
    }
}`,
			CppSC: `#include <vector>
#include <algorithm>
using namespace std;

class Solution {
public:
    vector<vector<int>> threeSum(vector<int>& nums) {
        sort(nums.begin(), nums.end());
        vector<vector<int>> res;
        for (int i = 0; i < (int)nums.size() - 2; i++) {
            if (i > 0 && nums[i] == nums[i - 1]) continue;
            int left = i + 1, right = nums.size() - 1;
            while (left < right) {
                int sum = nums[i] + nums[left] + nums[right];
                if (sum == 0) {
                    res.push_back({nums[i], nums[left], nums[right]});
                    while (left < right && nums[left] == nums[left + 1]) left++;
                    while (left < right && nums[right] == nums[right - 1]) right--;
                    left++;
                    right--;
                } else if (sum < 0) {
                    left++;
                } else {
                    right--;
                }
            }
        }
        return res;
    }
};`,
			GoSC: `package main

import "sort"

func threeSum(nums []int) [][]int {
    sort.Ints(nums)
    res := [][]int{}
    for i := 0; i < len(nums)-2; i++ {
        if i > 0 && nums[i] == nums[i-1] {
            continue
        }
        left, right := i+1, len(nums)-1
        for left < right {
            sum := nums[i] + nums[left] + nums[right]
            if sum == 0 {
                res = append(res, []int{nums[i], nums[left], nums[right]})
                for left < right && nums[left] == nums[left+1] { left++ }
                for left < right && nums[right] == nums[right-1] { right-- }
                left++
                right--
            } else if sum < 0 {
                left++
            } else {
                right--
            }
        }
    }
    return res
}`,
			TestCases: []TestCase{
				{Input: "[-1,0,1,2,-1,-4]", Expected: "[[-1,-1,2],[-1,0,1]]", IsHidden: false},
				{Input: "[0,1,1]", Expected: "[]", IsHidden: false},
				{Input: "[0,0,0]", Expected: "[[0,0,0]]", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0020",
			Slug:       "container-with-most-water",
			Title:      "Container With Most Water",
			Difficulty: "Medium",
			Topic:      "Array",
			XP:         100,
			Statement:  "You are given an integer array `height` of length `n`. There are `n` vertical lines drawn such that the two endpoints of the `i-th` line are `(i, 0)` and `(i, height[i])`.\n\nFind two lines that together with the x-axis form a container, such that the container contains the most water.\n\nReturn *the maximum amount of water a container can store*.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca002",
			Tags:       []string{"Array", "Two Pointers"},
			Examples: []Example{
				{
					Input:       "height = [1,8,6,2,5,4,8,3,7]",
					Output:      "49",
					Explanation: "The maximum area of water is 49 (heights 8 and 7, width 7).",
				},
			},
			Hints: []Hint{
				{
					Title: "Two Pointers",
					Body:  "Start pointers at both ends. Move the pointer pointing to the shorter line inward to try and find a taller boundary.",
				},
			},
			JavascriptSC: `/**
 * @param {number[]} height
 * @return {number}
 */
function maxArea(height) {
    let maxA = 0;
    let left = 0, right = height.length - 1;
    while (left < right) {
        const h = Math.min(height[left], height[right]);
        maxA = Math.max(maxA, h * (right - left));
        if (height[left] < height[right]) {
            left++;
        } else {
            right--;
        }
    }
    return maxA;
}`,
			PythonSC: `def maxArea(height: list[int]) -> int:
    max_a = 0
    left, right = 0, len(height) - 1
    while left < right:
        h = min(height[left], height[right])
        max_a = max(max_a, h * (right - left))
        if height[left] < height[right]:
            left += 1
        else:
            right -= 1
    return max_a`,
			JavaSC: `public class Solution {
    public int maxArea(int[] height) {
        int maxA = 0;
        int left = 0, right = height.length - 1;
        while (left < right) {
            int h = Math.min(height[left], height[right]);
            maxA = Math.max(maxA, h * (right - left));
            if (height[left] < height[right]) {
                left++;
            } else {
                right--;
            }
        }
        return maxA;
    }
}`,
			CppSC: `#include <vector>
#include <algorithm>
using namespace std;

class Solution {
public:
    int maxArea(vector<int>& height) {
        int maxA = 0;
        int left = 0, right = height.size() - 1;
        while (left < right) {
            int h = min(height[left], height[right]);
            maxA = max(maxA, h * (right - left));
            if (height[left] < height[right]) {
                left++;
            } else {
                right--;
            }
        }
        return maxA;
    }
};`,
			GoSC: `package main

func maxArea(height []int) int {
    maxA := 0
    left, right := 0, len(height)-1
    for left < right {
        h := height[left]
        if height[right] < h {
            h = height[right]
        }
        area := h * (right - left)
        if area > maxA {
            maxA = area
        }
        if height[left] < height[right] {
            left++
        } else {
            right--
        }
    }
    return maxA
}`,
			TestCases: []TestCase{
				{Input: "[1,8,6,2,5,4,8,3,7]", Expected: "49", IsHidden: false},
				{Input: "[1,1]", Expected: "1", IsHidden: false},
				{Input: "[4,3,2,1,4]", Expected: "16", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0021",
			Slug:       "median-of-two-sorted-arrays",
			Title:      "Median of Two Sorted Arrays",
			Difficulty: "Hard",
			Topic:      "Array",
			XP:         150,
			Statement:  "Given two sorted arrays `nums1` and `nums2` of size `m` and `n` respectively, return the median of the two sorted arrays.\n\nThe overall run time complexity should be $O(\\log(m+n))$.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Array", "Binary Search", "Divide and Conquer"},
			Examples: []Example{
				{
					Input:       "nums1 = [1,3], nums2 = [2]",
					Output:      "2.00000",
					Explanation: "merged array = [1,2,3] and median is 2.",
				},
				{
					Input:       "nums1 = [1,2], nums2 = [3,4]",
					Output:      "2.50000",
					Explanation: "merged array = [1,2,3,4] and median is (2 + 3) / 2 = 2.5.",
				},
			},
			Hints: []Hint{
				{
					Title: "Binary Search on Partition",
					Body:  "Partition the smaller array and the larger array such that the left half has the same number of elements as the right half. Use binary search to find the correct partition point.",
				},
			},
			JavascriptSC: `/**
 * @param {number[]} nums1
 * @param {number[]} nums2
 * @return {number}
 */
function findMedianSortedArrays(nums1, nums2) {
    if (nums1.length > nums2.length) {
        return findMedianSortedArrays(nums2, nums1);
    }
    const m = nums1.length;
    const n = nums2.length;
    let left = 0, right = m;
    while (left <= right) {
        const partitionX = Math.floor((left + right) / 2);
        const partitionY = Math.floor((m + n + 1) / 2) - partitionX;
        const maxLeftX = partitionX === 0 ? -Infinity : nums1[partitionX - 1];
        const minRightX = partitionX === m ? Infinity : nums1[partitionX];
        const maxLeftY = partitionY === 0 ? -Infinity : nums2[partitionY - 1];
        const minRightY = partitionY === n ? Infinity : nums2[partitionY];
        if (maxLeftX <= minRightY && maxLeftY <= minRightX) {
            if ((m + n) % 2 === 0) {
                return (Math.max(maxLeftX, maxLeftY) + Math.min(minRightX, minRightY)) / 2;
            } else {
                return Math.max(maxLeftX, maxLeftY);
            }
        } else if (maxLeftX > minRightY) {
            right = partitionX - 1;
        } else {
            left = partitionX + 1;
        }
    }
    return 0.0;
}`,
			PythonSC: `def findMedianSortedArrays(nums1: list[int], nums2: list[int]) -> float:
    if len(nums1) > len(nums2):
        nums1, nums2 = nums2, nums1
    m, n = len(nums1), len(nums2)
    left, right = 0, m
    while left <= right:
        partitionX = (left + right) // 2
        partitionY = (m + n + 1) // 2 - partitionX
        maxLeftX = float('-inf') if partitionX == 0 else nums1[partitionX - 1]
        minRightX = float('inf') if partitionX == m else nums1[partitionX]
        maxLeftY = float('-inf') if partitionY == 0 else nums2[partitionY - 1]
        minRightY = float('inf') if partitionY == n else nums2[partitionY]
        if maxLeftX <= minRightY and maxLeftY <= minRightX:
            if (m + n) % 2 == 0:
                return (max(maxLeftX, maxLeftY) + min(minRightX, minRightY)) / 2.0
            else:
                return float(max(maxLeftX, maxLeftY))
        elif maxLeftX > minRightY:
            right = partitionX - 1
        else:
            left = partitionX + 1
    return 0.0`,
			JavaSC: `public class Solution {
    public double findMedianSortedArrays(int[] nums1, int[] nums2) {
        if (nums1.length > nums2.length) {
            return findMedianSortedArrays(nums2, nums1);
        }
        int m = nums1.length;
        int n = nums2.length;
        int left = 0, right = m;
        while (left <= right) {
            int partitionX = (left + right) / 2;
            int partitionY = (m + n + 1) / 2 - partitionX;
            int maxLeftX = (partitionX == 0) ? Integer.MIN_VALUE : nums1[partitionX - 1];
            int minRightX = (partitionX == m) ? Integer.MAX_VALUE : nums1[partitionX];
            int maxLeftY = (partitionY == 0) ? Integer.MIN_VALUE : nums2[partitionY - 1];
            int minRightY = (partitionY == n) ? Integer.MAX_VALUE : nums2[partitionY];
            if (maxLeftX <= minRightY && maxLeftY <= minRightX) {
                if ((m + n) % 2 == 0) {
                    return ((double)Math.max(maxLeftX, maxLeftY) + Math.min(minRightX, minRightY)) / 2.0;
                } else {
                    return (double)Math.max(maxLeftX, maxLeftY);
                }
            } else if (maxLeftX > minRightY) {
                right = partitionX - 1;
            } else {
                left = partitionX + 1;
            }
        }
        return 0.0;
    }
}`,
			CppSC: `#include <vector>
#include <algorithm>
#include <climits>
using namespace std;

class Solution {
public:
    double findMedianSortedArrays(vector<int>& nums1, vector<int>& nums2) {
        if (nums1.size() > nums2.size()) {
            return findMedianSortedArrays(nums2, nums1);
        }
        int m = nums1.size();
        int n = nums2.size();
        int left = 0, right = m;
        while (left <= right) {
            int partitionX = (left + right) / 2;
            int partitionY = (m + n + 1) / 2 - partitionX;
            int maxLeftX = (partitionX == 0) ? INT_MIN : nums1[partitionX - 1];
            int minRightX = (partitionX == m) ? INT_MAX : nums1[partitionX];
            int maxLeftY = (partitionY == 0) ? INT_MIN : nums2[partitionY - 1];
            int minRightY = (partitionY == n) ? INT_MAX : nums2[partitionY];
            if (maxLeftX <= minRightY && maxLeftY <= minRightX) {
                if ((m + n) % 2 == 0) {
                    return (double(max(maxLeftX, maxLeftY)) + min(minRightX, minRightY)) / 2.0;
                } else {
                    return double(max(maxLeftX, maxLeftY));
                }
            } else if (maxLeftX > minRightY) {
                right = partitionX - 1;
            } else {
                left = partitionX + 1;
            }
        }
        return 0.0;
    }
};`,
			GoSC: `package main

import "math"

func findMedianSortedArrays(nums1 []int, nums2 []int) float64 {
    if len(nums1) > len(nums2) {
        return findMedianSortedArrays(nums2, nums1)
    }
    m := len(nums1)
    n := len(nums2)
    left, right := 0, m
    for left <= right {
        partitionX := (left + right) / 2
        partitionY := (m + n + 1) / 2 - partitionX
        var maxLeftX, minRightX, maxLeftY, minRightY float64
        if partitionX == 0 {
            maxLeftX = math.Inf(-1)
        } else {
            maxLeftX = float64(nums1[partitionX-1])
        }
        if partitionX == m {
            minRightX = math.Inf(1)
        } else {
            minRightX = float64(nums1[partitionX])
        }
        if partitionY == 0 {
            maxLeftY = math.Inf(-1)
        } else {
            maxLeftY = float64(nums2[partitionY-1])
        }
        if partitionY == n {
            minRightY = math.Inf(1)
        } else {
            minRightY = float64(nums2[partitionY])
        }
        if maxLeftX <= minRightY && maxLeftY <= minRightX {
            if (m+n)%2 == 0 {
                return (math.Max(maxLeftX, maxLeftY) + math.Min(minRightX, minRightY)) / 2.0
            }
            return math.Max(maxLeftX, maxLeftY)
        } else if maxLeftX > minRightY {
            right = partitionX - 1
        } else {
            left = partitionX + 1
        }
    }
    return 0.0
}`,
			TestCases: []TestCase{
				{Input: "[1,3]\n[2]", Expected: "2", IsHidden: false},
				{Input: "[1,2]\n[3,4]", Expected: "2.5", IsHidden: false},
				{Input: "[]\n[1]", Expected: "1", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0022",
			Slug:       "trapping-rain-water",
			Title:      "Trapping Rain Water",
			Difficulty: "Hard",
			Topic:      "Array",
			XP:         150,
			Statement:  "Given `n` non-negative integers representing an elevation map where the width of each bar is 1, compute how much water it can trap after raining.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Array", "Two Pointers", "Stack", "Dynamic Programming"},
			Examples: []Example{
				{
					Input:       "height = [0,1,0,2,1,0,1,3,2,1,2,1]",
					Output:      "6",
					Explanation: "The elevation map is represented by array [0,1,0,2,1,0,1,3,2,1,2,1]. In this case, 6 units of rain water are being trapped.",
				},
				{
					Input:       "height = [4,2,0,3,2,5]",
					Output:      "9",
					Explanation: "9 units of water are trapped.",
				},
			},
			Hints: []Hint{
				{
					Title: "Two Pointers Approach",
					Body:  "Keep two pointers: left at index 0 and right at index n-1. Maintain left_max and right_max. Move the pointer pointing to the smaller height, and add the trapped water at that position.",
				},
			},
			JavascriptSC: `/**
 * @param {number[]} height
 * @return {number}
 */
function trap(height) {
    let left = 0, right = height.length - 1;
    let leftMax = 0, rightMax = 0;
    let water = 0;
    while (left < right) {
        if (height[left] < height[right]) {
            if (height[left] >= leftMax) {
                leftMax = height[left];
            } else {
                water += leftMax - height[left];
            }
            left++;
        } else {
            if (height[right] >= rightMax) {
                rightMax = height[right];
            } else {
                water += rightMax - height[right];
            }
            right--;
        }
    }
    return water;
}`,
			PythonSC: `def trap(height: list[int]) -> int:
    left, right = 0, len(height) - 1
    left_max, right_max = 0, 0
    water = 0
    while left < right:
        if height[left] < height[right]:
            if height[left] >= left_max:
                left_max = height[left]
            else:
                water += left_max - height[left]
            left += 1
        else:
            if height[right] >= right_max:
                right_max = height[right]
            else:
                water += right_max - height[right]
            right -= 1
    return water`,
			JavaSC: `public class Solution {
    public int trap(int[] height) {
        int left = 0, right = height.length - 1;
        int leftMax = 0, rightMax = 0;
        int water = 0;
        while (left < right) {
            if (height[left] < height[right]) {
                if (height[left] >= leftMax) {
                    leftMax = height[left];
                } else {
                    water += leftMax - height[left];
                }
                left++;
            } else {
                if (height[right] >= rightMax) {
                    rightMax = height[right];
                } else {
                    water += rightMax - height[right];
                }
                right--;
            }
        }
        return water;
    }
}`,
			CppSC: `#include <vector>
#include <algorithm>
using namespace std;

class Solution {
public:
    int trap(vector<int>& height) {
        int left = 0, right = height.size() - 1;
        int leftMax = 0, rightMax = 0;
        int water = 0;
        while (left < right) {
            if (height[left] < height[right]) {
                if (height[left] >= leftMax) {
                    leftMax = height[left];
                } else {
                    water += leftMax - height[left];
                }
                left++;
            } else {
                if (height[right] >= rightMax) {
                    rightMax = height[right];
                } else {
                    water += rightMax - height[right];
                }
                right--;
            }
        }
        return water;
    }
};`,
			GoSC: `package main

func trap(height []int) int {
    left, right := 0, len(height)-1
    leftMax, rightMax := 0, 0
    water := 0
    for left < right {
        if height[left] < height[right] {
            if height[left] >= leftMax {
                leftMax = height[left]
            } else {
                water += leftMax - height[left]
            }
            left++
        } else {
            if height[right] >= rightMax {
                rightMax = height[right]
            } else {
                water += rightMax - height[right]
            }
            right--
        }
    }
    return water
}`,
			TestCases: []TestCase{
				{Input: "[0,1,0,2,1,0,1,3,2,1,2,1]", Expected: "6", IsHidden: false},
				{Input: "[4,2,0,3,2,5]", Expected: "9", IsHidden: false},
				{Input: "[1,0,2]", Expected: "1", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0023",
			Slug:       "first-missing-positive",
			Title:      "First Missing Positive",
			Difficulty: "Hard",
			Topic:      "Array",
			XP:         150,
			Statement:  "Given an unsorted integer array `nums`, return the smallest missing positive integer.\n\nYou must implement an algorithm that runs in $O(n)$ time and uses $O(1)$ auxiliary space.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Array", "Hash Table"},
			Examples: []Example{
				{
					Input:       "nums = [1,2,0]",
					Output:      "3",
					Explanation: "The numbers in the range [1,2] are all in the array, so the lowest missing positive is 3.",
				},
				{
					Input:       "nums = [3,4,-1,1]",
					Output:      "2",
					Explanation: "1 and 3 are in the array, but 2 is missing.",
				},
			},
			Hints: []Hint{
				{
					Title: "Index as Hash Key",
					Body:  "Rearrange the numbers such that number `i` is placed at index `i - 1`. Then scan the array from left to right, and the first index where `nums[i] != i + 1` represents the missing positive.",
				},
			},
			JavascriptSC: `/**
 * @param {number[]} nums
 * @return {number}
 */
function firstMissingPositive(nums) {
    const n = nums.length;
    for (let i = 0; i < n; i++) {
        while (nums[i] > 0 && nums[i] <= n && nums[nums[i] - 1] !== nums[i]) {
            const temp = nums[nums[i] - 1];
            nums[nums[i] - 1] = nums[i];
            nums[i] = temp;
        }
    }
    for (let i = 0; i < n; i++) {
        if (nums[i] !== i + 1) return i + 1;
    }
    return n + 1;
}`,
			PythonSC: `def firstMissingPositive(nums: list[int]) -> int:
    n = len(nums)
    for i in range(n):
        while 0 < nums[i] <= n and nums[nums[i] - 1] != nums[i]:
            correct_idx = nums[i] - 1
            nums[i], nums[correct_idx] = nums[correct_idx], nums[i]
    for i in range(n):
        if nums[i] != i + 1:
            return i + 1
    return n + 1`,
			JavaSC: `public class Solution {
    public int firstMissingPositive(int[] nums) {
        int n = nums.length;
        for (int i = 0; i < n; i++) {
            while (nums[i] > 0 && nums[i] <= n && nums[nums[i] - 1] != nums[i]) {
                int temp = nums[nums[i] - 1];
                nums[nums[i] - 1] = nums[i];
                nums[i] = temp;
            }
        }
        for (int i = 0; i < n; i++) {
            if (nums[i] != i + 1) return i + 1;
        }
        return n + 1;
    }
}`,
			CppSC: `#include <vector>
#include <algorithm>
using namespace std;

class Solution {
public:
    int firstMissingPositive(vector<int>& nums) {
        int n = nums.size();
        for (int i = 0; i < n; i++) {
            while (nums[i] > 0 && nums[i] <= n && nums[nums[i] - 1] != nums[i]) {
                swap(nums[i], nums[nums[i] - 1]);
            }
        }
        for (int i = 0; i < n; i++) {
            if (nums[i] != i + 1) return i + 1;
        }
        return n + 1;
    }
};`,
			GoSC: `package main

func firstMissingPositive(nums []int) int {
    n := len(nums)
    for i := 0; i < n; i++ {
        for nums[i] > 0 && nums[i] <= n && nums[nums[i]-1] != nums[i] {
            correctIdx := nums[i] - 1
            nums[i], nums[correctIdx] = nums[correctIdx], nums[i]
        }
    }
    for i := 0; i < n; i++ {
        if nums[i] != i+1 {
            return i + 1
        }
    }
    return n + 1
}`,
			TestCases: []TestCase{
				{Input: "[1,2,0]", Expected: "3", IsHidden: false},
				{Input: "[3,4,-1,1]", Expected: "2", IsHidden: false},
				{Input: "[7,8,9,11,12]", Expected: "1", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0024",
			Slug:       "largest-rectangle-in-histogram",
			Title:      "Largest Rectangle in Histogram",
			Difficulty: "Hard",
			Topic:      "Array",
			XP:         150,
			Statement:  "Given an array of integers `heights` representing the histogram's bar height where the width of each bar is 1, return *the area of the largest rectangle in the histogram*.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Array", "Stack"},
			Examples: []Example{
				{
					Input:       "heights = [2,1,5,6,2,3]",
					Output:      "10",
					Explanation: "The largest rectangle has area = 10, formed by bars 5 and 6.",
				},
				{
					Input:       "heights = [2,4]",
					Output:      "4",
					Explanation: "The area is 4, formed by bar at index 1.",
				},
			},
			Hints: []Hint{
				{
					Title: "Monotonic Stack",
					Body:  "Use a monotonic increasing stack to store indices of the bars. For each bar, pop from the stack and compute the area when the current bar is shorter than the bar at the top of the stack.",
				},
			},
			JavascriptSC: `/**
 * @param {number[]} heights
 * @return {number}
 */
function largestRectangleArea(heights) {
    const stack = [];
    let maxArea = 0;
    const n = heights.length;
    for (let i = 0; i <= n; i++) {
        const h = i === n ? 0 : heights[i];
        while (stack.length > 0 && heights[stack[stack.length - 1]] > h) {
            const height = heights[stack.pop()];
            const width = stack.length === 0 ? i : i - stack[stack.length - 1] - 1;
            maxArea = Math.max(maxArea, height * width);
        }
        stack.push(i);
    }
    return maxArea;
}`,
			PythonSC: `def largestRectangleArea(heights: list[int]) -> int:
    stack = []
    max_area = 0
    n = len(heights)
    for i in range(n + 1):
        h = heights[i] if i < n else 0
        while stack and heights[stack[-1]] > h:
            height = heights[stack.pop()]
            width = i if not stack else i - stack[-1] - 1
            max_area = max(max_area, height * width)
        stack.append(i)
    return max_area`,
			JavaSC: `import java.util.*;

public class Solution {
    public int largestRectangleArea(int[] heights) {
        Stack<Integer> stack = new Stack<>();
        int maxArea = 0;
        int n = heights.length;
        for (int i = 0; i <= n; i++) {
            int h = (i == n) ? 0 : heights[i];
            while (!stack.isEmpty() && heights[stack.peek()] > h) {
                int height = heights[stack.pop()];
                int width = stack.isEmpty() ? i : i - stack.peek() - 1;
                maxArea = Math.max(maxArea, height * width);
            }
            stack.push(i);
        }
        return maxArea;
    }
}`,
			CppSC: `#include <vector>
#include <stack>
#include <algorithm>
using namespace std;

class Solution {
public:
    int largestRectangleArea(vector<int>& heights) {
        stack<int> s;
        int maxArea = 0;
        int n = heights.size();
        for (int i = 0; i <= n; i++) {
            int h = (i == n) ? 0 : heights[i];
            while (!s.empty() && heights[s.top()] > h) {
                int height = heights[s.top()];
                s.pop();
                int width = s.empty() ? i : i - s.top() - 1;
                maxArea = max(maxArea, height * width);
            }
            s.push(i);
        }
        return maxArea;
    }
};`,
			GoSC: `package main

func largestRectangleArea(heights []int) int {
    stack := []int{}
    maxArea := 0
    n := len(heights)
    for i := 0; i <= n; i++ {
        h := 0
        if i < n {
            h = heights[i]
        }
        for len(stack) > 0 && heights[stack[len(stack)-1]] > h {
            height := heights[stack[len(stack)-1]]
            stack = stack[:len(stack)-1]
            width := i
            if len(stack) > 0 {
                width = i - stack[len(stack)-1] - 1
            }
            if area := height * width; area > maxArea {
                maxArea = area
            }
        }
        stack = append(stack, i)
    }
    return maxArea
}`,
			TestCases: []TestCase{
				{Input: "[2,1,5,6,2,3]", Expected: "10", IsHidden: false},
				{Input: "[2,4]", Expected: "4", IsHidden: false},
				{Input: "[1,1,1,1]", Expected: "4", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0025",
			Slug:       "sliding-window-maximum",
			Title:      "Sliding Window Maximum",
			Difficulty: "Hard",
			Topic:      "Array",
			XP:         150,
			Statement:  "You are given an array of integers `nums`, there is a sliding window of size `k` which is moving from the very left of the array to the very right. You can only see the `k` numbers in the window. Each time the sliding window moves right by one position.\n\nReturn *the max sliding window*.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Array", "Queue", "Sliding Window"},
			Examples: []Example{
				{
					Input:       "nums = [1,3,-1,-3,5,3,6,7], k = 3",
					Output:      "[3,3,5,5,6,7]",
					Explanation: "Window position:\n[1 3 -1] -3 5 3 6 7 -> max is 3\n1 [3 -1 -3] 5 3 6 7 -> max is 3\n1 3 [-1 -3 5] 3 6 7 -> max is 5\netc.",
				},
			},
			Hints: []Hint{
				{
					Title: "Deque approach",
					Body:  "Use a double-ended queue (deque) to store indices of array elements. Ensure the deque only contains indices inside the current window, and maintain a monotonic decreasing order of values.",
				},
			},
			JavascriptSC: `/**
 * @param {number[]} nums
 * @param {number} k
 * @return {number[]}
 */
function maxSlidingWindow(nums, k) {
    const q = [];
    const res = [];
    for (let i = 0; i < nums.length; i++) {
        while (q.length > 0 && q[0] < i - k + 1) {
            q.shift();
        }
        while (q.length > 0 && nums[q[q.length - 1]] < nums[i]) {
            q.pop();
        }
        q.push(i);
        if (i >= k - 1) {
            res.push(nums[q[0]]);
        }
    }
    return res;
}`,
			PythonSC: `from collections import deque

def maxSlidingWindow(nums: list[int], k: int) -> list[int]:
    q = deque()
    res = []
    for i, num in enumerate(nums):
        while q and q[0] < i - k + 1:
            q.popleft()
        while q and nums[q[-1]] < num:
            q.pop()
        q.append(i)
        if i >= k - 1:
            res.append(nums[q[0]])
    return res`,
			JavaSC: `import java.util.*;

public class Solution {
    public int[] maxSlidingWindow(int[] nums, int k) {
        if (nums == null || nums.length == 0) return new int[0];
        int n = nums.length;
        int[] res = new int[n - k + 1];
        Deque<Integer> q = new ArrayDeque<>();
        for (int i = 0; i < n; i++) {
            while (!q.isEmpty() && q.peek() < i - k + 1) {
                q.poll();
            }
            while (!q.isEmpty() && nums[q.peekLast()] < nums[i]) {
                q.pollLast();
            }
            q.offer(i);
            if (i >= k - 1) {
                res[i - k + 1] = nums[q.peek()];
            }
        }
        return res;
    }
}`,
			CppSC: `#include <vector>
#include <deque>
using namespace std;

class Solution {
public:
    vector<int> maxSlidingWindow(vector<int>& nums, int k) {
        vector<int> res;
        deque<int> q;
        for (int i = 0; i < nums.size(); i++) {
            while (!q.empty() && q.front() < i - k + 1) {
                q.pop_front();
            }
            while (!q.empty() && nums[q.back()] < nums[i]) {
                q.pop_back();
            }
            q.push_back(i);
            if (i >= k - 1) {
                res.push_back(nums[q.front()]);
            }
        }
        return res;
    }
};`,
			GoSC: `package main

func maxSlidingWindow(nums []int, k int) []int {
    if len(nums) == 0 {
        return []int{}
    }
    q := []int{}
    res := []int{}
    for i, num := range nums {
        for len(q) > 0 && q[0] < i-k+1 {
            q = q[1:]
        }
        for len(q) > 0 && nums[q[len(q)-1]] < num {
            q = q[:len(q)-1]
        }
        q = append(q, i)
        if i >= k-1 {
            res = append(res, nums[q[0]])
        }
    }
    return res
}`,
			TestCases: []TestCase{
				{Input: "[1,3,-1,-3,5,3,6,7]\n3", Expected: "[3,3,5,5,6,7]", IsHidden: false},
				{Input: "[1]\n1", Expected: "[1]", IsHidden: false},
				{Input: "[9,11]\n2", Expected: "[11]", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0026",
			Slug:       "best-time-to-buy-and-sell-stock-iii",
			Title:      "Best Time to Buy and Sell Stock III",
			Difficulty: "Hard",
			Topic:      "Array",
			XP:         150,
			Statement:  "You are given an array `prices` where `prices[i]` is the price of a given stock on the `i-th` day.\n\nFind the maximum profit you can achieve. You may complete at most **two transactions**.\n\n*Note*: You may not engage in multiple transactions simultaneously (i.e., you must sell the stock before you buy again).",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Array", "Dynamic Programming"},
			Examples: []Example{
				{
					Input:       "prices = [3,3,5,0,0,3,1,4]",
					Output:      "6",
					Explanation: "Buy on day 4 (price = 0) and sell on day 6 (price = 3), profit = 3-0 = 3. Then buy on day 7 (price = 1) and sell on day 8 (price = 4), profit = 4-1 = 3.",
				},
			},
			Hints: []Hint{
				{
					Title: "Four States DP",
					Body:  "Maintain 4 variables representing: first buy, first sell, second buy, second sell. Update them dynamically as you scan prices.",
				},
			},
			JavascriptSC: `/**
 * @param {number[]} prices
 * @return {number}
 */
function maxProfit(prices) {
    let buy1 = -Infinity, sell1 = 0;
    let buy2 = -Infinity, sell2 = 0;
    for (let i = 0; i < prices.length; i++) {
        buy1 = Math.max(buy1, -prices[i]);
        sell1 = Math.max(sell1, buy1 + prices[i]);
        buy2 = Math.max(buy2, sell1 - prices[i]);
        sell2 = Math.max(sell2, buy2 + prices[i]);
    }
    return sell2;
}`,
			PythonSC: `def maxProfit(prices: list[int]) -> int:
    buy1, sell1 = float('-inf'), 0
    buy2, sell2 = float('-inf'), 0
    for price in prices:
        buy1 = max(buy1, -price)
        sell1 = max(sell1, buy1 + price)
        buy2 = max(buy2, sell1 - price)
        sell2 = max(sell2, buy2 + price)
    return sell2`,
			JavaSC: `public class Solution {
    public int maxProfit(int[] prices) {
        int buy1 = Integer.MIN_VALUE, sell1 = 0;
        int buy2 = Integer.MIN_VALUE, sell2 = 0;
        for (int price : prices) {
            buy1 = Math.max(buy1, -price);
            sell1 = Math.max(sell1, buy1 + price);
            buy2 = Math.max(buy2, sell1 - price);
            sell2 = Math.max(sell2, buy2 + price);
        }
        return sell2;
    }
}`,
			CppSC: `#include <vector>
#include <algorithm>
#include <climits>
using namespace std;

class Solution {
public:
    int maxProfit(vector<int>& prices) {
        int buy1 = INT_MIN, sell1 = 0;
        int buy2 = INT_MIN, sell2 = 0;
        for (int price : prices) {
            buy1 = max(buy1, -price);
            sell1 = max(sell1, buy1 + price);
            buy2 = max(buy2, sell1 - price);
            sell2 = max(sell2, buy2 + price);
        }
        return sell2;
    }
};`,
			GoSC: `package main

import "math"

func maxProfit(prices []int) int {
    buy1, sell1 := math.MinInt32, 0
    buy2, sell2 := math.MinInt32, 0
    for _, price := range prices {
        if -price > buy1 {
            buy1 = -price
        }
        if buy1 + price > sell1 {
            sell1 = buy1 + price
        }
        if sell1 - price > buy2 {
            buy2 = sell1 - price
        }
        if buy2 + price > sell2 {
            sell2 = buy2 + price
        }
    }
    return sell2
}`,
			TestCases: []TestCase{
				{Input: "[3,3,5,0,0,3,1,4]", Expected: "6", IsHidden: false},
				{Input: "[1,2,3,4,5]", Expected: "4", IsHidden: false},
				{Input: "[7,6,4,3,1]", Expected: "0", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0027",
			Slug:       "best-time-to-buy-and-sell-stock-iv",
			Title:      "Best Time to Buy and Sell Stock IV",
			Difficulty: "Hard",
			Topic:      "Array",
			XP:         150,
			Statement:  "You are given an integer `k` and an array of integers `prices` where `prices[i]` is the price of a given stock on the `i-th` day.\n\nFind the maximum profit you can achieve. You may complete at most `k` transactions.\n\n*Note*: You may not engage in multiple transactions simultaneously (i.e., you must sell the stock before you buy again).",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Array", "Dynamic Programming"},
			Examples: []Example{
				{
					Input:       "k = 2, prices = [2,4,1]",
					Output:      "2",
					Explanation: "Buy on day 1 (price = 2) and sell on day 2 (price = 4), profit = 4-2 = 2.",
				},
			},
			Hints: []Hint{
				{
					Title: "DP table",
					Body:  "Maintain buy and sell arrays of size k+1. For each price, update the buy and sell states for transaction count 1 to k.",
				},
			},
			JavascriptSC: `/**
 * @param {number} k
 * @param {number[]} prices
 * @return {number}
 */
function maxProfit(k, prices) {
    if (k === 0 || prices.length === 0) return 0;
    const buy = new Array(k + 1).fill(-Infinity);
    const sell = new Array(k + 1).fill(0);
    for (let i = 0; i < prices.length; i++) {
        for (let j = 1; j <= k; j++) {
            buy[j] = Math.max(buy[j], sell[j - 1] - prices[i]);
            sell[j] = Math.max(sell[j], buy[j] + prices[i]);
        }
    }
    return sell[k];
}`,
			PythonSC: `def maxProfit(k: int, prices: list[int]) -> int:
    if k == 0 or not prices:
        return 0
    buy = [float('-inf')] * (k + 1)
    sell = [0] * (k + 1)
    for price in prices:
        for j in range(1, k + 1):
            buy[j] = max(buy[j], sell[j - 1] - price)
            sell[j] = max(sell[j], buy[j] + price)
    return sell[k]`,
			JavaSC: `import java.util.*;

public class Solution {
    public int maxProfit(int k, int[] prices) {
        if (k == 0 || prices.length == 0) return 0;
        int[] buy = new int[k + 1];
        int[] sell = new int[k + 1];
        Arrays.fill(buy, Integer.MIN_VALUE);
        for (int price : prices) {
            for (int j = 1; j <= k; j++) {
                if (sell[j - 1] != Integer.MIN_VALUE) {
                    buy[j] = Math.max(buy[j], sell[j - 1] - price);
                }
                if (buy[j] != Integer.MIN_VALUE) {
                    sell[j] = Math.max(sell[j], buy[j] + price);
                }
            }
        }
        return sell[k];
    }
}`,
			CppSC: `#include <vector>
#include <algorithm>
#include <climits>
using namespace std;

class Solution {
public:
    int maxProfit(int k, vector<int>& prices) {
        if (k == 0 || prices.empty()) return 0;
        vector<int> buy(k + 1, INT_MIN);
        vector<int> sell(k + 1, 0);
        for (int price : prices) {
            for (int j = 1; j <= k; j++) {
                if (sell[j - 1] != INT_MIN) {
                    buy[j] = max(buy[j], sell[j - 1] - price);
                }
                if (buy[j] != INT_MIN) {
                    sell[j] = max(sell[j], buy[j] + price);
                }
            }
        }
        return sell[k];
    }
};`,
			GoSC: `package main

import "math"

func maxProfit(k int, prices []int) int {
    if k == 0 || len(prices) == 0 {
        return 0
    }
    buy := make([]int, k+1)
    sell := make([]int, k+1)
    for i := range buy {
        buy[i] = math.MinInt32
    }
    for _, price := range prices {
        for j := 1; j <= k; j++ {
            if sell[j-1] - price > buy[j] {
                buy[j] = sell[j-1] - price
            }
            if buy[j] + price > sell[j] {
                sell[j] = buy[j] + price
            }
        }
    }
    return sell[k]
}`,
			TestCases: []TestCase{
				{Input: "2\n[2,4,1]", Expected: "2", IsHidden: false},
				{Input: "2\n[3,2,6,5,0,3]", Expected: "7", IsHidden: false},
				{Input: "1\n[1,2]", Expected: "1", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0028",
			Slug:       "maximum-gap",
			Title:      "Maximum Gap",
			Difficulty: "Hard",
			Topic:      "Array",
			XP:         150,
			Statement:  "Given an integer array `nums`, return *the maximum difference between two successive elements in its sorted form*. If the array contains less than two elements, return `0`.\n\nYou must write an algorithm that runs in linear time and uses linear extra space.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Array", "Sorting"},
			Examples: []Example{
				{
					Input:       "nums = [3,6,9,1]",
					Output:      "3",
					Explanation: "The sorted form of the array is [1,3,6,9], either (3,6) or (6,9) has the maximum difference 3.",
				},
			},
			Hints: []Hint{
				{
					Title: "Bucket Sort Principle (Pigeonhole)",
					Body:  "Use bucket sort. Calculate the minimum and maximum elements. Put numbers in buckets and find the gap between the maximum of one bucket and the minimum of the next.",
				},
			},
			JavascriptSC: `/**
 * @param {number[]} nums
 * @return {number}
 */
function maximumGap(nums) {
    if (nums.length < 2) return 0;
    const minVal = Math.min(...nums);
    const maxVal = Math.max(...nums);
    const n = nums.length;
    const bucketSize = Math.max(1, Math.floor((maxVal - minVal) / (n - 1)));
    const bucketCount = Math.floor((maxVal - minVal) / bucketSize) + 1;
    const bucketsMin = new Array(bucketCount).fill(Infinity);
    const bucketsMax = new Array(bucketCount).fill(-Infinity);
    for (let x of nums) {
        const idx = Math.floor((x - minVal) / bucketSize);
        bucketsMin[idx] = Math.min(bucketsMin[idx], x);
        bucketsMax[idx] = Math.max(bucketsMax[idx], x);
    }
    let maxGap = 0;
    let prev = minVal;
    for (let i = 0; i < bucketCount; i++) {
        if (bucketsMin[i] === Infinity) continue;
        maxGap = Math.max(maxGap, bucketsMin[i] - prev);
        prev = bucketsMax[i];
    }
    return maxGap;
}`,
			PythonSC: `import math

def maximumGap(nums: list[int]) -> int:
    if len(nums) < 2:
        return 0
    min_val, max_val = min(nums), max(nums)
    if min_val == max_val:
        return 0
    n = len(nums)
    bucket_size = max(1, (max_val - min_val) // (n - 1))
    bucket_count = (max_val - min_val) // bucket_size + 1
    buckets_min = [float('inf')] * bucket_count
    buckets_max = [float('-inf')] * bucket_count
    for x in nums:
        idx = (x - min_val) // bucket_size
        buckets_min[idx] = min(buckets_min[idx], x)
        buckets_max[idx] = max(buckets_max[idx], x)
    max_gap = 0
    prev = min_val
    for i in range(bucket_count):
        if buckets_min[i] == float('inf'):
            continue
        max_gap = max(max_gap, buckets_min[i] - prev)
        prev = buckets_max[i]
    return max_gap`,
			JavaSC: `import java.util.*;

public class Solution {
    public int maximumGap(int[] nums) {
        if (nums == null || nums.length < 2) return 0;
        int min = nums[0], max = nums[0];
        for (int x : nums) {
            min = Math.min(min, x);
            max = Math.max(max, x);
        }
        if (min == max) return 0;
        int n = nums.length;
        int bucketSize = Math.max(1, (max - min) / (n - 1));
        int bucketCount = (max - min) / bucketSize + 1;
        int[] bucketsMin = new int[bucketCount];
        int[] bucketsMax = new int[bucketCount];
        Arrays.fill(bucketsMin, Integer.MAX_VALUE);
        Arrays.fill(bucketsMax, Integer.MIN_VALUE);
        for (int x : nums) {
            int idx = (x - min) / bucketSize;
            bucketsMin[idx] = Math.min(bucketsMin[idx], x);
            bucketsMax[idx] = Math.max(bucketsMax[idx], x);
        }
        int maxGap = 0;
        int prev = min;
        for (int i = 0; i < bucketCount; i++) {
            if (bucketsMin[i] == Integer.MAX_VALUE) continue;
            maxGap = Math.max(maxGap, bucketsMin[i] - prev);
            prev = bucketsMax[i];
        }
        return maxGap;
    }
}`,
			CppSC: `#include <vector>
#include <algorithm>
#include <climits>
using namespace std;

class Solution {
public:
    int maximumGap(vector<int>& nums) {
        if (nums.size() < 2) return 0;
        int minVal = nums[0], maxVal = nums[0];
        for (int x : nums) {
            minVal = min(minVal, x);
            maxVal = max(maxVal, x);
        }
        if (minVal == maxVal) return 0;
        int n = nums.size();
        int bucketSize = max(1, (maxVal - minVal) / (n - 1));
        int bucketCount = (maxVal - minVal) / bucketSize + 1;
        vector<int> bucketsMin(bucketCount, INT_MAX);
        vector<int> bucketsMax(bucketCount, INT_MIN);
        for (int x : nums) {
            int idx = (x - minVal) / bucketSize;
            bucketsMin[idx] = min(bucketsMin[idx], x);
            bucketsMax[idx] = max(bucketsMax[idx], x);
        }
        int maxGap = 0;
        int prev = minVal;
        for (int i = 0; i < bucketCount; i++) {
            if (bucketsMin[i] == INT_MAX) continue;
            maxGap = max(maxGap, bucketsMin[i] - prev);
            prev = bucketsMax[i];
        }
        return maxGap;
    }
};`,
			GoSC: `package main

import "math"

func maximumGap(nums []int) int {
    if len(nums) < 2 {
        return 0
    }
    minVal, maxVal := nums[0], nums[0]
    for _, x := range nums {
        if x < minVal { minVal = x }
        if x > maxVal { maxVal = x }
    }
    if minVal == maxVal {
        return 0
    }
    n := len(nums)
    bucketSize := (maxVal - minVal) / (n - 1)
    if bucketSize < 1 {
        bucketSize = 1
    }
    bucketCount := (maxVal - minVal) / bucketSize + 1
    bucketsMin := make([]int, bucketCount)
    bucketsMax := make([]int, bucketCount)
    for i := range bucketsMin {
        bucketsMin[i] = math.MaxInt32
        bucketsMax[i] = math.MinInt32
    }
    for _, x := range nums {
        idx := (x - minVal) / bucketSize
        if x < bucketsMin[idx] { bucketsMin[idx] = x }
        if x > bucketsMax[idx] { bucketsMax[idx] = x }
    }
    maxGap := 0
    prev := minVal
    for i := 0; i < bucketCount; i++ {
        if bucketsMin[i] == math.MaxInt32 {
            continue
        }
        if bucketsMin[i] - prev > maxGap {
            maxGap = bucketsMin[i] - prev
        }
        prev = bucketsMax[i]
    }
    return maxGap
}`,
			TestCases: []TestCase{
				{Input: "[3,6,9,1]", Expected: "3", IsHidden: false},
				{Input: "[10]", Expected: "0", IsHidden: false},
				{Input: "[1,100]", Expected: "99", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0029",
			Slug:       "reverse-pairs",
			Title:      "Reverse Pairs",
			Difficulty: "Hard",
			Topic:      "Array",
			XP:         150,
			Statement:  "Given an integer array `nums`, return *the number of reverse pairs in the array*.\n\nA reverse pair is a pair `(i, j)` where `0 <= i < j < nums.length` and `nums[i] > 2 * nums[j]`.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Array", "Divide and Conquer"},
			Examples: []Example{
				{
					Input:       "nums = [1,3,2,3,1]",
					Output:      "2",
					Explanation: "There are two reverse pairs:\n(1, 4) -> nums[1] = 3 > 2 * nums[4] = 2\n(3, 4) -> nums[3] = 3 > 2 * nums[4] = 2",
				},
			},
			Hints: []Hint{
				{
					Title: "Merge Sort Modification",
					Body:  "Use merge sort. During the merge step, count pairs before merging the two sorted halves. For each element in the left half, count how many elements in the right half satisfy the condition.",
				},
			},
			JavascriptSC: `/**
 * @param {number[]} nums
 * @return {number}
 */
function reversePairs(nums) {
    function mergeSort(l, r) {
        if (l >= r) return 0;
        const mid = Math.floor((l + r) / 2);
        let count = mergeSort(l, mid) + mergeSort(mid + 1, r);
        let j = mid + 1;
        for (let i = l; i <= mid; i++) {
            while (j <= r && nums[i] > 2 * nums[j]) {
                j++;
            }
            count += (j - (mid + 1));
        }
        const temp = [];
        let p1 = l, p2 = mid + 1;
        while (p1 <= mid && p2 <= r) {
            if (nums[p1] <= nums[p2]) {
                temp.push(nums[p1++]);
            } else {
                temp.push(nums[p2++]);
            }
        }
        while (p1 <= mid) temp.push(nums[p1++]);
        while (p2 <= r) temp.push(nums[p2++]);
        for (let i = 0; i < temp.length; i++) {
            nums[l + i] = temp[i];
        }
        return count;
    }
    return mergeSort(0, nums.length - 1);
}`,
			PythonSC: `def reversePairs(nums: list[int]) -> int:
    def mergeSort(l: int, r: int) -> int:
        if l >= r:
            return 0
        mid = (l + r) // 2
        count = mergeSort(l, mid) + mergeSort(mid + 1, r)
        j = mid + 1
        for i in range(l, mid + 1):
            while j <= r and nums[i] > 2 * nums[j]:
                j += 1
            count += j - (mid + 1)
        nums[l:r+1] = sorted(nums[l:r+1])
        return count
    return mergeSort(0, len(nums) - 1)`,
			JavaSC: `import java.util.*;

public class Solution {
    public int reversePairs(int[] nums) {
        return mergeSort(nums, 0, nums.length - 1);
    }
    private int mergeSort(int[] nums, int l, int r) {
        if (l >= r) return 0;
        int mid = l + (r - l) / 2;
        int count = mergeSort(nums, l, mid) + mergeSort(nums, mid + 1, r);
        int j = mid + 1;
        for (int i = l; i <= mid; i++) {
            while (j <= r && (long)nums[i] > 2 * (long)nums[j]) {
                j++;
            }
            count += j - (mid + 1);
        }
        Arrays.sort(nums, l, r + 1);
        return count;
    }
}`,
			CppSC: `#include <vector>
#include <algorithm>
using namespace std;

class Solution {
public:
    int reversePairs(vector<int>& nums) {
        return mergeSort(nums, 0, nums.size() - 1);
    }
private:
    int mergeSort(vector<int>& nums, int l, int r) {
        if (l >= r) return 0;
        int mid = l + (r - l) / 2;
        int count = mergeSort(nums, l, mid) + mergeSort(nums, mid + 1, r);
        int j = mid + 1;
        for (int i = l; i <= mid; i++) {
            while (j <= r && (long long)nums[i] > 2LL * nums[j]) {
                j++;
            }
            count += j - (mid + 1);
        }
        sort(nums.begin() + l, nums.begin() + r + 1);
        return count;
    }
};`,
			GoSC: `package main

import "sort"

func reversePairs(nums []int) int {
    var mergeSort func(l, r int) int
    mergeSort = func(l, r int) int {
        if l >= r {
            return 0
        }
        mid := (l + r) / 2
        count := mergeSort(l, mid) + mergeSort(mid + 1, r)
        j := mid + 1
        for i := l; i <= mid; i++ {
            for j <= r && int64(nums[i]) > 2 * int64(nums[j]) {
                j++
            }
            count += j - (mid + 1)
        }
        sort.Ints(nums[l : r+1])
        return count
    }
    return mergeSort(0, len(nums)-1)
}`,
			TestCases: []TestCase{
				{Input: "[1,3,2,3,1]", Expected: "2", IsHidden: false},
				{Input: "[2,4,3,5,1]", Expected: "3", IsHidden: false},
				{Input: "[5,4,3,2,1]", Expected: "4", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0030",
			Slug:       "maximal-rectangle",
			Title:      "Maximal Rectangle",
			Difficulty: "Hard",
			Topic:      "Array",
			XP:         150,
			Statement:  "Given a `rows x cols` binary `matrix` filled with `0`'s and `1`'s, find the largest rectangle containing only `1`'s and return its area.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Array", "Dynamic Programming", "Stack"},
			Examples: []Example{
				{
					Input:       "matrix = [[\"1\",\"0\",\"1\",\"0\",\"0\"],[\"1\",\"0\",\"1\",\"1\",\"1\"],[\"1\",\"1\",\"1\",\"1\",\"1\"],[\"1\",\"0\",\"0\",\"1\",\"0\"]]",
					Output:      "6",
					Explanation: "The maximal rectangle is shown with area = 6.",
				},
			},
			Hints: []Hint{
				{
					Title: "Histogram Reduction",
					Body:  "Compute the maximum rectangle in histogram for each row, treating each row as a histogram base. Maintain heights of consecutive '1's for each column.",
				},
			},
			JavascriptSC: `/**
 * @param {string[][]} matrix
 * @return {number}
 */
function maximalRectangle(matrix) {
    if (matrix.length === 0) return 0;
    const m = matrix.length;
    const n = matrix[0].length;
    const heights = new Array(n).fill(0);
    let maxArea = 0;
    for (let i = 0; i < m; i++) {
        for (let j = 0; j < n; j++) {
            heights[j] = matrix[i][j] === '1' ? heights[j] + 1 : 0;
        }
        maxArea = Math.max(maxArea, largestRectangleArea(heights));
    }
    return maxArea;
}
function largestRectangleArea(heights) {
    const stack = [];
    let maxArea = 0;
    const n = heights.length;
    for (let i = 0; i <= n; i++) {
        const h = i === n ? 0 : heights[i];
        while (stack.length > 0 && heights[stack[stack.length - 1]] > h) {
            const height = heights[stack.pop()];
            const width = stack.length === 0 ? i : i - stack[stack.length - 1] - 1;
            maxArea = Math.max(maxArea, height * width);
        }
        stack.push(i);
    }
    return maxArea;
}`,
			PythonSC: `def maximalRectangle(matrix: list[list[str]]) -> int:
    if not matrix:
        return 0
    n = len(matrix[0])
    heights = [0] * n
    max_area = 0
    def largestRectangleArea(heights):
        stack = []
        max_area = 0
        for i in range(len(heights) + 1):
            h = heights[i] if i < len(heights) else 0
            while stack and heights[stack[-1]] > h:
                height = heights[stack.pop()]
                width = i if not stack else i - stack[-1] - 1
                max_area = max(max_area, height * width)
            stack.append(i)
        return max_area
    for row in matrix:
        for j in range(n):
            heights[j] = heights[j] + 1 if row[j] == '1' else 0
        max_area = max(max_area, largestRectangleArea(heights))
    return max_area`,
			JavaSC: `import java.util.*;

public class Solution {
    public int maximalRectangle(char[][] matrix) {
        if (matrix == null || matrix.length == 0) return 0;
        int n = matrix[0].length;
        int[] heights = new int[n];
        int maxArea = 0;
        for (char[] row : matrix) {
            for (int j = 0; j < n; j++) {
                heights[j] = (row[j] == '1') ? heights[j] + 1 : 0;
            }
            maxArea = Math.max(maxArea, largestRectangleArea(heights));
        }
        return maxArea;
    }
    private int largestRectangleArea(int[] heights) {
        Stack<Integer> stack = new Stack<>();
        int maxArea = 0;
        int n = heights.length;
        for (int i = 0; i <= n; i++) {
            int h = (i == n) ? 0 : heights[i];
            while (!stack.isEmpty() && heights[stack.peek()] > h) {
                int height = heights[stack.pop()];
                int width = stack.isEmpty() ? i : i - stack.peek() - 1;
                maxArea = Math.max(maxArea, height * width);
            }
            stack.push(i);
        }
        return maxArea;
    }
}`,
			CppSC: `#include <vector>
#include <stack>
#include <algorithm>
using namespace std;

class Solution {
public:
    int maximalRectangle(vector<vector<char>>& matrix) {
        if (matrix.empty()) return 0;
        int n = matrix[0].size();
        vector<int> heights(n, 0);
        int maxArea = 0;
        for (auto& row : matrix) {
            for (int j = 0; j < n; j++) {
                heights[j] = (row[j] == '1') ? heights[j] + 1 : 0;
            }
            maxArea = max(maxArea, largestRectangleArea(heights));
        }
        return maxArea;
    }
private:
    int largestRectangleArea(vector<int>& heights) {
        stack<int> s;
        int maxArea = 0;
        int n = heights.size();
        for (int i = 0; i <= n; i++) {
            int h = (i == n) ? 0 : heights[i];
            while (!s.empty() && heights[s.top()] > h) {
                int height = heights[s.top()];
                s.pop();
                int width = s.empty() ? i : i - s.top() - 1;
                maxArea = max(maxArea, height * width);
            }
            s.push(i);
        }
        return s.push(i), maxArea;
    }
};`,
			GoSC: `package main

func maximalRectangle(matrix [][]byte) int {
    if len(matrix) == 0 {
        return 0
    }
    n := len(matrix[0])
    heights := make([]int, n)
    maxArea := 0
    for _, row := range matrix {
        for j := 0; j < n; j++ {
            if row[j] == '1' {
                heights[j]++
            } else {
                heights[j] = 0
            }
        }
        if area := largestRectangleArea(heights); area > maxArea {
            maxArea = area
        }
    }
    return maxArea
}
func largestRectangleArea(heights []int) int {
    stack := []int{}
    maxArea := 0
    n := len(heights)
    for i := 0; i <= n; i++ {
        h := 0
        if i < n {
            h = heights[i]
        }
        for len(stack) > 0 && heights[stack[len(stack)-1]] > h {
            height := heights[stack[len(stack)-1]]
            stack = stack[:len(stack)-1]
            width := i
            if len(stack) > 0 {
                width = i - stack[len(stack)-1] - 1
            }
            if area := height * width; area > maxArea {
                maxArea = area
            }
        }
        stack = append(stack, i)
    }
    return maxArea
}`,
			TestCases: []TestCase{
				{Input: "[[\"1\",\"0\",\"1\",\"0\",\"0\"],[\"1\",\"0\",\"1\",\"1\",\"1\"],[\"1\",\"1\",\"1\",\"1\",\"1\"],[\"1\",\"0\",\"0\",\"1\",\"0\"]]", Expected: "6", IsHidden: false},
				{Input: "[[\"0\"]]", Expected: "0", IsHidden: false},
				{Input: "[[\"1\"]]", Expected: "1", IsHidden: true},
			},
		},
	}

	problems = append(problems, getHardProblems()...)

	// Dynamically append the other placeholder problems
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

		// Insert problem_tags
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
func getHardProblems() []Problem {
	return []Problem{
		// ─── STRING (2) ───
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0031",
			Slug:       "minimum-window-substring",
			Title:      "Minimum Window Substring",
			Difficulty: "Hard",
			Topic:      "String",
			XP:         150,
			Statement:  "Given two strings `s` and `t` of lengths `m` and `n` respectively, return the minimum window substring of `s` such that every character in `t` (including duplicates) is included in the window. If there is no such substring, return the empty string `\"\"`.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"String", "Sliding Window"},
			Examples: []Example{
				{
					Input:       "s = \"ADOBECODEBANC\", t = \"ABC\"",
					Output:      "\"BANC\"",
					Explanation: "The minimum window substring \"BANC\" includes 'A', 'B', and 'C' from string t.",
				},
			},
			Hints: []Hint{
				{
					Title: "Sliding Window",
					Body:  "Use two pointers to create a window. Expand the right pointer until all characters are found, then contract from the left to minimize the window.",
				},
			},
			JavascriptSC: `function minWindow(s, t) {
    if (!s || !t) return "";
    let map = {};
    for (let char of t) map[char] = (map[char] || 0) + 1;
    let left = 0, right = 0, count = Object.keys(map).length;
    let minLen = Infinity, minStart = 0;
    while (right < s.length) {
        let char = s[right];
        if (map[char] !== undefined) {
            map[char]--;
            if (map[char] === 0) count--;
        }
        right++;
        while (count === 0) {
            if (right - left < minLen) {
                minLen = right - left;
                minStart = left;
            }
            let leftChar = s[left];
            if (map[leftChar] !== undefined) {
                if (map[leftChar] === 0) count++;
                map[leftChar]++;
            }
            left++;
        }
    }
    return minLen === Infinity ? "" : s.substring(minStart, minStart + minLen);
}`,
			PythonSC: `from collections import Counter

def minWindow(s: str, t: str) -> str:
    if not s or not t: return ""
    dict_t = Counter(t)
    required = len(dict_t)
    l, r = 0, 0
    formed = 0
    window_counts = {}
    ans = float("inf"), None, None
    while r < len(s):
        character = s[r]
        window_counts[character] = window_counts.get(character, 0) + 1
        if character in dict_t and window_counts[character] == dict_t[character]:
            formed += 1
        while l <= r and formed == required:
            character = s[l]
            if r - l + 1 < ans[0]:
                ans = (r - l + 1, l, r)
            window_counts[character] -= 1
            if character in dict_t and window_counts[character] < dict_t[character]:
                formed -= 1
            l += 1
        r += 1
    return "" if ans[0] == float("inf") else s[ans[1] : ans[2] + 1]`,
			JavaSC: `import java.util.*;

public class Solution {
    public String minWindow(String s, String t) {
        if (s == null || t == null || s.length() < t.length()) return "";
        Map<Character, Integer> map = new HashMap<>();
        for (char c : t.toCharArray()) map.put(c, map.getOrDefault(c, 0) + 1);
        int left = 0, right = 0, count = map.size();
        int minLen = Integer.MAX_VALUE, minStart = 0;
        while (right < s.length()) {
            char c = s.charAt(right);
            if (map.containsKey(c)) {
                map.put(c, map.get(c) - 1);
                if (map.get(c) == 0) count--;
            }
            right++;
            while (count == 0) {
                if (right - left < minLen) {
                    minLen = right - left;
                    minStart = left;
                }
                char leftChar = s.charAt(left);
                if (map.containsKey(leftChar)) {
                    if (map.get(leftChar) == 0) count++;
                    map.put(leftChar, map.get(leftChar) + 1);
                }
                left++;
            }
        }
        return minLen == Integer.MAX_VALUE ? "" : s.substring(minStart, minStart + minLen);
    }
}`,
			CppSC: `#include <string>
#include <unordered_map>
#include <climits>
using namespace std;

class Solution {
public:
    string minWindow(string s, string t) {
        unordered_map<char, int> map;
        for (char c : t) map[c]++;
        int left = 0, right = 0, count = map.size();
        int minLen = INT_MAX, minStart = 0;
        while (right < s.length()) {
            char c = s[right];
            if (map.count(c)) {
                map[c]--;
                if (map[c] == 0) count--;
            }
            right++;
            while (count == 0) {
                if (right - left < minLen) {
                    minLen = right - left;
                    minStart = left;
                }
                char leftChar = s[left];
                if (map.count(leftChar)) {
                    if (map[leftChar] == 0) count++;
                    map[leftChar]++;
                }
                left++;
            }
        }
        return minLen == INT_MAX ? "" : s.substr(minStart, minLen);
    }
};`,
			GoSC: `package main

func minWindow(s string, t string) string {
    if len(s) == 0 || len(t) == 0 { return "" }
    tFreq := make(map[byte]int)
    for i := 0; i < len(t); i++ { tFreq[t[i]]++ }
    windowFreq := make(map[byte]int)
    required := len(tFreq)
    formed := 0
    left, right := 0, 0
    ans := []int{-1, 0, 0} // length, left, right
    for right < len(s) {
        c := s[right]
        windowFreq[c]++
        if tFreq[c] > 0 && windowFreq[c] == tFreq[c] {
            formed++
        }
        for left <= right && formed == required {
            cLeft := s[left]
            if ans[0] == -1 || right-left+1 < ans[0] {
                ans[0] = right - left + 1
                ans[1] = left
                ans[2] = right
            }
            windowFreq[cLeft]--
            if tFreq[cLeft] > 0 && windowFreq[cLeft] < tFreq[cLeft] {
                formed--
            }
            left++
        }
        right++
    }
    if ans[0] == -1 { return "" }
    return s[ans[1] : ans[2]+1]
}`,
			TestCases: []TestCase{
				{Input: "\"ADOBECODEBANC\"\n\"ABC\"", Expected: "\"BANC\"", IsHidden: false},
				{Input: "\"a\"\n\"a\"", Expected: "\"a\"", IsHidden: false},
				{Input: "\"a\"\n\"aa\"", Expected: "\"\"", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0032",
			Slug:       "integer-to-english-words",
			Title:      "Integer to English Words",
			Difficulty: "Hard",
			Topic:      "String",
			XP:         150,
			Statement:  "Convert a non-negative integer `num` to its English words representation.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"String", "Recursion"},
			Examples: []Example{
				{
					Input:       "num = 123",
					Output:      "\"One Hundred Twenty Three\"",
					Explanation: "Represent 123 in standard English.",
				},
			},
			Hints: []Hint{
				{
					Title: "Chunking",
					Body:  "Divide the number into chunks of three digits (thousands, millions, billions) and solve for each chunk recursively.",
				},
			},
			JavascriptSC: `function numberToWords(num) {
    if (num === 0) return "Zero";
    const ones = ["", "One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen"];
    const tens = ["", "Ten", "Twenty", "Thirty", "Forty", "Fifty", "Sixty", "Seventy", "Eighty", "Ninety"];
    const thousands = ["", "Thousand", "Million", "Billion"];
    function helper(n) {
        if (n === 0) return "";
        if (n < 20) return ones[n] + " ";
        if (n < 100) return tens[Math.floor(n / 10)] + " " + helper(n % 10);
        return ones[Math.floor(n / 100)] + " Hundred " + helper(n % 100);
    }
    let res = "";
    let i = 0;
    while (num > 0) {
        if (num % 1000 !== 0) {
            res = helper(num % 1000) + thousands[i] + " " + res;
        }
        num = Math.floor(num / 1000);
        i++;
    }
    return res.trim().replace(/\s+/g, ' ');
}`,
			PythonSC: `def numberToWords(num: int) -> str:
    if num == 0: return "Zero"
    ones = ["", "One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen"]
    tens = ["", "Ten", "Twenty", "Thirty", "Forty", "Fifty", "Sixty", "Seventy", "Eighty", "Ninety"]
    thousands = ["", "Thousand", "Million", "Billion"]
    def helper(n):
        if n == 0: return ""
        elif n < 20: return ones[n] + " "
        elif n < 100: return tens[n // 10] + " " + helper(n % 10)
        else: return ones[n // 100] + " Hundred " + helper(n % 100)
    res = ""
    i = 0
    while num > 0:
        if num % 1000 != 0:
            res = helper(num % 1000) + thousands[i] + " " + res
        num //= 1000
        i += 1
    return res.strip()`,
			JavaSC: `public class Solution {
    private final String[] ones = {"", "One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen"};
    private final String[] tens = {"", "Ten", "Twenty", "Thirty", "Forty", "Fifty", "Sixty", "Seventy", "Eighty", "Ninety"};
    private final String[] thousands = {"", "Thousand", "Million", "Billion"};
    public String numberToWords(int num) {
        if (num == 0) return "Zero";
        String res = "";
        int i = 0;
        while (num > 0) {
            if (num % 1000 != 0) {
                res = helper(num % 1000) + thousands[i] + " " + res;
            }
            num /= 1000;
            i++;
        }
        return res.trim();
    }
    private String helper(int n) {
        if (n == 0) return "";
        else if (n < 20) return ones[n] + " ";
        else if (n < 100) return tens[n / 10] + " " + helper(n % 10);
        else return ones[n / 100] + " Hundred " + helper(n % 100);
    }
}`,
			CppSC: `#include <string>
#include <vector>
using namespace std;

class Solution {
    vector<string> ones = {"", "One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen"};
    vector<string> tens = {"", "Ten", "Twenty", "Thirty", "Forty", "Fifty", "Sixty", "Seventy", "Eighty", "Ninety"};
    vector<string> thousands = {"", "Thousand", "Million", "Billion"};
public:
    string numberToWords(int num) {
        if (num == 0) return "Zero";
        string res = "";
        int i = 0;
        while (num > 0) {
            if (num % 1000 != 0) {
                res = helper(num % 1000) + thousands[i] + " " + res;
            }
            num /= 1000;
            i++;
        }
        while(!res.empty() && res.back() == ' ') res.pop_back();
        return res;
    }
    string helper(int n) {
        if (n == 0) return "";
        else if (n < 20) return ones[n] + " ";
        else if (n < 100) return tens[n / 10] + " " + helper(n % 10);
        else return ones[n / 100] + " Hundred " + helper(n % 100);
    }
};`,
			GoSC: `package main

import (
	"strings"
)

func numberToWords(num int) string {
	if num == 0 { return "Zero" }
	ones := []string{"", "One", "Two", "Three", "Four", "Five", "Six", "Seven", "Eight", "Nine", "Ten", "Eleven", "Twelve", "Thirteen", "Fourteen", "Fifteen", "Sixteen", "Seventeen", "Eighteen", "Nineteen"}
	tens := []string{"", "Ten", "Twenty", "Thirty", "Forty", "Fifty", "Sixty", "Seventy", "Eighty", "Ninety"}
	thousands := []string{"", "Thousand", "Million", "Billion"}
	var helper func(n int) string
	helper = func(n int) string {
		if n == 0 { return "" }
		if n < 20 { return ones[n] + " " }
		if n < 100 { return tens[n/10] + " " + helper(n%10) }
		return ones[n/100] + " Hundred " + helper(n%100)
	}
	res := ""
	i := 0
	for num > 0 {
		if num%1000 != 0 {
			res = helper(num%1000) + thousands[i] + " " + res
		}
		num /= 1000
		i++
	}
	return strings.TrimSpace(res)
}`,
			TestCases: []TestCase{
				{Input: "123", Expected: "\"One Hundred Twenty Three\"", IsHidden: false},
				{Input: "12345", Expected: "\"Twelve Thousand Three Hundred Forty Five\"", IsHidden: false},
				{Input: "0", Expected: "\"Zero\"", IsHidden: true},
			},
		},

		// ─── HASHMAP (2) ───
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0033",
			Slug:       "max-points-on-a-line",
			Title:      "Max Points on a Line",
			Difficulty: "Hard",
			Topic:      "HashMap",
			XP:         150,
			Statement:  "Given an array of `points` where `points[i] = [xi, yi]` represents a point on the X-Y plane, return the maximum number of points that lie on the same straight line.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"HashMap", "Math"},
			Examples: []Example{
				{
					Input:       "points = [[1,1],[2,2],[3,3]]",
					Output:      "3",
					Explanation: "All three points lie on the line y = x.",
				},
			},
			Hints: []Hint{
				{
					Title: "Slopes",
					Body:  "For each point, calculate the slopes to all other points. Use a hash map to count occurrences of each slope.",
				},
			},
			JavascriptSC: `function maxPoints(points) {
    if (points.length <= 2) return points.length;
    let maxPts = 1;
    for (let i = 0; i < points.length; i++) {
        let map = {};
        let duplicate = 0;
        let localMax = 0;
        for (let j = i + 1; j < points.length; j++) {
            let dx = points[j][0] - points[i][0];
            let dy = points[j][1] - points[i][1];
            if (dx === 0 && dy === 0) {
                duplicate++;
                continue;
            }
            let g = gcd(dx, dy);
            dx /= g;
            dy /= g;
            let slope = dy + "/" + dx;
            map[slope] = (map[slope] || 0) + 1;
            localMax = Math.max(localMax, map[slope]);
        }
        maxPts = Math.max(maxPts, localMax + duplicate + 1);
    }
    return maxPts;
}
function gcd(a, b) {
    return b === 0 ? a : gcd(b, a % b);
}`,
			PythonSC: `from math import gcd

def maxPoints(points: list[list[int]]) -> int:
    if len(points) <= 2: return len(points)
    max_pts = 1
    for i in range(len(points)):
        slopes = {}
        duplicates = 0
        local_max = 0
        for j in range(i + 1, len(points)):
            dx = points[j][0] - points[i][0]
            dy = points[j][1] - points[i][1]
            if dx == 0 and dy == 0:
                duplicates += 1
                continue
            g = gcd(dx, dy)
            dx //= g
            dy //= g
            slope = (dy, dx)
            slopes[slope] = slopes.get(slope, 0) + 1
            local_max = max(local_max, slopes[slope])
        max_pts = max(max_pts, local_max + duplicates + 1)
    return max_pts`,
			JavaSC: `import java.util.*;

public class Solution {
    public int maxPoints(int[][] points) {
        if (points.length <= 2) return points.length;
        int maxPts = 1;
        for (int i = 0; i < points.length; i++) {
            Map<String, Integer> map = new HashMap<>();
            int duplicates = 0, localMax = 0;
            for (int j = i + 1; j < points.length; j++) {
                int dx = points[j][0] - points[i][0];
                int dy = points[j][1] - points[i][1];
                if (dx == 0 && dy == 0) {
                    duplicates++;
                    continue;
                }
                int g = gcd(dx, dy);
                dx /= g;
                dy /= g;
                String slope = dy + "/" + dx;
                map.put(slope, map.getOrDefault(slope, 0) + 1);
                localMax = Math.max(localMax, map.get(slope));
            }
            maxPts = Math.max(maxPts, localMax + duplicates + 1);
        }
        return maxPts;
    }
    private int gcd(int a, int b) {
        return b == 0 ? a : gcd(b, a % b);
    }
}`,
			CppSC: `#include <vector>
#include <unordered_map>
#include <string>
#include <algorithm>
using namespace std;

class Solution {
public:
    int maxPoints(vector<vector<int>>& points) {
        if (points.size() <= 2) return points.size();
        int maxPts = 1;
        for (int i = 0; i < points.size(); i++) {
            unordered_map<string, int> map;
            int duplicates = 0, localMax = 0;
            for (int j = i + 1; j < points.size(); j++) {
                int dx = points[j][0] - points[i][0];
                int dy = points[j][1] - points[i][1];
                if (dx == 0 && dy == 0) {
                    duplicates++;
                    continue;
                }
                int g = gcd(dx, dy);
                dx /= g;
                dy /= g;
                string slope = to_string(dy) + "/" + to_string(dx);
                map[slope]++;
                localMax = max(localMax, map[slope]);
            }
            maxPts = max(maxPts, localMax + duplicates + 1);
        }
        return maxPts;
    }
private:
    int gcd(int a, int b) {
        return b == 0 ? a : gcd(b, a % b);
    }
};`,
			GoSC: `package main

func maxPoints(points [][]int) int {
	if len(points) <= 2 { return len(points) }
	var gcd func(a, b int) int
	gcd = func(a, b int) int {
		if b == 0 { return a }
		return gcd(b, a%b)
	}
	maxPts := 1
	for i := 0; i < len(points); i++ {
		slopes := make(map[[2]int]int)
		duplicates := 0
		localMax := 0
		for j := i + 1; j < len(points); j++ {
			dx := points[j][0] - points[i][0]
			dy := points[j][1] - points[i][1]
			if dx == 0 && dy == 0 {
				duplicates++
				continue
			}
			g := gcd(dx, dy)
			dx /= g
			dy /= g
			slope := [2]int{dy, dx}
			slopes[slope]++
			if slopes[slope] > localMax {
				localMax = slopes[slope]
			}
		}
		if localMax + duplicates + 1 > maxPts {
			maxPts = localMax + duplicates + 1
		}
	}
	return maxPts
}`,
			TestCases: []TestCase{
				{Input: "[[1,1],[2,2],[3,3]]", Expected: "3", IsHidden: false},
				{Input: "[[1,1],[3,2],[5,3],[4,1],[2,3],[1,4]]", Expected: "4", IsHidden: false},
				{Input: "[[0,0]]", Expected: "1", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0034",
			Slug:       "substring-with-concatenation-of-all-words",
			Title:      "Substring with Concatenation of All Words",
			Difficulty: "Hard",
			Topic:      "HashMap",
			XP:         150,
			Statement:  "You are given a string `s` and an array of strings `words`. All the strings of `words` are of the same length.\n\nReturn the starting indices of all the concatenated substrings in `s` that contain all the words in any order without any intervening characters.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"HashMap", "String", "Sliding Window"},
			Examples: []Example{
				{
					Input:       "s = \"barfoothefoobarman\", words = [\"foo\",\"bar\"]",
					Output:      "[0,9]",
					Explanation: "Substrings starting at 0 and 9 are \"barfoo\" and \"foobar\" respectively, which are permutations of words.",
				},
			},
			Hints: []Hint{
				{
					Title: "Word Map",
					Body:  "Count occurrences of each word in a map. Use a sliding window to check matching word counts for window size (length of words * word length).",
				},
			},
			JavascriptSC: `function findSubstring(s, words) {
    if (!s || !words || words.length === 0) return [];
    let wordLen = words[0].length;
    let wordCount = words.length;
    let totalLen = wordLen * wordCount;
    let counts = {};
    for (let w of words) counts[w] = (counts[w] || 0) + 1;
    let res = [];
    for (let i = 0; i <= s.length - totalLen; i++) {
        let sub = s.substring(i, i + totalLen);
        let tempCounts = {...counts};
        let j = 0;
        while (j < wordCount) {
            let word = sub.substring(j * wordLen, (j + 1) * wordLen);
            if (tempCounts[word] === undefined || tempCounts[word] === 0) break;
            tempCounts[word]--;
            j++;
        }
        if (j === wordCount) res.push(i);
    }
    return res;
}`,
			PythonSC: `def findSubstring(s: str, words: list[str]) -> list[int]:
    if not s or not words: return []
    word_len = len(words[0])
    word_count = len(words)
    total_len = word_len * word_count
    counts = {}
    for w in words: counts[w] = counts.get(w, 0) + 1
    res = []
    for i in range(len(s) - total_len + 1):
        temp = {}
        j = 0
        while j < word_count:
            word = s[i + j*word_len : i + (j+1)*word_len]
            if word in counts:
                temp[word] = temp.get(word, 0) + 1
                if temp[word] > counts[word]:
                    break
            else:
                break
            j += 1
        if j == word_count:
            res.append(i)
    return res`,
			JavaSC: `import java.util.*;

public class Solution {
    public List<Integer> findSubstring(String s, String[] words) {
        List<Integer> res = new ArrayList<>();
        if (s == null || words == null || words.length == 0) return res;
        int wordLen = words[0].length();
        int wordCount = words.length;
        int totalLen = wordLen * wordCount;
        Map<String, Integer> counts = new HashMap<>();
        for (String w : words) counts.put(w, counts.getOrDefault(w, 0) + 1);
        for (int i = 0; i <= s.length() - totalLen; i++) {
            Map<String, Integer> temp = new HashMap<>();
            int j = 0;
            while (j < wordCount) {
                String word = s.substring(i + j * wordLen, i + (j + 1) * wordLen);
                if (counts.containsKey(word)) {
                    temp.put(word, temp.getOrDefault(word, 0) + 1);
                    if (temp.get(word) > counts.get(word)) break;
                } else {
                    break;
                }
                j++;
            }
            if (j == wordCount) res.add(i);
        }
        return res;
    }
}`,
			CppSC: `#include <vector>
#include <string>
#include <unordered_map>
using namespace std;

class Solution {
public:
    vector<int> findSubstring(string s, vector<string>& words) {
        vector<int> res;
        if (s.empty() || words.empty()) return res;
        int wordLen = words[0].length();
        int wordCount = words.size();
        int totalLen = wordLen * wordCount;
        if (s.length() < totalLen) return res;
        unordered_map<string, int> counts;
        for (const string& w : words) counts[w]++;
        for (int i = 0; i <= (int)s.length() - totalLen; i++) {
            unordered_map<string, int> temp;
            int j = 0;
            while (j < wordCount) {
                string word = s.substr(i + j * wordLen, wordLen);
                if (counts.count(word)) {
                    temp[word]++;
                    if (temp[word] > counts[word]) break;
                } else {
                    break;
                }
                j++;
            }
            if (j == wordCount) res.push_back(i);
        }
        return res;
    }
};`,
			GoSC: `package main

func findSubstring(s string, words []string) []int {
	res := []int{}
	if len(s) == 0 || len(words) == 0 { return res }
	wordLen := len(words[0])
	wordCount := len(words)
	totalLen := wordLen * wordCount
	if len(s) < totalLen { return res }
	counts := make(map[string]int)
	for _, w := range words { counts[w]++ }
	for i := 0; i <= len(s)-totalLen; i++ {
		temp := make(map[string]int)
		j := 0
		for j < wordCount {
			word := s[i + j*wordLen : i + (j+1)*wordLen]
			if counts[word] > 0 {
				temp[word]++
				if temp[word] > counts[word] { break }
			} else {
				break
			}
			j++
		}
		if j == wordCount {
			res = append(res, i)
		}
	}
	return res
}`,
			TestCases: []TestCase{
				{Input: "\"barfoothefoobarman\"\n[\"foo\",\"bar\"]", Expected: "[0,9]", IsHidden: false},
				{Input: "\"wordgoodgoodgoodbestword\"\n[\"word\",\"good\",\"best\",\"word\"]", Expected: "[]", IsHidden: false},
				{Input: "\"abcdef\"\n[\"ab\",\"cd\",\"ef\"]", Expected: "[0]", IsHidden: true},
			},
		},

		// ─── LINKED LIST (2) ───
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0035",
			Slug:       "merge-k-sorted-lists",
			Title:      "Merge k Sorted Lists",
			Difficulty: "Hard",
			Topic:      "Linked List",
			XP:         150,
			Statement:  "You are given an array of `k` sorted integer lists (represented as arrays). Merge all the sorted lists into one sorted list and return it.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Linked List", "Heap", "Divide and Conquer"},
			Examples: []Example{
				{
					Input:       "lists = [[1,4,5],[1,3,4],[2,6]]",
					Output:      "[1,1,2,3,4,4,5,6]",
					Explanation: "Merging the lists: [1,1,2,3,4,4,5,6].",
				},
			},
			Hints: []Hint{
				{
					Title: "Min-Heap",
					Body:  "Use a min-heap or priority queue to repeatedly select the smallest element among the heads of all k lists.",
				},
			},
			JavascriptSC: `function mergeKLists(lists) {
    let arr = [];
    for (let list of lists) {
        if (list) arr.push(...list);
    }
    return arr.sort((a, b) => a - b);
}`,
			PythonSC: `def mergeKLists(lists: list[list[int]]) -> list[int]:
    res = []
    for l in lists:
        res.extend(l)
    res.sort()
    return res`,
			JavaSC: `import java.util.*;

public class Solution {
    public int[] mergeKLists(int[][] lists) {
        List<Integer> list = new ArrayList<>();
        for (int[] l : lists) {
            for (int val : l) list.add(val);
        }
        Collections.sort(list);
        int[] res = new int[list.size()];
        for (int i = 0; i < list.size(); i++) res[i] = list.get(i);
        return res;
    }
}`,
			CppSC: `#include <vector>
#include <algorithm>
using namespace std;

class Solution {
public:
    vector<int> mergeKLists(vector<vector<int>>& lists) {
        vector<int> res;
        for (auto& l : lists) {
            res.insert(res.end(), l.begin(), l.end());
        }
        sort(res.begin(), res.end());
        return res;
    }
};`,
			GoSC: `package main

import "sort"

func mergeKLists(lists [][]int) []int {
	res := []int{}
	for _, l := range lists {
		res = append(res, l...)
	}
	sort.Ints(res)
	return res
}`,
			TestCases: []TestCase{
				{Input: "[[1,4,5],[1,3,4],[2,6]]", Expected: "[1,1,2,3,4,4,5,6]", IsHidden: false},
				{Input: "[]", Expected: "[]", IsHidden: false},
				{Input: "[[]]", Expected: "[]", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0036",
			Slug:       "reverse-nodes-in-k-group",
			Title:      "Reverse Nodes in k-Group",
			Difficulty: "Hard",
			Topic:      "Linked List",
			XP:         150,
			Statement:  "Given an array of integers `head` representing node values, reverse the nodes of the list `k` at a time, and return the modified array. If the number of nodes is not a multiple of `k` left-out nodes, in the end, should remain as it is.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Linked List", "Recursion"},
			Examples: []Example{
				{
					Input:       "head = [1,2,3,4,5], k = 2",
					Output:      "[2,1,4,3,5]",
					Explanation: "First group [1,2] is reversed to [2,1]. Second group [3,4] is reversed to [4,3]. Last node [5] is left as is.",
				},
			},
			Hints: []Hint{
				{
					Title: "Iterative reversal",
					Body:  "Count nodes first to verify if a group of size k is possible. If yes, perform standard link reversals on k nodes.",
				},
			},
			JavascriptSC: `function reverseKGroup(head, k) {
    let res = [...head];
    let n = res.length;
    for (let i = 0; i <= n - k; i += k) {
        let left = i, right = i + k - 1;
        while (left < right) {
            let temp = res[left];
            res[left] = res[right];
            res[right] = temp;
            left++;
            right--;
        }
    }
    return res;
}`,
			PythonSC: `def reverseKGroup(head: list[int], k: int) -> list[int]:
    res = list(head)
    n = len(res)
    for i in range(0, n - k + 1, k):
        left, right = i, i + k - 1
        while left < right:
            res[left], res[right] = res[right], res[left]
            left += 1
            right -= 1
    return res`,
			JavaSC: `public class Solution {
    public int[] reverseKGroup(int[] head, int k) {
        int[] res = head.clone();
        int n = res.length;
        for (int i = 0; i <= n - k; i += k) {
            int left = i, right = i + k - 1;
            while (left < right) {
                int temp = res[left];
                res[left] = res[right];
                res[right] = temp;
                left++;
                right--;
            }
        }
        return res;
    }
}`,
			CppSC: `#include <vector>
#include <algorithm>
using namespace std;

class Solution {
public:
    vector<int> reverseKGroup(vector<int>& head, int k) {
        vector<int> res = head;
        int n = res.size();
        for (int i = 0; i <= n - k; i += k) {
            reverse(res.begin() + i, res.begin() + i + k);
        }
        return res;
    }
};`,
			GoSC: `package main

func reverseKGroup(head []int, k int) []int {
	res := make([]int, len(head))
	copy(res, head)
	n := len(res)
	for i := 0; i <= n-k; i += k {
		left, right := i, i+k-1
		for left < right {
			res[left], res[right] = res[right], res[left]
			left++
			right--
		}
	}
	return res
}`,
			TestCases: []TestCase{
				{Input: "[1,2,3,4,5]\n2", Expected: "[2,1,4,3,5]", IsHidden: false},
				{Input: "[1,2,3,4,5]\n3", Expected: "[3,2,1,4,5]", IsHidden: false},
				{Input: "[1]\n1", Expected: "[1]", IsHidden: true},
			},
		},

		// ─── TREE (2) ───
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0037",
			Slug:       "binary-tree-maximum-path-sum",
			Title:      "Binary Tree Maximum Path Sum",
			Difficulty: "Hard",
			Topic:      "Tree",
			XP:         150,
			Statement:  "A path in a binary tree is a sequence of nodes where each pair of adjacent nodes in the sequence has an edge connecting them. Given the level-order array `root` representing node values (where `-10000` represents a null node), return the maximum path sum of any non-empty path.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Tree", "DFS", "Dynamic Programming"},
			Examples: []Example{
				{
					Input:       "root = [1,2,3]",
					Output:      "6",
					Explanation: "The optimal path is 2 -> 1 -> 3 with path sum 2 + 1 + 3 = 6.",
				},
			},
			Hints: []Hint{
				{
					Title: "DFS Postorder",
					Body:  "Compute the maximum path sum from left and right children recursively. The max path passing through the current node is node.val + left_gain + right_gain.",
				},
			},
			JavascriptSC: `function maxPathSum(root) {
    let maxPath = -Infinity;
    function dfs(index) {
        if (index >= root.length || root[index] === -10000) return 0;
        let left = Math.max(0, dfs(2 * index + 1));
        let right = Math.max(0, dfs(2 * index + 2));
        maxPath = Math.max(maxPath, root[index] + left + right);
        return root[index] + Math.max(left, right);
    }
    dfs(0);
    return maxPath;
}`,
			PythonSC: `def maxPathSum(root: list[int]) -> int:
    max_path = float('-inf')
    def dfs(index):
        nonlocal max_path
        if index >= len(root) or root[index] == -10000:
            return 0
        left = max(0, dfs(2 * index + 1))
        right = max(0, dfs(2 * index + 2))
        max_path = max(max_path, root[index] + left + right)
        return root[index] + max(left, right)
    dfs(0)
    return max_path`,
			JavaSC: `public class Solution {
    private int maxPath = Integer.MIN_VALUE;
    public int maxPathSum(int[] root) {
        dfs(root, 0);
        return maxPath;
    }
    private int dfs(int[] root, int index) {
        if (index >= root.length || root[index] == -10000) return 0;
        int left = Math.max(0, dfs(root, 2 * index + 1));
        int right = Math.max(0, dfs(root, 2 * index + 2));
        maxPath = Math.max(maxPath, root[index] + left + right);
        return root[index] + Math.max(left, right);
    }
}`,
			CppSC: `#include <vector>
#include <algorithm>
#include <climits>
using namespace std;

class Solution {
    int maxPath = INT_MIN;
public:
    int maxPathSum(vector<int>& root) {
        dfs(root, 0);
        return maxPath;
    }
    int dfs(const vector<int>& root, int index) {
        if (index >= root.size() || root[index] == -10000) return 0;
        int left = max(0, dfs(root, 2 * index + 1));
        int right = max(0, dfs(root, 2 * index + 2));
        maxPath = max(maxPath, root[index] + left + right);
        return root[index] + max(left, right);
    }
};`,
			GoSC: `package main

import "math"

func maxPathSum(root []int) int {
	maxPath := math.MinInt32
	var dfs func(int) int
	dfs = func(index int) int {
		if index >= len(root) || root[index] == -10000 {
			return 0
		}
		left := dfs(2*index + 1)
		if left < 0 { left = 0 }
		right := dfs(2*index + 2)
		if right < 0 { right = 0 }
		if val := root[index] + left + right; val > maxPath {
			maxPath = val
		}
		maxGain := root[index]
		if left > right {
			maxGain += left
		} else {
			maxGain += right
		}
		return maxGain
	}
	dfs(0)
	return maxPath
}`,
			TestCases: []TestCase{
				{Input: "[1,2,3]", Expected: "6", IsHidden: false},
				{Input: "[-10,9,20,-10000,-10000,15,7]", Expected: "42", IsHidden: false},
				{Input: "[5]", Expected: "5", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0038",
			Slug:       "vertical-order-traversal-of-a-binary-tree",
			Title:      "Vertical Order Traversal of a Binary Tree",
			Difficulty: "Hard",
			Topic:      "Tree",
			XP:         150,
			Statement:  "Given the level-order array `root` representing node values (where `-10000` represents a null node), return the vertical order traversal of the binary tree.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Tree", "DFS", "Sorting"},
			Examples: []Example{
				{
					Input:       "root = [3,9,20,-10000,-10000,15,7]",
					Output:      "[[9],[3,15],[20],[7]]",
					Explanation: "Column -1 has [9]. Column 0 has [3, 15] (top-down). Column 1 has [20]. Column 2 has [7].",
				},
			},
			Hints: []Hint{
				{
					Title: "Coordinate system",
					Body:  "Set coordinate (row, col) for each node: root is (0,0). Left child is (row+1, col-1) and right is (row+1, col+1). Group nodes by column and sort by row then node value.",
				},
			},
			JavascriptSC: `function verticalTraversal(root) {
    let nodes = [];
    function dfs(index, row, col) {
        if (index >= root.length || root[index] === -10000) return;
        nodes.push({val: root[index], row, col});
        dfs(2 * index + 1, row + 1, col - 1);
        dfs(2 * index + 2, row + 1, col + 1);
    }
    dfs(0, 0, 0);
    nodes.sort((a, b) => {
        if (a.col !== b.col) return a.col - b.col;
        if (a.row !== b.row) return a.row - b.row;
        return a.val - b.val;
    });
    let res = [];
    let currColVal = -Infinity;
    let currCol = [];
    for (let node of nodes) {
        if (node.col !== currColVal) {
            if (currCol.length > 0) res.push(currCol);
            currColVal = node.col;
            currCol = [node.val];
        } else {
            currCol.push(node.val);
        }
    }
    if (currCol.length > 0) res.push(currCol);
    return res;
}`,
			PythonSC: `def verticalTraversal(root: list[int]) -> list[list[int]]:
    nodes = []
    def dfs(index, row, col):
        if index >= len(root) or root[index] == -10000:
            return
        nodes.append((col, row, root[index]))
        dfs(2 * index + 1, row + 1, col - 1)
        dfs(2 * index + 2, row + 1, col + 1)
    dfs(0, 0, 0)
    nodes.sort()
    res = []
    curr_col = None
    curr_list = []
    for col, row, val in nodes:
        if col != curr_col:
            if curr_list:
                res.append(curr_list)
            curr_col = col
            curr_list = [val]
        else:
            curr_list.append(val)
    if curr_list:
        res.append(curr_list)
    return res`,
			JavaSC: `import java.util.*;

public class Solution {
    private static class Node {
        int val, row, col;
        Node(int val, int row, int col) {
            this.val = val;
            this.row = row;
            this.col = col;
        }
    }
    public List<List<Integer>> verticalTraversal(int[] root) {
        List<Node> list = new ArrayList<>();
        dfs(root, 0, 0, 0, list);
        list.sort((a, b) -> {
            if (a.col != b.col) return Integer.compare(a.col, b.col);
            if (a.row != b.row) return Integer.compare(a.row, b.row);
            return Integer.compare(a.val, b.val);
        });
        List<List<Integer>> res = new ArrayList<>();
        int currColVal = Integer.MIN_VALUE;
        List<Integer> currCol = new ArrayList<>();
        for (Node n : list) {
            if (n.col != currColVal) {
                if (!currCol.isEmpty()) res.add(currCol);
                currColVal = n.col;
                currCol = new ArrayList<>();
                currCol.add(n.val);
            } else {
                currCol.add(n.val);
            }
        }
        if (!currCol.isEmpty()) res.add(currCol);
        return res;
    }
    private void dfs(int[] root, int index, int row, int col, List<Node> list) {
        if (index >= root.length || root[index] == -10000) return;
        list.add(new Node(root[index], row, col));
        dfs(root, 2 * index + 1, row + 1, col - 1, list);
        dfs(root, 2 * index + 2, row + 1, col + 1, list);
    }
}`,
			CppSC: `#include <vector>
#include <algorithm>
using namespace std;

class Solution {
    struct Node {
        int val, row, col;
    };
public:
    vector<vector<int>> verticalTraversal(vector<int>& root) {
        vector<Node> list;
        dfs(root, 0, 0, 0, list);
        sort(list.begin(), list.end(), [](const Node& a, const Node& b) {
            if (a.col != b.col) return a.col < b.col;
            if (a.row != b.row) return a.row < b.row;
            return a.val < b.val;
        });
        vector<vector<int>> res;
        int currColVal = -1e9;
        vector<int> currCol;
        for (const auto& n : list) {
            if (n.col != currColVal) {
                if (!currCol.empty()) res.push_back(currCol);
                currColVal = n.col;
                currCol = {n.val};
            } else {
                currCol.push_back(n.val);
            }
        }
        if (!currCol.empty()) res.push_back(currCol);
        return res;
    }
    void dfs(const vector<int>& root, int index, int row, int col, vector<Node>& list) {
        if (index >= root.size() || root[index] == -10000) return;
        list.push_back({root[index], row, col});
        dfs(root, 2 * index + 1, row + 1, col - 1, list);
        dfs(root, 2 * index + 2, row + 1, col + 1, list);
    }
};`,
			GoSC: `package main

import "sort"

type vNode struct {
	val, row, col int
}

func verticalTraversal(root []int) [][]int {
	list := []vNode{}
	var dfs func(int, int, int)
	dfs = func(index, row, col int) {
		if index >= len(root) || root[index] == -10000 {
			return
		}
		list = append(list, vNode{root[index], row, col})
		dfs(2*index+1, row+1, col-1)
		dfs(2*index+2, row+1, col+1)
	}
	dfs(0, 0, 0)
	sort.Slice(list, func(i, j int) bool {
		if list[i].col != list[j].col { return list[i].col < list[j].col }
		if list[i].row != list[j].row { return list[i].row < list[j].row }
		return list[i].val < list[j].val
	})
	res := [][]int{}
	if len(list) == 0 { return res }
	currColVal := -1000000
	var currCol []int
	for _, n := range list {
		if n.col != currColVal {
			if len(currCol) > 0 {
				res = append(res, currCol)
			}
			currColVal = n.col
			currCol = []int{n.val}
		} else {
			currCol = append(currCol, n.val)
		}
	}
	if len(currCol) > 0 {
		res = append(res, currCol)
	}
	return res
}`,
			TestCases: []TestCase{
				{Input: "[3,9,20,-10000,-10000,15,7]", Expected: "[[9],[3,15],[20],[7]]", IsHidden: false},
				{Input: "[1,2,3,4,5,6,7]", Expected: "[[4],[2],[1,5,6],[3],[7]]", IsHidden: false},
				{Input: "[1]", Expected: "[[1]]", IsHidden: true},
			},
		},

		// ─── GRAPH (2) ───
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0039",
			Slug:       "word-ladder",
			Title:      "Word Ladder",
			Difficulty: "Hard",
			Topic:      "Graph",
			XP:         150,
			Statement:  "Given two words, `beginWord` and `endWord`, and a dictionary `wordList`, return the number of words in the shortest transformation sequence from `beginWord` to `endWord`, or `0` if no such sequence exists.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Graph", "BFS"},
			Examples: []Example{
				{
					Input:       "beginWord = \"hit\", endWord = \"cog\", wordList = [\"hot\",\"dot\",\"dog\",\"lot\",\"log\",\"cog\"]",
					Output:      "5",
					Explanation: "As shortest transformation is \"hit\" -> \"hot\" -> \"dot\" -> \"dog\" -> \"cog\", return its length 5.",
				},
			},
			Hints: []Hint{
				{
					Title: "Breadth-First Search (BFS)",
					Body:  "Treat the words as graph nodes. If two words differ by exactly 1 character, they have an edge. Do a BFS from beginWord to endWord.",
				},
			},
			JavascriptSC: `function ladderLength(beginWord, endWord, wordList) {
    let wordSet = new Set(wordList);
    if (!wordSet.has(endWord)) return 0;
    let queue = [[beginWord, 1]];
    while (queue.length > 0) {
        let [curr, steps] = queue.shift();
        if (curr === endWord) return steps;
        for (let i = 0; i < curr.length; i++) {
            for (let c = 97; c <= 122; c++) {
                let char = String.fromCharCode(c);
                let next = curr.substring(0, i) + char + curr.substring(i + 1);
                if (wordSet.has(next)) {
                    wordSet.delete(next);
                    queue.push([next, steps + 1]);
                }
            }
        }
    }
    return 0;
}`,
			PythonSC: `from collections import deque

def ladderLength(beginWord: str, endWord: str, wordList: list[str]) -> int:
    word_set = set(wordList)
    if endWord not in word_set: return 0
    queue = deque([(beginWord, 1)])
    while queue:
        word, steps = queue.popleft()
        if word == endWord:
            return steps
        for i in range(len(word)):
            for c in "abcdefghijklmnopqrstuvwxyz":
                next_word = word[:i] + c + word[i+1:]
                if next_word in word_set:
                    word_set.remove(next_word)
                    queue.append((next_word, steps + 1))
    return 0`,
			JavaSC: `import java.util.*;

public class Solution {
    public int ladderLength(String beginWord, String endWord, List<String> wordList) {
        Set<String> wordSet = new HashSet<>(wordList);
        if (!wordSet.contains(endWord)) return 0;
        Queue<String> queue = new LinkedList<>();
        queue.add(beginWord);
        int steps = 1;
        while (!queue.isEmpty()) {
            int size = queue.size();
            for (int i = 0; i < size; i++) {
                String curr = queue.poll();
                if (curr.equals(endWord)) return steps;
                char[] chars = curr.toCharArray();
                for (int j = 0; j < chars.length; j++) {
                    char original = chars[j];
                    for (char c = 'a'; c <= 'z'; c++) {
                        chars[j] = c;
                        String next = new String(chars);
                        if (wordSet.contains(next)) {
                            wordSet.remove(next);
                            queue.add(next);
                        }
                    }
                    chars[j] = original;
                }
            }
            steps++;
        }
        return 0;
    }
}`,
			CppSC: `#include <string>
#include <vector>
#include <unordered_set>
#include <queue>
using namespace std;

class Solution {
public:
    int ladderLength(string beginWord, string endWord, vector<string>& wordList) {
        unordered_set<string> wordSet(wordList.begin(), wordList.end());
        if (!wordSet.count(endWord)) return 0;
        queue<pair<string, int>> q;
        q.push({beginWord, 1});
        while (!q.empty()) {
            auto [curr, steps] = q.front();
            q.pop();
            if (curr == endWord) return steps;
            for (int i = 0; i < curr.length(); ++i) {
                char orig = curr[i];
                for (char c = 'a'; c <= 'z'; ++c) {
                    curr[i] = c;
                    if (wordSet.count(curr)) {
                        wordSet.erase(curr);
                        q.push({curr, steps + 1});
                    }
                }
                curr[i] = orig;
            }
        }
        return 0;
    }
};`,
			GoSC: `package main

import "container/list"

func ladderLength(beginWord string, endWord string, wordList []string) int {
	wordSet := make(map[string]bool)
	for _, w := range wordList { wordSet[w] = true }
	if !wordSet[endWord] { return 0 }
	queue := list.New()
	type item struct {
		word  string
		steps int
	}
	queue.PushBack(item{beginWord, 1})
	for queue.Len() > 0 {
		element := queue.Front()
		queue.Remove(element)
		curr := element.Value.(item)
		if curr.word == endWord {
			return curr.steps
		}
		for i := 0; i < len(curr.word); i++ {
			original := curr.word[i]
			wordBytes := []byte(curr.word)
			for c := 'a'; c <= 'z'; c++ {
				wordBytes[i] = byte(c)
				nextWord := string(wordBytes)
				if wordSet[nextWord] {
					delete(wordSet, nextWord)
					queue.PushBack(item{nextWord, curr.steps + 1})
				}
			}
			wordBytes[i] = original
		}
	}
	return 0
}`,
			TestCases: []TestCase{
				{Input: "\"hit\"\n\"cog\"\n[\"hot\",\"dot\",\"dog\",\"lot\",\"log\",\"cog\"]", Expected: "5", IsHidden: false},
				{Input: "\"hit\"\n\"cog\"\n[\"hot\",\"dot\",\"dog\",\"lot\",\"log\"]", Expected: "0", IsHidden: false},
				{Input: "\"a\"\n\"c\"\n[\"a\",\"b\",\"c\"]", Expected: "2", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0040",
			Slug:       "longest-increasing-path-in-a-matrix",
			Title:      "Longest Increasing Path in a Matrix",
			Difficulty: "Hard",
			Topic:      "Graph",
			XP:         150,
			Statement:  "Given an `m x n` integers `matrix`, return the length of the longest increasing path in `matrix`. From each cell, you can move in four directions: left, right, up, down.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Graph", "DFS", "Dynamic Programming"},
			Examples: []Example{
				{
					Input:       "matrix = [[9,9,4],[6,6,8],[2,1,1]]",
					Output:      "4",
					Explanation: "The longest increasing path is [1, 2, 6, 9].",
				},
			},
			Hints: []Hint{
				{
					Title: "DFS with Memoization",
					Body:  "Compute the longest path starting at each cell using DFS. Cache the results in a 2D array to avoid computing the same state twice.",
				},
			},
			JavascriptSC: `function longestIncreasingPath(matrix) {
    if (!matrix || matrix.length === 0) return 0;
    let m = matrix.length, n = matrix[0].length;
    let memo = Array.from({length: m}, () => new Array(n).fill(0));
    let maxPath = 0;
    function dfs(r, c) {
        if (memo[r][c] !== 0) return memo[r][c];
        let dirs = [[0,1], [0,-1], [1,0], [-1,0]];
        let maxLen = 1;
        for (let [dr, dc] of dirs) {
            let nr = r + dr, nc = c + dc;
            if (nr >= 0 && nr < m && nc >= 0 && nc < n && matrix[nr][nc] > matrix[r][c]) {
                maxLen = Math.max(maxLen, 1 + dfs(nr, nc));
            }
        }
        memo[r][c] = maxLen;
        return maxLen;
    }
    for (let i = 0; i < m; i++) {
        for (let j = 0; j < n; j++) {
            maxPath = Math.max(maxPath, dfs(i, j));
        }
    }
    return maxPath;
}`,
			PythonSC: `def longestIncreasingPath(matrix: list[list[int]]) -> int:
    if not matrix: return 0
    m, n = len(matrix), len(matrix[0])
    memo = [[0] * n for _ in range(m)]
    def dfs(r, c):
        if memo[r][c] != 0: return memo[r][c]
        max_len = 1
        for dr, dc in [(0,1), (0,-1), (1,0), (-1,0)]:
            nr, nc = r + dr, c + dc
            if 0 <= nr < m and 0 <= nc < n and matrix[nr][nc] > matrix[r][c]:
                max_len = max(max_len, 1 + dfs(nr, nc))
        memo[r][c] = max_len
        return max_len
    res = 0
    for i in range(m):
        for j in range(n):
            res = max(res, dfs(i, j))
    return res`,
			JavaSC: `public class Solution {
    public int longestIncreasingPath(int[][] matrix) {
        if (matrix == null || matrix.length == 0) return 0;
        int m = matrix.length, n = matrix[0].length;
        int[][] memo = new int[m][n];
        int maxPath = 0;
        for (int i = 0; i < m; i++) {
            for (int j = 0; j < n; j++) {
                maxPath = Math.max(maxPath, dfs(matrix, i, j, memo));
            }
        }
        return maxPath;
    }
    private int dfs(int[][] matrix, int r, int c, int[][] memo) {
        if (memo[r][c] != 0) return memo[r][c];
        int m = matrix.length, n = matrix[0].length;
        int[][] dirs = {{0,1}, {0,-1}, {1,0}, {-1,0}};
        int maxLen = 1;
        for (int[] d : dirs) {
            int nr = r + d[0], nc = c + d[1];
            if (nr >= 0 && nr < m && nc >= 0 && nc < n && matrix[nr][nc] > matrix[r][c]) {
                maxLen = Math.max(maxLen, 1 + dfs(matrix, nr, nc, memo));
            }
        }
        memo[r][c] = maxLen;
        return maxLen;
    }
}`,
			CppSC: `#include <vector>
#include <algorithm>
using namespace std;

class Solution {
public:
    int longestIncreasingPath(vector<vector<int>>& matrix) {
        if (matrix.empty()) return 0;
        int m = matrix.size(), n = matrix[0].size();
        vector<vector<int>> memo(m, vector<int>(n, 0));
        int maxPath = 0;
        for (int i = 0; i < m; ++i) {
            for (int j = 0; j < n; ++j) {
                maxPath = max(maxPath, dfs(matrix, i, j, memo));
            }
        }
        return maxPath;
    }
    int dfs(const vector<vector<int>>& matrix, int r, int c, vector<vector<int>>& memo) {
        if (memo[r][c] != 0) return memo[r][c];
        int m = matrix.size(), n = matrix[0].size();
        int dirs[4][2] = {{0,1}, {0,-1}, {1,0}, {-1,0}};
        int maxLen = 1;
        for (auto& d : dirs) {
            int nr = r + d[0], nc = c + d[1];
            if (nr >= 0 && nr < m && nc >= 0 && nc < n && matrix[nr][nc] > matrix[r][c]) {
                maxLen = max(maxLen, 1 + dfs(matrix, nr, nc, memo));
            }
        }
        memo[r][c] = maxLen;
        return maxLen;
    }
};`,
			GoSC: `package main

func longestIncreasingPath(matrix [][]int) int {
	if len(matrix) == 0 { return 0 }
	m, n := len(matrix), len(matrix[0])
	memo := make([][]int, m)
	for i := range memo { memo[i] = make([]int, n) }
	var dfs func(int, int) int
	dfs = func(r, c int) int {
		if memo[r][c] != 0 { return memo[r][c] }
		dirs := [][]int{{0,1}, {0,-1}, {1,0}, {-1,0}}
		maxLen := 1
		for _, d := range dirs {
			nr, nc := r + d[0], c + d[1]
			if nr >= 0 && nr < m && nc >= 0 && nc < n && matrix[nr][nc] > matrix[r][c] {
				lenVal := 1 + dfs(nr, nc)
				if lenVal > maxLen {
					maxLen = lenVal
				}
			}
		}
		memo[r][c] = maxLen
		return maxLen
	}
	res := 0
	for i := 0; i < m; i++ {
		for j := 0; j < n; j++ {
			if val := dfs(i, j); val > res {
				res = val
			}
		}
	}
	return res
}`,
			TestCases: []TestCase{
				{Input: "[[9,9,4],[6,6,8],[2,1,1]]", Expected: "4", IsHidden: false},
				{Input: "[[3,4,5],[3,2,6],[2,2,1]]", Expected: "4", IsHidden: false},
				{Input: "[[1]]", Expected: "1", IsHidden: true},
			},
		},

		// ─── DP (2) ───
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0041",
			Slug:       "edit-distance",
			Title:      "Edit Distance",
			Difficulty: "Hard",
			Topic:      "DP",
			XP:         150,
			Statement:  "Given two strings `word1` and `word2`, return the minimum number of operations required to convert `word1` to `word2`. You have three operations: Insert, Delete, Replace.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"DP", "String"},
			Examples: []Example{
				{
					Input:       "word1 = \"horse\", word2 = \"ros\"",
					Output:      "3",
					Explanation: "horse -> rorse (replace 'h' with 'r') -> rose (remove 'r') -> ros (remove 'e').",
				},
			},
			Hints: []Hint{
				{
					Title: "2D DP Table",
					Body:  "Let dp[i][j] represent the edit distance of word1[0...i] and word2[0...j]. Find the recurrence relation based on whether word1[i] == word2[j].",
				},
			},
			JavascriptSC: `function minDistance(word1, word2) {
    let m = word1.length, n = word2.length;
    let dp = Array.from({length: m + 1}, () => new Array(n + 1).fill(0));
    for (let i = 0; i <= m; i++) dp[i][0] = i;
    for (let j = 0; j <= n; j++) dp[0][j] = j;
    for (let i = 1; i <= m; i++) {
        for (let j = 1; j <= n; j++) {
            if (word1[i - 1] === word2[j - 1]) {
                dp[i][j] = dp[i - 1][j - 1];
            } else {
                dp[i][j] = Math.min(dp[i - 1][j] + 1, dp[i][j - 1] + 1, dp[i - 1][j - 1] + 1);
            }
        }
    }
    return dp[m][n];
}`,
			PythonSC: `def minDistance(word1: str, word2: str) -> int:
    m, n = len(word1), len(word2)
    dp = [[0] * (n + 1) for _ in range(m + 1)]
    for i in range(m + 1): dp[i][0] = i
    for j in range(n + 1): dp[0][j] = j
    for i in range(1, m + 1):
        for j in range(1, n + 1):
            if word1[i - 1] == word2[j - 1]:
                dp[i][j] = dp[i - 1][j - 1]
            else:
                dp[i][j] = min(dp[i - 1][j] + 1, dp[i][j - 1] + 1, dp[i - 1][j - 1] + 1)
    return dp[m][n]`,
			JavaSC: `public class Solution {
    public int minDistance(String word1, String word2) {
        int m = word1.length(), n = word2.length();
        int[][] dp = new int[m + 1][n + 1];
        for (int i = 0; i <= m; i++) dp[i][0] = i;
        for (int j = 0; j <= n; j++) dp[0][j] = j;
        for (int i = 1; i <= m; i++) {
            for (int j = 1; j <= n; j++) {
                if (word1.charAt(i - 1) == word2.charAt(j - 1)) {
                    dp[i][j] = dp[i - 1][j - 1];
                } else {
                    dp[i][j] = Math.min(dp[i - 1][j] + 1, Math.min(dp[i][j - 1] + 1, dp[i - 1][j - 1] + 1));
                }
            }
        }
        return dp[m][n];
    }
}`,
			CppSC: `#include <string>
#include <vector>
#include <algorithm>
using namespace std;

class Solution {
public:
    int minDistance(string word1, string word2) {
        int m = word1.length(), n = word2.length();
        vector<vector<int>> dp(m + 1, vector<int>(n + 1, 0));
        for (int i = 0; i <= m; ++i) dp[i][0] = i;
        for (int j = 0; j <= n; ++j) dp[0][j] = j;
        for (int i = 1; i <= m; ++i) {
            for (int j = 1; j <= n; ++j) {
                if (word1[i - 1] == word2[j - 1]) {
                    dp[i][j] = dp[i - 1][j - 1];
                } else {
                    dp[i][j] = min(dp[i - 1][j] + 1, min(dp[i][j - 1] + 1, dp[i - 1][j - 1] + 1));
                }
            }
        }
        return dp[m][n];
    }
};`,
			GoSC: `package main

func minDistance(word1 string, word2 string) int {
	m, n := len(word1), len(word2)
	dp := make([][]int, m+1)
	for i := range dp { dp[i] = make([]int, n+1) }
	for i := 0; i <= m; i++ { dp[i][0] = i }
	for j := 0; j <= n; j++ { dp[0][j] = j }
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if word1[i-1] == word2[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				minVal := dp[i-1][j] + 1
				if dp[i][j-1] + 1 < minVal { minVal = dp[i][j-1] + 1 }
				if dp[i-1][j-1] + 1 < minVal { minVal = dp[i-1][j-1] + 1 }
				dp[i][j] = minVal
			}
		}
	}
	return dp[m][n]
}`,
			TestCases: []TestCase{
				{Input: "\"horse\"\n\"ros\"", Expected: "3", IsHidden: false},
				{Input: "\"intention\"\n\"execution\"", Expected: "5", IsHidden: false},
				{Input: "\"\"\n\"\"", Expected: "0", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0042",
			Slug:       "regular-expression-matching",
			Title:      "Regular Expression Matching",
			Difficulty: "Hard",
			Topic:      "DP",
			XP:         150,
			Statement:  "Given an input string `s` and a pattern `p`, implement regular expression matching with support for `'.'` and `'*'`. `'.'` matches any single character. `'*'` matches zero or more of the preceding element.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"DP", "String"},
			Examples: []Example{
				{
					Input:       "s = \"aa\", p = \"a*\"",
					Output:      "true",
					Explanation: "'*' means zero or more of the preceding element 'a'.",
				},
			},
			Hints: []Hint{
				{
					Title: "Dynamic Programming",
					Body:  "Define dp[i][j] representing if s[0...i] matches p[0...j]. Handle '*' separately by either checking zero characters matched or one character matched.",
				},
			},
			JavascriptSC: `function isMatch(s, p) {
    let m = s.length, n = p.length;
    let dp = Array.from({length: m + 1}, () => new Array(n + 1).fill(false));
    dp[0][0] = true;
    for (let j = 1; j <= n; j++) {
        if (p[j - 1] === '*' && dp[0][j - 2]) {
            dp[0][j] = true;
        }
    }
    for (let i = 1; i <= m; i++) {
        for (let j = 1; j <= n; j++) {
            if (p[j - 1] === s[i - 1] || p[j - 1] === '.') {
                dp[i][j] = dp[i - 1][j - 1];
            } else if (p[j - 1] === '*') {
                dp[i][j] = dp[i][j - 2];
                if (p[j - 2] === s[i - 1] || p[j - 2] === '.') {
                    dp[i][j] = dp[i][j] || dp[i - 1][j];
                }
            }
        }
    }
    return dp[m][n];
}`,
			PythonSC: `def isMatch(s: str, p: str) -> bool:
    m, n = len(s), len(p)
    dp = [[False] * (n + 1) for _ in range(m + 1)]
    dp[0][0] = True
    for j in range(1, n + 1):
        if p[j - 1] == '*' and dp[0][j - 2]:
            dp[0][j] = True
    for i in range(1, m + 1):
        for j in range(1, n + 1):
            if p[j - 1] == s[i - 1] or p[j - 1] == '.':
                dp[i][j] = dp[i - 1][j - 1]
            elif p[j - 1] == '*':
                dp[i][j] = dp[i][j - 2]
                if p[j - 2] == s[i - 1] or p[j - 2] == '.':
                    dp[i][j] = dp[i][j] or dp[i - 1][j]
    return dp[m][n]`,
			JavaSC: `public class Solution {
    public boolean isMatch(String s, String p) {
        int m = s.length(), n = p.length();
        boolean[][] dp = new boolean[m + 1][n + 1];
        dp[0][0] = true;
        for (int j = 1; j <= n; j++) {
            if (p.charAt(j - 1) == '*' && dp[0][j - 2]) {
                dp[0][j] = true;
            }
        }
        for (int i = 1; i <= m; i++) {
            for (int j = 1; j <= n; j++) {
                if (p.charAt(j - 1) == s.charAt(i - 1) || p.charAt(j - 1) == '.') {
                    dp[i][j] = dp[i - 1][j - 1];
                } else if (p.charAt(j - 1) == '*') {
                    dp[i][j] = dp[i][j - 2];
                    if (p.charAt(j - 2) == s.charAt(i - 1) || p.charAt(j - 2) == '.') {
                        dp[i][j] = dp[i][j] || dp[i - 1][j];
                    }
                }
            }
        }
        return dp[m][n];
    }
}`,
			CppSC: `#include <string>
#include <vector>
using namespace std;

class Solution {
public:
    bool isMatch(string s, string p) {
        int m = s.length(), n = p.length();
        vector<vector<bool>> dp(m + 1, vector<bool>(n + 1, false));
        dp[0][0] = true;
        for (int j = 1; j <= n; ++j) {
            if (p[j - 1] == '*' && dp[0][j - 2]) {
                dp[0][j] = true;
            }
        }
        for (int i = 1; i <= m; ++i) {
            for (int j = 1; j <= n; ++j) {
                if (p[j - 1] == s[i - 1] || p[j - 1] == '.') {
                    dp[i][j] = dp[i - 1][j - 1];
                } else if (p[j - 1] == '*') {
                    dp[i][j] = dp[i][j - 2];
                    if (p[j - 2] == s[i - 1] || p[j - 2] == '.') {
                        dp[i][j] = dp[i][j] || dp[i - 1][j];
                    }
                }
            }
        }
        return dp[m][n];
    }
};`,
			GoSC: `package main

func isMatch(s string, p string) bool {
	m, n := len(s), len(p)
	dp := make([][]bool, m+1)
	for i := range dp { dp[i] = make([]bool, n+1) }
	dp[0][0] = true
	for j := 1; j <= n; j++ {
		if p[j-1] == '*' && dp[0][j-2] {
			dp[0][j] = true
		}
	}
	for i := 1; i <= m; i++ {
		for j := 1; j <= n; j++ {
			if p[j-1] == s[i-1] || p[j-1] == '.' {
				dp[i][j] = dp[i-1][j-1]
			} else if p[j-1] == '*' {
				dp[i][j] = dp[i][j-2]
				if p[j-2] == s[i-1] || p[j-2] == '.' {
					dp[i][j] = dp[i][j] || dp[i-1][j]
				}
			}
		}
	}
	return dp[m][n]
}`,
			TestCases: []TestCase{
				{Input: "\"aa\"\n\"a*\"", Expected: "true", IsHidden: false},
				{Input: "\"ab\"\n\".*\"", Expected: "true", IsHidden: false},
				{Input: "\"aab\"\n\"c*a*b\"", Expected: "true", IsHidden: true},
			},
		},

		// ─── STACK/QUEUE (2) ───
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0043",
			Slug:       "basic-calculator",
			Title:      "Basic Calculator",
			Difficulty: "Hard",
			Topic:      "Stack/Queue",
			XP:         150,
			Statement:  "Given a string `s` representing a valid expression, implement a basic calculator to evaluate it. Expression contains `(`, `)`, `+`, `-`, non-negative integers and empty spaces.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Stack", "Math"},
			Examples: []Example{
				{
					Input:       "s = \"(1+(4+5+2)-3)+(6+8)\"",
					Output:      "23",
					Explanation: "Evaluating basic arithmetic operations inside parentheses.",
				},
			},
			Hints: []Hint{
				{
					Title: "Stack for Context",
					Body:  "Maintain current sum, sign, and a stack. When you encounter '(', push the current sum and sign to the stack and reset them. When you encounter ')', pop the sign and previous sum to compute values.",
				},
			},
			JavascriptSC: `function calculate(s) {
    let stack = [];
    let sum = 0;
    let sign = 1;
    for (let i = 0; i < s.length; i++) {
        let c = s[i];
        if (c >= '0' && c <= '9') {
            let num = 0;
            while (i < s.length && s[i] >= '0' && s[i] <= '9') {
                num = num * 10 + (s[i].charCodeAt(0) - 48);
                i++;
            }
            i--;
            sum += sign * num;
        } else if (c === '+') {
            sign = 1;
        } else if (c === '-') {
            sign = -1;
        } else if (c === '(') {
            stack.push(sum);
            stack.push(sign);
            sum = 0;
            sign = 1;
        } else if (c === ')') {
            sum = stack.pop() * sum + stack.pop();
        }
    }
    return sum;
}`,
			PythonSC: `def calculate(s: str) -> int:
    stack = []
    sum_val = 0
    sign = 1
    i = 0
    while i < len(s):
        c = s[i]
        if c.isdigit():
            num = 0
            while i < len(s) and s[i].isdigit():
                num = num * 10 + int(s[i])
                i += 1
            i -= 1
            sum_val += sign * num
        elif c == '+':
            sign = 1
        elif c == '-':
            sign = -1
        elif c == '(':
            stack.append(sum_val)
            stack.append(sign)
            sum_val = 0
            sign = 1
        elif c == ')':
            sum_val = stack.pop() * sum_val + stack.pop()
        i += 1
    return sum_val`,
			JavaSC: `import java.util.*;

public class Solution {
    public int calculate(String s) {
        Stack<Integer> stack = new Stack<>();
        int sum = 0, sign = 1;
        for (int i = 0; i < s.length(); i++) {
            char c = s.charAt(i);
            if (Character.isDigit(c)) {
                int num = 0;
                while (i < s.length() && Character.isDigit(s.charAt(i))) {
                    num = num * 10 + (s.charAt(i) - '0');
                    i++;
                }
                i--;
                sum += sign * num;
            } else if (c == '+') {
                sign = 1;
            } else if (c == '-') {
                sign = -1;
            } else if (c == '(') {
                stack.push(sum);
                stack.push(sign);
                sum = 0;
                sign = 1;
            } else if (c == ')') {
                sum = stack.pop() * sum + stack.pop();
            }
        }
        return sum;
    }
}`,
			CppSC: `#include <string>
#include <stack>
using namespace std;

class Solution {
public:
    int calculate(string s) {
        stack<int> st;
        int sum = 0, sign = 1;
        for (int i = 0; i < s.length(); i++) {
            char c = s[i];
            if (c >= '0' && c <= '9') {
                long long num = 0;
                while (i < s.length() && s[i] >= '0' && s[i] <= '9') {
                    num = num * 10 + (s[i] - '0');
                    i++;
                }
                i--;
                sum += sign * num;
            } else if (c == '+') {
                sign = 1;
            } else if (c == '-') {
                sign = -1;
            } else if (c == '(') {
                st.push(sum);
                st.push(sign);
                sum = 0;
                sign = 1;
            } else if (c == ')') {
                sum = st.top() * sum;
                st.pop();
                sum += st.top();
                st.pop();
            }
        }
        return sum;
    }
};`,
			GoSC: `package main

func calculate(s string) int {
	stack := []int{}
	sum, sign := 0, 1
	for i := 0; i < len(s); i++ {
		c := s[i]
		if c >= '0' && c <= '9' {
			num := 0
			for i < len(s) && s[i] >= '0' && s[i] <= '9' {
				num = num*10 + int(s[i]-'0')
				i++
			}
			i--
			sum += sign * num
		} else if c == '+' {
			sign = 1
		} else if c == '-' {
			sign = -1
		} else if c == '(' {
			stack = append(stack, sum)
			stack = append(stack, sign)
			sum = 0
			sign = 1
		} else if c == ')' {
			curSign := stack[len(stack)-1]
			prevSum := stack[len(stack)-2]
			stack = stack[:len(stack)-2]
			sum = curSign * sum + prevSum
		}
	}
	return sum
}`,
			TestCases: []TestCase{
				{Input: "\"(1+(4+5+2)-3)+(6+8)\"", Expected: "23", IsHidden: false},
				{Input: "\" 2-1 + 2 \"", Expected: "3", IsHidden: false},
				{Input: "\"- (3 + (2 - 1))\"", Expected: "-4", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0044",
			Slug:       "shortest-subarray-with-sum-at-least-k",
			Title:      "Shortest Subarray with Sum at Least K",
			Difficulty: "Hard",
			Topic:      "Stack/Queue",
			XP:         150,
			Statement:  "Given an integer array `nums` and an integer `k`, return the length of the shortest non-empty subarray of `nums` with a sum of at least `k`. If there is no such subarray, return `-1`.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Stack", "Deque", "Prefix Sum"},
			Examples: []Example{
				{
					Input:       "nums = [2,-1,2], k = 3",
					Output:      "3",
					Explanation: "The shortest subarray is [2,-1,2] with sum 3.",
				},
			},
			Hints: []Hint{
				{
					Title: "Prefix Sum and Monotonic Queue",
					Body:  "Construct the prefix sum array. Maintain a monotonic increasing queue of prefix sum indices to find the shortest subarray satisfying the condition.",
				},
			},
			JavascriptSC: `function shortestSubarray(nums, k) {
    let n = nums.length;
    let P = new Array(n + 1).fill(0);
    for (let i = 0; i < n; i++) P[i + 1] = P[i] + nums[i];
    let ans = n + 1;
    let monoq = [];
    for (let y = 0; y < P.length; y++) {
        while (monoq.length > 0 && P[y] <= P[monoq[monoq.length - 1]]) {
            monoq.pop();
        }
        while (monoq.length > 0 && P[y] - P[monoq[0]] >= k) {
            ans = Math.min(ans, y - monoq.shift());
        }
        monoq.push(y);
    }
    return ans <= n ? ans : -1;
}`,
			PythonSC: `from collections import deque

def shortestSubarray(nums: list[int], k: int) -> int:
    n = len(nums)
    P = [0] * (n + 1)
    for i in range(n): P[i + 1] = P[i] + nums[i]
    ans = n + 1
    monoq = deque()
    for y, py in enumerate(P):
        while monoq and py <= P[monoq[-1]]:
            monoq.pop()
        while monoq and py - P[monoq[0]] >= k:
            ans = min(ans, y - monoq.popleft())
        monoq.append(y)
    return ans if ans <= n else -1`,
			JavaSC: `import java.util.*;

public class Solution {
    public int shortestSubarray(int[] nums, int k) {
        int n = nums.length;
        long[] P = new long[n + 1];
        for (int i = 0; i < n; i++) P[i + 1] = P[i] + nums[i];
        int ans = n + 1;
        Deque<Integer> monoq = new ArrayDeque<>();
        for (int y = 0; y < P.length; y++) {
            while (!monoq.isEmpty() && P[y] <= P[monoq.peekLast()]) {
                monoq.pollLast();
            }
            while (!monoq.isEmpty() && P[y] - P[monoq.peekFirst()] >= k) {
                ans = Math.min(ans, y - monoq.pollFirst());
            }
            monoq.offerLast(y);
        }
        return ans <= n ? ans : -1;
    }
}`,
			CppSC: `#include <vector>
#include <deque>
#include <algorithm>
using namespace std;

class Solution {
public:
    int shortestSubarray(vector<int>& nums, int k) {
        int n = nums.size();
        vector<long long> P(n + 1, 0);
        for (int i = 0; i < n; ++i) P[i + 1] = P[i] + nums[i];
        int ans = n + 1;
        deque<int> monoq;
        for (int y = 0; y < P.size(); ++y) {
            while (!monoq.empty() && P[y] <= P[monoq.back()]) {
                monoq.pop_back();
            }
            while (!monoq.empty() && P[y] - P[monoq.front()] >= k) {
                ans = min(ans, y - monoq.front());
                monoq.pop_front();
            }
            monoq.push_back(y);
        }
        return ans <= n ? ans : -1;
    }
};`,
			GoSC: `package main

func shortestSubarray(nums []int, k int) int {
	n := len(nums)
	P := make([]int64, n+1)
	for i := 0; i < n; i++ { P[i+1] = P[i] + int64(nums[i]) }
	ans := n + 1
	monoq := []int{}
	for y := 0; y < len(P); y++ {
		for len(monoq) > 0 && P[y] <= P[monoq[len(monoq)-1]] {
			monoq = monoq[:len(monoq)-1]
		}
		for len(monoq) > 0 && P[y] - P[monoq[0]] >= int64(k) {
			diff := y - monoq[0]
			if diff < ans { ans = diff }
			monoq = monoq[1:]
		}
		monoq = append(monoq, y)
	}
	if ans <= n { return ans }
	return -1
}`,
			TestCases: []TestCase{
				{Input: "[2,-1,2]\n3", Expected: "3", IsHidden: false},
				{Input: "[1,2]\n4", Expected: "-1", IsHidden: false},
				{Input: "[1]\n1", Expected: "1", IsHidden: true},
			},
		},

		// ─── HEAP (2) ───
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0045",
			Slug:       "minimum-cost-to-hire-k-workers",
			Title:      "Minimum Cost to Hire K Workers",
			Difficulty: "Hard",
			Topic:      "Heap",
			XP:         150,
			Statement:  "There are `n` workers. Given two integer arrays `quality` and `wage` and an integer `k`, return the least amount of money needed to form a paid group of `k` workers satisfying the hiring conditions.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Heap", "Greedy"},
			Examples: []Example{
				{
					Input:       "quality = [10,20,5], wage = [70,50,30], k = 2",
					Output:      "105.00000",
					Explanation: "We pay 70 to 0-th worker and 35 to 2-nd worker.",
				},
			},
			Hints: []Hint{
				{
					Title: "Wage/Quality Ratio",
					Body:  "Sort workers by their wage/quality ratio. Iterate through workers, maintaining a max-heap of qualities to ensure we select the K workers with the lowest quality in the current ratio bounds.",
				},
			},
			JavascriptSC: `function mincostToHireWorkers(quality, wage, k) {
    let workers = [];
    for (let i = 0; i < quality.length; i++) {
        workers.push({ratio: wage[i] / quality[i], q: quality[i]});
    }
    workers.sort((a, b) => a.ratio - b.ratio);
    let maxHeap = [];
    let sumQ = 0;
    let minCost = Infinity;
    for (let w of workers) {
        maxHeap.push(w.q);
        maxHeap.sort((a, b) => b - a); // Maintain max heap order
        sumQ += w.q;
        if (maxHeap.length > k) {
            sumQ -= maxHeap.shift();
        }
        if (maxHeap.length === k) {
            minCost = Math.min(minCost, sumQ * w.ratio);
        }
    }
    return minCost;
}`,
			PythonSC: `import heapq

def mincostToHireWorkers(quality: list[int], wage: list[int], k: int) -> float:
    workers = sorted([ (w / q, q) for w, q in zip(wage, quality) ])
    pool = []
    sum_q = 0
    ans = float("inf")
    for ratio, q in workers:
        heapq.heappush(pool, -q)
        sum_q += q
        if len(pool) > k:
            sum_q += heapq.heappop(pool)
        if len(pool) == k:
            ans = min(ans, ratio * sum_q)
    return float(ans)`,
			JavaSC: `import java.util.*;

public class Solution {
    public double mincostToHireWorkers(int[] quality, int[] wage, int k) {
        int n = quality.length;
        double[][] workers = new double[n][2];
        for (int i = 0; i < n; i++) {
            workers[i] = new double[]{(double)wage[i] / quality[i], (double)quality[i]};
        }
        Arrays.sort(workers, Comparator.comparingDouble(a -> a[0]));
        PriorityQueue<Double> pool = new PriorityQueue<>(Comparator.reverseOrder());
        double sumQ = 0, ans = Double.MAX_VALUE;
        for (double[] w : workers) {
            pool.offer(w[1]);
            sumQ += w[1];
            if (pool.size() > k) sumQ -= pool.poll();
            if (pool.size() == k) ans = Math.min(ans, w[0] * sumQ);
        }
        return ans;
    }
}`,
			CppSC: `#include <vector>
#include <queue>
#include <algorithm>
using namespace std;

class Solution {
public:
    double mincostToHireWorkers(vector<int>& quality, vector<int>& wage, int k) {
        int n = quality.size();
        vector<pair<double, int>> workers(n);
        for (int i = 0; i < n; ++i) {
            workers[i] = {double(wage[i]) / quality[i], quality[i]};
        }
        sort(workers.begin(), workers.end());
        priority_queue<int> pool;
        int sumQ = 0;
        double ans = 1e18;
        for (auto& w : workers) {
            pool.push(w.second);
            sumQ += w.second;
            if (pool.size() > k) {
                sumQ -= pool.top();
                pool.pop();
            }
            if (pool.size() == k) {
                ans = min(ans, w.first * sumQ);
            }
        }
        return ans;
    }
};`,
			GoSC: `package main

import (
	"container/heap"
	"math"
	"sort"
)

type worker struct {
	ratio float64
	q     float64
}

type maxFloatHeap []float64
func (h maxFloatHeap) Len() int           { return len(h) }
func (h maxFloatHeap) Less(i, j int) bool { return h[i] > h[j] }
func (h maxFloatHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *maxFloatHeap) Push(x interface{}){ *h = append(*h, x.(float64)) }
func (h *maxFloatHeap) Pop() interface{}  {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func mincostToHireWorkers(quality []int, wage []int, k int) float64 {
	n := len(quality)
	workers := make([]worker, n)
	for i := 0; i < n; i++ {
		workers[i] = worker{float64(wage[i]) / float64(quality[i]), float64(quality[i])}
	}
	sort.Slice(workers, func(i, j int) bool { return workers[i].ratio < workers[j].ratio })
	pool := &maxFloatHeap{}
	heap.Init(pool)
	sumQ := 0.0
	ans := math.MaxFloat64
	for _, w := range workers {
		heap.Push(pool, w.q)
		sumQ += w.q
		if pool.Len() > k {
			sumQ -= heap.Pop(pool).(float64)
		}
		if pool.Len() == k {
			if cost := w.ratio * sumQ; cost < ans {
				ans = cost
			}
		}
	}
	return ans
}`,
			TestCases: []TestCase{
				{Input: "[10,20,5]\n[70,50,30]\n2", Expected: "105", IsHidden: false},
				{Input: "[3,1,10,10,1]\n[4,8,2,2,7]\n3", Expected: "30.66667", IsHidden: false},
				{Input: "[1]\n[1]\n1", Expected: "1", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0046",
			Slug:       "smallest-range-covering-elements-from-k-lists",
			Title:      "Smallest Range Covering Elements from K Lists",
			Difficulty: "Hard",
			Topic:      "Heap",
			XP:         150,
			Statement:  "You have `k` lists of sorted integers in non-decreasing order. Find the smallest range `[a, b]` that includes at least one number from each of the `k` lists.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Heap", "Sliding Window"},
			Examples: []Example{
				{
					Input:       "nums = [[4,10,15,24,26],[0,9,12,20],[5,18,22,30]]",
					Output:      "[20,24]",
					Explanation: "List 1: 24, List 2: 20, List 3: 22. Range [20,24] covers elements from all 3 lists.",
				},
			},
			Hints: []Hint{
				{
					Title: "Priority Queue",
					Body:  "Put the first element of each list in a min-heap. Track the maximum element currently in the heap. Repeatedly pop the minimum, calculate the range, and push the next element of that list.",
				},
			},
			JavascriptSC: `function smallestRange(nums) {
    let heap = [];
    let max = -Infinity;
    for (let i = 0; i < nums.length; i++) {
        heap.push({val: nums[i][0], row: i, col: 0});
        max = Math.max(max, nums[i][0]);
    }
    heap.sort((a, b) => a.val - b.val);
    let minRange = Infinity;
    let res = [0, 0];
    while (true) {
        let min = heap.shift();
        if (max - min.val < minRange) {
            minRange = max - min.val;
            res = [min.val, max];
        }
        if (min.col + 1 === nums[min.row].length) break;
        let nextVal = nums[min.row][min.col + 1];
        max = Math.max(max, nextVal);
        heap.push({val: nextVal, row: min.row, col: min.col + 1});
        heap.sort((a, b) => a.val - b.val);
    }
    return res;
}`,
			PythonSC: `import heapq

def smallestRange(nums: list[list[int]]) -> list[int]:
    heap = []
    max_val = float('-inf')
    for i in range(len(nums)):
        heapq.heappush(heap, (nums[i][0], i, 0))
        max_val = max(max_val, nums[i][0])
    min_range = float('inf')
    res = [0, 0]
    while True:
        val, r, c = heapq.heappop(heap)
        if max_val - val < min_range:
            min_range = max_val - val
            res = [val, max_val]
        if c + 1 == len(nums[r]):
            break
        next_val = nums[r][c + 1]
        max_val = max(max_val, next_val)
        heapq.heappush(heap, (next_val, r, c + 1))
    return res`,
			JavaSC: `import java.util.*;

public class Solution {
    private static class Element {
        int val, row, col;
        Element(int val, int row, int col) {
            this.val = val;
            this.row = row;
            this.col = col;
        }
    }
    public int[] smallestRange(int[][] nums) {
        PriorityQueue<Element> minHeap = new PriorityQueue<>(Comparator.comparingInt(a -> a.val));
        int max = Integer.MIN_VALUE;
        for (int i = 0; i < nums.length; i++) {
            minHeap.offer(new Element(nums[i][0], i, 0));
            max = Math.max(max, nums[i][0]);
        }
        int minRange = Integer.MAX_VALUE;
        int[] res = new int[2];
        while (true) {
            Element min = minHeap.poll();
            if (max - min.val < minRange) {
                minRange = max - min.val;
                res = new int[]{min.val, max};
            }
            if (min.col + 1 == nums[min.row].length) break;
            int nextVal = nums[min.row][min.col + 1];
            max = Math.max(max, nextVal);
            minHeap.offer(new Element(nextVal, min.row, min.col + 1));
        }
        return res;
    }
}`,
			CppSC: `#include <vector>
#include <queue>
#include <algorithm>
#include <climits>
using namespace std;

class Solution {
    struct Element {
        int val, row, col;
        bool operator>(const Element& other) const { return val > other.val; }
    };
public:
    vector<int> smallestRange(vector<vector<int>>& nums) {
        priority_queue<Element, vector<Element>, greater<Element>> minHeap;
        int maxVal = INT_MIN;
        for (int i = 0; i < nums.size(); ++i) {
            minHeap.push({nums[i][0], i, 0});
            maxVal = max(maxVal, nums[i][0]);
        }
        int minRange = INT_MAX;
        vector<int> res(2);
        while (true) {
            Element minNode = minHeap.top();
            minHeap.pop();
            if (maxVal - minNode.val < minRange) {
                minRange = maxVal - minNode.val;
                res = {minNode.val, maxVal};
            }
            if (minNode.col + 1 == nums[minNode.row].size()) break;
            int nextVal = nums[minNode.row][minNode.col + 1];
            maxVal = max(maxVal, nextVal);
            minHeap.push({nextVal, minNode.row, minNode.col + 1});
        }
        return res;
    }
};`,
			GoSC: `package main

import (
	"container/heap"
	"math"
)

type hElement struct {
	val, row, col int
}
type minHeap []hElement
func (h minHeap) Len() int           { return len(h) }
func (h minHeap) Less(i, j int) bool { return h[i].val < h[j].val }
func (h minHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *minHeap) Push(x interface{}){ *h = append(*h, x.(hElement)) }
func (h *minHeap) Pop() interface{}  {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func smallestRange(nums [][]int) []int {
	mh := &minHeap{}
	heap.Init(mh)
	maxVal := math.MinInt32
	for i := 0; i < len(nums); i++ {
		heap.Push(mh, hElement{nums[i][0], i, 0})
		if nums[i][0] > maxVal { maxVal = nums[i][0] }
	}
	minRange := math.MaxInt32
	res := []int{0, 0}
	for {
		minNode := heap.Pop(mh).(hElement)
		if maxVal - minNode.val < minRange {
			minRange = maxVal - minNode.val
			res = []int{minNode.val, maxVal}
		}
		if minNode.col + 1 == len(nums[minNode.row]) {
			break
		}
		nextVal := nums[minNode.row][minNode.col+1]
		if nextVal > maxVal { maxVal = nextVal }
		heap.Push(mh, hElement{nextVal, minNode.row, minNode.col + 1})
	}
	return res
}`,
			TestCases: []TestCase{
				{Input: "[[4,10,15,24,26],[0,9,12,20],[5,18,22,30]]", Expected: "[20,24]", IsHidden: false},
				{Input: "[[1,2,3],[1,2,3],[1,2,3]]", Expected: "[1,1]", IsHidden: false},
				{Input: "[[10],[11]]", Expected: "[10,11]", IsHidden: true},
			},
		},

		// ─── BACKTRACKING (2) ───
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0047",
			Slug:       "n-queens-ii",
			Title:      "N-Queens II",
			Difficulty: "Hard",
			Topic:      "Backtracking",
			XP:         150,
			Statement:  "The n-queens puzzle is the problem of placing `n` queens on an `n x n` chessboard such that no two queens attack each other. Given an integer `n`, return *the number of distinct solutions*.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Backtracking"},
			Examples: []Example{
				{
					Input:       "n = 4",
					Output:      "2",
					Explanation: "There are two distinct solutions to the 4-queens puzzle.",
				},
			},
			Hints: []Hint{
				{
					Title: "Backtracking with bitsets",
					Body:  "Use columns, diagonals (row - col), and anti-diagonals (row + col) lookups to check if a queen can be safely placed in a column.",
				},
			},
			JavascriptSC: `function totalNQueens(n) {
    let count = 0;
    let cols = new Set();
    let diag1 = new Set();
    let diag2 = new Set();
    function backtrack(row) {
        if (row === n) {
            count++;
            return;
        }
        for (let col = 0; col < n; col++) {
            if (cols.has(col) || diag1.has(row - col) || diag2.has(row + col)) continue;
            cols.add(col);
            diag1.add(row - col);
            diag2.add(row + col);
            backtrack(row + 1);
            cols.delete(col);
            diag1.delete(row - col);
            diag2.delete(row + col);
        }
    }
    backtrack(0);
    return count;
}`,
			PythonSC: `def totalNQueens(n: int) -> int:
    cols = set()
    diag1 = set()
    diag2 = set()
    count = 0
    def backtrack(row):
        nonlocal count
        if row == n:
            count += 1
            return
        for col in range(n):
            if col in cols or (row - col) in diag1 or (row + col) in diag2:
                continue
            cols.add(col)
            diag1.add(row - col)
            diag2.add(row + col)
            backtrack(row + 1)
            cols.remove(col)
            diag1.remove(row - col)
            diag2.remove(row + col)
    backtrack(0)
    return count`,
			JavaSC: `import java.util.*;

public class Solution {
    private int count = 0;
    private final Set<Integer> cols = new HashSet<>();
    private final Set<Integer> diag1 = new HashSet<>();
    private final Set<Integer> diag2 = new HashSet<>();
    public int totalNQueens(int n) {
        backtrack(0, n);
        return count;
    }
    private void backtrack(int row, int n) {
        if (row == n) {
            count++;
            return;
        }
        for (int col = 0; col < n; col++) {
            if (cols.contains(col) || diag1.contains(row - col) || diag2.contains(row + col)) continue;
            cols.add(col);
            diag1.add(row - col);
            diag2.add(row + col);
            backtrack(row + 1, n);
            cols.remove(col);
            diag1.remove(row - col);
            diag2.remove(row + col);
        }
    }
}`,
			CppSC: `#include <vector>
#include <unordered_set>
using namespace std;

class Solution {
    int count = 0;
    unordered_set<int> cols;
    unordered_set<int> diag1;
    unordered_set<int> diag2;
public:
    int totalNQueens(int n) {
        backtrack(0, n);
        return count;
    }
    void backtrack(int row, int n) {
        if (row == n) {
            count++;
            return;
        }
        for (int col = 0; col < n; ++col) {
            if (cols.count(col) || diag1.count(row - col) || diag2.count(row + col)) continue;
            cols.insert(col);
            diag1.insert(row - col);
            diag2.insert(row + col);
            backtrack(row + 1, n);
            cols.erase(col);
            diag1.erase(row - col);
            diag2.erase(row + col);
        }
    }
};`,
			GoSC: `package main

func totalNQueens(n int) int {
	cols := make(map[int]bool)
	diag1 := make(map[int]bool)
	diag2 := make(map[int]bool)
	count := 0
	var backtrack func(int)
	backtrack = func(row int) {
		if row == n {
			count++
			return
		}
		for col := 0; col < n; col++ {
			if cols[col] || diag1[row-col] || diag2[row+col] { continue }
			cols[col] = true
			diag1[row-col] = true
			diag2[row+col] = true
			backtrack(row + 1)
			delete(cols, col)
			delete(diag1, row-col)
			delete(diag2, row+col)
		}
	}
	backtrack(0)
	return count
}`,
			TestCases: []TestCase{
				{Input: "4", Expected: "2", IsHidden: false},
				{Input: "1", Expected: "1", IsHidden: false},
				{Input: "8", Expected: "92", IsHidden: true},
			},
		},
		{
			ID:         "54574a34-9a68-4e65-ab9a-af05db4d0048",
			Slug:       "sudoku-solver",
			Title:      "Sudoku Solver",
			Difficulty: "Hard",
			Topic:      "Backtracking",
			XP:         150,
			Statement:  "Write a program to solve a Sudoku puzzle by filling the empty cells. Input is a 9x9 board represented as a list of strings where `'.'` represents an empty cell. Return the solved board as a list of strings.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca003",
			Tags:       []string{"Backtracking"},
			Examples: []Example{
				{
					Input:       "board = [[\"5\",\"3\",\".\",\".\",\"7\",\".\",\".\",\".\",\".\"],...]",
					Output:      "[[\"5\",\"3\",\"4\",\"6\",\"7\",\"8\",\"9\",\"1\",\"2\"],...]",
					Explanation: "Outputs the completed sudoku board.",
				},
			},
			Hints: []Hint{
				{
					Title: "Row/Col/Box Constraint Checks",
					Body:  "For each cell, try numbers from '1' to '9'. If the number is valid (doesn't conflict in row, column, or 3x3 box), proceed recursively. Backtrack if it fails.",
				},
			},
			JavascriptSC: `function solveSudoku(board) {
    let charBoard = board.map(row => row.split(''));
    solve(charBoard);
    return charBoard.map(row => row.join(''));
}
function solve(board) {
    for (let i = 0; i < 9; i++) {
        for (let j = 0; j < 9; j++) {
            if (board[i][j] === '.') {
                for (let c = 1; c <= 9; c++) {
                    let char = c.toString();
                    if (isValid(board, i, j, char)) {
                        board[i][j] = char;
                        if (solve(board)) return true;
                        board[i][j] = '.';
                    }
                }
                return false;
            }
        }
    }
    return true;
}
function isValid(board, row, col, c) {
    for (let i = 0; i < 9; i++) {
        if (board[i][col] === c) return false;
        if (board[row][i] === c) return false;
        let boxRow = 3 * Math.floor(row / 3) + Math.floor(i / 3);
        let boxCol = 3 * Math.floor(col / 3) + i % 3;
        if (board[boxRow][boxCol] === c) return false;
    }
    return true;
}`,
			PythonSC: `def solveSudoku(board: list[str]) -> list[str]:
    char_board = [list(row) for row in board]
    def solve(b):
        for i in range(9):
            for j in range(9):
                if b[i][j] == '.':
                    for c in "123456789":
                        if isValid(b, i, j, c):
                            b[i][j] = c
                            if solve(b): return True
                            b[i][j] = '.'
                    return False
        return True
    def isValid(b, r, c, val):
        for i in range(9):
            if b[i][c] == val: return False
            if b[r][i] == val: return False
            if b[3*(r//3) + i//3][3*(c//3) + i%3] == val: return False
        return True
    solve(char_board)
    return ["".join(row) for row in char_board]`,
			JavaSC: `public class Solution {
    public String[] solveSudoku(String[] board) {
        char[][] charBoard = new char[9][9];
        for (int i = 0; i < 9; i++) charBoard[i] = board[i].toCharArray();
        solve(charBoard);
        String[] res = new String[9];
        for (int i = 0; i < 9; i++) res[i] = new String(charBoard[i]);
        return res;
    }
    private boolean solve(char[][] board) {
        for (int i = 0; i < 9; i++) {
            for (int j = 0; j < 9; j++) {
                if (board[i][j] == '.') {
                    for (char c = '1'; c <= '9'; c++) {
                        if (isValid(board, i, j, c)) {
                            board[i][j] = c;
                            if (solve(board)) return true;
                            board[i][j] = '.';
                        }
                    }
                    return false;
                }
            }
        }
        return true;
    }
    private boolean isValid(char[][] board, int row, int col, char c) {
        for (int i = 0; i < 9; i++) {
            if (board[i][col] == c) return false;
            if (board[row][i] == c) return false;
            if (board[3 * (row / 3) + i / 3][3 * (col / 3) + i % 3] == c) return false;
        }
        return true;
    }
}`,
			CppSC: `#include <vector>
#include <string>
using namespace std;

class Solution {
public:
    vector<string> solveSudoku(vector<string>& board) {
        vector<vector<char>> charBoard(9, vector<char>(9));
        for (int i = 0; i < 9; i++) {
            for (int j = 0; j < 9; j++) charBoard[i][j] = board[i][j];
        }
        solve(charBoard);
        vector<string> res(9, ".........");
        for (int i = 0; i < 9; i++) {
            for (int j = 0; j < 9; j++) res[i][j] = charBoard[i][j];
        }
        return res;
    }
    bool solve(vector<vector<char>>& board) {
        for (int i = 0; i < 9; i++) {
            for (int j = 0; j < 9; j++) {
                if (board[i][j] == '.') {
                    for (char c = '1'; c <= '9'; c++) {
                        if (isValid(board, i, j, c)) {
                            board[i][j] = c;
                            if (solve(board)) return true;
                            board[i][j] = '.';
                        }
                    }
                    return false;
                }
            }
        }
        return true;
    }
    bool isValid(const vector<vector<char>>& board, int row, int col, char c) {
        for (int i = 0; i < 9; i++) {
            if (board[i][col] == c) return false;
            if (board[row][i] == c) return false;
            if (board[3 * (row / 3) + i / 3][3 * (col / 3) + i % 3] == c) return false;
        }
        return true;
    }
};`,
			GoSC: `package main

func solveSudoku(board []string) []string {
	charBoard := make([][]byte, 9)
	for i := 0; i < 9; i++ {
		charBoard[i] = []byte(board[i])
	}
	var solve func() bool
	var isValid func(row, col int, c byte) bool
	isValid = func(row, col int, c byte) bool {
		for i := 0; i < 9; i++ {
			if charBoard[i][col] == c { return false }
			if charBoard[row][i] == c { return false }
			if charBoard[3*(row/3)+i/3][3*(col/3)+i%3] == c { return false }
		}
		return true
	}
	solve = func() bool {
		for i := 0; i < 9; i++ {
			for j := 0; j < 9; j++ {
				if charBoard[i][j] == '.' {
					for c := byte('1'); c <= byte('9'); c++ {
						if isValid(i, j, c) {
							charBoard[i][j] = c
							if solve() { return true }
							charBoard[i][j] = '.'
						}
					}
					return false
				}
			}
		}
		return true
	}
	solve()
	res := make([]string, 9)
	for i := 0; i < 9; i++ {
		res[i] = string(charBoard[i])
	}
	return res
}`,
			TestCases: []TestCase{
				{Input: "[\"53..7....\",\"6..195...\",\".98....6.\",\"8...6...3\",\"4..8.3..1\",\"7...2...6\",\".6....28.\",\"...419..5\",\"....8..79\"]", Expected: "[\"534678912\",\"672195348\",\"198342567\",\"859761423\",\"426853791\",\"713924856\",\"961537284\",\"287419635\",\"345286179\"]", IsHidden: false},
			},
		},
	}
}
