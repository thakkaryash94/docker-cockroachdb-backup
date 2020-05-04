# CockroachDB Auto-backup with Docker

CockroachDB is one of the most popular cloud-native databases. CockroachDB is an ACID compliant, relational database that’s wire compatible with PostgreSQL. CockroachDB delivers full ACID transactions at scale even in a distributed environment and guarantees serializable isolation in a cloud-neutral distributed database. We can deploy it using Docker and Kubernetes without any issue. CockroachDB delivers on the key cloud-native primitives of horizontal scale, no single points of failure, survivability, automatable operations, and no platform-specific encumbrances.

### How to run CockroachDB?

We can run CockroachDB using a single executable binary or using Docker and Kubernetes. We will see how to run CockroachDB in Docker and how to take backup as well. You can follow the blog or [Docs]([Start a Cluster in Docker (Insecure) | CockroachDB Docs](https://www.cockroachlabs.com/docs/stable/start-a-local-cluster-in-docker-mac.html)) and follow the backup instructions.

#### Start CockroachDB using Docker

We can run a single or multi node cluster with Docker but for development, we will be using single node cluster only.

Below command will pull the `latest` CockroachDB image from Docker hub, bind 26257, 8080 ports, bind the volume with cockroach_data/roach1. Run the command and go to http://localhost:8080, you will be able to see the CockroachDB Dashboard.

```shell
docker run -d \
      --name=roach1 \
      --hostname=roach1 \
      -p 26257:26257 -p 8080:8080  \
      -v "${PWD}/cockroach_data/roach1:/cockroach/cockroach-data"  \
      cockroachdb/cockroach:latest start \
      --insecure
```

Now, our CockroachDB database server is up and running. Now, Let's create some data using workload init command or you can dump existing database as well. Run below commands to do that. Here is the [docs](https://www.cockroachlabs.com/docs/stable/learn-cockroachdb-sql.html) explaining in details what we are doing in below command.

```shell
docker exec -it roach1 bash

./cockroach workload init movr 'postgresql://root@localhost:26257?sslmode=disable'
./cockroach sql --insecure --host=localhost:26257
USE movr;
SHOW tables;
```

As you can understand,  it will create a database `movr` with some tables and records. You should be able to see below output on your terminal.

```shell
          table_name
+----------------------------+
  promo_codes
  rides
  user_promo_codes
  users
  vehicle_location_histories
  vehicles
(6 rows)

Time: 9.579561ms
```

This means, we have successfully generated demo data and it's time to take a backup of it. Run `exit` 2 times to get out of the container. 1st exit command to get out of cockroach sql command line and 2nd is to exit from the container.

### CockroachDB Database backup

Our database is ready, now it's time to take a backup of it. To do that, there are 2 ways, we can take a backup of it.

#### 1. BACKUP

CockroachDB's BACKUP statement allows you to create full or incremental backups of your cluster's schema and data that are consistent as of a given timestamp. Backups can be with or without revision history.

There are many advantages of this process. We can setup whether we want to take a full backup or Incremental backup, automate backup with [JOBS](https://www.cockroachlabs.com/docs/stable/backup.html#viewing-and-controlling-backups-jobs), upload it to Amazon, Azure, Google Cloud, NFS or any S3-compatible services, backup a single table or view and many more. You can read more on it [here](https://www.cockroachlabs.com/docs/stable/backup.html).

This happens inside cockroachDB container environment, so cockroachDB has full control over it.

The only disadvantage is this is available only for [enterprise](https://www.cockroachlabs.com/product/cockroachdb/) users. This means that, if we are running cockroachDB locally or on a small server, where we may not want enterprise support, we can't use this feature. We can CockroachCloud, if we are running it on a small scale and planning to scale it in the future. CockroachCloud provides a Fully hosted and managed, Self-service platform with Enterprise features and basic support. You can read more [here](https://www.cockroachlabs.com/pricing/).

So how can we do it without BACKUP feature, which is available only for `CockroachCloud` and `CockroachDB Enterprise` users.

#### 2. cockroach dump

CockroachDB provides `dump` command, which is similar to pg_dump. The cockroach dump command outputs the SQL statements required to recreate tables, views, and sequences. This command can be used to back up or export each database in a cluster. The output should also be suitable for importing into other relational databases, with minimal adjustments. You can read more [here]([cockroach dump | CockroachDB Docs](https://www.cockroachlabs.com/docs/stable/cockroach-dump.html)).

That's great, so now, all we need to do it run a cron that executes `cockroach dump` whenever we want and that's it.

But, there are many things we have to think about. Like, how and where we are going to store our backup, how to take a backup on-demand etc etc.

Exactly, for that, I have created an open source docker image with the above features. So let's go through it one by one.

##### Features

- Customize cron with CRON_SCHEDULE env. https://godoc.org/github.com/robfig/cron
- Manual backup at any time
- Optional backup AWS S3/Spaces upload, if you provide ACCESS_KEY_ID, then it will take it as you want backup to uploaded on S3 or Spaces or anywhere compatible with s3 API.
- All cockroach image env variable support, you can override COCKROACH_USER, COCKROACH_INSECURE etc. [docs](https://www.cockroachlabs.com/docs/v19.2/cockroach-dump.html#client-connection)
- It exposes /data as volume, which contains backup zip file, so we can use backup from here if we don't want to upload it to any S3 services.

##### Environment Variables

###### Required:

- COCKROACH_DATABASE: Database name
- CRON_SCHEDULE: Cron value in double quotes. https://godoc.org/github.com/robfig/cron

###### Optional:

- ACCESS_KEY_ID: Spaces access key id
- BUCKET_NAME: Spaces bucket name
- S3_URL: AWS S3(s3.ap-south-1.amazonaws.com) or DO Spaces(nyc3.digitaloceanspaces.com)
- SECRET_ACCESS_KEY: Spaces secret access key

Now, let's look at how we can run our docker image to take a backup. Below, is the sample, how we can run the Docker container, which will take a backup of our movr database everyday and upload it to S3 service like AWS S3, Digital Ocean Spaces etc.

```shell
docker run -d \
      --name cockroach-backup \
      -v $(pwd)/data:/data \
      -v $(pwd)/cockroach-certs:/cockroach-certs \
      -e ACCESS_KEY_ID=ACCESS_KEY_ID \
      -e BUCKET_NAME=BUCKET_NAME \
      -e S3_URL=S3_URL \
      -e SECRET_ACCESS_KEY=SECRET_ACCESS_KEY \
      -e COCKROACH_DATABASE=movr \
      -e COCKROACH_HOST=localhost \
      -e CRON_SCHEDULE="0 0 * * *" \
      -e COCKROACH_INSECURE=true
      -e COCKROACH_USER=root
      thakkaryash94/cockroach-backup:latest
```

With the above command, we will run our backup container name `cockroach-backup` with `/data` volume, which will contain all the backup files and with ACCESS_KEY_ID, we can upload it to wherever we want. It supports every client connection params.

Now, let's say, you want to take a database backup of the current moment, you can run a below command to trigger manual backup as well.
```shell
curl -X POST http://localhost:9000/backup
```
These features are already available and few more already on the list. Feel free to open an issue to add more features.

Upcoming features:
- Flags support
- Multiple database backup support

**Note:** This is an open source project with MIT. This is my first golang project, I am a newbie, so I may have made mistakes. Issues and pull requests are most welcome.

Links:

- [Docker CockroachDB Backup GitHub](https://github.com/thakkaryash94/docker-cockroachdb-backup)
- [Learn CockroachDB SQL | CockroachDB Docs](https://www.cockroachlabs.com/docs/stable/learn-cockroachdb-sql.html)
