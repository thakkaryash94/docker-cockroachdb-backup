# Docker Cockroach DB Backup
Docker image for Cockroach DB backup

### Environment Variables

- ACCESS_KEY_ID: Spaces access key id
- BUCKET_NAME: Spaces bucket name
- CRON_SCHEDULE: Cron value in double quotes. https://godoc.org/github.com/robfig/cron
- S3_URL: Spaces url(nyc3.digitaloceanspaces.com)
- SECRET_ACCESS_KEY: Spaces secret access key


### Volumes:

- mount backup folder with `/data` path

- Run cockroach db using below command

```
docker run -d \
      --name=roach1 \
      --hostname=roach1 \
      -p 26257:26257 -p 8080:8080  \
      -v "${PWD}/cockroach_data/roach1:/cockroach/cockroach-data"  \
      cockroachdb/cockroach:v19.2.4 start \
      --insecure
```

#### Example:

```sh
docker run -d \
      --name cockroach-backup \
      -v $(pwd)/data:/data \
      -e ACCESS_KEY_ID=ACCESS_KEY_ID \
      -e BUCKET_NAME=BUCKET_NAME \
      -e S3_URL=S3_URL \
      -e SECRET_ACCESS_KEY=SECRET_ACCESS_KEY \
      -e COCKROACH_DATABASE=COCKROACH_DATABASE \
      -e COCKROACH_HOST=COCKROACH_HOST \
      -e CRON_SCHEDULE="0 0 * * *" \
      thakkaryash94/cockroach-backup:latest
```
