package twitterscraper

import (
	"context"
	"net/url"
)

// GetFollowers returns channel with followers for a given user
func (s *Scraper) GetFollowers(ctx context.Context, userId string, maxProfilesNbr int) <-chan *ProfileResult {
	return getUserTimeline(ctx, userId, maxProfilesNbr, s.FetchUserFollowers)
}

// FetchUserFollowers gets followers for a given userID, via the Twitter frontend API
func (s *Scraper) FetchUserFollowers(userId string, maxFollowersNbr int, cursor string) ([]*Profile, string, error) {
	timeline, err := s.FetchUserFollowersByUserID(userId, maxFollowersNbr, cursor)
	if err != nil {
		return nil, "", err
	}
	profile, nextcursor := timeline.parseUsers()

	return profile, nextcursor, nil
}

// FetchUserFollowersByUserID gets followers for a given userID, via the Twitter frontend GraphQL API.
func (s *Scraper) FetchUserFollowersByUserID(userID string, maxNbr int, cursor string) (*timeline, error) {

	if maxNbr > 50 {
		maxNbr = 50
	}

	req, err := s.newRequest("GET", "https://twitter.com/i/api/graphql/3yX7xr2hKjcZYnXt6cU6lQ/Followers")
	if err != nil {
		return nil, err
	}

	variables := map[string]interface{}{
		"userId":                 userID,
		"count":                  maxNbr,
		"includePromotedContent": false,
	}

	features := map[string]interface{}{
		"rweb_lists_timeline_redesign_enabled":                                    true,
		"responsive_web_graphql_exclude_directive_enabled":                        true,
		"verified_phone_label_enabled":                                            false,
		"creator_subscriptions_tweet_preview_api_enabled":                         true,
		"responsive_web_graphql_timeline_navigation_enabled":                      true,
		"responsive_web_graphql_skip_user_profile_image_extensions_enabled":       false,
		"tweetypie_unmention_optimization_enabled":                                true,
		"responsive_web_edit_tweet_api_enabled":                                   true,
		"graphql_is_translatable_rweb_tweet_is_translatable_enabled":              true,
		"view_counts_everywhere_api_enabled":                                      true,
		"longform_notetweets_consumption_enabled":                                 true,
		"responsive_web_twitter_article_tweet_consumption_enabled":                false,
		"tweet_awards_web_tipping_enabled":                                        false,
		"freedom_of_speech_not_reach_fetch_enabled":                               true,
		"standardized_nudges_misinfo":                                             true,
		"tweet_with_visibility_results_prefer_gql_limited_actions_policy_enabled": true,
		"longform_notetweets_rich_text_read_enabled":                              true,
		"longform_notetweets_inline_media_enabled":                                true,
		"responsive_web_media_download_video_enabled":                             false,
		"responsive_web_enhance_cards_enabled":                                    false,
	}

	if cursor != "" {
		variables["cursor"] = cursor
	}

	query := url.Values{}
	query.Set("variables", mapToJSONString(variables))
	query.Set("features", mapToJSONString(features))
	req.URL.RawQuery = query.Encode()

	var timeline timeline
	err = s.RequestAPI(req, &timeline)
	if err != nil {
		return nil, err
	}
	return &timeline, nil
}
