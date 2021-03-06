#include <Arduino.h>
#include <ESP8266WiFi.h> 
#include <PubSubClient.h> 
#include <Wire.h>
#include "Adafruit_Si7021.h"
#include <ArduinoJson.h>
#include <OneWire.h>

#include "phSensor.h"
#include "waterLevelSensor.h"
#include "ecMeter.h"

OneWire  ds(D4);
Adafruit_Si7021 sensor = Adafruit_Si7021();
PHSensor phSensor = PHSensor();
WaterLevelSensor waterLevelSensor = WaterLevelSensor();
ECMeter ecMeter = ECMeter();

const char* ssid = "NETGEAR59";
const char* wifi_password = "calmink680";

const char* mqtt_server = "selfhydro-base.local";
const char* mqtt_ambient_temperature_topic = "/state/ambient_temperature";
const char* mqtt_ambient_humidity_topic = "/state/ambient_humidity";
const char* mqtt_water_temperature_topic = "/state/water_temperature";
const char* mqtt_battery_voltage_topic = "/state/battery_voltage";
const char* mqtt_water_level_topic = "/state/water_level";
const char* mqtt_water_ec_topic = "/state/water_ec";
const char* mqtt_pH_topic = "/state/pH";
const char* mqtt_username = "";
const char* mqtt_password = "";

#define durationDeepSleep  900 // 15min 

#ifdef AMBIENT_TEMP
const char* clientID = "Ambient Humidity and Temperature and Water Temp Sensor";
#elif EC_METER
const char* clientID = "EC Meter";
#endif

WiFiClient wifiClient;
PubSubClient client(mqtt_server, 1883, wifiClient);

float waterTemperature = 0.0;

float getBatteryVoltage() {
  int rawVolatge = analogRead(A0);
  Serial.println(rawVolatge);
  float voltage = rawVolatge / 1023.0;
  Serial.println(voltage);

  voltage = voltage * 45/13;
  Serial.println(voltage);

  return voltage;
}

void waterTemperatureCallback(char* topic, byte* payload, unsigned int length) {
  Serial.print("Message arrived [");
  Serial.print(topic);
  Serial.print("] ");
  for (int i = 0; i < length; i++) {
    Serial.print((char)payload[i]);
  }
  Serial.println();

  const int capacity = JSON_OBJECT_SIZE(3);
  StaticJsonDocument<capacity> waterTemperatureJSON;
  deserializeJson(waterTemperatureJSON, payload);
  waterTemperature = waterTemperatureJSON["temperature"].as<float>();
}

float getWaterTemp() {
    byte i;
    byte present = 0;
    byte type_s;
    byte data[12];
    byte addr[8];
    float celsius = 0.0;
  
    if ( !ds.search(addr)) 
    {
      ds.reset_search();
      delay(250);
      return celsius;
    }
  
  
    if (OneWire::crc8(addr, 7) != addr[7]) 
    {
        Serial.println("CRC is not valid!");
        return celsius;
    }
    Serial.println();
  
    // the first ROM byte indicates which chip
    switch (addr[0]) 
    {
      case 0x10:
        type_s = 1;
        break;
      case 0x28:
        type_s = 0;
        break;
      case 0x22:
        type_s = 0;
        break;
      default:
        Serial.println("Device is not a DS18x20 family device.");
        return celsius;
    } 
  
    ds.reset();
    ds.select(addr);
    ds.write(0x44, 1);        // start conversion, with parasite power on at the end  
    delay(1000);
    present = ds.reset();
    ds.select(addr);    
    ds.write(0xBE);         // Read Scratchpad
  
    for ( i = 0; i < 9; i++) 
    {           
      data[i] = ds.read();
    }
  
    // Convert the data to actual temperature
    int16_t rawTemperature = (data[1] << 8) | data[0];
    if (type_s) {
      rawTemperature = rawTemperature << 3; // 9 bit resolution default
      if (data[7] == 0x10) 
      {
        rawTemperature = (rawTemperature & 0xFFF0) + 12 - data[6];
      }
    } 
    else 
    {
      byte cfg = (data[4] & 0x60);
      if (cfg == 0x00) rawTemperature = rawTemperature & ~7;  // 9 bit resolution, 93.75 ms
      else if (cfg == 0x20) rawTemperature = rawTemperature & ~3; // 10 bit res, 187.5 ms
      else if (cfg == 0x40) rawTemperature = rawTemperature & ~1; // 11 bit res, 375 ms
  
    }
    float waterTemperatureCelcius = (float)rawTemperature / 16.0;
    return waterTemperatureCelcius;
}

void reconnect() {
  while (!client.connected()) {
    Serial.print("Attempting MQTT connection...");
    if (client.connect(clientID)) {
      Serial.println("connected");
      client.subscribe(mqtt_water_temperature_topic);
    } else {
      Serial.print("failed, rc=");
      Serial.print(client.state());
      Serial.println(" try again in 5 seconds");
      delay(5000);
    }
  }
}

void setupAmbientTempAndHumidity() {
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


void setup() {
  Serial.begin(115200);

  while (! Serial) {
    delay(1);
  }
  pinMode(A0, INPUT);
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

  #ifdef AMBIENT_TEMP
    setupAmbientTempAndHumidity();
  #endif

  #ifdef WATER_LEVEL
    waterLevelSensor.Setup();
  #endif

  #ifdef PH_SENSOR
    phSensor.Setup();
    Serial.println("ph sensor setup");
  #endif

  #ifdef EC_METER
    ecMeter.Setup();
    client.setCallback(waterTemperatureCallback);
  #endif
 
}

void loop() {
  if (!client.connected()){
    reconnect();
  }

  #ifdef EC_METER
    float ecLevel;
    Serial.println(waterTemperature);
    while(waterTemperature == 0.0){
      delay(500);
      if (!client.connected()){
        reconnect();
      } else {
        client.loop();
      }
    }
    Serial.println(waterTemperature);

    ecLevel = ecMeter.GetReading(waterTemperature);
    const int capacity = JSON_OBJECT_SIZE(3);
    StaticJsonDocument<capacity> ecLevelJSON;
    ecLevelJSON["ecLevel"] = ecLevel;
    char ecLevelCStr[128];
    serializeJson(ecLevelJSON, ecLevelCStr);
    if (client.publish(mqtt_water_ec_topic, ecLevelCStr)) {
      Serial.println("EC level measured and message sent");
    } else {
      Serial.println("Message failed to send via mqtt");
      reconnect();
      client.publish(mqtt_water_ec_topic, ecLevelCStr);
    }
    delay(2000000);
  #endif

  #ifdef WATER_LEVEL
    float waterLevel = waterLevelSensor.GetReading();
    Serial.print("water level:"); Serial.print(waterLevel, 2); Serial.println("mm");

    const int capacity = JSON_OBJECT_SIZE(3);
    StaticJsonDocument<capacity> waterLevelJSON;
    waterLevelJSON["waterLevel"] = waterLevel;
    char waterLevelCStr[128];
    serializeJson(waterLevelJSON, waterLevelCStr);
    if (client.publish(mqtt_water_level_topic, waterLevelCStr)) {
      Serial.println("Water level measured and message sent");
    } else {
      Serial.println("Message failed to send via mqtt");
      reconnect();
      client.publish(mqtt_water_level_topic, waterLevelCStr);
    }
  #endif

  #ifdef WATER_TEMP

  #endif

  #ifdef AMBIENT_TEMP
    float humidity = sensor.readHumidity();
    float temperature = sensor.readTemperature();
    float waterTemperature = getWaterTemp();
    float batteryVolate = getBatteryVoltage();

    Serial.print("Humidity:    ");
    Serial.print(humidity, 2);
    Serial.print("\tTemperature: ");
    Serial.println(temperature, 2);
    Serial.print("\tWater Temperature: ");
    Serial.println(waterTemperature, 2);
    Serial.print("\tBattery Voltage: ");
    Serial.println(batteryVolate, 2);

    const int capacity = JSON_OBJECT_SIZE(3);
    StaticJsonDocument<capacity> ambientTemperatureJson;
    StaticJsonDocument<capacity> ambientHumidityJson;
    StaticJsonDocument<capacity> waterTemperatureJson;
    StaticJsonDocument<capacity> batteryVoltageJson;

    ambientHumidityJson["humidity"] = humidity;
    ambientTemperatureJson["temperature"] = temperature;
    waterTemperatureJson["temperature"] = waterTemperature;
    batteryVoltageJson["voltage"] = batteryVolate;
    char ambientTemperatureCStr[128];
    char ambientHumidityCStr[128];
    char waterTemperatureCStr[128];
    char batteryVoltageCStr[128];
    serializeJson(ambientTemperatureJson, ambientTemperatureCStr);
    serializeJson(ambientHumidityJson, ambientHumidityCStr);
    serializeJson(waterTemperatureJson, waterTemperatureCStr);
    serializeJson(batteryVoltageJson, batteryVoltageCStr);

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

    if (client.publish(mqtt_water_temperature_topic, waterTemperatureCStr)) {
      Serial.println("Water temp mesured and message sent");
    } else {
      Serial.println("Message failed to send via mqtt");
      reconnect();
      client.publish(mqtt_water_temperature_topic, waterTemperatureCStr);
    }

    if (client.publish(mqtt_battery_voltage_topic, batteryVoltageCStr)) {
      Serial.println("Battery voltage mesured and message sent");
    } else {
      Serial.println("Message failed to send via mqtt");
      reconnect();
      client.publish(mqtt_battery_voltage_topic, batteryVoltageCStr);
    }
  #endif

  #ifdef PH_SENSOR
    float ph = phSensor.GetReading();
    
    const int capacity = JSON_OBJECT_SIZE(3);
    StaticJsonDocument<capacity> phJson;
    phJson["pH"] = ph;
    char phCStr[128];
    serializeJson(phJson, phCStr);
    if (client.publish(mqtt_pH_topic, phCStr)) {
      Serial.println("pH measured and message sent");
    } else {
      Serial.println("Message failed to send via mqtt");
      reconnect();
      client.publish(mqtt_pH_topic, phCStr);
    }
  #endif
    delay(40);
    Serial.println("INFO: Closing the MQTT connection");
    client.disconnect();

    Serial.println("INFO: Closing the Wifi connection");
    WiFi.disconnect();

     while (client.connected() || (WiFi.status() == WL_CONNECTED))
    {
      Serial.println("Waiting for shutdown before sleeping");
      delay(10);
    }
    ESP.deepSleep(durationDeepSleep * 1000000);
}