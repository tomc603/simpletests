package main

func chanpubsub() {
	/* This test spawns a listener goroutine as well as a number of
	 * writer goroutines. Each writer sends data over a channel to the
	 * listener routine, which aggregates that data into a struct.
	 * Each routine counts the number of writes it performs to the channel
	 * and how much time it has spent waiting on a channel.
	 *
	 * At the end of the test, each routine reports min/max/avg wait
	 * as well as total events.
	 */

	var (
		pubcount uint
		quit     bool
		wg       sync.WaitGroup
	)

	c := make(chan roundtrip, 2048)
	stoptimer := time.NewTimer(30 * time.Second)
	spawntimer := time.NewTicker(1 * time.Second)

	fmt.Print("Starting channel test\n")

	// Spawn subscriber goroutine
	go chansub(c, &wg, &quit)

	// Spawn publisher goroutines
	//s := time.Now()
	for {
		select {
		case <-stoptimer.C:
			// Time has elapsed. Stop the routine.
			stoptimer.Stop()
			spawntimer.Stop()
			fmt.Print("Time elapsed...\n")
			quit = true
		case <-spawntimer.C:
			pubcount++
			// fmt.Printf("Adding channel publisher %d...\n", pubcount)
			go chanpub(c, &wg, &quit, pubcount)
		}
		if quit {
			break
		}
	}
	fmt.Print("Waiting for channel workers to complete...\n")
	wg.Wait()
	// fmt.Printf("Count: %d (%f/s)\nMin/Max/Avg RTT: %dns/%dns/%dns\n\n",
	//  r.count, (float64(r.count) / time.Since(s).Seconds()),
	//  r.min.Nanoseconds(), r.max.Nanoseconds(), r.delaytotal.Nanoseconds()/int64(r.count))
}

func chanpub(c chan roundtrip, w *sync.WaitGroup, q *bool, id uint) {
	// Add this routine to the waitgroup so we don't terminate too early,
	// defer marking it done until the routine quits.
	w.Add(1)
	defer w.Done()

	// var (
	//  txCount    int64
	//  dCurrent   time.Duration
	//  dMax       time.Duration
	//  dMin       time.Duration
	//  dTotal     time.Duration
	//  tStart     time.Time
	//  tLoopStart time.Time
	//  tTotal     time.Duration
	// )

	// tStart = time.Now()
	for {
		// tLoopStart = time.Now()
		c <- roundtrip{sent: time.Now()}
		// dCurrent = time.Since(tLoopStart)
		// dTotal += dCurrent
		// txCount += 1

		// switch {
		// case dMin == 0 || dCurrent < dMin:
		//  // fmt.Printf("(pub-%2.0d) Setting min: %d\n", id, dCurrent.Nanoseconds())
		//  dMin = dCurrent
		// case dCurrent > dMax:
		//  // fmt.Printf("(pub-%2.0d) Setting max: %d\n", id, dCurrent.Nanoseconds())
		//  dMax = dCurrent
		// }

		if *q {
			// tTotal = time.Since(tStart)
			break
		}
	}
	// fmt.Printf("(pub-%2.0d) Sent: %d (%2.2f/sec)\nWait total/min/max/avg: %dns/%dns/%dns/%dns\n",
	//  id, txCount, float64(txCount)/tTotal.Seconds(), dTotal.Nanoseconds(),
	//  dMin.Nanoseconds(), dMax.Nanoseconds(), dTotal.Nanoseconds()/txCount)
}

func chansub(c chan roundtrip, w *sync.WaitGroup, q *bool) {
	// Add this routine to the waitgroup so we don't terminate too early,
	// defer marking it done until the routine quits.
	w.Add(1)
	defer w.Done()

	var (
		rxCount  int64
		dCurrent time.Duration
		dMax     time.Duration
		dMin     time.Duration
		dTotal   time.Duration
		tStart   time.Time
		tTotal   time.Duration
	)

	tStart = time.Now()
	for {
		select {
		case r := <-c:
			// Receive struct from channel, calculate delay, min, max,
			// and increment the counter
			dCurrent = time.Since(r.sent)
			dTotal += dCurrent
			rxCount += 1

			switch {
			case dMin == 0 || dCurrent < dMin:
				// fmt.Printf("(sub) Setting min: %d\n", dCurrent.Nanoseconds())
				dMin = dCurrent
			case dCurrent > dMax:
				// fmt.Printf("(sub) Setting max: %d\n", dCurrent.Nanoseconds())
				dMax = dCurrent
			}
		default:
		}
		if *q {
			tTotal = time.Since(tStart)
			break
		}
	}
	fmt.Printf("(sub) Received: %d (%2.2f/sec)\nDelay total/min/max/avg: %dns/%dns/%dns/%dns\n",
		rxCount, float64(rxCount)/tTotal.Seconds(), dTotal.Nanoseconds(),
		dMin.Nanoseconds(), dMax.Nanoseconds(), dTotal.Nanoseconds()/rxCount)
}
