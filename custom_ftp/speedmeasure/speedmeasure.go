package speedmeasure

import (
	"fmt"
	"time"
)

type SpeedMeasurer struct {
	currentBPS       float64
	alreadyReadBytes uint64
}

type Traffic interface {
	GetAlreadyReadBytes() uint64
}

const (
	timeoutInSeconds = 3
)

func (sm *SpeedMeasurer) MeasureSpeed(traffic Traffic) {
	diff := traffic.GetAlreadyReadBytes() - sm.alreadyReadBytes
	sm.currentBPS = float64(diff) / float64(timeoutInSeconds)
	sm.alreadyReadBytes += diff
	fmt.Println("Current Speed (bps):", sm.GetCurrentBPS())

	go func() {
		for {
			select {
			case <-time.After(time.Second * 1):
				diff := traffic.GetAlreadyReadBytes() - sm.alreadyReadBytes
				sm.currentBPS = float64(diff) / float64(timeoutInSeconds)
				sm.alreadyReadBytes += diff
				fmt.Printf("Current Speed (bps): %f\r", sm.GetCurrentBPS())
			}
		}
	}()
}

func (sm *SpeedMeasurer) GetCurrentBPS() float64 {
	return sm.currentBPS
}
