package storage

import (
	"errors"
	"time"
)

type FeedItem struct {
	ID      int64
	FeedID  int64
	Created time.Time
	Updated time.Time
	Title   string
	Date    time.Time
	URL     string
	Content string
}

// Validate is used to assert Title and URL are set
func (item *FeedItem) Validate() error {
	if item.URL == "" {
		return errors.New("Missing FeedItem.URL")
	}

	if item.Title == "" {
		return errors.New("Missing FeedItem.Title")
	}

	if item.FeedID == 0 {
		return errors.New("Missing FeedItem.FeedID")
	}

	return nil
}

// ToBookmark converts the feed item to a bookmark
func (item *FeedItem) ToBookmark() *Bookmark {
	return &Bookmark{
		Title:   item.Title,
		URL:     item.URL,
		Content: item.Content,
	}
}

type ListFeedItemsOptions struct {
	Search   string
	Category string
	FeedID   string
	Limit    int
	Offset   int
}

// ListFeedItems fetches multiple feeds from the database
func (store *Store) ListFeedItems(options *ListFeedItemsOptions) (*[]*FeedItem, int) {
	query := store.db.Select("items i")

	if options.Search != "" {
		query.Where("(i.title LIKE ? OR i.url LIKE ? OR i.content LIKE ?)", "%"+options.Search+"%", "%"+options.Search+"%", "%"+options.Search+"%")
	}

	if options.Category != "" {
		query.Join("LEFT JOIN feeds f ON i.feed_id = f.id")
		query.Where("f.category = ?", options.Category)
	}

	if options.FeedID != "" {
		query.Where("i.feed_id = ?", options.FeedID)
	}

	feedItems := []*FeedItem{}
	totalCount := 0

	query.Columns("COUNT(i.id)")
	query.LoadValue(&totalCount)

	query.Columns("i.*")
	query.OrderBy("i.date", "DESC")
	query.Limit(options.Limit)
	query.Offset(options.Offset)
	query.Load(&feedItems)

	return &feedItems, totalCount
}

// GetFeedItem finds a single feed by ID or URL
func (store *Store) GetFeedItem(item *FeedItem) error {
	query := store.db.Select("items")
	query.Limit(1)

	if item.ID != 0 {
		query.Where("id = ?", item.ID)
	} else if item.URL != "" {
		query.Where("url = ?", item.URL)
	} else {
		return errors.New("Missing FeedItem.ID or FeedItem.URL")
	}

	if err := query.LoadValue(item); err != nil {
		return err
	}

	return nil
}

// DeleteFeedItem deletes the given feed item from the database
func (store *Store) DeleteFeedItem(item *FeedItem) error {
	if item.ID == 0 {
		return errors.New("Not an existing feed item")
	}

	query := store.db.Delete("items")
	query.Where("id = ?", item.ID)

	if _, err := query.Exec(); err != nil {
		return err
	}

	return nil
}
