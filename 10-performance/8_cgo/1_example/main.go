package main

/*
#include <math.h>
*/
import "C"
import "fmt"

func main() {
	num := C.double(16)
	result := C.sqrt(num)
	fmt.Printf("Квадратный корень из %.2f равен %.2f\n", num, result)
}
