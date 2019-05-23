#include <Arduino.h>

#include "phSensor.h"
#include "analogDigitalConverter.h"

#define SENSOR_CHANNEL 0 
#define Offset 32.97 //deviation compensate

void PHSensor::Setup(){
    adc = AnalogDigitalConverter();
}

float PHSensor::GetReading(){
    int buf[10]; 
    for(int i=0;i<10;i++) 
    {
        buf[i]=adc.GetValue(SENSOR_CHANNEL);
        delay(10);
    }
    for(int i=0;i<9;i++) 
    {
        for(int j=i+1;j<10;j++)
        {
        if(buf[i]>buf[j])
        {
            int temp=buf[i];
            buf[i]=buf[j];
            buf[j]=temp;
        }
        }
    }
    avgValue=0;
    for(int i=2;i<8;i++) 
        avgValue+=buf[i];

    float V = (float)avgValue*5.0/1024/6; 

    float phValue=-6.65*V+Offset; 
    Serial.print(" pH:");
    Serial.print(phValue,2);
    Serial.println(" ");
    return phValue;
}
