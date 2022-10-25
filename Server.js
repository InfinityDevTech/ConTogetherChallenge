const express = require("express");
const { parse } = require("url");
const { WebSocketServer } = require("ws");
const fs = require("fs");

const app = express();
const wss1 = new WebSocketServer({ noServer: true });
let current = "";

wss1.on("connection", function connection(ws) {
  fs.readFile("hello.cloned", "utf8", (err, data) => {
    ws.send(data);
  });
  ws.on("close", function close() {
    console.log("disconnected");
  });
  ws.on("message", (m) => {
    m = JSON.parse(m);
    let bufferObj = Buffer.from(m.Fdata, "base64");
    if (bufferObj != current) {
      current = bufferObj;
      fs.writeFile(`${m.Fname}.cloned`, bufferObj, () => {});
      wss1.clients.forEach(function each(client) {
        if (client !== ws) {
          client.send(bufferObj);
        }
      });
    }
  });
});

const server = app.listen(8080);

server.on("upgrade", function upgrade(request, socket, head) {
  const { pathname } = parse(request.url);

  if (pathname === "/") {
    wss1.handleUpgrade(request, socket, head, function done(ws) {
      wss1.emit("connection", ws, request);
    });
  } else {
    socket.destroy();
  }
});
