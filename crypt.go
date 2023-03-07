package main

import (
	"bufio"
	"crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"os"
	"strconv"
	"strings"
)

// Encoder
// BUG: sometimes decoding does not result in the actual secret
// Specifically when decoding with the coordinates (1, _) (2, _) (_, _)

type Point struct {
	x, y *big.Int
}

func encode(secret *big.Int, shares int, minimum int) ([]Point, *big.Int, error) {
	if minimum > shares {
		return nil, nil, errors.New("minimum shares cannot be larger than total shares")
	}

	poly := make([]*big.Int, minimum)
	for i := 1; i < minimum; i++ {
		n, err := rand.Int(rand.Reader, secret)
		if err != nil {
			return nil, nil, err
		}

		poly[i] = n
	}
	poly[0] = secret

	field := secret
	if field.Int64()&1 == 0 {
		field = new(big.Int).Add(field, big.NewInt(1))
	}

	for !field.ProbablyPrime(20) {
		field = new(big.Int).Add(field, big.NewInt(2))
	}

	points := make([]Point, shares)
	for i := 0; i < shares; i++ {
		points[i] = Point{big.NewInt(int64(i + 1)), big.NewInt(0)}

		a := big.NewInt(1)
		for _, c := range poly {
			points[i].y = new(big.Int).Add(points[i].y, new(big.Int).Mul(a, c))
			a = new(big.Int).Mul(a, points[i].x)
		}

		points[i].y = new(big.Int).Mod(points[i].y, field)
	}

	return points, field, nil
}

func decode(points []Point, field *big.Int) *big.Int {
	sum := big.NewFloat(0.0)
	for i := range points {
		xi := new(big.Float).SetInt(points[i].x)

		product := big.NewFloat(1.0)
		for j := range points {
			if j != i {
				xj := new(big.Float).SetInt(points[j].x)
				product = new(big.Float).Mul(product, new(big.Float).Quo(xj, new(big.Float).Sub(xj, xi)))
			}
		}

		sum = new(big.Float).Add(sum, new(big.Float).Mul(new(big.Float).SetInt(points[i].y), product))
	}

	result, _ := sum.Int(nil)
	return new(big.Int).Mod(result, field)
}

// Parser
func inputLine(scanner *bufio.Scanner) (string, error) {
	if !scanner.Scan() {
		return "", scanner.Err()
	}

	return scanner.Text(), nil
}

func parseBigInt(str string) (*big.Int, error) {
	n, ok := new(big.Int).SetString(strings.TrimSpace(str), 10)
	if !ok {
		return nil, errors.New("invalid number '" + str + "'")
	}
	return n, nil
}

func inputBigInt(scanner *bufio.Scanner) (*big.Int, error) {
	input, err := inputLine(scanner)
	if err != nil {
		return nil, err
	}
	return parseBigInt(input)
}

func inputIntMaybe(scanner *bufio.Scanner, fallback int) (int, error) {
	input, err := inputLine(scanner)
	if err != nil {
		return 0, err
	}

	if input == "" {
		return fallback, nil
	}

	return strconv.Atoi(input)
}

// CLI
func usage(file *os.File) {
	fmt.Fprintln(file, "usage:")
	fmt.Fprintln(file, "  "+os.Args[0]+" <mode>")
	fmt.Fprintln(os.Stderr)
	fmt.Fprintln(file, "modes:")
	fmt.Fprintln(file, "  help    Show this message and exit")
	fmt.Fprintln(file, "  decode  Decode the secret from shares")
	fmt.Fprintln(file, "  encode  Encode the secret into shares")
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func modeDecode() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the field: ")
	field, err := inputBigInt(scanner)
	handleError(err)

	fmt.Println("Enter the shares (^D to stop):")
	var shares []Point
	for {
		input, err := inputLine(scanner)
		if input == "" {
			break
		}
		handleError(err)

		coords := strings.Split(input, ",")
		if len(coords) != 2 {
			fmt.Fprintln(os.Stderr, "error: invalid share '"+input+"'")
			os.Exit(1)
		}

		var share Point
		share.x, err = parseBigInt(coords[0])
		handleError(err)

		share.y, err = parseBigInt(coords[1])
		handleError(err)

		shares = append(shares, share)
	}

	fmt.Println()
	fmt.Println("Recovered secret:", decode(shares, field))
}

func modeEncode() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the secret: ")
	secret, err := inputBigInt(scanner)
	handleError(err)

	fmt.Print("Enter the number of shares (default 6): ")
	shares, err := inputIntMaybe(scanner, 6)
	handleError(err)

	fmt.Print("Enter the minimum number of shares (default 3): ")
	minimum, err := inputIntMaybe(scanner, 3)
	handleError(err)

	points, field, err := encode(secret, shares, minimum)
	handleError(err)

	fmt.Println()
	fmt.Println("Field:", field)
	fmt.Println("Shares:")
	for _, point := range points {
		fmt.Printf("  %d, %d\n", point.x, point.y)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "error: mode not provided")
		fmt.Fprintln(os.Stderr)
		usage(os.Stderr)
		os.Exit(1)
	}

	switch mode := os.Args[1]; mode {
	case "help":
		usage(os.Stdout)

	case "decode":
		modeDecode()

	case "encode":
		modeEncode()

	default:
		fmt.Fprintln(os.Stderr, "error: invalid mode '"+mode+"'")
		fmt.Fprintln(os.Stderr)
		usage(os.Stderr)
		os.Exit(1)
	}
}
