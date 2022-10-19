// best practice to put all files of "x" package in "x" folder vice versa for "y"

package helper // saying that this file belongs to the "helper" package --> could also put it in "main" package

import (
	"strings"
)

// variables and functions at the package level can be referenced by any file in the package
// by capitalizing the first letter of the function, it will export the function which allows it to be
// ... usable by other packages within the module
func Validate_user_input(first_name string, last_name string, user_email string, user_tickets uint, remaining_tickets uint) (bool, bool, bool, bool) {
	valid_name := len(first_name) >= 2 && len(last_name) >= 2 // both first and last name >= 2 characters
	valid_email := strings.Contains(user_email, "@")          // email contains @
	valid_ticket_positive := user_tickets > 0
	valid_ticket_remain := user_tickets <= remaining_tickets
	return valid_name, valid_email, valid_ticket_positive, valid_ticket_remain
}
