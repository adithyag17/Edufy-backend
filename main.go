package main

import (
	"fmt"
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
	SRN         string `gorm:"column:srn;foreignKey:SRN"`
	UNID        int    `gorm:"column:unid;foreignKey:UniversityID"`
	ProgramName string `gorm:"column:program"`
	Joined      string `gorm:"column:joined"`
}

type AppliedForInterview struct {
	gorm.Model
	SRN                 string `gorm:"column:srn;foreignKey:SRN"`
	CID                 int    `gorm:"column:cid;foreignKey:CID"`
	DateOfInterview     string `gorm:"column:date_of_interview"`
	InterviewExperience string `gorm:"column:interview_experience"`
	CTC                 int    `gorm:"column:date_of_interview"`
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
type addplacementreq struct {
	srn             string
	company_name    string
	DateOfInterview string
	isAdmin         bool
	experience      string
	CTC             int
	selected        bool
}
type addstudentmastersreq struct {
	srn         string
	uname       string
	programname string
	isAdmin     bool
	Joined      bool
}
type detailstoupdatereq struct {
	name          string
	email         string
	srn           string
	oldsrn        string
	phone_number  string
	dob           string
	department_id int
	password      string
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
	app.Post("/updatestudentprofile", func(c *fiber.Ctx) error {
		var student Student
		var input detailstoupdatereq
		db.First(&student, "srn = ?", input.oldsrn)
		//changing profile
		student.Name = input.name
		student.Email = input.email
		student.DOB = input.dob
		student.Password = input.password
		student.SRN = input.srn
		student.PhoneNumber = input.phone_number
		db.Save(&student)
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Student profile updated successfully",
			"data":    student,
		})
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
	//func to find cid using cname
	findCompanyID := func(name string) (uint, error) {
		var company Company
		result := db.Where("Name = ?", name).First(&company)
		if result.Error != nil {
			return 0, result.Error
		}
		return company.ID, nil
	}
	//func to find uid using uname
	findUniversityID := func(name string) (uint, error) {
		var university University
		result := db.Where("name = ?", name).First(&university)
		if result.Error != nil {
			return 0, result.Error
		}
		return university.ID, nil
	}
	app.Post("/addstudentplacements", func(c *fiber.Ctx) error {
		var input addplacementreq
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		cname := input.company_name
		var Data AppliedForInterview
		companyID, err := findCompanyID(cname)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
		if input.selected {
			if input.isAdmin {
				Data = AppliedForInterview{
					SRN:                 input.srn,
					CID:                 int(companyID),
					DateOfInterview:     input.DateOfInterview,
					InterviewExperience: input.experience,
					CTC:                 input.CTC,
				}
				result := db.Create(&Data)
				if result.Error != nil {
					return c.Status(500).SendString(result.Error.Error())
				}

			}
		}
		return c.JSON(Data)
	})

	app.Post("/addstudentmasters", func(c *fiber.Ctx) error {
		var input addstudentmastersreq
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		uname := input.uname
		var Data GotIn
		UniversityID, err := findUniversityID(uname)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
		if input.Joined {
			if input.isAdmin {
				Data = GotIn{
					SRN:         input.srn,
					UNID:        int(UniversityID),
					ProgramName: input.programname,
				}
				result := db.Create(&Data)
				if result.Error != nil {
					return c.Status(500).SendString(result.Error.Error())
				}
			}
		}

		return c.JSON(Data)
	})

	// Start the server
	log.Fatal(app.Listen(":9090"))
}
