package defaults

func String(value, defaultValue string) string {
	if value != "" {
		return value
	}

	return defaultValue
}
