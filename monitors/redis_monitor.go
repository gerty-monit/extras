package monitors

import (
	"bufio"
	"fmt"
	core "github.com/gerty-monit/core"
	"net"
	"time"
)

type RedisMonitor struct {
	*core.BaseMonitor
	host   string
	port   int
	buffer core.CircularBuffer
	opts   *RedisMonitorOptions
}

type RedisMonitorOptions struct {
	Checks  int
	Timeout time.Duration
}

var DefaultRedisMonitorOptions = RedisMonitorOptions{
	Checks:  5,
	Timeout: 10 * time.Second,
}

func mergeRedisOpts(given *RedisMonitorOptions) *RedisMonitorOptions {
	if given == nil {
		return &DefaultRedisMonitorOptions
	}

	if given.Checks <= 0 {
		given.Checks = DefaultRedisMonitorOptions.Checks
	}

	if given.Timeout <= 0 {
		given.Timeout = DefaultRedisMonitorOptions.Timeout
	}

	return given
}

var _ core.Monitor = (*RedisMonitor)(nil)

func NewRedisMonitorWithOptions(title, description, host string, port int, _opts *RedisMonitorOptions) *RedisMonitor {
	opts := mergeRedisOpts(_opts)
	buffer := core.NewCircularBuffer(opts.Checks)
	return &RedisMonitor{core.NewBaseMonitor(title, description), host, port, buffer, opts}
}

func NewRedisMonitor(title, description, host string, port int) *RedisMonitor {
	return NewRedisMonitorWithOptions(title, description, host, port, nil)
}

func (monitor *RedisMonitor) Values() []core.ValueWithTimestamp {
	return monitor.buffer.GetValues()
}

func (monitor *RedisMonitor) Check() core.Result {
	address := fmt.Sprintf("%s:%d", monitor.host, monitor.port)
	conn, err := net.DialTimeout("tcp", address, monitor.opts.Timeout)

	if err != nil {
		// logger.Printf("tcp monitor check failed, error: %v", err)
		monitor.buffer.Append(core.NOK)
		return core.NOK
	}

	defer conn.Close()

	fmt.Fprintf(conn, "PING\n")
	status, err := bufio.NewReader(conn).ReadString('\n')
	pingResponse := "+PONG\r\n"
	if err != nil || status != pingResponse {
		monitor.buffer.Append(core.NOK)
		return core.NOK
	}

	monitor.buffer.Append(core.OK)
	return core.OK
}
