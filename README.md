# Go File Thing
## Client-Server app

It is not a real application. It was made just to play with golang and tcp connection.  
Working directory with files and configs is ~/GFS/ in your user home directory  

# How to use
## Server 
``` start-server [port] ``` - start listening specified port (1109 is no port specified)  

## Client  
``` help ``` - display all the comands  
``` set-server <host address> ``` - set host address in config file  
``` ping ``` -  send ping to specified server  
``` list ``` - list all files stored on server  
``` send <path> ``` - send specified file to server. Path must be absolute  
``` get <filename> ``` - download specified file from server. Filename must be short, like files from ``` list ``` output  
``` delete <filename> ``` - delete file from server. Filename must be short  

