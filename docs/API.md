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

Provisioning is asynchronous. A successful request returns HTTP `202 Accepted` and a job object.

Poll:

```sh
curl http://127.0.0.1:8080/api/v1/jobs/job-000001
```

Job states are currently `queued`, `running`, `completed` and `failed`.
The simulated provisioning state machine exposes steps including power-off,
attach-media, boot-installer, detach-media and verify.

Power, reset and virtual-media operations are behind backend interfaces. M1 can
therefore add physical implementations without changing the public lifecycle API.
