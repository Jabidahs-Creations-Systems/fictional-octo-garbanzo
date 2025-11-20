package gitlen

import (
	"context"
	"fmt"
	"strings"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/git"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type HistoryOptions struct {
	IO        *iostreams.IOStreams
	GitClient *git.Client
	FilePath  string
	MaxCount  int
}

func NewCmdHistory(f *cmdutil.Factory) *cobra.Command {
	opts := &HistoryOptions{
		IO:        f.IOStreams,
		GitClient: f.GitClient,
		MaxCount:  10,
	}

	cmd := &cobra.Command{
		Use:   "history <file>",
		Short: "Show enhanced file history",
		Long: heredoc.Doc(`
			Display the commit history for a file with enhanced formatting.
			Shows all commits that modified the specified file.
		`),
		Example: heredoc.Doc(`
			# Show history for a file
			$ gh gitlen history README.md

			# Show more commits
			$ gh gitlen history src/main.go --limit 20
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.FilePath = args[0]
			return historyRun(opts)
		},
	}

	cmd.Flags().IntVarP(&opts.MaxCount, "limit", "n", 10, "Maximum number of commits to show")

	return cmd
}

func historyRun(opts *HistoryOptions) error {
	ctx := context.Background()
	cs := opts.IO.ColorScheme()

	// Get file history
	logCmd, err := opts.GitClient.Command(ctx, "log",
		fmt.Sprintf("-%d", opts.MaxCount),
		"--follow",
		"--format=%H|%an|%ae|%ad|%s",
		"--date=short",
		"--",
		opts.FilePath)
	if err != nil {
		return fmt.Errorf("failed to create log command: %w", err)
	}
	
	logOutput, err := logCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get file history: %w", err)
	}

	if len(logOutput) == 0 {
		fmt.Fprintf(opts.IO.Out, "No history found for %s\n", opts.FilePath)
		return nil
	}

	fmt.Fprintf(opts.IO.Out, "%s %s\n\n", cs.Bold("History for:"), cs.Cyan(opts.FilePath))

	lines := strings.Split(string(logOutput), "\n")
	for _, line := range lines {
		if line == "" {
			continue
		}

		parts := strings.Split(line, "|")
		if len(parts) < 5 {
			continue
		}

		hash := parts[0][:8]
		author := parts[1]
		date := parts[3]
		message := parts[4]

		// Get stats for this commit
		statCmd, err := opts.GitClient.Command(ctx, "diff", "--shortstat",
			fmt.Sprintf("%s^", parts[0]),
			parts[0],
			"--",
			opts.FilePath)
		if err == nil {
			statOutput, _ := statCmd.Output()
			stats := strings.TrimSpace(string(statOutput))

			fmt.Fprintf(opts.IO.Out, "%s %s %s %s\n",
				cs.Yellow(hash),
				cs.Cyan(truncate(author, 20)),
				cs.Gray(date),
				message,
			)
			
			if stats != "" {
				fmt.Fprintf(opts.IO.Out, "  %s\n", cs.Gray(stats))
			}
			fmt.Fprintln(opts.IO.Out)
		}
	}

	return nil
}
