run:
	docker build -t go-casket-dev . \
	&& docker run --rm -it --name go-casket-dev go-casket-dev sh
