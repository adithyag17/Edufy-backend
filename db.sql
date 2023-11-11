-- Create the DEPARTMENT table first because other tables reference it
CREATE TABLE department (
    dept_id INT NOT NULL,
    dept_name VARCHAR(255) NOT NULL,
    dept_head VARCHAR(255) NOT NULL,
    PRIMARY KEY (dept_id)
);

-- Create the COMPANY table
CREATE TABLE company (
    cid INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    package NUMERIC,
    PRIMARY KEY (cid)
);

-- Create the STUDENT table
CREATE TABLE student (
    srn INT NOT NULL,
    name VARCHAR(255) NOT NULL,
    password VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    phone_number VARCHAR(255) NOT NULL,
    department_id INT NOT NULL,
    dob DATE,
    role_id INT NOT NULL,
    PRIMARY KEY (srn),
    FOREIGN KEY (department_id) REFERENCES department(dept_id) ON DELETE CASCADE
    FOREIGN KEY (role_id) REFERENCES roles(role_id)
);

-- Create the FACULTY_MENTOR table
CREATE TABLE faculty_mentor (
    faculty_mentorid INT PRIMARY KEY,
    fname VARCHAR(255),
    lname VARCHAR(255),
    dept INT,
    FOREIGN KEY (dept) REFERENCES department(dept_id)
);

-- Create the UNIVERSITY table
CREATE TABLE university (
    university_id INT PRIMARY KEY,
    uname VARCHAR(255),
    dept INT,
    ranking INT,
    location VARCHAR(255),
    FOREIGN KEY (dept) REFERENCES department(dept_id)
);

-- Create the APPLIES_FOR table
CREATE TABLE applies_for (
    srn INT,
    unid INT,
    recommended_by INT,
    PRIMARY KEY (srn, unid),
    FOREIGN KEY (srn) REFERENCES student(srn),
    FOREIGN KEY (unid) REFERENCES university(university_id),
    FOREIGN KEY (recommended_by) REFERENCES faculty_mentor(faculty_mentorid)
);

-- Create the GOTIN table
CREATE TABLE gotin(
    srn INT,
    unid INT,
    joined VARCHAR(255),
    PRIMARY KEY (srn, unid),
    FOREIGN KEY (srn) REFERENCES student(srn),
    FOREIGN KEY (unid) REFERENCES university(university_id)
);

-- Create the APPLIED_FOR_INTERVIEW table
CREATE TABLE applied_for_interview(
    srn INT,
    cid INT,
    company_name VARCHAR(255),
    date_of_interview DATE,
    interview_experience VARCHAR(255),
    selected BIT,  -- 1 for true , 0 for false
    PRIMARY KEY (srn, cid),
    FOREIGN KEY (srn) REFERENCES student(srn),
    FOREIGN KEY (cid) REFERENCES company(cid)
);
CREATE TABLE admin(
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) NOT NULL,
    role_id INT NOT NULL,
    FOREIGN KEY (role_id) REFERENCES roles(role_id)
    password VARCHAR(255) NOT NULL,
)
CREATE TABLE roles(
    rolename VARCHAR(255) NOT NULL,
    role_id INT NOT NULL,
)
