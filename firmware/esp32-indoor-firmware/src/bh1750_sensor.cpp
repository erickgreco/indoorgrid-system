#include "bh1750_sensor.h"
#include <Wire.h>

BH1750Sensor::BH1750Sensor(uint8_t address)
    : lightMeter(address), address(address), lux(0) {}

bool BH1750Sensor::begin() {
    if (!lightMeter.begin(BH1750::CONTINUOUS_HIGH_RES_MODE)) {
        Serial.println("BH1750: sensor not found");
        return false;
    }

    Serial.println("BH1750: initialized");
    return true;
}

bool BH1750Sensor::read() {
    if (!lightMeter.measurementReady()) {
        return false;
    }

    lux = lightMeter.readLightLevel();

    if (lux < 0) {
        Serial.println("BH1750: failed to read");
        return false;
    }

    return true;
}

float BH1750Sensor::getLux() { return lux; }