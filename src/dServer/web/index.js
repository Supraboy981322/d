const dom = document;
const body = dom.body;
const domain = window.location.origin;
const do_scroll = true;

function construct() {
  body.innerHTML = "";

  let page_container = dom.createElement("div");
  body.appendChild(page_container);
  page_container.id = "page_container";

  let board = dom.createElement("div");
  page_container.appendChild(board);
  board.id = "board";

  let msg_container = dom.createElement("div");
  board.appendChild(msg_container);
  msg_container.className = "container";

  let msg_box = dom.createElement("input");
  board.appendChild(msg_box);
  msg_box.className = "msg_box";
  msg_box.type = "text";
  msg_box.addEventListener("keydown", (event) => {
    if (event.key === "Enter") send();
  });

  let to_btm_btn = dom.createElement("button");
  body.appendChild(to_btm_btn);
  to_btm_btn.id = "to_btm";
  to_btm_btn.onclick = scroll;
  to_btm.innerText = "▼";

  update_board();
}
construct()

async function send(msg) {
  const d = new Date();
  let timestamp = `${d.getHours()}:${d.getMinutes()}:${d.getSeconds()}`;
  let input = dom.querySelector("#board > input.msg_box")
  if (msg === null || msg === undefined) msg = input.value;
  input.value = "";
  try {
    let resp = await fetch(`${domain}/post`, {
      method: "POST",
      body: msg,
    });
    if (!resp.ok) throw new Error(`SERVER ERR: ${resp}`);
  } catch (e) {
    alert(e);
    return;
  }
  new_msg_elem({
    Timestamp: timestamp,
    Msg: msg,
  });
}

async function update_board() {
  let resp;
  try {
    resp = await fetch(`${domain}/today`, {
      method: "GET",
    });
    if (!resp.ok) throw new Error(`SERVER ERR: ${resp}`);
  } catch (e) {
    alert(e);
    return;
  }

  let json = await resp.json();
  json.forEach((msg) => new_msg_elem(msg));
}

function scroll() {
  let board = dom.getElementById("board");
  if (do_scroll)
    board.scrollTop = board.scrollHeight; 
}

function new_msg_elem(msg) {
  let msg_container = dom.createElement("div");
  dom.querySelector("#board > .container").appendChild(msg_container);
  msg_container.id = msg.Timestamp;
  msg_container.className = "msg";

  let timestamp = dom.createElement("p");
  msg_container.appendChild(timestamp);
  timestamp.className = "timestamp";
  timestamp.innerText = msg.Timestamp;

  let msg_txt = dom.createElement("p");
  msg_container.appendChild(msg_txt);
  msg_txt.className = "txt";
  msg_txt.innerText = msg.Msg;
  scroll()
}

//setInterval(scroll, 100)
