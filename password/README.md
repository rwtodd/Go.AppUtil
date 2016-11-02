# terminal/password Go package 

A Go package to read passwords from the current
TTY or console, without echo.  It opens the tty/console
directly, so that it can be used even if `stdin` and `stdout`
have been redirected.

Some of the code was shamelessly lifted and modified from 
golang.org/x/crypto/ssh/terminal.  However, the code there
wasn't sufficient becuase it does not open a fresh handle
to the tty/console.

## Go Get It

you can use this in your code with a simple:

    go get github.com/rwtodd/apputil/password

and a:

    import "github.com/rwtodd/apputil/password"

    ...

    // ask them to type it twice for confirmation:
    pw, err := password.Read("Password: ", 2) 

