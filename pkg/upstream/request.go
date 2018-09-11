package upstream

import "fmt"

type Namespace []string

type Request struct {
	Transport     string
	Domain        string
	Namespace     Namespace
	Version       string
	Path          string
	GoGetRedirect bool
	Headers       map[string]string
}

func (r *Request) String() string {
	return fmt.Sprintf(
		"[%q %q %v %q %q %t]",
		r.Transport,
		r.Domain,
		r.Namespace,
		r.Version,
		r.Path,
		r.GoGetRedirect,
	)
}

// The URI is only valid AFTER a Request has passed through
// all of the Transform functors.
func (r *Request) URI() string {
	return fmt.Sprintf("%s://%s/%s", r.Transport, r.Domain, r.Path)
}
