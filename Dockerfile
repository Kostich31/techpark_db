FROM golang:1.17 AS build

ADD . /app
WORKDIR /app
RUN go build ./main.go

FROM ubuntu:20.04

RUN apt-get -y update && apt-get install -y tzdata
ENV TZ=Russia/Moscow
RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

ENV PGVER 12
RUN apt-get -y update && apt-get install -y postgresql-$PGVER
USER postgres

RUN /etc/init.d/postgresql start &&\
    psql --command "CREATE USER kostya WITH SUPERUSER PASSWORD 'kostya';" &&\
    createdb -O kostya kostya &&\
    /etc/init.d/postgresql stop

EXPOSE 5432
VOLUME  ["/etc/postgresql", "/var/log/postgresql", "/var/lib/postgresql"]
USER root

WORKDIR /usr/src/app
COPY . .
COPY --from=build /app/main/ .

EXPOSE 5000
ENV PGPASSWORD kostya
CMD service postgresql start && psql -h localhost -d kostya -U kostya -p 5432 -a -q -f ./db/db.sql && ./main
