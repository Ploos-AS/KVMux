# KVMux Controller Firmware

Firmware for the dedicated KVMux hardware controller.

M1 target: RP2350.

The controller owns bounded hardware operations such as target power, reset,
UART, watchdog/interlocks and service-path selection. Provisioning policy,
authentication and orchestration remain in KVMux Core on Linux.

Protocol: see `docs/M1-controller-protocol.md`.

Initial bring-up order:
1. PING / INFO
2. two LED-backed power outputs
3. two reset outputs
4. two UART channels
5. real protected power switches
6. USB/HDMI service selection

The firmware must boot into a safe state and must not toggle target power merely
because the Linux host resets or disconnects.
