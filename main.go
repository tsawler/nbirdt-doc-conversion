package main

import (
	"database/sql"
	"flag"
	"fmt"
	"github.com/gosimple/slug"
	"io/ioutil"
	"log"
	"os"
	"strings"

	_ "github.com/jackc/pgconn" // need this and next two for pgx
	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
)

var app *application

type application struct {
	db *sql.DB
}

// openDB opens a database connection
func openDB(dsn string) (*sql.DB, error) {
	d, err := sql.Open("pgx", dsn)
	if err != nil {
		panic(err)
	}

	return d, err
}

// main is main app function
func main() {
	// create loggers
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	dbHost := "127.0.01"
	dbPort := "5432"
	databaseName := "nbirdt"

	// read flags
	dbUser := flag.String("u", "", "DB Username")
	dbPass := flag.String("p", "", "DB Password")
	dbSsl := flag.String("s", "disable", "SSL Settings")
	flag.Parse()

	if *dbUser == "" {
		fmt.Println("Missing required flags.")
		os.Exit(1)
	}

	dsn := ""

	if *dbPass == "" {
		dsn = fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5", dbHost, dbPort, *dbUser, databaseName, *dbSsl)
	} else {
		dsn = fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=%s timezone=UTC connect_timeout=5", dbHost, dbPort, *dbUser, *dbPass, databaseName, *dbSsl)
	}

	// open connection to db
	db, err := openDB(dsn)
	if err != nil {
		errorLog.Fatal(err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Println("Pinged database successfully!")

	// create necessary foldres
	CreateDirIfNotExist("./ui/static/site-content/files/holding-documents")
	CreateDirIfNotExist("./ui/static/site-content/files/publication-documents")
	CreateDirIfNotExist("./ui/static/site-content/files/project-documents")

	// populate config
	app = &application{
		db: db,
	}

	err = app.addSlugToHoldings()
	if err != nil {
		errorLog.Fatal(err)
	}

	err = app.addSlugToPublications()
	if err != nil {
		errorLog.Fatal(err)
	}

	err = app.addSlugToProjects()
	if err != nil {
		errorLog.Fatal(err)
	}

	results, err := app.getAllHoldingDocs()
	if err != nil {
		errorLog.Fatal(err)
	}

	infoLog.Println("Starting Holdings")
	infoLog.Println()

	for _, x := range results {
		folderName := slug.Make(x.HoldingNameEn)
		CreateDirIfNotExist(fmt.Sprintf("./ui/static/site-content/files/holding-documents/%s-%d", folderName, x.HoldingID))
		destinationFolder := fmt.Sprintf("%s-%d", folderName, x.HoldingID)
		sourceFile := fmt.Sprintf("./client/clienthandlers/files/holdings/%d/%s", x.HoldingID, x.FileName)
		oldDisplayName := x.FileNameDisplay
		dot := strings.LastIndex(oldDisplayName, ".")
		rootName := oldDisplayName[0:dot]
		last4 := oldDisplayName[dot:len(oldDisplayName)]
		newFileName := fmt.Sprintf("%s%s", slug.Make(rootName), last4)
		destinationFile := fmt.Sprintf("./ui/static/site-content/files/holding-documents/%s/%s", destinationFolder, newFileName)

		// copy file
		input, err := ioutil.ReadFile(sourceFile)
		if err != nil {
			fmt.Println(err)
		}

		err = ioutil.WriteFile(destinationFile, input, 0644)
		if err != nil {
			fmt.Println("Error creating", destinationFile)
			fmt.Println(err)
		}
	}

	pubResults, err := app.getAllPublicationDocs()
	if err != nil {
		errorLog.Fatal(err)
	}

	infoLog.Println("Done Holdings")
	infoLog.Println()
	infoLog.Println("------------------------")
	infoLog.Println()
	infoLog.Println("Starting Publications")
	infoLog.Println()

	for _, x := range pubResults {
		folderName := slug.Make(x.PublicationNameEn)
		CreateDirIfNotExist(fmt.Sprintf("./ui/static/site-content/files/publication-documents/%s-%d", folderName, x.PublicationID))
		destinationFolder := fmt.Sprintf("%s-%d", folderName, x.PublicationID)
		sourceFile := fmt.Sprintf("./client/clienthandlers/files/publications/%d/%s", x.PublicationID, x.FileName)
		oldDisplayName := x.FileNameDisplay
		dot := strings.LastIndex(oldDisplayName, ".")
		rootName := oldDisplayName[0:dot]
		last4 := oldDisplayName[dot:len(oldDisplayName)]
		newFileName := fmt.Sprintf("%s%s", slug.Make(rootName), last4)
		destinationFile := fmt.Sprintf("./ui/static/site-content/files/publication-documents/%s/%s", destinationFolder, newFileName)

		// copy file
		input, err := ioutil.ReadFile(sourceFile)
		if err != nil {
			fmt.Println(err)
		}

		err = ioutil.WriteFile(destinationFile, input, 0644)
		if err != nil {
			fmt.Println("Error creating", destinationFile)
			fmt.Println(err)
		}
	}

	projResults, err := app.getAllProjectDocs()
	if err != nil {
		errorLog.Fatal(err)
	}

	infoLog.Println("Done Publications")
	infoLog.Println()
	infoLog.Println("------------------------")
	infoLog.Println()
	infoLog.Println("Starting Projects")
	infoLog.Println()

	for _, x := range projResults {
		folderName := slug.Make(x.ProjectNameEn)
		CreateDirIfNotExist(fmt.Sprintf("./ui/static/site-content/files/project-documents/%s-%d", folderName, x.ProjectID))
		destinationFolder := fmt.Sprintf("%s-%d", folderName, x.ProjectID)
		sourceFile := fmt.Sprintf("./client/clienthandlers/files/projects/%d/%s", x.ProjectID, x.FileName)
		oldDisplayName := x.FileNameDisplay
		dot := strings.LastIndex(oldDisplayName, ".")
		rootName := oldDisplayName[0:dot]
		last4 := oldDisplayName[dot:len(oldDisplayName)]
		newFileName := fmt.Sprintf("%s%s", slug.Make(rootName), last4)
		destinationFile := fmt.Sprintf("./ui/static/site-content/files/project-documents/%s/%s", destinationFolder, newFileName)

		// copy file
		input, err := ioutil.ReadFile(sourceFile)
		if err != nil {
			fmt.Println(err)
		}

		err = ioutil.WriteFile(destinationFile, input, 0644)
		if err != nil {
			fmt.Println("Error creating", destinationFile)
			fmt.Println(err)
		}
	}

	// update database
	err = app.updateFileNamesForHoldings()
	if err != nil {
		errorLog.Println(err)
	}

	err = app.updateFileNamesForProjects()
	if err != nil {
		errorLog.Println(err)
	}

	err = app.updateFileNamesForPublications()
	if err != nil {
		errorLog.Println(err)
	}

	infoLog.Println("Done!")
}

// CreateDirIfNotExist creates a directory if it does not exist
func CreateDirIfNotExist(path string) {
	const mode = 0755
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, mode)
		if err != nil {
			log.Fatal(err)
		}
	}
}
