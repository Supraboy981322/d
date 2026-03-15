package main

import rl "github.com/gen2brain/raylib-go/raylib"

type KeyPair struct {
	Key int32
	Val byte
}
var Keys = map[int32]Key {
	rl.KeySpace:        NewKey(' '),
	rl.KeyEscape:       NewKey(0),
	rl.KeyEnter:        NewKey(0),
	rl.KeyTab:          NewKey(0),
	rl.KeyBackspace:    NewKey(0),
	rl.KeyInsert:       NewKey(0),
	rl.KeyDelete:       NewKey(0),
	rl.KeyRight:        NewKey(0),
	rl.KeyLeft:         NewKey(0),
	rl.KeyDown:         NewKey(0),
	rl.KeyUp:           NewKey(0),
	rl.KeyPageUp:       NewKey(0),
	rl.KeyPageDown:     NewKey(0),
	rl.KeyHome:         NewKey(0),
	rl.KeyEnd:          NewKey(0),
	rl.KeyCapsLock:     NewKey(0),
	rl.KeyLeftShift:    NewKey(0),
	rl.KeyLeftControl:  NewKey(0),
	rl.KeyLeftAlt:      NewKey(0),
	rl.KeyLeftSuper:    NewKey(0),
	rl.KeyRightShift:   NewKey(0),
	rl.KeyRightControl: NewKey(0),
	rl.KeyRightAlt:     NewKey(0),
	rl.KeyRightSuper:   NewKey(0),
	rl.KeyLeftBracket:  NewKey('['),
	rl.KeyBackSlash:    NewKey('\\'),
	rl.KeyRightBracket: NewKey(']'),
	rl.KeyGrave:        NewKey('~'),
	rl.KeyApostrophe:   NewKey('\''),
	rl.KeyComma:        NewKey(','),
	rl.KeyMinus:        NewKey('-'),
	rl.KeyPeriod:       NewKey('.'),
	rl.KeySlash:        NewKey('/'),
	rl.KeyZero:         NewKey('0'),
	rl.KeyOne:          NewKey('1'),
	rl.KeyTwo:          NewKey('2'),
	rl.KeyThree:        NewKey('3'),
	rl.KeyFour:         NewKey('4'),
	rl.KeyFive:         NewKey('5'),
	rl.KeySix:          NewKey('6'),
	rl.KeySeven:        NewKey('7'),
	rl.KeyEight:        NewKey('8'),
	rl.KeyNine:         NewKey('9'),
	rl.KeySemicolon:    NewKey(';'),
	rl.KeyEqual:        NewKey('='),
	rl.KeyA:            NewKey('a'),
	rl.KeyB:            NewKey('b'),
	rl.KeyC:            NewKey('c'),
	rl.KeyD:            NewKey('d'),
	rl.KeyE:            NewKey('e'),
	rl.KeyF:            NewKey('f'),
	rl.KeyG:            NewKey('g'),
	rl.KeyH:            NewKey('h'),
	rl.KeyI:            NewKey('i'),
	rl.KeyJ:            NewKey('j'),
	rl.KeyK:            NewKey('k'),
	rl.KeyL:            NewKey('l'),
	rl.KeyM:            NewKey('m'),
	rl.KeyN:            NewKey('n'),
	rl.KeyO:            NewKey('o'),
	rl.KeyP:            NewKey('p'),
	rl.KeyQ:            NewKey('q'),
	rl.KeyR:            NewKey('r'),
	rl.KeyS:            NewKey('s'),
	rl.KeyT:            NewKey('t'),
	rl.KeyU:            NewKey('u'),
	rl.KeyV:            NewKey('v'),
	rl.KeyW:            NewKey('w'),
	rl.KeyX:            NewKey('x'),
	rl.KeyY:            NewKey('y'),
	rl.KeyZ:            NewKey('z'),
}

type Event int
const (
	NOP Event = iota
	ESC
	EXIT
	ERR
)

type Mode int
const (
	NORMAL Mode = iota
	INSERT
	VISUAL
	CMD
)

type (
	Key struct {
		Ticker *Ticker
		Byte byte
	}
	Events struct {
		Previous []Event
		Current []Event
	}
	Cursor struct {
		Visible bool
		Ticker Ticker
		X int32
		Y int32
	}
	Ticker struct {
		LastTriggered float64
		Current float64
		Delay float64
		Rate float64
	}
	Scrollback struct {
		History [][]rune
		View [][]rune
		Pos int32
	}
	KeysState struct {
		Keys map[int32]Key
		LastSeen []int32
	}
  State struct {
		Mode Mode
		PreviousMode Mode
		Buf []rune
		CmdBuf []rune
		Exit bool
		Keys KeysState
		Events Events
		Error error
		InputView []rune
		Cursor Cursor 
		Font rl.Font
		Scrollback Scrollback
	}
)

func NewKey(b byte) Key {
	return Key {
		Ticker: &Ticker {
			Rate: 0.05,
			Delay: 0.4,
		},
		Byte: b,
	}
}
