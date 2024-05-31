package processors

import(
	"math/rand"
)

var isPrimeSlice []string

func WasteTime(a int, b int) {

	isPrimeSlice = nil

	n := (a + rand.Intn(b - a)) * 100

	for i := 2; i < n; i++ {
        if isPrime(i) {
			isPrimeSlice = append(isPrimeSlice, "Is Prime.")
		} else {
			isPrimeSlice = append(isPrimeSlice, "Is Not Prime.")
		}
    }
}

func isPrime(n int) bool {
    if n <= 1 {
        return false
    }
    for i := 2; i < n; i++ {
        if n % i == 0 {
            return false
        }
    }
    return true
}

