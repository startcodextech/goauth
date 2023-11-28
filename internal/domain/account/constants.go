package account

const (
	// NameOrLastnameRegexp is the regular expression used to validate a name or lastname.
	NameOrLastnameRegexp = `^[a-zA-ZáéíóúÁÉÍÓÚñÑüÜ' -]{2,50}$`

	// EmailRegexp is the regular expression used to validate an email.
	EmailRegexp = `^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`
)
