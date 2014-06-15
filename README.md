#simpletests

##Purpose
The purpose of these tests is to measure and prove or disprove many common functions and idioms in Go.

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

