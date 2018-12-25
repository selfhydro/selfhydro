# Mini Sensors

### Water Distance Sensor

This sensor is used for detecting the height of the water in the tank.

Hardware:
- Wemo D1 mini
- ToF Sensor

Hooking up the hardware:
- SDA -> D2
- SCL -> D1

### MQTT
(Needs to be run from this directory)
1. Run: `docker build -t mosquitto:latest .`
2. Run 
````docker run -d -p 1883:1883 --name mosquitto -v $(pwd)/mosquitto:/var/lib/mosquitto/ --restart on-failure mosquitto````