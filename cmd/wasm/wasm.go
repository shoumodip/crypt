package main

import (
	"encoding/hex"
	"github.com/shoumodip/crypt/pkg/crypt"
	"strings"
	"syscall/js"
)

func textBoxShow(textBox js.Value, message string) {
	textBox.Get("style").Set("color", "")
	textBox.Set("value", message)
}

func textBoxError(textBox js.Value, message string) {
	textBox.Get("style").Set("color", "red")
	textBox.Set("value", message)
}

func getElementById(document js.Value, id string) js.Value {
	element := document.Call("getElementById", id)
	if !element.Truthy() {
		panic("Unable to get element by id '" + id + "'")
	}
	return element
}

func addEventListener(element js.Value, event string, callback func(js.Value, []js.Value) any) {
	element.Call("addEventListener", event, js.FuncOf(callback))
}

func encodeUpdateCallback(encodeInput js.Value, encodeOutput js.Value, encodeShares js.Value, encodeMinimum js.Value) {
	if !encodeShares.Get("validity").Get("valid").Bool() {
		textBoxError(encodeOutput, "Number of shares is not valid")
		return
	}

	if !encodeMinimum.Get("validity").Get("valid").Bool() {
		textBoxError(encodeOutput, "Number of minimum shares is not valid")
		return
	}

	n := byte(encodeShares.Get("valueAsNumber").Int())
	k := byte(encodeMinimum.Get("valueAsNumber").Int())
	secret := encodeInput.Get("value").String()
	if secret == "" {
		textBoxShow(encodeOutput, "")
		return
	}

	shares, err := crypt.Encode([]byte(secret), n, k)
	if err != nil {
		textBoxError(encodeOutput, err.Error())
		return
	}

	var result strings.Builder
	for _, share := range shares {
		result.WriteString(hex.EncodeToString(share))
		result.WriteByte('\n')
	}

	textBoxShow(encodeOutput, result.String())
}

func main() {
	document := js.Global().Get("document")
	if !document.Truthy() {
		panic("Unable to get document object")
	}

	encodeTab := getElementById(document, "encode-tab")
	encodeMain := getElementById(document, "encode-main")
	encodeInput := getElementById(document, "encode-input")
	encodeOutput := getElementById(document, "encode-output")
	encodeShares := getElementById(document, "encode-shares")
	encodeMinimum := getElementById(document, "encode-minimum")

	decodeTab := getElementById(document, "decode-tab")
	decodeMain := getElementById(document, "decode-main")
	decodeInput := getElementById(document, "decode-input")
	decodeOutput := getElementById(document, "decode-output")

	addEventListener(encodeTab, "click", func(this js.Value, args []js.Value) any {
		decodeTab.Set("className", "")
		decodeMain.Set("className", "")

		encodeTab.Set("className", "current")
		encodeMain.Set("className", "current")
		return nil
	})

	addEventListener(decodeTab, "click", func(this js.Value, args []js.Value) any {
		encodeTab.Set("className", "")
		encodeMain.Set("className", "")

		decodeTab.Set("className", "current")
		decodeMain.Set("className", "current")
		return nil
	})

	addEventListener(encodeInput, "input", func(this js.Value, args []js.Value) any {
		encodeUpdateCallback(encodeInput, encodeOutput, encodeShares, encodeMinimum)
		return nil
	})

	addEventListener(encodeShares, "input", func(this js.Value, args []js.Value) any {
		encodeUpdateCallback(encodeInput, encodeOutput, encodeShares, encodeMinimum)
		return nil
	})

	addEventListener(encodeMinimum, "input", func(this js.Value, args []js.Value) any {
		encodeUpdateCallback(encodeInput, encodeOutput, encodeShares, encodeMinimum)
		return nil
	})

	addEventListener(decodeInput, "input", func(this js.Value, args []js.Value) any {
		source := strings.Split(decodeInput.Get("value").String(), "\n")

		var shares [][]byte
		for _, line := range source {
			if line = strings.TrimSpace(line); len(line) > 0 {
				share, err := hex.DecodeString(line)
				if err != nil {
					textBoxError(decodeOutput, err.Error())
					return nil
				}

				shares = append(shares, share)
			}
		}

		secret, err := crypt.Decode(shares)
		if err != nil {
			textBoxError(decodeOutput, err.Error())
			return nil
		}

		textBoxShow(decodeOutput, string(secret))
		return nil
	})

	<-make(chan bool)
}
