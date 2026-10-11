package database

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// newAccessKey returns an unguessable, URL-safe key (96 random bits).
func newAccessKey() string {
	b := make([]byte, 12)
	rand.Read(b)
	return base64.RawURLEncoding.EncodeToString(b)
}

// SetVisibility makes a flipbook public, or link-only. Switching to link-only
// creates an access key if the flipbook has none; an existing key is kept so
// links already sent keep working.
func (d *DB) SetVisibility(id string, public bool) error {
	ctx := context.Background()
	set := bson.M{"is_public": public, "updated_at": time.Now()}
	if !public {
		fb, err := d.GetFlipbook(id)
		if err != nil {
			return err
		}
		if fb.AccessKey == "" {
			set["access_key"] = newAccessKey()
		}
	}
	_, err := d.flipbooks.UpdateByID(ctx, id, bson.M{"$set": set})
	return err
}

// RotateAccessKey issues a new access key, so every previously shared
// link-only URL stops working. It returns the new key.
func (d *DB) RotateAccessKey(id string) (string, error) {
	key := newAccessKey()
	_, err := d.flipbooks.UpdateByID(context.Background(), id,
		bson.M{"$set": bson.M{"access_key": key, "updated_at": time.Now()}})
	return key, err
}
