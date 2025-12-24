package viewmodel

type TokenViewModel struct {
	Success bool   `json:"success"`
	Expires int64  `json:"expires"`
	Title   string `json:"token"`
}

func ToTokenViewModel(success bool, expires int64, title string) TokenViewModel {
	return TokenViewModel{
		Success: success,
		Expires: expires,
		Title:   title,
	}
}
