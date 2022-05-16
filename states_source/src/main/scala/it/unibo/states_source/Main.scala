package it.unibo.states_source


import collection.JavaConverters._
import org.apache.flink.api.java.utils.ParameterTool
import org.apache.flink.api.scala._
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment
import org.apache.flink.api.common.serialization.SerializationSchema
import org.apache.flink.streaming.api.functions.sink.SinkFunction
import org.apache.flink.streaming.api.windowing.assigners.TumblingProcessingTimeWindows
import org.apache.flink.streaming.api.windowing.time.Time




import org.opensky.api.OpenSkyApi

object Main {
  def main(args: Array[String]): Unit = {
    val sourceFunction = new OpenSkySourceFunction(System.getenv("COORDINATES"))
    val params = ParameterTool.fromArgs(args)
    val env = StreamExecutionEnvironment.getExecutionEnvironment
    env.getConfig.setGlobalJobParameters(params)
    val source = env.addSource(sourceFunction)
    val minsSerializationSchema: SerializationSchema[MinimalState]  = new JSONMinStateSerializer()
    val aircraftsSerializationSchema: SerializationSchema[Aircrafts]  = new JSONAircraftsSerializer()
    val minsPubsubSink: SinkFunction[MinimalState] = 
      CustomPubSubSink
        .newBuilder()
        .withSerializationSchema(minsSerializationSchema)  
        .withProjectName(System.getenv("GOOGLE_CLOUD_PROJECT_ID"))
        .withTopicName(System.getenv("GOOGLE_PUBSUB_VECTORS_TOPIC_ID"))
        .build()
    
    source.map(
      sv => new MinimalState(
        sv.getIcao24(), 
        sv.getLatitude(), 
        sv.getLongitude(),
        (System.currentTimeMillis() / 1000L).toString()
      )
    )
      .addSink(minsPubsubSink)

    val aFirebaseSink : SinkFunction[Aircrafts] = new AircraftsFirebaseSink()
    
    source.map(ms => ms.getIcao24())
      .windowAll(TumblingProcessingTimeWindows.of(Time.minutes(2)))
      .aggregate(new AircraftsAggregator())
      .addSink(aFirebaseSink)
    
    val fMinutesFirebaseSink : SinkFunction[(Int,Int,String)] = new FiveMinutesFirebaseSink()
    
    source.map(
      sv => new MinimalState(
        sv.getIcao24(), 
        sv.getLatitude(), 
        sv.getLongitude(),
        (System.currentTimeMillis() / 1000L).toString()
      )
    )
      .keyBy(ms => ms.icao)
      .window(TumblingProcessingTimeWindows.of(Time.minutes(5)))
      .aggregate(new MinimalStateAggregator())
      .windowAll(TumblingProcessingTimeWindows.of(Time.minutes(5)))
      .aggregate(new FinalAggregator())
      .addSink(fMinutesFirebaseSink)
  
    env.execute("OpenSkyStreamApp") 

  }

}
