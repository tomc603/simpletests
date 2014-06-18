package main

import (
	"flag"
	"fmt"
	"sync"
	"time"
)

type mutexdata struct {
	sync.Mutex
	count      uint
	delaytotal time.Duration
	min        time.Duration
	max        time.Duration
}

type roundtrip struct {
	sent     time.Time
	received time.Time
}

func tstest() {
	/* Test speed of timestamp assignment by assigning
	 * as many timestamps as possible in 10 seconds */
	var (
		loopCount uint64
		timeDelta time.Duration
		timeMax   time.Duration
		timeMin   time.Duration
		timeTotal time.Duration
		quit      bool
	)

	data := roundtrip{}
	timeStart := time.Now()
	stopTimer := time.NewTimer(30 * time.Second)

	fmt.Print("Starting time stamp delay test\n")
	for {
		select {
		case _ <- stopTimer.C:
			stopTimer.Stop()
			quit = true
		}

		data.sent = time.Now()
		data.received = time.Now()

		timeDelta = data.received.Sub(data.sent)
		timeTotal += timeDelta
		switch {
		case timeMin == 0 || timeDelta < timeMin:
			timeMin = timeDelta
		case timeDelta > timeMax:
			timeMax = timeDelta
		}

		loopCount++
		if quit {
			break
		}
	}
	fmt.Printf("Count: %d (%f/s)\nMin/Max/Avg RTT: %dns/%dns/%dns\n\n",
		loopCount, (float64(loopCount) / time.Since(st).Seconds()),
		timeMin.Nanoseconds(), timeMax.Nanoseconds(), uint64(timeTotal.Nanoseconds())/loopCount)
}

func main() {
	v := *flag.Bool("v", false, "set v flag")
	flag.Parse()

	fmt.Printf("v: %v\n", v)
	for _, fl := range flag.Args() {
		fmt.Printf("Flag: %v\n", fl)
	}

	/* Run test functions... */
	// Time delta between writing time.Now() to a struct
	tstest()
	fmt.Print("\n\n")

	// Reflect a struct back to sender
	//chancommtest()
	//fmt.Print("\n\n")

	// Test writing to a struct, locking with a mutex
	mutextest()
	fmt.Print("\n\n")

	// Test writing to a struct via channel
	chanpubsub()
	fmt.Print("\n\n")
}
