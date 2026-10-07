# Architecture

## Principles

KVMux is an IP-KVM plus node lifecycle controller.

The architecture prioritizes:

1. unattended recovery
2. deterministic automation
3. modular hardware
4. SBC-specific provisioning without SBC-specific core code
5. inexpensive scaling
6. repairable/open hardware

## Control plane versus data plane

### Control plane

Runs on the management computer and owns:

- HTTP API
- authentication
- node inventory
- jobs/state machines
- image catalog
- audit/event log
- web UI
- CI integrations

### Hardware/data plane

Owns:

- HDMI/video switching and capture
- USB HID switching/emulation
- USB virtual-media presentation
- UART routing
- power switching
- reset lines
- optional voltage/current telemetry

## Controller split

A Linux management computer is appropriate for networking, storage and video.

A dedicated MCU is preferred for low-level control because power/reset/watchdog functionality should remain deterministic and recoverable even if Linux crashes.

Candidate responsibilities for the MCU:

- per-port power outputs
- reset outputs
- UART mux/control
- watchdog
- hardware interlocks
- sensor polling
- emergency all-off or defined recovery policy

Communication between Linux and MCU should use a documented protocol so the controller hardware can evolve independently.

## Shared video

KVMux does not require one capture device per target.

The initial design may switch one or a small number of capture pipelines among several HDMI sources while preserving always-on serial and power control for every port.

This is a major cost advantage over implementing eight independent PiKVM-class nodes.

## Virtual media

Virtual media needs two layers:

- image storage/catalog on the management computer
- USB-device presentation toward the selected target

The implementation may initially use Linux USB gadget-capable hardware or dedicated USB-device controllers. The abstraction must support multiple future implementations.

## Platform adapters

A generic KVM has no concept of SBC boot peculiarities. KVMux should.

Platform adapters may eventually describe:

```text
raspberry-pi
orange-pi
generic-uefi-arm64
generic-x86-uefi
custom-dev-board
```

A platform adapter can define:

- usable boot methods
- recovery/reset GPIO
- whether USB boot is available
- expected serial settings
- install/reimage strategy
- health detection hints

This keeps KVMux usable for both SBCs and ordinary tiny PCs.
