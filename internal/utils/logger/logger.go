package logger

import "log"

func Info(message string) {
	log.Printf("[LynxPilot][信息]:%s", message)
}

func Error(message string) {
	log.Printf("[LynxPilot][错误]:%s", message)
}
