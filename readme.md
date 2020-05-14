# File conversion for NBIRDT

A one shot program to migrate files for publications, data holdings, and projects from 
`./client/client-handers/files` to `./ui/static/site-content/files`


**Important** - run the migrations to add the `slug` field to the three db tables first!


## Usage

1: Generate the binary:

~~~
env GOOS=linux GOARCH=amd64  go build -o migrateDocs *.go
~~~


2: Copy the file to the **root directory** of the NBIRDT application.

3: Run with flags:

~~~bash
./migrateDocs -u username -p password -s ssl
~~~

where `username` and `password` are the db credentials, and `ssl` is the Postgres SSL setting
(e.g. disable).

Once the migration is complete and has been verified, the directory `./client/clientahandlers/files` can
be deleted.

## What Happens

Previously, all documents associated with publications, holdings, and projects were stored outside of
document root for the web application, so as to avoid inadvertent indexing by search engines for anything
that was under embargo.

The files were previously stored in `./client/clienthandlers/files/[id]/[filename]`, where `[id]` was the primary key for
the associated publication/holding/project, and `[filename]` as a randomly generated string, so as to avoid
name collisions.

This application moves all of the files to `./ui/static/site-content-files/files/[category/[name-id]`, where
`[category]` is one of publications, holdings, or projects, and `[name-id]` is a slugified version of the 
original name of the publication/project/holding, followed by its id from the database, so as to avoid name
collisions.
 
The application also renames the files using the 
following logic:

    Old filename: random_string_of_characters.extension
    New filename: slugified-version-of-display-name.original_extension

Example: a holding might have a file with the following characteristics:

    Display Name: 2020 Final Report for Really Important Project.pdf
    Actual file name: dgr1gw54tgrdafsdgyretwgrgreagfr.pdf
    Newly created name: 2020-final-report-for-really-important-project.pdf

Note that the display name will not change, and that the file actually downloaded to the user's system will be the
original display name. We slugify the actual file name so as to avoid problems with accents, non-standard characters, 
etc.