package auth

func sessionExtractor() string {
	return "loggedin"
}

func GetSession() string {
	return sessionExtractor()
}
