# M1 Two-Port Control Interface

## Scope

This is the first physical control-plane prototype. HDMI and USB service-path
hardware can be added after power/reset/UART is qualified.

## Per target port

Signals:
- PWR_EN: controller output to protected power switch or adapter
- RESET_N: open-drain reset output where supported
- UART_TX: 3.3 V MCU output to target RX
- UART_RX: 3.3 V target TX to MCU input
- GND: common signal reference
- optional PWR_SENSE
- optional PRESENT/ID

Never connect RS-232 voltage levels directly.

## Port 1 / Port 2 behavior

Both ports are independently addressable. UART is continuously available for
both ports; it is not tied to the selected HDMI/USB KVM service port.

Power operations:
- ON is idempotent
- OFF is idempotent
- CYCLE = OFF, enforced minimum off interval, ON
- RESET never substitutes for CYCLE

## Bench prototype

Recommended bring-up:
1. RP2350 development board
2. LEDs in place of PWR_EN and RESET_N
3. two USB-to-TTL adapters or two SBC UARTs
4. logic analyzer
5. protected MOSFET/load-switch modules for real target power

Do not power SBCs directly from RP2350 GPIO.

## Firmware abstraction

Each physical port exposes:
- setPower(bool)
- pulseReset(duration)
- readPowerState()
- uartRead()/uartWrite()
- readFaults()

The final PCB may use expanders or dedicated supervisors without changing the
host protocol.

## Qualification gates

- independently toggle two simulated power channels
- reset one port without affecting the other
- maintain both UART sessions
- reject invalid port numbers
- enforce minimum power-cycle off time
- survive Linux controller disconnect without target power changes
- reboot MCU without unintended target reset/power pulse
- operate two real SBC power channels safely
