package main

/*
#include <stdio.h>

long long factorialC(long long n) {
    long long result = 1;
    for (long long i = 2; i <= n; i++) {
        result *= i;
    }
    return result;
}
*/
import "C"
import (
	"fmt"
	"math/big"
)

func factorialCGo(n int64) int64 {
	return int64(C.factorialC(C.longlong(n)))
}

func factorialGo(n int64) *big.Int {
	result := big.NewInt(1)
	for i := int64(2); i <= n; i++ {
		result.Mul(result, big.NewInt(i))
	}
	return result
}

func main() {
	n := int64(20)
	fmt.Println("Факториал на C через CGO:", factorialCGo(n))
}
