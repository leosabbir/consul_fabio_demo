package main

import (
	"bufio"
	"consul_fabio_demo/utility"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"strings"
)

var reader = bufio.NewReader(os.Stdin)
var fabiopath = "http://13.59.178.243:9999"

func main() {
	c := choice()

	for true {
		if c == 1 {
			getUser()
		} else if c == 2 {
			createUser()
		} else if c == 3 {
			deleteUser()
		} else if c == 4 {
			orgUsers()
		} else if c == 5 {
			users()
		} else if c == 6 {
			break
		}
		c = choice()
	}
}

func choice() int {
	fmt.Println("Enter your choice:")
	fmt.Println("1. Get User")
	fmt.Println("2. Create User")
	fmt.Println("3. Delete User")
	fmt.Println("4. Get Organization Users")
	fmt.Println("5. Get All Users")
	fmt.Println("6. Exit")
	fmt.Print("\n\nEnter choice: ")

	ch, _ := reader.ReadString('\n')
	c := ch[0:1]

	switch c {
	case "1":
		return 1
	case "2":
		return 2
	case "3":
		return 3
	case "4":
		return 4
	case "5":
		return 5
	case "6":
		return 6
	default:
		return 0
	}
}

func getUser() {
	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	res, _, _ := utility.SendRequest(fmt.Sprintf("%s/service/user/%s", fabiopath, clean(email)), http.MethodGet, nil, nil)

	fmt.Println(string(*res))
}

func createUser() {
	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	fmt.Print("Enter first name: ")
	fname, _ := reader.ReadString('\n')
	fmt.Print("Enter last name: ")
	lname, _ := reader.ReadString('\n')
	fmt.Print("Enter organization: ")
	org, _ := reader.ReadString('\n')
	fmt.Print("Enter title: ")
	title, _ := reader.ReadString('\n')
	fmt.Print("Enter address: ")
	address, _ := reader.ReadString('\n')

	data := url.Values{}
	data.Set("email", clean(email))
	data.Set("fname", clean(fname))
	data.Set("lname", clean(lname))
	data.Set("organization", clean(org))
	data.Set("title", clean(title))
	data.Set("address", clean(address))

	header := http.Header{}
	header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, _, _ := utility.SendRequest(fabiopath+"/service/user", http.MethodPost, strings.NewReader(data.Encode()), &header)
	fmt.Println(string(*res))
}

func deleteUser() {
	fmt.Print("Enter email: ")
	email, _ := reader.ReadString('\n')
	res, _, _ := utility.SendRequest(fmt.Sprintf("%s/service/user/%s", fabiopath, clean(email)), http.MethodDelete, nil, nil)

	fmt.Println(string(*res))
}

func orgUsers() {
	fmt.Print("Enter organization: ")
	org, _ := reader.ReadString('\n')
	res, _, _ := utility.SendRequest(fmt.Sprintf("%s/service/user/org/%s", fabiopath, clean(org)), http.MethodGet, nil, nil)

	fmt.Println(string(*res))
}

func users() {
	res, _, _ := utility.SendRequest(fmt.Sprintf("%s/service/users", fabiopath), http.MethodGet, nil, nil)
	fmt.Println(string(*res))
}

func clean(in string) string {
	in = strings.Replace(in, "\n", "", -1)
	in = strings.Replace(in, "\r", "", -1)
	return in
}
