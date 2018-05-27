# osvc-rest

An (under development) CLI for using the [Oracle Service Cloud REST API](https://docs.oracle.com/cloud/latest/servicecs_gs/CXSVC/) written in Go.

## TODO 
		
### add functionality for 
		
1. file attachments upload on POST requests
2. osvc-crest-api-access-token
3. osvc-crest-next-request-after
4. Session Authorization
5. OAuth Authorization

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
	$ osvc-rest report --id 176 -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

	Report (with filters and limiting) Example:
	$ osvc-rest report --id 176 --limit 10 --filters '[{"name":"search_ex","values":"returns"}]' -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE
### HTTP Methods

In order to create a resource, you must use the _post_ command to send JSON data to the resource of your choice

	$ osvc-rest post "opportunities" --data '{"name":"TEST"}' -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

In order to fetch a resource, you must use the _get_ command to request JSON data from the resource of your choice
	
	$ osvc-rest get "opportunities/?q=name like 'TEST'" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE
	
In order to update a resource, you must use the _patch_ command to send JSON data to update the resource of your choice

	$ osvc-rest patch "opportunities/5" --data '{"name":"updated NAME for TEST"}'  -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

In order to delete a resource, you must use the _delete_ command to delete the resource of your choice
	
	$ osvc-rest delete "opportunities/5" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

To review the options of what HTTP verbs you can use against a resource, use the _options_ command
	
	$ osvc-rest options "opportunities" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE
## Required Flags:

	  Basic Authentication
	  -u, --username (string)  Username to use for basic authentication
	  -p, --password (string)  Password to use for basic authentication
	  -i, --interface (string) Oracle Service Cloud Interface to connect with

	  Session Authentication
	  -s, --session-auth

	  OAuth Authentication
	  -o, --oauth

## Optional Flags:
	      --demosite           Change the domain from 'custhelp' to 'rightnowdemo'
	  -v, --version (string)   Changes the CCOM version (default "v1.3")
	  -a, --annotate (string)  Adds a custom header that adds an annotation (CCOM version must be set to "v1.4" or "latest"); limited to 40 characters
	      --no-ssl-verify      Turns off SSL verification
	  -e, --exclude-null       Adds a custom header to excludes null from results
	  -s, --suppress-rules     Adds a header to suppress business rules
	  -t, --utcTime            Adds a custom header to return results using Coordinated Universal Time (UTC) format for time (Supported on November 2016+)
	      --debug              Prints request headers for debugging
	      --schema             Sets 'Accept' header to 'application/schema+json'
	
	Use "osvc-rest [command] --help" for more information about a command.
