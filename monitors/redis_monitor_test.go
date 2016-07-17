// +build redis

package monitors

import (
	core "github.com/gerty-monit/core"
	"testing"
)

func TestShouldPingRedis(t *testing.T) {
	monitor := NewRedisMonitor("redis", "redis monitor", "localhost", 6379)
	status := monitor.Check()
	if status != core.OK {
		t.Errorf("failed, must be OK")
	}
}
