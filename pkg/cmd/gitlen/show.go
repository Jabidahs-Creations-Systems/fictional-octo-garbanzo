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

type ShowOptions struct {
	IO        *iostreams.IOStreams
	GitClient *git.Client
	CommitSHA string
	ShowStats bool
}

func NewCmdShow(f *cmdutil.Factory) *cobra.Command {
	opts := &ShowOptions{
		IO:        f.IOStreams,
		GitClient: f.GitClient,
	}

	cmd := &cobra.Command{
		Use:   "show [<commit>]",
		Short: "Show enhanced commit details",
		Long: heredoc.Doc(`
			Display detailed information about a commit with enhanced formatting.
			Shows commit metadata, changes, and statistics.
		`),
		Example: heredoc.Doc(`
			# Show details for the latest commit
			$ gh gitlen show

			# Show details for a specific commit
			$ gh gitlen show abc123

			# Show commit with statistics
			$ gh gitlen show --stats
		`),
		Args: cobra.MaximumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.CommitSHA = args[0]
			} else {
				opts.CommitSHA = "HEAD"
			}
			return showRun(opts)
		},
	}

	cmd.Flags().BoolVar(&opts.ShowStats, "stats", false, "Show commit statistics")

	return cmd
}

func showRun(opts *ShowOptions) error {
	ctx := context.Background()
	cs := opts.IO.ColorScheme()

	// Get commit details
	showArgs := []string{"show", "--format=fuller", "--no-patch"}
	if opts.ShowStats {
		showArgs = append(showArgs, "--stat")
	}
	showArgs = append(showArgs, opts.CommitSHA)

	showCmd, err := opts.GitClient.Command(ctx, showArgs...)
	if err != nil {
		return fmt.Errorf("failed to create show command: %w", err)
	}
	showOutput, err := showCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to show commit: %w", err)
	}

	// Parse and format output
	lines := strings.Split(string(showOutput), "\n")
	for _, line := range lines {
		switch {
		case strings.HasPrefix(line, "commit "):
			hash := strings.TrimPrefix(line, "commit ")
			fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Commit:"), cs.Yellow(hash))
		case strings.HasPrefix(line, "Author: "):
			author := strings.TrimPrefix(line, "Author: ")
			fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Author:"), cs.Cyan(author))
		case strings.HasPrefix(line, "AuthorDate: "):
			date := strings.TrimPrefix(line, "AuthorDate: ")
			fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Date:  "), date)
		case strings.HasPrefix(line, "Commit: "):
			committer := strings.TrimPrefix(line, "Commit: ")
			fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("Committer:"), cs.Cyan(committer))
		case strings.HasPrefix(line, "CommitDate: "):
			date := strings.TrimPrefix(line, "CommitDate: ")
			fmt.Fprintf(opts.IO.Out, "%s %s\n", cs.Bold("CommitDate:"), date)
		case strings.HasPrefix(line, "    "):
			// Commit message lines
			fmt.Fprintf(opts.IO.Out, "%s\n", line)
		case line == "":
			fmt.Fprintln(opts.IO.Out)
		default:
			// Stats lines
			fmt.Fprintf(opts.IO.Out, "%s\n", line)
		}
	}

	// Get file changes
	diffCmd, err := opts.GitClient.Command(ctx, "diff-tree", "--no-commit-id", "--name-status", "-r", opts.CommitSHA)
	if err != nil {
		return fmt.Errorf("failed to create diff command: %w", err)
	}
	diffOutput, err := diffCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to get file changes: %w", err)
	}

	if len(diffOutput) > 0 {
		fmt.Fprintf(opts.IO.Out, "\n%s\n", cs.Bold("Files changed:"))
		lines := strings.Split(string(diffOutput), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			parts := strings.Fields(line)
			if len(parts) >= 2 {
				status := parts[0]
				file := parts[1]
				statusColor := cs.Green
				statusText := "Added"
				switch status {
				case "M":
					statusColor = cs.Yellow
					statusText = "Modified"
				case "D":
					statusColor = cs.Red
					statusText = "Deleted"
				case "R":
					statusColor = cs.Blue
					statusText = "Renamed"
				case "A":
					statusColor = cs.Green
					statusText = "Added"
				}
				fmt.Fprintf(opts.IO.Out, "  %s %s\n", statusColor(statusText), file)
			}
		}
	}

	return nil
}
