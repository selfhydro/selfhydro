#include <Arduino.h>
#include <ESP8266WiFi.h> 
#include <PubSubClient.h> 

#include <Wire.h>
#include "Adafruit_Si7021.h"
#include <ArduinoJson.h>

Adafruit_Si7021 sensor = Adafruit_Si7021();

const char* ssid = "ii52938Dprimary";
const char* wifi_password = "3dcd5fb5";

const char* mqtt_server = "water.local";
const char* mqtt_ambient_temperature_topic = "/state/ambient_temperature";
const char* mqtt_ambient_humidity_topic = "/state/ambient_humidity";
const char* mqtt_username = "";
const char* mqtt_password = "";

const char* clientID = "Ambient Temperature and Humidity Sensor";

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
  
  Serial.println("Si7021 test!");
  
  if (!sensor.begin()) {
    Serial.println("Did not find Si7021 sensor!");
    while (true);
  }

  Serial.print("Found model ");
  switch(sensor.getModel()) {
    case SI_Engineering_Samples:
      Serial.print("SI engineering samples"); break;
    case SI_7013:
      Serial.print("Si7013"); break;
    case SI_7020:
      Serial.print("Si7020"); break;
    case SI_7021:
      Serial.print("Si7021"); break;
    case SI_UNKNOWN:
    default:
      Serial.print("Unknown");
  }
  Serial.print(" Rev(");
  Serial.print(sensor.getRevision());
  Serial.print(")");
  Serial.print(" Serial #"); Serial.print(sensor.sernum_a, HEX); Serial.println(sensor.sernum_b, HEX); 
}

void loop() {
  if (!client.connected()){
    reconnect();
  }
    
  float humidity = sensor.readHumidity();
  float temperature = sensor.readTemperature();

  Serial.print("Humidity:    ");
  Serial.print(humidity, 2);
  Serial.print("\tTemperature: ");
  Serial.println(temperature, 2);

  const int capacity = JSON_OBJECT_SIZE(3);
  StaticJsonDocument<capacity> ambientTemperatureJson;
  StaticJsonDocument<capacity> ambientHumidityJson;
  ambientHumidityJson["humidity"] = humidity;
  ambientTemperatureJson["temperature"] = temperature;

  char ambientTemperatureCStr[128];
  char ambientHumidityCStr[128];

  serializeJson(ambientTemperatureJson, ambientTemperatureCStr);
  serializeJson(ambientHumidityJson, ambientHumidityCStr);

  if (client.publish(mqtt_ambient_temperature_topic, ambientTemperatureCStr)) {
    Serial.println("Temperature measured and message sent");
  } else {
    Serial.println("Message failed to send via mqtt");
    reconnect();
    client.publish(mqtt_ambient_temperature_topic, ambientTemperatureCStr);
  }

  if (client.publish(mqtt_ambient_humidity_topic, ambientHumidityCStr)) {
    Serial.println("Humidity mesured and message sent");
  } else {
    Serial.println("Message failed to send via mqtt");
    reconnect();
    client.publish(mqtt_ambient_humidity_topic, ambientHumidityCStr);
  }

  delay(5000);
}