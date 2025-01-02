package commons

type ClientCredentials struct {
	ClientName   string `json:"client_name"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}
