.PHONY: all clean test

all:
	@cd src && $(MAKE)

clean:
	@cd src && make clean
	@cd pkg/dbl && make clean

test:
	@cd pkg/dbl && make test

doc:
	@godoc -http :6060 -goroot .

docker-build:
	@docker-compose build

docker-up:
	@docker-compose up -d --scale kc=5

docker-clean:
	@docker-compose down