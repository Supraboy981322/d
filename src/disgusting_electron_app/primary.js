"use strict";

var conf; 
var start_ok = true;
var mode = "normal";
var total_entries = 0;

function exists(thing) {
  return thing !== undefined && thing !== null;
}

async function set_config() {
  console.log("getting config"); 
  document.body.innerHTML = "";

  try {
    conf = await window.api.get_config();
  } catch (e) {
    window.api.panic(e);
  }
  console.log("got config"); 

  if (!exists(conf))
    window.api.panic("CONFIG UNDEFINED");

  if (conf.server === "https://[your server address]") {
    start_ok = false;
    console.log("server url bad"); 

    let cont = document.createElement("div");
    document.body.appendChild(cont);
    cont.className = "conf_popup_cont";
    
    let msg = document.createElement("p");
    cont.appendChild(msg);
    msg.className = "msg";
    msg.innerText = "you don't appear to have set your server";

    let input_title = document.createElement("p");
    cont.appendChild(input_title);
    input_title.className = "input_label";
    input_title.innerText = "please enter your server address";

    let input_box = document.createElement("input");
    cont.appendChild(input_box);
    input_box.setAttribute("type", "text");
    input_box.addEventListener("keydown", event => {
      if (event.key === "Enter") set_server(cont);
    });

    let done_btn = document.createElement("button");
    done_btn.innerText = "done";
    done_btn.addEventListener("click", () => set_server(cont));
    cont.appendChild(done_btn);
    console.log("popup created"); 
  }

  if (conf.server[-1] === "/")
    conf.server = conf.server.slice(0, -1);

}

async function set_server(cont) {
  let resp_msg = document.querySelector(".conf_popup_cont > #resp_msg");
  if (!exists(resp_msg)) {
    resp_msg = document.createElement("p");
    cont.appendChild(resp_msg);
    resp_msg.id = "resp_msg";
  }

  let url = document.querySelector(`.conf_popup_cont > input[type="text"]`).value;
  if (url === "") {
    resp_msg.innerText = "url is empty";
    return;
  }
  console.log(url);

  try {
    let resp = await fetch(url); 
    if (!resp.ok)
      throw new Error("url is invalid or unreachable");
    resp_msg.innerText = "success reaching server";
    conf.server = url;
    await window.api.set_config(conf);
    resp_msg.innerText = "wrote config";
    start_ok = true;
    document.querySelector(".conf_popup_cont").remove();
    update_board();
  } catch (e) {
    resp_msg.innerText = e;
    return;
  }
}

async function startup () {
  await set_config();
  await construct();
  if (start_ok)
    await update_board();
}
startup();

async function construct() {
  let page_container = document.createElement("div");
  document.body.appendChild(page_container);
  page_container.id = "page";

  let board = document.createElement("div");
  page_container.appendChild(board);
  board.id = "board";

  let msg_container = document.createElement("div");
  board.appendChild(msg_container);
  msg_container.className = "msg_container";

  let msg_box = document.createElement("input");
  board.appendChild(msg_box);
  msg_box.className = "msg";
  msg_box.type = "text";
  msg_box.addEventListener("keydown", event => {
    if (event.key === "Enter") switch (mode) {
      case "insert": send();
      case "command": do_cmd();
    }
  });
  msg_box.addEventListener("click", () => set_mode("insert"));

  let to_btm_btn = document.createElement("button");
  document.body.appendChild(to_btm_btn);
  to_btm_btn.id = "to_btm";
  to_btm_btn.onclick = () => scroll(false);
  to_btm_btn.innerText = "▼";
  
  let mode_indicator = document.createElement("p");
  document.body.appendChild(mode_indicator);
  mode_indicator.id = "mode";
  set_mode(mode);
}

async function send(msg) {
  if (!start_ok) return;
  let msg_box = document.querySelector("input.msg");
  if (!exists(msg))
    msg = msg_box.value;

  //return if empty message
  if (!exists(msg) || msg === "")
    return;

  msg_box.value = "";

  var msg_rendered;
  try {
    let resp = await fetch(`${conf.server}/post`, {
      method: "POST",
      headers: { "echo": "HTML" },
      body: msg,
    });

    if (!resp.ok)
      throw new Error("SERVER ERR");

    let p = (new DOMParser())
          .parseFromString(await resp.text(), "text/html"),
        t = p.querySelector("p");
    msg_rendered = (t === null) ? p.body.innerHTML : t.innerHTML;
  } catch (e) {
    popup(e, true);
    console.error(`send(): err{${e}} server{${conf.server}}`);
    return;
  }

  //yikes, that's a lot of bytes just to get the time
  let timestamp = new Intl.DateTimeFormat("en-GB", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  }).format(
    new Date(Date.now())
  );

  new_msg_elem({
    Timestamp: timestamp,
    Msg: msg_rendered,
  });
}

setInterval(sync_board, 30000);

async function sync_board() {
  if (!start_ok) return;
  try {
    let resp = await fetch(`${conf.server}/sync`, {
      headers: {
        "have": total_entries,
      },
      method: "GET",
    });
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
            resp.code
          } ; ${
            resp.error
          })`
        );
    }
    try {
      let json = await resp.json();
      if (json.length > 0) json.forEach((msg) => {
        new_msg_elem(msg);
        console.log(msg);
      });
    } catch { return }
  } catch (e) {
    console.log(e);
    popup(e, true);
    console.error(`sync_board(): err{${e}} server{${conf.server}}`);
    return;
  }
}

async function update_board() {
  if (!start_ok) return;
  try {
    let resp = await fetch(`${conf.server}/today`, {
      method: "GET",
    });
    if (!resp.ok)
      throw new Error("SERVER ERR");

    let json = await resp.json();
    total_entries = 0;
    json.forEach(msg => new_msg_elem(msg));
  } catch (e) {
    popup(e, true);
    console.error(`update_board(): err{${e}} server{${conf.server}}`);
    return;
  }
}

function scroll(to_top){
  let board = document.querySelector("#board");
  board.setAttribute("style", "scroll-behavior: smooth;");
  if (to_top)
    document.querySelector("#board").scrollTop = 0;
  else
    board.scrollTop = board.scrollHeight; 
  board.removeAttribute("style");
}

function new_msg_elem(msg) {
  total_entries++;
  let msg_board = document.querySelector(".msg_container");
  let msg_container = document.createElement("div");
  msg_board.appendChild(msg_container);
  msg_container.id = msg.Timestamp;
  msg_container.className = "msg";

  let timestamp = document.createElement("p");
  msg_container.appendChild(timestamp);
  timestamp.className = "timestamp";
  timestamp.innerText = msg.Timestamp;

  let msg_txt = document.createElement("p");
  msg_container.appendChild(msg_txt);
  msg_txt.className = "txt";
  msg_txt.innerHTML = msg.Msg;
  scroll(false);
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
      document.querySelector("#board").scrollBy({
        top: (event.key === "j") ? 100 : -100,
        behavior: (event.repeat) ? undefined : "smooth",
      });
    } break sw;

    //to start or end of scrollback
    case "G": { scroll(false) } break sw;
    case "g": { if (last_key === event.key) scroll(true) } break sw;

    //close any and all popups
    case "q": {
      [ ...get_all_elem("div.popup"),
        ...get_all_elem("div.error")
      ].forEach(
        (elem) => elem?.remove()
      );
    } break sw;

    default: {
      last_key = event.key; //to bad this crappy language doesn't have defer
      return
    }
  }
  last_key = event.key;
});

function set_mode(m) {
  mode = m;
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
    case "set": {
      if (p.res.length < 3) {
        popup("missing args: need something to set", false);
        return;
      }
      switch (p.res[1]) {
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
    case "refresh": { update_board() } break sw;
    case "sync": { sync_board() } break sw;
    case "q": case ":q": { window.api.quit() } break sw;
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

function popup(msg, is_error) {
  let container = document.createElement("div");
  document.body.appendChild(container);
  container.className = (exists(is_error) && is_error) ? "error" : "popup";
  container.focus();
  
  let msg_elem = document.createElement("p");
  container.appendChild(msg_elem);
  msg_elem.innerText = msg;

  let close_btn = document.createElement("button");
  container.appendChild(close_btn);
  close_btn.innerText = "close (popup)";
  close_btn.addEventListener("click", () => container.remove());

  set_mode("normal");
}

function get_elem(selector) {
  return document.querySelector(selector);
}

function get_all_elem(selector) {
  return document.querySelectorAll(selector);
}
