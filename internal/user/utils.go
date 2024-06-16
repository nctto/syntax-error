package user


func UserNickName(user interface{}) string{
	if user == nil {
		return ""
	}
	nickname := user.(map[string]interface{})["nickname"].(string)
	return nickname
}