package upstream

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_parseGoGetMetadata(t *testing.T) {
	try := func(input string, exp goGetMeta, expErr bool) {
		output, err := parseGoGetMetadata(input)
		if expErr {
			require.Error(t, err)
			return
		}
		require.Equal(t, exp, output)
	}

	try(metaHTML1, goGetMeta{
		transport: "https",
		domain:    "github.com",
		path:      "apache/thrift",
	}, false)

	try(metaHTML2, goGetMeta{
		transport: "https",
		domain:    "github.com",
		path:      "go-yaml/yaml/tree/v2.2.1",
	}, false)

	try(metaHTML3, goGetMeta{
		transport: "https",
		domain:    "github.com",
		path:      "golang/net",
	}, false)

	try(metaHTML4, goGetMeta{
		transport: "https",
		domain:    "github.com",
		path:      "GoogleCloudPlatform/gcloud-golang",
	}, false)

	try(metaHTML5, goGetMeta{
		transport: "https",
		domain:    "github.com",
		path:      "census-ecosystem/opencensus-go-exporter-stackdriver",
	}, false)

	try(metaHTML6, goGetMeta{ // does this work?
		transport: "https",
		domain:    "dmitri.shuralyov.com",
		path:      "text/kebabcase",
	}, false)

	try(metaHTML7, goGetMeta{
		transport: "https",
		domain:    "github.com",
		path:      "census-instrumentation/opencensus-go",
	}, false)

	try(metaHTML8, goGetMeta{
		transport: "https",
		domain:    "github.com",
		path:      "uber-go/atomic",
	}, false)

	try(metaHTML9, goGetMeta{
		transport: "https",
		domain:    "github.com",
		path:      "google/google-api-go-client",
	}, false)
}

const (
	metaHTML1 = `<meta name="go-import" content="github.com/apache/thrift git https://github.com/apache/thrift.git">`
	metaHTML2 = `
<head>
<meta name="go-import" content="gopkg.in/yaml.v2 git https://gopkg.in/yaml.v2">
<meta name="go-source" content="gopkg.in/yaml.v2 _ https://github.com/go-yaml/yaml/tree/v2.2.1{/dir} https://github.com/go-yaml/yaml/blob/v2.2.1{/dir}/{file}#L{line}">
</head>
<body>
`
	metaHTML3 = `
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="golang.org/x/net git https://go.googlesource.com/net">
<meta name="go-source" content="golang.org/x/net https://github.com/golang/net/ https://github.com/golang/net/tree/master{/dir} https://github.com/golang/net/blob/master{/dir}/{file}#L{line}">
<meta http-equiv="refresh" content="0; url=https://godoc.org/golang.org/x/net">
</head>
`

	metaHTML4 = `
<head>
  <meta name="go-import" content="cloud.google.com/go git https://code.googlesource.com/gocloud">
  <meta name="go-source" content="cloud.google.com/go https://github.com/GoogleCloudPlatform/gcloud-golang https://github.com/GoogleCloudPlatform/gcloud-golang/tree/master{/dir} https://github.com/GoogleCloudPlatform/gcloud-golang/tree/master{/dir}/{file}#L{line}">
  <meta http-equiv="refresh" content="0; url=/go/home">
</head>
`

	metaHTML5 = `
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="contrib.go.opencensus.io/exporter/stackdriver git https://github.com/census-ecosystem/opencensus-go-exporter-stackdriver">
<meta name="go-source" content="contrib.go.opencensus.io/exporter/stackdriver https://github.com/census-ecosystem/opencensus-go-exporter-stackdriver https://github.com/census-ecosystem/opencensus-go-exporter-stackdriver/tree/master{/dir} https://github.com/census-ecosystem/opencensus-go-exporter-stackdriver/blob/master{/dir}/{file}#L{line}">
<meta http-equiv="refresh" content="0; url=https://godoc.org/contrib.go.opencensus.io/exporter/stackdriver/">
</head>
`

	metaHTML6 = `
<meta name="go-import" content="dmitri.shuralyov.com/text/kebabcase git https://dmitri.shuralyov.com/text/kebabcase">
<meta name="go-source" content="dmitri.shuralyov.com/text/kebabcase https://dmitri.shuralyov.com/text/kebabcase https://gotools.org/dmitri.shuralyov.com/text/kebabcase https://gotools.org/dmitri.shuralyov.com/text/kebabcase#{file}-L{line}">
`

	metaHTML7 = `
<head>
<meta http-equiv="Content-Type" content="text/html; charset=utf-8"/>
<meta name="go-import" content="go.opencensus.io git https://github.com/census-instrumentation/opencensus-go">
<meta name="go-source" content="go.opencensus.io https://github.com/census-instrumentation/opencensus-go https://github.com/census-instrumentation/opencensus-go/tree/master{/dir} https://github.com/census-instrumentation/opencensus-go/blob/master{/dir}/{file}#L{line}">
<meta http-equiv="refresh" content="0; url=https://godoc.org/go.opencensus.io/">
</head>
`

	metaHTML8 = `
<head>
    <meta name="go-import" content="go.uber.org/atomic git https://github.com/uber-go/atomic">
    <meta name="go-source" content="go.uber.org/atomic https://github.com/uber-go/atomic https://github.com/uber-go/atomic/tree/master{/dir} https://github.com/uber-go/atomic/tree/master{/dir}/{file}#L{line}">
    <meta http-equiv="refresh" content="0; url=https://godoc.org/go.uber.org/atomic">
</head>
`

	metaHTML9 = `
<head>
<meta name="go-import" content="google.golang.org/api git https://code.googlesource.com/google-api-go-client">
<meta name="go-source" content="google.golang.org/api https://github.com/google/google-api-go-client https://github.com/google/google-api-go-client/tree/master{/dir} https://github.com/google/google-api-go-client/tree/master{/dir}/{file}#L{line}">
<meta http-equiv="refresh" content="0; url=https://godoc.org/google.golang.org/api">
</head>

`
)

/*
match["golang.org"] = true
match["cloud.google.com"] = true
match["google.golang.org"] = true
match["gopkg.in"] = true
match["contrib.go.opencensus.io"] = true
match["go.opencensus.io"] = true
match["go.uber.org"] = true
match["git.apache.org"] = true
*/
