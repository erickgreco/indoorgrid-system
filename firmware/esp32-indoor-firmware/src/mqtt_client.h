#ifndef MQTT_CLIENT_H
#define MQTT_CLIENT_H

#include <Arduino.h>
#include <WiFi.h>
#include <PubSubClient.h>

class MQTTClient {
public:
    MQTTClient(const char* ssid, const char* password, const char* broker, int port);
    void connect();
    void publish(const char* topic, const char* payload);
    void loop();
    bool isConnected();

private:
    const char* ssid;
    const char* password;
    const char* broker;
    int port;
    WiFiClient wifiClient;
    PubSubClient client;

    void connectWiFi();
    void connectBroker();
};

#endif