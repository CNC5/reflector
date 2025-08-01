FROM python:3.13-alpine

WORKDIR /app
COPY requirements.txt requirements.txt
RUN apk add --no-cache gcc musl-dev linux-headers
RUN pip install -r requirements.txt
RUN apk add --no-cache bash curl

# Runtime dependencies
COPY download_binaries.sh download_binaries.sh
RUN bash download_binaries.sh
COPY download_camo_templates.sh download_camo_templates.sh
RUN bash download_camo_templates.sh
RUN apk add --no-cache openssl certbot
COPY serverops serverops
COPY test.config.yaml config.yaml
COPY LICENSE LICENSE
RUN adduser -D -G www-data www-data

ENTRYPOINT ["python", "-m", "serverops"]
