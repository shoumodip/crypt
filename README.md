# Crypt
A WIP cryptography system in Go using the [Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing) method.

## Quick Start
```console
$ go build crypt.go
$ ./crypt help
```

## Usage
```console
$ ./crypt encode
Enter the secret: 123
Enter the number of shares (default 6):
Enter the minimum number of shares (default 3):

Field: 127
Shares:
  1, 61
  2, 67
  3, 14
  4, 29
  5, 112
  6, 9

$ ./crypt decode
Enter the field: 127
Enter the shares (^D to stop):
6, 9
4, 29
5, 112

Recovered secret: 123
```
