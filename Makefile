all:
	@cd src && $(MAKE)

clean:
	@cd src && $(MAKE)

doc:
	@godoc -http :6060 -goroot .
