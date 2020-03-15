.DEFAULT_GOAL := build

mkfile_path := $(dir $(abspath $(lastword $(MAKEFILE_LIST))))
app_src_path := "$(mkfile_path)quoteservice"
app_dest_path := "/usr/local/bin/quoteservice"

service_src_path := "$(mkfile_path)quoteservice.service"
service_dest_path := "/etc/systemd/system/quoteservice.service"

clean:
	rm -rf quoteservice

build:
	go build

run: build
	./quoteservice 9009

install: build
	cp $(app_src_path) $(app_dest_path)
	echo $(mkfile_path)
	sudo ln -s $(service_src_path) $(service_dest_path)
	sudo systemctl enable quoteservice
	sudo systemctl daemon-reload
	sudo systemctl restart my_restart
