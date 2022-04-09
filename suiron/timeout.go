package suiron

// Time Out - functions for measuring elapsed execution time.
//
// Cleve Lendon

import (
    "time"
)

// Maximum execution time.
var suironMaxTime   int64 = 1000 * 1_000_000  // 300 millisecond default
var suironStartTime time.Time
var suironZeroTime time.Time = time.Time{}
var suironHasTimedOut = false

// SetMaxTimeMilliseconds - sets the maximum execution time.
// Param: maxTime (in milliseconds)
func SetMaxTimeMilliseconds(maxTime int64) {
    if maxTime < 0 { return }
    suironMaxTime = maxTime * 1_000_000 // convert to nanoseconds
}

// SetStartTime - sets the start time before starting the query.
func SetStartTime() {
    suironStartTime = time.Now()
    suironHasTimedOut = false
}

// ClearStartTime - clears the start time to 0.
// This will prevent a time-out.
func ClearStartTime() {
    suironStartTime = suironZeroTime
    suironHasTimedOut = false
}

// ElapsedTime - returns time (in nanoseconds) since the start of the query.
func ElapsedTime() int64 {
    return int64(time.Since(suironStartTime))
}

// HasTimedOut - returns true if the maximum execution time has been exceeded.
// If no start time was set (zero time), return false.
func HasTimedOut() bool {
    if suironHasTimedOut == true { return true }
    if suironStartTime == suironZeroTime { return false }
    if int64(time.Since(suironStartTime)) > suironMaxTime {
        suironHasTimedOut = true
        return true
    }
    return false
}

// MakeTimer - makes a timer to limit the runtime of the inference engine.
func MakeTimer() *time.Timer {
    return time.NewTimer(time.Duration(suironMaxTime))
}
