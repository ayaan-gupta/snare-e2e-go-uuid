// Copyright 2023 Google Inc.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package uuid

// UUID version 6 is a field-compatible version of UUIDv1, reordered for improved DB locality.
// It is expected that UUIDv6 will primarily be used in contexts where there are existing v1 UUIDs.
// Systems that do not involve legacy UUIDv1 SHOULD consider using UUIDv7 instead.
//
// see https://datatracker.ietf.org/doc/html/draft-peabody-dispatch-new-uuid-format-03#uuidv6
//
// NewV6 returns a Version 6 UUID based on the current NodeID and clock
// sequence, and the current time. If the NodeID has not been set by SetNodeID
// or SetNodeInterface then it will be set automatically. If the NodeID cannot
// be set NewV6 set NodeID is random bits automatically . If clock sequence has not been set by
// SetClockSequence then it will be set automatically. If GetTime fails to
// return the current NewV6 returns Nil and an error.
func NewV6() (UUID, error) {
	var uuid UUID
	now, seq, err := GetTime()
	if err != nil {
		return uuid, err
	}

	// RFC 9562 places the 60-bit timestamp in three fields around the
	// version nibble: timestamp[59:28], timestamp[27:12], timestamp[11:0].
	timestamp := uint64(now)
	uuid[0] = byte(timestamp >> 52)
	uuid[1] = byte(timestamp >> 44)
	uuid[2] = byte(timestamp >> 36)
	uuid[3] = byte(timestamp >> 28)
	uuid[4] = byte(timestamp >> 20)
	uuid[5] = byte(timestamp >> 12)
	uuid[6] = 0x60 | byte(timestamp>>8&0x0f)
	uuid[7] = byte(timestamp)

	uuid[8] = 0x80 | byte(seq>>8&0x3f)
	uuid[9] = byte(seq)

	nodeMu.Lock()
	if nodeID == zeroID {
		setNodeInterface("")
	}
	copy(uuid[10:], nodeID[:])
	nodeMu.Unlock()

	return uuid, nil
}
