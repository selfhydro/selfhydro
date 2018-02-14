package main

import (
	"testing"
	"time"
)

func TestShouldTurnOnPump(t *testing.T){

}

func TestShouldSayIfTimeIsBetweenTwoTimes(t *testing.T){
	startingTime := time.Now().Local().Add(-time.Hour)
	endTime := time.Now().Local().Add(time.Hour)
	if betweenTime(startingTime, endTime) {
		t.Error("Error: Current time should be between start and end time")
	}
}