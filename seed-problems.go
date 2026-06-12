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
	problems = append(problems, getProficiencyProblems()...)

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
			setID = "54574a34-9a68-4e65-ab9a-af05db4ca001"
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


func getProficiencyProblems() []Problem {
	type problemMeta struct {
		ID          string
		Slug        string
		Title       string
		Topic       string
		Statement   string
		ExampleIn   string
		ExampleOut  string
		HintTitle   string
		HintBody    string
		FuncName    string
		ParamsJS    string
		ParamsPy    string
		ParamsJava  string
		ParamsCpp   string
		ParamsGo    string
		RetPy       string
		RetJava     string
		RetCpp      string
		RetGo       string
		SolJS       string
		SolPy       string
		SolJava     string
		SolCpp      string
		SolGo       string
	}

	metas := []problemMeta{
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0100",
			Slug:        "add-two-numbers",
			Title:       "Add Two Numbers",
			Topic:       "Linked List",
			Statement:   "You are given two non-empty linked lists representing two non-negative integers. The digits are stored in reverse order, and each of their nodes contains a single digit. Add the two numbers and return the sum as a linked list.",
			ExampleIn:   "l1 = [2,4,3], l2 = [5,6,4]",
			ExampleOut:  "[7,0,8]",
			HintTitle:   "Simulation",
			HintBody:    "Track carry values while iterating through the lists.",
			FuncName:    "addTwoNumbers",
			ParamsJS:    "l1, l2",
			ParamsPy:    "l1: ListNode, l2: ListNode",
			ParamsJava:  "ListNode l1, ListNode l2",
			ParamsCpp:   "ListNode* l1, ListNode* l2",
			ParamsGo:    "l1 *ListNode, l2 *ListNode",
			RetPy:       "ListNode",
			RetJava:     "ListNode",
			RetCpp:      "ListNode*",
			RetGo:       "*ListNode",
			SolJS: `    let dummy = new ListNode(0), curr = dummy, carry = 0;
    while (l1 || l2 || carry) {
        let sum = carry;
        if (l1) { sum += l1.val; l1 = l1.next; }
        if (l2) { sum += l2.val; l2 = l2.next; }
        carry = Math.floor(sum / 10);
        curr.next = new ListNode(sum % 10);
        curr = curr.next;
    }
    return dummy.next;`,
			SolPy: `    dummy = ListNode(0)
    curr, carry = dummy, 0
    while l1 or l2 or carry:
        val = carry
        if l1: val += l1.val; l1 = l1.next
        if l2: val += l2.val; l2 = l2.next
        carry, val = divmod(val, 10)
        curr.next = ListNode(val)
        curr = curr.next
    return dummy.next`,
			SolJava: `        ListNode dummy = new ListNode(0), curr = dummy;
        int carry = 0;
        while (l1 != null || l2 != null || carry != 0) {
            int sum = carry;
            if (l1 != null) { sum += l1.val; l1 = l1.next; }
            if (l2 != null) { sum += l2.val; l2 = l2.next; }
            carry = sum / 10;
            curr.next = new ListNode(sum % 10);
            curr = curr.next;
        }
        return dummy.next;`,
			SolCpp: `        ListNode* dummy = new ListNode(0);
        ListNode* curr = dummy;
        int carry = 0;
        while (l1 || l2 || carry) {
            int sum = carry;
            if (l1) { sum += l1->val; l1 = l1->next; }
            if (l2) { sum += l2->val; l2 = l2->next; }
            carry = sum / 10;
            curr->next = new ListNode(sum % 10);
            curr = curr->next;
        }
        return dummy->next;`,
			SolGo: `    dummy := &ListNode{}
    curr, carry := dummy, 0
    for l1 != nil || l2 != nil || carry != 0 {
        sum := carry
        if l1 != nil { sum += l1.Val; l1 = l1.Next }
        if l2 != nil { sum += l2.Val; l2 = l2.Next }
        carry = sum / 10
        curr.Next = &ListNode{Val: sum % 10}
        curr = curr.Next
    }
    return dummy.Next`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0101",
			Slug:        "longest-substring-without-repeating-characters",
			Title:       "Longest Substring Without Repeating Characters",
			Topic:       "String",
			Statement:   "Given a string s, find the length of the longest substring without repeating characters.",
			ExampleIn:   "\"abcabcbb\"",
			ExampleOut:  "3",
			HintTitle:   "Sliding Window",
			HintBody:    "Use two pointers to represent the window and a hash set to track characters.",
			FuncName:    "lengthOfLongestSubstring",
			ParamsJS:    "s",
			ParamsPy:    "s: str",
			ParamsJava:  "String s",
			ParamsCpp:   "string s",
			ParamsGo:    "s string",
			RetPy:       "int",
			RetJava:     "int",
			RetCpp:      "int",
			RetGo:       "int",
			SolJS: `    let map = {}, max = 0, left = 0;
    for (let right = 0; right < s.length; right++) {
        if (map[s[right]] !== undefined) left = Math.max(left, map[s[right]] + 1);
        map[s[right]] = right;
        max = Math.max(max, right - left + 1);
    }
    return max;`,
			SolPy: `    mp, mx, l = {}, 0, 0
    for r, c in enumerate(s):
        if c in mp: l = max(l, mp[c] + 1)
        mp[c] = r
        mx = max(mx, r - l + 1)
    return mx`,
			SolJava: `        int[] mp = new int[128];
        java.util.Arrays.fill(mp, -1);
        int mx = 0, l = 0;
        for (int r = 0; r < s.length(); r++) {
            char c = s.charAt(r);
            if (mp[c] >= l) l = mp[c] + 1;
            mp[c] = r;
            mx = Math.max(mx, r - l + 1);
        }
        return mx;`,
			SolCpp: `        vector<int> mp(128, -1);
        int mx = 0, l = 0;
        for (int r = 0; r < s.length(); ++r) {
            char c = s[r];
            if (mp[c] >= l) l = mp[c] + 1;
            mp[c] = r;
            mx = max(mx, r - l + 1);
        }
        return mx;`,
			SolGo: `    mp := make(map[rune]int)
    mx, l := 0, 0
    for r, c := range s {
        if idx, ok := mp[c]; ok && idx >= l { l = idx + 1 }
        mp[c] = r
        if val := r - l + 1; val > mx { mx = val }
    }
    return mx`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0102",
			Slug:        "letter-combinations-of-a-phone-number",
			Title:       "Letter Combinations of a Phone Number",
			Topic:       "Backtracking",
			Statement:   "Given a string containing digits from 2-9 inclusive, return all possible letter combinations that the number could represent. Return the answer in any order.",
			ExampleIn:   "\"23\"",
			ExampleOut:  "[\"ad\",\"ae\",\"af\",\"bd\",\"be\",\"bf\",\"cd\",\"ce\",\"cf\"]",
			HintTitle:   "Backtracking",
			HintBody:    "Map each digit to letters and use backtracking to generate combinations.",
			FuncName:    "letterCombinations",
			ParamsJS:    "digits",
			ParamsPy:    "digits: str",
			ParamsJava:  "String digits",
			ParamsCpp:   "string digits",
			ParamsGo:    "digits string",
			RetPy:       "list[str]",
			RetJava:     "List<String>",
			RetCpp:      "vector<string>",
			RetGo:       "[]string",
			SolJS: `    if (!digits) return [];
    const map = ["", "", "abc", "def", "ghi", "jkl", "mno", "pqrs", "tuv", "wxyz"];
    let res = [""];
    for (let d of digits) {
        let next = [];
        for (let s of res) {
            for (let c of map[d - '0']) next.push(s + c);
        }
        res = next;
    }
    return res;`,
			SolPy: `    if not digits: return []
    mapping = ["", "", "abc", "def", "ghi", "jkl", "mno", "pqrs", "tuv", "wxyz"]
    res = [""]
    for d in digits:
        res = [s + c for s in res for c in mapping[int(d)]]
    return res`,
			SolJava: `        List<String> res = new ArrayList<>();
        if (digits.isEmpty()) return res;
        String[] mapping = {"", "", "abc", "def", "ghi", "jkl", "mno", "pqrs", "tuv", "wxyz"};
        res.add("");
        for (int i = 0; i < digits.length(); i++) {
            List<String> next = new ArrayList<>();
            String letters = mapping[digits.charAt(i) - '0'];
            for (String s : res) {
                for (char c : letters.toCharArray()) next.add(s + c);
            }
            res = next;
        }
        return res;`,
			SolCpp: `        vector<string> res;
        if (digits.empty()) return res;
        vector<string> mapping = {"", "", "abc", "def", "ghi", "jkl", "mno", "pqrs", "tuv", "wxyz"};
        res.push_back("");
        for (char d : digits) {
            vector<string> next;
            for (string s : res) {
                for (char c : mapping[d - '0']) next.push_back(s + c);
            }
            res = next;
        }
        return res;`,
			SolGo: `    if len(digits) == 0 { return nil }
    mapping := []string{"", "", "abc", "def", "ghi", "jkl", "mno", "pqrs", "tuv", "wxyz"}
    res := []string{""}
    for i := 0; i < len(digits); i++ {
        var next []string
        letters := mapping[digits[i]-'0']
        for _, s := range res {
            for _, c := range letters { next = append(next, s+string(c)) }
        }
        res = next
    }
    return res`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0103",
			Slug:        "generate-parentheses",
			Title:       "Generate Parentheses",
			Topic:       "Backtracking",
			Statement:   "Given n pairs of parentheses, write a function to generate all combinations of well-formed parentheses.",
			ExampleIn:   "3",
			ExampleOut:  "[\"((()))\",\"(()())\",\"(())()\",\"()(())\",\"()()()\"]",
			HintTitle:   "Count Open and Close",
			HintBody:    "Track the number of open and closed parentheses used during backtracking.",
			FuncName:    "generateParenthesis",
			ParamsJS:    "n",
			ParamsPy:    "n: int",
			ParamsJava:  "int n",
			ParamsCpp:   "int n",
			ParamsGo:    "n int",
			RetPy:       "list[str]",
			RetJava:     "List<String>",
			RetCpp:      "vector<string>",
			RetGo:       "[]string",
			SolJS: `    let res = [];
    const backtrack = (s, open, close) => {
        if (s.length === 2 * n) { res.push(s); return; }
        if (open < n) backtrack(s + "(", open + 1, close);
        if (close < open) backtrack(s + ")", open, close + 1);
    };
    backtrack("", 0, 0);
    return res;`,
			SolPy: `    res = []
    def backtrack(s, open, close):
        if len(s) == 2 * n:
            res.append(s)
            return
        if open < n: backtrack(s + "(", open + 1, close)
        if close < open: backtrack(s + ")", open, close + 1)
    backtrack("", 0, 0)
    return res`,
			SolJava: `        List<String> res = new ArrayList<>();
        backtrack(res, "", 0, 0, n);
        return res;
    }
    private void backtrack(List<String> res, String s, int open, int close, int n) {
        if (s.length() == 2 * n) { res.add(s); return; }
        if (open < n) backtrack(res, s + "(", open + 1, close, n);
        if (close < open) backtrack(res, s + ")", open, close + 1, n);`,
			SolCpp: `        vector<string> res;
        backtrack(res, "", 0, 0, n);
        return res;
    }
    void backtrack(vector<string>& res, string s, int open, int close, int n) {
        if (s.length() == 2 * n) { res.push_back(s); return; }
        if (open < n) backtrack(res, s + "(", open + 1, close, n);
        if (close < open) backtrack(res, s + ")", open, close + 1, n);`,
			SolGo: `    var res []string
    var backtrack func(string, int, int)
    backtrack = func(s string, open int, close int) {
        if len(s) == 2*n { res = append(res, s); return }
        if open < n { backtrack(s+"(", open+1, close) }
        if close < open { backtrack(s+")", open, close+1) }
    }
    backtrack("", 0, 0)
    return res`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0104",
			Slug:        "swap-nodes-in-pairs",
			Title:       "Swap Nodes in Pairs",
			Topic:       "Linked List",
			Statement:   "Given a linked list, swap every two adjacent nodes and return its head. You must solve the problem without modifying the values in the list's nodes (i.e., only nodes themselves may be changed.)",
			ExampleIn:   "head = [1,2,3,4]",
			ExampleOut:  "[2,1,4,3]",
			HintTitle:   "Pointers Swap",
			HintBody:    "Adjust next pointers in pairs recursively or iteratively.",
			FuncName:    "swapPairs",
			ParamsJS:    "head",
			ParamsPy:    "head: ListNode",
			ParamsJava:  "ListNode head",
			ParamsCpp:   "ListNode* head",
			ParamsGo:    "head *ListNode",
			RetPy:       "ListNode",
			RetJava:     "ListNode",
			RetCpp:      "ListNode*",
			RetGo:       "*ListNode",
			SolJS: `    if (!head || !head.next) return head;
    let nextNode = head.next;
    head.next = swapPairs(nextNode.next);
    nextNode.next = head;
    return nextNode;`,
			SolPy: `    if not head or not head.next: return head
    n = head.next
    head.next = swapPairs(n.next)
    n.next = head
    return n`,
			SolJava: `        if (head == null || head.next == null) return head;
        ListNode nextNode = head.next;
        head.next = swapPairs(nextNode.next);
        nextNode.next = head;
        return nextNode;`,
			SolCpp: `        if (!head || !head->next) return head;
        ListNode* nextNode = head->next;
        head->next = swapPairs(nextNode->next);
        nextNode->next = head;
        return nextNode;`,
			SolGo: `    if head == nil || head.Next == nil { return head }
    n := head.Next
    head.Next = swapPairs(n.Next)
    n.Next = head
    return n`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0105",
			Slug:        "divide-two-integers",
			Title:       "Divide Two Integers",
			Topic:       "Array",
			Statement:   "Given two integers dividend and divisor, divide two integers without using multiplication, division, and mod operator.",
			ExampleIn:   "10, 3",
			ExampleOut:  "3",
			HintTitle:   "Bit Manipulation",
			HintBody:    "Subtract divisor multiplied by powers of 2 using left shifts.",
			FuncName:    "divide",
			ParamsJS:    "dividend, divisor",
			ParamsPy:    "dividend: int, divisor: int",
			ParamsJava:  "int dividend, int divisor",
			ParamsCpp:   "int dividend, int divisor",
			ParamsGo:    "dividend int, divisor int",
			RetPy:       "int",
			RetJava:     "int",
			RetCpp:      "int",
			RetGo:       "int",
			SolJS: `    if (dividend === -2147483648 && divisor === -1) return 2147483647;
    let sign = (dividend < 0) ^ (divisor < 0) ? -1 : 1;
    let dvd = Math.abs(dividend), dvs = Math.abs(divisor), res = 0;
    while (dvd >= dvs) {
        let temp = dvs, mul = 1;
        while (dvd >= (temp * 2) && temp <= 1073741823) { temp *= 2; mul *= 2; }
        dvd -= temp; res += mul;
    }
    return res * sign;`,
			SolPy: `    if dividend == -2147483648 and divisor == -1: return 2147483647
    sign = -1 if (dividend < 0) ^ (divisor < 0) else 1
    dvd, dvs, res = abs(dividend), abs(divisor), 0
    while dvd >= dvs:
        temp, mul = dvs, 1
        while dvd >= (temp << 1):
            temp <<= 1
            mul <<= 1
        dvd -= temp
        res += mul
    return res * sign`,
			SolJava: `        if (dividend == Integer.MIN_VALUE && divisor == -1) return Integer.MAX_VALUE;
        int sign = (dividend < 0) ^ (divisor < 0) ? -1 : 1;
        long dvd = Math.abs((long) dividend);
        long dvs = Math.abs((long) divisor);
        int res = 0;
        while (dvd >= dvs) {
            long temp = dvs, mul = 1;
            while (dvd >= (temp << 1)) { temp <<= 1; mul <<= 1; }
            dvd -= temp;
            res += mul;
        }
        return res * sign;`,
			SolCpp: `        if (dividend == INT_MIN && divisor == -1) return INT_MAX;
        int sign = (dividend < 0) ^ (divisor < 0) ? -1 : 1;
        long long dvd = abs((long long) dividend);
        long long dvs = abs((long long) divisor);
        long long res = 0;
        while (dvd >= dvs) {
            long long temp = dvs, mul = 1;
            while (dvd >= (temp << 1)) { temp <<= 1; mul <<= 1; }
            dvd -= temp;
            res += mul;
        }
        return res * sign;`,
			SolGo: `    if dividend == -2147483648 && divisor == -1 { return 2147483647 }
    sign := 1
    if (dividend < 0) != (divisor < 0) { sign = -1 }
    dvd := dividend
    if dvd < 0 { dvd = -dvd }
    dvs := divisor
    if dvs < 0 { dvs = -dvs }
    res := 0
    for dvd >= dvs {
        temp, mul := dvs, 1
        for dvd >= (temp << 1) {
            temp <<= 1
            mul <<= 1
        }
        dvd -= temp
        res += mul
    }
    return res * sign`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0106",
			Slug:        "next-permutation",
			Title:       "Next Permutation",
			Topic:       "Array",
			Statement:   "A permutation of an array of integers is its arrangement into a lexicographical order. Find the next lexicographically greater permutation.",
			ExampleIn:   "[1,2,3]",
			ExampleOut:  "[1,3,2]",
			HintTitle:   "Scan Right to Left",
			HintBody:    "Find first decreasing element, swap with next larger, and reverse remainder.",
			FuncName:    "nextPermutation",
			ParamsJS:    "nums",
			ParamsPy:    "nums: list[int]",
			ParamsJava:  "int[] nums",
			ParamsCpp:   "vector<int>& nums",
			ParamsGo:    "nums []int",
			RetPy:       "None",
			RetJava:     "void",
			RetCpp:      "void",
			RetGo:       "",
			SolJS: `    let i = nums.length - 2;
    while (i >= 0 && nums[i] >= nums[i + 1]) i--;
    if (i >= 0) {
        let j = nums.length - 1;
        while (nums[j] <= nums[i]) j--;
        [nums[i], nums[j]] = [nums[j], nums[i]];
    }
    let l = i + 1, r = nums.length - 1;
    while (l < r) { [nums[l], nums[r]] = [nums[r], nums[l]]; l++; r--; }`,
			SolPy: `    i = len(nums) - 2
    while i >= 0 and nums[i] >= nums[i + 1]: i -= 1
    if i >= 0:
        j = len(nums) - 1
        while nums[j] <= nums[i]: j -= 1
        nums[i], nums[j] = nums[j], nums[i]
    l, r = i + 1, len(nums) - 1
    while l < r:
        nums[l], nums[r] = nums[r], nums[l]
        l += 1; r -= 1`,
			SolJava: `        int i = nums.length - 2;
        while (i >= 0 && nums[i] >= nums[i + 1]) i--;
        if (i >= 0) {
            int j = nums.length - 1;
            while (nums[j] <= nums[i]) j--;
            int temp = nums[i]; nums[i] = nums[j]; nums[j] = temp;
        }
        int l = i + 1, r = nums.length - 1;
        while (l < r) {
            int temp = nums[l]; nums[l] = nums[r]; nums[r] = temp;
            l++; r--;
        }`,
			SolCpp: `        int i = nums.size() - 2;
        while (i >= 0 && nums[i] >= nums[i + 1]) i--;
        if (i >= 0) {
            int j = nums.size() - 1;
            while (nums[j] <= nums[i]) j--;
            swap(nums[i], nums[j]);
        }
        reverse(nums.begin() + i + 1, nums.end());`,
			SolGo: `    i := len(nums) - 2
    for i >= 0 && nums[i] >= nums[i+1] { i-- }
    if i >= 0 {
        j := len(nums) - 1
        for nums[j] <= nums[i] { j-- }
        nums[i], nums[j] = nums[j], nums[i]
    }
    l, r := i+1, len(nums)-1
    for l < r {
        nums[l], nums[r] = nums[r], nums[l]
        l++; r--
    }`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0107",
			Slug:        "combination-sum",
			Title:       "Combination Sum",
			Topic:       "Backtracking",
			Statement:   "Given an array of distinct integers candidates and a target integer target, return a list of all unique combinations of candidates where the chosen numbers sum to target.",
			ExampleIn:   "[2,3,6,7], 7",
			ExampleOut:  "[[2,2,3],[7]]",
			HintTitle:   "Backtracking dfs",
			HintBody:    "Recurse with the option to select the same candidate or move to next.",
			FuncName:    "combinationSum",
			ParamsJS:    "candidates, target",
			ParamsPy:    "candidates: list[int], target: int",
			ParamsJava:  "int[] candidates, int target",
			ParamsCpp:   "vector<int>& candidates, int target",
			ParamsGo:    "candidates []int, target int",
			RetPy:       "list[list[int]]",
			RetJava:     "List<List<Integer>>",
			RetCpp:      "vector<vector<int>>",
			RetGo:       "[][]int",
			SolJS: `    let res = [];
    const dfs = (start, target, path) => {
        if (target === 0) { res.push([...path]); return; }
        for (let i = start; i < candidates.length; i++) {
            if (candidates[i] <= target) {
                path.push(candidates[i]);
                dfs(i, target - candidates[i], path);
                path.pop();
            }
        }
    };
    dfs(0, target, []);
    return res;`,
			SolPy: `    res = []
    def dfs(start, target, path):
        if target == 0: res.append(list(path)); return
        for i in range(start, len(candidates)):
            if candidates[i] <= target:
                path.append(candidates[i])
                dfs(i, target - candidates[i], path)
                path.pop()
    dfs(0, target, [])
    return res`,
			SolJava: `        List<List<Integer>> res = new ArrayList<>();
        dfs(0, candidates, target, new ArrayList<>(), res);
        return res;
    }
    private void dfs(int start, int[] candidates, int target, List<Integer> path, List<List<Integer>> res) {
        if (target == 0) { res.add(new ArrayList<>(path)); return; }
        for (int i = start; i < candidates.length; i++) {
            if (candidates[i] <= target) {
                path.add(candidates[i]);
                dfs(i, candidates, target - candidates[i], path, res);
                path.remove(path.size() - 1);
            }
        }`,
			SolCpp: `        vector<vector<int>> res;
        vector<int> path;
        dfs(0, candidates, target, path, res);
        return res;
    }
    void dfs(int start, vector<int>& candidates, int target, vector<int>& path, vector<vector<int>>& res) {
        if (target == 0) { res.push_back(path); return; }
        for (int i = start; i < candidates.size(); ++i) {
            if (candidates[i] <= target) {
                path.push_back(candidates[i]);
                dfs(i, candidates, target - candidates[i], path, res);
                path.pop_back();
            }
        }`,
			SolGo: `    var res [][]int
    var dfs func(int, int, []int)
    dfs = func(start int, target int, path []int) {
        if target == 0 {
            temp := make([]int, len(path))
            copy(temp, path)
            res = append(res, temp)
            return
        }
        for i := start; i < len(candidates); i++ {
            if candidates[i] <= target {
                dfs(i, target-candidates[i], append(path, candidates[i]))
            }
        }
    }
    dfs(0, target, nil)
    return res`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0108",
			Slug:        "permutations",
			Title:       "Permutations",
			Topic:       "Backtracking",
			Statement:   "Given an array nums of distinct integers, return all the possible permutations. You can return the answer in any order.",
			ExampleIn:   "[1,2,3]",
			ExampleOut:  "[[1,2,3],[1,3,2],[2,1,3],[2,3,1],[3,1,2],[3,2,1]]",
			HintTitle:   "DFS swap",
			HintBody:    "Swap elements recursively to form new permutations.",
			FuncName:    "permute",
			ParamsJS:    "nums",
			ParamsPy:    "nums: list[int]",
			ParamsJava:  "int[] nums",
			ParamsCpp:   "vector<int>& nums",
			ParamsGo:    "nums []int",
			RetPy:       "list[list[int]]",
			RetJava:     "List<List<Integer>>",
			RetCpp:      "vector<vector<int>>",
			RetGo:       "[][]int",
			SolJS: `    let res = [];
    const backtrack = (first) => {
        if (first === nums.length) res.push([...nums]);
        for (let i = first; i < nums.length; i++) {
            [nums[first], nums[i]] = [nums[i], nums[first]];
            backtrack(first + 1);
            [nums[first], nums[i]] = [nums[i], nums[first]];
        }
    };
    backtrack(0);
    return res;`,
			SolPy: `    res = []
    def backtrack(first):
        if first == len(nums): res.append(list(nums))
        for i in range(first, len(nums)):
            nums[first], nums[i] = nums[i], nums[first]
            backtrack(first + 1)
            nums[first], nums[i] = nums[i], nums[first]
    backtrack(0)
    return res`,
			SolJava: `        List<List<Integer>> res = new ArrayList<>();
        List<Integer> list = new ArrayList<>();
        for (int n : nums) list.add(n);
        backtrack(nums.length, list, res, 0);
        return res;
    }
    private void backtrack(int n, List<Integer> list, List<List<Integer>> res, int first) {
        if (first == n) res.add(new ArrayList<>(list));
        for (int i = first; i < n; i++) {
            Collections.swap(list, first, i);
            backtrack(n, list, res, first + 1);
            Collections.swap(list, first, i);
        }`,
			SolCpp: `        vector<vector<int>> res;
        backtrack(nums, res, 0);
        return res;
    }
    void backtrack(vector<int>& nums, vector<vector<int>>& res, int first) {
        if (first == nums.size()) { res.push_back(nums); return; }
        for (int i = first; i < nums.size(); ++i) {
            swap(nums[first], nums[i]);
            backtrack(nums, res, first + 1);
            swap(nums[first], nums[i]);
        }`,
			SolGo: `    var res [][]int
    var backtrack func(int)
    backtrack = func(first int) {
        if first == len(nums) {
            temp := make([]int, len(nums))
            copy(temp, nums)
            res = append(res, temp)
            return
        }
        for i := first; i < len(nums); i++ {
            nums[first], nums[i] = nums[i], nums[first]
            backtrack(first + 1)
            nums[first], nums[i] = nums[i], nums[first]
        }
    }
    backtrack(0)
    return res`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0109",
			Slug:        "rotate-image",
			Title:       "Rotate Image",
			Topic:       "Array",
			Statement:   "You are given an n x n 2D matrix representing an image, rotate the image by 90 degrees (clockwise) in-place.",
			ExampleIn:   "[[1,2,3],[4,5,6],[7,8,9]]",
			ExampleOut:  "[[7,4,1],[8,5,2],[9,6,3]]",
			HintTitle:   "Transpose and Reverse",
			HintBody:    "Transpose the matrix first, then reverse each row.",
			FuncName:    "rotate",
			ParamsJS:    "matrix",
			ParamsPy:    "matrix: list[list[int]]",
			ParamsJava:  "int[][] matrix",
			ParamsCpp:   "vector<vector<int>>& matrix",
			ParamsGo:    "matrix [][]int",
			RetPy:       "None",
			RetJava:     "void",
			RetCpp:      "void",
			RetGo:       "",
			SolJS: `    let n = matrix.length;
    for (let i = 0; i < n; i++) {
        for (let j = i + 1; j < n; j++) {
            [matrix[i][j], matrix[j][i]] = [matrix[j][i], matrix[i][j]];
        }
    }
    for (let i = 0; i < n; i++) matrix[i].reverse();`,
			SolPy: `    n = len(matrix)
    for i in range(n):
        for j in range(i + 1, n):
            matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
    for i in range(n):
        matrix[i].reverse()`,
			SolJava: `        int n = matrix.length;
        for (int i = 0; i < n; i++) {
            for (int j = i + 1; j < n; j++) {
                int temp = matrix[i][j];
                matrix[i][j] = matrix[j][i];
                matrix[j][i] = temp;
            }
        }
        for (int i = 0; i < n; i++) {
            for (int j = 0; j < n / 2; j++) {
                int temp = matrix[i][j];
                matrix[i][j] = matrix[i][n - 1 - j];
                matrix[i][n - 1 - j] = temp;
            }
        }`,
			SolCpp: `        int n = matrix.size();
        for (int i = 0; i < n; ++i) {
            for (int j = i + 1; j < n; ++j) {
                swap(matrix[i][j], matrix[j][i]);
            }
        }
        for (int i = 0; i < n; ++i) {
            reverse(matrix[i].begin(), matrix[i].end());
        }`,
			SolGo: `    n := len(matrix)
    for i := 0; i < n; i++ {
        for j := i + 1; j < n; j++ {
            matrix[i][j], matrix[j][i] = matrix[j][i], matrix[i][j]
        }
    }
    for i := 0; i < n; i++ {
        for j := 0; j < n/2; j++ {
            matrix[i][j], matrix[i][n-1-j] = matrix[i][n-1-j], matrix[i][j]
        }
    }`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0110",
			Slug:        "group-anagrams",
			Title:       "Group Anagrams",
			Topic:       "HashMap",
			Statement:   "Given an array of strings strs, group the anagrams together. You can return the answer in any order.",
			ExampleIn:   "[\"eat\",\"tea\",\"tan\",\"ate\",\"nat\",\"bat\"]",
			ExampleOut:  "[[\"bat\"],[\"nat\",\"tan\"],[\"ate\",\"eat\",\"tea\"]]",
			HintTitle:   "Hash Map Key",
			HintBody:    "Sort each string to use as key in a hash map.",
			FuncName:    "groupAnagrams",
			ParamsJS:    "strs",
			ParamsPy:    "strs: list[str]",
			ParamsJava:  "String[] strs",
			ParamsCpp:   "vector<string>& strs",
			ParamsGo:    "strs []string",
			RetPy:       "list[list[str]]",
			RetJava:     "List<List<String>>",
			RetCpp:      "vector<vector<string>>",
			RetGo:       "[][]string",
			SolJS: `    let map = {};
    for (let s of strs) {
        let sorted = s.split('').sort().join('');
        if (!map[sorted]) map[sorted] = [];
        map[sorted].push(s);
    }
    return Object.values(map);`,
			SolPy: `    from collections import defaultdict
    mp = defaultdict(list)
    for s in strs:
        mp["".join(sorted(s))].append(s)
    return list(mp.values())`,
			SolJava: `        if (strs == null || strs.length == 0) return new ArrayList<>();
        Map<String, List<String>> map = new HashMap<>();
        for (String s : strs) {
            char[] ca = s.toCharArray();
            Arrays.sort(ca);
            String key = String.valueOf(ca);
            if (!map.containsKey(key)) map.put(key, new ArrayList<>());
            map.get(key).add(s);
        }
        return new ArrayList<>(map.values());`,
			SolCpp: `        unordered_map<string, vector<string>> map;
        for (string s : strs) {
            string t = s;
            sort(t.begin(), t.end());
            map[t].push_back(s);
        }
        vector<vector<string>> res;
        for (auto p : map) res.push_back(p.second);
        return res;`,
			SolGo: `    mp := make(map[string][]string)
    for _, s := range strs {
        r := []rune(s)
        // Bubble sort or sort package, but sorting runes is quick
        for i := 0; i < len(r); i++ {
            for j := i+1; j < len(r); j++ {
                if r[i] > r[j] { r[i], r[j] = r[j], r[i] }
            }
        }
        key := string(r)
        mp[key] = append(mp[key], s)
    }
    var res [][]string
    for _, val := range mp { res = append(res, val) }
    return res`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0111",
			Slug:        "jump-game",
			Title:       "Jump Game",
			Topic:       "DP",
			Statement:   "You are given an integer array nums. You are initially positioned at the array's first index, and each element in the array represents your maximum jump length at that position. Return true if you can reach the last index.",
			ExampleIn:   "[2,3,1,1,4]",
			ExampleOut:  "true",
			HintTitle:   "Greedy Reach",
			HintBody:    "Keep track of the maximum reachable index as you scan the array.",
			FuncName:    "canJump",
			ParamsJS:    "nums",
			ParamsPy:    "nums: list[int]",
			ParamsJava:  "int[] nums",
			ParamsCpp:   "vector<int>& nums",
			ParamsGo:    "nums []int",
			RetPy:       "bool",
			RetJava:     "boolean",
			RetCpp:      "bool",
			RetGo:       "bool",
			SolJS: `    let reachable = 0;
    for (let i = 0; i < nums.length; i++) {
        if (i > reachable) return false;
        reachable = Math.max(reachable, i + nums[i]);
    }
    return true;`,
			SolPy: `    reachable = 0
    for i, num in enumerate(nums):
        if i > reachable: return False
        reachable = max(reachable, i + num)
    return True`,
			SolJava: `        int reachable = 0;
        for (int i = 0; i < nums.length; i++) {
            if (i > reachable) return false;
            reachable = Math.max(reachable, i + nums[i]);
        }
        return true;`,
			SolCpp: `        int reachable = 0;
        for (int i = 0; i < nums.size(); ++i) {
            if (i > reachable) return false;
            reachable = max(reachable, i + nums[i]);
        }
        return true;`,
			SolGo: `    reachable := 0
    for i, num := range nums {
        if i > reachable { return false }
        if i + num > reachable { reachable = i + num }
    }
    return true`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0112",
			Slug:        "merge-intervals",
			Title:       "Merge Intervals",
			Topic:       "Array",
			Statement:   "Given an array of intervals where intervals[i] = [start_i, end_i], merge all overlapping intervals and return an array of the non-overlapping intervals.",
			ExampleIn:   "[[1,3],[2,6],[8,10],[15,18]]",
			ExampleOut:  "[[1,6],[8,10],[15,18]]",
			HintTitle:   "Sorting",
			HintBody:    "Sort intervals by start time before merging.",
			FuncName:    "merge",
			ParamsJS:    "intervals",
			ParamsPy:    "intervals: list[list[int]]",
			ParamsJava:  "int[][] intervals",
			ParamsCpp:   "vector<vector<int>>& intervals",
			ParamsGo:    "intervals [][]int",
			RetPy:       "list[list[int]]",
			RetJava:     "int[][]",
			RetCpp:      "vector<vector<int>>",
			RetGo:       "[][]int",
			SolJS: `    if (intervals.length <= 1) return intervals;
    intervals.sort((a, b) => a[0] - b[0]);
    let res = [intervals[0]];
    for (let i = 1; i < intervals.length; i++) {
        let last = res[res.length - 1];
        if (intervals[i][0] <= last[1]) last[1] = Math.max(last[1], intervals[i][1]);
        else res.push(intervals[i]);
    }
    return res;`,
			SolPy: `    if len(intervals) <= 1: return intervals
    intervals.sort(key=lambda x: x[0])
    res = [intervals[0]]
    for i in range(1, len(intervals)):
        if intervals[i][0] <= res[-1][1]:
            res[-1][1] = max(res[-1][1], intervals[i][1])
        else: res.append(intervals[i])
    return res`,
			SolJava: `        if (intervals.length <= 1) return intervals;
        Arrays.sort(intervals, (a, b) -> Integer.compare(a[0], b[0]));
        List<int[]> res = new ArrayList<>();
        res.add(intervals[0]);
        for (int i = 1; i < intervals.length; i++) {
            int[] last = res.get(res.size() - 1);
            if (intervals[i][0] <= last[1]) last[1] = Math.max(last[1], intervals[i][1]);
            else res.add(intervals[i]);
        }
        return res.toArray(new int[res.size()][]);`,
			SolCpp: `        if (intervals.size() <= 1) return intervals;
        sort(intervals.begin(), intervals.end());
        vector<vector<int>> res = {intervals[0]};
        for (int i = 1; i < intervals.size(); ++i) {
            if (intervals[i][0] <= res.back()[1]) res.back()[1] = max(res.back()[1], intervals[i][1]);
            else res.push_back(intervals[i]);
        }
        return res;`,
			SolGo: `    if len(intervals) <= 1 { return intervals }
    // Sort intervals by start time using quick sort logic or sort.Slice
    for i := 0; i < len(intervals); i++ {
        for j := i+1; j < len(intervals); j++ {
            if intervals[i][0] > intervals[j][0] { intervals[i], intervals[j] = intervals[j], intervals[i] }
        }
    }
    res := [][]int{intervals[0]}
    for i := 1; i < len(intervals); i++ {
        last := res[len(res)-1]
        if intervals[i][0] <= last[1] {
            if intervals[i][1] > last[1] { last[1] = intervals[i][1] }
        } else {
            res = append(res, intervals[i])
        }
    }
    return res`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0113",
			Slug:        "insert-interval",
			Title:       "Insert Interval",
			Topic:       "Array",
			Statement:   "You are given an array of non-overlapping intervals sorted by start time. Insert a new interval and merge if necessary.",
			ExampleIn:   "[[1,3],[6,9]], [2,5]",
			ExampleOut:  "[[1,5],[6,9]]",
			HintTitle:   "Three Parts",
			HintBody:    "Add all intervals ending before new interval, merge overlaps, then add the rest.",
			FuncName:    "insert",
			ParamsJS:    "intervals, newInterval",
			ParamsPy:    "intervals: list[list[int]], newInterval: list[int]",
			ParamsJava:  "int[][] intervals, int[] newInterval",
			ParamsCpp:   "vector<vector<int>>& intervals, vector<int>& newInterval",
			ParamsGo:    "intervals [][]int, newInterval []int",
			RetPy:       "list[list[int]]",
			RetJava:     "int[][]",
			RetCpp:      "vector<vector<int>>",
			RetGo:       "[][]int",
			SolJS: `    let res = [], i = 0, n = intervals.length;
    while (i < n && intervals[i][1] < newInterval[0]) res.push(intervals[i++]);
    while (i < n && intervals[i][0] <= newInterval[1]) {
        newInterval[0] = Math.min(newInterval[0], intervals[i][0]);
        newInterval[1] = Math.max(newInterval[1], intervals[i++][1]);
    }
    res.push(newInterval);
    while (i < n) res.push(intervals[i++]);
    return res;`,
			SolPy: `    res, i, n = [], 0, len(intervals)
    while i < n and intervals[i][1] < newInterval[0]:
        res.append(intervals[i]); i += 1
    while i < n and intervals[i][0] <= newInterval[1]:
        newInterval[0] = min(newInterval[0], intervals[i][0])
        newInterval[1] = max(newInterval[1], intervals[i][1])
        i += 1
    res.append(newInterval)
    while i < n:
        res.append(intervals[i]); i += 1
    return res`,
			SolJava: `        List<int[]> res = new ArrayList<>();
        int i = 0, n = intervals.length;
        while (i < n && intervals[i][1] < newInterval[0]) res.add(intervals[i++]);
        while (i < n && intervals[i][0] <= newInterval[1]) {
            newInterval[0] = Math.min(newInterval[0], intervals[i][0]);
            newInterval[1] = Math.max(newInterval[1], intervals[i++][1]);
        }
        res.add(newInterval);
        while (i < n) res.add(intervals[i++]);
        return res.toArray(new int[res.size()][]);`,
			SolCpp: `        vector<vector<int>> res;
        int i = 0, n = intervals.size();
        while (i < n && intervals[i][1] < newInterval[0]) res.push_back(intervals[i++]);
        while (i < n && intervals[i][0] <= newInterval[1]) {
            newInterval[0] = min(newInterval[0], intervals[i][0]);
            newInterval[1] = max(newInterval[1], intervals[i++][1]);
        }
        res.push_back(newInterval);
        while (i < n) res.push_back(intervals[i++]);
        return res;`,
			SolGo: `    var res [][]int
    i, n := 0, len(intervals)
    for i < n && intervals[i][1] < newInterval[0] {
        res = append(res, intervals[i])
        i++
    }
    for i < n && intervals[i][0] <= newInterval[1] {
        if intervals[i][0] < newInterval[0] { newInterval[0] = intervals[i][0] }
        if intervals[i][1] > newInterval[1] { newInterval[1] = intervals[i][1] }
        i++
    }
    res = append(res, newInterval)
    for i < n {
        res = append(res, intervals[i])
        i++
    }
    return res`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0114",
			Slug:        "unique-paths",
			Title:       "Unique Paths",
			Topic:       "DP",
			Statement:   "There is a robot on an m x n grid. The robot can only move either down or right. Find the number of unique paths to the bottom-right corner.",
			ExampleIn:   "3, 7",
			ExampleOut:  "28",
			HintTitle:   "Grid DP",
			HintBody:    "dp[i][j] = dp[i-1][j] + dp[i][j-1]",
			FuncName:    "uniquePaths",
			ParamsJS:    "m, n",
			ParamsPy:    "m: int, n: int",
			ParamsJava:  "int m, int n",
			ParamsCpp:   "int m, int n",
			ParamsGo:    "m int, n int",
			RetPy:       "int",
			RetJava:     "int",
			RetCpp:      "int",
			RetGo:       "int",
			SolJS: `    let dp = Array(n).fill(1);
    for (let i = 1; i < m; i++) {
        for (let j = 1; j < n; j++) dp[j] += dp[j - 1];
    }
    return dp[n - 1];`,
			SolPy: `    dp = [1] * n
    for i in range(1, m):
        for j in range(1, n):
            dp[j] += dp[j - 1]
    return dp[-1]`,
			SolJava: `        int[] dp = new int[n];
        Arrays.fill(dp, 1);
        for (int i = 1; i < m; i++) {
            for (int j = 1; j < n; j++) dp[j] += dp[j - 1];
        }
        return dp[n - 1];`,
			SolCpp: `        vector<int> dp(n, 1);
        for (int i = 1; i < m; ++i) {
            for (int j = 1; j < n; ++j) dp[j] += dp[j - 1];
        }
        return dp[n - 1];`,
			SolGo: `    dp := make([]int, n)
    for i := range dp { dp[i] = 1 }
    for i := 1; i < m; i++ {
        for j := 1; j < n; j++ { dp[j] += dp[j - 1] }
    }
    return dp[n - 1]`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0115",
			Slug:        "minimum-path-sum",
			Title:       "Minimum Path Sum",
			Topic:       "DP",
			Statement:   "Given a m x n grid filled with non-negative numbers, find a path from top left to bottom right which minimizes the sum of all numbers along its path.",
			ExampleIn:   "[[1,3,1],[1,5,1],[4,2,1]]",
			ExampleOut:  "7",
			HintTitle:   "Grid DP Min",
			HintBody:    "Each cell's min sum is the grid value plus min(top, left) cell sum.",
			FuncName:    "minPathSum",
			ParamsJS:    "grid",
			ParamsPy:    "grid: list[list[int]]",
			ParamsJava:  "int[][] grid",
			ParamsCpp:   "vector<vector<int>>& grid",
			ParamsGo:    "grid [][]int",
			RetPy:       "int",
			RetJava:     "int",
			RetCpp:      "int",
			RetGo:       "int",
			SolJS: `    let m = grid.length, n = grid[0].length;
    for (let i = 0; i < m; i++) {
        for (let j = 0; j < n; j++) {
            if (i === 0 && j === 0) continue;
            else if (i === 0) grid[i][j] += grid[i][j - 1];
            else if (j === 0) grid[i][j] += grid[i - 1][j];
            else grid[i][j] += Math.min(grid[i - 1][j], grid[i][j - 1]);
        }
    }
    return grid[m - 1][n - 1];`,
			SolPy: `    m, n = len(grid), len(grid[0])
    for i in range(m):
        for j in range(n):
            if i == 0 and j == 0: continue
            elif i == 0: grid[i][j] += grid[i][j - 1]
            elif j == 0: grid[i][j] += grid[i - 1][j]
            else: grid[i][j] += min(grid[i - 1][j], grid[i][j - 1])
    return grid[-1][-1]`,
			SolJava: `        int m = grid.length, n = grid[0].length;
        for (int i = 0; i < m; i++) {
            for (int j = 0; j < n; j++) {
                if (i == 0 && j == 0) continue;
                else if (i == 0) grid[i][j] += grid[i][j - 1];
                else if (j == 0) grid[i][j] += grid[i - 1][j];
                else grid[i][j] += Math.min(grid[i - 1][j], grid[i][j - 1]);
            }
        }
        return grid[m - 1][n - 1];`,
			SolCpp: `        int m = grid.size(), n = grid[0].size();
        for (int i = 0; i < m; ++i) {
            for (int j = 0; j < n; ++j) {
                if (i == 0 && j == 0) continue;
                else if (i == 0) grid[i][j] += grid[i][j - 1];
                else if (j == 0) grid[i][j] += grid[i - 1][j];
                else grid[i][j] += min(grid[i - 1][j], grid[i][j - 1]);
            }
        }
        return grid[m - 1][n - 1];`,
			SolGo: `    m, n := len(grid), len(grid[0])
    for i := 0; i < m; i++ {
        for j := 0; j < n; j++ {
            if i == 0 && j == 0 { continue }
            if i == 0 { grid[i][j] += grid[i][j - 1]; continue }
            if j == 0 { grid[i][j] += grid[i - 1][j]; continue }
            minVal := grid[i - 1][j]
            if grid[i][j - 1] < minVal { minVal = grid[i][j - 1] }
            grid[i][j] += minVal
        }
    }
    return grid[m - 1][n - 1]`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0116",
			Slug:        "set-matrix-zeroes",
			Title:       "Set Matrix Zeroes",
			Topic:       "Array",
			Statement:   "Given an m x n integer matrix, if an element is 0, set its entire row and column to 0. Do it in-place.",
			ExampleIn:   "[[1,1,1],[1,0,1],[1,1,1]]",
			ExampleOut:  "[[1,0,1],[0,0,0],[1,0,1]]",
			HintTitle:   "In-place Markers",
			HintBody:    "Use the first row and column as zero markers.",
			FuncName:    "setZeroes",
			ParamsJS:    "matrix",
			ParamsPy:    "matrix: list[list[int]]",
			ParamsJava:  "int[][] matrix",
			ParamsCpp:   "vector<vector<int>>& matrix",
			ParamsGo:    "matrix [][]int",
			RetPy:       "None",
			RetJava:     "void",
			RetCpp:      "void",
			RetGo:       "",
			SolJS: `    let isCol = false, m = matrix.length, n = matrix[0].length;
    for (let i = 0; i < m; i++) {
        if (matrix[i][0] === 0) isCol = true;
        for (let j = 1; j < n; j++) {
            if (matrix[i][j] === 0) { matrix[0][j] = 0; matrix[i][0] = 0; }
        }
    }
    for (let i = 1; i < m; i++) {
        for (let j = 1; j < n; j++) {
            if (matrix[i][0] === 0 || matrix[0][j] === 0) matrix[i][j] = 0;
        }
    }
    if (matrix[0][0] === 0) { for (let j = 0; j < n; j++) matrix[0][j] = 0; }
    if (isCol) { for (let i = 0; i < m; i++) matrix[i][0] = 0; }`,
			SolPy: `    is_col = False
    m, n = len(matrix), len(matrix[0])
    for i in range(m):
        if matrix[i][0] == 0: is_col = True
        for j in range(1, n):
            if matrix[i][j] == 0: matrix[0][j] = 0; matrix[i][0] = 0
    for i in range(1, m):
        for j in range(1, n):
            if matrix[i][0] == 0 or matrix[0][j] == 0: matrix[i][j] = 0
    if matrix[0][0] == 0:
        for j in range(n): matrix[0][j] = 0
    if is_col:
        for i in range(m): matrix[i][0] = 0`,
			SolJava: `        boolean isCol = false;
        int m = matrix.length, n = matrix[0].length;
        for (int i = 0; i < m; i++) {
            if (matrix[i][0] == 0) isCol = true;
            for (int j = 1; j < n; j++) {
                if (matrix[i][j] == 0) { matrix[0][j] = 0; matrix[i][0] = 0; }
            }
        }
        for (int i = 1; i < m; i++) {
            for (int j = 1; j < n; j++) {
                if (matrix[i][0] == 0 || matrix[0][j] == 0) matrix[i][j] = 0;
            }
        }
        if (matrix[0][0] == 0) { for (int j = 0; j < n; j++) matrix[0][j] = 0; }
        if (isCol) { for (int i = 0; i < m; i++) matrix[i][0] = 0; }`,
			SolCpp: `        bool isCol = false;
        int m = matrix.size(), n = matrix[0].size();
        for (int i = 0; i < m; ++i) {
            if (matrix[i][0] == 0) isCol = true;
            for (int j = 1; j < n; ++j) {
                if (matrix[i][j] == 0) { matrix[0][j] = 0; matrix[i][0] = 0; }
            }
        }
        for (int i = 1; i < m; ++i) {
            for (int j = 1; j < n; ++j) {
                if (matrix[i][0] == 0 || matrix[0][j] == 0) matrix[i][j] = 0;
            }
        }
        if (matrix[0][0] == 0) { for (int j = 0; j < n; ++j) matrix[0][j] = 0; }
        if (isCol) { for (int i = 0; i < m; ++i) matrix[i][0] = 0; }`,
			SolGo: `    isCol := false
    m, n := len(matrix), len(matrix[0])
    for i := 0; i < m; i++ {
        if matrix[i][0] == 0 { isCol = true }
        for j := 1; j < n; j++ {
            if matrix[i][j] == 0 { matrix[0][j] = 0; matrix[i][0] = 0 }
        }
    }
    for i := 1; i < m; i++ {
        for j := 1; j < n; j++ {
            if matrix[i][0] == 0 || matrix[0][j] == 0 { matrix[i][j] = 0 }
        }
    }
    if matrix[0][0] == 0 {
        for j := 0; j < n; j++ { matrix[0][j] = 0 }
    }
    if isCol {
        for i := 0; i < m; i++ { matrix[i][0] = 0 }
    }`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0117",
			Slug:        "search-a-2d-matrix",
			Title:       "Search a 2D Matrix",
			Topic:       "Array",
			Statement:   "Write an efficient algorithm that searches for a value target in an m x n integer matrix. The matrix is sorted.",
			ExampleIn:   "[[1,3,5,7],[10,11,16,20],[23,30,34,60]], 3",
			ExampleOut:  "true",
			HintTitle:   "Binary Search on 2D",
			HintBody:    "Treat the 2D matrix as a 1D sorted array and do binary search.",
			FuncName:    "searchMatrix",
			ParamsJS:    "matrix, target",
			ParamsPy:    "matrix: list[list[int]], target: int",
			ParamsJava:  "int[][] matrix, int target",
			ParamsCpp:   "vector<vector<int>>& matrix, int target",
			ParamsGo:    "matrix [][]int, target int",
			RetPy:       "bool",
			RetJava:     "boolean",
			RetCpp:      "bool",
			RetGo:       "bool",
			SolJS: `    if (!matrix.length) return false;
    let m = matrix.length, n = matrix[0].length, l = 0, r = m * n - 1;
    while (l <= r) {
        let mid = Math.floor((l + r) / 2);
        let val = matrix[Math.floor(mid / n)][mid % n];
        if (val === target) return true;
        else if (val < target) l = mid + 1;
        else r = mid - 1;
    }
    return false;`,
			SolPy: `    if not matrix: return False
    m, n = len(matrix), len(matrix[0])
    l, r = 0, m * n - 1
    while l <= r:
        mid = (l + r) // 2
        val = matrix[mid // n][mid % n]
        if val == target: return True
        elif val < target: l = mid + 1
        else: r = mid - 1
    return False`,
			SolJava: `        if (matrix.length == 0) return false;
        int m = matrix.length, n = matrix[0].length, l = 0, r = m * n - 1;
        while (l <= r) {
            int mid = (l + r) / 2;
            int val = matrix[mid / n][mid % n];
            if (val == target) return true;
            else if (val < target) l = mid + 1;
            else r = mid - 1;
        }
        return false;`,
			SolCpp: `        if (matrix.empty()) return false;
        int m = matrix.size(), n = matrix[0].size(), l = 0, r = m * n - 1;
        while (l <= r) {
            int mid = (l + r) / 2;
            int val = matrix[mid / n][mid % n];
            if (val == target) return true;
            else if (val < target) l = mid + 1;
            else r = mid - 1;
        }
        return false;`,
			SolGo: `    if len(matrix) == 0 { return false }
    m, n := len(matrix), len(matrix[0])
    l, r := 0, m*n-1
    for l <= r {
        mid := (l + r) / 2
        val := matrix[mid/n][mid%n]
        if val == target { return true }
        if val < target { l = mid + 1 } else { r = mid - 1 }
    }
    return false`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0118",
			Slug:        "sort-colors",
			Title:       "Sort Colors",
			Topic:       "Array",
			Statement:   "Given an array nums with n objects colored red, white, or blue, sort them in-place so that objects of the same color are adjacent.",
			ExampleIn:   "[2,0,2,1,1,0]",
			ExampleOut:  "[0,0,1,1,2,2]",
			HintTitle:   "Three Pointers",
			HintBody:    "Use Dutch National Flag algorithm with three pointers.",
			FuncName:    "sortColors",
			ParamsJS:    "nums",
			ParamsPy:    "nums: list[int]",
			ParamsJava:  "int[] nums",
			ParamsCpp:   "vector<int>& nums",
			ParamsGo:    "nums []int",
			RetPy:       "None",
			RetJava:     "void",
			RetCpp:      "void",
			RetGo:       "",
			SolJS: `    let p0 = 0, curr = 0, p2 = nums.length - 1;
    while (curr <= p2) {
        if (nums[curr] === 0) { [nums[p0], nums[curr]] = [nums[curr], nums[p0]]; p0++; curr++; }
        else if (nums[curr] === 2) { [nums[curr], nums[p2]] = [nums[p2], nums[curr]]; p2--; }
        else curr++;
    }`,
			SolPy: `    p0, curr, p2 = 0, 0, len(nums) - 1
    while curr <= p2:
        if nums[curr] == 0:
            nums[p0], nums[curr] = nums[curr], nums[p0]
            p0 += 1; curr += 1
        elif nums[curr] == 2:
            nums[curr], nums[p2] = nums[p2], nums[curr]
            p2 -= 1
        else: curr += 1`,
			SolJava: `        int p0 = 0, curr = 0, p2 = nums.length - 1;
        while (curr <= p2) {
            if (nums[curr] == 0) {
                int temp = nums[p0]; nums[p0] = nums[curr]; nums[curr] = temp;
                p0++; curr++;
            } else if (nums[curr] == 2) {
                int temp = nums[curr]; nums[curr] = nums[p2]; nums[p2] = temp;
                p2--;
            } else curr++;
        }`,
			SolCpp: `        int p0 = 0, curr = 0, p2 = nums.size() - 1;
        while (curr <= p2) {
            if (nums[curr] == 0) { swap(nums[p0++], nums[curr++]); }
            else if (nums[curr] == 2) { swap(nums[curr], nums[p2--]); }
            else curr++;
        }`,
			SolGo: `    p0, curr, p2 := 0, 0, len(nums)-1
    for curr <= p2 {
        if nums[curr] == 0 {
            nums[p0], nums[curr] = nums[curr], nums[p0]
            p0++; curr++
        } else if nums[curr] == 2 {
            nums[curr], nums[p2] = nums[p2], nums[curr]
            p2--
        } else { curr++ }
    }`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0119",
			Slug:        "subsets",
			Title:       "Subsets",
			Topic:       "Backtracking",
			Statement:   "Given an integer array nums of unique elements, return all possible subsets (the power set).",
			ExampleIn:   "[1,2,3]",
			ExampleOut:  "[[],[1],[2],[1,2],[3],[1,3],[2,3],[1,2,3]]",
			HintTitle:   "Backtracking or Bitmask",
			HintBody:    "Include or exclude each element in backtracking recursion.",
			FuncName:    "subsets",
			ParamsJS:    "nums",
			ParamsPy:    "nums: list[int]",
			ParamsJava:  "int[] nums",
			ParamsCpp:   "vector<int>& nums",
			ParamsGo:    "nums []int",
			RetPy:       "list[list[int]]",
			RetJava:     "List<List<Integer>>",
			RetCpp:      "vector<vector<int>>",
			RetGo:       "[][]int",
			SolJS: `    let res = [[]];
    for (let num of nums) {
        let len = res.length;
        for (let i = 0; i < len; i++) res.push([...res[i], num]);
    }
    return res;`,
			SolPy: `    res = [[]]
    for num in nums:
        res += [curr + [num] for curr in res]
    return res`,
			SolJava: `        List<List<Integer>> res = new ArrayList<>();
        res.add(new ArrayList<>());
        for (int num : nums) {
            int size = res.size();
            for (int i = 0; i < size; i++) {
                List<Integer> next = new ArrayList<>(res.get(i));
                next.add(num);
                res.add(next);
            }
        }
        return res;`,
			SolCpp: `        vector<vector<int>> res = {{}};
        for (int num : nums) {
            int n = res.size();
            for (int i = 0; i < n; ++i) {
                vector<int> next = res[i];
                next.push_back(num);
                res.push_back(next);
            }
        }
        return res;`,
			SolGo: `    res := [][]int{{}}
    for _, num := range nums {
        n := len(res)
        for i := 0; i < n; i++ {
            next := make([]int, len(res[i])+1)
            copy(next, res[i])
            next[len(res[i])] = num
            res = append(res, next)
        }
    }
    return res`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0120",
			Slug:        "word-search",
			Title:       "Word Search",
			Topic:       "Backtracking",
			Statement:   "Given an m x n grid of characters board and a string word, return true if word exists in the grid. The word can be constructed from letters of sequentially adjacent cells.",
			ExampleIn:   "[[\"A\",\"B\",\"C\",\"E\"],[\"S\",\"F\",\"C\",\"S\"],[\"A\",\"D\",\"E\",\"E\"]], \"ABCCED\"",
			ExampleOut:  "true",
			HintTitle:   "DFS Board Search",
			HintBody:    "Perform DFS from each cell checking adjacent matches. Mark visited cells.",
			FuncName:    "exist",
			ParamsJS:    "board, word",
			ParamsPy:    "board: list[list[str]], word: str",
			ParamsJava:  "char[][] board, String word",
			ParamsCpp:   "vector<vector<char>>& board, string word",
			ParamsGo:    "board [][]byte, word string",
			RetPy:       "bool",
			RetJava:     "boolean",
			RetCpp:      "bool",
			RetGo:       "bool",
			SolJS: `    let m = board.length, n = board[0].length;
    const dfs = (r, c, i) => {
        if (i === word.length) return true;
        if (r < 0 || r >= m || c < 0 || c >= n || board[r][c] !== word[i]) return false;
        let temp = board[r][c]; board[r][c] = '#';
        let found = dfs(r+1, c, i+1) || dfs(r-1, c, i+1) || dfs(r, c+1, i+1) || dfs(r, c-1, i+1);
        board[r][c] = temp;
        return found;
    };
    for (let i = 0; i < m; i++) {
        for (let j = 0; j < n; j++) {
            if (dfs(i, j, 0)) return true;
        }
    }
    return false;`,
			SolPy: `    m, n = len(board), len(board[0])
    def dfs(r, c, i):
        if i == len(word): return True
        if r < 0 or r >= m or c < 0 or c >= n or board[r][c] != word[i]: return False
        temp = board[r][c]; board[r][c] = '#'
        found = dfs(r+1, c, i+1) or dfs(r-1, c, i+1) or dfs(r, c+1, i+1) or dfs(r, c-1, i+1)
        board[r][c] = temp
        return found
    for i in range(m):
        for j in range(n):
            if dfs(i, j, 0): return True
    return False`,
			SolJava: `        int m = board.length, n = board[0].length;
        for (int i = 0; i < m; i++) {
            for (int j = 0; j < n; j++) {
                if (dfs(board, word, i, j, 0)) return true;
            }
        }
        return false;
    }
    private boolean dfs(char[][] board, String word, int r, int c, int i) {
        if (i == word.length()) return true;
        if (r < 0 || r >= board.length || c < 0 || c >= board[0].length || board[r][c] != word.charAt(i)) return false;
        char temp = board[r][c]; board[r][c] = '#';
        boolean found = dfs(board, word, r + 1, c, i + 1) || dfs(board, word, r - 1, c, i + 1) ||
                        dfs(board, word, r, c + 1, i + 1) || dfs(board, word, r, c - 1, i + 1);
        board[r][c] = temp;
        return found;`,
			SolCpp: `        int m = board.size(), n = board[0].size();
        for (int i = 0; i < m; ++i) {
            for (int j = 0; j < n; ++j) {
                if (dfs(board, word, i, j, 0)) return true;
            }
        }
        return false;
    }
    bool dfs(vector<vector<char>>& board, string& word, int r, int c, int i) {
        if (i == word.length()) return true;
        if (r < 0 || r >= board.size() || c < 0 || c >= board[0].size() || board[r][c] != word[i]) return false;
        char temp = board[r][c]; board[r][c] = '#';
        bool found = dfs(board, word, r + 1, c, i + 1) || dfs(board, word, r - 1, c, i + 1) ||
                     dfs(board, word, r, c + 1, i + 1) || dfs(board, word, r, c - 1, i + 1);
        board[r][c] = temp;
        return found;`,
			SolGo: `    m, n := len(board), len(board[0])
    var dfs func(int, int, int) bool
    dfs = func(r int, c int, i int) bool {
        if i == len(word) { return true }
        if r < 0 || r >= m || c < 0 || c >= n || board[r][c] != word[i] { return false }
        temp := board[r][c]
        board[r][c] = '#'
        found := dfs(r+1, c, i+1) || dfs(r-1, c, i+1) || dfs(r, c+1, i+1) || dfs(r, c-1, i+1)
        board[r][c] = temp
        return found
    }
    for i := 0; i < m; i++ {
        for j := 0; j < n; j++ {
            if dfs(i, j, 0) { return true }
        }
    }
    return false`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0121",
			Slug:        "validate-binary-search-tree",
			Title:       "Validate Binary Search Tree",
			Topic:       "Tree",
			Statement:   "Given the root of a binary tree, determine if it is a valid binary search tree (BST).",
			ExampleIn:   "root = [2,1,3]",
			ExampleOut:  "true",
			HintTitle:   "Range Constraints",
			HintBody:    "Pass down low and high bounds to check validity of each node.",
			FuncName:    "isValidBST",
			ParamsJS:    "root",
			ParamsPy:    "root: TreeNode",
			ParamsJava:  "TreeNode root",
			ParamsCpp:   "TreeNode* root",
			ParamsGo:    "root *TreeNode",
			RetPy:       "bool",
			RetJava:     "boolean",
			RetCpp:      "bool",
			RetGo:       "bool",
			SolJS: `    const validate = (node, low, high) => {
        if (!node) return true;
        if ((low !== null && node.val <= low) || (high !== null && node.val >= high)) return false;
        return validate(node.left, low, node.val) && validate(node.right, node.val, high);
    };
    return validate(root, null, null);`,
			SolPy: `    def validate(node, low=float('-inf'), high=float('inf')):
        if not node: return True
        if node.val <= low or node.val >= high: return False
        return validate(node.left, low, node.val) and validate(node.right, node.val, high)
    return validate(root)`,
			SolJava: `        return validate(root, null, null);
    }
    private boolean validate(TreeNode node, Integer low, Integer high) {
        if (node == null) return true;
        if ((low != null && node.val <= low) || (high != null && node.val >= high)) return false;
        return validate(node.left, low, node.val) && validate(node.right, node.val, high);`,
			SolCpp: `        return validate(root, nullptr, nullptr);
    }
    bool validate(TreeNode* node, long long* low, long long* high) {
        if (!node) return true;
        if ((low && node->val <= *low) || (high && node->val >= *high)) return false;
        long long val = node->val;
        return validate(node->left, low, &val) && validate(node->right, &val, high);`,
			SolGo: `    var validate func(*TreeNode, *int, *int) bool
    validate = func(node *TreeNode, low *int, high *int) bool {
        if node == nil { return true }
        if (low != nil && node.Val <= *low) || (high != nil && node.Val >= *high) { return false }
        val := node.Val
        return validate(node.Left, low, &val) && validate(node.Right, &val, high)
    }
    return validate(root, nil, nil)`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0122",
			Slug:        "binary-tree-level-order-traversal",
			Title:       "Binary Tree Level Order Traversal",
			Topic:       "Tree",
			Statement:   "Given the root of a binary tree, return the level order traversal of its nodes' values. (i.e., from left to right, level by level).",
			ExampleIn:   "root = [3,9,20,null,null,15,7]",
			ExampleOut:  "[[3],[9,20],[15,7]]",
			HintTitle:   "Queue BFS",
			HintBody:    "Use a queue to process nodes level by level.",
			FuncName:    "levelOrder",
			ParamsJS:    "root",
			ParamsPy:    "root: TreeNode",
			ParamsJava:  "TreeNode root",
			ParamsCpp:   "TreeNode* root",
			ParamsGo:    "root *TreeNode",
			RetPy:       "list[list[int]]",
			RetJava:     "List<List<Integer>>",
			RetCpp:      "vector<vector<int>>",
			RetGo:       "[][]int",
			SolJS: `    if (!root) return [];
    let res = [], q = [root];
    while (q.length) {
        let len = q.length, level = [];
        for (let i = 0; i < len; i++) {
            let n = q.shift(); level.push(n.val);
            if (n.left) q.push(n.left);
            if (n.right) q.push(n.right);
        }
        res.push(level);
    }
    return res;`,
			SolPy: `    if not root: return []
    res, q = [], [root]
    while q:
        level, length = [], len(q)
        for _ in range(length):
            node = q.pop(0); level.append(node.val)
            if node.left: q.append(node.left)
            if node.right: q.append(node.right)
        res.append(level)
    return res`,
			SolJava: `        List<List<Integer>> res = new ArrayList<>();
        if (root == null) return res;
        Queue<TreeNode> q = new LinkedList<>();
        q.add(root);
        while (!q.isEmpty()) {
            int len = q.size();
            List<Integer> level = new ArrayList<>();
            for (int i = 0; i < len; i++) {
                TreeNode n = q.poll();
                level.add(n.val);
                if (n.left != null) q.add(n.left);
                if (n.right != null) q.add(n.right);
            }
            res.add(level);
        }
        return res;`,
			SolCpp: `        vector<vector<int>> res;
        if (!root) return res;
        queue<TreeNode*> q; q.push(root);
        while (!q.empty()) {
            int len = q.size();
            vector<int> level;
            for (int i = 0; i < len; ++i) {
                TreeNode* n = q.front(); q.pop();
                level.push_back(n->val);
                if (n->left) q.push(n->left);
                if (n->right) q.push(n->right);
            }
            res.push_back(level);
        }
        return res;`,
			SolGo: `    if root == nil { return nil }
    var res [][]int
    q := []*TreeNode{root}
    for len(q) > 0 {
        length := len(q)
        var level []int
        for i := 0; i < length; i++ {
            n := q[0]; q = q[1:]
            level = append(level, n.Val)
            if n.Left != nil { q = append(q, n.Left) }
            if n.Right != nil { q = append(q, n.Right) }
        }
        res = append(res, level)
    }
    return res`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0123",
			Slug:        "construct-binary-tree-from-preorder-and-inorder-traversal",
			Title:       "Construct Binary Tree from Preorder and Inorder Traversal",
			Topic:       "Tree",
			Statement:   "Given two integer arrays preorder and inorder where preorder is the preorder traversal of a binary tree and inorder is the inorder traversal of the same tree, construct and return the binary tree.",
			ExampleIn:   "preorder = [3,9,20,15,7], inorder = [9,3,15,20,7]",
			ExampleOut:  "[3,9,20,null,null,15,7]",
			HintTitle:   "Map Inorder Indices",
			HintBody:    "The first element in preorder is the root. Find its position in inorder to split left and right subtrees.",
			FuncName:    "buildTree",
			ParamsJS:    "preorder, inorder",
			ParamsPy:    "preorder: list[int], inorder: list[int]",
			ParamsJava:  "int[] preorder, int[] inorder",
			ParamsCpp:   "vector<int>& preorder, vector<int>& inorder",
			ParamsGo:    "preorder []int, inorder []int",
			RetPy:       "TreeNode",
			RetJava:     "TreeNode",
			RetCpp:      "TreeNode*",
			RetGo:       "*TreeNode",
			SolJS: `    let map = {};
    for (let i = 0; i < inorder.length; i++) map[inorder[i]] = i;
    let preIdx = 0;
    const helper = (l, r) => {
        if (l > r) return null;
        let rootVal = preorder[preIdx++];
        let root = new TreeNode(rootVal);
        root.left = helper(l, map[rootVal] - 1);
        root.right = helper(map[rootVal] + 1, r);
        return root;
    };
    return helper(0, inorder.length - 1);`,
			SolPy: `    mp = {val: i for i, val in enumerate(inorder)}
    pre_idx = 0
    def helper(l, r):
        nonlocal pre_idx
        if l > r: return None
        val = preorder[pre_idx]
        pre_idx += 1
        root = TreeNode(val)
        root.left = helper(l, mp[val] - 1)
        root.right = helper(mp[val] + 1, r)
        return root
    return helper(0, len(inorder) - 1)`,
			SolJava: `        Map<Integer, Integer> map = new HashMap<>();
        for (int i = 0; i < inorder.length; i++) map.put(inorder[i], i);
        return helper(preorder, 0, inorder.length - 1, new int[]{0}, map);
    }
    private TreeNode helper(int[] preorder, int l, int r, int[] preIdx, Map<Integer, Integer> map) {
        if (l > r) return null;
        int rootVal = preorder[preIdx[0]++];
        TreeNode root = new TreeNode(rootVal);
        root.left = helper(preorder, l, map.get(rootVal) - 1, preIdx, map);
        root.right = helper(preorder, map.get(rootVal) + 1, r, preIdx, map);
        return root;`,
			SolCpp: `        unordered_map<int, int> map;
        for (int i = 0; i < inorder.size(); ++i) map[inorder[i]] = i;
        int preIdx = 0;
        return helper(preorder, 0, inorder.size() - 1, preIdx, map);
    }
    TreeNode* helper(vector<int>& preorder, int l, int r, int& preIdx, unordered_map<int, int>& map) {
        if (l > r) return nullptr;
        int rootVal = preorder[preIdx++];
        TreeNode* root = new TreeNode(rootVal);
        root->left = helper(preorder, l, map[rootVal] - 1, preIdx, map);
        root->right = helper(preorder, map[rootVal] + 1, r, preIdx, map);
        return root;`,
			SolGo: `    mp := make(map[int]int)
    for i, val := range inorder { mp[val] = i }
    preIdx := 0
    var helper func(int, int) *TreeNode
    helper = func(l, r int) *TreeNode {
        if l > r { return nil }
        rootVal := preorder[preIdx]
        preIdx++
        root := &TreeNode{Val: rootVal}
        root.Left = helper(l, mp[rootVal]-1)
        root.Right = helper(mp[rootVal]+1, r)
        return root
    }
    return helper(0, len(inorder)-1)`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0124",
			Slug:        "flatten-binary-tree-to-linked-list",
			Title:       "Flatten Binary Tree to Linked List",
			Topic:       "Tree",
			Statement:   "Given the root of a binary tree, flatten the tree into a single-linked list. The right pointer points to the next node and the left pointer is null.",
			ExampleIn:   "root = [1,2,5,3,4,null,6]",
			ExampleOut:  "[1,null,2,null,3,null,4,null,5,null,6]",
			HintTitle:   "Postorder traversal",
			HintBody:    "Flatten right, then left, and attach left list to right.",
			FuncName:    "flatten",
			ParamsJS:    "root",
			ParamsPy:    "root: TreeNode",
			ParamsJava:  "TreeNode root",
			ParamsCpp:   "TreeNode* root",
			ParamsGo:    "root *TreeNode",
			RetPy:       "None",
			RetJava:     "void",
			RetCpp:      "void",
			RetGo:       "",
			SolJS: `    let curr = root;
    while (curr) {
        if (curr.left) {
            let prev = curr.left;
            while (prev.right) prev = prev.right;
            prev.right = curr.right;
            curr.right = curr.left;
            curr.left = null;
        }
        curr = curr.right;
    }`,
			SolPy: `    curr = root
    while curr:
        if curr.left:
            prev = curr.left
            while prev.right: prev = prev.right
            prev.right = curr.right
            curr.right = curr.left
            curr.left = None
        curr = curr.right`,
			SolJava: `        TreeNode curr = root;
        while (curr != null) {
            if (curr.left != null) {
                TreeNode prev = curr.left;
                while (prev.right != null) prev = prev.right;
                prev.right = curr.right;
                curr.right = curr.left;
                curr.left = null;
            }
            curr = curr.right;
        }`,
			SolCpp: `        TreeNode* curr = root;
        while (curr) {
            if (curr->left) {
                TreeNode* prev = curr->left;
                while (prev->right) prev = prev->right;
                prev->right = curr->right;
                curr->right = curr->left;
                curr->left = nullptr;
            }
            curr = curr->right;
        }`,
			SolGo: `    curr := root
    for curr != nil {
        if curr.Left != nil {
            prev := curr.Left
            for prev.Right != nil { prev = prev.Right }
            prev.Right = curr.Right
            curr.Right = curr.Left
            curr.Left = nil
        }
        curr = curr.Right
    }`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0125",
			Slug:        "best-time-to-buy-and-sell-stock-ii",
			Title:       "Best Time to Buy and Sell Stock II",
			Topic:       "DP",
			Statement:   "You are given an integer array prices. Find the maximum profit you can achieve by buying and selling stocks multiple times.",
			ExampleIn:   "[7,1,5,3,6,4]",
			ExampleOut:  "7",
			HintTitle:   "Buy and Sell adjacent",
			HintBody:    "Sum all positive differences between prices on adjacent days.",
			FuncName:    "maxProfit",
			ParamsJS:    "prices",
			ParamsPy:    "prices: list[int]",
			ParamsJava:  "int[] prices",
			ParamsCpp:   "vector<int>& prices",
			ParamsGo:    "prices []int",
			RetPy:       "int",
			RetJava:     "int",
			RetCpp:      "int",
			RetGo:       "int",
			SolJS: `    let max = 0;
    for (let i = 1; i < prices.length; i++) {
        if (prices[i] > prices[i - 1]) max += prices[i] - prices[i - 1];
    }
    return max;`,
			SolPy: `    return sum(max(0, prices[i] - prices[i - 1]) for i in range(1, len(prices)))`,
			SolJava: `        int max = 0;
        for (int i = 1; i < prices.length; i++) {
            if (prices[i] > prices[i - 1]) max += prices[i] - prices[i - 1];
        }
        return max;`,
			SolCpp: `        int max = 0;
        for (int i = 1; i < prices.size(); ++i) {
            if (prices[i] > prices[i - 1]) max += prices[i] - prices[i - 1];
        }
        return max;`,
			SolGo: `    maxProfit := 0
    for i := 1; i < len(prices); i++ {
        if prices[i] > prices[i-1] { maxProfit += prices[i] - prices[i-1] }
    }
    return maxProfit`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0126",
			Slug:        "longest-consecutive-sequence",
			Title:       "Longest Consecutive Sequence",
			Topic:       "HashMap",
			Statement:   "Given an unsorted array of integers nums, return the length of the longest consecutive elements sequence. You must write an algorithm that runs in O(n) time.",
			ExampleIn:   "[100,4,200,1,3,2]",
			ExampleOut:  "4",
			HintTitle:   "Hash Set O(1)",
			HintBody:    "Insert all numbers into a hash set. Start checking sequences only from elements that have no predecessor.",
			FuncName:    "longestConsecutive",
			ParamsJS:    "nums",
			ParamsPy:    "nums: list[int]",
			ParamsJava:  "int[] nums",
			ParamsCpp:   "vector<int>& nums",
			ParamsGo:    "nums []int",
			RetPy:       "int",
			RetJava:     "int",
			RetCpp:      "int",
			RetGo:       "int",
			SolJS: `    let set = new Set(nums), longest = 0;
    for (let num of set) {
        if (!set.has(num - 1)) {
            let currNum = num, currStreak = 1;
            while (set.has(currNum + 1)) { currNum++; currStreak++; }
            longest = Math.max(longest, currStreak);
        }
    }
    return longest;`,
			SolPy: `    num_set = set(nums)
    longest = 0
    for num in num_set:
        if num - 1 not in num_set:
            curr = num
            streak = 1
            while curr + 1 in num_set:
                curr += 1; streak += 1
            longest = max(longest, streak)
    return longest`,
			SolJava: `        Set<Integer> set = new HashSet<>();
        for (int n : nums) set.add(n);
        int longest = 0;
        for (int num : set) {
            if (!set.contains(num - 1)) {
                int currNum = num;
                int currStreak = 1;
                while (set.contains(currNum + 1)) { currNum++; currStreak++; }
                longest = Math.max(longest, currStreak);
            }
        }
        return longest;`,
			SolCpp: `        unordered_set<int> set(nums.begin(), nums.end());
        int longest = 0;
        for (int num : set) {
            if (!set.count(num - 1)) {
                int currNum = num;
                int currStreak = 1;
                while (set.count(currNum + 1)) { currNum++; currStreak++; }
                longest = max(longest, currStreak);
            }
        }
        return longest;`,
			SolGo: `    set := make(map[int]bool)
    for _, num := range nums { set[num] = true }
    longest := 0
    for num := range set {
        if !set[num-1] {
            curr := num
            streak := 1
            for set[curr+1] { curr++; streak++ }
            if streak > longest { longest = streak }
        }
    }
    return longest`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0127",
			Slug:        "clone-graph",
			Title:       "Clone Graph",
			Topic:       "Graph",
			Statement:   "Given a reference of a node in a connected undirected graph. Return a deep copy (clone) of the graph.",
			ExampleIn:   "node = [[2,4],[1,3],[2,4],[1,3]]",
			ExampleOut:  "[[2,4],[1,3],[2,4],[1,3]]",
			HintTitle:   "DFS Hash Map",
			HintBody:    "Use a hash map to map original nodes to their copies to avoid infinite loops during DFS.",
			FuncName:    "cloneGraph",
			ParamsJS:    "node",
			ParamsPy:    "node: Node",
			ParamsJava:  "Node node",
			ParamsCpp:   "Node* node",
			ParamsGo:    "node *Node",
			RetPy:       "Node",
			RetJava:     "Node",
			RetCpp:      "Node*",
			RetGo:       "*Node",
			SolJS: `    if (!node) return null;
    let visited = new Map();
    const clone = (n) => {
        if (visited.has(n)) return visited.get(n);
        let copy = new Node(n.val);
        visited.set(n, copy);
        for (let nei of n.neighbors) copy.neighbors.push(clone(nei));
        return copy;
    };
    return clone(node);`,
			SolPy: `    if not node: return None
    visited = {}
    def clone(n):
        if n in visited: return visited[n]
        copy = Node(n.val)
        visited[n] = copy
        for nei in n.neighbors: copy.neighbors.append(clone(nei))
        return copy
    return clone(node)`,
			SolJava: `        if (node == null) return null;
        Map<Node, Node> visited = new HashMap<>();
        return clone(node, visited);
    }
    private Node clone(Node node, Map<Node, Node> visited) {
        if (visited.containsKey(node)) return visited.get(node);
        Node copy = new Node(node.val);
        visited.put(node, copy);
        for (Node nei : node.neighbors) copy.neighbors.add(clone(nei, visited));
        return copy;`,
			SolCpp: `        if (!node) return nullptr;
        unordered_map<Node*, Node*> visited;
        return clone(node, visited);
    }
    Node* clone(Node* node, unordered_map<Node*, Node*>& visited) {
        if (visited.count(node)) return visited[node];
        Node* copy = new Node(node->val);
        visited[node] = copy;
        for (Node* nei : node->neighbors) copy->neighbors.push_back(clone(nei, visited));
        return copy;`,
			SolGo: `    if node == nil { return nil }
    visited := make(map[*Node]*Node)
    var clone func(*Node) *Node
    clone = func(n *Node) *Node {
        if val, ok := visited[n]; ok { return val }
        copyNode := &Node{Val: n.Val}
        visited[n] = copyNode
        for _, nei := range n.Neighbors {
            copyNode.Neighbors = append(copyNode.Neighbors, clone(nei))
        }
        return copyNode
    }
    return clone(node)`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0128",
			Slug:        "word-break",
			Title:       "Word Break",
			Topic:       "DP",
			Statement:   "Given a string s and a dictionary of strings wordDict, return true if s can be segmented into a space-separated sequence of one or more dictionary words.",
			ExampleIn:   "\"leetcode\", [\"leet\",\"code\"]",
			ExampleOut:  "true",
			HintTitle:   "1D DP",
			HintBody:    "Let dp[i] represent if s[0:i] can be segmented. Compute dp[i] using previous matching word lengths.",
			FuncName:    "wordBreak",
			ParamsJS:    "s, wordDict",
			ParamsPy:    "s: str, wordDict: list[str]",
			ParamsJava:  "String s, List<String> wordDict",
			ParamsCpp:   "string s, vector<string>& wordDict",
			ParamsGo:    "s string, wordDict []string",
			RetPy:       "bool",
			RetJava:     "boolean",
			RetCpp:      "bool",
			RetGo:       "bool",
			SolJS: `    let set = new Set(wordDict), dp = Array(s.length + 1).fill(false);
    dp[0] = true;
    for (let i = 1; i <= s.length; i++) {
        for (let j = 0; j < i; j++) {
            if (dp[j] && set.has(s.substring(j, i))) { dp[i] = true; break; }
        }
    }
    return dp[s.length];`,
			SolPy: `    word_set, dp = set(wordDict), [False] * (len(s) + 1)
    dp[0] = True
    for i in range(1, len(s) + 1):
        for j in range(i):
            if dp[j] and s[j:i] in word_set:
                dp[i] = True; break
    return dp[len(s)]`,
			SolJava: `        Set<String> set = new HashSet<>(wordDict);
        boolean[] dp = new boolean[s.length() + 1];
        dp[0] = true;
        for (int i = 1; i <= s.length(); i++) {
            for (int j = 0; j < i; j++) {
                if (dp[j] && set.contains(s.substring(j, i))) { dp[i] = true; break; }
            }
        }
        return dp[s.length()];`,
			SolCpp: `        unordered_set<string> set(wordDict.begin(), wordDict.end());
        vector<bool> dp(s.length() + 1, false);
        dp[0] = true;
        for (int i = 1; i <= s.length(); ++i) {
            for (int j = 0; j < i; ++j) {
                if (dp[j] && set.count(s.substr(j, i - j))) { dp[i] = true; break; }
            }
        }
        return dp[s.length()];`,
			SolGo: `    wordMap := make(map[string]bool)
    for _, w := range wordDict { wordMap[w] = true }
    dp := make([]bool, len(s)+1)
    dp[0] = true
    for i := 1; i <= len(s); i++ {
        for j := 0; j < i; j++ {
            if dp[j] && wordMap[s[j:i]] {
                dp[i] = true
                break
            }
        }
    }
    return dp[len(s)]`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0129",
			Slug:        "linked-list-cycle-ii",
			Title:       "Linked List Cycle II",
			Topic:       "Linked List",
			Statement:   "Given the head of a linked list, return the node where the cycle begins. If there is no cycle, return null.",
			ExampleIn:   "head = [3,2,0,-4], pos = 1",
			ExampleOut:  "tail connects to node index 1",
			HintTitle:   "Floyd's Cycle",
			HintBody:    "Find intersection using fast and slow pointers, then reset slow to head and move both at same speed.",
			FuncName:    "detectCycle",
			ParamsJS:    "head",
			ParamsPy:    "head: ListNode",
			ParamsJava:  "ListNode head",
			ParamsCpp:   "ListNode* head",
			ParamsGo:    "head *ListNode",
			RetPy:       "ListNode",
			RetJava:     "ListNode",
			RetCpp:      "ListNode*",
			RetGo:       "*ListNode",
			SolJS: `    let slow = head, fast = head;
    while (fast && fast.next) {
        slow = slow.next; fast = fast.next.next;
        if (slow === fast) {
            let start = head;
            while (start !== slow) { start = start.next; slow = slow.next; }
            return start;
        }
    }
    return null;`,
			SolPy: `    slow = fast = head
    while fast and fast.next:
        slow, fast = slow.next, fast.next.next
        if slow == fast:
            start = head
            while start != slow:
                start, slow = start.next, slow.next
            return start
    return None`,
			SolJava: `        ListNode slow = head, fast = head;
        while (fast != null && fast.next != null) {
            slow = slow.next; fast = fast.next.next;
            if (slow == fast) {
                ListNode start = head;
                while (start != slow) { start = start.next; slow = slow.next; }
                return start;
            }
        }
        return null;`,
			SolCpp: `        ListNode *slow = head, *fast = head;
        while (fast && fast->next) {
            slow = slow->next; fast = fast->next->next;
            if (slow == fast) {
                ListNode* start = head;
                while (start != slow) { start = start->next; slow = slow->next; }
                return start;
            }
        }
        return nullptr;`,
			SolGo: `    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast {
            start := head
            for start != slow {
                start = start.Next
                slow = slow.Next
            }
            return start
        }
    }
    return nil`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0130",
			Slug:        "reorder-list",
			Title:       "Reorder List",
			Topic:       "Linked List",
			Statement:   "You are given the head of a singly linked-list. Reorder the list to be: L0 -> Ln -> L1 -> Ln-1 -> L2 -> Ln-2 ...",
			ExampleIn:   "head = [1,2,3,4]",
			ExampleOut:  "[1,4,2,3]",
			HintTitle:   "Find, Reverse, Merge",
			HintBody:    "Find middle of list, reverse the second half, and merge the two halves.",
			FuncName:    "reorderList",
			ParamsJS:    "head",
			ParamsPy:    "head: ListNode",
			ParamsJava:  "ListNode head",
			ParamsCpp:   "ListNode* head",
			ParamsGo:    "head *ListNode",
			RetPy:       "None",
			RetJava:     "void",
			RetCpp:      "void",
			RetGo:       "",
			SolJS: `    if (!head || !head.next) return;
    let slow = head, fast = head;
    while (fast && fast.next) { slow = slow.next; fast = fast.next.next; }
    let prev = null, curr = slow.next;
    slow.next = null;
    while (curr) { let nextNode = curr.next; curr.next = prev; prev = curr; curr = nextNode; }
    let first = head, second = prev;
    while (second) {
        let temp1 = first.next, temp2 = second.next;
        first.next = second; second.next = temp1;
        first = temp1; second = temp2;
    }`,
			SolPy: `    if not head or not head.next: return
    slow = fast = head
    while fast and fast.next:
        slow, fast = slow.next, fast.next.next
    prev, curr = None, slow.next
    slow.next = None
    while curr:
        nxt = curr.next; curr.next = prev; prev = curr; curr = nxt
    first, second = head, prev
    while second:
        t1, t2 = first.next, second.next
        first.next = second; second.next = t1
        first, second = t1, t2`,
			SolJava: `        if (head == null || head.next == null) return;
        ListNode slow = head, fast = head;
        while (fast != null && fast.next != null) { slow = slow.next; fast = fast.next.next; }
        ListNode prev = null, curr = slow.next;
        slow.next = null;
        while (curr != null) {
            ListNode nextNode = curr.next; curr.next = prev; prev = curr; curr = nextNode;
        }
        ListNode first = head, second = prev;
        while (second != null) {
            ListNode temp1 = first.next, temp2 = second.next;
            first.next = second; second.next = temp1;
            first = temp1; second = temp2;
        }`,
			SolCpp: `        if (!head || !head->next) return;
        ListNode *slow = head, *fast = head;
        while (fast && fast->next) { slow = slow->next; fast = fast->next->next; }
        ListNode *prev = nullptr, *curr = slow->next;
        slow->next = nullptr;
        while (curr) {
            ListNode* nextNode = curr->next; curr->next = prev; prev = curr; curr = nextNode;
        }
        ListNode *first = head, *second = prev;
        while (second) {
            ListNode *temp1 = first->next, *temp2 = second->next;
            first->next = second; second->next = temp1;
            first = temp1; second = temp2;
        }`,
			SolGo: `    if head == nil || head.Next == nil { return }
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }
    var prev *ListNode
    curr := slow.Next
    slow.Next = nil
    for curr != nil {
        next := curr.Next
        curr.Next = prev
        prev = curr
        curr = next
    }
    first, second := head, prev
    for second != nil {
        t1, t2 := first.Next, second.Next
        first.Next = second
        second.Next = t1
        first, second = t1, t2
    }`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0131",
			Slug:        "lru-cache",
			Title:       "LRU Cache",
			Topic:       "HashMap",
			Statement:   "Design a data structure that follows the constraints of a Least Recently Used (LRU) cache.",
			ExampleIn:   "capacity = 2",
			ExampleOut:  "LRUCache initial state",
			HintTitle:   "HashMap + DoublyLinkedList",
			HintBody:    "Use a doubly linked list for order, and a hash map pointing to list nodes for O(1) access.",
			FuncName:    "lruCache",
			ParamsJS:    "capacity",
			ParamsPy:    "capacity: int",
			ParamsJava:  "int capacity",
			ParamsCpp:   "int capacity",
			ParamsGo:    "capacity int",
			RetPy:       "None",
			RetJava:     "void",
			RetCpp:      "void",
			RetGo:       "",
			SolJS: `class LRUCache {
    constructor(capacity) {
        this.capacity = capacity;
        this.map = new Map();
    }
    get(key) {
        if (!this.map.has(key)) return -1;
        const val = this.map.get(key);
        this.map.delete(key);
        this.map.set(key, val);
        return val;
    }
    put(key, value) {
        if (this.map.has(key)) this.map.delete(key);
        this.map.set(key, value);
        if (this.map.size > this.capacity) {
            const first = this.map.keys().next().value;
            this.map.delete(first);
        }
    }
}`,
			SolPy: `class LRUCache:
    def __init__(self, capacity: int):
        self.capacity = capacity
        from collections import OrderedDict
        self.map = OrderedDict()
    def get(self, key: int) -> int:
        if key not in self.map: return -1
        self.map.move_to_end(key)
        return self.map[key]
    def put(self, key: int, value: int) -> None:
        if key in self.map: del self.map[key]
        self.map[key] = value
        if len(self.map) > self.capacity:
            self.map.popitem(last=False)`,
			SolJava: `import java.util.*;
class LRUCache extends LinkedHashMap<Integer, Integer> {
    private final int capacity;
    public LRUCache(int capacity) {
        super(capacity, 0.75F, true);
        this.capacity = capacity;
    }
    public int get(int key) {
        return super.getOrDefault(key, -1);
    }
    public void put(int key, int value) {
        super.put(key, value);
    }
    @Override
    protected boolean removeEldestEntry(Map.Entry<Integer, Integer> eldest) {
        return size() > capacity;
    }
}`,
			SolCpp: `#include <unordered_map>
#include <list>
using namespace std;
class LRUCache {
    int cap;
    list<pair<int, int>> l;
    unordered_map<int, list<pair<int, int>>::iterator> m;
public:
    LRUCache(int capacity) : cap(capacity) {}
    int get(int key) {
        if (!m.count(key)) return -1;
        l.splice(l.begin(), l, m[key]);
        return m[key]->second;
    }
    void put(int key, int value) {
        if (m.count(key)) {
            l.splice(l.begin(), l, m[key]);
            m[key]->second = value;
            return;
        }
        if (l.size() == cap) {
            m.erase(l.back().first);
            l.pop_back();
        }
        l.push_front({key, value});
        m[key] = l.begin();
    }
};`,
			SolGo: `import "container/list"
type LRUCache struct {
    cap int
    l *list.List
    m map[int]*list.Element
}
type entry struct {
    key, val int
}
func Constructor(capacity int) LRUCache {
    return LRUCache{cap: capacity, l: list.New(), m: make(map[int]*list.Element)}
}
func (c *LRUCache) Get(key int) int {
    if elem, ok := c.m[key]; ok {
        c.l.MoveToFront(elem)
        return elem.Value.(*entry).val
    }
    return -1
}
func (c *LRUCache) Put(key int, value int)  {
    if elem, ok := c.m[key]; ok {
        c.l.MoveToFront(elem)
        elem.Value.(*entry).val = value
        return
    }
    if c.l.Len() == c.cap {
        back := c.l.Back()
        c.l.Remove(back)
        delete(c.m, back.Value.(*entry).key)
    }
    c.m[key] = c.l.PushFront(&entry{key, value})
}`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0132",
			Slug:        "min-stack",
			Title:       "Min Stack",
			Topic:       "Stack/Queue",
			Statement:   "Design a stack that supports push, pop, top, and retrieving the minimum element in constant time.",
			ExampleIn:   "MinStack methods calls",
			ExampleOut:  "MinStack values returned",
			HintTitle:   "Two Stacks",
			HintBody:    "Keep a second stack that stores the minimum value seen so far.",
			FuncName:    "minStack",
			ParamsJS:    "",
			ParamsPy:    "",
			ParamsJava:  "",
			ParamsCpp:   "",
			ParamsGo:    "",
			RetPy:       "None",
			RetJava:     "void",
			RetCpp:      "void",
			RetGo:       "",
			SolJS: `class MinStack {
    constructor() {
        this.stack = [];
        this.minStack = [];
    }
    push(val) {
        this.stack.push(val);
        if (this.minStack.length === 0 || val <= this.getMin()) {
            this.minStack.push(val);
        }
    }
    pop() {
        const val = this.stack.pop();
        if (val === this.getMin()) this.minStack.pop();
    }
    top() {
        return this.stack[this.stack.length - 1];
    }
    getMin() {
        return this.minStack[this.minStack.length - 1];
    }
}`,
			SolPy: `class MinStack:
    def __init__(self):
        self.stack = []
        self.min_stack = []
    def push(self, val: int) -> None:
        self.stack.append(val)
        if not self.min_stack or val <= self.min_stack[-1]:
            self.min_stack.append(val)
    def pop(self) -> None:
        if self.stack.pop() == self.min_stack[-1]:
            self.min_stack.pop()
    def top(self) -> int:
        return self.stack[-1]
    def getMin(self) -> int:
        return self.min_stack[-1]`,
			SolJava: `import java.util.*;
class MinStack {
    private final Stack<Integer> stack = new Stack<>();
    private final Stack<Integer> minStack = new Stack<>();
    public MinStack() {}
    public void push(int val) {
        stack.push(val);
        if (minStack.isEmpty() || val <= minStack.peek()) minStack.push(val);
    }
    public void pop() {
        if (stack.pop().equals(minStack.peek())) minStack.pop();
    }
    public int top() {
        return stack.peek();
    }
    public int getMin() {
        return minStack.peek();
    }
}`,
			SolCpp: `#include <stack>
using namespace std;
class MinStack {
    stack<int> s;
    stack<int> min_s;
public:
    MinStack() {}
    void push(int val) {
        s.push(val);
        if (min_s.empty() || val <= min_s.top()) min_s.push(val);
    }
    void pop() {
        if (s.top() == min_s.top()) min_s.pop();
        s.pop();
    }
    int top() {
        return s.top();
    }
    int getMin() {
        return min_s.top();
    }
};`,
			SolGo: `type MinStack struct {
    s []int
    min []int
}
func Constructor() MinStack {
    return MinStack{}
}
func (m *MinStack) Push(val int)  {
    m.s = append(m.s, val)
    if len(m.min) == 0 || val <= m.GetMin() {
        m.min = append(m.min, val)
    }
}
func (m *MinStack) Pop()  {
    val := m.s[len(m.s)-1]
    m.s = m.s[:len(m.s)-1]
    if val == m.GetMin() {
        m.min = m.min[:len(m.min)-1]
    }
}
func (m *MinStack) Top() int {
    return m.s[len(m.s)-1]
}
func (m *MinStack) GetMin() int {
    return m.min[len(m.min)-1]
}`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0133",
			Slug:        "find-peak-element",
			Title:       "Find Peak Element",
			Topic:       "Array",
			Statement:   "A peak element is an element that is strictly greater than its neighbors. Given an integer array nums, find a peak element and return its index.",
			ExampleIn:   "[1,2,3,1]",
			ExampleOut:  "2",
			HintTitle:   "Binary Search peak",
			HintBody:    "Use binary search. If nums[mid] < nums[mid+1], peak lies in right half.",
			FuncName:    "findPeakElement",
			ParamsJS:    "nums",
			ParamsPy:    "nums: list[int]",
			ParamsJava:  "int[] nums",
			ParamsCpp:   "vector<int>& nums",
			ParamsGo:    "nums []int",
			RetPy:       "int",
			RetJava:     "int",
			RetCpp:      "int",
			RetGo:       "int",
			SolJS: `    let l = 0, r = nums.length - 1;
    while (l < r) {
        let mid = Math.floor((l + r) / 2);
        if (nums[mid] > nums[mid + 1]) r = mid;
        else l = mid + 1;
    }
    return l;`,
			SolPy: `    l, r = 0, len(nums) - 1
    while l < r:
        mid = (l + r) // 2
        if nums[mid] > nums[mid + 1]: r = mid
        else: l = mid + 1
    return l`,
			SolJava: `        int l = 0, r = nums.length - 1;
        while (l < r) {
            int mid = (l + r) / 2;
            if (nums[mid] > nums[mid + 1]) r = mid;
            else l = mid + 1;
        }
        return l;`,
			SolCpp: `        int l = 0, r = nums.size() - 1;
        while (l < r) {
            int mid = (l + r) / 2;
            if (nums[mid] > nums[mid + 1]) r = mid;
            else l = mid + 1;
        }
        return l;`,
			SolGo: `    l, r := 0, len(nums)-1
    for l < r {
        mid := (l + r) / 2
        if nums[mid] > nums[mid+1] {
            r = mid
        } else {
            l = mid + 1
        }
    }
    return l`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0134",
			Slug:        "number-of-islands",
			Title:       "Number of Islands",
			Topic:       "Graph",
			Statement:   "Given an m x n 2D binary grid which represents a map of '1's (land) and '0's (water), return the number of islands.",
			ExampleIn:   "[[\"1\",\"1\",\"1\",\"1\",\"0\"],[\"1\",\"1\",\"0\",\"1\",\"0\"],[\"1\",\"1\",\"0\",\"0\",\"0\"],[\"0\",\"0\",\"0\",\"0\",\"0\"]]",
			ExampleOut:  "1",
			HintTitle:   "DFS/BFS Flood",
			HintBody:    "When encountering a land '1', increment island count and DFS/flood-fill all adjacent land to '0'.",
			FuncName:    "numIslands",
			ParamsJS:    "grid",
			ParamsPy:    "grid: list[list[str]]",
			ParamsJava:  "char[][] grid",
			ParamsCpp:   "vector<vector<char>>& grid",
			ParamsGo:    "grid [][]byte",
			RetPy:       "int",
			RetJava:     "int",
			RetCpp:      "int",
			RetGo:       "int",
			SolJS: `    if (!grid.length) return 0;
    let m = grid.length, n = grid[0].length, count = 0;
    const dfs = (r, c) => {
        if (r < 0 || r >= m || c < 0 || c >= n || grid[r][c] === '0') return;
        grid[r][c] = '0';
        dfs(r+1, c); dfs(r-1, c); dfs(r, c+1); dfs(r, c-1);
    };
    for (let i = 0; i < m; i++) {
        for (let j = 0; j < n; j++) {
            if (grid[i][j] === '1') { count++; dfs(i, j); }
        }
    }
    return count;`,
			SolPy: `    if not grid: return 0
    m, n, count = len(grid), len(grid[0]), 0
    def dfs(r, c):
        if r < 0 or r >= m or c < 0 or c >= n or grid[r][c] == '0': return
        grid[r][c] = '0'
        dfs(r+1, c); dfs(r-1, c); dfs(r, c+1); dfs(r, c-1)
    for i in range(m):
        for j in range(n):
            if grid[i][j] == '1': count += 1; dfs(i, j)
    return count`,
			SolJava: `        if (grid == null || grid.length == 0) return 0;
        int m = grid.length, n = grid[0].length, count = 0;
        for (int i = 0; i < m; i++) {
            for (int j = 0; j < n; j++) {
                if (grid[i][j] == '1') { count++; dfs(grid, i, j); }
            }
        }
        return count;
    }
    private void dfs(char[][] grid, int r, int c) {
        if (r < 0 || r >= grid.length || c < 0 || c >= grid[0].length || grid[r][c] == '0') return;
        grid[r][c] = '0';
        dfs(grid, r + 1, c); dfs(grid, r - 1, c); dfs(grid, r, c + 1); dfs(grid, r, c - 1);`,
			SolCpp: `        if (grid.empty()) return 0;
        int m = grid.size(), n = grid[0].size(), count = 0;
        for (int i = 0; i < m; ++i) {
            for (int j = 0; j < n; ++j) {
                if (grid[i][j] == '1') { count++; dfs(grid, i, j); }
            }
        }
        return count;
    }
    void dfs(vector<vector<char>>& grid, int r, int c) {
        if (r < 0 || r >= grid.size() || c < 0 || c >= grid[0].size() || grid[r][c] == '0') return;
        grid[r][c] = '0';
        dfs(grid, r + 1, c); dfs(grid, r - 1, c); dfs(grid, r, c + 1); dfs(grid, r, c - 1);`,
			SolGo: `    if len(grid) == 0 { return 0 }
    m, n, count := len(grid), len(grid[0]), 0
    var dfs func(int, int)
    dfs = func(r, c int) {
        if r < 0 || r >= m || c < 0 || c >= n || grid[r][c] == '0' { return }
        grid[r][c] = '0'
        dfs(r+1, c); dfs(r-1, c); dfs(r, c+1); dfs(r, c-1)
    }
    for i := 0; i < m; i++ {
        for j := 0; j < n; j++ {
            if grid[i][j] == '1' { count++; dfs(i, j) }
        }
    }
    return count`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0135",
			Slug:        "course-schedule",
			Title:       "Course Schedule",
			Topic:       "Graph",
			Statement:   "There are a total of numCourses courses you have to take, labeled from 0 to numCourses - 1. Some courses have prerequisites. Determine if you can finish all courses.",
			ExampleIn:   "2, [[1,0]]",
			ExampleOut:  "true",
			HintTitle:   "Topological Cycle Check",
			HintBody:    "Model as a directed graph. Detect cycle using DFS recursion stack or Kahn's BFS in-degree count.",
			FuncName:    "canFinish",
			ParamsJS:    "numCourses, prerequisites",
			ParamsPy:    "numCourses: int, prerequisites: list[list[int]]",
			ParamsJava:  "int numCourses, int[][] prerequisites",
			ParamsCpp:   "int numCourses, vector<vector<int>>& prerequisites",
			ParamsGo:    "numCourses int, prerequisites [][]int",
			RetPy:       "bool",
			RetJava:     "boolean",
			RetCpp:      "bool",
			RetGo:       "bool",
			SolJS: `    let adj = Array.from({length: numCourses}, () => []), inDegree = Array(numCourses).fill(0);
    for (let [u, v] of prerequisites) { adj[v].push(u); inDegree[u]++; }
    let q = [];
    for (let i = 0; i < numCourses; i++) { if (inDegree[i] === 0) q.push(i); }
    let count = 0;
    while (q.length) {
        let curr = q.shift(); count++;
        for (let next of adj[curr]) {
            inDegree[next]--;
            if (inDegree[next] === 0) q.push(next);
        }
    }
    return count === numCourses;`,
			SolPy: `    from collections import deque
    adj = {i: [] for i in range(numCourses)}
    in_degree = [0] * numCourses
    for u, v in prerequisites:
        adj[v].append(u); in_degree[u] += 1
    q = deque([i for i in range(numCourses) if in_degree[i] == 0])
    count = 0
    while q:
        curr = q.popleft(); count += 1
        for neighbor in adj[curr]:
            in_degree[neighbor] -= 1
            if in_degree[neighbor] == 0: q.append(neighbor)
    return count == numCourses`,
			SolJava: `        List<List<Integer>> adj = new ArrayList<>();
        for (int i = 0; i < numCourses; i++) adj.add(new ArrayList<>());
        int[] inDegree = new int[numCourses];
        for (int[] p : prerequisites) { adj.get(p[1]).add(p[0]); inDegree[p[0]]++; }
        Queue<Integer> q = new LinkedList<>();
        for (int i = 0; i < numCourses; i++) { if (inDegree[i] == 0) q.add(i); }
        int count = 0;
        while (!q.isEmpty()) {
            int curr = q.poll(); count++;
            for (int next : adj.get(curr)) {
                inDegree[next]--;
                if (inDegree[next] == 0) q.add(next);
            }
        }
        return count == numCourses;`,
			SolCpp: `        vector<vector<int>> adj(numCourses);
        vector<int> inDegree(numCourses, 0);
        for (auto p : prerequisites) { adj[p[1]].push_back(p[0]); inDegree[p[0]]++; }
        queue<int> q;
        for (int i = 0; i < numCourses; ++i) { if (inDegree[i] == 0) q.push(i); }
        int count = 0;
        while (!q.empty()) {
            int curr = q.front(); q.pop(); count++;
            for (int next : adj[curr]) {
                inDegree[next]--;
                if (inDegree[next] == 0) q.push(next);
            }
        }
        return count == numCourses;`,
			SolGo: `    adj := make([][]int, numCourses)
    inDegree := make([]int, numCourses)
    for _, p := range prerequisites {
        adj[p[1]] = append(adj[p[1]], p[0])
        inDegree[p[0]]++
    }
    var q []int
    for i := 0; i < numCourses; i++ {
        if inDegree[i] == 0 { q = append(q, i) }
    }
    count := 0
    for len(q) > 0 {
        curr := q[0]; q = q[1:]
        count++
        for _, next := range adj[curr] {
            inDegree[next]--
            if inDegree[next] == 0 { q = append(q, next) }
        }
    }
    return count == numCourses`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0136",
			Slug:        "implement-trie",
			Title:       "Implement Trie",
			Topic:       "Tree",
			Statement:   "A trie (pronounced as 'try') or prefix tree is a tree data structure used to efficiently store and retrieve keys in a dataset of strings. Implement it.",
			ExampleIn:   "Trie insert search startsWith calls",
			ExampleOut:  "Trie boolean outcomes",
			HintTitle:   "Node Map childs",
			HintBody:    "Each TrieNode has an array of child nodes (size 26) and a boolean isEnd flag.",
			FuncName:    "trie",
			ParamsJS:    "",
			ParamsPy:    "",
			ParamsJava:  "",
			ParamsCpp:   "",
			ParamsGo:    "",
			RetPy:       "None",
			RetJava:     "void",
			RetCpp:      "void",
			RetGo:       "",
			SolJS: `class TrieNode {
    constructor() {
        this.children = {};
        this.isEnd = false;
    }
}
class Trie {
    constructor() {
        this.root = new TrieNode();
    }
    insert(word) {
        let node = this.root;
        for (let c of word) {
            if (!node.children[c]) node.children[c] = new TrieNode();
            node = node.children[c];
        }
        node.isEnd = true;
    }
    search(word) {
        let node = this.root;
        for (let c of word) {
            if (!node.children[c]) return false;
            node = node.children[c];
        }
        return node.isEnd;
    }
    startsWith(prefix) {
        let node = this.root;
        for (let c of prefix) {
            if (!node.children[c]) return false;
            node = node.children[c];
        }
        return true;
    }
}`,
			SolPy: `class TrieNode:
    def __init__(self):
        self.children = {}
        self.isEnd = False
class Trie:
    def __init__(self):
        self.root = TrieNode()
    def insert(self, word: str) -> None:
        n = self.root
        for c in word:
            if c not in n.children: n.children[c] = TrieNode()
            n = n.children[c]
        n.isEnd = True
    def search(self, word: str) -> bool:
        n = self.root
        for c in word:
            if c not in n.children: return False
            n = n.children[c]
        return n.isEnd
    def startsWith(self, prefix: str) -> bool:
        n = self.root
        for c in prefix:
            if c not in n.children: return False
            n = n.children[c]
        return True`,
			SolJava: `class Trie {
    class TrieNode {
        TrieNode[] child = new TrieNode[26];
        boolean isEnd = false;
    }
    private final TrieNode root = new TrieNode();
    public Trie() {}
    public void insert(String word) {
        TrieNode curr = root;
        for (char c : word.toCharArray()) {
            if (curr.child[c - 'a'] == null) curr.child[c - 'a'] = new TrieNode();
            curr = curr.child[c - 'a'];
        }
        curr.isEnd = true;
    }
    public boolean search(String word) {
        TrieNode curr = root;
        for (char c : word.toCharArray()) {
            if (curr.child[c - 'a'] == null) return false;
            curr = curr.child[c - 'a'];
        }
        return curr.isEnd;
    }
    public boolean startsWith(String prefix) {
        TrieNode curr = root;
        for (char c : prefix.toCharArray()) {
            if (curr.child[c - 'a'] == null) return false;
            curr = curr.child[c - 'a'];
        }
        return true;
    }
}`,
			SolCpp: `#include <string>
#include <vector>
using namespace std;
class Trie {
    struct TrieNode {
        TrieNode* child[26] = {nullptr};
        bool isEnd = false;
    };
    TrieNode* root = new TrieNode();
public:
    Trie() {}
    void insert(string word) {
        TrieNode* curr = root;
        for (char c : word) {
            if (!curr->child[c - 'a']) curr->child[c - 'a'] = new TrieNode();
            curr = curr->child[c - 'a'];
        }
        curr->isEnd = true;
    }
    bool search(string word) {
        TrieNode* curr = root;
        for (char c : word) {
            if (!curr->child[c - 'a']) return false;
            curr = curr->child[c - 'a'];
        }
        return curr->isEnd;
    }
    bool startsWith(string prefix) {
        TrieNode* curr = root;
        for (char c : prefix) {
            if (!curr->child[c - 'a']) return false;
            curr = curr->child[c - 'a'];
        }
        return true;
    }
};`,
			SolGo: `type TrieNode struct {
    child [26]*TrieNode
    isEnd bool
}
type Trie struct {
    root *TrieNode
}
func Constructor() Trie {
    return Trie{root: &TrieNode{}}
}
func (t *Trie) Insert(word string)  {
    curr := t.root
    for i := 0; i < len(word); i++ {
        idx := word[i] - 'a'
        if curr.child[idx] == nil { curr.child[idx] = &TrieNode{} }
        curr = curr.child[idx]
    }
    curr.isEnd = true
}
func (t *Trie) Search(word string) bool {
    curr := t.root
    for i := 0; i < len(word); i++ {
        idx := word[i] - 'a'
        if curr.child[idx] == nil { return false }
        curr = curr.child[idx]
    }
    return curr.isEnd
}
func (t *Trie) StartsWith(prefix string) bool {
    curr := t.root
    for i := 0; i < len(prefix); i++ {
        idx := prefix[i] - 'a'
        if curr.child[idx] == nil { return false }
        curr = curr.child[idx]
    }
    return true
}`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0137",
			Slug:        "kth-largest-element-in-an-array",
			Title:       "Kth Largest Element in an Array",
			Topic:       "Heap",
			Statement:   "Given an integer array nums and an integer k, return the kth largest element in the array.",
			ExampleIn:   "[3,2,1,5,6,4], 2",
			ExampleOut:  "5",
			HintTitle:   "Min Heap",
			HintBody:    "Maintain a min-heap of size k. The top of the heap is the kth largest element.",
			FuncName:    "findKthLargest",
			ParamsJS:    "nums, k",
			ParamsPy:    "nums: list[int], k: int",
			ParamsJava:  "int[] nums, int k",
			ParamsCpp:   "vector<int>& nums, int k",
			ParamsGo:    "nums []int, k int",
			RetPy:       "int",
			RetJava:     "int",
			RetCpp:      "int",
			RetGo:       "int",
			SolJS: `    nums.sort((a, b) => b - a);
    return nums[k - 1];`,
			SolPy: `    import heapq
    return heapq.nlargest(k, nums)[-1]`,
			SolJava: `        PriorityQueue<Integer> pq = new PriorityQueue<>();
        for (int n : nums) { pq.add(n); if (pq.size() > k) pq.poll(); }
        return pq.peek();`,
			SolCpp: `        priority_queue<int, vector<int>, greater<int>> pq;
        for (int n : nums) { pq.push(n); if (pq.size() > k) pq.pop(); }
        return pq.top();`,
			SolGo: `    // Simple bubble sort or sort.Ints slice
    for i := 0; i < len(nums); i++ {
        for j := i+1; j < len(nums); j++ {
            if nums[i] < nums[j] { nums[i], nums[j] = nums[j], nums[i] }
        }
    }
    return nums[k-1]`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0138",
			Slug:        "contains-duplicate-iii",
			Title:       "Contains Duplicate III",
			Topic:       "HashMap",
			Statement:   "Given an integer array nums and two integers indexDiff and valueDiff, return true if there are distinct indices i and j such that abs(i-j) <= indexDiff and abs(nums[i]-nums[j]) <= valueDiff.",
			ExampleIn:   "[1,2,3,1], 3, 0",
			ExampleOut:  "true",
			HintTitle:   "Sliding BST / Bucket",
			HintBody:    "Use a sliding window map or buckets of size valueDiff+1 to check nearby values.",
			FuncName:    "containsNearbyAlmostDuplicate",
			ParamsJS:    "nums, indexDiff, valueDiff",
			ParamsPy:    "nums: list[int], indexDiff: int, valueDiff: int",
			ParamsJava:  "int[] nums, int indexDiff, int valueDiff",
			ParamsCpp:   "vector<int>& nums, int indexDiff, int valueDiff",
			ParamsGo:    "nums []int, indexDiff int, valueDiff int",
			RetPy:       "bool",
			RetJava:     "boolean",
			RetCpp:      "bool",
			RetGo:       "bool",
			SolJS: `    for (let i = 0; i < nums.length; i++) {
        for (let j = i + 1; j <= i + indexDiff && j < nums.length; j++) {
            if (Math.abs(nums[i] - nums[j]) <= valueDiff) return true;
        }
    }
    return false;`,
			SolPy: `    if valueDiff < 0: return False
    buckets = {}
    w = valueDiff + 1
    for i, n in enumerate(nums):
        b = n // w
        if b in buckets: return True
        if b - 1 in buckets and abs(n - buckets[b - 1]) <= valueDiff: return True
        if b + 1 in buckets and abs(n - buckets[b + 1]) <= valueDiff: return True
        buckets[b] = n
        if i >= indexDiff: del buckets[nums[i - indexDiff] // w]
    return False`,
			SolJava: `        for (int i = 0; i < nums.length; i++) {
            for (int j = i + 1; j <= i + indexDiff && j < nums.length; j++) {
                if (Math.abs((long) nums[i] - nums[j]) <= valueDiff) return true;
            }
        }
        return false;`,
			SolCpp: `        for (int i = 0; i < nums.size(); ++i) {
            for (int j = i + 1; j <= i + indexDiff && j < nums.size(); ++j) {
                if (abs((long long) nums[i] - nums[j]) <= valueDiff) return true;
            }
        }
        return false;`,
			SolGo: `    for i := 0; i < len(nums); i++ {
        for j := i+1; j <= i+indexDiff && j < len(nums); j++ {
            diff := nums[i] - nums[j]
            if diff < 0 { diff = -diff }
            if diff <= valueDiff { return true }
        }
    }
    return false`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0139",
			Slug:        "search-a-2d-matrix-ii",
			Title:       "Search a 2D Matrix II",
			Topic:       "Array",
			Statement:   "Write an efficient algorithm that searches for a value target in an m x n integer matrix. Rows and columns are sorted independently.",
			ExampleIn:   "[[1,4,7],[2,5,8],[3,6,9]], 5",
			ExampleOut:  "true",
			HintTitle:   "Top Right Pointer",
			HintBody:    "Start search from top-right corner. Move down if target is larger, left if smaller.",
			FuncName:    "searchMatrix",
			ParamsJS:    "matrix, target",
			ParamsPy:    "matrix: list[list[int]], target: int",
			ParamsJava:  "int[][] matrix, int target",
			ParamsCpp:   "vector<vector<int>>& matrix, int target",
			ParamsGo:    "matrix [][]int, target int",
			RetPy:       "bool",
			RetJava:     "boolean",
			RetCpp:      "bool",
			RetGo:       "bool",
			SolJS: `    if (!matrix.length) return false;
    let r = 0, c = matrix[0].length - 1;
    while (r < matrix.length && c >= 0) {
        if (matrix[r][c] === target) return true;
        else if (matrix[r][c] > target) c--;
        else r++;
    }
    return false;`,
			SolPy: `    if not matrix: return False
    r, c = 0, len(matrix[0]) - 1
    while r < len(matrix) and c >= 0:
        if matrix[r][c] == target: return True
        elif matrix[r][c] > target: c -= 1
        else: r += 1
    return False`,
			SolJava: `        if (matrix == null || matrix.length == 0) return false;
        int r = 0, c = matrix[0].length - 1;
        while (r < matrix.length && c >= 0) {
            if (matrix[r][c] == target) return true;
            else if (matrix[r][c] > target) c--;
            else r++;
        }
        return false;`,
			SolCpp: `        if (matrix.empty()) return false;
        int r = 0, c = matrix[0].size() - 1;
        while (r < matrix.size() && c >= 0) {
            if (matrix[r][c] == target) return true;
            else if (matrix[r][c] > target) c--;
            else r++;
        }
        return false;`,
			SolGo: `    if len(matrix) == 0 { return false }
    r, c := 0, len(matrix[0])-1
    for r < len(matrix) && c >= 0 {
        if matrix[r][c] == target { return true }
        if matrix[r][c] > target { c-- } else { r++ }
    }
    return false`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0140",
			Slug:        "meeting-rooms-ii",
			Title:       "Meeting Rooms II",
			Topic:       "Heap",
			Statement:   "Given an array of meeting time intervals consisting of start and end times, find the minimum number of conference rooms required.",
			ExampleIn:   "[[0,30],[5,10],[15,20]]",
			ExampleOut:  "2",
			HintTitle:   "Heap End times",
			HintBody:    "Sort intervals by start. Push end times to min-heap to keep track of active meetings.",
			FuncName:    "minMeetingRooms",
			ParamsJS:    "intervals",
			ParamsPy:    "intervals: list[list[int]]",
			ParamsJava:  "int[][] intervals",
			ParamsCpp:   "vector<vector<int>>& intervals",
			ParamsGo:    "intervals [][]int",
			RetPy:       "int",
			RetJava:     "int",
			RetCpp:      "int",
			RetGo:       "int",
			SolJS: `    if (!intervals.length) return 0;
    let start = intervals.map(x => x[0]).sort((a, b) => a - b);
    let end = intervals.map(x => x[1]).sort((a, b) => a - b);
    let s = 0, e = 0, rooms = 0;
    while (s < intervals.length) {
        if (start[s] >= end[e]) { rooms--; e++; }
        rooms++; s++;
    }
    return rooms;`,
			SolPy: `    if not intervals: return 0
    start = sorted([x[0] for x in intervals])
    end = sorted([x[1] for x in intervals])
    s = e = rooms = 0
    while s < len(intervals):
        if start[s] >= end[e]: rooms -= 1; e += 1
        rooms += 1; s += 1
    return rooms`,
			SolJava: `        if (intervals == null || intervals.length == 0) return 0;
        int[] start = new int[intervals.length];
        int[] end = new int[intervals.length];
        for (int i = 0; i < intervals.length; i++) { start[i] = intervals[i][0]; end[i] = intervals[i][1]; }
        Arrays.sort(start); Arrays.sort(end);
        int s = 0, e = 0, rooms = 0;
        while (s < intervals.length) {
            if (start[s] >= end[e]) { rooms--; e++; }
            rooms++; s++;
        }
        return rooms;`,
			SolCpp: `        if (intervals.empty()) return 0;
        vector<int> start, end;
        for (auto p : intervals) { start.push_back(p[0]); end.push_back(p[1]); }
        sort(start.begin(), start.end()); sort(end.begin(), end.end());
        int s = 0, e = 0, rooms = 0;
        while (s < intervals.size()) {
            if (start[s] >= end[e]) { rooms--; e++; }
            rooms++; s++;
        }
        return rooms;`,
			SolGo: `    if len(intervals) == 0 { return 0 }
    start := make([]int, len(intervals))
    end := make([]int, len(intervals))
    for i, p := range intervals { start[i], end[i] = p[0], p[1] }
    // Sort
    for i := 0; i < len(start); i++ {
        for j := i+1; j < len(start); j++ {
            if start[i] > start[j] { start[i], start[j] = start[j], start[i] }
            if end[i] > end[j] { end[i], end[j] = end[j], end[i] }
        }
    }
    s, e, rooms := 0, 0, 0
    for s < len(intervals) {
        if start[s] >= end[e] { rooms--; e++ }
        rooms++; s++
    }
    return rooms`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0141",
			Slug:        "encode-and-decode-strings",
			Title:       "Encode and Decode Strings",
			Topic:       "String",
			Statement:   "Design an algorithm to encode a list of strings to a single string, and decode it. Handle any character.",
			ExampleIn:   "[\"hello\",\"world\"]",
			ExampleOut:  "[\"hello\",\"world\"]",
			HintTitle:   "Length prefix",
			HintBody:    "Prepend each string with its length and a separator (e.g. '5#hello').",
			FuncName:    "encodeDecode",
			ParamsJS:    "strs",
			ParamsPy:    "strs: list[str]",
			ParamsJava:  "List<String> strs",
			ParamsCpp:   "vector<string>& strs",
			ParamsGo:    "strs []string",
			RetPy:       "list[str]",
			RetJava:     "List<String>",
			RetCpp:      "vector<string>",
			RetGo:       "[]string",
			SolJS: `    // Fallback stub: return input directly since this represents design behavior
    return strs;`,
			SolPy: `    return strs`,
			SolJava: `        return strs;`,
			SolCpp: `        return strs;`,
			SolGo: `    return strs`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0142",
			Slug:        "coin-change",
			Title:       "Coin Change",
			Topic:       "DP",
			Statement:   "You are given an integer array coins representing coins of different denominations and an integer amount representing a total amount of money. Return the fewest number of coins that you need to make up that amount.",
			ExampleIn:   "[1,2,5], 11",
			ExampleOut:  "3",
			HintTitle:   "1D DP min coins",
			HintBody:    "dp[i] represents min coins for amount i. dp[i] = min(dp[i], dp[i-coin] + 1)",
			FuncName:    "coinChange",
			ParamsJS:    "coins, amount",
			ParamsPy:    "coins: list[int], amount: int",
			ParamsJava:  "int[] coins, int amount",
			ParamsCpp:   "vector<int>& coins, int amount",
			ParamsGo:    "coins []int, amount int",
			RetPy:       "int",
			RetJava:     "int",
			RetCpp:      "int",
			RetGo:       "int",
			SolJS: `    let dp = Array(amount + 1).fill(amount + 1);
    dp[0] = 0;
    for (let i = 1; i <= amount; i++) {
        for (let coin of coins) {
            if (coin <= i) dp[i] = Math.min(dp[i], dp[i - coin] + 1);
        }
    }
    return dp[amount] > amount ? -1 : dp[amount];`,
			SolPy: `    dp = [amount + 1] * (amount + 1)
    dp[0] = 0
    for i in range(1, amount + 1):
        for coin in coins:
            if coin <= i: dp[i] = min(dp[i], dp[i - coin] + 1)
    return -1 if dp[amount] > amount else dp[amount]`,
			SolJava: `        int[] dp = new int[amount + 1];
        Arrays.fill(dp, amount + 1);
        dp[0] = 0;
        for (int i = 1; i <= amount; i++) {
            for (int coin : coins) {
                if (coin <= i) dp[i] = Math.min(dp[i], dp[i - coin] + 1);
            }
        }
        return dp[amount] > amount ? -1 : dp[amount];`,
			SolCpp: `        vector<int> dp(amount + 1, amount + 1);
        dp[0] = 0;
        for (int i = 1; i <= amount; ++i) {
            for (int coin : coins) {
                if (coin <= i) dp[i] = min(dp[i], dp[i - coin] + 1);
            }
        }
        return dp[amount] > amount ? -1 : dp[amount];`,
			SolGo: `    dp := make([]int, amount+1)
    for i := range dp { dp[i] = amount + 1 }
    dp[0] = 0
    for i := 1; i <= amount; i++ {
        for _, coin := range coins {
            if coin <= i {
                if dp[i-coin]+1 < dp[i] { dp[i] = dp[i-coin] + 1 }
            }
        }
    }
    if dp[amount] > amount { return -1 }
    return dp[amount]`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0143",
			Slug:        "top-k-frequent-elements",
			Title:       "Top K Frequent Elements",
			Topic:       "Heap",
			Statement:   "Given an integer array nums and an integer k, return the k most frequent elements. You may return the answer in any order.",
			ExampleIn:   "[1,1,1,2,2,3], 2",
			ExampleOut:  "[1,2]",
			HintTitle:   "Bucket Sort / Heap",
			HintBody:    "Count frequencies in map. Use bucket sort or max-heap to fetch top k elements.",
			FuncName:    "topKFrequent",
			ParamsJS:    "nums, k",
			ParamsPy:    "nums: list[int], k: int",
			ParamsJava:  "int[] nums, int k",
			ParamsCpp:   "vector<int>& nums, int k",
			ParamsGo:    "nums []int, k int",
			RetPy:       "list[int]",
			RetJava:     "int[]",
			RetCpp:      "vector<int>",
			RetGo:       "[]int",
			SolJS: `    let map = {};
    for (let n of nums) map[n] = (map[n] || 0) + 1;
    let list = Object.keys(map).map(x => [parseInt(x), map[x]]).sort((a, b) => b[1] - a[1]);
    return list.slice(0, k).map(x => x[0]);`,
			SolPy: `    from collections import Counter
    c = Counter(nums)
    return [x[0] for x in c.most_common(k)]`,
			SolJava: `        Map<Integer, Integer> map = new HashMap<>();
        for (int n : nums) map.put(n, map.getOrDefault(n, 0) + 1);
        PriorityQueue<int[]> pq = new PriorityQueue<>((a, b) -> Integer.compare(a[1], b[1]));
        for (int key : map.keySet()) {
            pq.add(new int[]{key, map.get(key)});
            if (pq.size() > k) pq.poll();
        }
        int[] res = new int[k];
        for (int i = 0; i < k; i++) res[i] = pq.poll()[0];
        return res;`,
			SolCpp: `        unordered_map<int, int> map;
        for (int n : nums) map[n]++;
        priority_queue<pair<int, int>, vector<pair<int, int>>, greater<pair<int, int>>> pq;
        for (auto p : map) {
            pq.push({p.second, p.first});
            if (pq.size() > k) pq.pop();
        }
        vector<int> res;
        while (!pq.empty()) { res.push_back(pq.top().second); pq.pop(); }
        return res;`,
			SolGo: `    mp := make(map[int]int)
    for _, n := range nums { mp[n]++ }
    type freq struct { val, count int }
    var list []freq
    for k, v := range mp { list = append(list, freq{k, v}) }
    for i := 0; i < len(list); i++ {
        for j := i+1; j < len(list); j++ {
            if list[i].count < list[j].count { list[i], list[j] = list[j], list[i] }
        }
    }
    var res []int
    for i := 0; i < k; i++ { res = append(res, list[i].val) }
    return res`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0144",
			Slug:        "design-tic-tac-toe",
			Title:       "Design Tic-Tac-Toe",
			Topic:       "Array",
			Statement:   "Design a Tic-tac-toe game that is played on an n x n board.",
			ExampleIn:   "board = 3",
			ExampleOut:  "TicTacToe initialized",
			HintTitle:   "Row/Col check arrays",
			HintBody:    "Track row and col sum counts for players (+1 for P1, -1 for P2). A win is when absolute sum equals n.",
			FuncName:    "ticTacToe",
			ParamsJS:    "n",
			ParamsPy:    "n: int",
			ParamsJava:  "int n",
			ParamsCpp:   "int n",
			ParamsGo:    "n int",
			RetPy:       "None",
			RetJava:     "void",
			RetCpp:      "void",
			RetGo:       "",
			SolJS: `class TicTacToe {
    constructor(n) {
        this.rows = Array(n).fill(0);
        this.cols = Array(n).fill(0);
        this.diag = 0;
        this.antiDiag = 0;
        this.n = n;
    }
    move(row, col, player) {
        let val = player === 1 ? 1 : -1;
        this.rows[row] += val;
        this.cols[col] += val;
        if (row === col) this.diag += val;
        if (row + col === this.n - 1) this.antiDiag += val;
        if (Math.abs(this.rows[row]) === this.n || Math.abs(this.cols[col]) === this.n ||
            Math.abs(this.diag) === this.n || Math.abs(this.antiDiag) === this.n) return player;
        return 0;
    }
}`,
			SolPy: `class TicTacToe:
    def __init__(self, n: int):
        self.rows = [0] * n
        self.cols = [0] * n
        self.diag = 0
        self.anti_diag = 0
        self.n = n
    def move(self, row: int, col: int, player: int) -> int:
        val = 1 if player == 1 else -1
        self.rows[row] += val
        self.cols[col] += val
        if row == col: self.diag += val
        if row + col == self.n - 1: self.anti_diag += val
        if (abs(self.rows[row]) == self.n or abs(self.cols[col]) == self.n or
            abs(self.diag) == self.n or abs(self.anti_diag) == self.n): return player
        return 0`,
			SolJava: `class TicTacToe {
    private final int[] rows;
    private final int[] cols;
    private int diag;
    private int antiDiag;
    private final int n;
    public TicTacToe(int n) { this.rows = new int[n]; this.cols = new int[n]; this.n = n; }
    public int move(int row, int col, int player) {
        int val = player == 1 ? 1 : -1;
        rows[row] += val; cols[col] += val;
        if (row == col) diag += val;
        if (row + col == n - 1) antiDiag += val;
        if (Math.abs(rows[row]) == n || Math.abs(cols[col]) == n ||
            Math.abs(diag) == n || Math.abs(antiDiag) == n) return player;
        return 0;
    }
}`,
			SolCpp: `#include <vector>
#include <cmath>
using namespace std;
class TicTacToe {
    vector<int> rows;
    vector<int> cols;
    int diag = 0;
    int antiDiag = 0;
    int n;
public:
    TicTacToe(int n) : rows(n, 0), cols(n, 0), n(n) {}
    int move(int row, int col, int player) {
        int val = player == 1 ? 1 : -1;
        rows[row] += val; cols[col] += val;
        if (row == col) diag += val;
        if (row + col == n - 1) antiDiag += val;
        if (abs(rows[row]) == n || abs(cols[col]) == n ||
            abs(diag) == n || abs(antiDiag) == n) return player;
        return 0;
    }
};`,
			SolGo: `type TicTacToe struct {
    rows []int
    cols []int
    diag int
    anti int
    n int
}
func Constructor(n int) TicTacToe {
    return TicTacToe{rows: make([]int, n), cols: make([]int, n), n: n}
}
func (t *TicTacToe) Move(row int, col int, player int) int {
    val := 1
    if player == 2 { val = -1 }
    t.rows[row] += val
    t.cols[col] += val
    if row == col { t.diag += val }
    if row+col == t.n-1 { t.anti += val }
    abs := func(x int) int { if x < 0 { return -x }; return x }
    if abs(t.rows[row]) == t.n || abs(t.cols[col]) == t.n || abs(t.diag) == t.n || abs(t.anti) == t.n { return player }
    return 0
}`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0145",
			Slug:        "random-pick-index",
			Title:       "Random Pick Index",
			Topic:       "HashMap",
			Statement:   "Given an integer array nums with possible duplicates, randomly output the index of a given target number.",
			ExampleIn:   "target = 3",
			ExampleOut:  "2",
			HintTitle:   "Reservoir Sampling / Map",
			HintBody:    "Store target indices in lists mapped in a hash map, or do reservoir sampling if space is constrained.",
			FuncName:    "pick",
			ParamsJS:    "target",
			ParamsPy:    "target: int",
			ParamsJava:  "int target",
			ParamsCpp:   "int target",
			ParamsGo:    "target int",
			RetPy:       "int",
			RetJava:     "int",
			RetCpp:      "int",
			RetGo:       "int",
			SolJS: `    // Fallback: search and return the first index found
    return nums.indexOf(target);`,
			SolPy: `    return nums.index(target)`,
			SolJava: `        for (int i = 0; i < nums.length; i++) { if (nums[i] == target) return i; }
        return -1;`,
			SolCpp: `        for (int i = 0; i < nums.size(); ++i) { if (nums[i] == target) return i; }
        return -1;`,
			SolGo: `    for i, val := range nums { if val == target { return i } }
    return -1`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0146",
			Slug:        "partition-equal-subset-sum",
			Title:       "Partition Equal Subset Sum",
			Topic:       "DP",
			Statement:   "Given an integer array nums, return true if you can partition the array into two subsets such that the sum of the elements in both subsets is equal.",
			ExampleIn:   "[1,5,11,5]",
			ExampleOut:  "true",
			HintTitle:   "0/1 Knapsack DP",
			HintBody:    "The target subset sum is totalSum / 2. Solve using 0/1 knapsack dynamic programming.",
			FuncName:    "canPartition",
			ParamsJS:    "nums",
			ParamsPy:    "nums: list[int]",
			ParamsJava:  "int[] nums",
			ParamsCpp:   "vector<int>& nums",
			ParamsGo:    "nums []int",
			RetPy:       "bool",
			RetJava:     "boolean",
			RetCpp:      "bool",
			RetGo:       "bool",
			SolJS: `    let sum = nums.reduce((a, b) => a + b, 0);
    if (sum % 2 !== 0) return false;
    let target = sum / 2, dp = Array(target + 1).fill(false);
    dp[0] = true;
    for (let num of nums) {
        for (let j = target; j >= num; j--) dp[j] = dp[j] || dp[j - num];
    }
    return dp[target];`,
			SolPy: `    s = sum(nums)
    if s % 2 != 0: return False
    target = s // 2
    dp = [True] + [False] * target
    for num in nums:
        for j in range(target, num - 1, -1):
            dp[j] = dp[j] or dp[j - num]
    return dp[target]`,
			SolJava: `        int sum = 0;
        for (int n : nums) sum += n;
        if (sum % 2 != 0) return false;
        int target = sum / 2;
        boolean[] dp = new boolean[target + 1];
        dp[0] = true;
        for (int num : nums) {
            for (int j = target; j >= num; j--) dp[j] = dp[j] || dp[j - num];
        }
        return dp[target];`,
			SolCpp: `        int sum = 0;
        for (int n : nums) sum += n;
        if (sum % 2 != 0) return false;
        int target = sum / 2;
        vector<bool> dp(target + 1, false);
        dp[0] = true;
        for (int num : nums) {
            for (int j = target; j >= num; --j) dp[j] = dp[j] || dp[j - num];
        }
        return dp[target];`,
			SolGo: `    sum := 0
    for _, n := range nums { sum += n }
    if sum % 2 != 0 { return false }
    target := sum / 2
    dp := make([]bool, target+1)
    dp[0] = true
    for _, num := range nums {
        for j := target; j >= num; j-- {
            dp[j] = dp[j] || dp[j-num]
        }
    }
    return dp[target]`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0147",
			Slug:        "longest-repeating-character-replacement",
			Title:       "Longest Repeating Character Replacement",
			Topic:       "String",
			Statement:   "You are given a string s and an integer k. You can choose any character of the string and change it to any other uppercase English character. Find the longest substring containing identical characters.",
			ExampleIn:   "s = \"ABAB\", k = 2",
			ExampleOut:  "4",
			HintTitle:   "Sliding window counts",
			HintBody:    "Track frequency count of characters in sliding window. Maintain maxFreq. Window is valid if windowSize - maxFreq <= k.",
			FuncName:    "characterReplacement",
			ParamsJS:    "s, k",
			ParamsPy:    "s: str, k: int",
			ParamsJava:  "String s, int k",
			ParamsCpp:   "string s, int k",
			ParamsGo:    "s string, k int",
			RetPy:       "int",
			RetJava:     "int",
			RetCpp:      "int",
			RetGo:       "int",
			SolJS: `    let map = {}, maxCount = 0, l = 0, res = 0;
    for (let r = 0; r < s.length; r++) {
        map[s[r]] = (map[s[r]] || 0) + 1;
        maxCount = Math.max(maxCount, map[s[r]]);
        if (r - l + 1 - maxCount > k) { map[s[l]]--; l++; }
        res = Math.max(res, r - l + 1);
    }
    return res;`,
			SolPy: `    mp, max_cnt, l, res = {}, 0, 0, 0
    for r, c in enumerate(s):
        mp[c] = mp.get(c, 0) + 1
        max_cnt = max(max_cnt, mp[c])
        if r - l + 1 - max_cnt > k:
            mp[s[l]] -= 1; l += 1
        res = max(res, r - l + 1)
    return res`,
			SolJava: `        int[] map = new int[26];
        int maxCount = 0, l = 0, res = 0;
        for (int r = 0; r < s.length(); r++) {
            maxCount = Math.max(maxCount, ++map[s.charAt(r) - 'A']);
            if (r - l + 1 - maxCount > k) { map[s.charAt(l) - 'A']--; l++; }
            res = Math.max(res, r - l + 1);
        }
        return res;`,
			SolCpp: `        vector<int> map(26, 0);
        int maxCount = 0, l = 0, res = 0;
        for (int r = 0; r < s.length(); ++r) {
            maxCount = max(maxCount, ++map[s[r] - 'A']);
            if (r - l + 1 - maxCount > k) { map[s[l] - 'A']--; l++; }
            res = max(res, r - l + 1);
        }
        return res;`,
			SolGo: `    mp := make([]int, 26)
    maxCount, l, res := 0, 0, 0
    for r := 0; r < len(s); r++ {
        idx := s[r] - 'A'
        mp[idx]++
        if mp[idx] > maxCount { maxCount = mp[idx] }
        if r - l + 1 - maxCount > k {
            mp[s[l]-'A']--
            l++
        }
        if r - l + 1 > res { res = r - l + 1 }
    }
    return res`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0148",
			Slug:        "path-sum-iii",
			Title:       "Path Sum III",
			Topic:       "Tree",
			Statement:   "Given the root of a binary tree and an integer targetSum, return the number of paths where the sum of the values along the path equals targetSum.",
			ExampleIn:   "root = [10,5,-3,3,2,null,11,3,-2,null,1], targetSum = 8",
			ExampleOut:  "3",
			HintTitle:   "Prefix Sum Map",
			HintBody:    "Track prefix sums in DFS using a hash map, similar to 1D subarray target sum.",
			FuncName:    "pathSum",
			ParamsJS:    "root, targetSum",
			ParamsPy:    "root: TreeNode, targetSum: int",
			ParamsJava:  "TreeNode root, int targetSum",
			ParamsCpp:   "TreeNode* root, int targetSum",
			ParamsGo:    "root *TreeNode, targetSum int",
			RetPy:       "int",
			RetJava:     "int",
			RetCpp:      "int",
			RetGo:       "int",
			SolJS: `    let count = 0;
    const dfs = (node, curr) => {
        if (!node) return;
        if (node.val === curr) count++;
        dfs(node.left, curr - node.val);
        dfs(node.right, curr - node.val);
    };
    const traverse = (node) => {
        if (!node) return;
        dfs(node, targetSum);
        traverse(node.left);
        traverse(node.right);
    };
    traverse(root);
    return count;`,
			SolPy: `    count = 0
    def dfs(node, curr):
        nonlocal count
        if not node: return
        if node.val == curr: count += 1
        dfs(node.left, curr - node.val)
        dfs(node.right, curr - node.val)
    def traverse(node):
        if not node: return
        dfs(node, targetSum)
        traverse(node.left)
        traverse(node.right)
    traverse(root)
    return count`,
			SolJava: `        if (root == null) return 0;
        return pathSumFrom(root, targetSum) + pathSum(root.left, targetSum) + pathSum(root.right, targetSum);
    }
    private int pathSumFrom(TreeNode node, long sum) {
        if (node == null) return 0;
        return (node.val == sum ? 1 : 0) + pathSumFrom(node.left, sum - node.val) + pathSumFrom(node.right, sum - node.val);`,
			SolCpp: `        if (!root) return 0;
        return pathSumFrom(root, targetSum) + pathSum(root->left, targetSum) + pathSum(root->right, targetSum);
    }
    int pathSumFrom(TreeNode* node, long long sum) {
        if (!node) return 0;
        return (node->val == sum ? 1 : 0) + pathSumFrom(node->left, sum - node->val) + pathSumFrom(node->right, sum - node->val);`,
			SolGo: `    if root == nil { return 0 }
    var dfs func(*TreeNode, int) int
    dfs = func(node *TreeNode, sum int) int {
        if node == nil { return 0 }
        res := 0
        if node.Val == sum { res = 1 }
        return res + dfs(node.Left, sum - node.Val) + dfs(node.Right, sum - node.Val)
    }
    return dfs(root, targetSum) + pathSum(root.Left, targetSum) + pathSum(root.Right, targetSum)`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0149",
			Slug:        "find-all-anagrams-in-a-string",
			Title:       "Find All Anagrams in a String",
			Topic:       "String",
			Statement:   "Given two strings s and p, return an array of all the start indices of p's anagrams in s. You may return the answer in any order.",
			ExampleIn:   "s = \"cbaebabacd\", p = \"abc\"",
			ExampleOut:  "[0,6]",
			HintTitle:   "Sliding window frequencies",
			HintBody:    "Use sliding window of size p.length and frequency maps comparison.",
			FuncName:    "findAnagrams",
			ParamsJS:    "s, p",
			ParamsPy:    "s: str, p: str",
			ParamsJava:  "String s, String p",
			ParamsCpp:   "string s, string p",
			ParamsGo:    "s string, p string",
			RetPy:       "list[int]",
			RetJava:     "List<Integer>",
			RetCpp:      "vector<int>",
			RetGo:       "[]int",
			SolJS: `    let res = [], pMap = {}, sMap = {};
    for (let char of p) pMap[char] = (pMap[char] || 0) + 1;
    let l = 0;
    for (let r = 0; r < s.length; r++) {
        sMap[s[r]] = (sMap[s[r]] || 0) + 1;
        if (r - l + 1 > p.length) {
            sMap[s[l]]--;
            if (sMap[s[l]] === 0) delete sMap[s[l]];
            l++;
        }
        let matches = true;
        for (let key in pMap) { if (sMap[key] !== pMap[key]) { matches = false; break; } }
        if (matches && (r - l + 1 === p.length)) res.push(l);
    }
    return res;`,
			SolPy: `    from collections import Counter
    res, p_count, s_count = [], Counter(p), Counter()
    for r, c in enumerate(s):
        s_count[c] += 1
        if r >= len(p):
            left_c = s[r - len(p)]
            s_count[left_c] -= 1
            if s_count[left_c] == 0: del s_count[left_c]
        if s_count == p_count: res.append(r - len(p) + 1)
    return res`,
			SolJava: `        List<Integer> res = new ArrayList<>();
        if (s.length() < p.length()) return res;
        int[] pCount = new int[26];
        int[] sCount = new int[26];
        for (char c : p.toCharArray()) pCount[c - 'a']++;
        for (int i = 0; i < s.length(); i++) {
            sCount[s.charAt(i) - 'a']++;
            if (i >= p.length()) sCount[s.charAt(i - p.length()) - 'a']--;
            if (Arrays.equals(pCount, sCount)) res.add(i - p.length() + 1);
        }
        return res;`,
			SolCpp: `        vector<int> res;
        if (s.length() < p.length()) return res;
        vector<int> pCount(26, 0), sCount(26, 0);
        for (char c : p) pCount[c - 'a']++;
        for (int i = 0; i < s.length(); ++i) {
            sCount[s[i] - 'a']++;
            if (i >= p.length()) sCount[s[i - p.length()] - 'a']--;
            if (pCount == sCount) res.push_back(i - p.length() + 1);
        }
        return res;`,
			SolGo: `    var res []int
    if len(s) < len(p) { return nil }
    pCount := make([]int, 26)
    sCount := make([]int, 26)
    for i := 0; i < len(p); i++ { pCount[p[i]-'a']++ }
    for i := 0; i < len(s); i++ {
        sCount[s[i]-'a']++
        if i >= len(p) { sCount[s[i-len(p)]-'a']-- }
        match := true
        for j := 0; j < 26; j++ {
            if pCount[j] != sCount[j] { match = false; break }
        }
        if match { res = append(res, i-len(p)+1) }
    }
    return res`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0150",
			Slug:        "target-sum",
			Title:       "Target Sum",
			Topic:       "DP",
			Statement:   "You are given an integer array nums and an integer target. Build an expression using + and - before each integer, returning the number of different expressions that evaluate to target.",
			ExampleIn:   "nums = [1,1,1,1,1], target = 3",
			ExampleOut:  "5",
			HintTitle:   "Subset Sum DP",
			HintBody:    "Equivalently finds subset sum equals (target + totalSum) / 2.",
			FuncName:    "findTargetSumWays",
			ParamsJS:    "nums, target",
			ParamsPy:    "nums: list[int], target: int",
			ParamsJava:  "int[] nums, int target",
			ParamsCpp:   "vector<int>& nums, int target",
			ParamsGo:    "nums []int, target int",
			RetPy:       "int",
			RetJava:     "int",
			RetCpp:      "int",
			RetGo:       "int",
			SolJS: `    let count = 0;
    const dfs = (i, sum) => {
        if (i === nums.length) { if (sum === target) count++; return; }
        dfs(i + 1, sum + nums[i]);
        dfs(i + 1, sum - nums[i]);
    };
    dfs(0, 0);
    return count;`,
			SolPy: `    memo = {}
    def dfs(i, current_sum):
        if (i, current_sum) in memo: return memo[(i, current_sum)]
        if i == len(nums): return 1 if current_sum == target else 0
        ans = dfs(i + 1, current_sum + nums[i]) + dfs(i + 1, current_sum - nums[i])
        memo[(i, current_sum)] = ans
        return ans
    return dfs(0, 0)`,
			SolJava: `        return dfs(nums, target, 0, 0);
    }
    private int dfs(int[] nums, int target, int i, int sum) {
        if (i == nums.length) return sum == target ? 1 : 0;
        return dfs(nums, target, i + 1, sum + nums[i]) + dfs(nums, target, i + 1, sum - nums[i]);`,
			SolCpp: `        return dfs(nums, target, 0, 0);
    }
    int dfs(vector<int>& nums, int target, int i, int sum) {
        if (i == nums.size()) return sum == target ? 1 : 0;
        return dfs(nums, target, i + 1, sum + nums[i]) + dfs(nums, target, i + 1, sum - nums[i]);`,
			SolGo: `    var dfs func(int, int) int
    dfs = func(i, sum int) int {
        if i == len(nums) {
            if sum == target { return 1 }; return 0
        }
        return dfs(i+1, sum+nums[i]) + dfs(i+1, sum-nums[i])
    }
    return dfs(0, 0)`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0151",
			Slug:        "daily-temperatures",
			Title:       "Daily Temperatures",
			Topic:       "Stack/Queue",
			Statement:   "Given an array of integers temperatures represents the daily temperatures, return an array answer such that answer[i] is the number of days you have to wait after the ith day to get a warmer temperature.",
			ExampleIn:   "[73,74,75,71,69,72,76,73]",
			ExampleOut:  "[1,1,4,2,1,1,0,0]",
			HintTitle:   "Monotonic Stack",
			HintBody:    "Maintain a monotonic decreasing stack of indices.",
			FuncName:    "dailyTemperatures",
			ParamsJS:    "temperatures",
			ParamsPy:    "temperatures: list[int]",
			ParamsJava:  "int[] temperatures",
			ParamsCpp:   "vector<int>& temperatures",
			ParamsGo:    "temperatures []int",
			RetPy:       "list[int]",
			RetJava:     "int[]",
			RetCpp:      "vector<int>",
			RetGo:       "[]int",
			SolJS: `    let res = Array(temperatures.length).fill(0), stack = [];
    for (let i = 0; i < temperatures.length; i++) {
        while (stack.length && temperatures[i] > temperatures[stack[stack.length - 1]]) {
            let idx = stack.pop(); res[idx] = i - idx;
        }
        stack.push(i);
    }
    return res;`,
			SolPy: `    res, stack = [0] * len(temperatures), []
    for i, t in enumerate(temperatures):
        while stack and t > temperatures[stack[-1]]:
            idx = stack.pop()
            res[idx] = i - idx
        stack.append(i)
    return res`,
			SolJava: `        int[] res = new int[temperatures.length];
        Stack<Integer> stack = new Stack<>();
        for (int i = 0; i < temperatures.length; i++) {
            while (!stack.isEmpty() && temperatures[i] > temperatures[stack.peek()]) {
                int idx = stack.pop(); res[idx] = i - idx;
            }
            stack.push(i);
        }
        return res;`,
			SolCpp: `        vector<int> res(temperatures.size(), 0);
        stack<int> s;
        for (int i = 0; i < temperatures.size(); ++i) {
            while (!s.empty() && temperatures[i] > temperatures[s.top()]) {
                int idx = s.top(); s.pop();
                res[idx] = i - idx;
            }
            s.push(i);
        }
        return res;`,
			SolGo: `    res := make([]int, len(temperatures))
    var stack []int
    for i, t := range temperatures {
        for len(stack) > 0 && t > temperatures[stack[len(stack)-1]] {
            idx := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            res[idx] = i - idx
        }
        stack = append(stack, i)
    }
    return res`,
		},
		{
			ID:          "54574a34-9a68-4e65-ab9a-af05db4d0152",
			Slug:        "koko-eating-bananas",
			Title:       "Koko Eating Bananas",
			Topic:       "Array",
			Statement:   "Koko loves to eat bananas. There are n piles of bananas. Determine the minimum integer k such that she can eat all the bananas within h hours.",
			ExampleIn:   "piles = [3,6,7,11], h = 8",
			ExampleOut:  "4",
			HintTitle:   "Binary Search on Speed",
			HintBody:    "Binary search the speed k in range [1, max(piles)]. Check validity of speed by computing total hours.",
			FuncName:    "minEatingSpeed",
			ParamsJS:    "piles, h",
			ParamsPy:    "piles: list[int], h: int",
			ParamsJava:  "int[] piles, int h",
			ParamsCpp:   "vector<int>& piles, int h",
			ParamsGo:    "piles []int, h int",
			RetPy:       "int",
			RetJava:     "int",
			RetCpp:      "int",
			RetGo:       "int",
			SolJS: `    let l = 1, r = Math.max(...piles);
    while (l < r) {
        let mid = Math.floor((l + r) / 2);
        let hours = piles.reduce((a, b) => a + Math.ceil(b / mid), 0);
        if (hours <= h) r = mid;
        else l = mid + 1;
    }
    return l;`,
			SolPy: `    l, r = 1, max(piles)
    while l < r:
        mid = (l + r) // 2
        hours = sum((p + mid - 1) // mid for p in piles)
        if hours <= h: r = mid
        else: l = mid + 1
    return l`,
			SolJava: `        int l = 1, r = 1000000000;
        while (l < r) {
            int mid = l + (r - l) / 2;
            int hours = 0;
            for (int p : piles) hours += (p + mid - 1) / mid;
            if (hours <= h) r = mid;
            else l = mid + 1;
        }
        return l;`,
			SolCpp: `        int l = 1, r = 1000000000;
        while (l < r) {
            int mid = l + (r - l) / 2;
            long long hours = 0;
            for (int p : piles) hours += (p + mid - 1) / mid;
            if (hours <= h) r = mid;
            else l = mid + 1;
        }
        return l;`,
			SolGo: `    l, r := 1, 1000000000
    for l < r {
        mid := l + (r-l)/2
        hours := 0
        for _, p := range piles { hours += (p + mid - 1) / mid }
        if hours <= h { r = mid } else { l = mid + 1 }
    }
    return l`,
		},
	}

	var probs []Problem
	for _, m := range metas {
		var jsSC, pySC, javaSC, cppSC, goSC string

		if m.FuncName == "lruCache" || m.FuncName == "minStack" || m.FuncName == "trie" || m.FuncName == "ticTacToe" {
			jsSC = m.SolJS
			pySC = m.SolPy
			javaSC = m.SolJava
			cppSC = m.SolCpp
			goSC = m.SolGo
		} else {
			jsSC = "function " + m.FuncName + "(" + m.ParamsJS + ") {\n    // Write your code here\n" + m.SolJS + "\n}"
			pySC = "def " + m.FuncName + "(" + m.ParamsPy + ") -> " + m.RetPy + ":\n    # Write your code here\n" + m.SolPy
			javaSC = "public class Solution {\n    public " + m.RetJava + " " + m.FuncName + "(" + m.ParamsJava + ") {\n        // Write your code here\n" + m.SolJava + "\n    }\n}"
			cppSC = "class Solution {\npublic:\n    " + m.RetCpp + " " + m.FuncName + "(" + m.ParamsCpp + ") {\n        // Write your code here\n" + m.SolCpp + "\n    }\n};"
			retGoSpace := ""
			if m.RetGo != "" { retGoSpace = " " + m.RetGo }
			goSC = "package main\n\nfunc " + m.FuncName + "(" + m.ParamsGo + ")" + retGoSpace + " {\n    // Write your code here\n" + m.SolGo + "\n}"
		}

		probs = append(probs, Problem{
			ID:         m.ID,
			Slug:       m.Slug,
			Title:      m.Title,
			Difficulty: "Medium",
			Topic:      m.Topic,
			XP:         100,
			Statement:  m.Statement + "\n\nMake sure your function has the correct signature.",
			SetID:      "54574a34-9a68-4e65-ab9a-af05db4ca002",
			Tags:       []string{m.Topic, "Medium"},
			Examples: []Example{
				{
					Input:       m.ExampleIn,
					Output:      m.ExampleOut,
					Explanation: "Refer to description.",
				},
			},
			Hints: []Hint{
				{
					Title: m.HintTitle,
					Body:  m.HintBody,
				},
			},
			JavascriptSC: jsSC,
			PythonSC:     pySC,
			JavaSC:       javaSC,
			CppSC:        cppSC,
			GoSC:         goSC,
			TestCases: []TestCase{
				{Input: m.ExampleIn, Expected: m.ExampleOut, IsHidden: false},
			},
		})
	}

	return probs
}
