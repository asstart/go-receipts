package base

func Add(numbers ...int) int {
	res := 0
	for _, n := range numbers {
		res += n
	}
	return res
}

func Sub(numbers ...int) int {
	if len(numbers) == 0 {
		return 0
	}
	res := numbers[0]
	for i := 1; i < len(numbers); i++ {
		res -= numbers[i]
	}
	return res
}

func Mult(numbers ...int) int {
	if len(numbers) == 0 {
		return 0
	}
	res := numbers[0]
	for i := 1; i < len(numbers); i++ {
		res *= numbers[i]
	}
	return res
}

func Div(numbers ...int) int {
	if len(numbers) == 0 {
		return 0
	}
	res := numbers[0]
	for i := 1; i < len(numbers); i++ {
		res /= numbers[i]
	}
	return res
}
