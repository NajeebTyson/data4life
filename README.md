## Coding task - Data4Life

A project in Go Lang which consists consist two sub projects:

* ### Generator
  Generator that creates a file with 10 million random tokens, one per line, each consisting of lowercase letters a-z.

* ### Reader
  Reader that reads the file and stores the tokens in the database. Some tokens will occur multiple time so there must be only one entry of token in the DB.

## Solution

In this solution I have written two executables named *generator* and *reader* in `cmd` directory. The `internal` directory contains the utils, repository and common stuff. `pkg` contains the actual library of reader and generator which does the main task. `script` directory contains the script to create and drop SQL table.

### Token generation

Generator takes **1.01 seconds** to generate and write 10 million tokens of size 7 in to file 

### Token reading

Reader takes aroung **8-11 seconds** to read 10 million the tokens and create in memory map for unique tokens

### Database

I have tried both PostreSQL (SQL) and MongoDB (NoSQL) to store the tokens. Here is the benchmarks of token insertion I have got using both dbs.

| Batch/Go-routines | PostgresSql | MongoDB |
| ----------------- | ----------- | ------- |
| 1000/50           | 60 s        | 23 s    |
| 60,000/40         | 50 s        | 34 s    |
| 5,000/40          | 61 s        | 21 s    |

I have tried with many options of batch size and no of go-routiens but he optimal solution which worked for me is using **MongoDB** with **batch size 5000** and **40 go routines** to insert batches in parallel.

### DB Schema

  - Database name is ***data4life***
  - Collection/table name is ***tokens***
  
* **MongoDB**
  With MongoDB I created a single collection named tokens to store the tokens. The attribute to hold the token is **token**. So one entry looks like this:

  `
  {
      token: string, _id: BSON ID
  }
  `
* **PostgreSQL**
  With PostgreSQL I created a table named tokens with single column named **token**. Here is the schema looks like
  
  | tokens                           |
  | -------------------------------- |
  | token VARCHAR(7) NOT NULL UNIQUE |


## How to build

* To build generator run this command from the root of the project
  
   `go build ./cmd/generator`   

* To build reader run this command
  
  `go build ./cmd/reader`

* To run directly

   `go run ./cmd/generator`

   `go run ./cmd/reader`