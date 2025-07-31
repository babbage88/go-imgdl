BINARY_NAME = imgdl
ARTIFACT_DIR = dist/
INSTALL_DIR = $(HOME)/go/bin

.PHONY: all build install clean

all: build

build:
	@mkdir -p $(ARTIFACT_DIR)
	@echo Removing old build artifacts...
	@rm -f $(ARTIFACT_DIR)/$(BINARY_NAME)
	@echo "Building $(BINARY_NAME)..."
	go build -v -o $(ARTIFACT_DIR)/$(BINARY_NAME) .

install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	mkdir -p $(INSTALL_DIR)
	cp $(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Installed $(BINARY_NAME) to $(INSTALL_DIR)/$(BINARY_NAME)"

clean:
	@echo "Cleaning..." 
	rm -f $(ARTIFACT_DIR)/$(BINARY_NAME)
