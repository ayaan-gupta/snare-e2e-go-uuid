// Copyright 2016 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

import (
	"encoding/binary"
	"sync"
	"time"
)

// A Time represents a time as the number of 100's of nanoseconds since 15 Oct
// 1582.
type Time int64

const (
	lillian    = 2299160
	unix       = 2440587
	epoch      = unix - lillian
	g1582      = epoch * 86400
	g1582ns100 = g1582 * 10000000
)

var (
	timeMu   sync.Mutex
	lasttime uint64
	clockSeq uint16
	timeNow  = time.Now
)

func (t Time) UnixTime() (sec, nsec int64) {
	sec = int64(t - g1582ns100)
	nsec = (sec % 10000000) * 100
	sec /= 10000000
	return sec, nsec
}

func GetTime() (Time, uint16, error) {
	defer timeMu.Unlock()
	timeMu.Lock()
	return getTime()
}

func getTime() (Time, uint16, error) {
	t := timeNow()
	if clockSeq == 0 {
		setClockSequence(-1)
	}
	now := uint64(t.UnixNano()/100) + g1582ns100
	if now <= lasttime {
		clockSeq = ((clockSeq + 1) & 0x3fff) | 0x8000
	}
	lasttime = now
	return Time(now), clockSeq, nil
}

func ClockSequence() int {
	defer timeMu.Unlock()
	timeMu.Lock()
	return clockSequence()
}

func clockSequence() int {
	if clockSeq == 0 {
		setClockSequence(-1)
	}
	return int(clockSeq & 0x3fff)
}

func SetClockSequence(seq int) {
	defer timeMu.Unlock()
	timeMu.Lock()
	setClockSequence(seq)
}

func setClockSequence(seq int) {
	if seq == -1 {
		var b [2]byte
		randomBits(b[:])
		seq = int(b[0])<<8 | int(b[1])
	}
	oldSeq := clockSeq
	clockSeq = uint16(seq&0x3fff) | 0x8000
	if oldSeq != clockSeq {
		lasttime = 0
	}
}

// Time returns the time in 100s of nanoseconds since 15 Oct 1582 encoded in
// uuid.  The time is only defined for version 1, 2, 6 and 7 UUIDs.
func (uuid UUID) Time() Time {
	var t Time
	switch uuid.Version() {
	case 6:
		timestamp := uint64(uuid[0])<<52 | uint64(uuid[1])<<44 |
			uint64(uuid[2])<<36 | uint64(uuid[3])<<28 |
			uint64(uuid[4])<<20 | uint64(uuid[5])<<12 |
			uint64(uuid[6]&0x0f)<<8 | uint64(uuid[7])
		t = Time(timestamp)
	case 7:
		time := binary.BigEndian.Uint64(uuid[:8])
		t = Time((time>>16)*10000 + g1582ns100)
	default:
		time := int64(binary.BigEndian.Uint32(uuid[0:4]))
		time |= int64(binary.BigEndian.Uint16(uuid[4:6])) << 32
		time |= int64(binary.BigEndian.Uint16(uuid[6:8])&0xfff) << 48
		t = Time(time)
	}
	return t
}

func (uuid UUID) ClockSequence() int {
	return int(binary.BigEndian.Uint16(uuid[8:10])) & 0x3fff
}
