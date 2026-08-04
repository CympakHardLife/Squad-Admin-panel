@echo off
REM Сборка SquadAdmin.exe. Требуется Go 1.24 или новее (https://go.dev/dl/).
set GOOS=windows
set GOARCH=amd64
set CGO_ENABLED=0
go build -trimpath -ldflags "-s -w" -o SquadAdmin.exe .
echo Готово: SquadAdmin.exe
pause
