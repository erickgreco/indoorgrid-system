#ifndef BME680_SENSOR_H
#define BME680_SENSOR_H

#include <Arduino.h>
#include <Adafruit_BME680.h>

class BME680Sensor {
public:
    BME680Sensor(uint8_t address);
    bool begin();
    bool read();
    float getTemperature();
    float getHumidity();
    float getPressure();
    float getGasResistance();

private:
    Adafruit_BME680 bme;
    uint8_t address;
    float temperature;
    float humidity;
    float pressure;
    float gasResistance;
};

#endif