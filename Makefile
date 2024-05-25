.DEFAULT_GOAL := build

fmt: 
	go fmt ./...

lint: fmt
	golint ./...

vet: fmt
	go vet ./...

build: vet
	go build -o bin/main main.go

# run: build
# 	go run main.go


# watch:
# 	@if [ -x "$(GOPATH)/bin/air" ]; then \
# 	    "$(GOPATH)/bin/air"; \
# 		@echo "Watching...";\
# 	else \
# 	    read -p "air is not installed. Do you want to install it now? (y/n) " choice; \
# 	    if [ "$$choice" = "y" ]; then \
# 			go install github.com/cosmtrek/air@latest; \
# 	        "$(GOPATH)/bin/air"; \
# 				@echo "Watching...";\
# 	    else \
# 	        echo "You chose not to install air. Exiting..."; \
# 	        exit 1; \
# 	    fi; \
# 	fi


clean:
	go clean
	rm bin/main

.PHONY: fmt lint vet build run clean