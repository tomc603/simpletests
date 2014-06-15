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
		t1        time.Time
		t2        time.Time
		total     time.Duration
		min       time.Duration
		max       time.Duration
	)

	st := time.Now()

	fmt.Print("Starting time stamp delay test\n")
	for time.Since(st).Seconds() < 10 {
		t1 = time.Now()
		t2 = time.Now()

		delay := t2.Sub(t1)
		total += delay
		switch {
		case min == 0:
			//fmt.Printf("Setting minimum: %d\n", delay.Nanoseconds())
			min = delay
		case delay > max:
			//fmt.Printf("Setting maximum: %d\n", delay.Nanoseconds())
			max = delay
		case delay < min:
			//fmt.Printf("Setting minimum: %d\n", delay.Nanoseconds())
			min = delay
		}

		loopCount++
	}
	fmt.Printf("Count: %d (%f/s)\nMin/Max/Avg RTT: %dns/%dns/%dns\n\n",
		loopCount, (float64(loopCount) / time.Since(st).Seconds()),
		min.Nanoseconds(), max.Nanoseconds(), uint64(total.Nanoseconds())/loopCount)
}

func chancommtest() {
	/* Test communication speed of channels by sending a struct
	 * with sent and received timestamps over a channel for 10 seconds.
	 * Structure to pass over channels */
	var (
		loopCount uint64
		total     time.Duration
		min       time.Duration
		max       time.Duration
	)

	c := make(chan roundtrip)
	st := time.Now()

	fmt.Print("Starting Channel communication test\n")
	/* Start goroutine to listen for chan activity */
	go receiver(c)

	for time.Since(st).Seconds() < 10 {
		/* Send an item over the channel, increment counter */
		c <- roundtrip{sent: time.Now()}
		r := <-c
		delay := time.Since(r.sent)
		total += delay

		switch {
		case min == 0:
			//fmt.Printf("Setting minimum RTT: %d\n", delay.Nanoseconds())
			min = delay
		case delay > max:
			//fmt.Printf("Setting maximum RTT: %d\n", delay.Nanoseconds())
			max = delay
		case delay < min:
			//fmt.Printf("Setting minimum RTT: %d\n", delay.Nanoseconds())
			min = delay
		}

		loopCount++
	}
	close(c)
	fmt.Printf("Count: %d (%f/s)\nMin/Max/Avg RTT: %dns/%dns/%dns\n\n",
		loopCount, (float64(loopCount) / time.Since(st).Seconds()),
		min.Nanoseconds(), max.Nanoseconds(), uint64(total.Nanoseconds())/loopCount)
}

func receiver(c chan roundtrip) {
	for {
		select {
		case r := <-c:
			r.received = time.Now()
			c <- r
		}
	}
}

func mutextest() {
	/* This test spawns a number of goroutines, each writing to
	 * a struct locked by a mutex lock. Each routine counts the number
	 * of writes it performs to the strutc, and how much time it has
	 * spent waiting on a lock.
	 *
	 * At the end of the test, each routine reports min/max/avg wait
	 * as well as total events.
	 */

	var (
		done bool
		r    mutexdata
		w    sync.WaitGroup
		wc   uint
	)

	fmt.Print("Starting mutex lock test\n")
	t := time.NewTimer(30 * time.Second)
	c := time.NewTicker(1 * time.Second)
	s := time.Now()
	for {
		select {
		case <-t.C:
			// Time has elapsed. Stop the routine.
			t.Stop()
			c.Stop()
			fmt.Print("Time elapsed...\n")
			done = true
		case <-c.C:
			wc++
			fmt.Printf("Adding mutexworker %d...\n", wc)
			w.Add(1)
			go mutexworker(&r, &w, wc)
		}
		if done {
			break
		}
	}
	fmt.Print("Waiting for mutexworkers to complete...\n")
	w.Wait()
	fmt.Print("mutexworkers done...\n")
	fmt.Printf("Count: %d (%f/s)\nMin/Max/Avg RTT: %dns/%dns/%dns\n\n",
		r.count, (float64(r.count) / time.Since(s).Seconds()),
		r.min.Nanoseconds(), r.max.Nanoseconds(), r.delaytotal.Nanoseconds()/int64(r.count))
}

func mutexworker(r *mutexdata, w *sync.WaitGroup, i uint) {
	/* This is a goroutine that accesses roundtrip{} */
	var (
		delay time.Duration
		done  bool
		s     time.Time
	)

	t := time.NewTimer(30 * time.Second)
	for {
		select {
		case <-t.C:
			// Time has elapsed. Stop the routine.
			t.Stop()
			done = true
		default:
			s = time.Now()
			r.Lock()
			delay = time.Since(s)
			r.delaytotal += delay

			switch {
			case r.min == 0:
				fmt.Printf("(%d) Setting mutex min: %d\n", i, delay.Nanoseconds())
				r.min = delay
			case delay > r.max:
				fmt.Printf("(%d) Setting mutex max: %d\n", i, delay.Nanoseconds())
				r.max = delay
			case delay < r.min:
				fmt.Printf("(%d) Setting mutex min: %d\n", i, delay.Nanoseconds())
				r.min = delay
			}

			r.count += 1
			r.Unlock()
		}
		if done {
			break
		}
	}
	fmt.Printf("mutexworker(%d): time elapsed\n", i)
	w.Done()
}

func chanwaittest() {
	/* This test spawns a listener goroutine as well as a number of
	 * writer goroutines. Each writer sends data over a channel to the
	 * listener routine, which aggregates that data into a struct.
	 * Each routine counts the number of writes it performs to the channel
	 * and how much time it has spent waiting on a channel.
	 *
	 * At the end of the test, each routine reports min/max/avg wait
	 * as well as total events.
	 */
}

func main() {
	v := *flag.Bool("v", false, "set v flag")
	flag.Parse()

	fmt.Printf("v: %v\n", v)
	for _, fl := range flag.Args() {
		fmt.Printf("Flag: %v\n", fl)
	}

	/* Run test functions... */
	//tstest()       // Run timestamp test
	//chancommtest() // Run channel test
	mutextest() // Run mutex test
}
