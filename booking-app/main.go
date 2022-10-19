/* to run main.go file: in terminal:
		go run main.go
- to run package that uses multiple files: in terminal:
		go run file1.go file2.go
- to run package that uses multiple files: in terminal:
		go run .
		- dot says to run all files in the current folder
- to break/stop application in termal: Ctrl + C
*/

// Misc learnings:
/*
- in most programming languages can only return 1 value from a function, in GO you can do any # (ex: validate_user_input call in main() function)
- 3 levels of variable scope: global, package, local
	- local: within specific function or block of code
	- package: declaration outside of all functions
	- global: can be used across packages --> have to capitalize the function/variable to make it global
*/

/* Arrays:
- arrays have a fixed size in GO
- cannot mix data types in an array
var bookings = [conference_tickets]string{}
	- var bookings = [max # of elements] data type {list of elements or blank if none}
	- could also define empty array with:
	var bookings = [conference_tickets] string
*/

/* Slices:
- Like arrays but don't need predefined size --> more efficient and dynamic
*/

/* For Loops:
- Go only has for loops, but can accomplish all the Python functionality with a for loop
- an infinite for loop:
	- for {} ; or
	- for true {}
- a definite for loop:
	- for condition1 && condition2 && condition^n {}
- each for loop gives you an index and an element in each loop
*/

/* Switch statements:
- similar to if/else just a bit cleaner
- ex code:
	city := "London"

	switch city {
		case "New York":
			// execute code for NY
		case "London", "Berlin", "San Francisco":
			// execute code for London OR Berlin OR San Francisco
		default:
			// execute code if none of the switches
	}
*/

/* What is a pointer:
- when we have a variable and a value the value is stored in memory
- the program needs to know where the value is stored
- a 'pointer' is a variable that points to the value of the variable
	- This way if the variable has multiple different values associated to it, the pointer can correctly map the variable to the instance of the value for that instance
fmt.Println(remaining_tickets)  // print's value of variable
fmt.Println(&remaining_tickets) // print's the pointer (memory location) of the variable
*/

/* Maps
- data type that allows storing multiple key-value pairs per item (ex: if want to store all of persons contact info and map to that person)
- all keys have the same data type, all values have the same data type
- syntax: var var_name = map[<key data type>]<value data type>
*/

// All code must belong to a package, the 1st statement in Go file must be "package ..."
package main // standard name for main package

// Go programs are organized in packages
/* Go's standard library has different core packages ('fmt' is one of them)
'fmt' package, different functions for input/output:
- print messages
- collect user input
- write into a file
*/
import (
	"booking-app/helper" // import "helper" package from the "booking-app" module
	"fmt"                // have to import packages you use the functions of ('fmt' contains the 'print' function)
	"sync"
	"time"
)

/* package level variables:
- can be accessed inside any of the functions and all files in the package
- have to be defined with syntax: type name = value
- !!! best practice is to define variables as a locally as possible --> probably don't want all variables to be package level
*/

const conference_tickets = 50         // declare constant --> constants are like variables but cannot be changed later in program.
var conference_name = "Go Conference" // declare a variable for the conference name
var remaining_tickets uint = 50       // can define type explicitly --> make sense if auto-detection is incorrect
// uint data type does not allow negative #s
/* can also assign variables like so:
conference_tickets := 50
- cannot declare constants
- cannot explicitly define a data type
*/

// var bookings = []string{} // List of all bookings using slice
// var bookings = make([]map[string]string, 0) // create empty list of maps --> that's why the [] is before "map". Define initial size of zero, but slice will expand dynamically
var bookings = make([]user_data, 0) // create empty list of user_data structs

// used to store list of varying data types
// creating custom data type "user_data"
type user_data struct {
	first_name   string
	last_name    string
	user_email   string
	user_tickets uint
}

/* Could also do:
var bookings []string{}
bookings := []string{}
*/

// waits for launched go routine to finish.
/* ex: if removed for loop the whole ticket booking might be complete before the "send_ticket()" function
- this would cause the application to close before the send_ticket() function is fully executed
*/
var wg = sync.WaitGroup{}

// GO needs to know where it starts executing code (entry point)
// For 1 GO application you will have 1 main function because you can only have 1 entry point
func main() {

	greet_users()

	for {
		first_name, last_name, user_email, user_tickets := get_user_input()

		// user validation
		// define variables returned from the function
		valid_name, valid_email, valid_ticket_positive, valid_ticket_remain := helper.Validate_user_input(first_name, last_name, user_email, user_tickets, remaining_tickets)

		if valid_name && valid_email && valid_ticket_positive && valid_ticket_remain {

			book_ticket(first_name, last_name, user_email, user_tickets) // updates remaining_tickets variable and the bookings slice variable

			wg.Add(1)                                                       // add # of goroutines to wait for before executing a new thread
			go send_ticket(first_name, last_name, user_email, user_tickets) // the "go" keyword enables concurrency. Execution of all other code continues independently of the send_ticket()

			first_names := get_first_names() // defining the functions return value as variable "first_names" so that we can call it below
			fmt.Printf("List of first names of ticket holders: %v\n", first_names)
			/*
				fmt.Printf("First ticket holder: %v\n", bookings[0]) // retrieving value for array or slice need to index list
				fmt.Printf("Data type of bookings slice: %T\n", bookings)
				fmt.Printf("Length of bookings slice: %v\n", len(bookings))
			*/

			// var no_tickets_remain bool = remaining_tickets == 0 --> could make a boolean variable
			if remaining_tickets == 0 {
				fmt.Println("Our conference is fully booked. Sorry.")
				break // ends the loop
			}

		} else { // since just a chain of if's below, will execute all of them
			if !valid_name { // the "!" negates, so is saying if invalid name
				fmt.Printf("First name and last name must each be at least 2 characters long.\n")
			}
			if !valid_email {
				fmt.Printf("Email address must contain an @\n")
			}
			if !valid_ticket_positive {
				fmt.Printf("Number of tickets must be greater than zero.\n")
			}
			if !valid_ticket_remain {
				fmt.Printf("We only have %v tickets remaining, so you can't book %v tickets.\n", remaining_tickets, user_tickets)
			}
			//break // if true then will stop the program here and not continue with below.
			//continue // will skip all below and move to next iteration of loop
		}
	}
	wg.Wait() // blocks continued execution until WaitGroup counter is zero
}

func greet_users() { // don't have to use same variable name as in main() function, but has to be exact if variable is called from package level
	// also if inputs are a package level variable don't have to list them can leave blank
	//fmt.Printf("Variable types: conference_name = %T, conference_tickets = %T, remaining_tickets = %T\n", conference_name, conference_tickets, remaining_tickets)
	fmt.Println("Welcome to", conference_name, "booking application")                                                 // Prints with a new line following, 'Print()' just prints
	fmt.Printf("There are a total of %v tickets and %v are still available\n", conference_tickets, remaining_tickets) // can concat text with variables using 'Printf()' and '%v' to signify variable placeholder
	fmt.Println("Get your tickets here to attend")
}

func get_first_names() []string { //syntax: func func_name(inputs and their types, , ,) outputs and their types {code}
	first_names := []string{} // slice variable to store list of only first names
	// range allows iteration over elements of a different data structure (for array/slices provides the index and value)
	/* get an error if code is: for index, booking := range bookings {
	- this is because we aren't using the index
	- can remove error by using blank identifier '_' --> used for variables you don't want to use
	*/
	for _, booking := range bookings {
		first_names = append(first_names, booking.first_name)
		// first_names = append(first_names, booking["first_name"]) // when bookings variable was a map

		/* when "bookings" was a slice and not a map
		// Fields() function in the string package splits strings at space and returns a slice with the split elements
		// ex of Fields(): "Noah Smith" string would return ["Noah", "Smith"] slice
		var names = strings.Fields(booking)         // Split whole name into a 2 element sice of [first name, last name]
		first_names = append(first_names, names[0]) // add first name in the names slice to the first_names variable
		*/
	}
	return first_names // when use return function need to list in the functions output parameters
}

// doesn't need any input parameters because we are asking for input from the user
func get_user_input() (string, string, string, uint) {
	// Go is a statically typed language --> need to specify datatype with variable (can be implied if on same line)
	var first_name string
	var last_name string
	var user_email string
	var user_tickets uint

	// ask user for their name, email, and quantity of tickets
	fmt.Println("Enter your first name:")
	fmt.Scan(&first_name) // save user's input into first_name variable in memory --> use '&' for pointer
	fmt.Println("Enter your last name:")
	fmt.Scan(&last_name)
	fmt.Println("Enter your email address:")
	fmt.Scan(&user_email)
	fmt.Println("Enter your quantity of tickets:")
	fmt.Scan(&user_tickets) // If float inputted rounds down, if non-conforming data type returns zero

	return first_name, last_name, user_email, user_tickets
}

func book_ticket(first_name string, last_name string, user_email string, user_tickets uint) {
	remaining_tickets = remaining_tickets - user_tickets // to perform calculations, variable must have the same type. Also, redefining variable.

	// create struct for each user
	var user_data = user_data{
		first_name:   first_name,
		last_name:    last_name,
		user_email:   user_email,
		user_tickets: user_tickets,
	}

	/* when bookings variable was a map
	// create a map for a user
	var user_data = make(map[string]string) // use make() to actually create the map of the defined type
	// define map keys and values
	user_data["first_name"] = first_name
	user_data["last_name"] = last_name
	user_data["user_email"] = user_email
	user_data["user_tickets"] = strconv.FormatUint(uint64(user_tickets), 10) // convert uint64() where the base is "10" (base 10 just means normal decimal #)
	*/

	bookings = append(bookings, user_data)
	//fmt.Printf("List of bookings is %v\n", bookings) // print out list of bookings after each booking

	/* when bookings variable was a slice
	// to add new booking to an array need to indicate index: bookings[0] = first_name + " " + last_name
	// to add new booking to a slice will just append element:
	//bookings = append(bookings, first_name+" "+last_name)
	*/

	fmt.Printf("Thank you %v %v for booking %v tickets. You will receive a confirmation at %v.\n", first_name, last_name, user_tickets, user_email)
	fmt.Printf("%v tickets remaining for %v\n", remaining_tickets, conference_name)
}

func send_ticket(first_name string, last_name string, user_email string, user_tickets uint) {
	time.Sleep(10 * time.Second) // simulate waiting 10 seconds
	// Can't save "Printf" into a variable, could only return the same results w/ "Printf" by concating at every variable inputted
	var ticket = fmt.Sprintf("%v tickets for %v %v", user_tickets, first_name, last_name)
	fmt.Println("############") // just a visual divider
	fmt.Printf("Sending ticket: \n%v \nto email address %v\n", ticket, user_email)
	fmt.Println("############")
	wg.Done() // removes the wg.Add() from the Wait Group counter --> so once completed will zero the counter out
}
