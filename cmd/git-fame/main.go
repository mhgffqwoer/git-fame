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

func main() {
	rootCmd := &cobra.Command{
		Use:   "gitfame",
		Short: "gitfame — анализ статистики авторов git-репозитория",
		Long:  `GitFame показывает, кто сколько строк кода написал в репозитории.`,
		RunE:  runRoot,
	}

	rootCmd.Flags().String("repository", ".", "Путь до Git-репозитория (по умолчанию текущая директория)")
	rootCmd.Flags().String("revision", "HEAD", "Указатель на коммит (по умолчанию HEAD)")
	rootCmd.Flags().String("order-by", "lines", "Порядок сортировки (lines|files|commits)")
	rootCmd.Flags().Bool("use-committer", false, "Использовать committer вместо author")
	rootCmd.Flags().String("format", "tabular", "Формат вывода (tabular|json|csv|json-lines)")
	rootCmd.Flags().String("extensions", "", "список расширений, сужающий список файлов в расчёте (например '.go,.md')")
	rootCmd.Flags().String("languages", "", "список языков, сужающий список файлов в расчёте (например 'go,markdown')")
	rootCmd.Flags().String("exclude", "", "набор Glob паттернов, исключающих файлы из расчёта (например 'foo/*,bar/*')")
	rootCmd.Flags().String("restrict-to", "", "набор Glob паттернов, исключающий все файлы, не удовлетворяющие ни одному из паттернов набора")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprint(os.Stderr, err)
		os.Exit(1)
	}
}

func runRoot(cmd *cobra.Command, args []string) error {
	ctx := context.Background()

	repoPath, _ := cmd.Flags().GetString("repository")
	commit, _ := cmd.Flags().GetString("revision")
	orderBy, _ := cmd.Flags().GetString("order-by")
	useCommitter, _ := cmd.Flags().GetBool("use-committer")
	formatType, _ := cmd.Flags().GetString("format")
	extensions, _ := cmd.Flags().GetString("extensions")
	languages, _ := cmd.Flags().GetString("languages")
	exclude, _ := cmd.Flags().GetString("exclude")
	restrictTo, _ := cmd.Flags().GetString("restrict-to")

	gitClient, err := git_client.NewBuilder().WithPath(repoPath).WithCommit(commit).Build()
	if err != nil {
		return err
	}
	f, err := filter.NewFileFilterBuilder().
		WithExtensions(extensions).
		WithLanguages(languages).
		WithExcludes(exclude).
		WithRestrictTo(restrictTo).
		Build()
	if err != nil {
		return err
	}
	analyzer := analysis_app.NewAnalyzer(
		git_adapters.NewFileStreamer(gitClient),
		git_adapters.NewFileHistoryStreamer(gitClient, useCommitter),
		f,
		storage.NewStorage(),
	)
	data, err := analyzer.Analyze(ctx)
	if err != nil {
		return err
	}

	s, err := sorter.NewSorter(orderBy)
	if err != nil {
		return err
	}
	s.Sort(data)
	fmtr, err := formatter.NewFormatter(formatType)
	if err != nil {
		return err
	}
	result, err := fmtr.Format(data)
	if err != nil {
		return err
	}
	fmt.Print(result)
	return nil
}
