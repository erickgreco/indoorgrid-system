#include "bme680_sensor.h"

BME680Sensor::BME680Sensor(uint8_t address)
    : address(address), temperature(0), humidity(0), pressure(0), gasResistance(0) {}

bool BME680Sensor::begin() {
    if (!bme.begin(address)) {
        Serial.println("BME680: sensor not found");
        return false;
    }

    bme.setTemperatureOversampling(BME680_OS_8X);
    bme.setHumidityOversampling(BME680_OS_2X);
    bme.setPressureOversampling(BME680_OS_4X);
    bme.setIIRFilterSize(BME680_FILTER_SIZE_3);
    bme.setGasHeater(320, 150);

    Serial.println("BME680: initialized");
    return true;
}

bool BME680Sensor::read() {
    if (!bme.performReading()) {
        Serial.println("BME680: failed to read");
        return false;
    }

    temperature = bme.temperature;
    humidity = bme.humidity;
    pressure = bme.pressure / 100.0;
    gasResistance = bme.gas_resistance;

    return true;
}

float BME680Sensor::getTemperature() { return temperature; }
float BME680Sensor::getHumidity() { return humidity; }
float BME680Sensor::getPressure() { return pressure; }
float BME680Sensor::getGasResistance() { return gasResistance; }