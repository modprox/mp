package stats

import "oss.indeed.com/go/modprox/pkg/netservice"

type Statsd struct {
	Agent netservice.Instance `json:"agent"`
}
