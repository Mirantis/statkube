.PHONY: all .FORCE

all: migrate getemployment

migrate: .FORCE
	go build -o migrate cmd/migrate.go

getemployment: .FORCE
	go build -o getemployment cmd/getemployment.go

.FORCE:
