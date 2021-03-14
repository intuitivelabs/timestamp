// Copyright 2021 Intuitive Labs GmbH. All rights reserved.
//
// Use of this source code is governed by a BSD-style license
// that can be found in the LICENSE file in the root of the source
// tree.

package timestamp

import (
	"math/rand"
	"os"
	"testing"
	"time"
)

var seed int64

func TestMain(m *testing.M) {
	seed = time.Now().UnixNano()
	rand.Seed(seed)
	res := m.Run()
	os.Exit(res)
}

func TestNow(t *testing.T) {
	u := time.Now()
	ts := Now()
	diff := ts.SubTime(u)
	if diff < 0 {
		diff = -diff
	}
	if diff > time.Millisecond/10 {
		t.Errorf("Now() test failed, diff: %v\n", diff)
	}
}

func tstTruncateEq(t *testing.T, prefix string,
	ts TS, u time.Time, d time.Duration) bool {
	if ts.Truncate(d).Time() != u.UTC().Truncate(d) {
		t.Errorf(prefix+"failed Truncate equal test for: %s <> %s :"+
			"truncated to %v: %s <> %s\n",
			ts, u,
			d,
			ts.Truncate(d).Time(), u.UTC().Truncate(d))
		return false
	}
	return true
}

func tstAdd(t *testing.T, prefix string,
	ts TS, d time.Duration) int {
	var errs int
	u := ts.Time()
	ts2 := ts.Add(d)
	u2 := u.Add(d)
	errs = tstCmpEq(t, prefix, ts2, u2)
	if d != 0 && u.Equal(u2) {
		t.Errorf(prefix+"failed equal tst: %s <> %s \n", ts, ts2)
		errs++
	}
	if d != 0 && (ts2.Sub(ts) != d) {
		t.Errorf(prefix+"failed sub tst: %s <> %s \n", ts, ts2)
		errs++
	}
	if d != 0 && (ts2.SubTime(u) != d) {
		t.Errorf(prefix+"failed subTime tst: %s <> %s \n", u, ts2)
		errs++
	}
	return errs
}

// returns the number of failed tests
func tstCmpEq(t *testing.T, prefix string, ts TS, u time.Time) int {
	var errs int
	if OutOfRange(u) {
		t.Errorf(prefix+"failed OutOfRange test: %s\n", u)
		errs++
	}
	if !ts.EqualTime(u) || !u.Equal(ts.Time()) {
		t.Errorf(prefix+"failed Now equal test: ts %s <> t %s\n", ts, u)
		errs++
	}

	if ts.Unix() != u.Unix() || ts.UnixNano() != u.UnixNano() {
		t.Errorf(prefix+"failed Unix equal test: ts %s <> t %s\n", ts, u)
		errs++
	}

	if ts.AfterTime(u) || u.After(ts.Time()) {
		t.Errorf(prefix+"failed After now test: %s <> %s\n", ts, u)
		errs++
	}

	if ts.BeforeTime(u) || u.Before(ts.Time()) {
		t.Errorf(prefix+"failed Before now test: %s <> %s\n", ts, u)
		errs++
	}

	if !tstTruncateEq(t, prefix, ts, u, time.Hour) ||
		!tstTruncateEq(t, prefix, ts, u, time.Minute) ||
		!tstTruncateEq(t, prefix, ts, u, time.Second) ||
		!tstTruncateEq(t, prefix, ts, u, time.Millisecond) ||
		!tstTruncateEq(t, prefix, ts, u, time.Microsecond) {
		errs++
	}

	if ts.String() != u.UTC().String() {
		t.Errorf(prefix+"failed String equal test: ts %s <> t %s\n", ts, u)
		errs++
	}
	return errs
}

func TestTimeToTS(t *testing.T) {
	u := time.Time{}
	ts := Timestamp(u)
	if errs := tstCmpEq(t, "time ZERO equal: ", ts, u); errs != 0 {
		t.Errorf("time ZERO equal %d errors\n", errs)
	}
	if !ts.Time().IsZero() || !ts.IsZero() {
		t.Errorf("failed zero test: %s\n", ts)
	}

	if !OutOfRange(u.Add(time.Microsecond)) {
		t.Errorf("failed out of range test for: %s\n", u.Add(time.Microsecond))
	}

	u = time.Now()
	ts = Timestamp(u)
	if ts.Time().IsZero() || ts.IsZero() || Zero().Equal(ts) {
		t.Errorf("time.Now(): failed non zero test: %s\n", ts)
	}
	if errs := tstCmpEq(t, "time.Now equal: ", ts, u); errs != 0 {
		t.Errorf("time.Now() equal %d errors\n", errs)
	}

	u.Add(time.Second)
	ts.Add(time.Second)
	if ts.Time().IsZero() || ts.IsZero() || Zero().Equal(ts) {
		t.Errorf("add 1s: failed non zero test: %s\n", ts)
	}
	if errs := tstCmpEq(t, "add 1s: ", ts, u); errs != 0 {
		t.Errorf("add 1s: equal %d errors\n", errs)
	}

	u.Add(time.Minute)
	ts.Add(time.Minute)
	if ts.Time().IsZero() || ts.IsZero() || Zero().Equal(ts) {
		t.Errorf("add 1m: failed non zero test: %s\n", ts)
	}
	if errs := tstCmpEq(t, "add 1m: ", ts, u); errs != 0 {
		t.Errorf("add 1m: equal %d errors\n", errs)
	}

	u.Add(time.Hour)
	ts.Add(time.Hour)
	if ts.Time().IsZero() || ts.IsZero() || Zero().Equal(ts) {
		t.Errorf("add 1h: failed non zero test: %s\n", ts)
	}
	if errs := tstCmpEq(t, "add 1m: ", ts, u); errs != 0 {
		t.Errorf("add 1h: equal %d errors\n", errs)
	}
}

func TestRandTS(t *testing.T) {
	const cfgIterations = 1000
	for i := uint(0); i < cfgIterations; i++ {
		ts1 := TS(rand.Int63n(int64(MaxTS - TS(time.Hour))))
		u1 := ts1.Time()
		if errs := tstCmpEq(t, "rand ts+: ", ts1, u1); errs != 0 {
			t.Errorf("rand ts+: %d errors, rand seed %d\n", errs, seed)
		}
		ts2 := TS(-rand.Int63n(int64(-(MinTS + 1) - TS(time.Hour))))
		u2 := ts2.Time()
		if errs := tstCmpEq(t, "rand ts-: ", ts2, u2); errs != 0 {
			t.Errorf("rand ts2-: equal %d errors, rand seed %d\n", errs, seed)
		}
		d1 := time.Duration(rand.Int63n(int64(MaxTS - ts1)))
		d2 := time.Duration(-rand.Int63n(int64(-(MinTS + 1))))
		if errs := tstAdd(t, "rand add+: ", ts1, d1); errs != 0 {
			t.Errorf("rand add+: %d errors, rand seed %d\n", errs, seed)
		}
		if errs := tstAdd(t, "rand add-: ", ts1, d2); errs != 0 {
			t.Errorf("rand add-: %d errors, rand seed %d\n", errs, seed)
		}
	}
}

func TestRandT(t *testing.T) {
	const cfgIterations = 1000
	for i := uint(0); i < cfgIterations; i++ {
		ns := rand.Int63n(int64(MaxTS - 1))
		u1 := time.Unix(ns/1e9, ns%1e9)
		ts1 := Timestamp(u1)
		if errs := tstCmpEq(t, "rand ts+: ", ts1, u1); errs != 0 {
			t.Errorf("rand ts+: %d errors, rand seed %d\n", errs, seed)
		}
		d1 := time.Duration(rand.Int63n(int64(MaxTS - ts1)))
		d2 := time.Duration(-rand.Int63n(int64(-(MinTS + 1))))
		if errs := tstAdd(t, "rand add+: ", ts1, d1); errs != 0 {
			t.Errorf("rand add+: %d errors, rand seed %d\n", errs, seed)
		}
		if errs := tstAdd(t, "rand add-: ", ts1, d2); errs != 0 {
			t.Errorf("rand add-: %d errors, rand seed %d\n", errs, seed)
		}
	}
}
