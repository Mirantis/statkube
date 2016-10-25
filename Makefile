.PHONY: all .FORCE

all: migrate getemployment getprs

migrate: .FORCE
	go build -o migrate cmd/migrate.go

getemployment: .FORCE
	go build -o getemployment cmd/getemployment.go

getprs: .FORCE
	go build -o getprs cmd/getprs.go

.FORCE:
