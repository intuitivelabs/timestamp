// Copyright 2021 Intuitive Labs GmbH. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE.txt file in the root of the source
// tree.

// Package timestamp provides a smaller pointer-less equivalent for time.Time.

package timestamp

import (
	"sync/atomic"
	"time"
)

// A TS represents a time stamp with nanosecond precision.
// The zero value (Zero() method) correspond to the 0 unix time:
// 01.01.1970 00:00:00  UTC.
//
// The timestamp range is approx. the zero unix time (1970) +/- 292 years.
//
// All conversion to time.Time use UTC.
type TS time.Duration

const (
	MaxTS TS = 1<<63 - 1
	MinTS TS = -1 << 63
)

// reference time.Time for TS 0
// TS can represent a value between tZero - 292 years and tZero + 292 years.
var tZero time.Time = time.Unix(0, 0).UTC()

// Timestamp returns the time stamp corresponding to the given time (time.Time).// If it exceeds the maximum TS range it will return MaxTS or MinTS.
// The zero time.Time is specially handled and converted to a 0 timestamp
// (otherwise would be an out-of-range value).
func Timestamp(t time.Time) TS {
	if t.IsZero() {
		// 0 time.Time corresponds to year 1, and exceeds TS representation
		//  ( zero unix time +/- 292 years) => force return 0 in this case
		return 0
	}
	return TS(t.Sub(tZero))
}

// OutOfRange returns true if t is out of the timestamp representation
// range ( ~ 1970 +/- 292 years).
func OutOfRange(t time.Time) bool {
	d := t.Sub(tZero)
	return !t.IsZero() &&
		(d <= time.Duration(MinTS) || d >= time.Duration(MaxTS))
}

// DurationToTS returns the time stamp corresponding to a given time.Duration.
func DurationToTS(d time.Duration) TS {
	return TS(d)
}

// Zero returns the time stamp corresponding to the zero value
func Zero() TS {
	return Timestamp(tZero)
}

// Now returns the current time as a time stamp.
func Now() TS {
	return Timestamp(time.Now())
}

// Unix returns the local time stamp corresponding to the given unix time.
// See time.Unix for more details.
func Unix(sec int64, nsec int64) TS {
	return Timestamp(time.Unix(sec, nsec))
}

// AtomicStore changes atomically ts value.
func AtomicStore(ts *TS, v TS) {
	atomic.StoreInt64((*int64)(ts), int64(v))
}

// AtomicLoad reads atomically the ts value.
func AtomicLoad(ts *TS) TS {
	return TS(atomic.LoadInt64((*int64)(ts)))
}

// AtomicSwap stores atomically a new TS value and returns the previous one.
func AtomicSwap(ts *TS, v TS) TS {
	return TS(atomic.SwapInt64((*int64)(ts), int64(v)))
}

// AtomicCompareAndSwap stores atomically newv if the current value
// is oldv. It returns true on success (values swapped).
func AtomicCompareAndSwap(ts *TS, oldv TS, newv TS) bool {
	return atomic.CompareAndSwapInt64((*int64)(ts), int64(oldv), int64(newv))
}

// Duration returns the time stamp converted to duration.
func (ts TS) Duration() time.Duration {
	return time.Duration(ts)
}

// Time returns the time stamp converted to time.
// The 0 timestamp value is specially handled and converted to the zero
// time.Time (since otherwise the 0 timestamp would correpond to
// time.Unix(0,0).UTC()).
func (ts TS) Time() time.Time {
	if ts.IsZero() {
		return time.Time{}
	}
	return tZero.Add(ts.Duration())
}

// In returns the time stamp converted to time in the specified location.
// See also Time() and time.In(loc).
func (ts TS) In(loc *time.Location) time.Time {
	return ts.Time().In(loc)
}

// Location returns the time zone information associated to the timestamp.
// It will always return UTC, since timestamps use only UTC.
func (ts TS) Location() *time.Location {
	return time.UTC
}

// Add returns the time stamp corresponding to ts+d.
func (ts TS) Add(d time.Duration) TS {
	return ts + TS(d)
}

// AddTS returns the time stamp corresponding to ts+ts2.
func (ts TS) AddTS(ts2 TS) TS {
	return ts + ts2
}

// AddDate returns the time stamp corresponding to ts + date.
func (ts TS) AddDate(years, months, days int) TS {
	date := time.Date(years, time.Month(months), days, 0, 0, 0, 0, nil)
	return ts.AddTS(Timestamp(date))
}

// Sub returns the difference ts - ts2.
func (ts TS) Sub(ts2 TS) time.Duration {
	return (ts - ts2).Duration()
}

// Sub returns the difference ts - t., where t is a time.Time.
// Note that the result is undefined if t is outside the timestamp
// representation range.
func (ts TS) SubTime(t time.Time) time.Duration {
	return ts.Time().Sub(t)
}

// After returns true if ts > ts.
func (ts TS) After(ts2 TS) bool {
	return ts > ts2
}

// After returns true if ts > t, where t is a time.Time.
// Note that the result is undefined if t is outside the timestamp
// representation range.
func (ts TS) AfterTime(t time.Time) bool {
	return ts.Time().After(t)
}

// Before returns if ts < ts2.
func (ts TS) Before(ts2 TS) bool {
	return ts < ts2
}

// Before returns if ts < t, where t is a time.Time
// Note that the result is undefined if t is outside the timestamp
// representation range.
func (ts TS) BeforeTime(t time.Time) bool {
	return ts.Time().Before(t)
}

// Equal returns if ts == ts2.
func (ts TS) Equal(ts2 TS) bool {
	return ts == ts2
}

// EqualTime returns if ts == t, where t is a time.Time.
// Note that the result is undefined if t is outside the timestamp
// representation range.
func (ts TS) EqualTime(t time.Time) bool {
	return ts.Time().Equal(t)
}

// EqualTS returns if ts == ts2
func (ts TS) EqualTS(ts2 TS) bool {
	return ts == ts2
}

// Format returns a string representation of ts using time.Format(layout).
func (ts TS) Format(layout string) string {
	return ts.Time().Format(layout)
}

// IsZero returns if ts represents the zero timestamp.
func (ts TS) IsZero() bool {
	return ts == 0
}

// Truncate returns the result of rounding ts down to d (see time.Truncate()
// for more details).
func (ts TS) Truncate(d time.Duration) TS {
	if d <= 0 {
		return ts
	}
	// take into account rounding down negative numbers (and not toward 0)
	// => use the sign of the remainder
	return (ts/TS(d) + (ts%TS(d))>>63) * TS(d)
	//return Timestamp(ts.Time().Truncate(d))
}

// Truncate returns the result of rounding ts down to d (see time.Truncate()
// for more details).
func (ts TS) TruncateTime(d time.Duration) time.Time {
	return ts.Time().Truncate(d)
}

// Unix returns ts as Unix time (number of seconds since January 1, 1970 UTC)
func (ts TS) Unix() int64 {
	return ts.Time().Unix()
}

// UnixNanno returns ts as Unix time (number of nanoseconds since
// January 1, 1970 UTC)
func (ts TS) UnixNano() int64 {
	return ts.Time().UnixNano()
}

// Sting returns the default string representation for the timestamp.
// (see time.String for more details)
// Since timestamps do not store a timezone, the string representation will
// always use UTC.
func (ts TS) String() string {
	return ts.Time().String()
}

// MarshalBinary implements the encoding.BinaryMarshaler interface.
func (ts TS) MarshalBinary() ([]byte, error) {
	return ts.Time().MarshalBinary()
}

// MarshalJSON implements the json.Marshaler interface.
func (ts TS) MarshalJSON() ([]byte, error) {
	return ts.Time().MarshalJSON()
}

// MarshalText implements the encoding.TextMarshaler interface.
func (ts TS) MarshalText() ([]byte, error) {
	return ts.Time().MarshalText()
}

/*

// UnmarshalBinary implements the encoding.BinaryUnmarshaler interface.
func (ts *TS) UnmarshalBinary(data []byte) error {
	var t time.Time
	var err error
	if err = t.UnmarshalBinary(data); err == nil {
		*ts = Timestamp(t)
	}
	return err
}

// UnmarshalJSON implements the json.Marshaler interface.
func (ts *TS) UnmarshalJSON(data []byte) error {
	var t time.Time
	var err error
	if err = t.UnmarshalJSON(data); err == nil {
		*ts = Timestamp(t)
	}
	return err
}

// UnmarshalText implements the encoding.TextMarshaler interface.
func (ts *TS) UnmarshalText(data []byte) error {
	var t time.Time
	var err error
	if err = t.UnmarshalText(data); err == nil {
		*ts = Timestamp(t)
	}
	return err
}
*/
