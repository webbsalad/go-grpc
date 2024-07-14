package main

import "fmt"

func abs(a int) int {
	if a >= 0 {
		return a
	}

	return -a
}

func findClosestNumber(nums []int) int {
	res := nums[0]

	for _, char := range nums {
		if abs(char) == abs(res) && char > res {
			res = char
		} else if abs(char) < abs(res) {
			res = char
		}
	}

	return res
}

func main() {
	fmt.Println(findClosestNumber([]int{2, -1, 1}))
}
