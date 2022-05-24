FROM maven:3.8.5-jdk-11 AS builder

COPY src /home/app/src
COPY pom.xml /home/app

RUN git clone https://github.com/openskynetwork/opensky-api.git \
 && mvn -f ./opensky-api/java/pom.xml install -DskipTests -Dmaven.javadoc.skip=true

RUN mvn -f /home/app/pom.xml clean package

FROM flink:1.15.0-scala_2.12
COPY --from=builder /home/app/target/states_source-0.1.jar /opt/flink/usrlib/job.jar
