build:
	go build

dev: build
	./redspot-finder-scraper \
		--reddit-user=accountfortests --reddit-password=passwordfortests \
		--subreddit=spotify

elastic:
	docker-compose up -d