# timestamp

[![Go Reference](https://pkg.go.dev/badge/github.com/intuitivelabs/timestamp.svg)](https://pkg.go.dev/github.com/intuitivelabs/timestamp)

The timestamp package provides a time.Time "compatible" way of keeping time
stamps, using less space (single int64) and no pointers.

It implements most of time.Time functions, so switching to it should require
 only minimal work (most of the time the only thing required is
 replacing time.Time with timestamp.TS and time.Now() with timestamp.Now()).

Functions are provided for converting between time.Time and timestamp.TS
 (timestamp.Timestamp(t)) and vice-versa (timestamp.Time()).

The time stamp is kept in nanoseconds.

## Limitations

* timestamps are limited to ~ 1970 +/- 292 years

* to check for out-of-range when converting from a time.Time value, one 
 has to compare both with timestamp.MinTS and timestamp.MaxTS

* string representation uses always UTC (since timezones are not supported).
To use a different representation one has to go through time.Time, for
 example via timestamp.In:

  ```
   ts := timestamp.Now()
   fmt.Printf("ts UTC: %s , local: %s\n", ts, ts.In(time.Local))
  ```
