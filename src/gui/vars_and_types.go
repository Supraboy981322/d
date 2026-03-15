package main

import rl "github.com/gen2brain/raylib-go/raylib"

type KeyPair struct {
	Key int32
	Val byte
}
var Keys = map[int32]byte {
	rl.KeySpace: ' ',
	rl.KeyEscape: 0,
	rl.KeyEnter: 0,
	rl.KeyTab: 0,
	rl.KeyBackspace: 0,
	rl.KeyInsert: 0,
	rl.KeyDelete: 0,
	rl.KeyRight: 0,
	rl.KeyLeft: 0,
	rl.KeyDown: 0,
	rl.KeyUp: 0,
	rl.KeyPageUp: 0,
	rl.KeyPageDown: 0,
	rl.KeyHome: 0,
	rl.KeyEnd: 0,
	rl.KeyCapsLock: 0,
	rl.KeyLeftShift: 0,
	rl.KeyLeftControl: 0,
	rl.KeyLeftAlt: 0,
	rl.KeyLeftSuper: 0,
	rl.KeyRightShift: 0,
	rl.KeyRightControl: 0,
	rl.KeyRightAlt: 0,
	rl.KeyRightSuper: 0,
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
		LastSeen int32
		Ticker Ticker
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
