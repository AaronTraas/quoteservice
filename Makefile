.DEFAULT_GOAL := build

mkfile_path := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
app_src_path := "$(mkfile_path)quoteservice"
app_dest_path := "/usr/local/bin/quoteservice"

service_src_path := "$(mkfile_path)quoteservice.service"
service_dest_path := "/etc/systemd/system/quoteservice.service"

coverage_file := coverage.out

clean:
	rm -rf quoteservice $(coverage_file)

build:
	go build

test:
	go test -coverprofile=$(coverage_file)

coverage: test
	go tool cover -func=$(coverage_file)

cover: coverage

coverage-html: test
	go tool cover -html=$(coverage_file)

cover-html: coverage-html

run: build
	./quoteservice 9009

install: build
	cp $(app_src_path) $(app_dest_path)
	echo $(mkfile_path)
	sudo ln -sf $(service_src_path) $(service_dest_path)
	sudo systemctl enable quoteservice
	sudo systemctl daemon-reload
	sudo systemctl restart quoteservice
