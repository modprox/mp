package finder

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func Test_githubCommit_Pseudo(t *testing.T) {
	var gc githubCommit
	gc.SHA = "c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc"
	gc.Commit.Author.Date = "2011-01-26T19:06:43Z"

	pseudo, err := gc.Pseudo()
	require.NoError(t, err)
	require.Equal(
		t,
		"v0.0.0-201101260700-c5b97d5ae6c1+incompatible",
		pseudo,
	)
}

func Test_githubCommit_Pseudo_bad_time(t *testing.T) {
	var gc githubCommit
	gc.SHA = "c5b97d5ae6c19d5c5df71a34c7fbeeda2479ccbc"
	gc.Commit.Author.Date = "2011-01-26T19:06:43" // no Z

	_, err := gc.Pseudo()
	require.Error(t, err)
}
