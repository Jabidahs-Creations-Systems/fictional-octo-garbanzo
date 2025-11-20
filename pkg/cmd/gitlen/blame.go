package gitlen

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/MakeNowJust/heredoc"
	"github.com/cli/cli/v2/git"
	"github.com/cli/cli/v2/pkg/cmdutil"
	"github.com/cli/cli/v2/pkg/iostreams"
	"github.com/spf13/cobra"
)

type BlameOptions struct {
	IO        *iostreams.IOStreams
	GitClient *git.Client
	FilePath  string
}

func NewCmdBlame(f *cmdutil.Factory) *cobra.Command {
	opts := &BlameOptions{
		IO:        f.IOStreams,
		GitClient: f.GitClient,
	}

	cmd := &cobra.Command{
		Use:   "blame <file>",
		Short: "Show enhanced blame with commit details",
		Long: heredoc.Doc(`
			Display git blame for a file with enhanced formatting and commit information.
			Shows who last modified each line, when, and the commit message.
		`),
		Example: heredoc.Doc(`
			# Show blame for a file
			$ gh gitlen blame README.md

			# Show blame for a file with line numbers
			$ gh gitlen blame src/main.go
		`),
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			opts.FilePath = args[0]
			return blameRun(opts)
		},
	}

	return cmd
}

func blameRun(opts *BlameOptions) error {
	ctx := context.Background()
	cs := opts.IO.ColorScheme()

	// Get git blame output
	blameCmd, err := opts.GitClient.Command(ctx, "blame", "--line-porcelain", opts.FilePath)
	if err != nil {
		return fmt.Errorf("failed to create git blame command: %w", err)
	}
	blameOutput, err := blameCmd.Output()
	if err != nil {
		return fmt.Errorf("failed to run git blame: %w", err)
	}

	// Parse and display blame output
	lines := strings.Split(string(blameOutput), "\n")
	var currentCommit, currentAuthor, currentDate string
	lineNumber := 1

	for i := 0; i < len(lines); i++ {
		line := lines[i]
		
		if len(line) == 0 {
			continue
		}

		// Parse commit hash (first line of each blame block)
		if len(line) > 40 && line[40] == ' ' {
			currentCommit = line[:8] // Short hash
		} else if strings.HasPrefix(line, "author ") {
			currentAuthor = strings.TrimPrefix(line, "author ")
		} else if strings.HasPrefix(line, "author-time ") {
			timeStr := strings.TrimPrefix(line, "author-time ")
			if timestamp, err := time.Parse("1136239445", timeStr); err == nil {
				currentDate = timestamp.Format("2006-01-02")
			}
		} else if strings.HasPrefix(line, "\t") {
			// This is the actual code line
			codeLine := strings.TrimPrefix(line, "\t")
			
			// Format and display the enhanced blame line
			fmt.Fprintf(opts.IO.Out, "%s %s %s %s | %s\n",
				cs.Gray(fmt.Sprintf("%4d", lineNumber)),
				cs.Yellow(currentCommit),
				cs.Cyan(truncate(currentAuthor, 15)),
				cs.Gray(currentDate),
				codeLine,
			)
			
			lineNumber++
		}
	}

	return nil
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return fmt.Sprintf("%-*s", maxLen, s)
	}
	return s[:maxLen-1] + "…"
}
