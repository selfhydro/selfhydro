package main

import (
	"time"
	"errors"
	"github.com/stianeikeland/go-rpio"
	"log"
)

const (
	COLLECTING_PERIOD  = 2 * time.Second
	LOGICAL_1_TRESHOLD = 50 * time.Microsecond
)

var (
	ChecksumError    = errors.New("error: checksum error")
	HumidityError    = errors.New("error: humidity range error")
	TemperatureError = errors.New("error: temperature range error")
)

type DHT22 struct {
	pin         int
	temperature float32
	humidity    float32
	readAt      time.Time
	err         error
}

func NewDHT22(pin int) *DHT22 {
	return &DHT22{pin: pin}
}

func (d *DHT22) Temperature() (float32, error) {
	if err := d.read(); err != nil {
		d.err = err
		return 0, err
	}

	return d.temperature, nil
}

func (d *DHT22) Humidity() (float32, error) {
	if err := d.read(); err != nil {
		d.err = err
		return 0, err
	}

	return d.humidity, nil
}

func (d *DHT22) read() error {
	if d.readAt.Add(COLLECTING_PERIOD).After(time.Now()) {
		return d.err
	}

	d.err = nil

	d.readAt = time.Now()

	lengths := make([]time.Duration, 40)
	iterator := 0

	pin := rpio.Pin(d.pin)
	pin.Mode(rpio.Output)

	var (
		humidity    uint16
		temperature uint16
	)

	initSensor(&pin)
	pin.Mode(rpio.Input)

	for {
		d.readData(&pin, lengths, iterator)

		bytes := generateBytes(lengths)

		if err := d.checksum(bytes); err != nil {
			if err != nil {
				return err
			}
		}

		d.calculateHumidity(humidity, bytes)

		// calculate temperature
		temperature |= uint16(bytes[2])
		temperature <<= 8
		temperature |= uint16(bytes[3])

	}

	// check for negative temperature
	if temperature&0x8000 > 0 {
		d.temperature = float32(temperature&0x7FFF) / -10
	} else {
		d.temperature = float32(temperature) / 10
	}

	// datasheet operating range
	if d.temperature < -40 || d.temperature > 80 {
		return TemperatureError
	}

	return nil
}

func generateBytes(lengths []time.Duration) []uint8 {
	bytes := make([]uint8, 5)
	for i := range bytes {
		for j := 0; j < 8; j++ {
			bytes[i] <<= 1
			if lengths[i*8+j] > LOGICAL_1_TRESHOLD {
				bytes[i] |= 0x01
				log.Print("1")
			} else {
				log.Print("0")
			}
		}
	}
	return bytes
}

func initSensor(pin *rpio.Pin) {
	pin.High()
	time.Sleep(250 * time.Millisecond)
	pin.Low()
	time.Sleep(5 * time.Millisecond)
	pin.High()
	time.Sleep(20 * time.Microsecond)
}

func (d *DHT22) readData(pin *rpio.Pin, lengths []time.Duration, iterator int) {
	for {

		duration := d.timePulse(pin, rpio.High)

		lengths[iterator] = duration
		iterator++
		if iterator >= 40 {
			break
		}
	}

	for dur := range lengths {
		log.Printf("%v", lengths[dur])
	}
}

func (d *DHT22) timePulse(pin *rpio.Pin, state rpio.State) (time.Duration) {

	aroundState := rpio.Low
	if state == rpio.Low {
		aroundState = rpio.High
	}

	for {
		v := pin.Read()

		if v == aroundState {
			break
		}
	}

	for {
		v := pin.Read()
		log.Printf("State: %s",v)
		if v == state {
			break
		}
	}

	startTime := time.Now()

	for {
		v := pin.Read()
		log.Printf("State: %s",v)
		if v == aroundState {
			break
		}
	}

	return time.Since(startTime)
}

func (d *DHT22) calculateHumidity(humidity uint16, bytes []uint8) error {
	humidity |= uint16(bytes[0])
	humidity <<= 8
	humidity |= uint16(bytes[1])

	if humidity < 0 || humidity > 1000 {
		return HumidityError
	}

	d.humidity = float32(humidity) / 10
	return nil
}

func (d *DHT22) checksum(bytes []uint8) error {
	var sum uint8

	for i := 0; i < 4; i++ {
		sum += bytes[i]
	}

	if sum != bytes[4] {
		return ChecksumError
	}

	return nil
}
