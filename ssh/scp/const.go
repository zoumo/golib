package scp

import "os"

const (
	//FileToken crate file token
	FileToken = 'C'
	//DirToken create directory token
	DirToken = 'D'
	//TimestampToken timestamp token
	TimestampToken = 'T'
	//EndDirToken end of dir token
	EndDirToken = 'E'
	//WarningToken warning token
	WarningToken = 0x1
	//ErrorToken error token
	ErrorToken = 0x2
	// StatusOK
	StatusOK = 0x0
)

const (
	//DefaultPort default SSH port
	DefaultPort = 22
)

const (
	//DefaultDirMode folder mode default
	DefaultDirMode = os.ModeDir | 0755
)
