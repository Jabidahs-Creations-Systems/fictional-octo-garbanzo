package gitlen

import (
	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/spf13/cobra"
)

func NewCmdGitlen(f *cmdutil.Factory) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "gitlen <command>",
		Short: "Supercharge Git with enhanced insights",
		Long: heredoc.Doc(`
			Supercharge your Git workflow with enhanced visualization and insights,
			inspired by GitLens. View detailed commit information, blame annotations,
			and file history with rich context.
		`),
		Example: heredoc.Doc(`
			# Show enhanced blame for a file
			$ gh gitlen blame <file>

			# Show enhanced commit details
			$ gh gitlen show <commit>

			# Show file history with details
			$ gh gitlen history <file>
		`),
		GroupID: "core",
	}

	cmdutil.DisableAuthCheck(cmd)

	cmd.AddCommand(NewCmdBlame(f))
	cmd.AddCommand(NewCmdShow(f))
	cmd.AddCommand(NewCmdHistory(f))

	return cmd
}
