build:
	go build

dev: build
	./reddit-spotify-playlists --reddit-user=accountfortests --reddit-password=passwordfortests --subreddit=spotify

elastic:
	docker-compose up -d