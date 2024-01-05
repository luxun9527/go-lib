package main

import (
	"fmt"
	"github.com/jinzhu/copier"
	"testing"
)

func TestCopier(t *testing.T) {

	// Tags in the destination Struct provide instructions to copier.Copy to ignore
	// or enforce copying and to panic or return an error if a field was not copied.
	type Employee struct {
		// Tell copier.Copy to panic if this field is not copied.
		Name string `copier:"must"`

		// Tell copier.Copy to return an error if this field is not copied.
		Age int32 `copier:"must,nopanic"`

		// Tell copier.Copy to explicitly ignore copying this field.
		Salary int `copier:"-"`

		DoubleAge  int32
		EmployeeId int64 `copier:"EmployeeNum"` // specify field name
		SuperRole  string
	}

	type User struct {
		Name         string
		Role         string
		Age          int32
		EmployeeCode int64 `copier:"EmployeeNum"` // specify field name

		// Explicitly ignored in the destination struct.
		Salary int
	}
	var (
		user      = User{Name: "Jinzhu", Age: 18, Role: "Admin", Salary: 200000}
		users     = []User{{Name: "Jinzhu", Age: 18, Role: "Admin", Salary: 100000}, {Name: "jinzhu 2", Age: 30, Role: "Dev", Salary: 60000}}
		employee  = Employee{Salary: 150000}
		employees = []Employee{}
	)

	copier.Copy(&employee, &user)

	fmt.Printf("%#v \n", employee)
	// Employee{
	//    varName: "Jinzhu",           // Copy from field
	//    Age: 18,                  // Copy from field
	//    Salary:150000,            // Copying explicitly ignored
	//    DoubleAge: 36,            // Copy from method
	//    EmployeeId: 0,            // Ignored
	//    SuperRole: "Super Admin", // Copy to method
	// }
	
	// Copy struct to slice
	copier.Copy(&employees, &user)

	fmt.Printf("%#v \n", employees)
	// []Employee{
	//   {varName: "Jinzhu", Age: 18, Salary:0, DoubleAge: 36, EmployeeId: 0, SuperRole: "Super Admin"}
	// }

	// Copy slice to slice
	employees = []Employee{}
	copier.Copy(&employees, &users)

	fmt.Printf("%#v \n", employees)

}
