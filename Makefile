all: test clean watch

watch:
	air .

test:
	@echo "tested"

clean:
	@echo "cleaned"

.PHONY: all watch test clean
