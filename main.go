package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Student struct {
	gorm.Model
	SRN          string `gorm:"column:srn;primaryKey"`
	Name         string `gorm:"column:name"`
	Password     string `gorm:"column:password"`
	Email        string `gorm:"column:email"`
	PhoneNumber  string `gorm:"column:phone_number"`
	DepartmentID int    `gorm:"foreignKey:DeptID"`
	DOB          string `gorm:"column:dob"`
	Desc         string `gorm:"column:desc"`
	IsAdmin      bool   `gorm:"column:isadmin"`
}

// For the 'department' table
type Department struct {
	gorm.Model
	DeptID   int    `gorm:"column:dept_id;primaryKey"`
	DeptName string `gorm:"column:dept_name"`
	DeptHead string `gorm:"column:dept_head"`
}

// For the 'company' table
type Company struct {
	gorm.Model
	CID     int     `gorm:"column:cid;primaryKey"`
	Name    string  `gorm:"column:name"`
	Package float64 `gorm:"column:package"`
}

// For the 'faculty_mentor' table
type FacultyMentor struct {
	gorm.Model
	FacultyMentorID int    `gorm:"column:faculty_mentorid;primaryKey"`
	FName           string `gorm:"column:fname"`
	LName           string `gorm:"column:lname"`
	Dept            int    `gorm:"foreignKey:DeptID"`
}

// For the 'university' table
type University struct {
	gorm.Model
	UniversityID int    `gorm:"column:university_id;primaryKey"`
	UName        string `gorm:"column:uname"`
	Dept         int    `gorm:"foreignKey:DeptID"`
	Ranking      int    `gorm:"column:ranking"`
	Location     string `gorm:"column:location"`
}

// For the 'applies_for' table
type AppliesFor struct {
	gorm.Model
	SRN           int `gorm:"column:srn;primaryKey"`
	UNID          int `gorm:"column:unid;primaryKey"`
	RecommendedBy int `gorm:"column:recommended_by;foreignKey:FacultyMentorID"`
}

// For the 'gotin' table
type GotIn struct {
	gorm.Model
	SRN         string `gorm:"column:srn;primaryKey"`
	UNID        int    `gorm:"column:unid;primaryKey"`
	ProgramName string `gorm:"column:program"`
	Joined      string `gorm:"column:joined"`
}

// For the 'applied_for_interview' table
type AppliedForInterview struct {
	gorm.Model
	SRN                 string `gorm:"column:srn;primaryKey"`
	CID                 int    `gorm:"column:cid;primaryKey"`
	DateOfInterview     string `gorm:"column:date_of_interview"`
	InterviewExperience string `gorm:"column:interview_experience"`
	CTC                 int    `gorm:"column:date_of_interview"`
	Selected            bool   `gorm:"column:selected"`
}

// For the 'admin' table
// type Admin struct {
// 	gorm.Model
// 	Name        string `gorm:"column:name;primaryKey"`
// 	Email       string `gorm:"column:email"`
// 	Password    string `gorm:"column:password"`
// 	DOB         string `gorm:"column:dob"`
// 	Phonenumber string `gorm:"column:phonenumber"`
// 	Deptid      int    `gorm:"column:deptid"`
// 	description string `gorm:"column:description"`
// }

// For the 'roles' table
type Roles struct {
	gorm.Model
	RoleName string `gorm:"column:rolename"`
	RoleID   int    `gorm:"column:role_id;primaryKey"`
}

type SignUpRequest struct {
	IsAdmin      bool   `json:"isAdmin"`
	Name         string `json:"name"`
	Email        string `json:"email"`
	SRN          string `json:"srn"`
	PhoneNumber  string `json:"phone_number"`
	DOB          string `json:"dob"`
	DepartmentID string `json:"department_id"`
	Password     string `json:"password"`
	Desc         string `json:"desc"`
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
}

type getStudentDetailsreq struct {
	Name string `json:"name"`
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
	Name          string `json:"name"`
	Email         string `json:"email"`
	Oldsrn        string `json:"srn"`
	Phone_number  string `json:"phone_number"`
	Dob           string `json:"dob"`
	Department_id int    `json:"dept_id"`
	Password      string `json:"password"`
}

func main() {
	// Connect to the database
	db, err := gorm.Open(postgres.Open("user=postgres port=5433 password=Adi2012@$ dbname=edufy sslmode=disable"), &gorm.Config{})
	if err != nil {
		log.Fatalln(err)
	}
	fmt.Println("Connected to database")

	// Migrate the models
	db.AutoMigrate(&Department{}, &Company{}, &Student{}, &FacultyMentor{}, &University{}, &AppliesFor{}, &GotIn{}, &AppliedForInterview{}, &Roles{})
	// if err := db.Exec(`
	// ALTER TABLE applied_for_interviews DROP COLUMN id;
	// ALTER TABLE applies_fors DROP COLUMN id;
	// ALTER TABLE admins DROP COLUMN id;
	// ALTER TABLE companies DROP COLUMN id;
	// ALTER TABLE departments DROP COLUMN id;
	// ALTER TABLE faculty_mentors DROP COLUMN id;
	// ALTER TABLE got_ins DROP COLUMN id;
	// ALTER TABLE roles DROP COLUMN id;
	// ALTER TABLE universities DROP COLUMN id;
	// ALTER TABLE students DROP COLUMN id;
	// `).Error; err != nil {
	// 	log.Fatalln(err)
	// }
	//declaring sql functions
	if err := db.Exec(`
        CREATE OR REPLACE FUNCTION avgsal() RETURNS numeric AS $$
        DECLARE
          total_sal numeric;
          count_rows integer;
        BEGIN
          SELECT sum(ctc), count(*) INTO total_sal, count_rows FROM AppliedForInterview;
          
          IF count_rows > 0 THEN
            RETURN total_sal / count_rows;
          ELSE
            RETURN 0;
          END IF;
        END;
        $$ LANGUAGE plpgsql;
    `).Error; err != nil {
		log.Fatalln(err)
	}
	// Declaring SQL functions
	if err := db.Exec(`
CREATE OR REPLACE FUNCTION findtotalplaced() RETURNS integer AS $$
DECLARE
  total_placed integer;
BEGIN
  SELECT count(*) INTO total_placed FROM applied_for_interview WHERE selected = true;
  RETURN total_placed;
END;
$$ LANGUAGE plpgsql;
`).Error; err != nil {
		log.Fatalln(err)
	}

	if err := db.Exec(`
CREATE OR REPLACE FUNCTION findtotalmasterjoined() RETURNS integer AS $$
DECLARE
  total_joined integer;
BEGIN
  SELECT count(*) INTO total_joined FROM got_in WHERE joined = true;
  RETURN total_joined;
END;
$$ LANGUAGE plpgsql;
`).Error; err != nil {
		log.Fatalln(err)
	}

	if err := db.Exec(`
CREATE OR REPLACE FUNCTION totalcollegeoffers() RETURNS integer AS $$
DECLARE
  total_offers integer;
BEGIN
  SELECT count(*) INTO total_offers FROM applied_for_interview UNION ALL SELECT count(*) FROM got_in;
  RETURN total_offers;
END;
$$ LANGUAGE plpgsql;
`).Error; err != nil {
		log.Fatalln(err)
	}
	app := fiber.New()
	app.Use(cors.New())
	app.Post("/signup",
		func(c *fiber.Ctx) error {
			var input SignUpRequest
			if err := c.BodyParser(&input); err != nil {
				return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
					"error": "Invalid request body",
				})
			}

			fmt.Printf("Received SignUp request:\n%+v\n", input)

			// if err != nil {
			// 	// ... handle error
			// 	panic(err)
			// }
			student := Student{
				SRN:         input.SRN,
				Name:        input.Name,
				Email:       input.Email,
				PhoneNumber: input.PhoneNumber,
				DOB:         input.DOB,
				Desc:        input.Desc,
				Password:    input.Password,
				IsAdmin:     input.IsAdmin,
			}
			result := db.Create(&student)
			if result.Error != nil {
				return c.Status(500).SendString(result.Error.Error())
			}
			// Add your logic to interact with the database for student creation here
			fmt.Println("Student created:", student)
			// Return the created student as JSON
			return c.JSON(student)
			// Add your logic to process the SignUp request here
		})
	app.Post("/login", func(c *fiber.Ctx) error {
		var input LoginInput
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		if input.IsAdmin {
			// Admin login logic
			var admin Student
			result := db.Where("email = ? AND password = ?", input.Email, input.Password).First(&admin)
			if result.Error != nil {
				return c.Status(401).SendString("Invalid credentials")
			}
			return c.JSON(admin.Name)
		} else {
			// Student login logic
			var student Student
			result := db.Where("email = ? AND password = ?", input.Email, input.Password).First(&student)
			if result.Error != nil {
				return c.Status(401).SendString("Invalid credentials")
			}
			return c.JSON(student.Name)
		}
	})

	app.Get("/hi", func(c *fiber.Ctx) error {
		return c.SendString("I'm a GET request!")
	})
	app.Post("/updatestudentprofile", func(c *fiber.Ctx) error {
		var student Student
		var input detailstoupdatereq
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		db.First(&student, "srn = ?", input.Oldsrn)

		if len(input.Name) != 0 {
			student.Name = input.Name
		}
		if len(input.Email) != 0 {
			student.Email = input.Email
		}
		if len(input.Dob) != 0 {
			student.DOB = input.Dob
		}
		if len(input.Password) != 0 {
			student.Password = input.Password
		}
		if len(input.Phone_number) != 0 {
			student.PhoneNumber = input.Phone_number
		}
		i := strconv.Itoa(input.Department_id)
		if len(i) != 0 {
			student.DepartmentID = input.Department_id
		}

		db.Save(&student)
		return c.JSON(fiber.Map{
			"status":  "success",
			"message": "Student profile updated successfully",
			"data":    student,
		})
	})
	app.Post("/getDetails", func(c *fiber.Ctx) error {
		var input getStudentDetailsreq
		var student Student
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		result := db.Where("name = ?", input.Name).First(&student)
		if result.Error != nil {
			return c.Status(500).SendString(result.Error.Error())
		}
		student.IsAdmin = true
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
		return uint(company.CID), nil
	}
	//func to find uid using uname
	findUniversityID := func(name string) (uint, error) {
		var university University
		result := db.Where("name = ?", name).First(&university)
		if result.Error != nil {
			return 0, result.Error
		}
		return uint(university.UniversityID), nil
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
			if !input.isAdmin {
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
	app.Get("/findavgsalary", func(c *fiber.Ctx) error {
		// Execute the SQL function to find the average salary
		var result struct {
			AverageSalary float64 `json:"average_salary"`
		}
		if err := db.Raw("SELECT avgsal() AS average_salary").Scan(&result).Error; err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(result)
	})
	// Routes
	app.Get("/findtotalplaced", func(c *fiber.Ctx) error {
		// Execute the SQL function to find the total placed students
		var result struct {
			TotalPlaced int `json:"total_placed"`
		}
		if err := db.Raw("SELECT findtotalplaced() AS total_placed").Scan(&result).Error; err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(result)
	})

	app.Get("/findtotalmasterjoined", func(c *fiber.Ctx) error {
		// Execute the SQL function to find the total students joined universities
		var result struct {
			TotalJoined int `json:"total_joined"`
		}
		if err := db.Raw("SELECT findtotalmasterjoined() AS total_joined").Scan(&result).Error; err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(result)
	})

	app.Get("/totalcollegeoffers", func(c *fiber.Ctx) error {
		// Execute the SQL function to find the total college offers
		var result struct {
			TotalOffers int `json:"total_offers"`
		}
		if err := db.Raw("SELECT totalcollegeoffers() AS total_offers").Scan(&result).Error; err != nil {
			return c.Status(500).SendString(err.Error())
		}
		return c.JSON(result)
	})
	// Start the server
	log.Fatal(app.Listen(":9090"))
}
