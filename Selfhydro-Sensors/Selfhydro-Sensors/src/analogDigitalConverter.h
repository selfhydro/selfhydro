#ifndef ANALOG_DIGITAL_CONVERTER_H
#define ANALOG_DIGITAL_CONVERTER_H

#include <Arduino.h>
#include <Adafruit_MCP3008.h>

class AnalogDigitalConverter {

    #define CS_PIN D8
    #define CLOCK_PIN D5
    #define MOSI_PIN D7
    #define MISO_PIN D6

    public:
        AnalogDigitalConverter() {
            adc.begin(CLOCK_PIN, MOSI_PIN, MISO_PIN, CS_PIN); 
        }

        int GetValue(int channel) {
            return adc.readADC(channel);
        } 

    private:
        Adafruit_MCP3008 adc;
};

#endif