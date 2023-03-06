package main

import (
	"crypto/rand"
	"errors"
	"fmt"
	"math"
	"math/big"
)

// TODO: enable arbitrary precision
type Point struct {
	x, y int
}

func solve(poly []int, x int) Point {
	a := 1
	y := 0
	for _, c := range poly {
		y += a * c
		a *= x
	}
	return Point{x, y}
}

func random() (int, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt))
	if err != nil {
		return 0, err
	}

	// TODO: enable decoding beyond numeric limits
	return int(n.Int64()) >> 32, nil
}

// TODO: use finite field to eliminate brute force attacks
func encode(secret int, shares int, minimum int) ([]Point, error) {
	if minimum > shares {
		return nil, errors.New("share subset cannot be larger than total shares")
	}

	poly := make([]int, minimum)
	for i := 1; i < minimum; i++ {
		n, err := random()
		if err != nil {
			return nil, err
		}

		poly[i] = n
	}
	poly[0] = secret

	points := make([]Point, shares)
	for i := 0; i < shares; i++ {
		points[i] = solve(poly, i+1)
	}

	return points, nil
}

func decode(points []Point) int {
	sum := 0.0
	for j := range points {
		product := 1.0
		for m := range points {
			if m != j {
				// TODO: enable arbitrary precision
				product *= float64(points[m].x) / float64(points[m].x-points[j].x)
			}
		}
		sum += float64(points[j].y) * product
	}

	return int(sum)
}

// TODO: customize at runtime
const (
	SECRET  = 123456789
	SHARES  = 6
	MINIMUM = 3
)

func main() {
	shares, err := encode(SECRET, SHARES, MINIMUM)
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
