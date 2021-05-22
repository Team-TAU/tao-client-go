package go_tau

type Event struct {
	ID          string `json:"id"`
	EventID     string `json:"event_id"`
	EventType   string `json:"event_type"`
	EventSource string `json:"event_source"`
	Created     string `json:"created"`
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
		BroadcasterLogin string `json:"broadcaster_login"`
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
		Bits             string `json:"bits"`
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
		Viewers              string `json:"viewers"`
	} `json:"event_data"`
}

type SubscriptionMsg struct {
	*Event
	EventData struct {
		Type string `json:"type"`
		Data struct {
			Topic              string `json:"topic"`
			SubPlan            string `json:"sub_plan"`
			SubPlanName        string `json:"sub_plan_name"`
			Months             int    `json:"months"`
			CumulativeMonths   int    `json:"cumulative_months"`
			Context            string `json:"context"`
			IsGift             bool   `json:"is_gift"`
			MultiMonthDuration int    `json:"multi_month_duration"`
			StreakMonths       int    `json:"streak_months"`
			Message            struct {
				BenefitEndMonth int    `json:"benefit_end_month"`
				UserName        string `json:"user_name"`
				DisplayName     string `json:"display_name"`
				ChannelName     string `json:"channel_name"`
				UserID          string `json:"user_id"`
				ChannelID       string `json:"channel_id"`
				Time            string `json:"time"`
				SubMessage      struct {
					Message string `json:"message"`
					Emotes  string `json:"emotes"`
				} `json:"sub_message"`
			} `json:"message"`
		} `json:"data"`
	}
}
