//go:build js
// +build js

package jsUtil

import (
	"fmt"
	"regexp"
	"strconv"
	"syscall/js"
)

var MOBILE_BROWSER_REGEX = regexp.MustCompile("(?i)Android|webOS|iPhone|iPad|iPod|BlackBerry|Windows Phone")

var document js.Value

var cb InsertCallBack

var started bool

var offsetTop int
var offsetLeft int

func init() {
	document = js.Global().Get("document")
	p := document.Call("createElement", "input")
	p.Set("id", "tempInput")
	p.Set("style", "height:0px; width:1px; margin:0px; position: fixed; overflow:hidden; top:-10px; border:0px; padding:0px")
	document.Get("body").Call("appendChild", p)

	p.Call("addEventListener", "input", js.FuncOf(handleInput), false)
	canvas := document.Get("body").Call("getElementsByTagName", "canvas").Index(0)
	offsetTop = canvas.Get("offsetTop").Int()
	offsetLeft = canvas.Get("offsetLeft").Int()
	canvas.Call("addEventListener", "touchstart", js.FuncOf(handleClick), false)
	canvas.Call("addEventListener", "touchend", js.FuncOf(handleClick), false)
}

func IsMobileBrowser() bool {
	navigator := js.Global().Get("navigator")
	userAgent := navigator.Get("userAgent")
	return MOBILE_BROWSER_REGEX.Match([]byte(userAgent.String()))
}

func Prompt(title string, value string, cursorPos int, yPos int, callback InsertCallBack) (string, bool) {
	fmt.Println("Prompt")
	/*	prompt := js.Global().Get("prompt")
		result := prompt.Invoke(title, value)
		if !result.IsNull() && !result.IsUndefined() {
			return result.String(), true
		} else {
			return "", false
		}
	*/
	cb = callback
	p := document.Call("getElementById", "tempInput")
	p.Call("setAttribute", "inputmode", "text")
	p.Set("value", value)
	p.Call("setSelectionRange", cursorPos, cursorPos)

	started = true
	fmt.Println(offsetTop + yPos)
	p.Get("style").Call("setProperty", "top", strconv.Itoa(offsetTop+yPos)+"px")

	return "", false
}

func SetCursorPosition(cursorPos int) {
	p := document.Call("getElementById", "tempInput")
	p.Call("setSelectionRange", cursorPos, cursorPos)
}

func GetCursorPosition() int {
	p := document.Call("getElementById", "tempInput")
	return p.Get("selectionStart").Int()
}

func handleClick(this js.Value, args []js.Value) any {
	if started {
		p := document.Call("getElementById", "tempInput")
		p.Call("focus")
		started = false
	}
	return nil
}

func handleInput(this js.Value, args []js.Value) any {
	lastTypedChar := args[0].Get("target").Get("value").String()
	fmt.Println(lastTypedChar)
	if cb != nil {
		cb(lastTypedChar)
	}
	return nil
}
