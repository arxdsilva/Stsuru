package persist

import "github.com/asaskevich/govalidator"

// ValidateURL takes an URL and validades If It's real or not
func ValidateURL(l string) bool {
	isURL := govalidator.IsURL(l)
	validURL := govalidator.IsRequestURL(l)
	if isURL == false || validURL == false {
		return false
	}
	return true
}
