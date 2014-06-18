package main

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
			// fmt.Print("Time elapsed...\n")
			done = true
		case <-c.C:
			wc++
			// fmt.Printf("Adding mutexworker %d...\n", wc)
			w.Add(1)
			go mutexworker(&r, &w, wc)
		}
		if done {
			break
		}
	}
	fmt.Print("Waiting for mutexworkers to complete...")
	w.Wait()
	fmt.Print("done...\n")
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
				// fmt.Printf("(%d) Setting mutex min: %d\n", i, delay.Nanoseconds())
				r.min = delay
			case delay > r.max:
				// fmt.Printf("(%d) Setting mutex max: %d\n", i, delay.Nanoseconds())
				r.max = delay
			case delay < r.min:
				// fmt.Printf("(%d) Setting mutex min: %d\n", i, delay.Nanoseconds())
				r.min = delay
			}

			r.count += 1
			r.Unlock()
		}
		if done {
			break
		}
	}
	// fmt.Printf("mutexworker(%d): time elapsed\n", i)
	w.Done()
}
