@echo off

REM Build the Go binary
go build -o act.exe main.go

REM Move the binary to a directory in the PATH
move act.exe %SystemRoot%\System32\

REM Verify installation
where act.exe
if %errorlevel%==0 (
    echo act successfully installed!
) else (
    echo Installation failed.
)
