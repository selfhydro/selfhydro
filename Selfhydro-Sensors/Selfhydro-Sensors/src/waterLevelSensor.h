#include <Arduino.h>
#include <Wire.h>
#include <SparkFun_VL6180X.h>

#define VL6180X_ADDRESS 0x29

class WaterLevelSensor
{
    public:
        WaterLevelSensor() {};
        void Setup() {
            sensor.getIdentification(&identification); 
            printIdentification(&identification); 
            
            if(sensor.VL6180xInit() != 0){
                Serial.println("FAILED TO INITALIZE"); 
            }; 
            sensor.VL6180xDefautSettings(); 
            delay(1000);
        };

        float GetReading(){
            float level = sensor.getDistance();
            Serial.print("Distance measured (mm) = ");
            Serial.println( sensor.getDistance() ); 

            return level;
        };

    private:
        VL6180xIdentification identification;
        VL6180x sensor = VL6180x(VL6180X_ADDRESS);

        void printIdentification(struct VL6180xIdentification *temp){
            Serial.print("Model ID = ");
            Serial.println(temp->idModel);

            Serial.print("Model Rev = ");
            Serial.print(temp->idModelRevMajor);
            Serial.print(".");
            Serial.println(temp->idModelRevMinor);

            Serial.print("Module Rev = ");
            Serial.print(temp->idModuleRevMajor);
            Serial.print(".");
            Serial.println(temp->idModuleRevMinor);  

            Serial.print("Manufacture Date = ");
            Serial.print((temp->idDate >> 3) & 0x001F);
            Serial.print("/");
            Serial.print((temp->idDate >> 8) & 0x000F);
            Serial.print("/1");
            Serial.print((temp->idDate >> 12) & 0x000F);
            Serial.print(" Phase: ");
            Serial.println(temp->idDate & 0x0007);

            Serial.print("Manufacture Time (s)= ");
            Serial.println(temp->idTime * 2);
            Serial.println();
            Serial.println();
        }
};