// SPDX-License-Identifier: GPL-3.0-or-later

#include "quantum.h"
#include "usb_main.h"

#include <string.h>

#define LAYER_MONITOR_HEARTBEAT_MS 500
#define LAYER_MONITOR_MAX_LAYER 6

extern usb_endpoint_in_t usb_endpoints_in[USB_ENDPOINT_IN_COUNT];

static bool send_layer(uint8_t layer) {
    const uint8_t           report[CONSOLE_EPSIZE] = {'D', 'L', 'M', '1', ':', (uint8_t)('0' + layer), '\n'};
    usb_endpoint_in_t      *endpoint               = &usb_endpoints_in[USB_ENDPOINT_IN_CONSOLE];
    output_buffers_queue_t *queue                  = &endpoint->obqueue;
    bool                    sent                   = false;

    // sendchar() can wait for 100 ms. Reserve one entire console report
    // without waiting, and leave any partially buffered QMK output untouched.
    osalSysLock();
    if (usbGetDriverStateI(endpoint->config.usbp) == USB_ACTIVE && !queue->suspended && queue->ptr == NULL && !obqIsFullI(queue)) {
        if (obqGetEmptyBufferTimeoutS(queue, TIME_IMMEDIATE) == MSG_OK) {
            memcpy(queue->ptr, report, sizeof(report));
            obqPostFullBufferS(queue, sizeof(report));
            sent = true;
        }
    }
    osalSysUnlock();
    return sent;
}

void housekeeping_task_user(void) {
    static uint8_t  observed_layer = UINT8_MAX;
    static uint32_t last_sent_at;
    static bool     pending = true;

    if (!is_keyboard_master()) {
        return;
    }

    uint8_t layer = get_highest_layer(layer_state | default_layer_state);
    if (layer != observed_layer) {
        observed_layer = layer;
        pending        = true;
    }

    // VIA can select layer 7, which has no ID in the monitor protocol.
    if (layer > LAYER_MONITOR_MAX_LAYER) {
        return;
    }

    if ((pending || timer_elapsed32(last_sent_at) >= LAYER_MONITOR_HEARTBEAT_MS) && send_layer(layer)) {
        last_sent_at = timer_read32();
        pending      = false;
    }
}
