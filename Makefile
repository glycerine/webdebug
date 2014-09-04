
normalbuild:
	go build
	go install

debugbuild:
	go build  -gcflags "-N -l"
	go install

testbuild:
	go test -c -gcflags "-N -l" -v

clean:
	rm -f webdebug *~ *.o

test:
	go test -v

startdebug:
	kill -USR1 $(shell cat webdebug.pid)

stopdebug:
	kill -USR2 $(shell cat webdebug.pid)
