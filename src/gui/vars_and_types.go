package main

import rl "github.com/gen2brain/raylib-go/raylib"

type KeyPair struct {
	Key int32
	Val byte
}
var Keys = map[int32]byte {
	rl.KeySpace: ' ',
	rl.KeyEscape: 19,
	rl.KeyEnter: 18,
	rl.KeyTab: 17,
	rl.KeyBackspace: 16,
	rl.KeyInsert: 15,
	rl.KeyDelete: 14,
	rl.KeyRight: 13,
	rl.KeyLeft: 12,
	rl.KeyDown: 11,
	rl.KeyUp: 10,
	rl.KeyPageUp: 9,
	rl.KeyPageDown: 8,
	rl.KeyHome: 7,
	rl.KeyEnd: 6,
	rl.KeyCapsLock: 5,
	rl.KeyLeftShift: 4,
	rl.KeyLeftControl: 3,
	rl.KeyLeftAlt: 2,
	rl.KeyLeftSuper: 1,
	rl.KeyRightShift: 4,
	rl.KeyRightControl: 3,
	rl.KeyRightAlt: 2,
	rl.KeyRightSuper: 1,
	rl.KeyLeftBracket: '[',
	rl.KeyBackSlash: '\\',
	rl.KeyRightBracket: ']',
	rl.KeyGrave: '~',
	rl.KeyApostrophe: '\'',
	rl.KeyComma: ',',
	rl.KeyMinus: '-',
	rl.KeyPeriod: '.',
	rl.KeySlash: '/',
	rl.KeyZero: '0',
	rl.KeyOne: '1',
	rl.KeyTwo: '2',
	rl.KeyThree: '3',
	rl.KeyFour: '4',
	rl.KeyFive: '5',
	rl.KeySix: '6',
	rl.KeySeven: '7',
	rl.KeyEight: '8',
	rl.KeyNine: '9',
	rl.KeySemicolon: ';',
	rl.KeyEqual: '=',
	rl.KeyA: 'a',
	rl.KeyB: 'b',
	rl.KeyC: 'c',
	rl.KeyD: 'd',
	rl.KeyE: 'e',
	rl.KeyF: 'f',
	rl.KeyG: 'g',
	rl.KeyH: 'h',
	rl.KeyI: 'i',
	rl.KeyJ: 'j',
	rl.KeyK: 'k',
	rl.KeyL: 'l',
	rl.KeyM: 'm',
	rl.KeyN: 'n',
	rl.KeyO: 'o',
	rl.KeyP: 'p',
	rl.KeyQ: 'q',
	rl.KeyR: 'r',
	rl.KeyS: 's',
	rl.KeyT: 't',
	rl.KeyU: 'u',
	rl.KeyV: 'v',
	rl.KeyW: 'w',
	rl.KeyX: 'x',
	rl.KeyY: 'y',
	rl.KeyZ: 'z',
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
		LastSeen []byte
		Repeat KeyRepeat
	}
	KeyRepeat struct {
		Delay float32
		Rate float32
		Timer float32
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
	}
	Scrollback struct {
		History [][]rune
		View [][]rune
		Pos int32
	}
  State struct {
		Mode Mode
		PreviousMode Mode
		Buf []rune
		CmdBuf []rune
		Exit bool
		Key Key
		Events Events
		Error error
		InputView []rune
		Cursor Cursor 
		Font rl.Font
		Scrollback Scrollback
	}
)
