#include <Arduino.h>
#include <Wire.h>
#include <time.h>
#include "credentials.h"
#include "mqtt_client.h"
#include "bme680_sensor.h"
#include "bh1750_sensor.h"

const char* TOPIC_BME680 = "indoorgrid/sensors/bme680";
const char* TOPIC_BH1750 = "indoorgrid/sensors/bh1750";

const unsigned long READ_INTERVAL = 60000;

MQTTClient mqtt(WIFI_SSID, WIFI_PASSWORD, MQTT_BROKER, MQTT_PORT);
BME680Sensor bme(0x77);
BH1750Sensor bh(0x23);

unsigned long lastRead = 0;

String getTimestamp() {
    struct tm timeinfo;
    if (!getLocalTime(&timeinfo)) {
        Serial.println("NTP: failed to get time");
        return "";
    }

    char buffer[30];
    strftime(buffer, sizeof(buffer), "%Y-%m-%dT%H:%M:%SZ", &timeinfo);
    return String(buffer);
}

void syncNTP() {
    configTime(-21600, 0, "pool.ntp.org");
    Serial.print("NTP: syncing");

    struct tm timeinfo;
    while (!getLocalTime(&timeinfo)) {
        delay(500);
        Serial.print(".");
    }

    Serial.println();
    Serial.println("NTP: synced");
}

void publishBME680(const String& timestamp) {
    if (!bme.read()) return;

    char payload[256];
    snprintf(payload, sizeof(payload),
        "{\"temperature_celcius\":%.2f,"
        "\"humidity_percent\":%.2f,"
        "\"pressure_hpa\":%.2f,"
        "\"gas_resistance_ohms\":%.2f,"
        "\"sensor_at\":\"%s\"}",
        bme.getTemperature(),
        bme.getHumidity(),
        bme.getPressure(),
        bme.getGasResistance(),
        timestamp.c_str()
    );

    mqtt.publish(TOPIC_BME680, payload);
    Serial.println(payload);
}

void publishBH1750(const String& timestamp) {
    if (!bh.read()) return;

    char payload[128];
    snprintf(payload, sizeof(payload),
        "{\"illuminance_lux\":%.2f,"
        "\"sensor_at\":\"%s\"}",
        bh.getLux(),
        timestamp.c_str()
    );

    mqtt.publish(TOPIC_BH1750, payload);
    Serial.println(payload);
}

void setup() {
    Serial.begin(115200);
    Wire.begin(21, 22);

    mqtt.connect();
    syncNTP();

    if (!bme.begin()) {
        Serial.println("BME680: halting");
        while (true) delay(1000);
    }

    if (!bh.begin()) {
        Serial.println("BH1750: halting");
        while (true) delay(1000);
    }
}

void loop() {
    mqtt.loop();

    if (millis() - lastRead >= READ_INTERVAL) {
        lastRead = millis();

        String timestamp = getTimestamp();
        if (timestamp.isEmpty()) return;

        publishBME680(timestamp);
        publishBH1750(timestamp);
    }
}