@echo off

del .\bin\goc.exe

setlocal

if exist install.bat goto ok
echo install.bat must be run from its folder
goto end

:ok

set GOBIN=
set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0

gofmt -w src

go install goc

:end
echo [FINISHED]
echo.

cls
.\bin\goc.exe -REDIS_PWD=shuoleniyebudong -CLIENT_ID=backend_1 -REDIS_SHA_AUTH=a0ad12f31d7de75a5153bdff954caf5bc15b9501