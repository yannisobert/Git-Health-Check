package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/pterm/pterm"
	"github.com/spf13/cobra"

	"github.com/yannisobert/git-health-check/internal/analyzer"
	githubclient "github.com/yannisobert/git-health-check/internal/github"
)

func newCheckCmd() *cobra.Command {
	var (
		jsonOut bool
		noColor bool
		period  string
	)

	cmd := &cobra.Command{
		Use:   "check <owner/repo>",
		Short: "Analyze the health of a GitHub repository",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			owner, repo, err := parseRepo(args[0])
			if err != nil {
				return err
			}

			if noColor {
				pterm.DisableStyling()
			}

			client := githubclient.NewClient()
			a := analyzer.New(client)

			spinner, _ := pterm.DefaultSpinner.Start(fmt.Sprintf("Analyzing %s/%s…", owner, repo))
			report, err := a.Analyze(context.Background(), owner, repo)
			if err != nil {
				spinner.Fail(err.Error())
				return err
			}
			spinner.Success("Done")

			if jsonOut {
				enc := json.NewEncoder(os.Stdout)
				enc.SetIndent("", "  ")
				return enc.Encode(report)
			}

			printReport(report, period)
			return nil
		},
	}

	cmd.Flags().BoolVar(&jsonOut, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&noColor, "no-color", false, "Disable colored output")
	cmd.Flags().StringVar(&period, "period", "weekly", "History period for trend display (weekly|monthly)")

	return cmd
}

func parseRepo(input string) (owner, repo string, err error) {
	input = strings.TrimPrefix(input, "https://github.com/")
	input = strings.TrimPrefix(input, "github.com/")
	input = strings.TrimSuffix(input, "/")
	parts := strings.SplitN(input, "/", 2)
	if len(parts) != 2 || parts[0] == "" || parts[1] == "" {
		return "", "", fmt.Errorf("invalid repo: use owner/repo or https://github.com/owner/repo")
	}
	return parts[0], parts[1], nil
}

func printReport(r *analyzer.Report, _ string) {
	pterm.Println()

	scoreColor := pterm.FgRed
	switch {
	case r.Score >= 70:
		scoreColor = pterm.FgGreen
	case r.Score >= 40:
		scoreColor = pterm.FgYellow
	}

	pterm.DefaultHeader.WithFullWidth().WithBackgroundStyle(pterm.NewStyle(pterm.BgDarkGray)).
		Printf("%s   %s", r.FullName, scoreColor.Sprintf("%d / %d", r.Score, r.MaxScore))
	pterm.Println()

	var suggestions []string

	for _, cat := range r.Categories {
		catRatio := float64(cat.Score) / float64(cat.MaxScore)
		catColor := pterm.FgRed
		switch {
		case catRatio >= 0.7:
			catColor = pterm.FgGreen
		case catRatio >= 0.4:
			catColor = pterm.FgYellow
		}

		pterm.DefaultSection.Printf("%s  %s", cat.Name, catColor.Sprintf("%d/%d", cat.Score, cat.MaxScore))

		for _, c := range cat.Checks {
			icon := pterm.FgGreen.Sprint("✓")
			nameStr := c.Name
			switch c.Status {
			case analyzer.StatusFail:
				icon = pterm.FgRed.Sprint("✗")
				nameStr = pterm.FgRed.Sprint(c.Name)
			case analyzer.StatusWarn:
				icon = pterm.FgYellow.Sprint("~")
				nameStr = pterm.FgYellow.Sprint(c.Name)
			}

			scoreStr := fmt.Sprintf("%d/%d", c.Score, c.MaxScore)
			pterm.Printf("  %s  %-30s %5s   %s\n",
				icon, nameStr, scoreStr, pterm.FgGray.Sprint(c.Detail))

			if c.Suggestion != "" {
				pterm.Printf("       %s %s\n",
					pterm.FgYellow.Sprint("→"), c.Suggestion)
				suggestions = append(suggestions, c.Suggestion)
			}
		}
		pterm.Println()
	}

	if len(suggestions) > 0 {
		pterm.DefaultSection.Println("Top improvements")
		max := 5
		if len(suggestions) < max {
			max = len(suggestions)
		}
		for _, s := range suggestions[:max] {
			pterm.Printf("  %s %s\n", pterm.FgCyan.Sprint("•"), s)
		}
		pterm.Println()
	}
}
