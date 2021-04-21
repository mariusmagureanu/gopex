all:
	@cd src && $(MAKE)

clean:
	@cd src && make clean

doc:
	@godoc -http :6060 -goroot .
