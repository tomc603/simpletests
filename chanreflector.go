package main

func chancommtest() {
	/* Test communication speed of channels by sending a struct
	 * with sent and received timestamps over a channel for 10 seconds.
	 * Structure to pass over channels */
	var (
		loopCount uint64
		total     time.Duration
		min       time.Duration
		max       time.Duration
		wg        sync.WaitGroup
	)

	send := make(chan *roundtrip)
	receive := make(chan *roundtrip)
	rt := &roundtrip{}
	st := time.Now()

	fmt.Print("Starting Channel communication test\n")
	/* Start goroutine to listen for chan activity */
	go receiver(send, receive, &wg)

	for time.Since(st).Seconds() < 10 {
		/* Send an item over the channel, increment counter */
		rt.sent = time.Now()
		fmt.Printf("0: Sending: %v", rt)
		send <- rt
		r := <-receive
		fmt.Print(" - RECEIVED\n")
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
	close(send)
	wg.Wait()
	fmt.Printf("Count: %d (%f/s)\nMin/Max/Avg RTT: %dns/%dns/%dns\n\n",
		loopCount, (float64(loopCount) / time.Since(st).Seconds()),
		min.Nanoseconds(), max.Nanoseconds(), uint64(total.Nanoseconds())/loopCount)
}

func receiver(receive chan *roundtrip, send chan *roundtrip, wg *sync.WaitGroup) {
	wg.Add(1)
	defer wg.Done()
	for {
		select {
		case r := <-receive:
			fmt.Print(" - RECEIVED\n")
			r.received = time.Now()
			fmt.Printf("1: Sending: %v", r)
			send <- r
		}
	}
}
