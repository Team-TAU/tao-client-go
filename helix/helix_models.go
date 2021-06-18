package helix

import "time"

// TwitchPagination represents pagination data from twitch on endpoints that support multi-paged responses.
type TwitchPagination struct {
	Cursor string `json:"cursor"`
}

// ExtensionAnalytics represents the response from Get Extension Analytics, see https://dev.twitch.tv/docs/api/reference#get-extension-analytics
type ExtensionAnalytics struct {
	Data []struct {
		ExtensionID string `json:"extension_id"`
		URL         string `json:"URL"`
		Type        string `json:"type"`
		DateRange   struct {
			StartedAt time.Time `json:"started_at"`
			EndedAt   time.Time `json:"ended_at"`
		} `json:"date_range"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// GameAnalytics represents the response from Get Game Analytics, see https://dev.twitch.tv/docs/api/reference#get-game-analytics
type GameAnalytics struct {
	Data []struct {
		GameID    string `json:"game_id"`
		URL       string `json:"URL"`
		Type      string `json:"type"`
		DateRange struct {
			StartedAt time.Time `json:"started_at"`
			EndedAt   time.Time `json:"ended_at"`
		} `json:"date_range"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// BitsLeaderboard represents the response from Get Bits Leaderboard, see https://dev.twitch.tv/docs/api/reference#get-bits-leaderboard
type BitsLeaderboard struct {
	Data []struct {
		UserID    string `json:"user_id"`
		UserLogin string `json:"user_login"`
		UserName  string `json:"user_name"`
		Rank      int    `json:"rank"`
		Score     int    `json:"score"`
	} `json:"data"`
	DateRange struct {
		StartedAt time.Time `json:"started_at"`
		EndedAt   time.Time `json:"ended_at"`
	} `json:"date_range"`
	Total int `json:"total"`
}

// CheermotesList represents the response from Get Cheermotes, see https://dev.twitch.tv/docs/api/reference#get-cheermotes
type CheermotesList struct {
	Data []struct {
		Prefix string `json:"prefix"`
		Tiers  []struct {
			MinBits int    `json:"min_bits"`
			ID      string `json:"id"`
			Color   string `json:"color"`
			Images  struct {
				Dark struct {
					Animated struct {
						One          string `json:"1"`
						OnePointFive string `json:"1.5"`
						Two          string `json:"2"`
						Three        string `json:"3"`
						Four         string `json:"4"`
					} `json:"animated"`
					Static struct {
						One          string `json:"1"`
						OnePointFive string `json:"1.5"`
						Two          string `json:"2"`
						Three        string `json:"3"`
						Four         string `json:"4"`
					} `json:"static"`
				} `json:"dark"`
				Light struct {
					Animated struct {
						One          string `json:"1"`
						OnePointFive string `json:"1.5"`
						Two          string `json:"2"`
						Three        string `json:"3"`
						Four         string `json:"4"`
					} `json:"animated"`
					Static struct {
						One          string `json:"1"`
						OnePointFive string `json:"1.5"`
						Two          string `json:"2"`
						Three        string `json:"3"`
						Four         string `json:"4"`
					} `json:"static"`
				} `json:"light"`
			} `json:"images"`
			CanCheer       bool `json:"can_cheer"`
			ShowInBitsCard bool `json:"show_in_bits_card"`
		} `json:"tiers"`
		Type         string    `json:"type"`
		Order        int       `json:"order"`
		LastUpdated  time.Time `json:"last_updated"`
		IsCharitable bool      `json:"is_charitable"`
	} `json:"data"`
}

// ExtensionTransactions represents the response from Get Extension Transactions, see https://dev.twitch.tv/docs/api/reference#get-extension-transactions
type ExtensionTransactions struct {
	Data []struct {
		ID               string    `json:"id"`
		Timestamp        time.Time `json:"timestamp"`
		BroadcasterID    string    `json:"broadcaster_id"`
		BroadcasterLogin string    `json:"broadcaster_login"`
		BroadcasterName  string    `json:"broadcaster_name"`
		UserID           string    `json:"user_id"`
		UserLogin        string    `json:"user_login"`
		UserName         string    `json:"user_name"`
		ProductType      string    `json:"product_type"`
		ProductData      struct {
			Sku  string `json:"sku"`
			Cost struct {
				Amount int    `json:"amount"`
				Type   string `json:"type"`
			} `json:"cost"`
			DisplayName   string `json:"displayName"`
			InDevelopment bool   `json:"inDevelopment"`
		} `json:"product_data"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// CustomRewardImage represents the various images associated with a custom reward.
type CustomRewardImage struct {
	Url1X string `json:"url_1x"`
	Url2X string `json:"url_2x"`
	Url4X string `json:"url_4x"`
}

// CustomRewards represents the response from Get Custom Rewards, see https://dev.twitch.tv/docs/api/reference#get-custom-reward
type CustomRewards struct {
	Data []struct {
		BroadcasterName     string             `json:"broadcaster_name"`
		BroadcasterLogin    string             `json:"broadcaster_login"`
		BroadcasterId       string             `json:"broadcaster_id"`
		ID                  string             `json:"id"`
		Image               *CustomRewardImage `json:"image"`
		BackgroundColor     string             `json:"background_color"`
		IsEnabled           bool               `json:"is_enabled"`
		Cost                int                `json:"cost"`
		Title               string             `json:"title"`
		Prompt              string             `json:"prompt"`
		IsUserInputRequired bool               `json:"is_user_input_required"`
		MaxPerStreamSetting struct {
			IsEnabled    bool `json:"is_enabled"`
			MaxPerStream int  `json:"max_per_stream"`
		} `json:"max_per_stream_setting"`
		MaxPerUserPerStreamSetting struct {
			IsEnabled           bool `json:"is_enabled"`
			MaxPerUserPerStream int  `json:"max_per_user_per_stream"`
		} `json:"max_per_user_per_stream_setting"`
		GlobalCooldownSetting struct {
			IsEnabled             bool `json:"is_enabled"`
			GlobalCooldownSeconds int  `json:"global_cooldown_seconds"`
		} `json:"global_cooldown_setting"`
		IsPaused                          bool               `json:"is_paused"`
		IsInStock                         bool               `json:"is_in_stock"`
		DefaultImage                      *CustomRewardImage `json:"default_image"`
		ShouldRedemptionsSkipRequestQueue bool               `json:"should_redemptions_skip_request_queue"`
		RedemptionsRedeemedCurrentStream  int                `json:"redemptions_redeemed_current_stream"`
		CooldownExpiresAt                 *time.Time         `json:"cooldown_expires_at"`
	} `json:"data"`
}

// CustomRewardRedemptions represents the response from Get Custom Reward Redemption, see https://dev.twitch.tv/docs/api/reference#get-custom-reward-redemption
type CustomRewardRedemptions struct {
	Data []struct {
		BroadcasterName  string    `json:"broadcaster_name"`
		BroadcasterLogin string    `json:"broadcaster_login"`
		BroadcasterID    string    `json:"broadcaster_id"`
		ID               string    `json:"id"`
		UserLogin        string    `json:"user_login"`
		UserID           string    `json:"user_id"`
		UserName         string    `json:"user_name"`
		UserInput        string    `json:"user_input"`
		Status           string    `json:"status"`
		RedeemedAt       time.Time `json:"redeemed_at"`
		Reward           struct {
			ID     string `json:"id"`
			Title  string `json:"title"`
			Prompt string `json:"prompt"`
			Cost   int    `json:"cost"`
		} `json:"reward"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// ChannelInformation represents the response from Get Channel Information, see https://dev.twitch.tv/docs/api/reference#get-channel-information
type ChannelInformation struct {
	Data []struct {
		BroadcasterID       string `json:"broadcaster_id"`
		BroadcasterLogin    string `json:"broadcaster_login"`
		BroadcasterName     string `json:"broadcaster_name"`
		BroadcasterLanguage string `json:"broadcaster_language"`
		GameID              string `json:"game_id"`
		GameName            string `json:"game_name"`
		Title               string `json:"title"`
		Delay               int    `json:"delay"`
	} `json:"data"`
}

// ChannelEditors represents the response from Get Channel Editors, see https://dev.twitch.tv/docs/api/reference#get-channel-editors
type ChannelEditors struct {
	Data []struct {
		UserID    string    `json:"user_id"`
		UserName  string    `json:"user_name"`
		CreatedAt time.Time `json:"created_at"`
	} `json:"data"`
}

// ChannelChatBadges represents the respons from Get Channel Chat Badges or Get Global Chat Badges, see https://dev.twitch.tv/docs/api/reference#get-channel-chat-badges
type ChannelChatBadges struct {
	Data []struct {
		SetID    string `json:"set_id"`
		Versions []struct {
			ID         string `json:"id"`
			ImageUrl1X string `json:"image_url_1x"`
			ImageUrl2X string `json:"image_url_2x"`
			ImageUrl4X string `json:"image_url_4x"`
		} `json:"versions"`
	} `json:"data"`
}

// Clips represents the response from Get Clips, see https://dev.twitch.tv/docs/api/reference#get-clips
type Clips struct {
	Data []struct {
		ID              string    `json:"id"`
		Url             string    `json:"url"`
		EmbedUrl        string    `json:"embed_url"`
		BroadcasterID   string    `json:"broadcaster_id"`
		BroadcasterName string    `json:"broadcaster_name"`
		CreatorID       string    `json:"creator_id"`
		CreatorName     string    `json:"creator_name"`
		VideoID         string    `json:"video_id"`
		GameID          string    `json:"game_id"`
		Language        string    `json:"language"`
		Title           string    `json:"title"`
		ViewCount       int       `json:"view_count"`
		CreatedAt       time.Time `json:"created_at"`
		ThumbnailUrl    string    `json:"thumbnail_url"`
		Duration        float64   `json:"duration"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// CodeStatus represents the response from Get Code Status, see https://dev.twitch.tv/docs/api/reference#get-code-status
type CodeStatus struct {
	Data []struct {
		Code   string `json:"code"`
		Status string `json:"status"`
	} `json:"data"`
}

// DropEntitlements represents the response from Get Drop Entitlements, see https://dev.twitch.tv/docs/api/reference#get-drops-entitlements
type DropEntitlements struct {
	Data []struct {
		ID        string    `json:"id"`
		BenefitID string    `json:"benefit_id"`
		Timestamp time.Time `json:"timestamp"`
		UserID    string    `json:"user_id"`
		GameID    string    `json:"game_id"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// EventSubSubscriptions represents the response from Get EventSub Subscriptions, see https://dev.twitch.tv/docs/api/reference#get-eventsub-subscriptions
type EventSubSubscriptions struct {
	Total int `json:"total"`
	Data  []struct {
		ID        string `json:"id"`
		Status    string `json:"status"`
		Type      string `json:"type"`
		Version   string `json:"version"`
		Condition struct {
			BroadcasterUserID string `json:"broadcaster_user_id,omitempty"`
			UserID            string `json:"user_id,omitempty"`
		} `json:"condition"`
		CreatedAt time.Time `json:"created_at"`
		Transport struct {
			Method   string `json:"method"`
			Callback string `json:"callback"`
		} `json:"transport"`
		Cost int `json:"cost"`
	} `json:"data"`
	TotalCost    int               `json:"total_cost"`
	MaxTotalCost int               `json:"max_total_cost"`
	Pagination   *TwitchPagination `json:"pagination"`
}

// Games represents the response from Get Top Games, Get Games, and Search Categories see https://dev.twitch.tv/docs/api/reference#get-top-games
type Games struct {
	Data []struct {
		ID        string `json:"id"`
		Name      string `json:"name"`
		BoxArtUrl string `json:"box_art_url"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// HypeTrainEvents represents the response from Get Hype Train Events, see https://dev.twitch.tv/docs/api/reference#get-hype-train-events
type HypeTrainEvents struct {
	Data []struct {
		ID             string    `json:"id"`
		EventType      string    `json:"event_type"`
		EventTimestamp time.Time `json:"event_timestamp"`
		Version        string    `json:"version"`
		EventData      struct {
			BroadcasterId    string    `json:"broadcaster_id"`
			CooldownEndTime  time.Time `json:"cooldown_end_time"`
			ExpiresAt        time.Time `json:"expires_at"`
			Goal             int       `json:"goal"`
			ID               string    `json:"id"`
			LastContribution struct {
				Total int    `json:"total"`
				Type  string `json:"type"`
				User  string `json:"user"`
			} `json:"last_contribution"`
			Level            int       `json:"level"`
			StartedAt        time.Time `json:"started_at"`
			TopContributions []struct {
				Total int    `json:"total"`
				Type  string `json:"type"`
				User  string `json:"user"`
			} `json:"top_contributions"`
			Total int `json:"total"`
		} `json:"event_data"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// BannedEvents represents the response from Get Banned Events, see https://dev.twitch.tv/docs/api/reference#get-banned-events
type BannedEvents struct {
	Data []struct {
		ID             string    `json:"id"`
		EventType      string    `json:"event_type"`
		EventTimestamp time.Time `json:"event_timestamp"`
		Version        string    `json:"version"`
		EventData      struct {
			BroadcasterID    string `json:"broadcaster_id"`
			BroadcasterLogin string `json:"broadcaster_login"`
			BroadcasterName  string `json:"broadcaster_name"`
			UserID           string `json:"user_id"`
			UserLogin        string `json:"user_login"`
			UserName         string `json:"user_name"`
			ExpiresAt        string `json:"expires_at"`
			// ExpiresAt is not a time.Time because "" will cause parsing errors
		} `json:"event_data"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// BannedUsers represents the response from Get Banned Users, see https://dev.twitch.tv/docs/api/reference#get-banned-users
type BannedUsers struct {
	Data []struct {
		UserID    string    `json:"user_id"`
		UserLogin string    `json:"user_login"`
		UserName  string    `json:"user_name"`
		ExpiresAt time.Time `json:"expires_at"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// Moderators represents the response from Get Moderators, see https://dev.twitch.tv/docs/api/reference#get-moderators
type Moderators struct {
	Data []struct {
		UserID    string `json:"user_id"`
		UserLogin string `json:"user_login"`
		UserName  string `json:"user_name"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// ModeratorEvents represents the response from Get Moderator Events, see https://dev.twitch.tv/docs/api/reference#get-moderator-events
type ModeratorEvents struct {
	Data []struct {
		ID             string    `json:"id"`
		EventType      string    `json:"event_type"`
		EventTimestamp time.Time `json:"event_timestamp"`
		Version        string    `json:"version"`
		EventData      struct {
			BroadcasterID    string `json:"broadcaster_id"`
			BroadcasterLogin string `json:"broadcaster_login"`
			BroadcasterName  string `json:"broadcaster_name"`
			UserID           string `json:"user_id"`
			UserLogin        string `json:"user_login"`
			UserName         string `json:"user_name"`
		} `json:"event_data"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// Polls represents the response from Get Polls, see https://dev.twitch.tv/docs/api/reference#get-polls
type Polls struct {
	Data []struct {
		ID               string `json:"id"`
		BroadcasterID    string `json:"broadcaster_id"`
		BroadcasterName  string `json:"broadcaster_name"`
		BroadcasterLogin string `json:"broadcaster_login"`
		Title            string `json:"title"`
		Choices          []struct {
			ID                 string `json:"id"`
			Title              string `json:"title"`
			Votes              int    `json:"votes"`
			ChannelPointsVotes int    `json:"channel_points_votes"`
			BitsVotes          int    `json:"bits_votes"`
		} `json:"choices"`
		BitsVotingEnabled          bool       `json:"bits_voting_enabled"`
		BitsPerVote                int        `json:"bits_per_vote"`
		ChannelPointsVotingEnabled bool       `json:"channel_points_voting_enabled"`
		ChannelPointsPerVote       int        `json:"channel_points_per_vote"`
		Status                     string     `json:"status"`
		Duration                   int        `json:"duration"`
		StartedAt                  time.Time  `json:"started_at"`
		EndedAt                    *time.Time `json:"ended_at"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// TopPredictors represents users who bet the most on their Predictions and won
type TopPredictors struct {
	UserID            string `json:"id"`
	UserName          string `json:"name"`
	UserLogin         string `json:"login"`
	ChannelPointsUsed int    `json:"channel_points_used"`
	ChannelPointsWon  int    `json:"channel_points_won"`
}

// Predictions represents the response from Get Predictions, see https://dev.twitch.tv/docs/api/reference#get-predictions
type Predictions struct {
	Data []struct {
		ID               string `json:"id"`
		BroadcasterID    string `json:"broadcaster_id"`
		BroadcasterName  string `json:"broadcaster_name"`
		BroadcasterLogin string `json:"broadcaster_login"`
		Title            string `json:"title"`
		WinningOutcomeId string `json:"winning_outcome_id"`
		Outcomes         []struct {
			ID            string           `json:"id"`
			Title         string           `json:"title"`
			Users         int              `json:"users"`
			ChannelPoints int              `json:"channel_points"`
			TopPredictors []*TopPredictors `json:"top_predictors"`
			Color         string           `json:"color"`
		} `json:"outcomes"`
		PredictionWindow int        `json:"prediction_window"`
		Status           string     `json:"status"`
		CreatedAt        time.Time  `json:"created_at"`
		EndedAt          *time.Time `json:"ended_at"`
		LockedAt         *time.Time `json:"locked_at"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// ChannelSearchResults represents the response from Search Channels, see https://dev.twitch.tv/docs/api/reference#search-channels
type ChannelSearchResults struct {
	Data []struct {
		BroadcasterLanguage string   `json:"broadcaster_language"`
		BroadcasterLogin    string   `json:"broadcaster_login"`
		DisplayName         string   `json:"display_name"`
		GameID              string   `json:"game_id"`
		GameName            string   `json:"game_name"`
		ID                  string   `json:"id"`
		IsLive              bool     `json:"is_live"`
		TagIDs              []string `json:"tag_ids"`
		ThumbnailUrl        string   `json:"thumbnail_url"`
		Title               string   `json:"title"`
		StartedAt           string   `json:"started_at"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// StreamKey represents the response from Get Stream Key, see https://dev.twitch.tv/docs/api/reference#get-stream-key
type StreamKey struct {
	Data []struct {
		StreamKey string `json:"stream_key"`
	} `json:"data"`
}

// Streams represents the response from Get Streams and Get Followed Streams, see https://dev.twitch.tv/docs/api/reference#get-streams
type Streams struct {
	Data []struct {
		ID           string    `json:"id"`
		UserID       string    `json:"user_id"`
		UserLogin    string    `json:"user_login"`
		UserName     string    `json:"user_name"`
		GameID       string    `json:"game_id"`
		GameName     string    `json:"game_name"`
		Type         string    `json:"type"`
		Title        string    `json:"title"`
		ViewerCount  int       `json:"viewer_count"`
		StartedAt    time.Time `json:"started_at"`
		Language     string    `json:"language"`
		ThumbnailUrl string    `json:"thumbnail_url"`
		TagIDs       []string  `json:"tag_ids"`
		IsMature     bool      `json:"is_mature"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// StreamMarkers represents the response from Get Stream Markers, see https://dev.twitch.tv/docs/api/reference#get-stream-markers
type StreamMarkers struct {
	Data []struct {
		UserID    string `json:"user_id"`
		UserName  string `json:"user_name"`
		UserLogin string `json:"user_login"`
		Videos    []struct {
			VideoId string `json:"video_id"`
			Markers []struct {
				ID              string    `json:"id"`
				CreatedAt       time.Time `json:"created_at"`
				Description     string    `json:"description"`
				PositionSeconds int       `json:"position_seconds"`
				URL             string    `json:"URL"`
			} `json:"markers"`
		} `json:"videos"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// Subscriptions represents the response from Get Broadcaster Subscriptions, see https://dev.twitch.tv/docs/api/reference#get-broadcaster-subscriptions
type Subscriptions struct {
	Data []struct {
		BroadcasterID    string `json:"broadcaster_id"`
		BroadcasterLogin string `json:"broadcaster_login"`
		BroadcasterName  string `json:"broadcaster_name"`
		GifterID         string `json:"gifter_id"`
		GifterLogin      string `json:"gifter_login"`
		GifterName       string `json:"gifter_name"`
		IsGift           bool   `json:"is_gift"`
		Tier             string `json:"tier"`
		PlanName         string `json:"plan_name"`
		UserID           string `json:"user_id"`
		UserName         string `json:"user_name"`
		UserLogin        string `json:"user_login"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
	Total      int               `json:"total"`
}

// UserSubscriptions represents the response from Check User Subscription, see https://dev.twitch.tv/docs/api/reference#check-user-subscription
type UserSubscriptions struct {
	Data []struct {
		BroadcasterId    string `json:"broadcaster_id"`
		BroadcasterName  string `json:"broadcaster_name"`
		BroadcasterLogin string `json:"broadcaster_login"`
		GifterID         string `json:"gifter_id"`
		GifterName       string `json:"gifter_name"`
		GifterLogin      string `json:"gifter_login"`
		IsGift           bool   `json:"is_gift"`
		Tier             string `json:"tier"`
	} `json:"data"`
}

// StreamTags represents the response from Get All Stream Tags and Get Stream Tags, see https://dev.twitch.tv/docs/api/reference#get-all-stream-tags
type StreamTags struct {
	Data []struct {
		TagID                    string            `json:"tag_id"`
		IsAuto                   bool              `json:"is_auto"`
		LocalizationNames        map[string]string `json:"localization_names"`
		LocalizationDescriptions map[string]string `json:"localization_descriptions"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// ChannelTeams represents the response from Get Channel Teams, see https://dev.twitch.tv/docs/api/reference#get-channel-teams
type ChannelTeams struct {
	Data []struct {
		BroadcasterID      string    `json:"broadcaster_id"`
		BroadcasterName    string    `json:"broadcaster_name"`
		BroadcasterLogin   string    `json:"broadcaster_login"`
		BackgroundImageUrl string    `json:"background_image_url"`
		Banner             string    `json:"banner"`
		CreatedAt          time.Time `json:"created_at"`
		UpdatedAt          time.Time `json:"updated_at"`
		Info               string    `json:"info"`
		ThumbnailUrl       string    `json:"thumbnail_url"`
		TeamName           string    `json:"team_name"`
		TeamDisplayName    string    `json:"team_display_name"`
		ID                 string    `json:"id"`
	} `json:"data"`
}

// Teams represents the response from Get Teams, see https://dev.twitch.tv/docs/api/reference#get-teams
type Teams struct {
	Data []struct {
		Users []struct {
			UserID    string `json:"user_id"`
			UserName  string `json:"user_name"`
			UserLogin string `json:"user_login"`
		} `json:"users"`
		BackgroundImageUrl interface{} `json:"background_image_url"`
		Banner             interface{} `json:"banner"`
		CreatedAt          time.Time   `json:"created_at"`
		UpdatedAt          time.Time   `json:"updated_at"`
		Info               string      `json:"info"`
		ThumbnailUrl       string      `json:"thumbnail_url"`
		TeamName           string      `json:"team_name"`
		TeamDisplayName    string      `json:"team_display_name"`
		ID                 string      `json:"id"`
	} `json:"data"`
}

// Users represents the response from Get Users, see https://dev.twitch.tv/docs/api/reference#get-users
type Users struct {
	Data []struct {
		ID              string    `json:"id"`
		Login           string    `json:"login"`
		DisplayName     string    `json:"display_name"`
		Type            string    `json:"type"`
		BroadcasterType string    `json:"broadcaster_type"`
		Description     string    `json:"description"`
		ProfileImageUrl string    `json:"profile_image_url"`
		OfflineImageUrl string    `json:"offline_image_url"`
		ViewCount       int       `json:"view_count"`
		Email           string    `json:"email"`
		CreatedAt       time.Time `json:"created_at"`
	} `json:"data"`
}

// UserFollows represents the response from Get User Follows, see https://dev.twitch.tv/docs/api/reference#get-users-follows
type UserFollows struct {
	Total int `json:"total"`
	Data  []struct {
		FromId     string    `json:"from_id"`
		FromLogin  string    `json:"from_login"`
		FromName   string    `json:"from_name"`
		ToId       string    `json:"to_id"`
		ToName     string    `json:"to_name"`
		FollowedAt time.Time `json:"followed_at"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// UserBlockList represents the response from Get User Block List, see https://dev.twitch.tv/docs/api/reference#get-user-block-list
type UserBlockList struct {
	Data []struct {
		UserID      string `json:"user_id"`
		UserLogin   string `json:"user_login"`
		DisplayName string `json:"display_name"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// UserExtensions represents the response from Get User Extensions, see https://dev.twitch.tv/docs/api/reference#get-user-extensions
type UserExtensions struct {
	Data []struct {
		ID          string   `json:"id"`
		Version     string   `json:"version"`
		Name        string   `json:"name"`
		CanActivate bool     `json:"can_activate"`
		Type        []string `json:"type"`
	} `json:"data"`
}

// UserActiveExtensions represents the response to Get User Active Extensions, see https://dev.twitch.tv/docs/api/reference#get-user-active-extensions
type UserActiveExtensions struct {
	Data struct {
		Panel struct {
			Field1 struct {
				Active  bool   `json:"active"`
				ID      string `json:"id"`
				Version string `json:"version"`
				Name    string `json:"name"`
			} `json:"1"`
			Field2 struct {
				Active  bool   `json:"active"`
				ID      string `json:"id"`
				Version string `json:"version"`
				Name    string `json:"name"`
			} `json:"2"`
			Field3 struct {
				Active  bool   `json:"active"`
				ID      string `json:"id"`
				Version string `json:"version"`
				Name    string `json:"name"`
			} `json:"3"`
		} `json:"panel"`
		Overlay struct {
			Field1 struct {
				Active  bool   `json:"active"`
				ID      string `json:"id"`
				Version string `json:"version"`
				Name    string `json:"name"`
			} `json:"1"`
		} `json:"overlay"`
		Component struct {
			Field1 struct {
				Active  bool   `json:"active"`
				ID      string `json:"id"`
				Version string `json:"version"`
				Name    string `json:"name"`
				X       int    `json:"x"`
				Y       int    `json:"y"`
			} `json:"1"`
			Field2 struct {
				Active  bool   `json:"active"`
				ID      string `json:"id"`
				Version string `json:"version"`
				Name    string `json:"name"`
				X       int    `json:"x"`
				Y       int    `json:"y"`
			} `json:"2"`
		} `json:"component"`
	} `json:"data"`
}

// Video represents the response from Get Videos, see https://dev.twitch.tv/docs/api/reference#get-videos
type Video struct {
	Data []struct {
		ID            string    `json:"id"`
		StreamID      string    `json:"stream_id"`
		UserID        string    `json:"user_id"`
		UserLogin     string    `json:"user_login"`
		UserName      string    `json:"user_name"`
		Title         string    `json:"title"`
		Description   string    `json:"description"`
		CreatedAt     time.Time `json:"created_at"`
		PublishedAt   time.Time `json:"published_at"`
		Url           string    `json:"url"`
		ThumbnailUrl  string    `json:"thumbnail_url"`
		Viewable      string    `json:"viewable"`
		ViewCount     int       `json:"view_count"`
		Language      string    `json:"language"`
		Type          string    `json:"type"`
		Duration      string    `json:"duration"`
		MutedSegments []struct {
			Duration int `json:"duration"`
			Offset   int `json:"offset"`
		} `json:"muted_segments"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// WebhookSubscriptions represents the response from Get Webhook Subscriptions, see https://dev.twitch.tv/docs/api/reference#get-webhook-subscriptions
type WebhookSubscriptions struct {
	Total int `json:"total"`
	Data  []struct {
		Topic     string    `json:"topic"`
		Callback  string    `json:"callback"`
		ExpiresAt time.Time `json:"expires_at"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}

// Vacation Represents a vacation object as part of a ChannelStreamSchedule
type Vacation struct {
	StartTime time.Time `json:"start_time"`
	EndTime   time.Time `json:"end_time"`
}

// ChannelStreamSchedule represents the stream schedule data
type ChannelStreamSchedule struct {
	Data struct {
		Segments []struct {
			ID            string     `json:"id"`
			StartTime     time.Time  `json:"start_time"`
			EndTime       time.Time  `json:"end_time"`
			Title         string     `json:"title"`
			CanceledUntil *time.Time `json:"canceled_until"`
			Category      struct {
				Id   string `json:"id"`
				Name string `json:"name"`
			} `json:"category"`
			IsRecurring bool `json:"is_recurring"`
		} `json:"segments"`
		BroadcasterId    string    `json:"broadcaster_id"`
		BroadcasterName  string    `json:"broadcaster_name"`
		BroadcasterLogin string    `json:"broadcaster_login"`
		Vacation         *Vacation `json:"vacation"`
	} `json:"data"`
	Pagination *TwitchPagination `json:"pagination"`
}
