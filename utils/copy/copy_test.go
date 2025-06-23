package copy


import (
	"fmt"
	"testing"

	"github.com/jinzhu/copier"
)

// 基础结构体定义
type User struct {
	Name     string
	Age      int
	Address  string
	Nickname string
}

type UserDTO struct {
	Name    string
	Age     int
	Address string
}

// 不同名字段映射
type Employee struct {
	Name       string
	WorkAge    int    `copier:"Age"`      // 通过 tag 映射到不同名字段
	WorkPlace  string `copier:"Address"`
	Department string
}

// 嵌套结构体
type Company struct {
	Name     string
	Location string
	Boss     User
}

type CompanyDTO struct {
	Name     string
	Location string
	Boss     UserDTO
}

// 指针字段
type Department struct {
	Name     string
	Manager  *User
	Location string
}

type DepartmentDTO struct {
	Name     string
	Manager  *UserDTO
	Location string
}

// 切片字段
type Team struct {
	Name    string
	Members []User
}

type TeamDTO struct {
	Name    string
	Members []UserDTO
}

func TestBasicCopy(t *testing.T) {
	user := User{
		Name:     "张三",
		Age:      25,
		Address:  "北京市",
		Nickname: "小张",
	}

	var userDTO UserDTO
	err := copier.Copy(&userDTO, &user)
	if err != nil {
		fmt.Printf("复制出错: %v\n", err)
		return
	}
	
	fmt.Printf("基础复制结果:\n")
	fmt.Printf("姓名: %s\n", userDTO.Name)
	fmt.Printf("年龄: %d\n", userDTO.Age)
	fmt.Printf("地址: %s\n", userDTO.Address)
}

func TestDifferentFieldNameCopy(t *testing.T) {
	user := User{
		Name:    "李四",
		Age:     30,
		Address: "上海市",
	}

	var employee Employee
	err := copier.Copy(&employee, &user)
	if err != nil {
		fmt.Printf("复制出错: %v\n", err)
		return
	}
	
	fmt.Printf("\n不同字段名复制结果:\n")
	fmt.Printf("姓名: %s\n", employee.Name)
	fmt.Printf("工作年限: %d\n", employee.WorkAge)
	fmt.Printf("工作地点: %s\n", employee.WorkPlace)
}

func TestNestedStructCopy(t *testing.T) {
	company := Company{
		Name:     "科技公司",
		Location: "广州市",
		Boss: User{
			Name:    "王五",
			Age:     40,
			Address: "广州市天河区",
		},
	}

	var companyDTO CompanyDTO
	err := copier.Copy(&companyDTO, &company)
	if err != nil {
		fmt.Printf("复制出错: %v\n", err)
		return
	}
	
	fmt.Printf("\n嵌套结构体复制结果:\n")
	fmt.Printf("公司名称: %s\n", companyDTO.Name)
	fmt.Printf("公司地址: %s\n", companyDTO.Location)
	fmt.Printf("老板信息:\n")
	fmt.Printf("  姓名: %s\n", companyDTO.Boss.Name)
	fmt.Printf("  年龄: %d\n", companyDTO.Boss.Age)
}

func TestPointerFieldCopy(t *testing.T) {
	manager := &User{
		Name:    "赵六",
		Age:     35,
		Address: "深圳市",
	}
	
	dept := Department{
		Name:     "研发部",
		Manager:  manager,
		Location: "深圳市南山区",
	}

	var deptDTO DepartmentDTO
	err := copier.Copy(&deptDTO, &dept)
	if err != nil {
		fmt.Printf("复制出错: %v\n", err)
		return
	}
	
	fmt.Printf("\n指针字段复制结果:\n")
	fmt.Printf("部门名称: %s\n", deptDTO.Name)
	fmt.Printf("部门地址: %s\n", deptDTO.Location)
	if deptDTO.Manager != nil {
		fmt.Printf("经理信息:\n")
		fmt.Printf("  姓名: %s\n", deptDTO.Manager.Name)
		fmt.Printf("  年龄: %d\n", deptDTO.Manager.Age)
	}
}

func TestSliceFieldCopy(t *testing.T) {
	team := Team{
		Name: "开发团队",
		Members: []User{
			{Name: "成员1", Age: 25},
			{Name: "成员2", Age: 28},
		},
	}

	var teamDTO TeamDTO
	err := copier.Copy(&teamDTO, &team)
	if err != nil {
		fmt.Printf("复制出错: %v\n", err)
		return
	}
	
	fmt.Printf("\n切片字段复制结果:\n")
	fmt.Printf("团队名称: %s\n", teamDTO.Name)
	fmt.Printf("团队成员数量: %d\n", len(teamDTO.Members))
	for i, member := range teamDTO.Members {
		fmt.Printf("成员 %d:\n", i+1)
		fmt.Printf("  姓名: %s\n", member.Name)
		fmt.Printf("  年龄: %d\n", member.Age)
	}
}
