# kindle-words
CLI tool to interact with Kindle vocabulary builder

Connect the Kindle to your USB port and run the binary.

```
$ chmod +x kindle-words
$ ./kindle-words -h

NAME:
   kindle-words - Provides methods to work with vocabulary builder

USAGE:
   kindle-words [global options] command [command options] [arguments...]

VERSION:
   0.1.0

COMMANDS:
     export       Export words
     delete-book  Delete book
     delete-word  Delete word
     books        show book title
     help, h      Shows a list of commands or help for one command

GLOBAL OPTIONS:
   --verbose      Show more output
   --help, -h     show help
   --version, -v  print the version

```

**Important**: This has been only tested on Mac for the moment.

**Note**: Your Kindle must be restarted to see the changes.
