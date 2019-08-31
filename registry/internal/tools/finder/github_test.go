package finder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_githubCommit_Pseudo(t *testing.T) {
	var gc githubCommit
	gc.SHA = "c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc"
	gc.Commit.Author.Date = "2011-01-26T19:06:43Z"

	pseudo, err := gc.Pseudo([]Tag{}, true)
	require.NoError(t, err)
	require.Equal(
		t,
		"v0.0.0-20110126190643-c5b97d5ae6c1",
		pseudo,
	)
}

func Test_githubCommit_Pseudo_bad_time(t *testing.T) {
	var gc githubCommit
	gc.SHA = "c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc"
	gc.Commit.Author.Date = "2011-01-26T19:06:43" // no Z

	_, err := gc.Pseudo([]Tag{}, true)
	require.Error(t, err)
}

func Test_githubCommit_Pseudo_previous_pre_version(t *testing.T) {
	var gc githubCommit
	gc.SHA = "c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc"
	gc.Commit.Author.Date = "2011-01-26T19:06:43Z"

	pseudo, err := gc.Pseudo([]Tag{{SemVer: "v1.2.4-pre"}, {SemVer: "v1.2.3"}}, true)
	require.NoError(t, err)
	require.Equal(t, "v1.2.4-pre.0.20110126190643-c5b97d5ae6c1", pseudo)
}

func Test_githubCommit_Pseudo_previous_version(t *testing.T) {
	var gc githubCommit
	gc.SHA = "c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc"
	gc.Commit.Author.Date = "2011-01-26T19:06:43Z"

	pseudo, err := gc.Pseudo([]Tag{{SemVer: "v1.2.4"}, {SemVer: "v1.2.3"}}, true)
	require.NoError(t, err)
	require.Equal(t, "v1.2.5-0.20110126190643-c5b97d5ae6c1", pseudo)
}

// v3.0.0-rc.2 v3.0.0-rc.1 v2.2.0 2.2.0 2.1.0 2.0.0 1.1.0 1.0.0

func Test_parseSemVer(t *testing.T) {
	require.Equal(t, &SemVer{1, 2, 3, false}, parseSemVer("v1.2.3"))
	require.Equal(t, &SemVer{1, 2, 0, false}, parseSemVer("v1.2"))
	require.Equal(t, &SemVer{1, 0, 0, false}, parseSemVer("v1"))
	require.Nil(t, parseSemVer("foo"))
}
