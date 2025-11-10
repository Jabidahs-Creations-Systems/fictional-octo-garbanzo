package gitlen

import (
	"bytes"
	"io"
	"testing"

	"github.com/cli/cli/v2/internal/config"
	"github.com/cli/cli/v2/internal/gh"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/google/shlex"
	"github.com/stretchr/testify/assert"
)

func TestNewCmdGitlen(t *testing.T) {
	tests := []struct {
		name     string
		cli      string
		wantsErr bool
	}{
		{
			name:     "no args",
			cli:      "",
			wantsErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ios, _, _, _ := iostreams.Test()

			f := &cmdutil.Factory{
				IOStreams: ios,
				Config: func() (gh.Config, error) {
					return config.NewBlankConfig(), nil
				},
			}

			cmd := NewCmdGitlen(f)
			argv, err := shlex.Split(tt.cli)
			assert.NoError(t, err)

			cmd.SetArgs(argv)
			cmd.SetIn(&bytes.Buffer{})
			cmd.SetOut(io.Discard)
			cmd.SetErr(io.Discard)

			_, err = cmd.ExecuteC()
			if tt.wantsErr {
				assert.Error(t, err)
			}
		})
	}
}

func TestNewCmdBlame(t *testing.T) {
	ios, _, _, _ := iostreams.Test()

	f := &cmdutil.Factory{
		IOStreams: ios,
		Config: func() (gh.Config, error) {
			return config.NewBlankConfig(), nil
		},
	}

	cmd := NewCmdBlame(f)
	assert.NotNil(t, cmd)
	assert.Equal(t, "blame <file>", cmd.Use)
}

func TestNewCmdShow(t *testing.T) {
	ios, _, _, _ := iostreams.Test()

	f := &cmdutil.Factory{
		IOStreams: ios,
		Config: func() (gh.Config, error) {
			return config.NewBlankConfig(), nil
		},
	}

	cmd := NewCmdShow(f)
	assert.NotNil(t, cmd)
	assert.Equal(t, "show [<commit>]", cmd.Use)
}

func TestNewCmdHistory(t *testing.T) {
	ios, _, _, _ := iostreams.Test()

	f := &cmdutil.Factory{
		IOStreams: ios,
		Config: func() (gh.Config, error) {
			return config.NewBlankConfig(), nil
		},
	}

	cmd := NewCmdHistory(f)
	assert.NotNil(t, cmd)
	assert.Equal(t, "history <file>", cmd.Use)
}
