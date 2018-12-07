package stats

import "github.com/modprox/mp/pkg/netservice"

type Statsd struct {
	Agent netservice.Instance `json:"agent"`
}
