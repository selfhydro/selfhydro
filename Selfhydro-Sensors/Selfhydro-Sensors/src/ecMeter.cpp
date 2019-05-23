#include "ecMeter.h"

void ECMeter::Setup() {
    adc = AnalogDigitalConverter();
}

float ECMeter::GetReading(float temperature) {
  if(millis()-AnalogSampleTime>=AnalogSampleInterval)
  {
    AnalogSampleTime=millis();
    AnalogValueTotal = AnalogValueTotal - readings[index];
    readings[index] = adc.GetValue(SENSOR_CHANNEL);
    AnalogValueTotal = AnalogValueTotal + readings[index];
    index = index + 1;
    if (index >= NumberOfReadings)
    index = 0;
    AnalogAverage = AnalogValueTotal / NumberOfReadings;
  }

  if(millis()-printTime>=printInterval)
  {
    printTime=millis();
    averageVoltage=AnalogAverage*(float)5000/1024;
    
    Serial.print("Analog value:");
    Serial.print(AnalogAverage);   
    Serial.print("    Voltage:");
    Serial.print(averageVoltage);  
    Serial.print("mV    ");
    Serial.print("temp:");
    Serial.print(temperature);    
    Serial.print("^C     EC:");
    
    float TempCoefficient=1.0+0.0185*(temperature-25.0);    
    float CoefficientVolatge=(float)averageVoltage/TempCoefficient;
    
    if(CoefficientVolatge<150)Serial.println("No solution!");   
    else if(CoefficientVolatge>3300)Serial.println("Out of the range!");  
    else
    {
        if(CoefficientVolatge<=448)ECcurrent=6.84*CoefficientVolatge-64.32;   
        else if(CoefficientVolatge<=1457)ECcurrent=6.98*CoefficientVolatge-127;  
        else ECcurrent=5.3*CoefficientVolatge+2278;                           
    
        ECcurrent/=1000;    
        Serial.print(ECcurrent,2);  
        Serial.println("ms/cm");
    }
    }
    return 0.0;
}