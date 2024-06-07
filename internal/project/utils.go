package project


func RequiredFields(newProject Project) bool {
	if newProject.Title == "" {
		return false
	}
	if newProject.Content == "" {
		return false
	}
	if newProject.Link == "" {
		return false
	}
	if newProject.Tags == "" {
		return false
	}
	return true
}