# Docker Cockroach DB Backup ![Build and Push Docker](https://github.com/thakkaryash94/docker-cockroachdb-backup/workflows/Build%20and%20Push%20Docker/badge.svg)

## Features

- Customize cron with CRON_SCHEDULE env. https://godoc.org/github.com/robfig/cron
- Manual backup at any time, run `curl -X POST http://localhost:9000/backup` to take current data backup
- Optional backup AWS S3/Spaces upload, if you provide ACCESS_KEY_ID, then it will take it as you want backup to uploaded on S3 or Spaces or anywhere compatible with s3 API.
- All cockroach image env variable support, you can override COCKROACH_USER, COCKROACH_INSECURE etc. [docs](https://www.cockroachlabs.com/docs/v19.2/cockroach-dump.html#client-connection)
- It exposes /data as volume, which contains backup zip file, so can use backup from here if you are not uploading it to any S3 services.

### Environment Variables

#### Required

- COCKROACH_DATABASE: Database name
- CRON_SCHEDULE: Cron value in double quotes. https://godoc.org/github.com/robfig/cron

#### Optional

- ACCESS_KEY_ID: Spaces access key id
- BUCKET_NAME: Spaces bucket name
- S3_URL: AWS S3(s3.ap-south-1.amazonaws.com) or DO Spaces(nyc3.digitaloceanspaces.com)
- SECRET_ACCESS_KEY: Spaces secret access key

### Volumes

- mount backup folder with `/data` path
- Run cockroach db using below command

```sh
docker run -d \
      --name=roach1 \
      --hostname=roach1 \
      -p 26257:26257 -p 8080:8080  \
      -v "${PWD}/cockroach_data/roach1:/cockroach/cockroach-data"  \
      cockroachdb/cockroach:v19.2.4 start \
      --insecure
```

#### Example

```sh
docker run -d \
      --name cockroach-backup \
      -v $(pwd)/data:/data \
      -v $(pwd)/cockroach-certs:/cockroach-certs \
      -e ACCESS_KEY_ID=ACCESS_KEY_ID \
      -e BUCKET_NAME=BUCKET_NAME \
      -e S3_URL=S3_URL \
      -e SECRET_ACCESS_KEY=SECRET_ACCESS_KEY \
      -e COCKROACH_DATABASE=COCKROACH_DATABASE \
      -e COCKROACH_HOST=COCKROACH_HOST \
      -e CRON_SCHEDULE="0 0 * * *" \
      -e COCKROACH_INSECURE=true
      -e COCKROACH_USER=root
      thakkaryash94/cockroach-backup:latest
```
