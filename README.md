# Reverse ATS - GAMMA

A gamma prototype for a tool to organize track what roles you've applied to, the companies those roles belong to, and interviews related to the role, and relevant contacts for each company.

## Overview

[Previously](https://github.com/rustyrohbot/reverse-ats-alpha) I got a system working with an Excel workbook that had five sheets: Companies, Roles, Contacts, Interviews, and InterviewContacts.

[Then](https://github.com/rustyrohbot/reverse-ats-beta) I had a Python script migrate the data from Excel into a SQLite database.

Now we are buliding a CLI app on top of the database.


## Requirement

- Go 1.18 or later
- SQLite3 (via Go module: 'github.com/mattn/go-sqlite3')


## Installation

1. Clone this repository:
   ```
   git clone https://github.com/rustyrohbot/reverse-ats-gamma.git
   cd reverse-ats-gamma
   ```

2. Install dependencies:
   ```
   go mod tidy
   ```

## Usage

To create a new empty database, run

```
go run main.go -create database.sqlite
```

to generate a file called `database.sqlite`

To run the full application, run

```
go run main.go database.sqlite
```

to run the applicaiton using `database.sqlite` as a datasource


## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.