"use strict";

var conf; 
var start_ok = true;
var mode = "normal";

async function set_config() {
  console.log("getting config"); 
  document.body.innerHTML = "";

  try {
    conf = await window.api.get_config();
  } catch (e) {
    window.api.panic(e);
  }
  console.log("got config"); 

  if (conf === undefined)
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

  if (conf.server[-1] !== "/")
    conf.server = conf.server.slice(0, -1);

}

async function set_server(cont) {
  let resp_msg = document.querySelector(".conf_popup_cont > #resp_msg");
  if (resp_msg === undefined || resp_msg === null) {
    resp_msg = document.createElement("p");
    cont.appendChild(resp_msg);
    resp_msg.id = "resp_msg";
  }

  let url = document.querySelector(`.conf_popup_cont > input[type="text"]`).value;
  if (url === "") {
    resp_msg.innerText = "url is empty";
    return
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
    return
  }
}

async function startup () {
  await set_config();
  await construct();
  if (start_ok)
    await update_board();
}
startup()

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
  document.body.appendChild(to_btm_btn)
  to_btm_btn.id = "to_btm";
  to_btm_btn.onclick = () => scroll(false);
  to_btm_btn.innerText = "▼";
  
  let mode_indicator = document.createElement("p");
  document.body.appendChild(mode_indicator);
  mode_indicator.id = "mode";
  set_mode(mode);
}

async function send(msg) {
  if (!start_ok) return
  let msg_box = document.querySelector("input.msg");
  if (msg === undefined)
    msg = msg_box.value;

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
    msg_rendered = t === null ? p.body.innerHTML : t.innerHTML;
  } catch (e) {
    alert(e);
    return;
  }

  //yikes, that's a lot of bytes just to get the time
  let timestamp = new Intl.DateTimeFormat("en-GB", {
    hour: "2-digit",
    minute: "2-digit",
    second: "2-digit",
  }).format(
      new Date( Date.now() )
  );

  new_msg_elem({
    Timestamp: timestamp,
    Msg: msg_rendered,
  });
}

async function update_board() {
  if (!start_ok) return
  try {
    let resp = await fetch(`${conf.server}/today`, {
      method: "GET",
    });
    if (!resp.ok)
      throw new Error("SERVER ERR");

    let json = await resp.json();
    json.forEach(msg => new_msg_elem(msg));
  } catch (e) {
    alert(e);
    return;
  }
}

function scroll(to_top, n){
  let board = document.querySelector("#board");
  board.setAttribute("style", "scroll-behavior: smooth;");
  if (to_top)
    document.querySelector("#board").scrollTop = 0;
  else
    board.scrollTop = board.scrollHeight; 
  board.removeAttribute("style");
}

function new_msg_elem(msg) {
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
  scroll(false)
}


//Vim bindings
var last_key = undefined; //keeps track for longer than 'event.repeat'
document.addEventListener("keydown", (event) => {
  //ignore stupid JS behavior 
  if (event.key === undefined || event.key === null)
    return;
  
  let currently_focused = document.activeElement.tagName.toLowerCase();
  if (["input", "textarea", "select"].includes(currently_focused)) {
    //unfocus element if escape key, otherwise back-off 
    if (event.key === "Escape") {
      event.target.blur();
      set_mode("normal")
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
  document.querySelector("#mode").innerText = `--${mode}--`;
}

//string parsing in JS? HERESY
function do_cmd() {
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
        if (p.str_type === char || p.str_type === undefined) {
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
    p.res.push(p.mem)
  input.value = "";
  if (p.res.length < 1) 
    return

  set_mode("normal"); 
  sw: switch (p.res[0]) {
    case "set": {
      if (p.res.length < 2) {
        alert("missing args: need something to set");
        return
      }
      switch (p.res[1]) {
        case "server": { set_config() } break sw;
        default: { alert(`invalid arg: I don't know how to set "${p.res[1]}"`) };
      }
    } break sw;
    case "q": { window.api.quit() } break sw;
    default: { alert(`invalid command: "${p.res[0]}"`) }
  }
}
