package helpers

import (
	"log"
	"regexp"
)

type Validator struct{}

var Validators = Validator{}

func (f *Validator) Email(email string) bool {
	regex := `^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`
	match, err := regexp.MatchString(regex, email)
	if err != nil {
		log.Println(err)
		return false
	}
	return match
}

func (f *Validator) Age(age uint8) bool {
	if age >= 18 && age <= 120 {
		return true
	}

	return false
}
