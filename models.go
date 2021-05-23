package go_tau

import "time"

type Event struct {
	ID          string `json:"id"`
	EventID     string `json:"event_id"`
	EventType   string `json:"event_type"`
	EventSource string `json:"event_source"`
	Created     Time   `json:"created"`
	Origin      string `json:"origin"`
}

type FollowMsg struct {
	*Event
	EventData struct {
		UserName         string `json:"user_name"`
		UserID           string `json:"user_id"`
		UserLogin        string `json:"user_login"`
		BroadcasterID    string `json:"broadcaster_user_id"`
		BroadcasterName  string `json:"broadcaster_user_name"`
		BroadcasterLogin string `json:"broadcaster_user_login"`
	} `json:"event_data"`
}

type StreamUpdateMsg struct {
	*Event
	EventData struct {
		Title            string `json:"title"`
		Language         string `json:"language"`
		IsMature         bool   `json:"is_mature"`
		CategoryID       int    `json:"category_id"`
		CategoryName     string `json:"category_name"`
		BroadcasterID    string `json:"broadcaster_user_id"`
		BroadcasterName  string `json:"broadcaster_user_name"`
		BroadcasterLogin string `json:"broadcaster_user_login"`
	} `json:"event_data"`
}

type CheerMsg struct {
	*Event
	EventData struct {
		IsAnonymous      bool   `json:"is_anonymous"`
		UserID           string `json:"user_id"`
		UserName         string `json:"user_name"`
		UserLogin        string `json:"user_login"`
		BroadcasterID    string `json:"broadcaster_user_id"`
		BroadcasterName  string `json:"broadcaster_user_name"`
		BroadcasterLogin string `json:"broadcaster_user_login"`
		Bits             int    `json:"bits"`
		Message          string `json:"message"`
	} `json:"event_data"`
}

type RaidMsg struct {
	*Event
	EventData struct {
		FromBroadcasterName  string `json:"from_broadcaster_user_name"`
		FromBroadcasterID    string `json:"from_broadcaster_user_id"`
		FromBroadcasterLogin string `json:"from_broadcaster_user_login"`
		ToBroadcasterName    string `json:"to_broadcaster_user_name"`
		ToBroadcasterID      string `json:"to_broadcaster_user_id"`
		ToBroadcasterLogin   string `json:"to_broadcaster_user_login"`
		Viewers              int    `json:"viewers"`
	} `json:"event_data"`
}

type SubscriptionMsg struct {
	*Event
	EventData struct {
		Type string `json:"type"`
		Data struct {
			Topic   string `json:"topic"`
			Message struct {
				BenefitEndMonth    int       `json:"benefit_end_month"`
				UserName           string    `json:"user_name"`
				DisplayName        string    `json:"display_name"`
				ChannelName        string    `json:"channel_name"`
				UserID             string    `json:"user_id"`
				ChannelID          string    `json:"channel_id"`
				Time               time.Time `json:"time"`
				SubPlan            string    `json:"sub_plan"`
				SubPlanName        string    `json:"sub_plan_name"`
				Months             int       `json:"months"`
				CumulativeMonths   int       `json:"cumulative_months"`
				Context            string    `json:"context"`
				IsGift             bool      `json:"is_gift"`
				MultiMonthDuration int       `json:"multi_month_duration"`
				StreakMonths       int       `json:"streak_months"`
				SubMessage         struct {
					Message string `json:"message"`
					Emotes  []struct {
						Start int `json:"start"`
						End   int `json:"end"`
						ID    int `json:"id"`
					} `json:"emotes"`
				} `json:"sub_message"`
			} `json:"message"`
		} `json:"data"`
	} `json:"event_data"`
}

type HypeTrainBeginMsg struct {
	*Event
	EventData struct {
		BroadcasterID    string    `json:"broadcaster_user_id"`
		BroadcasterName  string    `json:"broadcaster_user_name"`
		BroadcasterLogin string    `json:"broadcaster_user_login"`
		Total            int       `json:"total"`
		Progress         int       `json:"progress"`
		Goal             int       `json:"goal"`
		StartedAt        time.Time `json:"started_at"`
		ExpiresAt        time.Time `json:"expires_at"`
		TopContributions []struct {
			UserID    string `json:"user_id"`
			UserLogin string `json:"user_login"`
			UserName  string `json:"user_name"`
			Type      string `json:"type"`
			Total     int    `json:"total"`
		} `json:"top_contributions"`
		LastContribution struct {
			UserID    string `json:"user_id"`
			UserLogin string `json:"user_login"`
			UserName  string `json:"user_name"`
			Type      string `json:"type"`
			Total     int    `json:"total"`
		} `json:"last_contribution"`
	} `json:"event_data"`
}

type HypeTrainProgressMsg struct {
	*Event
	EventData struct {
		BroadcasterID    string    `json:"broadcaster_user_id"`
		BroadcasterName  string    `json:"broadcaster_user_name"`
		BroadcasterLogin string    `json:"broadcaster_user_login"`
		Level            int       `json:"level"`
		Total            int       `json:"total"`
		Progress         int       `json:"progress"`
		Goal             int       `json:"goal"`
		StartedAt        time.Time `json:"started_at"`
		ExpiresAt        time.Time `json:"expires_at"`
		TopContributions []struct {
			UserID    string `json:"user_id"`
			UserLogin string `json:"user_login"`
			UserName  string `json:"user_name"`
			Type      string `json:"type"`
			Total     int    `json:"total"`
		} `json:"top_contributions"`
		LastContribution struct {
			UserID    string `json:"user_id"`
			UserLogin string `json:"user_login"`
			UserName  string `json:"user_name"`
			Type      string `json:"type"`
			Total     int    `json:"total"`
		} `json:"last_contribution"`
	} `json:"event_data"`
}

type HypeTrainEndedMsg struct {
	*Event
	EventData struct {
		BroadcasterID    string    `json:"broadcaster_user_id"`
		BroadcasterName  string    `json:"broadcaster_user_name"`
		BroadcasterLogin string    `json:"broadcaster_user_login"`
		Level            int       `json:"level"`
		Total            int       `json:"total"`
		Progress         int       `json:"progress"`
		StartedAt        time.Time `json:"started_at"`
		EndedAt          time.Time `json:"ended_at"`
		CooldownEndsAt   time.Time `json:"cooldown_ends_at"`
		TopContributions []struct {
			UserID    string `json:"user_id"`
			UserLogin string `json:"user_login"`
			UserName  string `json:"user_name"`
			Type      string `json:"type"`
			Total     int    `json:"total"`
		} `json:"top_contributions"`
	} `json:"event_data"`
}

type StreamOnlineMsg struct {
	*Event
	EventData struct {
		ID               string    `json:"id"`
		BroadcasterID    string    `json:"broadcaster_user_id"`
		BroadcasterName  string    `json:"broadcaster_user_name"`
		BroadcasterLogin string    `json:"broadcaster_user_login"`
		Type             string    `json:"type"`
		StartedAt        time.Time `json:"started_at"`
	} `json:"event_data"`
}

type StreamOfflineMsg struct {
	*Event
	EventData struct {
		BroadcasterID    string `json:"broadcaster_user_id"`
		BroadcasterName  string `json:"broadcaster_user_name"`
		BroadcasterLogin string `json:"broadcaster_user_login"`
	} `json:"event_data"`
}

type PointsRedemptionMsg struct {
	*Event
	EventData struct {
		BroadcasterID    string    `json:"broadcaster_user_id"`
		BroadcasterName  string    `json:"broadcaster_user_name"`
		BroadcasterLogin string    `json:"broadcaster_user_login"`
		ID               string    `json:"id"`
		UserID           string    `json:"user_id"`
		UserLogin        string    `json:"user_login"`
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
	} `json:"event_data"`
}
