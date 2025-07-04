# Stage 1: Build Go binary
FROM golang:1.24-alpine AS go-builder

# Required for shared volume 
WORKDIR /usr/share/banditsecret

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/ ./cmd/
COPY internal/ ./internal/

RUN go build -v -o /usr/local/bin/banditsecret ./cmd/server/main.go


# Stage 2: Final runtime image with Python + Go binary
FROM python:3.12-alpine AS runtime

# Install pip
RUN apk add --no-cache py3-pip

WORKDIR /usr/share/banditsecret

# Create and activate virtualenv
ENV VIRTUAL_ENV=/usr/share/banditsecret/venv
ENV PATH="$VIRTUAL_ENV/bin:$PATH"

# Copy only requirements first (to leverage cache)
COPY scripts/requirements.txt ./scripts/requirements.txt

RUN python3 -m venv $VIRTUAL_ENV && \
    pip install --no-cache-dir -r ./scripts/requirements.txt

# Copy the python scripts
# This way, if only the python code changes (not requirements.txt), then docker will cache the pip install layer
COPY scripts/ ./scripts/

# Copy Go binary
COPY --from=go-builder /usr/local/bin/banditsecret /usr/local/bin/banditsecret

CMD ["banditsecret"]
