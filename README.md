# librato [![Build Status](https://github.com/nyaruka/librato/workflows/CI/badge.svg)](https://github.com/nyaruka/librato/actions?query=workflow%3ACI) [![codecov](https://codecov.io/gh/nyaruka/librato/branch/main/graph/badge.svg)](https://codecov.io/gh/nyaruka/librato) [![Go Report Card](https://goreportcard.com/badge/github.com/nyaruka/librato)](https://goreportcard.com/report/github.com/nyaruka/librato)

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
