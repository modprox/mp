package finder

import (
	"testing"

	"github.com/stretchr/testify/require"

	"gophers.dev/pkgs/semantic"
)

func Test_githubCommit_Pseudo(t *testing.T) {
	var gc githubCommit
	gc.SHA = "c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc"
	gc.Commit.Author.Date = "2011-01-26T19:06:43Z"

	pseudo, err := gc.Pseudo([]semantic.Tag{}, true)
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

	_, err := gc.Pseudo([]semantic.Tag{}, true)
	require.Error(t, err)
}

func Test_githubCommit_Pseudo_previous_pre_version(t *testing.T) {
	var gc githubCommit
	gc.SHA = "c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc"
	gc.Commit.Author.Date = "2011-01-26T19:06:43Z"

	pseudo, err := gc.Pseudo([]semantic.Tag{
		{Major: 1, Minor: 2, Patch: 4, Extension: "pre"},
		{Major: 1, Minor: 2, Patch: 3},
	}, true)
	require.NoError(t, err)
	require.Equal(t, "v1.2.4-pre.0.20110126190643-c5b97d5ae6c1", pseudo)
}

func Test_githubCommit_Pseudo_previous_version(t *testing.T) {
	var gc githubCommit
	gc.SHA = "c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc"
	gc.Commit.Author.Date = "2011-01-26T19:06:43Z"

	pseudo, err := gc.Pseudo([]semantic.Tag{
		{Major: 1, Minor: 2, Patch: 4},
		{Major: 1, Minor: 2, Patch: 3},
	}, true)
	require.NoError(t, err)
	require.Equal(t, "v1.2.5-0.20110126190643-c5b97d5ae6c1", pseudo)
}

func Test_githubCommit_Pseudo_incompatible(t *testing.T) {
	var gc githubCommit
	gc.SHA = "c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc"
	gc.Commit.Author.Date = "2011-01-26T19:06:43Z"

	pseudo, err := gc.Pseudo([]semantic.Tag{
		{Major: 2, Minor: 2, Patch: 4},
		{Major: 1, Minor: 2, Patch: 3},
	}, false)
	require.NoError(t, err)
	require.Equal(t, "v2.2.5-0.20110126190643-c5b97d5ae6c1+incompatible", pseudo)
}

func Test_githubCommit_Pseudo_incompatible_semver1(t *testing.T) {
	var gc githubCommit
	gc.SHA = "c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc"
	gc.Commit.Author.Date = "2011-01-26T19:06:43Z"

	pseudo, err := gc.Pseudo([]semantic.Tag{
		{Major: 1, Minor: 2, Patch: 4},
		{Major: 1, Minor: 2, Patch: 3},
	}, false)
	require.NoError(t, err)
	require.Equal(t, "v1.2.5-0.20110126190643-c5b97d5ae6c1", pseudo)
}

// v3.0.0-rc.2 v3.0.0-rc.1 v2.2.0 2.2.0 2.1.0 2.0.0 1.1.0 1.0.0
