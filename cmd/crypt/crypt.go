package main

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"github.com/shoumodip/crypt/pkg/crypt"
	"os"
	"strconv"
)

func inputLine(scanner *bufio.Scanner) (string, error) {
	if !scanner.Scan() {
		return "", scanner.Err()
	}

	return scanner.Text(), nil
}

func inputByte(scanner *bufio.Scanner, fallback byte) (byte, error) {
	input, err := inputLine(scanner)
	if err != nil {
		return 0, err
	}

	if input == "" {
		return fallback, nil
	}

	result, err := strconv.ParseUint(input, 10, 8)
	return byte(result), err
}

func handleError(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func modeHelp(file *os.File) {
	fmt.Fprintln(file, "usage:")
	fmt.Fprintln(file, "  "+os.Args[0]+" <mode>")
	fmt.Fprintln(file)
	fmt.Fprintln(file, "modes:")
	fmt.Fprintln(file, "  help    Show this message and exit")
	fmt.Fprintln(file, "  decode  Decode the secret from shares")
	fmt.Fprintln(file, "  encode  Encode the secret into shares")
}

func modeDecode() {
	shares := [][]byte{}
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("Enter the shares (^D to stop):")
	for {
		input, err := inputLine(scanner)
		if input == "" {
			break
		}
		handleError(err)

		share, err := hex.DecodeString(input)
		handleError(err)

		shares = append(shares, share)
	}

	secret, err := crypt.Decode(shares)
	handleError(err)

	fmt.Println()
	fmt.Println("Recovered secret:", string(secret))
}

func modeEncode() {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Print("Enter the secret: ")
	secret, err := inputLine(scanner)
	handleError(err)

	fmt.Print("Enter the number of shares (default 4): ")
	n, err := inputByte(scanner, 4)
	handleError(err)

	fmt.Print("Enter the minimum number of shares (default 2): ")
	k, err := inputByte(scanner, 2)
	handleError(err)

	shares, err := crypt.Encode([]byte(secret), n, k)
	handleError(err)

	fmt.Println()
	fmt.Println("Shares:")
	for _, share := range shares {
		fmt.Println(hex.EncodeToString(share))
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintln(os.Stderr, "error: mode not provided")
		fmt.Fprintln(os.Stderr)
		modeHelp(os.Stderr)
		os.Exit(1)
	}

	switch mode := os.Args[1]; mode {
	case "help":
		modeHelp(os.Stdout)

	case "decode":
		modeDecode()

	case "encode":
		modeEncode()

	default:
		fmt.Fprintln(os.Stderr, "error: invalid mode '"+mode+"'")
		fmt.Fprintln(os.Stderr)
		modeHelp(os.Stderr)
		os.Exit(1)
	}
}
