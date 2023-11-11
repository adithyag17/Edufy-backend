package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Department struct {
	gorm.Model
	DeptID   int    `gorm:"column:dept_id"`
	DeptName string `gorm:"column:dept_name"`
	DeptHead string `gorm:"column:dept_head"`
}

type Company struct {
	gorm.Model
	CID     int     `gorm:"column:cid"`
	Name    string  `gorm:"column:name"`
	Package float64 `gorm:"column:package"`
}

type Student struct {
	gorm.Model
	SRN          string `gorm:"column:srn;primaryKey"`
	Name         string `gorm:"column:name"`
	Password     string `gorm:"column:password"`
	Email        string `gorm:"column:email"`
	PhoneNumber  string `gorm:"column:phone_number"`
	DepartmentID int    `gorm:"foreignKey:DeptID"`
	DOB          string `gorm:"column:dob"`
	RoleID       int    `gorm:"column:role_id"`
}

type FacultyMentor struct {
	gorm.Model
	FacultyMentorID int    `gorm:"column:faculty_mentorid"`
	FName           string `gorm:"column:fname"`
	LName           string `gorm:"column:lname"`
	Dept            int    `gorm:"foreignKey:DeptID"`
}

type University struct {
	gorm.Model
	UniversityID int    `gorm:"column:university_id"`
	UName        string `gorm:"column:uname"`
	Dept         int    `gorm:"foreignKey:DeptID"`
	Ranking      int    `gorm:"column:ranking"`
	Location     string `gorm:"column:location"`
}

type AppliesFor struct {
	gorm.Model
	SRN           int `gorm:"column:srn;foreignKey:SRN"`
	UNID          int `gorm:"column:unid;foreignKey:UniversityID"`
	RecommendedBy int `gorm:"column:recommended_by;foreignKey:FacultyMentorID"`
}

type GotIn struct {
	gorm.Model
	SRN    int    `gorm:"column:srn;foreignKey:SRN"`
	UNID   int    `gorm:"column:unid;foreignKey:UniversityID"`
	Joined string `gorm:"column:joined"`
}

type AppliedForInterview struct {
	gorm.Model
	SRN                 int    `gorm:"column:srn;foreignKey:SRN"`
	CID                 int    `gorm:"column:cid;foreignKey:CID"`
	DateOfInterview     string `gorm:"column:date_of_interview"`
	InterviewExperience string `gorm:"column:interview_experience"`
	Selected            bool   `gorm:"column:selected"`
}

type Admin struct {
	gorm.Model
	Name     string `gorm:"column:name"`
	Email    string `gorm:"column:email"`
	RoleID   int    `gorm:"column:role_id"`
	Password string `gorm:"column:password"`
}

type Roles struct {
	gorm.Model
	RoleName string `gorm:"column:rolename"`
	RoleID   int    `gorm:"column:role_id"`
}

type SignupInput struct {
	IsAdmin       bool
	name          string
	email         string
	srn           string
	phone_number  string
	dob           string
	department_id int
	password      string
	// Add other fields as needed
}
type LoginInput struct {
	IsAdmin  bool
	Email    string
	Password string
}
type getStudentDetailsreq struct {
	srn string
}

func main() {
	// Connect to the database
	db, err := gorm.Open(postgres.Open("user=postgres port=5433 password=Adi2012@$ dbname=edufy sslmode=disable"), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}

	// Migrate the models
	db.AutoMigrate(&Department{}, &Company{}, &Student{}, &FacultyMentor{}, &University{}, &AppliesFor{}, &GotIn{}, &AppliedForInterview{}, &Admin{}, &Roles{})

	// Create a new Fiber app
	app := fiber.New()

	// Add routes here
	// Define the signup route
	app.Post("/signup", func(c *fiber.Ctx) error {
		var input SignupInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		if input.IsAdmin {
			admin := Admin{
				Name:     input.name,
				Email:    input.email,
				Password: input.password,
			}
			result := db.Create(&admin)
			if result.Error != nil {
				return c.Status(500).SendString(result.Error.Error())
			}
			return c.JSON(admin)
		} else {
			student := Student{
				SRN:          input.srn,
				Name:         input.name,
				Password:     input.password,
				Email:        input.email,
				PhoneNumber:  input.phone_number,
				DepartmentID: input.department_id,
				DOB:          input.dob,
				// Set other fields as needed
			}
			result := db.Create(&student)
			if result.Error != nil {
				return c.Status(500).SendString(result.Error.Error())
			}
			return c.JSON(student)
		}
	})
	// Define the login route
	app.Post("/login", func(c *fiber.Ctx) error {
		var input LoginInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).SendString(err.Error())
		}

		if input.IsAdmin {
			var admin Admin
			result := db.Where("email = ? AND password = ?", input.Email, input.Password).First(&admin)
			if result.Error != nil {
				return c.Status(500).SendString(result.Error.Error())
			}
			return c.JSON(admin)
		} else {
			var student Student
			result := db.Where("email = ? AND password = ?", input.Email, input.Password).First(&student)
			if result.Error != nil {
				return c.Status(500).SendString(result.Error.Error())
			}
			return c.JSON(student)
		}
	})
	app.Get("/getStudentDetails", func(c *fiber.Ctx) error {
		var input getStudentDetailsreq
		var student Student
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		result := db.Where("srn = ?", input.srn).First(&student)
		if result.Error != nil {
			return c.Status(500).SendString(result.Error.Error())
		}
		return c.JSON(student)

	})

	app.Get("/getstudentplacements", func(c *fiber.Ctx) error {
		srn := c.Params("srn")

		var interviews []AppliedForInterview
		result := db.Find(&interviews, "srn = ?", srn)

		if result.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Record not found!"})
		}

		return c.JSON(fiber.Map{"data": interviews})
	})
	app.Get("/getstudentuniversities", func(c *fiber.Ctx) error {
		srn := c.Params("srn")
		var universities []GotIn
		result := db.Find(&universities, "srn = ?", srn)
		if result.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Record not found!"})
		}
		return c.JSON(fiber.Map{"data": universities})
	})
	app.Get("/allstudentplacementdetails", func(c *fiber.Ctx) error {
		var allstudents []AppliedForInterview
		result := db.Find(&allstudents)
		if result.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Record not found!"})
		}
		return c.JSON(fiber.Map{"data": allstudents})
	})
	app.Get("/allstudentsuniversitydetails", func(c *fiber.Ctx) error {
		var allstudents []GotIn
		result := db.Find(&allstudents)
		if result.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Record not found!"})
		}
		return c.JSON(fiber.Map{"data": allstudents})
	})

	// Start the server
	log.Fatal(app.Listen(":9090"))
}
