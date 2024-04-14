# Go File Thing
## Client-Server app

It is not a real application. It was made just to play with golang and tcp connection.
Working directory with files and configs is ~/GFS/ in your user home directory 

# How to use
## Server 
``` start-server [port] ``` - will start listening specified port (1109 is no port specified)

## Client 
``` help ``` - will display all the comands
``` set-server <host address> ``` - will set host address in config file
``` ping ``` -  sends ping to specified server
``` list ``` - will list all files stored on server
``` send <path> ``` - will send specified file to server. Path must be absolute
``` get <filename> ``` - will download specified file from server. Filename must be short, like files from ``` list ``` output
``` delete <filename> ``` - will delete file from server. Filename must be short

