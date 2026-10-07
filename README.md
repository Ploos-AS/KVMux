# KVMux

Open multiport IP KVM and lifecycle controller for SBCs, compute nodes and lab hardware.

KVMux is designed for Raspberry Pi, Orange Pi and other SBCs, but is intentionally not limited to ARM. It combines remote console access with the functions that matter in unattended labs and CI/render farms:

- HDMI/video KVM
- USB HID keyboard and mouse
- serial/UART console
- per-node power on/off/cycle
- per-node reset
- virtual USB mass storage
- OS image mounting and reinstall/provisioning
- automation-friendly HTTP API
- web UI and remote console
- telemetry and watchdog integration

## Project goal

A single KVMux controller should manage multiple nodes without requiring one complete IP-KVM computer per target.

The initial product target is **8 managed ports**, with a modular architecture that can later scale to 16/32 nodes.

Each fully managed port is intended to expose:

```text
HDMI/video
USB HID
USB virtual media
UART/serial
POWER
RESET
optional power/current telemetry
```

KVMux treats power control and OS reinstall as core functions rather than optional extensions.

## Node lifecycle

The intended unattended workflow is:

```text
select node
    ↓
power off / reset
    ↓
attach OS image as virtual USB storage
    ↓
select appropriate boot path
    ↓
power on / reboot
    ↓
watch serial/video console
    ↓
install or reimage
    ↓
detach installation media
    ↓
boot installed OS
```

This makes KVMux suitable for self-hosted CI runners, render nodes, test clusters and SBC labs.

## Architecture

```text
                    +---------------------------+
Ethernet ---------->| KVMux management computer |
                    | Linux                     |
                    | Web UI / API / SSH        |
                    +-------------+-------------+
                                  |
        +-------------------------+-------------------------+
        |                         |                         |
   video switch              USB subsystem           controller MCU
        |                   HID + storage             power/reset/
        |                         |                   watchdog/UART
   +----+----------------------------------------------------+
   |                    managed node ports                    |
   +----+-----------+-----------+-----------+----------------+
        |           |           |           |
      node 1      node 2      node 3      ... node 8
```

The management computer handles networking, UI, streaming, image management and API services. A dedicated microcontroller should eventually own safety-critical low-level functions such as power switching, reset and watchdog behaviour so these remain available even if the Linux management service fails.

## M0

M0 defines the architecture and provides a hardware-independent control-plane prototype.

See [docs/M0.md](docs/M0.md).

### M0 requirements

- node inventory and state model
- REST API skeleton
- power on/off/cycle semantics
- reset semantics
- virtual-media attach/detach model
- OS image catalog model
- install/reinstall workflow
- serial-console model
- video/HID abstraction
- simulated hardware backend
- hardware/backend interfaces for later boards

Physical switching and capture hardware begin after M0.

## Planned milestones

| Milestone | Scope |
| --- | --- |
| M0 | Architecture, API, node model and simulator |
| M1 | 2-port functional prototype |
| M2 | 4-port prototype / first PCB |
| M3 | 8-port KVMux |
| M4 | Power telemetry and watchdogs |
| M5 | Advanced virtual media and provisioning |
| M6 | Forgejo/GitHub runner integration |
| M7 | Rack/backplane version |

## API direction

Example future operations:

```text
GET  /api/v1/nodes
GET  /api/v1/nodes/{node}
POST /api/v1/nodes/{node}/power/on
POST /api/v1/nodes/{node}/power/off
POST /api/v1/nodes/{node}/power/cycle
POST /api/v1/nodes/{node}/reset

GET  /api/v1/images
POST /api/v1/nodes/{node}/media/attach
POST /api/v1/nodes/{node}/media/detach
POST /api/v1/nodes/{node}/provision

GET  /api/v1/nodes/{node}/serial
GET  /api/v1/nodes/{node}/video
```

The API is intended to be usable directly from CI orchestration, allowing a failed runner to be power-cycled or reinstalled automatically.

## Licensing

Planned project policy:

- software: MIT
- hardware/PCB/HDL: CERN-OHL-P-2.0
- documentation: CC BY 4.0
