@echo off
echo Building futuapi4go-demo...

set OUTPUT=cmd\demo\futuapi4go-demo.exe
go build -o %OUTPUT% ./cmd/demo

if %ERRORLEVEL% EQU 0 (
    echo Build successful: %OUTPUT%
) else (
    echo Build failed!
    exit /b %ERRORLEVEL%
)
