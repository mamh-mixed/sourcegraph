package upgradestore

import (
	"context"
	"testing"

	"github.com/Masterminds/semver"
	"github.com/google/go-cmp/cmp"

	"github.com/sourcegraph/log/logtest"

	"github.com/sourcegraph/sourcegraph/internal/database"
	"github.com/sourcegraph/sourcegraph/internal/database/dbtest"
	"github.com/sourcegraph/sourcegraph/lib/errors"
)

func TestGetFirstServiceVersion(t *testing.T) {
	ctx := context.Background()
	logger := logtest.Scoped(t)
	db := database.NewDB(logger, dbtest.NewDB(logger, t))
	store := New(db, nil)

	if err := store.UpdateServiceVersion(ctx, "service", "1.2.3"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if err := store.UpdateServiceVersion(ctx, "service", "1.2.4"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if err := store.UpdateServiceVersion(ctx, "service", "1.3.0"); err != nil {
		t.Fatalf("unexpected error: %s", err)
	}

	firstVersion, err := store.GetFirstServiceVersion(ctx, "service")
	if err != nil {
		t.Fatalf("unexpected error: %s", err)
	}
	if firstVersion != "1.2.3" {
		t.Errorf("unexpected first version. want=%s have=%s", "1.2.3", firstVersion)
	}
}

func TestUpdateServiceVersion(t *testing.T) {
	ctx := context.Background()
	logger := logtest.Scoped(t)
	db := database.NewDB(logger, dbtest.NewDB(logger, t))
	store := New(db, nil)

	for _, tc := range []struct {
		version string
		err     error
	}{
		{"0.0.0", nil},
		{"0.0.1", nil},
		{"0.1.0", nil},
		{"0.2.0", nil},
		{"1.0.0", nil},
		{"1.2.0", &UpgradeError{
			Service:  "service",
			Previous: semver.MustParse("1.0.0"),
			Latest:   semver.MustParse("1.2.0"),
		}},
		{"2.1.0", &UpgradeError{
			Service:  "service",
			Previous: semver.MustParse("1.0.0"),
			Latest:   semver.MustParse("2.1.0"),
		}},
		{"0.3.0", nil}, // rollback
		{"non-semantic-version-is-always-valid", nil},
		{"1.0.0", nil}, // back to semantic version is allowed
		{"2.1.0", &UpgradeError{
			Service:  "service",
			Previous: semver.MustParse("1.0.0"),
			Latest:   semver.MustParse("2.1.0"),
		}}, // upgrade policy violation returns
	} {
		have := store.UpdateServiceVersion(ctx, "service", tc.version)
		want := tc.err

		if !errors.Is(have, want) {
			t.Fatal(cmp.Diff(have, want))
		}

		t.Logf("version = %q", tc.version)
	}
}
