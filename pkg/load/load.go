//
// Load average resource.
//
// This collector reports on the following stat metrics:
//
//  - "load1" (gauge)
//  - "load5" (gauge)
//  - "load15" (gauge)
//
package load

import "github.com/statsd/client-interface"
import "github.com/c9s/goprocinfo/linux"
import "github.com/segmentio/go-log"
import "time"

// Load resource.
type Load struct {
	Path     string
	Interval time.Duration
	Extended bool
	client   statsd.Client
	exit     chan struct{}
}

// New Load resource.
func New(interval time.Duration) *Load {
	return &Load{
		Path:     "/proc/loadavg",
		Interval: interval,
		exit:     make(chan struct{}),
	}
}

// Name of the resource.
func (c *Load) Name() string {
	return "load"
}

// Start resource collection.
func (c *Load) Start(client statsd.Client) error {
	c.client = client
	go c.Report()
	return nil
}

// Report resource collection.
func (c *Load) Report() {
	tick := time.Tick(c.Interval)

	for {
		select {
		case <-tick:
			stat, err := linux.ReadLoadAvg(c.Path)

			if err != nil {
				log.Error("load: %s", err)
				continue
			}

			c.client.Gauge("load1", int(stat.Last1Min * 100))
			c.client.Gauge("load5", int(stat.Last5Min * 100))
			c.client.Gauge("load15", int(stat.Last15Min * 100))

		case <-c.exit:
			log.Info("load: exiting")
			return
		}
	}
}

// Stop resource collection.
func (c *Load) Stop() error {
	println("stopping load")
	close(c.exit)
	return nil
}
