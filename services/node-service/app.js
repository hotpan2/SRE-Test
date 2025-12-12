const http = require('http');
const client = require('prom-client');

const port = process.env.PORT || 8081;
const name = process.env.NODE_SERVICE_NAME || 'node-service';

// Enable default Node system metrics
client.collectDefaultMetrics();

// Custom counter
const httpRequests = new client.Counter({
  name: 'node_service_http_requests_total',
  help: 'Total number of HTTP requests received'
});

const server = http.createServer(async (req, res) => {
  if (req.url === '/healthz') {
    httpRequests.inc();
    res.writeHead(200, { 'Content-Type': 'text/plain' });
    res.end('OK');
    return;
  }

  if (req.url === '/metrics') {
    // Expose Prometheus metrics
    res.setHeader('Content-Type', client.register.contentType);
    return res.end(await client.register.metrics());
  }

  httpRequests.inc();
  res.writeHead(200, { 'Content-Type': 'text/plain' });
  res.end(`Hello from ${name}\n`);
});

server.listen(port, () => {
  console.log(`Node service listening on ${port}`);
});
