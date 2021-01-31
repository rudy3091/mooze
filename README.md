# Mooze (**UNDER DEVELOPMENT**)
A command-line REST api test tool  
  
## Features available
- Responsive terminal UI
- Http request - GET, POST  
  
![0-0-1-image](./asset/image/0-0-1.gif)
  
## Install
mooze requires go to be installed  
only supports linux system for now (tested under WSL-Ubuntu, windows terminal)  
```
$ git clone https://github.com/RudyPark3091/mooze.git ~/.mooze
$ cd ~/.mooze
$ go build
$ sudo ln -s ./mooze /usr/bin/mooze
$ mooze
```
  
## Keybindings
- \[u\]: Enter input mode for target url  
- \[m\]: Enter input mode for http method  
- \[b\]: Enter input mode for request body(as json)  
- \[Ctrl + s\]: Send request  
- \[q\]: Exit application  
  
## TODOS
Add test codes  
Check response timeout  
Let error don't kill application  
Navigate through response text  
Support another methods - PUT, PATCH, DELETE ...  
Add history mngment (with additional tui)  
Add bulk request feature  
