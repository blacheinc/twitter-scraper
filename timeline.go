package twitterscraper

type results struct {
	Typename       string     `json:"__typename"`
	IsBlueVerified bool       `json:"is_blue_verified"`
	RestID         string     `json:"rest_id"`
	Legacy         legacyUser `json:"legacy"`
}

type entrys struct {
	EntryId string `json:"entryId"`
	Content struct {
		CursorType  string `json:"cursorType"`
		Value       string `json:"value"`
		ItemContent struct {
			UserDisplayType string `json:"userDisplayType"`
			UserResults     struct {
				Result results `json:"results"`
			} `json:"user_results"`
		} `json:"itemContent"`
	} `json:"content"`
}

// timeline JSON object
type timeline struct {
	Data struct {
		User struct {
			Result struct {
				Timeline struct {
					Timeline struct {
						Instructions []struct {
							Entries []entrys `json:"entries"`
							Type    string   `json:"type"`
							Entry   entry    `json:"entry,omitempty"`
						} `json:"instructions"`
					} `json:"timeline"`
				} `json:"timeline"`
			} `json:"result"`
		} `json:"user"`
	} `json:"data"`
}

// parseUsers parses the timeline data and extracts follower entryID and cursor value.
func (timeline *timeline) parseUsers() ([]*Profile, string) {
	followers := make([]*Profile, 0)
	cursor := ""

	for _, instruction := range timeline.Data.User.Result.Timeline.Timeline.Instructions {
		if instruction.Type == "TimelineAddEntries" {
			for _, entry := range instruction.Entries {
				if entry.Content.CursorType == "Bottom" {
					cursor = entry.Content.Value
					continue
				}
			}

			for _, entry := range instruction.Entries {
				if entry.Content.ItemContent.UserDisplayType == "User" {
					follower := parseProfiles(entry)
					if follower.UserID == "" {
						follower.UserID = entry.EntryId
					}
					followers = append(followers, &follower)

				} else if entry.Content.CursorType == "Bottom" {
					cursor = entry.Content.Value
				}

			}
		}
	}
	return followers, cursor
}
