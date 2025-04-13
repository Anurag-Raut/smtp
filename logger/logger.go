package logger

import (
	"log"
	"os"
)

var ServerLogger = log.New(os.Stdout, "Server: ", log.LstdFlags)
var ClientLogger = log.New(os.Stdout, "Client: ", log.LstdFlags)
