package main

import (
	"time"
	"bytes"
	"slices"
	"strings"
 	rl "github.com/gen2brain/raylib-go/raylib"
)

func calc_v_centered(h int) int32 {
	height := rl.GetScreenHeight()
	return int32((height/2)-((h/2)*20))
}

func calc_w_centered(w int) int32 {
	width := rl.GetScreenWidth()
	return int32((width-(10 * w))/2)
}

func longest_line_len(s []rune) int {
	longest := 0
	for _, line := range strings.Split(string(s), "\n") {
		cur := 0
		for range line { cur++ }
		if cur > longest { longest = cur }
	}
	return longest
}

func pop[S ~[]T, T any](buf *S) {
	if len(*buf) > 0 {
		*buf = (*buf)[:len(*buf)-1]
	}
}
func shift[S ~[]T, T any](buf *S) {
	if len(*buf) > 0 {
		*buf = (*buf)[1:]
	}
}

func wait(f func()bool) {
	for !f() {
		time.Sleep(1 * time.Millisecond)
	}
}

func IsShiftDown() bool {
	return rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)
}

func is_ctrl_pressed() bool {
	return rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl)
}

func GetKeysDown() []int32 {
	var res []int32
	for k := range Keys {
		if rl.IsKeyDown(k) { res = append(res, k) }
	}
	return res
}

func GetCharsDown() []byte {
	var res []byte
	for k, b := range Keys {
		if rl.IsKeyDown(k) { res = append(res, b.Byte) }
	}
	return res
}

func IsCharPressed(b byte) bool {
	return bytes.Contains(GetCharsDown(), []byte{b})
}

func (es *Events) Add(e ...Event) {
	es.Current = append(es.Current, e...)
}

func (es *Events) Did(e Event) bool {
	return slices.Contains(es.Previous , e)
}
func (es *Events) Has(e Event) bool {
	return slices.Contains(es.Current, e)
}
