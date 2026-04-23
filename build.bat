@echo off
setlocal
cd /d "%~dp0"

echo Building futuapi4go-demo...

go build ./...

if %ERRORLEVEL% EQU 0 (
    echo Build successful.
) else (
    echo Build failed!
    exit /b %ERRORLEVEL%
)
