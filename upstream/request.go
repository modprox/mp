package upstream

import "fmt"

type Namespace []string

type Request struct {
	Transport string
	Domain    string
	Namespace Namespace
	Version   string
	Path      string
}

func (r *Request) String() string {
	return fmt.Sprintf(
		"[%q %q %v %q %q]",
		r.Transport,
		r.Domain,
		r.Namespace,
		r.Version,
		r.Path,
	)
}

func (r *Request) URI() string {
	return fmt.Sprintf("%s://%s/%s", r.Transport, r.Domain, r.Path)
}
