package models

type GitHubProfile struct {
	Login     string `json:"login"`
	ID        int    `json:"id"`
	AvatarURL string `json:"avatar_url"`
	Name      string `json:"name"`
	Blog      string `json:"blog"`
	Bio       string `json:"bio"`
}

type UserProfile struct {
	ID          string        `json:"_id"`
	ProfileID   string        `json:"profileId"`
	AccessToken string        `json:"accessToken"`
	Avatar      string        `json:"avatar"`
	Name        string        `json:"name"`
	Username    string        `json:"username"`
	RawProfile  GitHubProfile `json:"_json"`
}
