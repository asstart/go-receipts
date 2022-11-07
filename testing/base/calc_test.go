package base_test

import (
	"fmt"
	"testing"

	"github.com/asstart/go-receipts/testing/base"
)

/*

Example of the simplest test in go

Can run all tests with: go test ./...
Or particular test with: go test -v -run Addition

*/
func TestAddition(t *testing.T) {
	a, b := 2, 3

	res := base.Add(a, b)

	if res != 5 {
		t.Fatalf("result of %v+%v should be equals to 5, but it's %v", a, b, res)
	}
}

/*

Example of using test tables for testing multiple cases in one test
tt - stands for test table
tc - stands for test case

*/
func TestSubstraction(t *testing.T) {
	tt := []struct {
		numbers  []int
		expected int
	}{
		{[]int{5, 4}, 1},
		{[]int{}, 0},
	}

	for _, tc := range tt {
		res := base.Sub(tc.numbers...)
		if res != tc.expected {
			t.Errorf("result of subtracting %v is %v, but should be %v", tc.numbers, res, tc.expected)
		}
	}
}

/*

Example of using test tables for testing multiple cases in one test
But adding additional information as name to identify which case was failed
tt - stands for test table
tc - stands for test case

*/
func TestMuliplication(t *testing.T) {
	tt := []struct {
		name     string
		numbers  []int
		expected int
	}{
		{"1 to 3 mult", []int{1, 2, 3}, 6},
		{"no numbers mult", []int{}, 0},
		{"1 mult", []int{1}, 1},
	}

	for _, tc := range tt {
		res := base.Mult(tc.numbers...)
		if res != tc.expected {
			t.Errorf("result of %v is %v, but should be %v", tc.name, res, tc.expected)
		}
	}
}

/*

Example of using test tables for testing multiple cases in one test
But adding additional information as name to identify which case was failed
And run each tc as subtest, so now we're able to run only certain set of test instead of all
with command: go test -v -run Division/8 (will start first and third test)
tt - stands for test table
tc - stands for test case

*/
func TestDivision(t *testing.T) {
	tt := []struct {
		name     string
		numbers  []int
		expected int
	}{
		{"8 to 2 div", []int{8, 4, 2}, 1},
		{"no numbers div", []int{}, 0},
		{"8 divison", []int{8}, 8},
	}

	for _, tc := range tt {
		t.Run(
			tc.name,
			func(t *testing.T) {
				res := base.Div(tc.numbers...)
				if res != tc.expected {
					t.Fatalf("result of %v is %v, but should be %v", tc.name, res, tc.expected)
				}
			},
		)
	}
}

/*

Test will appear in documentation
For testing run: godoc -http=:6060

*/
func ExampleAdd() {
	s := base.Add(1,2,3)
	fmt.Printf("sum of 1,2,3 is %v", s)
	// Output:
	// sum of 1,2,3 is 6
}