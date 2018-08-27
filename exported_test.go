package librato_test

import (
	"sync"
	"testing"
	"time"

	"github.com/nyaruka/librato"
)

func TestExported(t *testing.T) {
	// all methods are NOOPs if the libray hasn't been configured
	librato.Start()
	librato.Gauge("foo.bar", 123.45)
	librato.Stop()

	librato.Configure("bob", "1234567", "foo.com", time.Second*30, &sync.WaitGroup{})
	librato.Start()
	librato.Stop()
}
