package speedmeasure

import (
	"fmt"
	"time"
)

type SpeedMeasurer struct {
	outputChannel    chan string
	currentBPS       float64
	alreadyReadBytes int64
	initTime         time.Time
	trafficProvider  Traffic
}

type Traffic interface {
	GetAlreadyReadBytes() int64
}

const (
	timeoutInSeconds = 3
)

type ConnectionTrafficInfo struct {
	Speed    float64
	AvgSpeed float64
	Uptime   int64
}

func NewSpeedMeasurer(traffic Traffic, ch chan string) *SpeedMeasurer {
	return &SpeedMeasurer{
		outputChannel:    ch,
		currentBPS:       0,
		alreadyReadBytes: 0,
		initTime:         time.Now(),
		trafficProvider:  traffic,
	}
}

func (sm *SpeedMeasurer) MeasureSpeed() {
	go func() {
		for {
			select {
			case <-time.After(time.Second * timeoutInSeconds):
				alreadyReadBytes := sm.trafficProvider.GetAlreadyReadBytes()
				if alreadyReadBytes < 0 {
					return
				}
				diff := alreadyReadBytes - sm.alreadyReadBytes
				sm.currentBPS = float64(diff) / float64(timeoutInSeconds)
				sm.alreadyReadBytes = alreadyReadBytes

				sm.outputChannel <- fmt.Sprintf("Speed: %f bps. avg: %f, uptime: %f s\r", sm.GetCurrentBPS(), sm.GetAverageBPS(), sm.GetTimeDeltaInSeconds())
			}
		}
	}()
}

func (sm *SpeedMeasurer) GetCurrentBPS() float64 {
	return sm.currentBPS
}

func (sm *SpeedMeasurer) GetAverageBPS() float64 {
	return float64(sm.alreadyReadBytes) / sm.GetTimeDeltaInSeconds()
}

func (sm *SpeedMeasurer) GetTimeDeltaInSeconds() float64 {
	return time.Now().Sub(sm.initTime).Seconds()
}
