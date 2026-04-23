@echo off
setlocal
cd /d "%~dp0"

echo Cleaning build artifacts...

if exist futuapi4go-demo.exe del /q futuapi4go-demo.exe
if exist examples\futuapi4go-demo.exe del /q examples\futuapi4go-demo.exe

echo Done.
