# Bandit Secret
A tool that allows you to quickly find your favourite moments in youtube videos

## Features
- Search through captions of YouTube videos
- Automatically downloads English subtitles via yt-dlp
- Converts .vtt captions to structured JSON for easy processing
- Built with Go for speed

## Requirements
- Go 1.18+
- yt-dlp installed and the folder containing the executable is added to PATH (on Windows)
- Python 3 (converting .vtt caption file to JSON)

## Setting up Python environment and installing dependencies (Windows)
```
python -m venv venv
.\venv\Scripts\activate

pip install -r requirements.txt
```

## Running
Creating the database
```
mysql -u root -p < internal/storage/schema.sql
```


Build the binary 
```
go build -o bin/server.exe ./cmd/server
```

```
.\bin\server.exe
```

## Running with Docker
Completely remove network, volume mount, and container
```
docker-compose down -v
```


```
docker-compose up --build
```



## License
MIT - use freely, give credit where it's due