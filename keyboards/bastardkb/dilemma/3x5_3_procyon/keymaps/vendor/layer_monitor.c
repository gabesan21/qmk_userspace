// SPDX-License-Identifier: GPL-3.0-or-later

#include "quantum.h"
#include "usb_main.h"
#include "via.h"

#include <string.h>

#define LAYER_MONITOR_HEARTBEAT_MS 500
#define LAYER_MONITOR_MAX_LAYER 6

extern usb_endpoint_in_t usb_endpoints_in[USB_ENDPOINT_IN_COUNT];

// Consulta avulsa pelo VIA; o monitor normal nunca abre essa interface.
enum monitor_status { NOT_RUN, SECONDARY, UNKNOWN_LAYER, USB_INACTIVE, SUSPENDED, PARTIAL, FULL, UNAVAILABLE, POSTED };
static uint8_t observed_layer = UINT8_MAX;
static enum monitor_status last_status;
static uint32_t posted;

void via_custom_value_command_kb(uint8_t *data, uint8_t length) {
    if (length < 16 || data[0] != id_custom_get_value || data[1] != id_custom_channel || data[2] != 0x4C) {
        if (length) data[0] = id_unhandled;
        return;
    }
    memcpy(data + 3, "DLD1", 4);
    data[7] = last_status;
    data[8] = observed_layer;
    for (uint8_t i = 0; i < 4; ++i) data[9 + i] = (uint8_t)(posted >> (8 * i));
}

static enum monitor_status send_layer(uint8_t layer) {
    const uint8_t           report[CONSOLE_EPSIZE] = {'D', 'L', 'M', '1', ':', (uint8_t)('0' + layer), '\n'};
    usb_endpoint_in_t      *endpoint               = &usb_endpoints_in[USB_ENDPOINT_IN_CONSOLE];
    output_buffers_queue_t *queue                  = &endpoint->obqueue;
    enum monitor_status status;

    // Nunca espera por um leitor USB nem descarta saída parcial de outro produtor.
    osalSysLock();
    if (usbGetDriverStateI(endpoint->config.usbp) != USB_ACTIVE) {
        status = USB_INACTIVE;
    } else if (queue->suspended) {
        status = SUSPENDED;
    } else if (queue->ptr != NULL) {
        status = PARTIAL;
    } else if (obqIsFullI(queue)) {
        status = FULL;
    } else if (obqGetEmptyBufferTimeoutS(queue, TIME_IMMEDIATE) != MSG_OK) {
        status = UNAVAILABLE;
    } else {
        memcpy(queue->ptr, report, sizeof(report));
        obqPostFullBufferS(queue, sizeof(report));
        status = POSTED;
    }
    osalSysUnlock();
    return status;
}

void housekeeping_task_user(void) {
    static uint32_t last_sent_at;
    static bool     pending = true;

    if (!is_keyboard_master()) {
        last_status = SECONDARY;
        return;
    }

    uint8_t layer = get_highest_layer(layer_state | default_layer_state);
    if (layer != observed_layer) {
        observed_layer = layer;
        pending        = true;
    }

    // A camada 7 do VIA não possui ID no protocolo do console.
    if (layer > LAYER_MONITOR_MAX_LAYER) {
        last_status = UNKNOWN_LAYER;
        return;
    }

    if (pending || timer_elapsed32(last_sent_at) >= LAYER_MONITOR_HEARTBEAT_MS) {
        last_status = send_layer(layer);
        if (last_status == POSTED) {
            ++posted;
            last_sent_at = timer_read32();
            pending = false;
        }
    }
}
