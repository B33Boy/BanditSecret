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

Build the binary 
```
go build -o bin/search.exe ./cmd/search
```

```
.\bin\search.exe
```

## License
MIT - use freely, give credit where it's due