package mgr

import (
	"testing"
	"time"
)

func TestIsNotSend(t *testing.T) {
	for i, v := range []struct {
		to     *ToRedis
		not    bool
		addsec int64
	}{
		{
			to: &ToRedis{
				toExpireMinSec: 0,
				toExpireMaxSec: 0,
				isToDB:         false,
			},
			not:    true,
			addsec: 3600,
		},
		{
			to: &ToRedis{
				toExpireMinSec: 0,
				toExpireMaxSec: 0,
				isToDB:         true,
			},
			not:    false,
			addsec: 3600,
		},
		{
			to: &ToRedis{
				toExpireMinSec: 0,
				toExpireMaxSec: 0,
				isToDB:         true,
			},
			not:    true,
			addsec: -3600,
		},
		{
			to: &ToRedis{
				toExpireMinSec: 3600,
				toExpireMaxSec: 0,
				isToDB:         true,
			},
			not:    true,
			addsec: 3600,
		},
		{
			to: &ToRedis{
				toExpireMinSec: 0,
				toExpireMaxSec: 3600,
				isToDB:         true,
			},
			not:    false,
			addsec: 1800,
		},
		{
			to: &ToRedis{
				toExpireMinSec: 1800,
				toExpireMaxSec: 3600,
				isToDB:         true,
			},
			not:    true,
			addsec: 900,
		},
	} {
		if isNot := v.to.isNotSend(expiryMs(v.addsec)); isNot != v.not {
			t.Fatalf("Error: %v %#v %v", i, v, expiryMs(v.addsec))
		}
	}
}

func expiryMs(sec int64) int64 {
	return now.Add(time.Duration(sec)*time.Second).UnixNano() / 1000000
}
