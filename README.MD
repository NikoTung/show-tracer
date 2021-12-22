## show tracer

Trace and download your TV shows.It traces your TV shows by a rss feed,once the feed updated,it will submit the new 
show episode to the download client,currently aria2 supported only,and send you a message by telegram.

### Configuration

It reads configuration from a file in the following format.

    {
        "rss": "rss feed",
        "dt": "your download client type(eg.aria2)",
        "api": "your download client api endpoint",
        "secret": "aria2 rpc secret",
        "telegram_token": "telegram bot token",
        "telegram_chat_id": "telegram chat id"
    }

For the rss feed,you can use [rarbg.to](https://rarbg.to/rssdd.php),or [showrss.info](https://showrss.info/).

Just build from the code `go build` then `./show-tracer -f config.json`


### docker 

    docker run -d  -v $PWD/config.json:/app/config.json --name show-tracer show-tracer:latest