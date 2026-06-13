#include "mqtt_client.h"

MQTTClient::MQTTClient(const char* ssid, const char* password, const char* broker, int port)
    : ssid(ssid), password(password), broker(broker), port(port), client(wifiClient) {}

void MQTTClient::connect() {
    connectWiFi();
    connectBroker();
}

void MQTTClient::connectWiFi() {
    WiFi.begin(ssid, password);
    Serial.print("WiFi: connecting");

    while (WiFi.status() != WL_CONNECTED) {
        delay(500);
        Serial.print(".");
    }

    Serial.println();
    Serial.print("WiFi: connected - IP: ");
    Serial.println(WiFi.localIP());
}

void MQTTClient::connectBroker() {
    client.setServer(broker, port);

    while (!client.connected()) {
        Serial.print("MQTT: connecting to broker...");

        if (client.connect("indoorgrid-esp32")) {
            Serial.println(" connected");
        } else {
            Serial.print(" failed, rc=");
            Serial.println(client.state());
            delay(5000);
        }
    }
}

void MQTTClient::publish(const char* topic, const char* payload) {
    client.publish(topic, payload);
}

void MQTTClient::loop() {
    if (!client.connected()) {
        connectBroker();
    }
    client.loop();
}

bool MQTTClient::isConnected() {
    return client.connected();
}