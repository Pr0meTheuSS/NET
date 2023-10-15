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
	close            bool
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
		close:            false,
	}
}

func (sm *SpeedMeasurer) MeasureSpeed() {
	go func() {
		alreadyReadBytes := sm.trafficProvider.GetAlreadyReadBytes()
		if alreadyReadBytes < 0 {
			return
		}
		diff := alreadyReadBytes - sm.alreadyReadBytes
		sm.currentBPS = float64(diff) / float64(timeoutInSeconds)
		sm.alreadyReadBytes = alreadyReadBytes

		sm.outputChannel <- fmt.Sprintf("Speed: %f bps. avg: %f, uptime: %f s\r", sm.getCurrentBPS(), sm.getAverageBPS(), sm.getTimeDeltaInSeconds())

		for !sm.isClose() {
			select {
			case <-time.After(time.Second * timeoutInSeconds):
				alreadyReadBytes := sm.trafficProvider.GetAlreadyReadBytes()
				if alreadyReadBytes < 0 {
					return
				}
				diff := alreadyReadBytes - sm.alreadyReadBytes
				sm.currentBPS = float64(diff) / float64(timeoutInSeconds)
				sm.alreadyReadBytes = alreadyReadBytes

				sm.outputChannel <- fmt.Sprintf("Speed: %f bps. avg: %f, uptime: %f s\r", sm.getCurrentBPS(), sm.getAverageBPS(), sm.getTimeDeltaInSeconds())
			}
		}
	}()
}

func (sm *SpeedMeasurer) getCurrentBPS() float64 {
	return sm.currentBPS
}

func (sm *SpeedMeasurer) getAverageBPS() float64 {
	return float64(sm.alreadyReadBytes) / sm.getTimeDeltaInSeconds()
}

func (sm *SpeedMeasurer) getTimeDeltaInSeconds() float64 {
	return time.Now().Sub(sm.initTime).Seconds()
}

func (sm *SpeedMeasurer) isClose() bool {
	return sm.close
}

func (sm *SpeedMeasurer) Close() {
	sm.close = true
}
