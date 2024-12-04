package main

// Зачем нужны?
func equalInt(a, b int) bool {
	return a == b
}

func equalFloat(a, b float64) bool {
	return a == b
}

func equal[T int | float64](a, b T) bool {
	return a == b
}

// Алиасы
type MyType int

func equalWithAlias[T ~int | ~float64](a, b T) bool {
	return a == b
}

func KeyExists[T any](m map[string]T, key string) bool {
	_, exists := m[key]
	return exists
}

func main() {
	println(equalInt(1, 2))
	println(equalFloat(1.0, 2.0))

	println(equal(1, 2))
	println(equal(1.0, 2.0))

	// Cannot use MyType as the type interface{ int | float64 }
	// Type does not implement constraint interface{ int | float64 } because type is not included in type set (int, float64)
	//
	//
	//println(equal(MyType(1),MyType(2)))
	println(equalWithAlias(MyType(1), MyType(2)))

	m1 := make(map[string]int)
	m2 := make(map[string]bool)

	println(KeyExists(m2, "test"))
	println(KeyExists(m1, "test"))

}
