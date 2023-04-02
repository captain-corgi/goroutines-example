package news

type (
	Article struct {
		Title   string `json:"title"`
		Link    string `json:"link"`
		Id      string `json:"id"`
		MongoId string `json:"_id"`
	}
	News struct {
		Source   string
		Articles []*Article `json:"articles"`
	}
)
