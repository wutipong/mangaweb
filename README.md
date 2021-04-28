# mangaweb

## What is mangaweb??
`mangaweb` is a web application which let the user read a manga from zip file without extracting it. Basically it scans for files in a designated directory and display 
it as a web application.

## Why??
I have a local server which downloads and keeps a bunch of stuffs. This server is inaccessible from the internet. Now I need a way to view those contents on my server, 
and manga are among the thing I have. I need something that works with iOS, Android, Windows, etc, so desktop application is kinda out of the window. Also I want something
that simple enough to suite my need without too many steps. 

## Requirements
`mangaweb` require MongoDB to keeps the metadata and stuffs.

## Parameter/Configuration

`mangaweb` takes a few parameter to run. Each parameter is accomapnied by an environment parameter to ease the settings when using it as a docker image.

* `-address`/`MANGAWEB_ADDRESS` is the server address. The default value is `:80` which indicates running the web server at port 80. If you want to change the port number, 
  change the numer **without** removing the `:`. eg `:8080` for port 8080.
  
* `-data`/`MANGAWEB_DATA_PATH` is the path where the manga web looks for file. By default, it looks at the `./data` which contains some test data. Definitely override this value.

* `-database`/`MANGAWEB_DB` is the MongoDB server address.

## Development
The souce code contains `docker-compose` file which has mangaweb, MongoDB and Mongo-Express. To start those service, runs `docker-compose up -d` at the project directory.

The mangaweb service runs at port 8080 and mongo-express at 8081.

You can debug mangaweb inside the container, or run it externally using different port than 8080 (which is already different than the default port 80). 
