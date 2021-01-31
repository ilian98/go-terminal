# go-terminal
Terminal with basic commands such as pwd, cat, cd and other, implemented on Go.

This terminal has basic functionality like: 
- starting one command from the list: <code> pwd, cd, ls, cat, cp, mv, mkdir, rm, find and ping </code>
- exiting with one command from the list: <code> exit, logout and bye </code>
- running command in background mode (by writing '&')
- make pipe of commands (with the standard '|' between commands)
- command execution can be stopped by Ctrl+C
- standard input and output streams can be redirected to files (by < and > respectively)

Also for escaping certain characters on command line, one can use " " around the property.
When there is some error in parsing or command execution, appropriate messages are given.

To build and run locally:
<pre>
git clone https://github.com/ilian98/go-terminal.git
go run main.go
</pre>
