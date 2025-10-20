package main

import "github.com/oxmies/oxmies"

// User is a sample model used in examples.
type User struct {
	oxmies.Model
	ID    int    `orm:"primary_key,column:id"`
	Name  string `orm:"column:name"`
	Email string `orm:"column:email"`
}
