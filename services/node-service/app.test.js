const request = require('supertest');
const server = require('./app');

afterAll(() => server.close());

describe('GET /healthz', () => {
  it('returns 200 with OK', async () => {
    const res = await request(server).get('/healthz');
    expect(res.status).toBe(200);
    expect(res.text).toBe('OK');
  });
});

describe('GET /', () => {
  it('returns 200 with greeting', async () => {
    const res = await request(server).get('/');
    expect(res.status).toBe(200);
    expect(res.text).toContain('Hello from');
  });

  it('uses NODE_SERVICE_NAME env var', async () => {
    process.env.NODE_SERVICE_NAME = 'my-test-service';
    const res = await request(server).get('/');
    expect(res.text).toContain('my-test-service');
    delete process.env.NODE_SERVICE_NAME;
  });
});

describe('GET /metrics', () => {
  it('returns 200 with prometheus metrics', async () => {
    const res = await request(server).get('/metrics');
    expect(res.status).toBe(200);
    expect(res.text).toContain('node_service_http_requests_total');
  });
});
