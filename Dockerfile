FROM python:3.13-alpine

WORKDIR /app
COPY requirements.txt requirements.txt
RUN apk add --no-cache gcc musl-dev linux-headers
RUN pip install -r requirements.txt
RUN apk add --no-cache bash curl
COPY download_binaries.sh download_binaries.sh
RUN bash download_binaries.sh
RUN apk add --no-cache openssl
COPY serverops serverops
COPY test.config.yaml config.yaml
RUN adduser -D -G www-data www-data

ENTRYPOINT ["python", "-m", "serverops"]
