package chucknorris

type JokeResult struct {
	IconURL string `json:"icon_url"`
	ID      string `json:"id"`
	URL     string `json:"url"`
	Value   string `json:"value"`
}
