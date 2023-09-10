compile:
	go build -o build/gcat

benchmark:
	go test -bench=.

debug:
	make compile
	./build/gcat -n test/1.txt