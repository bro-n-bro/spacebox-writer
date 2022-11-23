package utils

import (
	"math/rand"
	"time"
)

func RandomUsers(list []string) (u1, u2 string) {
	rand.Seed(time.Now().UnixNano())
	u1 = list[rand.Intn(len(list))]

	rand.Seed(time.Now().UnixNano())
	u2 = list[rand.Intn(len(list))]

	return
}

func RandomFloat64(min, max float64) float64 {
	rand.Seed(time.Now().UnixNano())
	return min + rand.Float64()*(max-min)
}
