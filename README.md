# Crypt
[Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing) in Go

## Quick Start
```console
$ go build crypt.go
$ ./crypt test
```

## Example Usage
```console
$ ./crypt encode
Enter the secret: Hello, world!
Enter the number of shares (default 6):
Enter the minimum number of shares (default 3):

Shares:
  1, 13145232784294297638295594869514
  2, 27056236747980358260317628304707
  3, 47468828654132036819454247543500
  4, 74383008502749333315705452585893
  5, 107798776293832247749071243431886
  6, 147716132027380780119551620081479

$ ./crypt decode
Enter the shares (^D to stop):
6, 147716132027380780119551620081479
4, 74383008502749333315705452585893
3, 47468828654132036819454247543500

Recovered secret: Hello, world!
```
