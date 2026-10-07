# M1 Controller Protocol

KVMux Core delegates bounded hardware operations to the controller MCU. The
Linux side owns policy, jobs, authentication and provisioning.

## Transport

M1 uses a framed request/response protocol over USB CDC. UART is a fallback
transport. The message model is transport-independent.

## Commands

- PING
- INFO
- PORT_STATUS <port>
- POWER_ON <port>
- POWER_OFF <port>
- POWER_CYCLE <port> <off_ms>
- RESET <port> <pulse_ms>
- UART_SELECT <port>
- SERVICE_SELECT <port>
- ALL_SAFE

Ports are numbered from 1.

## Rules

Commands carry a monotonically increasing request ID. Responses echo that ID.
Destructive commands are idempotent where possible. POWER_CYCLE enforces a
minimum off interval in firmware. Unknown ports and invalid timing values fail
closed.

The MCU must never infer provisioning policy. It must not reboot a target merely
because Linux disconnects.

## Safe startup

On MCU reset:
- reset outputs inactive
- service-path switches disconnected or in documented default state
- no target power transition is generated
- watchdog does not power-cycle targets by default

## Events

The MCU may asynchronously report:
- controller boot
- target power-state change
- reset assertion/release
- overcurrent/fault
- watchdog warning
- UART framing/overflow diagnostics

M1 begins with newline-delimited ASCII frames for bring-up and debugging. A
versioned binary framing format can replace it later without changing KVMux Core
backend semantics.
