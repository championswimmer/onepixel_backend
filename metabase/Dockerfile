ARG METABASE_VERSION=v0.48.6
FROM alpine as builder

ARG METABASE_CLICKHOUSE_DRIVER_VERSION=1.3.2

RUN apk add --no-cache wget
RUN apk add --no-cache ca-certificates
RUN update-ca-certificates


# Move to working directory (/build).
WORKDIR /build

# Download metabase-clickhouse-driver from github
RUN wget "https://github.com/ClickHouse/metabase-clickhouse-driver/releases/download/$METABASE_CLICKHOUSE_DRIVER_VERSION/clickhouse.metabase-driver.jar"

ARG METABASE_VERSION
FROM metabase/metabase:$METABASE_VERSION

# Copy the clickhouse.jar into the plugins folder
COPY --from=builder /build/clickhouse.metabase-driver.jar /plugins/clickhouse.metabase-driver.jar

# Set the permissions for the clickhouse.jar
RUN chmod 744 /plugins/clickhouse.metabase-driver.jar