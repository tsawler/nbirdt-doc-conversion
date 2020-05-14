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
./migrateDocs -u <username> -p <password>
~~~

where `<username>` and `<password>` are the db credentials


Done.