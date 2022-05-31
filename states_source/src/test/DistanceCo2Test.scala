package it.unibo.states_source

import org.opensky.model.StateVector
import org.scalatest.FlatSpec
import org.scalatest.Matchers
import org.scalatest.BeforeAndAfter
import org.apache.flink.test.util.MiniClusterWithClientResource
import org.apache.flink.test.util.MiniClusterResourceConfiguration
import org.apache.flink.streaming.api.environment.StreamExecutionEnvironment
import java.util.Date
import java.text.SimpleDateFormat
import java.time.Instant
import org.apache.flink.streaming.api.functions.sink.SinkFunction
import org.apache.flink.streaming.api.windowing.assigners.TumblingProcessingTimeWindows
import org.apache.flink.streaming.api.windowing.time.Time

class DistanceCo2Test extends FlatSpec with Matchers {

  val flinkCluster = new MiniClusterWithClientResource(new MiniClusterResourceConfiguration.Builder()
      .setNumberSlotsPerTaskManager(2)
      .setNumberTaskManagers(1)
      .build)

  "DistanceCo2" should "aggregate distance and calculate co2" in {

  val env = StreamExecutionEnvironment.getExecutionEnvironment

  // configure your test environment
  env.setParallelism(1)

  // values are collected in  a static variable
  CollectDistanceCo2Sink.km=0
  CollectDistanceCo2Sink.co2=(-1)
  val testSource = new TestDistanceCo2Source()
  val source = env.addSource(testSource)
  // create a stream of custom elements and apply transformations
  source.map(
      sv => new MinimalState(
        sv.getIcao24(), 
        sv.getLatitude(), 
        sv.getLongitude(),
        (System.currentTimeMillis() / 1000L).toString()
      )
    )
      .keyBy(ms => ms.icao)
      .window(TumblingProcessingTimeWindows.of(Time.seconds(1)))
      .aggregate(new MinimalStateAggregator())
      .windowAll(TumblingProcessingTimeWindows.of(Time.seconds(2)))
      .aggregate(new FinalAggregator())
      .addSink(new CollectDistanceCo2Sink())

  // execute
  env.execute()

  CollectDistanceCo2Sink.km should equal (302)
  CollectDistanceCo2Sink.co2 should equal (0)
  }
}

// create a testing sink
class CollectDistanceCo2Sink extends SinkFunction[(Int,Int,String)] {

  override def invoke(elems: (Int,Int,String), context: SinkFunction.Context): Unit = {
    CollectDistanceCo2Sink.km=elems._2
    CollectDistanceCo2Sink.co2=elems._1
  }
}

object CollectDistanceCo2Sink {
// must be static
  var km: Int = 0
  var co2: Int = (-1)
}


