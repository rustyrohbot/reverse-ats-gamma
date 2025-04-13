package main

import (
	"bufio"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	createFlag := flag.String("create", "", "Create a new SQLite database with the given name")
	flag.Parse()

	if *createFlag != "" {
		if err := createDatabase(*createFlag); err != nil {
			log.Fatalf("Failed to create database: %v", err)
		}
		fmt.Printf("Database created at %s\n", *createFlag)
		return
	}

	if len(flag.Args()) != 1 {
		fmt.Println("Usage:")
		fmt.Println("  go run main.go -create <databasename>")
		fmt.Println("  go run main.go path/to/db")
		return
	}

	dbPath := flag.Arg(0)
	if _, err := os.Stat(dbPath); os.IsNotExist(err) {
		log.Fatalf("Database file not found: %s", dbPath)
	}

	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		log.Fatalf("Failed to open database: %v", err)
	}
	defer db.Close()

	runMenu(db)
}

func createDatabase(path string) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return err
	}
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return err
	}
	defer db.Close()

	schema := []string{
		`CREATE TABLE IF NOT EXISTS Companies (
			companyID INTEGER PRIMARY KEY,
			name TEXT NOT NULL,
			description TEXT,
			url TEXT,
			hqCity TEXT,
			hqState TEXT
		);`,
		`CREATE TABLE IF NOT EXISTS Roles (
			roleID INTEGER PRIMARY KEY,
			companyID INTEGER NOT NULL,
			name TEXT NOT NULL,
			url TEXT,
			description TEXT,
			coverLetter TEXT,
			applied TEXT,
			appliedDate TEXT,
			postedRangeMin INTEGER,
			postedRangeMax INTEGER,
			equity BOOLEAN,
			workCity TEXT,
			workState TEXT,
			location TEXT,
			status TEXT,
			discovery TEXT,
			referral BOOLEAN,
			notes TEXT,
			FOREIGN KEY (companyID) REFERENCES Companies(companyID)
		);`,
		`CREATE TABLE IF NOT EXISTS Interviews (
			interviewID INTEGER PRIMARY KEY,
			roleID INTEGER NOT NULL,
			date TEXT,
			start TEXT,
			end TEXT,
			notes TEXT,
			type TEXT,
			FOREIGN KEY (roleID) REFERENCES Roles(roleID)
		);`,
		`CREATE TABLE IF NOT EXISTS Contacts (
			contactID INTEGER PRIMARY KEY,
			companyID INTEGER NOT NULL,
			firstName TEXT,
			lastName TEXT,
			role TEXT,
			email TEXT,
			phone TEXT,
			linkedin TEXT,
			notes TEXT,
			FOREIGN KEY (companyID) REFERENCES Companies(companyID)
		);`,
		`CREATE TABLE IF NOT EXISTS InterviewsContacts (
			interviewContactId INTEGER PRIMARY KEY,
			interviewId INTEGER NOT NULL,
			contactId INTEGER NOT NULL,
			FOREIGN KEY (interviewId) REFERENCES Interviews(interviewID),
			FOREIGN KEY (contactId) REFERENCES Contacts(contactID)
		);`,
	}

	for _, stmt := range schema {
		if _, err := db.Exec(stmt); err != nil {
			return err
		}
	}
	return nil
}

func runMenu(db *sql.DB) {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("\nSelect an option:")
		fmt.Println("1. List companies")
		fmt.Println("2. List roles with company name")
		fmt.Println("3. List interviews with role and company")
		fmt.Println("4. List contacts with company")
		fmt.Println("5. Run custom SQL query")
		fmt.Println("6. Exit")
		fmt.Print("> ")

		var choice int
		fmt.Scan(&choice)

		switch choice {
		case 1:
			queryAndPrintTable(db, `
				SELECT companyID, name, description, url, hqCity, hqState FROM Companies
			`)
		case 2:
			queryAndPrintTable(db, `
				SELECT r.roleID, r.name, r.description, r.status, r.postedRangeMin, r.postedRangeMax, c.name AS companyName
				FROM Roles r
				JOIN Companies c ON r.companyID = c.companyID
				ORDER BY r.roleID
			`)
		case 3:
			queryAndPrintTable(db, `
				SELECT i.interviewID, i.date, i.start, i.end, i.notes, i.type,
					   r.name AS roleName, c.name AS companyName
				FROM Interviews i
				JOIN Roles r ON i.roleID = r.roleID
				JOIN Companies c ON r.companyID = c.companyID
				ORDER BY i.interviewID
			`)
		case 4:
			queryAndPrintTable(db, `
				SELECT con.contactID, con.firstName, con.lastName, con.role, con.email, con.phone,
					   con.linkedin, con.notes, c.name AS companyName
				FROM Contacts con
				JOIN Companies c ON con.companyID = c.companyID
				ORDER BY con.contactID
			`)
		case 5:
			fmt.Print("Enter SQL query:\n> ")
			sqlQuery, _ := reader.ReadString('\n')
			sqlQuery = strings.TrimSpace(sqlQuery)
			if strings.HasPrefix(strings.ToLower(sqlQuery), "select") {
				queryAndPrintTable(db, sqlQuery)
			} else {
				_, err := db.Exec(sqlQuery)
				if err != nil {
					fmt.Println("Execution error:", err)
				} else {
					fmt.Println("Query executed.")
				}
			}
		case 6:
			fmt.Println("Exiting.")
			return
		default:
			fmt.Println("Invalid choice")
		}
	}
}

func queryAndPrintTable(db *sql.DB, query string) {
	rows, err := db.Query(query)
	if err != nil {
		log.Printf("Query failed: %v\n", err)
		return
	}
	defer rows.Close()

	cols, err := rows.Columns()
	if err != nil {
		log.Printf("Failed to get columns: %v\n", err)
		return
	}

	fmt.Println()
	fmt.Println(strings.Join(cols, " | "))
	fmt.Println(strings.Repeat("-", len(strings.Join(cols, " | "))))

	vals := make([]interface{}, len(cols))
	valPtrs := make([]interface{}, len(cols))
	for i := range vals {
		valPtrs[i] = &vals[i]
	}

	for rows.Next() {
		if err := rows.Scan(valPtrs...); err != nil {
			log.Printf("Row scan failed: %v\n", err)
			continue
		}

		strRow := make([]string, len(cols))
		for i, val := range vals {
			if val == nil {
				strRow[i] = "NULL"
			} else {
				strRow[i] = fmt.Sprintf("%v", val)
			}
		}
		fmt.Println(strings.Join(strRow, " | "))
	}
}
