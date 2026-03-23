"use strict";



//I know some of the choices of syntax is (probably)
//  a little strange, but a lot of it was done to optimize
//    for compression size after minification



function construct() {
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

  update_board();
}

construct();

async function send(msg) {
  let msg_box = document.querySelector(".M");
  if (msg === undefined)
    msg = msg_box.value;

  msg_box.value = "";

  var msg_rendered;
  try {
    let resp = await fetch(`${window.location.origin}/post`, {
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
  try {
    let resp = await fetch(`${window.location.origin}/today`, {
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
