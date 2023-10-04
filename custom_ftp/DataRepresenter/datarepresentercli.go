package datarepresenter

import (
	"fmt"
	"os"
	"os/exec"
	"sync"
)

type DataRepresenter interface {
	Represent()
	Register() chan string
	Unregister(chan string)
}

type DataRepresenterCli struct {
	actualData map[chan string]string
	mtx        sync.Mutex
}

func NewDataRepresenterCli() *DataRepresenterCli {
	ret := &DataRepresenterCli{
		actualData: map[chan string]string{},
	}
	ret.Represent()
	return ret
}

var clearCommand = "clear"

func clearScreen() {
	cmd := exec.Command(clearCommand)
	cmd.Stdout = os.Stdout
	cmd.Run()
}

func (drcli *DataRepresenterCli) Represent() {
	go func() {
		for {
			drcli.mtx.Lock()
			for ch := range drcli.actualData {
				select {
				case data, ok := <-ch:
					if !ok {
						break
					}

					drcli.actualData[ch] = data
					for _, v := range drcli.actualData {
						fmt.Println(v)
					}
				default:
					break
				}
			}
			drcli.mtx.Unlock()
		}
	}()
}

func (drcli *DataRepresenterCli) Register() chan string {
	fmt.Println("Somebody registered")
	ch := make(chan string)

	drcli.mtx.Lock()
	drcli.actualData[ch] = ""
	drcli.mtx.Unlock()

	return ch
}

func (drcli *DataRepresenterCli) Unregister(ch chan string) {
	drcli.mtx.Lock()
	delete(drcli.actualData, ch)
	drcli.mtx.Unlock()
}
