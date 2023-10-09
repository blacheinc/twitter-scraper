package twitterscraper_test

import (
	"context"
	"testing"
	twitterscraper "github.com/blacheinc/twitter-scraper"
)

func TestGetUserFollowers(t *testing.T) {
	count := 0
	maxProfilesNbr := 4
	dupcheck := make(map[string]bool)
	scraper := twitterscraper.New()
	if err := scraper.Login("funmiloge4194", "funmi4194"); err != nil {
		t.Fatal(err)
	}

	for follower := range scraper.GetFollowers(context.Background(), "1008617870643875840", maxProfilesNbr) {
		if follower.Error != nil {
			t.Error(follower.Error)
		} else {
			count++
			if follower.UserID == "" {
				t.Error("Expected UserID is empty")
			} else {
				if dupcheck[follower.UserID] {
					t.Errorf("Detect duplicated UserID: %s", follower.UserID)
				} else {
					dupcheck[follower.UserID] = true
				}
			}
		}
	}

	if count != maxProfilesNbr {
		t.Errorf("Expected profiles count=%v, got: %v", maxProfilesNbr, count)
	}
}
