package system

import (
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
)

// ButtonState is the current state of an input button
type ButtonState int

// KeyID is a key ID!
type KeyID int

// Keyboard constants
const (
	KeyA KeyID = iota
	KeyB
	KeyC
	KeyD
	KeyE
	KeyF
	KeyG
	KeyH
	KeyI
	KeyJ
	KeyK
	KeyL
	KeyM
	KeyN
	KeyO
	KeyP
	KeyQ
	KeyR
	KeyS
	KeyT
	KeyU
	KeyV
	KeyW
	KeyX
	KeyY
	KeyZ
	KeyEsc // 26
	KeyCtrl

	KeyCt = 28 // Total number of keys to be checked
)

// ButtonState constants
const (
	IDPressed ButtonState = iota
	IDHeld
	IDReleased
	IDNotHeld
)

// Input is a set of current mouse and keyboard inputs
type Input struct {
	MousePos pixel.Vec
	BtnL     ButtonState
	BtnR     ButtonState
	Keys     [KeyCt]ButtonState
}

func inputLogic(prev ButtonState, now bool) ButtonState {
	if now { // if currently pressed
		switch prev {
		case IDPressed:
			return IDHeld
		case IDHeld:
			return IDHeld
		case IDReleased:
			return IDPressed
		case IDNotHeld:
			return IDPressed
		}
	} else { // if not currently pressed
		switch prev {
		case IDPressed:
			return IDReleased
		case IDHeld:
			return IDReleased
		case IDReleased:
			return IDNotHeld
		case IDNotHeld:
			return IDNotHeld
		}
	}
	return -1
}

// ToBool takes ButtonState and just tells you if it means the button is currently pressed
func ToBool(state ButtonState) bool {
	if state == IDPressed || state == IDHeld {
		return true
	}
	return false
}

// GetInput extracts input from a Window and returns it
func GetInput(win *pixelgl.Window, prev Input) Input {
	prev.Keys[KeyA] = inputLogic(prev.Keys[KeyA], win.Pressed(pixelgl.KeyA)) // this is cumbersome
	prev.Keys[KeyB] = inputLogic(prev.Keys[KeyB], win.Pressed(pixelgl.KeyB))
	prev.Keys[KeyC] = inputLogic(prev.Keys[KeyC], win.Pressed(pixelgl.KeyC))
	prev.Keys[KeyD] = inputLogic(prev.Keys[KeyD], win.Pressed(pixelgl.KeyD))
	prev.Keys[KeyE] = inputLogic(prev.Keys[KeyE], win.Pressed(pixelgl.KeyE))
	prev.Keys[KeyF] = inputLogic(prev.Keys[KeyF], win.Pressed(pixelgl.KeyF))
	prev.Keys[KeyG] = inputLogic(prev.Keys[KeyG], win.Pressed(pixelgl.KeyG))
	prev.Keys[KeyH] = inputLogic(prev.Keys[KeyH], win.Pressed(pixelgl.KeyH))
	prev.Keys[KeyI] = inputLogic(prev.Keys[KeyI], win.Pressed(pixelgl.KeyI))
	prev.Keys[KeyJ] = inputLogic(prev.Keys[KeyJ], win.Pressed(pixelgl.KeyJ))
	prev.Keys[KeyK] = inputLogic(prev.Keys[KeyK], win.Pressed(pixelgl.KeyK))
	prev.Keys[KeyL] = inputLogic(prev.Keys[KeyL], win.Pressed(pixelgl.KeyL))
	prev.Keys[KeyM] = inputLogic(prev.Keys[KeyM], win.Pressed(pixelgl.KeyM))
	prev.Keys[KeyN] = inputLogic(prev.Keys[KeyN], win.Pressed(pixelgl.KeyN))
	prev.Keys[KeyO] = inputLogic(prev.Keys[KeyO], win.Pressed(pixelgl.KeyO))
	prev.Keys[KeyP] = inputLogic(prev.Keys[KeyP], win.Pressed(pixelgl.KeyP))
	prev.Keys[KeyQ] = inputLogic(prev.Keys[KeyQ], win.Pressed(pixelgl.KeyQ))
	prev.Keys[KeyR] = inputLogic(prev.Keys[KeyR], win.Pressed(pixelgl.KeyR))
	prev.Keys[KeyS] = inputLogic(prev.Keys[KeyS], win.Pressed(pixelgl.KeyS))
	prev.Keys[KeyT] = inputLogic(prev.Keys[KeyT], win.Pressed(pixelgl.KeyT))
	prev.Keys[KeyU] = inputLogic(prev.Keys[KeyU], win.Pressed(pixelgl.KeyU))
	prev.Keys[KeyV] = inputLogic(prev.Keys[KeyV], win.Pressed(pixelgl.KeyV))
	prev.Keys[KeyW] = inputLogic(prev.Keys[KeyW], win.Pressed(pixelgl.KeyW))
	prev.Keys[KeyX] = inputLogic(prev.Keys[KeyX], win.Pressed(pixelgl.KeyX))
	prev.Keys[KeyY] = inputLogic(prev.Keys[KeyY], win.Pressed(pixelgl.KeyY))
	prev.Keys[KeyZ] = inputLogic(prev.Keys[KeyZ], win.Pressed(pixelgl.KeyZ))

	prev.Keys[KeyEsc] = inputLogic(prev.Keys[KeyEsc], win.Pressed(pixelgl.KeyEscape))
	prev.Keys[KeyCtrl] = inputLogic(prev.Keys[KeyCtrl], win.Pressed(pixelgl.KeyLeftControl))

	return Input{
		MousePos: win.MousePosition(),
		BtnL:     inputLogic(prev.BtnL, win.Pressed(pixelgl.MouseButton1)),
		BtnR:     inputLogic(prev.BtnR, win.Pressed(pixelgl.MouseButton2)),
		Keys:     prev.Keys,
	}
}

// BlankInput returns a 'blank' input (nothing is pressed).
func BlankInput() Input {
	keys := [KeyCt]ButtonState{}
	for i := 0; i < KeyCt; i++ {
		keys[i] = IDNotHeld
	}

	return Input{
		MousePos: pixel.V(0, 0),
		BtnL:     IDNotHeld,
		BtnR:     IDNotHeld,
		Keys:     keys,
	}
}

// JustPressed tells you if a specific combination of keys has just been pressed.
// They need not all have been pressed at the same time, but at least one of them must
// be equal to IDPressed while the others can either be IDPressed or IDHeld to return true.
func (i *Input) JustPressed(keys ...KeyID) bool {
	held := true
	pressed := false

	for _, k := range keys {
		held = held && ToBool(i.Keys[k])            // all keys must currently be held/pressed
		pressed = pressed || i.Keys[k] == IDPressed // at least one must have been JUST pressed
	}
	if !held || !pressed {
		return false
	}

	return true
}
