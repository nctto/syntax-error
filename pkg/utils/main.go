package utils

import (
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

func DateToString(date primitive.DateTime) string {
	d := date.Time()
	currentTime := time.Now()
	diff := currentTime.Sub(d)
	days := int(diff.Hours() / 24)
	hours := int(diff.Hours())
	minutes := int(diff.Minutes())
	seconds := int(diff.Seconds())

	if days > 0 {
		return fmt.Sprintf("%d days ago", days)
	} else if hours > 0 {
		return fmt.Sprintf("%d hours ago", hours)
	} else if minutes > 0 {
		return fmt.Sprintf("%d minutes ago", minutes)
	} else if seconds > 0 {
		return fmt.Sprintf("%d seconds ago", seconds)
	}
	return date.Time().String()
}