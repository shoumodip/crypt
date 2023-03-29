# Crypt
[Shamir's Secret Sharing](https://en.wikipedia.org/wiki/Shamir%27s_secret_sharing) in Go

## Quick Start
```console
$ go build cmd/crypt/crypt.go
```

## Example Usage
```console
$ ./crypt encode
Enter the secret: Hello, world!
Enter the number of shares (default 6):
Enter the minimum number of shares (default 3):

Shares:
005a98af88a98a9b2e3690779b10
01eb2c14982fa9a3715200746ec0
02f9d1d77ce90f18280be26f91f1
033f6125039e58a89d1f14e4e1c2
042d9ce6e758fe13c446f6ff1ef3
059c285df7dedd2b9b2266fceb23

$ ./crypt decode
Enter the shares (^D to stop):
01eb2c14982fa9a3715200746ec0
033f6125039e58a89d1f14e4e1c2
059c285df7dedd2b9b2266fceb23

Recovered secret: Hello, world!
```

## Web Version
```console
$ GOOS=js GOARCH=wasm go build -o web/crypt.wasm cmd/wasm/wasm.go
$ go run cmd/server/server.go
Serving HTTP on localhost:8000
```

Open your web browser on `localhost:8000`
