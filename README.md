# osvc-rest

An (under development) CLI for using the [Oracle Service Cloud REST API](https://docs.oracle.com/cloud/latest/servicecs_gs/CXSVC/) written in Go.

## Installing Go (for Windows)
[Installation options are available on the Go website](https://golang.org/doc/install#windows).

Make sure that you follow the instructions about adding system variables.

You will want to create a folder exclusively for your Go projects.

[Read more about what you will need to set up to get Go-ing. (I'm sorry for the bad pun)](https://github.com/golang/go/wiki/SettingGOPATH)


## TODO 
		
### add functionality for 
		
1. file attachments 
	1. download
	2.  upload
2. exclude null data
3. osvc-crest-api-access-token
4. osvc-crest-next-request-after
5. Session Authorization
6. OAuth Authorization
	
### update 
1. reports to take filters
2. add validations
   
<!-- ## Installation

    $ cd ..<go projects folder>
    $ git clone https://github.com/rajangdavis/osvc-rest.git
    $ go build

## Compatibility

Go works everywhere; [learn how at this link.](https://dave.cheney.net/2015/08/22/cross-compilation-with-go-1-5) -->

## Usage
	$osvc-rest [command]

## Available Commands:

### Running one or more ROQL queries
Runs one or more ROQL queries and returns parsed results
	
	Single Query Example:
	$ osvc-rest query "DESCRIBE" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE
	
	Multiple Queries Example: (Queries should be wrapped in quotes and space separated)
	$ osvc-rest query "SELECT * FROM INCIDENTS LIMIT 100" "SELECT * FROM SERVICEPRODUCTS LIMIT 100" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

### Running Reports
Runs an analytics report and returns parsed results

	Report (without filters) Example:
	$ osvc-rest report --id 186 -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

### HTTP Methods
HTTP Methods Example:

In order to create a resource, you must use the *post* command to send JSON data to the resource of your choice

	$ osvc-rest post "opportunities" --data '{"name":"TEST"}' -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

In order to fetch a resource, you must use the *get* command to request JSON data from the resource of your choice
	
	$ osvc-rest get "opportunities/?q=name like 'TEST'" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE
	
In order to update a resource, you must use the *patch* command to send JSON data to update the resource of your choice

	$ osvc-rest patch "opportunities/5" --data '{"name":"updated NAME for TEST"}'  -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

In order to delete a resource, you must use the *delete* command to delete the resource of your choice
	
	$ osvc-rest delete "opportunities/5" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE


## Flags:
	  -a, --annotate string    Adds a custom header that adds an annotation
	      --demosite           Change the domain from 'custhelp' to 'rightnowdemo'
	  -e, --exclude-null       Adds a custom header to excludes null from results
	  -h, --help               help for osvc-rest
	  -i, --interface string   Oracle Service Cloud Interface to connect with
	      --no-ssl-verify      Turns off SSL verification
	  -p, --password string    Password to use for basic authentication
	  -s, --suppress-rules     Adds a header to suppress business rules
	  -u, --username string    Username to use for basic authentication
	  -t, --utcTime            Adds a custom header to return results using Coordinated Universal Time (UTC) format for time
	  -v, --version string     Changes the CCOM version (default "v1.3")

	Use "osvc-rest [command] --help" for more information about a command.
