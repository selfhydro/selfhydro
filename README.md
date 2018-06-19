# SelfHydro

This is an automated hydroponic system. 

It currently consists of an automated air pump, grow LEDs and temp. sensors for both water temperature and ambient temperature.

It all runs on a Raspberry Pi.

Every 4 hours google cloud is updated via a MQTT bridge on the state of the device.

### Wifi-Connect
For first time setup of device, wifi connect is used to setup the wifi network.

Install using: 
``` 
bash <(curl -L https://github.com/resin-io/resin-wifi-connect/raw/master/scripts/raspbian-install.sh)
```

[]


### Setting Selfhydro up as a service

1. ``cp selfhydro.serivce /etc/systemd/system/selfhydro.service``
2. ``sudo systemctl enable selfhydro.service``