package main

import (
	"io"
	"fmt"
	"time"
	"bytes"
	"errors"
	"strings"
	"net/http"
	rl "github.com/gen2brain/raylib-go/raylib"
)

var (
	state = State{
		Scrollback: Scrollback {
			History: [][]rune{},
			View: [][]rune{},
		},
		Cursor: Cursor {
			Visible: true,
			Ticker: Ticker {
				Delay: 0.5,
			},
		},
		CmdBuf: []rune{':'},
		Key: Key {
			Ticker: Ticker {
				Delay: 0.05,
			},
		},
	}
)

func _(){fmt.Print()}

func main() {
  rl.InitWindow(800, 450, "foo")
  defer rl.CloseWindow()

	// TODO: change this to something not hard-coded
	state.Font = rl.LoadFontEx("/run/current-system/sw/share/X11/fonts/CascadiaCodeNF-Regular.ttf", 20, nil, 0)

  rl.SetTargetFPS(60)
	rl.SetExitKey(rl.KeyNull)

	//print "closing..." just before closing 
	defer func() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		left_pad := float32((800-(10 * 10))/2)
		top_pad := float32(rl.GetScreenHeight()/2)
		rl.DrawTextEx(
			state.Font,
			"closing...",
			rl.NewVector2(left_pad, top_pad),
			20,
			0,
			rl.RayWhite,
		)
		rl.EndDrawing()
		time.Sleep(75 * time.Millisecond)
	}()

	state.Cursor.Ticker.LastTriggered = rl.GetTime()
	state.Key.Ticker.LastTriggered = rl.GetTime()
	loop: for !rl.WindowShouldClose() {
		goto start
		done: { //defer doesn't run until loop breaks, which is stupid
			state.PreviousMode = state.Mode
			state.Events.Previous = state.Events.Current
			state.Events.Current = []Event{}
			rl.EndDrawing()
			continue loop
		}
		start: {
			if state.Exit { break loop }

			rl.BeginDrawing()
			rl.ClearBackground(rl.Black)

			state.Cursor.Ticker.Current = rl.GetTime()
			time_diff := state.Cursor.Ticker.Current - state.Cursor.Ticker.LastTriggered
			if time_diff >= state.Cursor.Ticker.Delay {
				state.Cursor.Visible = !state.Cursor.Visible
				state.Cursor.Ticker.LastTriggered = state.Cursor.Ticker.Current
			}

			switch state.Mode {
			  case INSERT, CMD: { state.Cursor.Y = int32(rl.GetScreenHeight() / 20) - 1 }
			}
		}

		var events []Event
		previous_err := state.Error
		state.Error, events = handle_keys()
		state.Events.Add(events...)
		if previous_err != nil && state.Error == nil {
			state.Error = previous_err
		}

		if state.Events.Has(ESC) { state.Error = nil }
		if state.Events.Has(EXIT) && !state.Events.Has(ERR) {
			state.Exit = true
		}

		if state.Error != nil {

			left_pad := float32(calc_w_centered(len(state.Error.Error())))

			rl.DrawTextEx(
				state.Font,
				state.Error.Error(),
				rl.NewVector2(left_pad, float32(rl.GetScreenHeight()/2)),
				20,
				0,
				rl.Red,
			)

			goto done
		}
		
		state.Scrollback.View = state.Scrollback.History
		max_lines := (rl.GetScreenHeight() / 20)
		exceeds_screen := len(state.Scrollback.View) > max_lines - 1
		exceeds_scrollback := int32(len(state.Scrollback.View)) < state.Scrollback.Pos
		for exceeds_screen || exceeds_scrollback {
			shift(&state.Scrollback.View)
		}
		for i, line := range state.Scrollback.View {
			rl.DrawTextEx(
				state.Font,
				string(line), 
				rl.NewVector2(10, float32(20 * i)),
				20,
				0,
				rl.RayWhite,
			)
		}

		if state.Mode == Mode(CMD) {
			rl.DrawTextEx(
				state.Font,
				string(state.CmdBuf),
				rl.NewVector2(0, float32((rl.GetScreenHeight()/20)-1)*20.0),
				20,
				0,
				rl.RayWhite,
			)
			state.Cursor.X = int32(len(state.CmdBuf))
		} else if state.Mode == Mode(INSERT) {

			state.InputView = state.Buf
			input_len := int32(longest_line_len(state.InputView))

			if (input_len * 10) > int32(rl.GetScreenWidth()) {
				diff := input_len - int32(rl.GetScreenWidth() / 10)
				state.InputView = state.InputView[diff:]
			}
			input_len = int32(longest_line_len(state.InputView))


			left_pad := float32(calc_w_centered(int(input_len)))
			rl.DrawTextEx(
				state.Font,
				string(state.InputView),
				rl.NewVector2(left_pad, float32((rl.GetScreenHeight()/20)-1)*20.0),
				20,
				0,
				rl.RayWhite,
			)
			state.Cursor.X = int32(left_pad / 10) + (input_len)
		}

		//draw cursor
		if state.Cursor.Visible || state.Mode == Mode(NORMAL) || state.Mode == Mode(VISUAL) {
			rl.DrawRectangle(
				state.Cursor.X * 10,
				(state.Cursor.Y * 20),
				10,
				20,
				rl.RayWhite,
			)
		}
		goto done
  }
}

func cmd() (error, Event) {
	defer func() {
		state.CmdBuf = []rune{':'}
	}()

	if len(strings.TrimSpace(string(state.CmdBuf))) < 1 {
		return nil, Event(NOP)
	}

	first_trimmed := strings.Split(string(state.CmdBuf), " ")[0][1:]
	switch first_trimmed {
		case "q": { state.Exit = true }
	  default: {
			return errors.New("unknown command"), Event(ERR)
		}
	}

	return nil, Event(NOP)
}

// WARN: spaget
func handle_keys() (error, []Event) {
	if is_ctrl_pressed() && rl.IsKeyDown(rl.KeyC) {
		state.Exit = true
		return nil, []Event{ Event(EXIT) }
	} else if rl.IsKeyPressed(rl.KeyEscape) {
		defer func() {
			state.Mode = Mode(NORMAL)
			if len(state.CmdBuf) > 1 { state.CmdBuf = []rune{':'} }
			if state.Error != nil { state.Error = nil }
		}()
		return nil, []Event{ Event(ESC) }
	} else if state.Error != nil {
		return nil, []Event{ Event(NOP) }
	} else if state.Mode == Mode(NORMAL) || state.Mode == Mode(VISUAL) {
		loop: for _, k := range []int32{ rl.KeyH, rl.KeyJ, rl.KeyK, rl.KeyL } {
			if rl.IsKeyDown(k) {

				state.Key.Ticker.Current = rl.GetTime()
				time_diff := state.Key.Ticker.Current - state.Key.Ticker.LastTriggered
				shouldRepeat := time_diff >= state.Key.Ticker.Delay
				isRepeat := state.Key.LastSeen == k
				if (isRepeat && shouldRepeat) || !isRepeat {
					state.Key.Ticker.LastTriggered = state.Key.Ticker.Current
					goto sw
				} else if isRepeat { continue loop }

				sw: switch k {
					// TODO: scroll one row horizontally
					case rl.KeyH: if state.Cursor.X > 0 { state.Cursor.X-- }
				
					// TODO: fix vertical scrollback
					case rl.KeyK: if state.Cursor.Y > 0 { state.Cursor.Y-- } else {
						if state.Scrollback.Pos > 0 {
							fmt.Println("up")
							state.Scrollback.Pos--
						}
					}

					// TODO: fix vertical scrollback
					case rl.KeyJ: if state.Cursor.Y < int32(rl.GetScreenHeight() / 20) {
						state.Cursor.Y++
					} else {
						if state.Scrollback.Pos < int32(len(state.Scrollback.History)) {
							fmt.Println("down")
							state.Scrollback.Pos++
						}
					}

					// TODO: scroll one row horizontally
					case rl.KeyL: if state.Cursor.X < int32(rl.GetScreenWidth() / 10) {
						state.Cursor.X++
					}
				}

				state.Key.LastSeen = k
			}
		}
	}

	switch state.Mode {

	  case Mode(INSERT): if k := rl.GetCharPressed(); k != 0 {
			state.Buf = append(state.Buf, rune(k))	
		} else {
			switch rl.GetKeyPressed() {
				case rl.KeyBackspace: { pop(&state.Buf) }
				case rl.KeyEnter: { return state.post() }
			}
		}

	  case Mode(CMD): if k := rl.GetCharPressed(); k != 0 {
			state.CmdBuf = append(state.CmdBuf, rune(k))	
		} else {
			switch rl.GetKeyPressed() {
				case rl.KeyEnter: {
					state.Mode = Mode(NORMAL)
					e, event := cmd()
					return e, []Event{ Event(CMD), event }
				}
				case rl.KeyBackspace: if len(state.CmdBuf) > 1 {
					pop(&state.CmdBuf)
				}
			}
		}

		case Mode(NORMAL): if k := rl.GetCharPressed(); k != 0 {
			switch k {
				case 'i': { state.Mode = Mode(INSERT) }
				case 'v': { state.Mode = Mode(VISUAL) }
		    case ':': { state.Mode = Mode(CMD) }
			}
		} else { /*  TODO: */ }

	}
	return nil, []Event{ Event(NOP) }
}

func (s *State) post() (error, []Event) {
	var e error

	buf := bytes.NewBuffer([]byte(string(s.Buf)))
	req, e := http.NewRequest("POST", "http://[::1]:8008/post", buf)
	if e != nil { return e, []Event{ Event(ERR) } }

	req.Header.Set("Content-Type", "text/plain")
	client := &http.Client{}
	res, e := client.Do(req)
	if e != nil { return e, []Event{ Event(ERR) } }
	defer res.Body.Close()

	_, e = io.ReadAll(res.Body)
	if e != nil { return e, []Event{ Event(ERR) } }

	s.Scrollback.History = append(s.Scrollback.History, s.Buf)
	s.Buf = []rune{}

	return nil, []Event{ Event(NOP) }
}
