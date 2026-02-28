package main

import (
	"context"
	_ "embed"
	"fmt"
	"os"

	analysis_app "github.com/mhgffqwoer/git-fame/internal/analyzer"
	"github.com/mhgffqwoer/git-fame/internal/analyzer/filter"
	git_adapters "github.com/mhgffqwoer/git-fame/internal/analyzer/git_adapters"
	"github.com/mhgffqwoer/git-fame/internal/analyzer/storage"
	git_client "github.com/mhgffqwoer/git-fame/internal/git_client"
	formatter "github.com/mhgffqwoer/git-fame/internal/output/formatter"
	"github.com/mhgffqwoer/git-fame/internal/output/sorter"
	"github.com/spf13/cobra"
)

// main initializes the CLI application, sets up the available flags,
// and executes the root command.
func main() {
	rootCmd := &cobra.Command{
		Use:   "git-fame",
		Short: "git-fame — git repository authorship statistics tool",
		Long:  `Git-Fame calculates and displays authorship statistics for a given git repository.`,
		RunE:  runRoot,
	}

	// Flag definitions for customizing the analysis and output.
	rootCmd.Flags().String("repository", ".", "path to the git repository (default is current directory)")
	rootCmd.Flags().String("revision", "HEAD", "revision pointer (default is HEAD)")
	rootCmd.Flags().String("order-by", "lines", "sort order (lines|files|commits)")
	rootCmd.Flags().Bool("use-committer", false, "use committer instead of author")
	rootCmd.Flags().String("format", "tabular", "output format (tabular|json|csv|json-lines)")
	rootCmd.Flags().String("extensions", "", "comma-separated list of extensions to include (e.g., '.go,.md')")
	rootCmd.Flags().String("languages", "", "comma-separated list of languages to include (e.g., 'go,markdown')")
	rootCmd.Flags().String("exclude", "", "set of glob patterns to exclude files (e.g., 'foo/*,bar/*')")
	rootCmd.Flags().String("restrict-to", "", "set of glob patterns; exclude all files that do not match at least one pattern")

	// Execute the command and handle fatal startup errors.
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

// runRoot implements the core logic for the root command.
// It orchestrates the entire flow: parsing flags, initializing the git client and filters,
// running the analyzer, and finally sorting and formatting the output.
func runRoot(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	// Extract flag values.
	repoPath, _ := cmd.Flags().GetString("repository")
	commit, _ := cmd.Flags().GetString("revision")
	orderBy, _ := cmd.Flags().GetString("order-by")
	useCommitter, _ := cmd.Flags().GetBool("use-committer")
	formatType, _ := cmd.Flags().GetString("format")
	extensions, _ := cmd.Flags().GetString("extensions")
	languages, _ := cmd.Flags().GetString("languages")
	exclude, _ := cmd.Flags().GetString("exclude")
	restrictTo, _ := cmd.Flags().GetString("restrict-to")

	// Initialize the Git client to interact with the repository at a specific revision.
	gitClient, err := git_client.NewBuilder().WithPath(repoPath).WithCommit(commit).Build()
	if err != nil {
		return err
	}

	// Build the file filter based on extensions, languages, and glob patterns.
	f, err := filter.NewFileFilterBuilder().
		WithExtensions(extensions).
		WithLanguages(languages).
		WithExcludes(exclude).
		WithRestrictTo(restrictTo).
		Build()
	if err != nil {
		return err
	}

	// Set up the Analyzer with necessary dependencies for data streaming and storage.
	analyzer := analysis_app.NewAnalyzer(
		git_adapters.NewFileStreamer(gitClient),
		git_adapters.NewFileHistoryStreamer(gitClient, useCommitter),
		f,
		storage.NewStorage(),
	)

	// Perform the analysis (running git blame on filtered files).
	data, err := analyzer.Analyze(ctx)
	if err != nil {
		return err
	}

	// Initialize the sorter and order the results based on the chosen metric.
	s, err := sorter.NewSorter(orderBy)
	if err != nil {
		return err
	}
	s.Sort(data)

	// Initialize the formatter and generate the final output string.
	fmtr, err := formatter.NewFormatter(formatType)
	if err != nil {
		return err
	}
	result, err := fmtr.Format(data)
	if err != nil {
		return err
	}

	// Print the final report to standard output.
	fmt.Print(result)
	return nil
}
