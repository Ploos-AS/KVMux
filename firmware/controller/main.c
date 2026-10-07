#include <stdio.h>
#include <stdlib.h>
#include <string.h>
#include "pico/stdlib.h"

#define PORTS 2
#define MIN_CYCLE_MS 100
#define MAX_CYCLE_MS 60000
static const uint power_gpio[PORTS] = {2, 3};
static const uint reset_gpio[PORTS] = {4, 5};
static bool power_state[PORTS] = {false, false};

static void respond(const char *id, const char *result) {
    printf("%s %s\n", id, result);
}
static bool number(const char *s, long *out) {
    if (!s || !*s) return false;
    char *end;
    long n = strtol(s, &end, 10);
    if (*end) return false;
    *out = n;
    return true;
}
static void execute(char *line) {
    char *id = strtok(line, " \t");
    char *cmd = strtok(NULL, " \t");
    if (!id || !cmd) { respond("0", "ERR MALFORMED"); return; }
    char *p = strtok(NULL, " \t");
    char *duration = strtok(NULL, " \t");
    char *extra = strtok(NULL, " \t");
    if (!strcmp(cmd, "PING") && !p) { respond(id, "OK PONG"); return; }
    if (!strcmp(cmd, "ALL_SAFE") && !p) { respond(id, "OK SAFE"); return; }
    long port;
    if (!number(p, &port) || port < 1 || port > PORTS) {
        respond(id, "ERR INVALID_PORT"); return;
    }
    unsigned idx = (unsigned)(port - 1);
    if (!strcmp(cmd, "POWER_ON") && !duration) {
        gpio_put(power_gpio[idx], 1); power_state[idx] = true;
    } else if (!strcmp(cmd, "POWER_OFF") && !duration) {
        gpio_put(power_gpio[idx], 0); power_state[idx] = false;
    } else if (!strcmp(cmd, "POWER_CYCLE") && duration && !extra) {
        long ms;
        if (!number(duration, &ms) || ms < MIN_CYCLE_MS || ms > MAX_CYCLE_MS) {
            respond(id, "ERR INVALID_INTERVAL"); return;
        }
        gpio_put(power_gpio[idx], 0); power_state[idx] = false;
        sleep_ms((uint32_t)ms);
        gpio_put(power_gpio[idx], 1); power_state[idx] = true;
    } else if (!strcmp(cmd, "RESET") && duration && !extra) {
        long ms;
        if (!number(duration, &ms) || ms < 1 || ms > 5000) {
            respond(id, "ERR INVALID_PULSE"); return;
        }
        if (!power_state[idx]) { respond(id, "ERR POWERED_OFF"); return; }
        gpio_put(reset_gpio[idx], 0);
        sleep_ms((uint32_t)ms);
        gpio_put(reset_gpio[idx], 1);
    } else {
        respond(id, "ERR UNKNOWN_COMMAND"); return;
    }
    respond(id, "OK DONE");
}
int main(void) {
    // Set output latch before enabling output: no reset pulse during boot.
    for (unsigned i = 0; i < PORTS; ++i) {
        gpio_init(power_gpio[i]); gpio_put(power_gpio[i], 0);
        gpio_set_dir(power_gpio[i], GPIO_OUT);
        gpio_init(reset_gpio[i]); gpio_put(reset_gpio[i], 1);
        gpio_set_dir(reset_gpio[i], GPIO_OUT);
    }
    stdio_init_all();
    char line[160];
    size_t used = 0;
    while (true) {
        int ch = getchar_timeout_us(10000);
        if (ch == PICO_ERROR_TIMEOUT) continue;
        if (ch == '\n' || ch == '\r') {
            if (used) { line[used] = 0; execute(line); used = 0; }
        } else if (used < sizeof(line)-1) {
            line[used++] = (char)ch;
        } else {
            used = 0;
            respond("0", "ERR LINE_TOO_LONG");
        }
    }
}
