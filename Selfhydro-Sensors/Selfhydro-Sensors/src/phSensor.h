#include <Arduino.h>

#include "analogDigitalConverter.h"

class PHSensor
{
    public:
        PHSensor() {};
        void Setup();
        float GetReading();

    private:
        unsigned long int avgValue;
        AnalogDigitalConverter adc;

};