FROM python:alpine3.22

RUN apk add --no-cache curl

# Install yt-dlp binary
RUN curl -L https://github.com/yt-dlp/yt-dlp/releases/latest/download/yt-dlp \
    -o /usr/local/bin/yt-dlp && \
    chmod +x /usr/local/bin/yt-dlp

WORKDIR /usr/share/banditsecret

COPY requirements.txt .
RUN pip install --no-cache-dir -r requirements.txt

COPY ytdlp.py .

# default port, overriden with docker-compose env var
ENV YTDLP_PORT=5000
EXPOSE ${YTDLP_PORT}

# Run app using gunicorn
CMD gunicorn -b ":${YTDLP_PORT}" ytdlp:app