#include <Arduino.h>
#include "Adafruit_VL53L0X.h"
#include <ESP8266WiFi.h> 
#include <PubSubClient.h> 

Adafruit_VL53L0X lox = Adafruit_VL53L0X();

const char* ssid = "ii52938Dprimary";
const char* wifi_password = "3dcd5fb5";

const char* mqtt_server = "10.1.1.3";
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

  Serial.println("Adafruit VL53L0X test");
  if (!lox.begin()) {
    Serial.println(F("Failed to boot VL53L0X"));
    while(1);
  }
  Serial.println(F("VL53L0X API Simple Ranging example\n\n"));
}

void loop() {
  VL53L0X_RangingMeasurementData_t measure;

  if (!client.connected()){
    reconnect();
  }

  Serial.print("Reading a measurement... ");
  lox.rangingTest(&measure, false);
  String mqttMessage = String("\"Water Level\":" + measure.RangeMilliMeter);
  char cstr[16];
  if (measure.RangeStatus != 4) { 
    Serial.print("Distance (mm): "); Serial.println(measure.RangeMilliMeter);
    if (client.publish(mqtt_topic, itoa(measure.RangeMilliMeter, cstr, 10))) {
      Serial.println("Distance measured and message sent");
    } else {
      Serial.println("Message failed to send via mqtt");
      reconnect();
      client.publish(mqtt_topic, itoa(measure.RangeMilliMeter, cstr, 10));
    }
  } else {
  Serial.println(" out of range ");
  }

  delay(5000);
}