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
	UniversityID string `gorm:"column:university_id;primaryKey"`
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
	UName       string `gorm:"column:uname;primaryKey"`
	ProgramName string `gorm:"column:program"`
	Joined      bool   `gorm:"column:joined"`
	Desc        string `gorm:"column:desc"`
}

// For the 'applied_for_interview' table
type AppliedForInterview struct {
	gorm.Model
	SRN                 string `gorm:"column:srn;primaryKey"`
	CName               string `gorm:"column:cname;primaryKey"`
	DateOfInterview     string `gorm:"column:date_of_interview"`
	InterviewExperience string `gorm:"column:interview_experience"`
	TestExperience      string `gorm:"column:test_experience"`
	TestQualified       bool   `gorm:"column:test_qualified"`
	InternOrPlaced      string `gorm:"column:intern_or_placed"`
	CTC                 string `gorm:"column:ctc"`
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
	Srn              string `json:"srn"`
	Company_name     string `json:"company_name"`
	DateOfInterview  string `json:"doi"`
	Testexpr         string `json:"testexpr"`
	Internorplaced   string `json:"iorp"`
	Selectedforint   bool   `json:"selectedforint"`
	Inter_experience string `json:"intexp"`
	CTC              string `json:"ctc"`
	Selected         bool   `json:"selected"`
}

type addstudentmastersreq struct {
	Srn         string `json:"srn"`
	Uname       string `json:"uname"`
	Programname string `json:"pname"`
	Joined      bool   `json:"joined"`
	Desc        string `json:"desc"`
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
	CREATE OR REPLACE FUNCTION avgsal()
	RETURNS NUMERIC AS $$
	DECLARE
		total_ctc NUMERIC;
		num_rows INT;
	BEGIN
		-- Calculate the sum of the CTC column
		SELECT SUM(ctc) INTO total_ctc FROM applied_for_interviews;
	
		-- Get the number of rows in the table
		SELECT COUNT(*) INTO num_rows FROM applied_for_interviews;
	
		-- Avoid division by zero
		IF num_rows = 0 THEN
			RETURN 0;
		END IF;
	
		-- Calculate the average and return it
		RETURN total_ctc / num_rows;
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
	//db triggers:
	if err := db.Exec(`
	CREATE OR REPLACE FUNCTION check_srn_before_insert() RETURNS TRIGGER AS $$
BEGIN
   IF NEW.SRN NOT IN (SELECT SRN FROM Student) THEN
      RAISE EXCEPTION 'SRN does not exist in Student table';
   END IF;
   RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE  TRIGGER check_srn_before_insert
BEFORE INSERT ON GOT_INS
FOR EACH ROW EXECUTE FUNCTION check_srn_before_insert();

`).Error; err != nil {
		log.Fatalln(err)
	}
	//db trigger check for placement
	//db triggers:
	if err := db.Exec(`
	 CREATE OR REPLACE FUNCTION check_srn_before_insert_into_appliedforinterview() RETURNS TRIGGER AS $$
 BEGIN
	IF NEW.SRN NOT IN (SELECT SRN FROM Student) THEN
	   RAISE EXCEPTION 'SRN does not exist in Student table';
	END IF;
	RETURN NEW;
 END;
 $$ LANGUAGE plpgsql;
 
 CREATE OR REPLACE TRIGGER check_srn_before_insert
 BEFORE INSERT ON applied_for_interviews
 FOR EACH ROW EXECUTE FUNCTION check_srn_before_insert_into_appliedforinterview();
 
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

		// Admin login logic
		var admin Student
		result := db.Where("email = ? AND password = ?", input.Email, input.Password).First(&admin)
		if result.Error != nil {
			return c.Status(401).SendString("Invalid credentials")
		}
		return c.JSON(admin.Name)

		// Student login logic
		// var student Student
		// result := db.Where("email = ? AND password = ?", input.Email, input.Password).First(&student)
		// if result.Error != nil {
		// 	return c.Status(401).SendString("Invalid credentials")
		// }
		// return c.JSON(student.Name)

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
		student.SRN = input.Oldsrn
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
		return c.JSON(student)

	})
	type getall struct {
		Srn string `json:"srn"`
	}
	app.Post("/getstudentplacements", func(c *fiber.Ctx) error {
		var getall1 getall
		if err := c.BodyParser(&getall1); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		var interviews []AppliedForInterview
		result := db.Find(&interviews, "srn = ?", getall1.Srn)

		if result.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Record not found!"})
		}
		fmt.Println(interviews)

		return c.JSON(fiber.Map{"data": interviews})
	})
	app.Post("/getstudentuniversities", func(c *fiber.Ctx) error {
		var getall1 getall
		if err := c.BodyParser(&getall1); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		var universities []GotIn
		result := db.Find(&universities, "srn = ?", getall1.Srn)
		if result.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Record not found!"})
		}
		fmt.Println(universities)
		return c.JSON(fiber.Map{"data": universities})
	})
	app.Get("/allstudentplacementdetails", func(c *fiber.Ctx) error {
		var allstudents []AppliedForInterview
		result := db.Find(&allstudents)
		if result.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Record not found!"})
		}
		fmt.Println(result)
		return c.JSON(fiber.Map{"data": allstudents})
	})
	app.Get("/allstudentsuniversitydetails", func(c *fiber.Ctx) error {
		var allstudents []GotIn
		result := db.Find(&allstudents)
		if result.Error != nil {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Record not found!"})
		}
		fmt.Println(result)
		return c.JSON(fiber.Map{"data": allstudents})
	})
	//func to find cid using cname
	// findCompanyID := func(name string) (uint, error) {
	// 	var company Company
	// 	result := db.Where("Name = ?", name).First(&company)
	// 	if result.Error != nil {
	// 		return 0, result.Error
	// 	}
	// 	return uint(company.CID), nil
	// }
	app.Post("/addstudentplacements", func(c *fiber.Ctx) error {
		var input addplacementreq
		if err := c.BodyParser(&input); err != nil {
			fmt.Println(input)
			return c.Status(400).SendString(err.Error())
		}
		//cname := input.Company_name
		var Data AppliedForInterview
		//companyID, err := findCompanyID(cname)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			return err
		}
		Data = AppliedForInterview{
			SRN:                 input.Srn,
			CName:               input.Company_name,
			DateOfInterview:     input.DateOfInterview,
			InterviewExperience: input.Inter_experience,
			InternOrPlaced:      input.Internorplaced, //1 for placed 2 for intern
			TestQualified:       input.Selectedforint,
			TestExperience:      input.Testexpr,
			CTC:                 input.CTC, // in lpa
		}
		result := db.Create(&Data)
		if result.Error != nil {
			return c.Status(500).SendString(result.Error.Error())
		}
		return c.JSON(Data)
	})

	app.Post("/addstudentmasters", func(c *fiber.Ctx) error {
		var input addstudentmastersreq
		if err := c.BodyParser(&input); err != nil {
			return c.Status(400).SendString(err.Error())
		}
		var Data GotIn
		if input.Joined {
			Data = GotIn{
				SRN:         input.Srn,
				UName:       input.Uname,
				ProgramName: input.Programname,
				Desc:        input.Desc,
				Joined:      input.Joined,
			}
			result := db.Create(&Data)
			if result.Error != nil {
				return c.Status(500).SendString(result.Error.Error())
			}
		}
		var student Student
		result := db.Where("srn = ?", input.Srn).First(&student)
		if result.Error != nil {
			fmt.Println(result.Error)
		}
		return c.JSON(student.Name)
	})
	app.Get("/findavgsalary", func(c *fiber.Ctx) error {
		// Execute the SQL function to find the average salary
		var result struct {
			AverageSalary float64 `json:"average_salary"`
		}
		if err := db.Raw("SELECT avgsal() AS average_salary").Scan(&result).Error; err != nil {
			return c.Status(500).SendString(err.Error())
		}
		fmt.Println(result)
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
