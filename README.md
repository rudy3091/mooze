# Mooze (**UNDER DEVELOPMENT**)
A command-line REST api test tool  
  
## Preview
  
![0-0-2-image](./asset/image/0-0-2.gif)
  
## Install
mooze requires go to be installed  
only supports linux system for now (tested under WSL-Ubuntu, windows terminal)  
```
$ git clone https://github.com/rudy3091/mooze.git ~/.mooze
$ cd ~/.mooze
$ go build
$ sudo ln -s ./mooze /usr/bin/mooze
$ mooze
```
  
## Keybindings
- \[u\]: Enter input mode for target url  
- \[m\]: Enter input mode for http method  
- \[b\]: Enter input mode for request body(as json)  
- \[s\]: Send request  
- \[q\]: Quit application  
  
## TODOS
Add test codes  
Check response time  
Add history mngment (with additional tui)  
Add bulk request feature  
