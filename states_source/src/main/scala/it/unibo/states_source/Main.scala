package it.unibo.states_source

import java.time.LocalDateTime
import java.util.concurrent.TimeUnit

import collection.JavaConverters._

import org.apache.flink.api.java.utils.ParameterTool
import org.apache.flink.api.common.eventtime._
import org.apache.flink.api.common.serialization.SimpleStringSchema
import org.apache.flink.api.scala._
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment
import org.apache.flink.core.fs.Path
import org.apache.flink.connector.file.sink.FileSink
import org.apache.flink.streaming.api.functions.sink.filesystem.rollingpolicies.DefaultRollingPolicy
import org.apache.flink.api.common.serialization.SimpleStringEncoder
import org.apache.flink.api.common.serialization.SerializationSchema
import org.apache.flink.streaming.api.functions.sink.SinkFunction
import org.apache.flink.streaming.connectors.gcp.pubsub.PubSubSink


import org.opensky.api.OpenSkyApi

object Main {
  def main(args: Array[String]): Unit = {
    val sourceFunction = new OpenSkySourceFunction(System.getenv("COORDINATES"))

    val params = ParameterTool.fromArgs(args)
    val env = StreamExecutionEnvironment.getExecutionEnvironment
    env.getConfig.setGlobalJobParameters(params)
    val sourceVectors = env.addSource(sourceFunction)
    val sourceAircrafts = env.addSource(sourceFunction)
    val vectorsSerializationSchema: SerializationSchema[Vectors]  = new JSONVectorsSerializer();
    val aircraftsSerializationSchema: SerializationSchema[Aircrafts]  = new JSONAircraftsSerializer();
    val vectorsPubsubSink: SinkFunction[Vectors] = PubSubSink.newBuilder()
                                              .withSerializationSchema(vectorsSerializationSchema)
                                              .withCredentials(System.getenv("GOOGLE_APPLICATION_CREDENTIALS"))  
                                              .withProjectName(System.getenv("GOOGLE_CLOUD_PROJECT_ID"))
                                              .withTopicName(System.getenv("GOOGLE_PUBSUB_VECTORS_TOPIC_ID"))
                                              .build()
    sourceVectors.map(lsv => new Vectors(
                                      lsv.map(sv => new MinimalState(sv.getIcao24(), 
                                                                     sv.getLatitude(), 
                                                                     sv.getLongitude())),
                                      (System.currentTimeMillis() / 1000L).toString()))
                      .addSink(vectorsPubsubSink)

    val aircraftsPubsubSink: SinkFunction[Aircrafts] = PubSubSink.newBuilder()
                                                      .withSerializationSchema(aircraftsSerializationSchema)
                                                      .withCredentials("GOOGLE_APPLICATION_CREDENTIALS")  
                                                      .withProjectName(System.getenv("GOOGLE_CLOUD_PROJECT_ID"))
                                                      .withTopicName(System.getenv("GOOGLE_PUBSUB_AIRCRAFT_LIST_TOPIC_ID"))
                                                      .build()
    sourceAircrafts.map(lsv => new Aircrafts(
                                      (System.currentTimeMillis() / 1000L).toString(),
                                      lsv.map(sv => sv.getIcao24())))
                        .addSink(aircraftsPubsubSink)

    env.execute("OpenSkyStreamApp") 

  }

}
