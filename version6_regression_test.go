package uuid

import "testing"

func TestVersion6RFC9562Timestamp(t *testing.T) {
	uuid := Must(Parse("1EC9414C-232A-6B00-B3C8-9E6BDECED846"))
	if got, want := uint64(uuid.Time()), uint64(0x1ec9414c232ab00); got != want {
		t.Fatalf("Time() = %#x, want %#x", got, want)
	}
}

func TestVersion6TimestampRoundTripAndOrdering(t *testing.T) {
	for _, timestamp := range []Time{0, 1, 0x123456789abcdef, (1 << 60) - 1} {
		uuid := UUID{}
		uuid[0] = byte(uint64(timestamp) >> 52)
		uuid[1] = byte(uint64(timestamp) >> 44)
		uuid[2] = byte(uint64(timestamp) >> 36)
		uuid[3] = byte(uint64(timestamp) >> 28)
		uuid[4] = byte(uint64(timestamp) >> 20)
		uuid[5] = byte(uint64(timestamp) >> 12)
		uuid[6] = 0x60 | byte(uint64(timestamp)>>8&0xf)
		uuid[7] = byte(timestamp)
		uuid[8] = 0x80
		if uuid.Time() != timestamp {
			t.Errorf("timestamp %#x round trip = %#x", timestamp, uuid.Time())
		}
	}

	u1 := Must(Parse("00000000-0000-6000-8000-000000000000"))
	u2 := Must(Parse("00000000-0000-6001-8000-000000000000"))
	if u1.String() >= u2.String() {
		t.Fatalf("UUID ordering does not follow timestamp: %s >= %s", u1, u2)
	}
}
