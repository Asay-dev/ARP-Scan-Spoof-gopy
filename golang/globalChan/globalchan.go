package globalChan

import (
	"sync"
)

var chanInstance chan int
var chanString chan string

var chanOnceManager sync.Once

func GetABGlobalChanInt() chan int {

	chanOnceManager.Do(func() {

		chanInstance = make(chan int, 100)

	})

	return chanInstance

}
func GetABGlobalChanString() chan string {

	chanOnceManager.Do(func() {

		chanString = make(chan string)

	})

	return chanString

}
