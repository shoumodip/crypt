package main

import (
	"encoding/hex"
	"syscall/js"
	"strings"
	"github.com/shoumodip/crypt/pkg/crypt"
)

func alert(message string) {
	js.Global().Call("alert", message)
}

func getElementById(document js.Value, id string) js.Value {
	element := document.Call("getElementById", id)
	if !element.Truthy() {
		panic("Unable to get element by id '" + id + "'")
	}
	return element
}

func main() {
	document := js.Global().Get("document")
	if !document.Truthy() {
		panic("Unable to get document object")
	}

	input := getElementById(document, "input")
	output := getElementById(document, "output")
	decode := getElementById(document, "decode")
	encode := getElementById(document, "encode")
	shares := getElementById(document, "shares")
	minimum := getElementById(document, "minimum")

	decode.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
		source := strings.Split(input.Get("value").String(), "\n")

		var shares [][]byte
		for _, line := range source {
			if line = strings.TrimSpace(line); len(line) > 0 {
				share, err := hex.DecodeString(line)
				if err != nil {
					alert(err.Error())
					return nil
				}

				shares = append(shares, share)
			}
		}

		if len(shares) < 2 {
			alert("Need 2 or more shares to decode")
			return nil
		}

		output.Set("value", string(crypt.Decode(shares)))
		return nil
	}))

	encode.Call("addEventListener", "click", js.FuncOf(func(this js.Value, args []js.Value) any {
		if !shares.Get("validity").Get("valid").Bool() {
			alert("Number of shares is not valid")
			return nil
		}

		if !minimum.Get("validity").Get("valid").Bool() {
			alert("Number of minimum shares is not valid")
			return nil
		}

		n := byte(shares.Get("valueAsNumber").Int())
		k := byte(minimum.Get("valueAsNumber").Int())
		secret := input.Get("value").String()
		if secret == "" {
			return nil
		}

		shares, err := crypt.Encode([]byte(secret), n, k)
		if err != nil {
			alert(err.Error())
			return nil
		}

		var result strings.Builder
		for _, share := range shares {
			result.WriteString(hex.EncodeToString(share))
			result.WriteByte('\n')
		}

		output.Set("value", result.String())
		return nil
	}))

	<-make(chan bool)
}
