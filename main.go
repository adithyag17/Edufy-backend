package main

import (
	"fmt"
	"log"

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
type Admin struct {
	gorm.Model
	Name     string `gorm:"column:name;primaryKey"`
	Email    string `gorm:"column:email"`
	Password string `gorm:"column:password"`
}

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
}

type LoginInput struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	IsAdmin  bool   `json:"isAdmin"`
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
	fmt.Println("Connected to database")

	// Migrate the models
	db.AutoMigrate(&Department{}, &Company{}, &Student{}, &FacultyMentor{}, &University{}, &AppliesFor{}, &GotIn{}, &AppliedForInterview{}, &Admin{}, &Roles{})

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
			if input.IsAdmin {
				admin := Admin{
					Name:     input.Name,
					Email:    input.Email,
					Password: input.Password,
				}
				// Add your logic to interact with the database for admin creation here
				fmt.Println("Admin created:", admin)
				// Return the created admin as JSON
				return c.JSON(admin)
			}

			student := Student{
				SRN:         input.SRN,
				Name:        input.Name,
				Email:       input.Email,
				PhoneNumber: input.PhoneNumber,
				DOB:         input.DOB,
				// DepartmentID: input.DepartmentID,
				Password: input.Password,
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
			var admin Admin
			result := db.Where("email = ? AND password = ?", input.Email, input.Password).First(&admin)
			if result.Error != nil {
				return c.Status(401).SendString("Invalid credentials")
			}
			return c.JSON(admin)
		} else {
			// Student login logic
			var student Student
			result := db.Where("email = ? AND password = ?", input.Email, input.Password).First(&student)
			if result.Error != nil {
				return c.Status(401).SendString("Invalid credentials")
			}
			return c.JSON(student)
		}
	})

	app.Get("/hi", func(c *fiber.Ctx) error {
		return c.SendString("I'm a GET request!")
	})
	app.Post("/updatestudentprofile", func(c *fiber.Ctx) error {
		var student Student
		var input detailstoupdatereq
		db.First(&student, "srn = ?", input.oldsrn)
		student.Name = input.name
		student.Email = input.email
		student.DOB = input.dob
		student.Password = input.password
		student.SRN = input.srn
		student.PhoneNumber = input.phone_number
		student.DepartmentID = input.department_id
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
