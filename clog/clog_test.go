package clog

import "testing"
import "log"

func TestLogger(t *testing.T) {
	defer Stop()

	Debug("Debug")
	Info("Info")
	Warn("Warn")
	Error("Error")
	Fatal("Fatal")

	log.SetFlags(0)
	Setup(nil) // force getting the log.Flags()
	log.Println("Standard log no flags")

	log.SetFlags(log.Ldate)
	Setup(nil) // force getting the log.Flags()
	log.Println("Standard log date flag")

	log.SetFlags(log.Ltime)
	Setup(nil) // force getting the log.Flags()
	log.Println("Standard log time flag")
	
	log.SetFlags(log.Ltime|log.Lmicroseconds)
	Setup(nil) // force getting the log.Flags()
	log.Println("Standard log time+micro flag")
	
	log.SetFlags(log.Lmicroseconds)
	Setup(nil) // force getting the log.Flags()
	log.Println("Standard log time flag")

	log.SetFlags(log.Llongfile)
	Setup(nil) // force getting the log.Flags()
	log.Println("Standard log longfile flag")

	log.SetFlags(log.Lshortfile)
	Setup(nil) // force getting the log.Flags()
	log.Println("Standard log shortfile flag")

	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
	Setup(nil) // force getting the log.Flags()
	log.Println("Standard log shortfile+Ldate+Ltime flag")
}
