package adminanalytics

import (
	"context"

	"github.com/sourcegraph/sourcegraph/internal/database"
)

type Search struct {
	DateRange string
	DB        database.DB
}

func (s *Search) Searches() (*AnalyticsFetcher, error) {
	nodesQuery, summaryQuery, err := makeEventLogsQueries(s.DateRange, []string{"SearchResultsQueried"})
	if err != nil {
		return nil, err
	}

	return &AnalyticsFetcher{
		db:           s.DB,
		dateRange:    s.DateRange,
		nodesQuery:   nodesQuery,
		summaryQuery: summaryQuery,
		group:        "Search:Searches",
	}, nil
}

func (s *Search) FileViews() (*AnalyticsFetcher, error) {
	nodesQuery, summaryQuery, err := makeEventLogsQueries(s.DateRange, []string{"ViewBlob"})
	if err != nil {
		return nil, err
	}

	return &AnalyticsFetcher{
		db:           s.DB,
		dateRange:    s.DateRange,
		nodesQuery:   nodesQuery,
		summaryQuery: summaryQuery,
		group:        "Search:FileViews",
	}, nil
}

func (s *Search) FileOpens() (*AnalyticsFetcher, error) {
	nodesQuery, summaryQuery, err := makeEventLogsQueries(s.DateRange, []string{
		"GoToCodeHostClicked",
		"vscode.open.file",
		"openInAtom.open.file",
		"openineditor.open.file",
		"openInWebstorm.open.file",
	})
	if err != nil {
		return nil, err
	}

	return &AnalyticsFetcher{
		db:           s.DB,
		dateRange:    s.DateRange,
		nodesQuery:   nodesQuery,
		summaryQuery: summaryQuery,
		group:        "Search:FileOpens",
	}, nil
}

func (s *Search) CacheAll(ctx context.Context) error {
	fetcherBuilders := []func() (*AnalyticsFetcher, error){s.Searches, s.FileViews, s.FileOpens}
	for _, buildFetcher := range fetcherBuilders {
		fetcher, err := buildFetcher()
		if err != nil {
			return err
		}

		if _, err := fetcher.GetNodes(ctx, false); err != nil {
			return err
		}

		if _, err := fetcher.GetSummary(ctx, false); err != nil {
			return err
		}
	}
	return nil
}
