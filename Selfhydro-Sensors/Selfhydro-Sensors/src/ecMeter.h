#include <Arduino.h>
#include <OneWire.h>

#include "analogDigitalConverter.h"

#define StartConvert 0
#define ReadTemperature 1
#define SENSOR_CHANNEL 0

#define NumberOfReadings 20
#define ECSensorPin 1

class ECMeter
{
    public:
        ECMeter() {};
        void Setup();
        float GetReading(float temperature);

    private:
        AnalogDigitalConverter adc;
        unsigned int AnalogSampleInterval=25,printInterval=700,tempSampleInterval=850; 
        unsigned int readings[NumberOfReadings];     
        byte index = 0;                  
        unsigned long AnalogValueTotal = 0;                 
        unsigned int AnalogAverage = 0,averageVoltage=0;               
        unsigned long AnalogSampleTime,printTime,tempSampleTime;
        float temperature,ECcurrent;

};