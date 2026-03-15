package main

import (
	"time"
	"bytes"
	"slices"
	"strings"
 	rl "github.com/gen2brain/raylib-go/raylib"
)

/*
 * rendering related helpers
 */
//helper to calculate the vertical padding to center text
func calc_v_centered(h int) int32 {
	height := rl.GetScreenHeight()
	return int32((height/2)-((h/2)*20))
}
//helper to calculate the horizontal padding to center text
func calc_w_centered(w int) int32 {
	width := rl.GetScreenWidth()
	return int32((width-(10 * w))/2)
}


/*
 * slice related helpers
 */
//helper to remove last item from slice
//   NOTE: having to do this is moronic:
//   'slice = slice[:len(slice)-1]'
func pop[S ~[]T, T any](buf *S) {
	if len(*buf) > 0 {
		*buf = (*buf)[:len(*buf)-1]
	}
}
//helper to remove first item from slice
//   NOTE: see 'pop(...)'
func shift[S ~[]T, T any](buf *S) {
	if len(*buf) > 0 {
		*buf = (*buf)[1:]
	}
}


/*
 * keyboard input related helpers
 */
//helper to check if 'shift' is being held
//  regardless of which
func IsShiftDown() bool {
	return rl.IsKeyDown(rl.KeyLeftShift) || rl.IsKeyDown(rl.KeyRightShift)
}
//helper to check if 'ctrl' is being held
//  regardless of which
func is_ctrl_pressed() bool {
	return rl.IsKeyDown(rl.KeyLeftControl) || rl.IsKeyDown(rl.KeyRightControl)
}
//helper to get a slice of keys currently down
func GetKeysDown() []int32 {
	var res []int32
	for k := range Keys {
		if rl.IsKeyDown(k) { res = append(res, k) }
	}
	return res
}
//helper to get just the chars of keys being held 
func GetCharsDown() []byte {
	var res []byte
	for k, b := range Keys {
		if rl.IsKeyDown(k) { res = append(res, b.Byte) }
	}
	return res
}
//helper to check if a character key is being pressed
//  from the char
func IsCharPressed(b byte) bool {
	return bytes.Contains(GetCharsDown(), []byte{b})
}

/*
 * event related helpers
 */
//helper to add an event
func (es *Events) Add(e ...Event) {
	es.Current = append(es.Current, e...)
}
//helper to check if an event happened previously
//  (previous frame)
func (es *Events) Did(e Event) bool {
	return slices.Contains(es.Previous , e)
}
//helper to check if an event just happened
//  (current frame)
func (es *Events) Has(e Event) bool {
	return slices.Contains(es.Current, e)
}


/*
 * Random helpers
 */
//crappy helper to wait for a bool to be true
//  NOTE: why is this not a built-in?
func wait(f func()bool) {
	for !f() {
		time.Sleep(1 * time.Millisecond)
	}
}
//helper to get the longest line in a buffer containing newlines
func longest_line_len(s []rune) int {
	longest := 0
	//range over each line
	for _, line := range strings.Split(string(s), "\n") {
		// TODO: probably better to just do 'len(line)'
		cur := 0
		for range line { cur++ }
		if cur > longest { longest = cur }
	}
	//return the longest line found
	return longest
}
