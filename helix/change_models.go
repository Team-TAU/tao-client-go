package helix

import "time"

type CustomRewardsUpdate struct {
	Title                             *string `json:"title,omitempty"`
	Prompt                            *string `json:"prompt,omitempty"`
	Cost                              *int    `json:"cost,omitempty"`
	BackgroundColor                   *string `json:"background_color,omitempty"`
	IsEnabled                         *bool   `json:"is_enabled,omitempty"`
	IsUserInputRequired               *bool   `json:"is_user_input_required,omitempty"`
	IsMaxPerStreamEnabled             *bool   `json:"is_max_per_stream_enabled,omitempty"`
	MaxPerStream                      *int    `json:"max_per_stream,omitempty"`
	IsMaxPerUserPerStreamEnabled      *bool   `json:"is_max_per_user_per_stream_enabled,omitempty"`
	MaxPerUserPerStream               *int    `json:"max_per_user_per_stream,omitempty"`
	IsGlobalCooldownEnabled           *bool   `json:"is_global_cooldown_enabled,omitempty"`
	GlobalCooldownSeconds             *int    `json:"global_cooldown_seconds,omitempty"`
	IsPaused                          *bool   `json:"is_paused,omitempty"`
	ShouldRedemptionsSkipRequestQueue *bool   `json:"should_redemptions_skip_request_queue,omitempty"`
}

type StreamScheduleSegmentUpdate struct {
	StartTime   *time.Time `json:"start_time,omitempty"`
	Timezone    *string    `json:"timezone,omitempty"`
	IsRecurring *bool      `json:"is_recurring,omitempty"`
	Duration    *string    `json:"duration,omitempty"`
	CategoryID  *string    `json:"category_id,omitempty"`
	Title       *string    `json:"title,omitempty"`
}
