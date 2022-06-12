lg: lg.go pkg/**/*
	go build -ldflags="-s -w" -o lg .