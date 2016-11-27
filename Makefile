IMPORT_PATH = $(shell echo `pwd` | sed "s|^$(GOPATH)/src/||g")
APP_NAME = $(shell echo $(IMPORT_PATH) | sed 's:.*/::')
APP_VERSION = 0.1
TARGET = ./$(APP_NAME)-$(APP_VERSION)
DIST_TARGET = ./$(APP_NAME)-$(APP_VERSION)-dist
DIST_EXTERNAL_TARGET = ./$(APP_NAME)-external-$(APP_VERSION)-dist
GO_FILES = $(shell find . -type f -name "*.go")
BUNDLE = public/bundles
ASSETS = $(shell find assets -type f)
PID = .pid
NODE_BIN = $(shell npm bin)
#go server port
PORT ?= 9000
#webpack-dev-server port
DEV_HOT_PORT ?= 8090

build: clean $(BUNDLE) $(TARGET) $(DIST_TARGET) $(DIST_EXTERNAL_TARGET)

clean:
	@rm -rf public/bundles
	@rm -rf $(TARGET)
	@rm -rf $(DIST_TARGET)
	@rm -rf $(DIST_EXTERNAL_TARGET)
	@rm -rf $(APP_NAME)-$(APP_VERSION).zip

$(BUNDLE): $(ASSETS)
	@$(NODE_BIN)/webpack --progress --colors

$(TARGET): $(GO_FILES)
	@printf "Building go binary ......"
	@go build -race -o $@

$(DIST_TARGET): $(GO_FILES)
	@printf "Building dist go binary ......"
	@env GOOS=linux GOARCH=amd64 go build -o $@

$(DIST_EXTERNAL_TARGET): $(GO_FILES)
	@printf "Building dist external go binary ......"
	@env GOOS=linux GOARCH=amd64 go build -o $@ ./external

kill:
	@kill `cat $(PID)` || true

dev: clean $(BUNDLE) restart
	@DEV_HOT=true NODE_ENV=development $(NODE_BIN)/webpack-dev-server --config webpack.config.js &
	@printf "\n\nWaiting for the file change\n\n"
	@fswatch --one-per-batch $(GO_FILES) | xargs -n1 -I{} make restart || make kill

restart: kill $(TARGET)
	@printf "\n\nrestart the app .........\n\n"
	@$(TARGET) -debug --web=:$(PORT) --devWeb=:$(DEV_HOT_PORT) & echo $$! > $(PID)

dist: clean $(TARGET) $(DIST_TARGET) $(DIST_EXTERNAL_TARGET)
	@NODE_ENV=production $(NODE_BIN)/webpack --progress --colors
	@zip -r -v $(APP_NAME)-$(APP_VERSION).zip $(DIST_TARGET) $(DIST_EXTERNAL_TARGET) \
    	webpack-assets.json public templates static deploy
