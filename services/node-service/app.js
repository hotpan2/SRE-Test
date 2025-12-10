const http = require('http');

const port = process.env.PORT || 8081;
const name = process.env.NODE_SERVICE_NAME || 'node-service';

const server = http.createServer((req, res) => {
  if (req.url === '/healthz') {
    res.writeHead(200, {'Content-Type': 'text/plain'});
    res.end('OK');
    return;
  }
  res.writeHead(200, {'Content-Type': 'text/plain'});
  res.end(`Hello from ${name}\n`);
});

server.listen(port, () => {
  console.log(`Node service listening on ${port}`);
});
