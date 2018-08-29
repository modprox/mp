package upstream

// what we really need is a thing which
// transforms a module into a URI usable for an http
// request - by applying each of the types of transforms:
//
// - domain alias
// - URL path creation based on domain
// - authentication / authorization configuration
