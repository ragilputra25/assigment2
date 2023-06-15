package users

type UserFormatter struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Token string `json:"token"`
}

func FormatUser(user User, token string) UserFormatter {
	formatter := UserFormatter{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Token: token,
	}

	return formatter
}

type UserRegisterFormatter struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func FormatRegister(user User) UserRegisterFormatter {
	formatter := UserRegisterFormatter{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}

	return formatter
}
