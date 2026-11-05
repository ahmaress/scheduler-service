package postgres

// normalizeTime converts database time format (HH:MM:SS.microseconds) to HH:MM format
func normalizeTime(timeStr string) string {
	// Handle empty string
	if timeStr == "" {
		return ""
	}
	// Take only first 5 characters (HH:MM)
	if len(timeStr) >= 5 {
		return timeStr[:5]
	}
	return timeStr
}
