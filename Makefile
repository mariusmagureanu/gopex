all:
	@cd src && $(MAKE)

clean:
	@cd src && make clean
	@cd pkg/dbl && make clean

test:
	@cd pkg/dbl && make test

doc:
	@godoc -http :6060 -goroot .
