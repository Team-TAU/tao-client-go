package go_tau

import "time"

type Event struct {
	ID          string `json:"id"`
	EventID     string `json:"event_id"`
	EventType   string `json:"event_type"`
	EventSource string `json:"event_source"`
	Created     string `json:"created"`
	Origin      string `json:"origin"`
}

func (e *Event) CreatedAsTime() (time.Time, error) {
	// 2021-05-22T05:20:06.120452+00:00
	return time.Parse("2006-01-02T15:04:05.999999999-07:00", e.Created)
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
					Emotes  []struct {
						Start int `json:"start"`
						End   int `json:"end"`
						ID    int `json:"id"`
					} `json:"emotes"`
				} `json:"sub_message"`
			} `json:"message"`
		} `json:"data"`
	}
}

type HypeTrainBeginMsg struct {
	*Event
	EventData struct {
		BroadcasterID    string `json:"broadcaster_user_id"`
		BroadcasterName  string `json:"broadcaster_user_name"`
		BroadcasterLogin string `json:"broadcaster_user_login"`
		Total            int    `json:"total"`
		Progress         int    `json:"progress"`
		Goal             int    `json:"goal"`
		StartedAt        string `json:"started_at"`
		ExpiresAt        string `json:"expires_at"`
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

// StartedAtAsTime returns a Time object for the started at (which isn't needed often hence why
// it's not doing the conversion by default)
func (h *HypeTrainBeginMsg) StartedAtAsTime() (time.Time, error) {
	//2020-07-15T17:16:03.17106713Z
	return time.Parse("2006-01-02T15:04:05.999999999Z", h.EventData.StartedAt)
}

// ExpiresAtAsTime returns a Time object for the started at (which isn't needed often hence why
// it's not doing the conversion by default)
func (h *HypeTrainBeginMsg) ExpiresAtAsTime() (time.Time, error) {
	//2020-07-15T17:16:11.17106713Z
	return time.Parse("2006-01-02T15:04:05.999999999Z", h.EventData.ExpiresAt)
}

type HypeTrainProgressMsg struct {
	*Event
	EventData struct {
		BroadcasterID    string `json:"broadcaster_user_id"`
		BroadcasterName  string `json:"broadcaster_user_name"`
		BroadcasterLogin string `json:"broadcaster_user_login"`
		Level            int    `json:"level"`
		Total            int    `json:"total"`
		Progress         int    `json:"progress"`
		Goal             int    `json:"goal"`
		StartedAt        string `json:"started_at"`
		ExpiresAt        string `json:"expires_at"`
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

// StartedAtAsTime returns a Time object for the started at (which isn't needed often hence why
// it's not doing the conversion by default)
func (h *HypeTrainProgressMsg) StartedAtAsTime() (time.Time, error) {
	// 2020-07-15T17:16:03.17106713Z
	return time.Parse("2006-01-02T15:04:05.999999999Z", h.EventData.StartedAt)
}

// ExpiresAtAsTime returns a Time object for the started at (which isn't needed often hence why
// it's not doing the conversion by default)
func (h *HypeTrainProgressMsg) ExpiresAtAsTime() (time.Time, error) {
	// 2020-07-15T17:16:03.17106713Z
	return time.Parse("2006-01-02T15:04:05.999999999Z", h.EventData.ExpiresAt)
}

type HypeTrainEndedMsg struct {
	*Event
	EventData struct {
		BroadcasterID    string `json:"broadcaster_user_id"`
		BroadcasterName  string `json:"broadcaster_user_name"`
		BroadcasterLogin string `json:"broadcaster_user_login"`
		Level            int    `json:"level"`
		Total            int    `json:"total"`
		Progress         int    `json:"progress"`
		StartedAt        string `json:"started_at"`
		EndedAt          string `json:"ended_at"`
		CooldownEndsAt   string `json:"cooldown_ends_at"`
		TopContributions []struct {
			UserID    string `json:"user_id"`
			UserLogin string `json:"user_login"`
			UserName  string `json:"user_name"`
			Type      string `json:"type"`
			Total     int    `json:"total"`
		} `json:"top_contributions"`
	} `json:"event_data"`
}

// StartedAtAsTime returns a Time object for the started at (which isn't needed often hence why
// it's not doing the conversion by default)
func (h *HypeTrainEndedMsg) StartedAtAsTime() (time.Time, error) {
	// 2020-07-15T17:16:03.17106713Z
	return time.Parse("2006-01-02T15:04:05.999999999Z", h.EventData.StartedAt)
}

// EndedAtAsTime returns a Time object for the started at (which isn't needed often hence why
// it's not doing the conversion by default)
func (h *HypeTrainEndedMsg) EndedAtAsTime() (time.Time, error) {
	// 2020-07-15T17:16:03.17106713Z
	return time.Parse("2006-01-02T15:04:05.999999999Z", h.EventData.EndedAt)
}

// CooldownEndsAtAsTime returns a Time object for the started at (which isn't needed often hence why
// it's not doing the conversion by default)
func (h *HypeTrainEndedMsg) CooldownEndsAtAsTime() (time.Time, error) {
	// 2020-07-15T18:16:11.17106713Z
	return time.Parse("2006-01-02T15:04:05.999999999Z", h.EventData.CooldownEndsAt)
}
