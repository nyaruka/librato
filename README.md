# Librato

Basic Librato client library with batching of events. Thus far it only supports sending gauge values
because that's all we need, but contributions welcome.

## Usage

You can either instantiate a collector and use that:

```go
import "github.com/nyaruka/librato"

collector := librato.NewCollector(...)
collector.Start()
collector.Gauge("awesomeness.level", 10)
collector.Gauge("foo.count", 123.45)
collector.Stop()
```

Or configure the default collector and use it like:

```go
librato.Configure(...)
librato.Start()
librato.Gauge("awesomeness.level", 10)
librato.Gauge("foo.count", 123.45)
librato.Stop()
```
