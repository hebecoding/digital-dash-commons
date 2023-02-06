package profiles

type NewAccountProfile struct {
	ProfileID     string `json:"profile_id"`
	Username      string `json:"username"`
	FirstName     string `json:"first_name"`
	LastName      string `json:"last_name"`
	Email         string `json:"email"`
	EmailToken    string `json:"email_token"`
	ValidationURL string `json:"validation_url"`
}
