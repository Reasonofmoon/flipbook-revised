package database

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/jonradoff/flipbook/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
)

// openTestDB connects to FLIPBOOK_TEST_MONGO_URI using a throwaway database.
// Tests are skipped when the variable is unset (e.g. plain `go test` in CI).
func openTestDB(t *testing.T) *DB {
	t.Helper()
	uri := os.Getenv("FLIPBOOK_TEST_MONGO_URI")
	if uri == "" {
		t.Skip("FLIPBOOK_TEST_MONGO_URI not set; skipping MongoDB integration test")
	}
	ctx := context.Background()
	name := fmt.Sprintf("flipbook_test_%d", time.Now().UnixNano())
	d, err := Open(ctx, uri, name)
	if err != nil {
		t.Fatalf("open: %v", err)
	}
	t.Cleanup(func() {
		d.client.Database(name).Drop(ctx)
		d.Close(ctx)
	})
	return d
}

func create(t *testing.T, d *DB, id, slug string) {
	t.Helper()
	if err := d.CreateFlipbook(&models.Flipbook{ID: id, Title: slug, Slug: slug, Status: models.StatusReady}); err != nil {
		t.Fatalf("create %s: %v", slug, err)
	}
}

func TestDeleteReleasesSlug(t *testing.T) {
	d := openTestDB(t)
	create(t, d, "id-1", "수능-영어-독해")

	if got := d.EnsureUniqueSlug("수능-영어-독해"); got != "수능-영어-독해-2" {
		t.Fatalf("while live, EnsureUniqueSlug = %q, want suffix -2", got)
	}
	if err := d.DeleteFlipbook("id-1"); err != nil {
		t.Fatal(err)
	}
	if got := d.EnsureUniqueSlug("수능-영어-독해"); got != "수능-영어-독해" {
		t.Fatalf("after delete, EnsureUniqueSlug = %q, want the original slug back", got)
	}
	// The unique index must accept a new flipbook with the same slug.
	create(t, d, "id-2", "수능-영어-독해")
	fb, err := d.GetFlipbookBySlug("수능-영어-독해")
	if err != nil || fb.ID != "id-2" {
		t.Fatalf("GetFlipbookBySlug = %+v, %v; want the new flipbook id-2", fb, err)
	}
}

func TestDeleteTwiceKeepsSingleSuffix(t *testing.T) {
	d := openTestDB(t)
	create(t, d, "id-1", "unit-3")
	d.DeleteFlipbook("id-1")
	d.DeleteFlipbook("id-1") // second call must not append the marker again

	var doc flipbookDoc
	if err := d.flipbooks.FindOne(context.Background(), bson.M{"_id": "id-1"}).Decode(&doc); err != nil {
		t.Fatal(err)
	}
	if want := "unit-3" + DeletedSlugMarker + "id-1"; doc.Slug != want {
		t.Fatalf("slug = %q, want %q", doc.Slug, want)
	}
}

func TestReleaseDeletedSlugsMigratesOldRecords(t *testing.T) {
	d := openTestDB(t)
	ctx := context.Background()
	create(t, d, "old", "legacy-deck")
	// Simulate a record deleted by the old code: deleted_at set, slug untouched.
	if _, err := d.flipbooks.UpdateByID(ctx, "old", bson.M{"$set": bson.M{"deleted_at": time.Now()}}); err != nil {
		t.Fatal(err)
	}

	d.releaseDeletedSlugs(ctx)
	d.releaseDeletedSlugs(ctx) // idempotent

	var doc flipbookDoc
	if err := d.flipbooks.FindOne(ctx, bson.M{"_id": "old"}).Decode(&doc); err != nil {
		t.Fatal(err)
	}
	if strings.Count(doc.Slug, DeletedSlugMarker) != 1 || !strings.HasPrefix(doc.Slug, "legacy-deck") {
		t.Fatalf("migrated slug = %q", doc.Slug)
	}
	if got := d.EnsureUniqueSlug("legacy-deck"); got != "legacy-deck" {
		t.Fatalf("EnsureUniqueSlug after migration = %q, want legacy-deck", got)
	}
}
