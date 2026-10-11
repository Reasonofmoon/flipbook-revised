package database

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// ViewRecord is one recorded flipbook view, as used by reports.
type ViewRecord struct {
	FlipbookID string    `bson:"flipbook_id"`
	ViewedAt   time.Time `bson:"viewed_at"`
	Referrer   string    `bson:"referrer"`
	UserAgent  string    `bson:"user_agent"`
}

// ListViews returns views in [from, to). A nil flipbookIDs means all
// flipbooks; an empty (non-nil) slice matches nothing.
func (d *DB) ListViews(ctx context.Context, flipbookIDs []string, from, to time.Time) ([]ViewRecord, error) {
	if flipbookIDs != nil && len(flipbookIDs) == 0 {
		return nil, nil
	}
	filter := bson.M{"viewed_at": bson.M{"$gte": from, "$lt": to}}
	if flipbookIDs != nil {
		filter["flipbook_id"] = bson.M{"$in": flipbookIDs}
	}
	opts := options.Find().SetProjection(bson.M{"_id": 0, "flipbook_id": 1, "viewed_at": 1, "referrer": 1, "user_agent": 1})
	cur, err := d.views.Find(ctx, filter, opts)
	if err != nil {
		return nil, err
	}
	var out []ViewRecord
	if err := cur.All(ctx, &out); err != nil {
		return nil, err
	}
	return out, nil
}

// StorageUsage reports the MongoDB database size in bytes (data + indexes).
func (d *DB) StorageUsage(ctx context.Context) (int64, error) {
	var res struct {
		StorageSize float64 `bson:"storageSize"`
		IndexSize   float64 `bson:"indexSize"`
	}
	if err := d.flipbooks.Database().RunCommand(ctx, bson.D{{Key: "dbStats", Value: 1}}).Decode(&res); err != nil {
		return 0, err
	}
	return int64(res.StorageSize + res.IndexSize), nil
}
