#ifndef BH1750_SENSOR_H
#define BH1750_SENSOR_H

#include <Arduino.h>
#include <BH1750.h>

class BH1750Sensor {
public:
    BH1750Sensor(uint8_t address);
    bool begin();
    bool read();
    float getLux();

private:
    BH1750 lightMeter;
    uint8_t address;
    float lux;
};

#endif