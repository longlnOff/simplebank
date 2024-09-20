package utils

import "golang.org/x/exp/rand"

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"


func RandomInt(min int, max int) int {
	return min + rand.Intn(max-min)
}

func randomString(length int) string {
	bytes := make([]byte, length)
	for i := 0; i < length; i++ {
		bytes[i] = letters[RandomInt(0, len(letters)-1)]
	}
	return string(bytes)
}

func RandomOwner() string {
	return randomString(6)
}

func RandomMoney() int64 {
	return int64(RandomInt(0, 1000))
}

func RandomCurrency() string {
	currencies := []string{"EUR", "USD", "CAD"}
	n := len(currencies)
	return currencies[RandomInt(0, n-1)]
}
