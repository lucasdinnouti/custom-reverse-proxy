package processors

import (
	"math/rand"
	"strconv"
	"strings"
)

func WasteTime(a int, b int) {

	var isPrimeSlice = []string{"Is Prime.", "Is Prime."}
	isPrimeSlice = nil
	primes := 0

	baseMultiplier := 45
	memMultiplier := 1

	n := (a + rand.Intn(b-a)) * baseMultiplier

	for i := 2; i < n; i++ {
		if isPrime(i) {
			isPrimeSlice = append(isPrimeSlice, strings.Repeat(strconv.Itoa(i)+" is Prime.", memMultiplier))
			primes++
		}
	}

	println(isPrimeSlice[rand.Intn(primes-1)])
}

func isPrime(n int) bool {
	if n <= 1 {
		return false
	}
	for i := 2; i < n; i++ {
		if n%i == 0 {
			return false
		}
	}
	return true
}
