# apputil-go/cmdline

It's current only feature is `cmdline.GlobArgs()`, which 
handles the difference between Unix and Windows shell behavior.
On windows, globs aren't expanded by the shell, so the program will get a
literal "\*.\*" as an argument, rather than a list of matching
files. 

This isn't good for portable behavior between operating systems. So,
`GlobArgs()` will perform expansion on all arguments on Windows, and 
do nothing on other platforms.

## Get 

   go get github.com/rwtodd/apputil/cmdline
