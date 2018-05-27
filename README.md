# osvc-rest

A command line interface for using the [Oracle Service Cloud REST API](https://docs.oracle.com/cloud/latest/servicecs_gs/CXSVC/).

osvc-rest is written in [golang](https://golang.org).

The basic formula for this CLI is the following:
	
	$ osvc-rest <command to run> <something the command needs> <optional flags to change various settings> <some way to authenticate>

The commands come in the following flavors:

1. [HTTP Methods](#http-methods)
2. [Running one or more ROQL queries](running-one-or-more-roql-queries)
3. [Running reports](running-reports)


## TODO 
1. Session Authorization
2. OAuth Authorization

## Authentication:
Use the following flags to authenticate

	  Basic Authentication
	  -u, --username (string)  Username to use for basic authentication
	  -p, --password (string)  Password to use for basic authentication
	  -i, --interface (string) Oracle Service Cloud Interface to connect with

	  Session Authentication
	  -s, --session-auth

	  OAuth Authentication
	  -o, --oauth

## HTTP Methods
All the of HTTP Methods have the following formula:
	
	$ osvc-rest <http-verb> <resource-url> (optional flags for POST and PATCH requests: --data , --attach-file) <authentication-method>

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

## Uploading File Attachments

In order to upload a file attachment, use the --attach-file (or -f) flag to attach a file with the file location

	$ osvc-rest post "opportunities" --data '{"name":"TEST"}' --attach-file "./proof_of_purchase.jpg" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

To attach multiple files, use the --attach-file (or -f) flag for each file you wish to attach

	$ osvc-rest patch "incidents/302" -f "front_angle.png" -f "back_angle.png" -f "side_angle.png" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

## Downloading File Attachments

In order to download a file attachment from a given resource, [add "?download" to the file attachment URL](https://docs.oracle.com/en/cloud/saas/service/18b/cxsvc/c_osvc_managing_file_attachments.html#ManagingFileAttachments-07BABEF6__concept-406-3A92801C). The file will be downloaded in the same directory that the command is run in.

	$ osvc-rest get "incidents/24898/fileAttachments/253?download" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

To download all file attachmentss from a given resource, add ["?download" to the file Attachments URL](https://docs.oracle.com/en/cloud/saas/service/18b/cxsvc/c_osvc_managing_file_attachments.html#ManagingFileAttachments-07BABEF6__concept-410-3A92801F). A file called "downloadedAttachment.tgz" will be downloaded to your computer. 

	$ osvc-rest get "incidents/24898/fileAttachments?download" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

You can extract the file using [tar](https://askubuntu.com/questions/499807/how-to-unzip-tgz-file-using-the-terminal/499809#499809)
    
    $ tar -xvzf ./downloadedAttachment.tgz

## Running one or more ROQL queries
Runs one or more ROQL queries and returns parsed results
	
	Single Query Example:
	$ osvc-rest query "DESCRIBE" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE
	
	Multiple Queries Example: (Queries should be wrapped in quotes and space separated)
	$ osvc-rest query "SELECT * FROM INCIDENTS LIMIT 100" "SELECT * FROM SERVICEPRODUCTS LIMIT 100" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

## Running Reports
Runs an analytics report and returns parsed results

	Report (without filters) Example:
	$ osvc-rest report --id 176 -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

	Report (with filters and limiting) Example:
	$ osvc-rest report --id 176 --limit 10 --filters '[{"name":"search_ex","values":"returns"}]' -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE

## Bulk Delete
This CLI provides a very simple interface to use the Bulk Delete feature within the latest versions of the REST API. Before you can use this feature, make sure that you have the [correct permissions set up for your profile](https://docs.oracle.com/en/cloud/saas/service/18b/cxsvc/c_osvc_bulk_delete.html#BulkDelete-10689704__concept-212-37785F91).

	Bulk Delete Example: 
	$ osvc-rest query "DELETE from incidents limit 1000" "DELETE from incidents limit 1000" -u $OSC_ADMIN -p $OSC_PASSWORD -i $OSC_SITE --demosite -v latest -a "Testing bulk delete multiple requests"

## Optional Flags:
	    --access-token (string) 	Adds an access token to ensure quality of service
	-a, --annotate (string)     	Adds a custom header that adds an annotation (CCOM version must be set to "v1.4" or "latest"); limited to 40 characters
	    --debug                 	Prints request headers for debugging
	    --demosite              	Change the domain from 'custhelp' to 'rightnowdemo'
	-e, --exclude-null          	Adds a custom header to excludes null from results
	-h, --help                  	help for osvc-rest
	    --next-request (int)      	Number of milliseconds before another HTTP request can be made; this is an anti-DDoS measure
	    --no-ssl-verify         	Turns off SSL verification
	    --schema                	Sets 'Accept' header to 'application/schema+json'
	    --suppress-rules        	Adds a header to suppress business rules
	-t, --utc-time              	Adds a custom header to return results using Coordinated Universal Time (UTC) format for time (Supported on November 2016+)
	-v, --version (string)      	Changes the CCOM version (default "v1.3")
	
	Use "osvc-rest [command] --help" for more information about a command.
