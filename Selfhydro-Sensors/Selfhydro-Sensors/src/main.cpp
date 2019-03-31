#include <Arduino.h>
#include "Adafruit_VL53L0X.h"
#include <ESP8266WiFi.h> 
#include <PubSubClient.h> 

#include <Wire.h>
#include <VL53L0X.h>

VL53L0X sensor;
Adafruit_VL53L0X lox = Adafruit_VL53L0X();
#define HIGH_ACCURACY

const char* ssid = "ii52938Dprimary";
const char* wifi_password = "3dcd5fb5";

const char* mqtt_server = "water.local";
const char* mqtt_topic = "/sensors/water_level";
const char* mqtt_username = "";
const char* mqtt_password = "";

const char* clientID = "Water Level Sensor";

WiFiClient wifiClient;
PubSubClient client(mqtt_server, 1883, wifiClient);

void reconnect() {
  while (!client.connected()) {
    Serial.print("Attempting MQTT connection...");
    if (client.connect(clientID)) {
      Serial.println("connected");
    } else {
      Serial.print("failed, rc=");
      Serial.print(client.state());
      Serial.println(" try again in 5 seconds");
      delay(5000);
    }
  }
}

void setup() {
  Serial.begin(115200);

  while (! Serial) {
    delay(1);
  }

  Serial.print("Connecting to ");
  Serial.println(ssid);

  // Connect to the WiFi
  WiFi.begin(ssid, wifi_password);

  while (WiFi.status() != WL_CONNECTED) {
    delay(500);
    Serial.print(".");
  }

  Serial.println("WiFi connected");
  Serial.print("IP address: ");
  Serial.println(WiFi.localIP());

  if (client.connect(clientID, mqtt_username, mqtt_password)) {
    Serial.println("Connected to MQTT Broker!");
  }
  else {
    Serial.println("Connection to MQTT Broker failed...");
  }

  Serial.println("VL53L0X setup");
  Wire.begin();

  sensor.init();
  sensor.setTimeout(500);

  #if defined HIGH_SPEED
    // reduce timing budget to 20 ms (default is about 33 ms)
    sensor.setMeasurementTimingBudget(20000);
  #elif defined HIGH_ACCURACY
    // increase timing budget to 200 ms
    sensor.setMeasurementTimingBudget(200000);
  #endif
}

void loop() {
  if (!client.connected()){
    reconnect();
  }

  Serial.print("Reading a measurement... ");
  double range = sensor.readRangeSingleMillimeters();
  double adjustedRange = range - double(100);
  char cstr[16];
  if (!sensor.timeoutOccurred()) { 
    Serial.print("Distance (mm): "); Serial.println(adjustedRange);
    if (client.publish(mqtt_topic, itoa(adjustedRange, cstr, 10))) {
      Serial.println("Distance measured and message sent");
    } else {
      Serial.println("Message failed to send via mqtt");
      reconnect();
      client.publish(mqtt_topic, itoa(adjustedRange, cstr, 10));
    }
  } else {
      Serial.println(" out of range ");
      client.publish(mqtt_topic, "0");
  }

  delay(2000);
}