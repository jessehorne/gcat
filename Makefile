compile:
	go build -o build/gcat

benchmark:
	go test -bench=.

debug:
	make compile
	./build/gcat -s test/1.txt