package main

import (
	"io"
	"fmt"
	"time"
	"bytes"
	"slices"
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
		Keys: Key {
			Keys: Keys,
		},
	}
)

//so I can just willy-nilly (is that spelled correctly?)
//  insert a printf somewhere without worrying about it
func _(){fmt.Print()}

func main() {
	//create a window TODO: uhhh.... make better 
  rl.InitWindow(800, 450, "foo")
  defer rl.CloseWindow()
  rl.SetTargetFPS(60)

	// TODO: change this to something not hard-coded
	state.Font = rl.LoadFontEx("/run/current-system/sw/share/X11/fonts/CascadiaCodeNF-Regular.ttf", 20, nil, 0)

	//remove RayLib's exit key
	rl.SetExitKey(rl.KeyNull)

	//print "closing..." just before exiting
	defer func() {
		rl.BeginDrawing()
		rl.ClearBackground(rl.Black)
		//calculate the padding to center 
		left_pad := float32((800-(10 * 10))/2)
		top_pad := float32(rl.GetScreenHeight()/2)
		//print message 
		rl.DrawTextEx(
			state.Font,
			"closing...",
			rl.NewVector2(left_pad, top_pad),
			20,
			0,
			rl.RayWhite,
		)
		rl.EndDrawing()
		//wait 75ms before exiting
		time.Sleep(75 * time.Millisecond)
	}()

	//set the cursor's ticker
	state.Cursor.Ticker.LastTriggered = rl.GetTime()
	loop: for !rl.WindowShouldClose() {
		//skip the done label (used for cleanup at end of loop)
		goto start
		done: { //defer doesn't run until loop breaks, which is stupid
			state.PreviousMode = state.Mode
			state.Events.Previous = state.Events.Current
			state.Events.Current = []Event{}
			rl.EndDrawing()
			continue loop
		}
		//some stuff to prep for the loop
		start: {
			//if previous loop requested to exit, stop loop
			if state.Exit { break loop }

			rl.BeginDrawing()
			rl.ClearBackground(rl.Black)

			//set the cursor ticker's time 
			state.Cursor.Ticker.Current = rl.GetTime()
			//flip the cursor's visiblity on tick
			time_diff := state.Cursor.Ticker.Current - state.Cursor.Ticker.LastTriggered
			if time_diff >= state.Cursor.Ticker.Delay {
				state.Cursor.Visible = !state.Cursor.Visible
				state.Cursor.Ticker.LastTriggered = state.Cursor.Ticker.Current
			}

			//force the cursor's Y position to the input line if INSERT or CMD mode 
			switch state.Mode {
			  case INSERT, CMD: { state.Cursor.Y = int32(rl.GetScreenHeight() / 20) - 1 }
			}
		}

		var events []Event //holds current events 
		previous_err := state.Error //note previous error state 
		//handle key input and set current events and error state
		state.Error, events = handle_keys()
		//add events to state 
		state.Events.Add(events...)
		//make sure any previous error persists
		if previous_err != nil && state.Error == nil {
			state.Error = previous_err
		}

		//reset error state if escape triggered 
		if state.Events.Has(ESC) { state.Error = nil }

		//set exit on next frame if an event called for it
		//  BUG: for some reason, an error event seems to
		//  falsely trigger an exit event (this shouldn't happen)
		if state.Events.Has(EXIT) && !state.Events.Has(ERR) {
			state.Exit = true
		}

		//draw any errors 
		if state.Error != nil {
			//calculate the padding
			left_pad := float32(calc_w_centered(len(state.Error.Error())))
			rl.DrawTextEx(
				state.Font,
				state.Error.Error(),
				//print vertically and horizontally centered
				rl.NewVector2(left_pad, float32(rl.GetScreenHeight()/2)),
				20,
				0,
				rl.Red,
			)
			//skip all other rendering 
			goto done
		}
		
		//update the scrollback view
		state.Scrollback.View = state.Scrollback.History
		//update the max lines (incase window resized)
		state.Scrollback.MaxLines = int32(rl.GetScreenHeight() / 20)
		// TODO: cleanup
		for i := state.Scrollback.Pos; i > 0 && int32(len(state.Scrollback.View)) > state.Scrollback.MaxLines; i-- {
			pop(&state.Scrollback.View)
		}

		// TODO: remove this (for debugging)
		fmt.Printf(
			"pos{%d} screen{%d:%d} max{%d}\n",
			state.Scrollback.Pos,
			state.Scrollback.Pos, state.Scrollback.Pos + state.Scrollback.MaxLines,
			state.Scrollback.MaxLines,
		)

		//shift shift the scrollback view so the last MaxLines lines are shown 
		for len(state.Scrollback.View) > int(state.Scrollback.MaxLines) - 1 {
			shift(&state.Scrollback.View)
		}

		//only print the scrollback view if not hidden
		//  TODO: probably a better idea to only update scrollback view if not hidden
		if !state.Scrollback.Hide {
			//print each entry
			for i, line := range state.Scrollback.View {
				rl.DrawTextEx(
					state.Font,
					string(line), 
					// TODO: print from bottom up
					rl.NewVector2(10, float32(20 * i)),
					20,
					0,
					rl.RayWhite,
				)
			}
		}

		//print command buffer if CMD mode
		if state.Mode == Mode(CMD) {
			// TODO: horizontally scroll buffer
			rl.DrawTextEx(
				state.Font,
				string(state.CmdBuf),
				rl.NewVector2(0, float32((rl.GetScreenHeight()/20)-1)*20.0),
				20,
				0,
				rl.RayWhite,
			)
			state.Cursor.X = int32(len(state.CmdBuf))
		//print insert buffer if insert mode
		} else if state.Mode == Mode(INSERT) {

			//start with the full input view buffer
			// TODO: probably a better way to do this
			state.InputView = state.Buf
			//get the length buffer
			input_len := int32(longest_line_len(state.InputView))
			//shift view buffer so the end is within the window's bounds 
			if (input_len * 10) > int32(rl.GetScreenWidth()) {
				diff := input_len - int32(rl.GetScreenWidth() / 10)
				state.InputView = state.InputView[diff:]
			}
			//get the new length
			input_len = int32(longest_line_len(state.InputView))

			//calculate the left padding (to center it)
			left_pad := float32(calc_w_centered(int(input_len)))

			//draw the buffer
			rl.DrawTextEx(
				state.Font,
				string(state.InputView),
				//print at bottom center of window
				rl.NewVector2(left_pad, float32((rl.GetScreenHeight()/20)-1)*20.0),
				20,
				0,
				rl.RayWhite,
			)
			
			//move the cursor to the end of the buffer
			//  BUG: misaligned (I think something with the padding logic)
			state.Cursor.X = int32(left_pad / 10) + (input_len)
		}

		//draw cursor if should be visible
		if state.Cursor.Visible || state.Mode == Mode(NORMAL) || state.Mode == Mode(VISUAL) {
			rl.DrawRectangle(
				state.Cursor.X * 10,
				(state.Cursor.Y * 20),
				10,
				20,
				rl.RayWhite,
			)
		}

		//cleanup loop iteration
		goto done
  }
}

func cmd() (error, Event) {
	//reset the buffer on return 
	defer func() {
		state.CmdBuf = []rune{':'}
	}()

	//NOP on empty input
	if len(strings.TrimSpace(string(state.CmdBuf))) < 1 {
		return nil, Event(NOP)
	}

	//get the first arg (ignores the ':')  TODO: proper parsing
	first_trimmed := strings.Split(string(state.CmdBuf), " ")[0][1:]
	//switch on command name
	switch first_trimmed {
		//quit (ignores any remaining args)
		case "q": { state.Exit = true }
	  case "hide": { state.Scrollback.Hide = true }
	  case "show": { state.Scrollback.Hide = false }

		// TODO: more commands

	  default: {
			return errors.New("unknown command"), Event(ERR)
		}
	}

	//most commands are a NOP (at the moment)
	return nil, Event(NOP)
}

// WARN: spaget
func handle_keys() (error, []Event) {
	// NOTE: these are here to ensure that the somewhat janky logic
	//  for key repeating doesn't interfere with it
	//
	//'ctrl'+'c' exits as well TODO: remove this one
	if is_ctrl_pressed() && rl.IsKeyDown(rl.KeyC) {
		state.Exit = true
		return nil, []Event{ Event(EXIT) }
	//escape key sets the mode to NORMAL regardless of current mode
	} else if rl.IsKeyPressed(rl.KeyEscape) {
		defer func() {
			state.Mode = Mode(NORMAL)
			if len(state.CmdBuf) > 1 { state.CmdBuf = []rune{':'} }
			if state.Error != nil { state.Error = nil }
		}()
		return nil, []Event{ Event(ESC) }
	//NOP on no error
	} else if state.Error != nil {
		return nil, []Event{ Event(NOP) }
	}

	current_keys := GetKeysDown() 
	//loop over each key
	loop2: for rlKey, k := range state.Keys.Keys {
		//reset key state and skip if not pressed 
		if !rl.IsKeyDown(rlKey) {
			k.Ticker.LastTriggered = 0
			continue loop2
		}

		//local helper (might add more here) to update last seen 
		cleanup := func() {
			state.Keys.LastSeen = current_keys 
		}
		defer cleanup()

		//update the key's ticker timer
		k.Ticker.Current = rl.GetTime()
		//check if it's a repeated key
		isRepeat := slices.Contains(state.Keys.LastSeen, rlKey)
		if !isRepeat {
			//update last triggered and set the key's delay to a slower rate
			k.Ticker.LastTriggered = k.Ticker.Current
			k.Ticker.Delay = 0.4
		} else {
			//if enough time has passed on the key's ticker, set the delay
			//  to the ticker's faster rate
			time_diff := k.Ticker.Current - k.Ticker.LastTriggered
			if time_diff >= k.Ticker.Delay {
				k.Ticker.LastTriggered = k.Ticker.Current
				k.Ticker.Delay = k.Ticker.Rate
			//otherwise the it's still being delayed
			} else { goto done }
		}

		//ignore state for control key
		if is_ctrl_pressed() {
			//might (probably will) add more here
			switch rlKey {
			  case rl.KeyH: { state.Scrollback.Hide = !state.Scrollback.Hide }
			}
			goto done
		}

		//switch on the mode
		switch state.Mode {

			case Mode(INSERT): if k.Byte != 0 {
				//only add key to buffer if character
				state.Buf = append(state.Buf, rune(k.Byte))	
			} else {
				//otherwise (not a character), determine the action
				switch rlKey {
					case rl.KeyBackspace: { pop(&state.Buf) }
					case rl.KeyEnter: { return state.post() }
				}
			}

			case Mode(CMD): if k.Byte != 0 {
				//only add to command buffer if character
				state.CmdBuf = append(state.CmdBuf, rune(k.Byte))	
			} else {
				//otherwise (not a character), determine the action
				switch rlKey {
					//execute the command
					case rl.KeyEnter: {
						state.Mode = Mode(NORMAL)
						e, event := cmd()
						return e, []Event{ Event(CMD), event }
					}

					//delete from the buffer
					case rl.KeyBackspace: if len(state.CmdBuf) > 1 {
						pop(&state.CmdBuf)
					}
				}
			}

			case Mode(NORMAL): if k.Byte != 0 {
				//switch on the key character
				switch k.Byte {
					//insert mode
					case 'i': { state.Mode = Mode(INSERT) }
					//visual mode TODO: visual mode
					case 'v': { state.Mode = Mode(VISUAL) }
					//command mode (colon)
					case ';': if IsShiftDown() { state.Mode = Mode(CMD) }

					//basic cursor movement (left, down, up, right)
					case 'h', 'j', 'k', 'l': { goto vim_movements }

					// NOTE: not everything will be implemented
					default: {
						return errors.New("unknown action: " + string(k.Byte)), []Event{ Event(ERR) }
					}
				}
			} else { /*  TODO: */ }

			
		  case Mode(VISUAL): if k.Byte != 0 {
				// TODO: VISUAL mode actions
				switch k.Byte {
					//basic cursor movement (left, down, up, right)
					case 'h', 'j', 'k', 'l': { goto vim_movements }

					// NOTE: not everything will be implemented
					default: if k.Byte != 0 {
						return errors.New("unknown action: " + string(k.Byte)), []Event{ Event(ERR) }
					}
				}
			} else { /* TODO: */ }
		}
		
		//skip vim movements
		goto done

		vim_movements: {
			//left
			if rl.IsKeyDown(rl.KeyH) {
				// TODO: scroll one row horizontally
				if state.Cursor.X > 0 { state.Cursor.X-- }
			//down
			} else if rl.IsKeyDown(rl.KeyJ) {
				// TODO: fix vertical scrollback
				if state.Cursor.Y < state.Scrollback.MaxLines {
					state.Cursor.Y++
				} else {
					//only scroll of conditions are met
					scrollback_exceeds := len(state.Scrollback.History) > int(state.Scrollback.MaxLines)
					position_allows := state.Scrollback.Pos < state.Scrollback.MaxLines 
					can_scroll := position_allows && scrollback_exceeds
					if can_scroll {
						state.Scrollback.Pos++
						fmt.Println("down")
					}
				}
			//up
			} else if rl.IsKeyDown(rl.KeyK) {
				// TODO: fix vertical scrollback
				if state.Cursor.Y > 0 {
					state.Cursor.Y--
				} else {
					//only scroll of conditions are met
					scrollback_exceeds := len(state.Scrollback.History) > int(state.Scrollback.MaxLines)
					position_allows := state.Scrollback.Pos < state.Scrollback.MaxLines 
					can_scroll := position_allows && scrollback_exceeds
					if can_scroll {
						fmt.Println("up")
						state.Scrollback.Pos++
					}
				}
			//right
			} else if rl.IsKeyDown(rl.KeyL) {
				// TODO: scroll one row horizontally
				if state.Cursor.X < int32(rl.GetScreenWidth() / 10) {
					state.Cursor.X++
				}
			}
		}

		//cleanup and end loop NOTE: might change this
		done: {
			cleanup()
			break loop2
		}
	}
	return nil, []Event{ Event(NOP) }
}

func (s *State) post() (error, []Event) {
	var e error

	//create a byte buffer with the message 
	buf := bytes.NewBuffer([]byte(string(s.Buf)))
	// TODO: read from config/arg for server url
	req, e := http.NewRequest("POST", "http://[::1]:8008/post", buf)
	if e != nil { return e, []Event{ Event(ERR) } }
	req.Header.Set("Content-Type", "text/plain")

	//create a client and make request
	client := &http.Client{}
	res, e := client.Do(req)
	if e != nil { return e, []Event{ Event(ERR) } }
	defer res.Body.Close()

	//read the response  TODO: do something with this
	_, e = io.ReadAll(res.Body)
	if e != nil { return e, []Event{ Event(ERR) } }

	//add to scrollback
	s.Scrollback.History = append(s.Scrollback.History, s.Buf)
	s.Buf = []rune{}

	//return ok
	return nil, []Event{ Event(NOP) }
}
