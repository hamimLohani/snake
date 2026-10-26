# Snake IoT Hardware Controller Protocol Specification

## 1. Overview
The Snake CLI game engine is strictly decoupled from physical I/O. The `internal/input` package provides pluggable hardware controller support over **UDP/WiFi** or **USB Serial**.

---

## 2. Packet Formats

The controller driver accepts either **Single-Byte Binary** or **ASCII Plaintext** payloads.

### A. Single-Byte Binary Protocol (Recommended for Low Latency)
A minimal, 1-byte binary datagram:

| Byte Value | Hex    | Action / Direction |
|------------|--------|---------------------|
| `1`        | `0x01` | **UP**              |
| `2`        | `0x02` | **DOWN**            |
| `3`        | `0x03` | **LEFT**            |
| `4`        | `0x04` | **RIGHT**           |
| `5`        | `0x05` | **PAUSE / RESUME**  |
| `6`        | `0x06` | **BOSS SCREEN**     |

### B. ASCII Plaintext Protocol (Easy Debugging / Serial Monitor)
Text strings terminated by newline `\n` or single tokens (case-insensitive, whitespace stripped):

| Token / Word | Direction / Action | Notes |
|---|---|---|
| `"UP"`, `"U"`, `"W"` | **UP** | WASD `W` or directional initial |
| `"DOWN"`, `"S"` | **DOWN** | WASD `S` |
| `"LEFT"`, `"L"`, `"A"` | **LEFT** | WASD `A` or directional initial |
| `"RIGHT"`, `"R"`, `"D"` | **RIGHT** | WASD `D` or directional initial |

---

## 3. Communication Transports

### 1. UDP / WiFi
- **Default Port**: `4242` (UDP)
- **Host Target**: IP address of the macOS / Linux / Windows computer running `snake`.
- **Latency**: Sub-2ms over local 2.4GHz WiFi (802.11 b/g/n).
- **Packet Rate**: Recommended 20–50 packets/sec on joystick movement to prevent UDP socket buffer flooding.
- **CLI Flag**:
  ```bash
  snake play --hardware-udp :4242
  ```

### 2. USB Serial
- **Baud Rate**: `115200 bps`, 8 data bits, no parity, 1 stop bit (`8-N-1`).
- Connect ESP8266 / ESP32 via Micro-USB / USB-C to feed serial telemetry.

---

## 4. Hardware Pinout & Wiring

### A. ESP8266 (NodeMCU / D1 Mini)
| Component | ESP8266 Pin | Notes |
|---|---|---|
| Joystick X-axis (VRx) | `A0` (ADC0) | Analog input (0 – 1023) |
| Joystick Y-axis (VRy) | `GPIO4` (D2) | Or multiplexed / external ADC |
| Joystick Button (SW)  | `GPIO14` (D5) | Internal pull-up (`INPUT_PULLUP`) |
| VCC (Power)           | `3V3` or `5V` | Standard 3.3V power rail |
| GND (Ground)          | `GND` | Common system ground |

### B. ESP32 (DevKit / NodeMCU-32S)
| Component | ESP32 Pin | Notes |
|---|---|---|
| Joystick X-axis (VRx) | `GPIO34` (ADC1_CH6) | 12-bit Analog input (0 – 4095) |
| Joystick Y-axis (VRy) | `GPIO35` (ADC1_CH7) | 12-bit Analog input (0 – 4095) |
| Joystick Button (SW)  | `GPIO18` | Internal pull-up (`INPUT_PULLUP`) |
| VCC (Power)           | `3V3` | 3.3V rail |
| GND (Ground)          | `GND` | Common system ground |

---

## 5. Quick Start
1. Flash `hardware/esp8266_snake.ino` to your ESP8266 using Arduino IDE or PlatformIO.
2. Edit WiFi credentials `WIFI_SSID` and `TARGET_IP` in the sketch.
3. On your computer, run:
   ```bash
   snake play --hardware-udp :4242
   ```
4. Move the analog joystick to control the snake!
