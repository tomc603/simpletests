#simpletests
simpletests is a small suite of tests of the Go language to prove or disprove commonly repeated advice.

##Purpose
The purpose of this test will be to measure the difference between locking a struct with a mutex, and sending data over a channel to be stored by a listening goroutine in a struct. Each iteration of M modifications will be followed by an increase in goroutines writing to the struct/channel. At the conclusion of the test, N routines will be writing M data to the struct.

The measurement gathered will be the time a routine waited on a mutex lock, or the time between inserting data on a channel and that data being processed by the listening goroutine.

Future testing should include incrementing the CPU count that the routines may be spread among to reveal what delay, if any, is introduced by NUMA interactions.

##Functions
###tstest
This test calculates the speed of assigning time.Now() to a variable

###chancommtest
Tests the speed of communication over a channel by sending a struct to a listening routine, which adds a
timestamp and sends the struct back.

###mutextest
Tests the effectiveness of multiple goroutines modifying the same struct, using a sync.Mutex to lock the
struct for concurrency safety.

###chanwaittest
Tests the effectiveness of multiple goroutines modifying the same struct by sending data through a channel
to a listening goroutine, which handles the actual modification of the struct.

