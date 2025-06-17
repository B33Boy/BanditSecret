FROM golang:1.24-alpine

WORKDIR /usr/share/banditsecret

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN apk add --no-cache python3 py3-pip python3-dev \
    && rm -rf /var/cache/apk/*


ENV VIRTUAL_ENV=/usr/share/banditsecret/venv/
RUN python3 -m venv $VIRTUAL_ENV

COPY requirements.txt .
# set -e ensures that if any command fails, the build stops
RUN set -e && \
    . $VIRTUAL_ENV/bin/activate && \
    pip install --no-cache-dir -r requirements.txt

RUN go build -v -o /usr/local/bin/banditsecret ./cmd/server/main.go

CMD ["banditsecret"]
