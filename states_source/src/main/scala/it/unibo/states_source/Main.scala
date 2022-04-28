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
    val sourceFunction = new OpenSkySourceFunction(36.619987291, 47.1153931748, 6.7499552751, 18.4802470232)

    val params = ParameterTool.fromArgs(args)
    val env = StreamExecutionEnvironment.getExecutionEnvironment
    env.getConfig.setGlobalJobParameters(params)
    // TODO: make bounding box coordinates more clear

    // val sourceFunction = new OpenSkySourceFunction(36.619987291, 47.1153931748, 6.7499552751, 18.4802470232)
    val source = env.addSource(sourceFunction)
/*
  //Case with file as sink
    val fileSink: FileSink[Vectors] = 
      FileSink.forRowFormat(new Path("/usr/local/flink/output"), 
                            new SimpleStringEncoder[Vectors]("UTF-8"))
              .withRollingPolicy(DefaultRollingPolicy.builder()
                                                     .withRolloverInterval(TimeUnit.MINUTES.toMillis(1))
                                                     .withInactivityInterval(TimeUnit.MINUTES.toMillis(1))
                                                     .withMaxPartSize(1024 * 1024 * 1024)
                                                     .build())
	            .build()  

      source.map(lsv => new Vectors(
                                      lsv.map(sv => new MinimalState(sv.getIcao24(), 
                                                                     sv.getLatitude(), 
                                                                     sv.getLongitude(), 
                                                                    /* sv.isOnGround(), 
                                                                     LocalDateTime.now().toString()*/)),
                                      LocalDateTime.now().toString()))
          .sinkTo(fileSink)
*/

// Case with pubsub as sink
    val serializationSchema: SerializationSchema[Vectors]  = new CustomJSONSerializer();
    val pubsubSink: SinkFunction[Vectors] = PubSubSink.newBuilder()
                                              .withSerializationSchema(serializationSchema)
                                              //Use this with the real pubsub, comment with the emulator
                                            //.withCredentials("pathToCredentials")  
                                              .withProjectName("projectname")
                                              .withTopicName("topicname")
                                              //Use this with the emulator, comment with the real pubsub
                                              .withHostAndPortForEmulator("host.docker.internal:8085")
                                              .build()
    source.map(lsv => new Vectors(
                                      lsv.map(sv => new MinimalState(sv.getIcao24(), 
                                                                     sv.getLatitude(), 
                                                                     sv.getLongitude(), 
                                                                    /* sv.isOnGround(), 
                                                                     LocalDateTime.now().toString()*/)),
                                      LocalDateTime.now().toString()))
          .addSink(pubsubSink)

    env.execute("OpenSkyStreamApp") // needed to avoid the No Job Found error

  }

}
