run: build
	@./dist/tergom


build:
	@go build -o dist/tergom .
	@go build -o dist/tergom.exe .