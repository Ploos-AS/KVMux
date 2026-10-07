# M0 API quick start

Run:

```sh
go run ./cmd/kvmux
```

Default endpoint: `http://127.0.0.1:8080`.

Examples:

```sh
curl http://127.0.0.1:8080/api/v1/nodes
curl http://127.0.0.1:8080/api/v1/images

curl -X POST http://127.0.0.1:8080/api/v1/nodes/node-01/power/on
curl -X POST http://127.0.0.1:8080/api/v1/nodes/node-01/power/cycle
curl -X POST http://127.0.0.1:8080/api/v1/nodes/node-01/reset

curl -X POST \
  -H 'Content-Type: application/json' \
  -d '{"image":"demo-os"}' \
  http://127.0.0.1:8080/api/v1/nodes/node-01/media/attach

curl -X POST \
  -H 'Content-Type: application/json' \
  -d '{"image":"demo-os"}' \
  http://127.0.0.1:8080/api/v1/nodes/node-01/provision
```

The current M0 provisioning endpoint is deliberately synchronous and simulated.
The next M0 increment replaces it with an asynchronous job/state-machine API
and moves hardware operations behind backend interfaces.
