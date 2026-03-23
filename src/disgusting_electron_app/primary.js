"use strict";

var conf; 
var start_ok = true;

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
  page_container.id = "p";

  let board = document.createElement("div");
  page_container.appendChild(board);
  board.id = "b";

  let msg_container = document.createElement("div");
  board.appendChild(msg_container);
  msg_container.className = "c";

  let msg_box = document.createElement("input");
  board.appendChild(msg_box);
  msg_box.className = "M";
  msg_box.type = "text";
  msg_box.addEventListener("keydown", event => {
    if (event.key === "Enter") send();
  });

  let to_btm_btn = document.createElement("button");
  document.body.appendChild(to_btm_btn)
  to_btm_btn.id = "B";
  to_btm_btn.onclick = scroll;
  to_btm_btn.innerText = "▼";
}

async function send(msg) {
  if (!start_ok) return
  let msg_box = document.querySelector(".M");
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

function scroll(){
  let board = document.querySelector("#b");
  board.scrollTop = board.scrollHeight; 
}

function new_msg_elem(msg) {
  let msg_board = document.querySelector(".c");
  let msg_container = document.createElement("div");
  msg_board.appendChild(msg_container);
  msg_container.id = msg.Timestamp;
  msg_container.className = "m";

  let timestamp = document.createElement("p");
  msg_container.appendChild(timestamp);
  timestamp.className = "T";
  timestamp.innerText = msg.Timestamp;

  let msg_txt = document.createElement("p");
  msg_container.appendChild(msg_txt);
  msg_txt.className = "t";
  msg_txt.innerHTML = msg.Msg;
  scroll()
}
