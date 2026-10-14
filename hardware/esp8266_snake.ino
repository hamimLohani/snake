/*
 * ESP8266 / ESP32 WiFi & Serial Joystick Controller for Snake CLI
 *
 * Transmits real-time direction commands to the Snake game engine
 * over UDP (WiFi) and Serial (USB).
 *
 * Pin Connections:
 *   - Joystick VRx -> A0 (Analog 0)
 *   - Joystick SW  -> D5 (GPIO14, with internal pullup)
 *   - VCC          -> 3.3V
 *   - GND          -> GND
 */

#if defined(ESP8266)
  #include <ESP8266WiFi.h>
#elif defined(ESP32)
  #include <WiFi.h>
#endif
#include <WiFiUdp.h>

// ================= USER CONFIGURATION =================
const char* WIFI_SSID     = "YOUR_WIFI_SSID";
const char* WIFI_PASS     = "YOUR_WIFI_PASSWORD";
const char* TARGET_HOST   = "192.168.1.100"; // IP address of computer running Snake
const uint16_t TARGET_PORT = 4242;

// Protocol Command Bytes
#define CMD_UP     0x01
#define CMD_DOWN   0x02
#define CMD_LEFT   0x03
#define CMD_RIGHT  0x04
#define CMD_PAUSE  0x05
#define CMD_BOSS   0x06

// Hardware Pins
#define PIN_JOYSTICK_X A0
#define PIN_BUTTON     14 // D5 on NodeMCU

// Joystick Thresholds
#define THRESHOLD_LOW  300
#define THRESHOLD_HIGH 700

WiFiUDP udp;
uint8_t lastDirection = 0;
unsigned long lastSendTime = 0;
const unsigned long SEND_INTERVAL_MS = 60; // 60ms throttle

void setup() {
  Serial.begin(115200);
  delay(100);
  Serial.println("\n[Snake Controller] Booting...");

  pinMode(PIN_BUTTON, INPUT_PULLUP);

  // Connect to WiFi
  WiFi.mode(WIFI_STA);
  WiFi.begin(WIFI_SSID, WIFI_PASS);

  Serial.print("[Snake Controller] Connecting to WiFi");
  int attempts = 0;
  while (WiFi.status() != WL_CONNECTED && attempts < 20) {
    delay(500);
    Serial.print(".");
    attempts++;
  }

  if (WiFi.status() == WL_CONNECTED) {
    Serial.println("\n[Snake Controller] WiFi Connected! IP: " + WiFi.localIP().toString());
    udp.begin(TARGET_PORT);
  } else {
    Serial.println("\n[Snake Controller] WiFi connection failed. Fallback to Serial USB mode.");
  }
}

void sendCommand(uint8_t cmd) {
  // 1. Send via UDP if WiFi is active
  if (WiFi.status() == WL_CONNECTED) {
    udp.beginPacket(TARGET_HOST, TARGET_PORT);
    udp.write(cmd);
    udp.endPacket();
  }

  // 2. Also output to Serial (USB)
  Serial.write(cmd);
  Serial.println();
}

void loop() {
  unsigned long now = millis();

  // Read analog joystick
  int xVal = analogRead(PIN_JOYSTICK_X);
  uint8_t dir = 0;

  if (xVal < THRESHOLD_LOW) {
    dir = CMD_LEFT;
  } else if (xVal > THRESHOLD_HIGH) {
    dir = CMD_RIGHT;
  }

  // Check button press (Boss key or Pause)
  if (digitalRead(PIN_BUTTON) == LOW) {
    dir = CMD_BOSS;
    delay(200); // Debounce button
  }

  if (dir != 0 && (dir != lastDirection || (now - lastSendTime > SEND_INTERVAL_MS))) {
    sendCommand(dir);
    lastDirection = dir;
    lastSendTime = now;
  } else if (dir == 0) {
    lastDirection = 0;
  }

  delay(10);
}
