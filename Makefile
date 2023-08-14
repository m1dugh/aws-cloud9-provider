CC=go build
OUT=terraform-provider-awscloud9

.PHONY: clean build

build: $(OUT)

$(OUT):
	$(CC) -o $(OUT) .

clean:
	rm -f $(OUT)

fmt:
	@for f in `find . -type f -name "*.go"`; do go fmt "$$f"; done
