const express = require('express')
const { parse } = require('url')
const { WebSocketServer } = require('ws')
const fs = require('fs')

const app = express()
const wss1 = new WebSocketServer({ noServer: true });
const clients = []

wss1.on('connection', function connection(ws) {
  clients.push({id: ws.id, ws: ws})
  ws.on('close', function close() {
    console.log('disconnected');
    clients.filter((client) => client.id !== ws.id)
  })
   ws.on('message', (m) => {
    m = JSON.parse(m)
    let bufferObj = Buffer.from(m.Fdata, 'base64');
    fs.writeFile(`${m.Fname}.cloned`, bufferObj, () => {})
    clients.forEach((client) => {
      if (client.id !== ws.id) {
        client.ws.send(m)
      }
    })
   })
});

const server = app.listen(8080);

server.on('upgrade', function upgrade(request, socket, head) {
  const { pathname } = parse(request.url);

  if (pathname === '/') {
    wss1.handleUpgrade(request, socket, head, function done(ws) {
      wss1.emit('connection', ws, request);
    });
  } else {
    socket.destroy();
  }
});