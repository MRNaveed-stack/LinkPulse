package models

import "time"

// Request/response DTOs used by handlers.

type RegisterRequest struct {
	Username string `json:"username" binding:"required,min=2,max=50"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	Token        string `json:"token"`
}

type RefreshRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email" binding:"required,email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token" binding:"required"`
	NewPassword string `json:"new_password" binding:"required,min=8"`
}

type CreateLinkRequest struct {
	Title          string `json:"title" binding:"required"`
	Slug           string `json:"slug" binding:"required"`
	DestinationURL string `json:"destination_url" binding:"required,url"`
}

type AnalyticsOverview struct {
	TotalLinks  int64       `json:"total_links"`
	TotalClicks int64       `json:"total_clicks"`
	ClicksToday int64       `json:"clicks_today"`
	ActiveLinks int64       `json:"active_links"`
	TopLink     *TopLinkDTO `json:"top_link"`
}

type TopLinkDTO struct {
	Title  string `json:"title"`
	Slug   string `json:"slug"`
	Clicks int64  `json:"clicks"`
}

type PublicProfileResponse struct {
	Username    string          `json:"username"`
	DisplayName string          `json:"display_name"`
	Bio         string          `json:"bio"`
	AvatarURL   string          `json:"avatar_url"`
	Links       []PublicLinkDTO `json:"links"`
}

type PublicLinkDTO struct {
	Title string `json:"title"`
	Slug  string `json:"slug"`
}

type LinkAnalyticsDTO struct {
	Title  string `json:"title"`
	Slug   string `json:"slug"`
	Clicks int64  `json:"clicks"`
}

type UpdateLinkRequest struct {
	Title          string `json:"title"`
	Slug           string `json:"slug"`
	DestinationURL string `json:"destination_url"`
}

type DailyAnalyticsDTO struct {
	Date   string `json:"date"`
	Clicks int64  `json:"clicks"`
}

type AnalyticsQuery struct {
	Days int `form:"days"`
}

type RecentActivityDTO struct {
	LinkTitle string    `json:"link_title"`
	Slug      string    `json:"slug"`
	ClickedAt time.Time `json:"clicked_at"`
	IPAddress string    `json:"ip_address"`
	Referrer  string    `json:"referrer"`
}

type RecentActivityQuery struct {
	Limit int `form:"limit"`
}

type ReferrerAnalyticsDTO struct {
	Source string `json:"referrer"`
	Clicks int64  `json:"clicks"`
}
