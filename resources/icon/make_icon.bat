@ECHO OFF

IF "%1"=="" GOTO NOICO
IF NOT EXIST %1 GOTO BADFILE
ECHO Creating iconwin.go
ECHO //+build windows > iconwin.go
ECHO. >> iconwin.go
TYPE %1 | 2goarray Data icon >> iconwin.go
GOTO DONE

:CREATEFAIL
ECHO Unable to create output file
GOTO DONE

:NOICO
ECHO Please specify a .ico file
GOTO DONE

:BADFILE
ECHO %1 is not a valid file
GOTO DONE

:DONE