"use strict";

//forgive me for putting everything in one file
//  the "system(s)" in JS for importing and exporting
//    files is absolute hell 

var conf; //set by 'set_config()'
var start_ok = true; //determines if background processes should run
var mode = "normal"; //tracks the mode
var total_entries = 0; //keeps of the current entries
var scroll_pos = 0;

//helper to check if JS is being dumb 
function exists(thing) {
  return thing !== undefined && thing !== null;
}

//sets the config at start
async function set_config() {
  console.log("getting config");
  if (!start_ok)
    document.body.innerHTML = ""; //clears the DOM

  //calls ipc handler to read config file
  try {
    conf = await window.api.get_config();
  } catch (e) {
    window.api.panic(e);
  }
  console.log("got config");

  //print an error and call 'process.exit(1);' if uncaught error
  if (!exists(conf))
    window.api.panic("CONFIG UNDEFINED");

  //makes sure user changed from the default value
  if (!conf.server || conf.server?.includes("your server address")) {
    //makes sure background processes didn't set to true
    start_ok = false;
    console.log(`server url is ${(conf.server) ? "bad" : "not set"}`);

    /* creates a popup element */
    
    //container
    let cont = document.createElement("div");
    document.body.appendChild(cont);
    cont.className = "conf_popup_cont";
    
    //message
    let msg = document.createElement("p");
    cont.appendChild(msg);
    msg.className = "msg";
    msg.innerText = "you don't appear to have set your server";

    //input label
    let input_title = document.createElement("p");
    cont.appendChild(input_title);
    input_title.className = "input_label";
    input_title.innerText = "please enter your server address";

    //input box
    let input_box = document.createElement("input");
    cont.appendChild(input_box);
    input_box.setAttribute("type", "text");
    input_box.addEventListener("keydown", (event) => {
      //<enter> key server is set and window closes 
      if (event.key === "Enter") set_server(cont);
    });

    //incase the user insists on pressing a button
    let done_btn = document.createElement("button");
    done_btn.innerText = "done";
    done_btn.addEventListener("click", () => set_server(cont)); //calls helper to verify
    cont.appendChild(done_btn);
    console.log("popup created"); 
  }

  //removes any trailing '/'
  if (conf.server.at(-1) === "/")
    conf.server = conf.server.slice(0, -1);

  //set the custom color palette (is that spelled correctly?)
  Object.entries(
    conf.colors ?? {}
  ).forEach(([name, val]) => {
    document.documentElement.style.setProperty(`--${name}`, val);
  });
}

//helper set the server in config
async function set_server(cont) {
  //makes sure that there's a validation result element
  let resp_msg = document.querySelector(".conf_popup_cont > #resp_msg");
  if (!resp_msg) {
    resp_msg = document.createElement("p");
    cont.appendChild(resp_msg);
    resp_msg.id = "resp_msg";
  }

  //makes sure the url isn't empty
  let url = document.querySelector(`.conf_popup_cont > input[type="text"]`).value;
  if (url === "") {
    resp_msg.innerText = "url is empty";
    return;
  }
  console.log(url);

  //attempts to make a request to the provided url
  try {
    let resp = await fetch(url);
    if (!resp.ok)
      throw new Error("url is invalid or unreachable");
    //sets validation result to ok
    resp_msg.innerText = "success reaching server";
    conf.server = url;
    //calls ipc handler to write config file
    await window.api.set_config(conf);
    resp_msg.innerText = "wrote config";
    start_ok = true;
    //removes the popup
    document.querySelector(".conf_popup_cont").remove();
    //and updates the board
    update_board();
  } catch (e) {
    resp_msg.innerText = e;
    return;
  }
}

//helper to startup the program
async function startup () {
  //checks and sets the config in memory
  await set_config();
  //builds to UI
  await construct();
  //only fetches board if config ok
  if (start_ok)
    await update_board();
  //updates the clock element every 1 second
  setInterval(clock, 1000);
  //sync board every 30 seconds by default 
  if (!conf.options?.reduce_requests)
    setInterval(() => sync_board(false), 30000);
  //set to insert mode if enabled in config
  if (conf.options?.start_inserted)
    set_mode("insert");
}
startup();

//helper to build the UI
async function construct() {
  //resets the body
  document.body.innerHTML = "";

  //creates the page container
  let page_container = document.createElement("div");
  document.body.appendChild(page_container);
  page_container.id = "page";

  //creates the board
  let board = document.createElement("div");
  page_container.appendChild(board);
  board.id = "board";

  //adds the message container to the board
  let msg_container = document.createElement("div");
  board.appendChild(msg_container);
  msg_container.className = "msg_container";

  //creates an entry box
  let msg_box = document.createElement("input");
  board.appendChild(msg_box);
  msg_box.className = "msg";
  msg_box.type = "text";
  msg_box.addEventListener("keydown", (event) => {
    //either send or run a command on enter key (depending on mode)
    if (event.key === "Enter") switch (mode) {
      case "insert": send(); break
      case "command": do_cmd(); break
    }
  });
  //set insert mode if clicked
  msg_box.addEventListener("click", () => set_mode("insert"));

  //button to go to bottom of board
  let to_btm_btn = document.createElement("button");
  document.body.appendChild(to_btm_btn);
  to_btm_btn.id = "to_btm";
  to_btm_btn.onclick = () => scroll(false, null, true);
  to_btm_btn.innerText = "\u{25BC}";
  
  //container for mode indicator
  let mode_indicator = document.createElement("p");
  document.body.appendChild(mode_indicator);
  mode_indicator.id = "mode";

  //sets the mode (to default)
  set_mode(mode);
  clock(); //renders the clock
}

//helper to send a message
async function send(msg) {
  //noop config is bad
  if (!start_ok) return;

  //get the message box
  let msg_box = document.querySelector("input.msg");
  //sets message if not provided, sets from input
  if (!exists(msg))
    msg = msg_box.value;

  //return if empty message
  if (!exists(msg) || msg === "")
    return;

  //reset message box
  msg_box.value = "";

  var msg_rendered; //holds the rendered (HTML) message
  try {
    //sends the message (entry)
    let resp = await fetch(`${conf.server}/post`, {
      method: "POST",
      headers: { "echo": "HTML" }, //tells server to send HTML
      body: msg,
    });

    //format error
    if (!resp.ok)
      throw new Error(
        `SERVER ERR: ${
          await resp.text()
        } (${
          resp.status
        } ; ${
          resp.error
        })`
      );

    //gets the repsonse as HTML
    let p = (new DOMParser())
          .parseFromString(await resp.text(), "text/html"),
        t = p.querySelector("p");
    //if the message is not a '<p>', it uses the '<body>' contents
    msg_rendered = (!t) ? p.body.innerHTML : t.innerHTML;
  } catch (e) {
    popup(e, true); //show err to user
    console.error(`send(): err{${e}} server{${conf.server}}`);
    return;
  }

  //creates a new entry on the board 
  new_msg_elem({
    Timestamp: mk_timestamp(),
    Msg: msg_rendered,
  });
}

//helper to sync the board
async function sync_board(force) {
  //no-op if not forcing and can't start or config disallows
  if ((!start_ok || conf.options.reduce_requests) && !force) return;

  //makes the request to server
  try {
    let resp = await fetch(`${conf.server}/sync`, {
      headers: {
        "have": total_entries,
      },
      method: "GET",
    });
    //formatted err if any problems 
    if (!resp.ok) {
      //return if server has same amount
      if (JSON.parse(resp.headers.get("have")) === total_entries)
        return;
      else if (exists(resp.headers.get("have"))) {
        update_board();
        return
      } else
        throw new Error(
          `SERVER ERR: ${
            await resp.text()
          } (${
            resp.status
          } ; ${
            resp.error
          })`
        );
    }
    //creates new elements 
    try {
      let json = await resp.json();
      if (json.length > 0) json.forEach(
        (msg) => new_msg_elem(msg)
      );
    } catch { return }
  } catch (e) {
    popup(e, true);
    console.error(`sync_board(): err{${e}} server{${conf.server}}`);
    return;
  }
}

//helper update (reset) the board
async function update_board() {
  if (!start_ok) return;
  try {
    //makes the request
    let resp = await fetch(`${conf.server}/today`, {
      method: "GET",
    });
    //formatted err
    if (!resp.ok)
      throw new Error(
        `SERVER ERR: ${
          await resp.text()
        } (${
          resp.status
        } ; ${
          resp.error
        })`
      );

    document.querySelectorAll(".msg_container > div.msg")?.forEach((e) => e.remove());
    //creates new elements
    let json = await resp.json();
    total_entries = 0;
    json.forEach((msg) => new_msg_elem(msg));
  } catch (e) {
    popup(e, true);
    console.error(`update_board(): err{${e}} server{${conf.server}}`);
    return;
  }
}

//helper to scroll the board
function scroll(to_top, n, smooth) {
  let board = document.querySelector("#board");
  if (smooth)
    board.setAttribute("style", "scroll-behavior: smooth;");
  if (!exists(n)) {
    n = (to_top) ? 0 : total_entries-1;
    scroll_pos = n;
  }

  document.querySelectorAll("div.msg[selected]").forEach((elem) => {
    elem.removeAttribute("selected");
  });

  let container = document.querySelector(".msg_container");
  let selected = container.children[n];
  selected?.scrollIntoView();
  selected?.setAttribute("selected", "");
  if (smooth)
    board.removeAttribute("style");

  let btn = document.getElementById("to_btm");

  //get the the last entry
  let last; {
    var i = 1; while (!exists(last) && i >= 0) {
      last = container.children[container.children.length - i];
      i--;
    }
  }

  if (container.children[scroll_pos] === last)
    btn?.remove();
  else if (!exists(btn)) {
    //button to go to bottom of board
    let btn = document.createElement("button");
    document.body.appendChild(btn);
    btn.id = "to_btm";
    btn.onclick = () => scroll(false, null, true);
    btn.innerText = "\u{25BC}";
  }
}

//helper to create a new message
function new_msg_elem(msg) {
  total_entries++;

  let msg_board = document.querySelector(".msg_container");
  let msg_container = document.createElement("div");
  msg_board.appendChild(msg_container);
  msg_container.id = msg.Timestamp;
  msg_container.className = "msg";
  msg_container.setAttribute("selected", "");

  let timestamp = document.createElement("p");
  msg_container.appendChild(timestamp);
  timestamp.className = "timestamp";
  timestamp.innerText = msg.Timestamp;

  let msg_txt = document.createElement("p");
  msg_container.appendChild(msg_txt);
  msg_txt.className = "txt";
  msg_txt.innerHTML = msg.Msg;

  scroll(false, null, true); //scrolls to the bottom
  scroll_pos = total_entries - 1;
}

//set the mode to insert if the input box is clicked (or normal if not)
document.addEventListener("click", (event) => {
  if (event.target.tagName === "INPUT" && event.target.className === "msg")
    set_mode("insert");
  else
    set_mode("normal");
});

//Vim bindings
var last_key = undefined; //keeps track for longer than 'event.repeat'
document.addEventListener("keydown", (event) => {
  //ignore stupid JS behavior 
  if (!exists(event.key))
    return;
  
  //ignores if the user is typing
  let currently_focused = document.activeElement.tagName.toLowerCase();
  if (["input", "textarea", "select"].includes(currently_focused)) {
    //unfocus element if escape key, otherwise back-off 
    if (event.key === "Escape") {
      event.target.blur();
      set_mode("normal");
    } else
      return
  }

  event.preventDefault();
  sw: switch (event.key) {
    //modes
    case "i": { set_mode("insert") } break sw;
    case ":": { set_mode("command") } break sw;

    //basic movement
    case "j": case "k": {
      if (event.key === "j" && scroll_pos < total_entries-1)
        scroll_pos += (event.repeat && scroll_pos - total_entries > 2) ? 2 : 1;
      else if (event.key === "k" && scroll_pos > 0)
        scroll_pos -= (event.repeat && scroll_pos > 1) ? 2 : 1;
      scroll(null, scroll_pos, !event.repeat);
    } break sw;

    //to start or end of scrollback
    case "G": { scroll(false, null, true) } break sw;
    case "g": { if (last_key === event.key) scroll(true, null, true) } break sw;

    //close any and all popups
    case "q": {
      [ ...get_all_elem("div.popup"),
        ...get_all_elem("div.error")
      ].forEach(
        (elem) => elem?.remove()
      );
    } break sw;
  }
  last_key = event.key;
});

//helper to set the mode
function set_mode(m) {
  mode = m;
  //focues or unfocuses the input bar if in inser to command mode
  let input = document.querySelector("input.msg");
  if (mode === "insert" || mode === "command")
    input.focus();
  else
    input.blur();
  document.querySelector("#mode").innerText = `--${mode}--`;
}

//string parsing in JS? HERESY
async function do_cmd() {
  let input = document.querySelector("input.msg");
  let p = {
    esc: false,
    stringing: false,
    str_type: undefined,
    res: [],
    mem: "",
    in: input.value.split(""),
  };
  loop: while (p.in.length > 0) {
    let char = p.in.shift();
    if (p.esc) {
      p.mem += char;
      p.esc = false;
      continue loop;
    }
    sw: switch (char) {
      case "'": case "\"": {
        if (p.str_type === char || !exists(p.str_type)) {
          if (p.stringing) {
            p.res.push(p.mem);
            p.mem = "";
            p.str_type = undefined;
            p.stringing = false;
          } else {
            p.stringing = true;
            p.str_type = char;
          }
        } else
          p.mem += char;
      } break sw;
      
      case " ": case "\t": case "\n": {
        if (!p.stringing) {
          p.res.push(p.mem);
          p.mem = "";
        } else
          p.mem += char;
      } break sw;

      case "\\": { p.esc = !p.esc } break sw;

      default: { p.mem += char }
    }
  }

  if (p.mem.length > 0)
    p.res.push(p.mem);
  input.value = "";
  if (p.res.length < 1)
    return;

  set_mode("normal"); 
  sw: switch (p.res[0]) {

    // TODO: save altered config options
    case "set": {
      if (p.res.length < 3)
        return popup("missing args: need something to set", false);
      switch (p.res[1]) {
        case "color": {
          if (p.res.length < 4)
            return popup(`missing arg: need a value for ${p.res[2]}`, false);
          document.documentElement.style.setProperty(`--${p.res[2]}`, p.res[3]);
        } break sw;
        case "server": {
          let old = p.res[2];
          conf.server = p.res[2];
          if (!(await chk_server())) {
            popup("couldn't reach server; reverting to previous server", false); 
            conf.server = old;
          }
        } break sw;
        default: { popup(`invalid arg: I don't know how to set "${p.res[1]}"`, false) };
      }
    } break sw;

    case "sync": { sync_board(true) } break sw;

    case "q": case ":q": { window.api.quit() } break sw;

    case "reload": case "refresh": {
      if (p.res.length < 2)
        return popup("missing args: need something to reload", false);
      switch (p.res[1]) {
        case "config", "conf": { set_config() } break sw;
        case "board": case "scrollback": case "scroll": { update_board() } break sw;
        default: { popup(`I do not know how to reload "${p.res[1]}"`, false) }
      }
    } break sw;

    default: { popup(`invalid command: "${p.res[0]}"`, false) }
  }
}

//helper to check if the configured server address is valid
async function chk_server() {
  let url = conf.server;
  if (!exists(url))
    return false;
  try {
    let resp = await fetch(url); 
    if (!resp.ok)
      throw new Error("url is invalid or unreachable");
    return true;
  } catch (e) {
    return false;
  }
}

//helper for the popup element
function popup(msg, is_error) {
  //creates the container
  let container = document.createElement("div");
  document.body.appendChild(container);
  container.className = (exists(is_error) && is_error) ? "error" : "popup";
  container.focus();
  
  //sets message to what's provided
  let msg_elem = document.createElement("p");
  container.appendChild(msg_elem);
  msg_elem.innerText = msg;

  //a button to close the popup if user doesn't use the vim-bindings
  let close_btn = document.createElement("button");
  container.appendChild(close_btn);
  close_btn.innerText = "close (popup)";
  close_btn.addEventListener("click", () => container.remove());

  //sets mode to normal
  set_mode("normal");
}

//helper to get an element
//  (mildly dumb, I know, I don't like typing the long name)
function get_elem(selector) {
  return document.querySelector(selector);
}

//helper to get an element
//  (mildly dumb, I know, I don't like typing the long name)
function get_all_elem(selector) {
  return document.querySelectorAll(selector);
}

//helper for the clock
function clock() {
  //creates the clock element if not present
  var time_elem = document.querySelector("#board > #time");
  if (!exists(time_elem)) {
    time_elem = document.createElement("div");
    document.querySelector("#board").appendChild(time_elem);
    time_elem.id = "time";

    let time_txt = document.createElement("p");
    time_elem.appendChild(time_txt);
    time_txt.className = "txt";
  }
  //sets the clock
  let time_txt = time_elem?.querySelector("p.txt");
  time_txt.innerText = mk_timestamp();
}

//helper to create a timestamp
function mk_timestamp() {
  return new Intl.DateTimeFormat("en-GB", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  }).format(
    Date.now()
  );
}
