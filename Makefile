build: generate
	go build -o bin/chatdev .

generate:
	templ generate 

run: build
	air
