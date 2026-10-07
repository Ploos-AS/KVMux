# M1 firmware qualification

Status: PENDING

The RP2350 build gate compiles Pico SDK firmware for `pico2` and uploads the
resulting UF2. This proves compilation only, not safe physical switching.

## Required evidence

- [ ] RP2350 workflow completed successfully on a pull request
- [ ] UF2 artifact exists and is nonempty
- [ ] Host Go tests and race detector pass
- [ ] USB CDC PING and invalid-port behavior verified on a Pico 2
- [ ] GPIO2/3 power outputs tested with LEDs, not directly with SBC loads
- [ ] GPIO4/5 reset outputs tested with a logic analyzer
- [ ] MCU reboot does not cause an unintended reset pulse
- [ ] Power-cycle off interval measured on real hardware

## Safety blockers

The current firmware uses push-pull GPIO for RESET_N. Do not connect this to
an unknown target reset line; use an open-drain transistor/level-safe interface.

The current firmware defaults PWR_EN low at boot, which may turn off an
already-running target. The final controller needs a defined external power
latch or persistent control architecture before unattended operation.

The current ALL_SAFE handler acknowledges the command but does not yet
physically disconnect the service path. SERVICE_SELECT is not implemented.

Do not label M1 hardware qualified until these are resolved and measured.
