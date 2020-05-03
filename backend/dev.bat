@ECHO OFF
IF NOT EXIST "config.json" COPY "default_config.json" "config.json"
SET MODE=development
go run main.go