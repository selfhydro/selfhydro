# SelfHydro - Automated Hydroponics

This is an automated open-source hydroponic system. 

The system currently consists of the following:
- Lights on a timed cycle 
- Air pump turned on every 30 mins
- Water and ambient temperature sensors, which take readings every 3 hours
- Water level of the tank

Both the lights and the air pump are controlled via relays currently but could also use transistor

Currently the water temperature sensor is a DS18B20 sensor - communicating over one wire.
The ambient temperature sensor is a MCP9808 sensor, communicating using I2C. The water level sensor is Ultrasonic HC-SR04 sensor.

The system is designed to work on the Raspberry Pi 3.

Every 3 hours the system will send a json message with the telemetry to google cloud IoT core.
### Wifi-Connect
For first time setup of device, wifi connect is used to setup the wifi network.

Install using: 
``` 
bash <(curl -L https://github.com/resin-io/resin-wifi-connect/raw/master/scripts/raspbian-install.sh)
```


### Setting Selfhydro up as a service

1. ``cp selfhydro.serivce /etc/systemd/system/selfhydro.service``
2. ``sudo systemctl enable selfhydro.service``