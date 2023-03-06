package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
)

type Point struct {
	x, y *big.Int
}

func solve(poly []*big.Int, x *big.Int) Point {
	a := big.NewInt(1)
	y := big.NewInt(0)
	for _, c := range poly {
		y = new(big.Int).Add(y, new(big.Int).Mul(a, c))
		a = new(big.Int).Mul(a, x)
	}
	return Point{x, y}
}

// TODO: use finite field to eliminate brute force attacks
func encode(secret *big.Int, shares int, minimum int) ([]Point, error) {
	if minimum > shares {
		return nil, errors.New("share subset cannot be larger than total shares")
	}

	poly := make([]*big.Int, minimum)
	for i := 1; i < minimum; i++ {
		n, err := rand.Int(rand.Reader, secret)
		if err != nil {
			return nil, err
		}

		poly[i] = n
	}
	poly[0] = secret

	points := make([]Point, shares)
	for i := 0; i < shares; i++ {
		points[i] = solve(poly, big.NewInt(int64(i+1)))
	}

	return points, nil
}

func decode(points []Point) int {
	sum := big.NewFloat(0.0)
	for j := range points {
		product := big.NewFloat(1.0)
		for m := range points {
			if m != j {
				xm := new(big.Float).SetInt(points[m].x)
				xj := new(big.Float).SetInt(points[j].x)
				product = new(big.Float).Mul(product, new(big.Float).Quo(xm, new(big.Float).Sub(xm, xj)))
			}

		}
		sum = new(big.Float).Add(sum, new(big.Float).Mul(new(big.Float).SetInt(points[j].y), product))
	}

	result, _ := sum.Int64()
	return int(result)
}

// TODO: customize at runtime
const (
	SECRET  = 123456789
	SHARES  = 6
	MINIMUM = 3
)

func main() {
	shares, err := encode(big.NewInt(SECRET), SHARES, MINIMUM)
	if err != nil {
		panic(err)
	}

	fmt.Println("Shares:")
	for _, share := range shares {
		fmt.Printf("  %d, %d\n", share.x, share.y)
	}
	fmt.Println()

	fmt.Println("Secret:", SECRET)
	fmt.Println("Decode:", decode(shares))
}
