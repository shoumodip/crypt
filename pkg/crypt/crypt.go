package crypt

import (
	"crypto/rand"
	"errors"
	"github.com/shoumodip/crypt/pkg/field"
)

func Decode(shares [][]byte) ([]byte, error) {
	if len(shares) < 2 {
		return nil, errors.New("crypt: too few shares")
	}

	if len(shares[0]) < 2 {
		return nil, errors.New("crypt: invalid shares")
	}

	for _, share := range shares {
		if len(share) != len(shares[0]) {
			return nil, errors.New("crypt: invalid shares")
		}
	}

	xs := make([]byte, len(shares))
	ys := make([]byte, len(shares))
	secret := make([]byte, len(shares[0])-1)

	for i := range secret {
		for j, v := range shares {
			xs[j] = v[0] + 1
			ys[j] = v[i+1]
		}

		result := byte(0)
		for i := range shares {
			weight := byte(1)
			for j := range shares {
				if i != j {
					weight = field.Mul(weight, field.Div(xs[j], field.Sub(xs[j], xs[i])))
				}
			}
			result = field.Add(result, field.Mul(weight, ys[i]))
		}

		secret[i] = result
	}

	return secret, nil
}

func Encode(secret []byte, n, k byte) ([][]byte, error) {
	if n <= k {
		return nil, errors.New("crypt: minimum shares must be less than total shares")
	}

	shares := make([][]byte, n)
	for i := range shares {
		shares[i] = append(shares[i], byte(i))
	}

	for _, x := range secret {
		poly := make([]byte, k)
		poly[0] = x

		buffer := make([]byte, k-2)
		if _, err := rand.Read(buffer); err != nil {
			return nil, err
		}

		for i := byte(1); i < k-1; i++ {
			poly[i] = buffer[i-1]
		}

		for {
			buffer = make([]byte, 1)
			if _, err := rand.Read(buffer); err != nil {
				return nil, err
			}

			if buffer[0] != 0 {
				poly[k-1] = buffer[0]
				break
			}
		}

		for x := byte(1); x <= n; x++ {
			result := byte(0)
			for i := 1; i <= len(poly); i++ {
				result = field.Add(field.Mul(result, x), poly[len(poly)-i])
			}

			shares[x-1] = append(shares[x-1], result)
		}
	}

	return shares, nil
}
