# RSS Aggregator

A CLI tool which you can use to collect and view RSS feeds from multiple sources.

Built as part of [Build a Blog Aggregator in Go - Boot.dev](https://www.boot.dev/courses/build-blog-aggregator-golang) course.

## Learning goals

1. Learn how to integrate a Go application with a PostgreSQL database.
2. Practice using your SQL skills to query and migrate a database (using sqlc and goose, two lightweight tools for typesafe SQL in Go).
3. Learn how to write a long-running service that continuously fetches new posts from RSS feeds and stores them in the database.

## Installation

You need to have Go and Postgres installed to use this CLI tool.

Install the `go-rss-gator` by running

```
go install https://github.com/Zigelzi/go-rss-gator
```

### Configuration

You need to create `.gatorconfig.json` file to your home directory (e.g `$HOME` in Unix systems) . It is used for storing the configuration of the program.

Add `db_url` key to the config file with the Postgres connection string (e.g `postgres://username:password@localhost:5432/DB_NAME?sslmode=disable`)

## Usage

You have the following commands to use.
````
Usage: go-rss-gator COMMAND [ARGUMENTS]

register [USERNAME] - Register new user
login [USERNAME] - Log in as existing user
users - List all existing users

addfeed [NAME URL] - Add and follow new feed
agg [TIME_BETWEEN_REQUESTS] - Aggregate all the content from the existing feeds with given interval
feeds - List all available feeds
follow [URL] - Follow an existing feed
unfollow [URL] - Unfollow a feed
following - List all feeds that you're following

browse [LIMIT] - Browse posts from the feeds you're following
```