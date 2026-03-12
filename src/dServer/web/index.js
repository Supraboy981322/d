const dom = document;
const body = dom.body;

function construct() {
  let page_container = dom.createElement("div");
  body.appendChild(page_container);
  page_container.id = "page_container";

  let board = dom.createElement("div");
  page_container.appendChild(board);
  board.id = "board";

  let msg_box = dom.createElement("input");
  board.appendChild(msg_box);
  msg_box.className = "msg_box";
  msg_box.type = "text";
  msg_box.addEventListener("keydown", (event) => {
    if (event.key === "Enter") send();
  });
}
construct()

function send(msg) {
  let input = dom.querySelector("#board > input.msg_box")
  if (msg === null || msg === undefined) msg = input.value;
  input.value = "";
}

function new_msg_elm(msg) {
  let msg_container = dom.createElement("div");
  msg_container.id = msg.timestamp;
  msg_container.className = "msg";
}
