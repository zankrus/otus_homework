package main

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	tests := []struct {
		name         string
		command      []string
		expectedCode int
	}{
		{
			name:         "Успешное выполнение",
			command:      []string{"ls"},
			expectedCode: 0,
		},
		{
			name:         "Ошибка",
			command:      []string{"ls", "ls"},
			expectedCode: 1,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			x := RunCmd(test.command, nil)
			require.Equal(t, test.expectedCode, x)
		})
	}
}

func TestPrepareEnvs(t *testing.T) {
	t.Run("no deletions", func(t *testing.T) {
		ed := Environment{}
		envs := []string{"user=Ro"}

		res := prepareEnvs(envs, ed)

		require.Len(t, res, 1)
		require.Equal(t, []string{"user=Ro"}, res)

	})

	t.Run("c заменой", func(t *testing.T) {
		ed := Environment{}
		ed["user"] = EnvValue{
			Value:      "clone",
			NeedRemove: false,
		}
		envs := []string{"user=Ro"}

		res := prepareEnvs(envs, ed)

		require.Len(t, res, 1)
		require.Equal(t, []string{"user=clone"}, res)

	})

	t.Run("c удалением", func(t *testing.T) {
		ed := Environment{}
		ed["user"] = EnvValue{
			Value:      "clone",
			NeedRemove: true,
		}
		envs := []string{"user=Ro"}

		res := prepareEnvs(envs, ed)

		require.Len(t, res, 0)
		require.Equal(t, []string{}, res)

	})
}
