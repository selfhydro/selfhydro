#include <Arduino.h>
#include "Adafruit_VL53L0X.h"
#include <ESP8266WiFi.h> 
#include <PubSubClient.h> 

Adafruit_VL53L0X lox = Adafruit_VL53L0X();

const char* ssid = "ii52938Dprimary";
const char* wifi_password = "3dcd5fb5";

const char* mqtt_server = "";
const char* mqtt_topic = "/sensors/reading";
const char* mqtt_username = "";
const char* mqtt_password = "";

const char* clientID = "Water Level Sensor";

WiFiClient wifiClient;
PubSubClient client(mqtt_server, 1883, wifiClient);

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

  // if (client.connect(clientID, mqtt_username, mqtt_password)) {
  //   Serial.println("Connected to MQTT Broker!");
  // }
  // else {
  //   Serial.println("Connection to MQTT Broker failed...");
  // }

  Serial.println("Adafruit VL53L0X test");
  if (!lox.begin()) {
    Serial.println(F("Failed to boot VL53L0X"));
    while(1);
  }
  Serial.println(F("VL53L0X API Simple Ranging example\n\n"));
}

void loop() {
  VL53L0X_RangingMeasurementData_t measure;

  Serial.print("Reading a measurement... ");
  lox.rangingTest(&measure, false);

  if (measure.RangeStatus != 4) { 
    Serial.print("Distance (mm): "); Serial.println(measure.RangeMilliMeter);
    // if (client.publish(mqtt_topic, "Water Level: ")) {
    //   Serial.println("Button pushed and message sent!");
    // }
  } else {
  Serial.println(" out of range ");
  }

  delay(100);
}