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

class ListFunctionTest extends FlatSpec with Matchers {

  val flinkCluster = new MiniClusterWithClientResource(new MiniClusterResourceConfiguration.Builder()
      .setNumberSlotsPerTaskManager(2)
      .setNumberTaskManagers(1)
      .build)

  "ListFunction" should "aggregateAList" in {

  val env = StreamExecutionEnvironment.getExecutionEnvironment

  // configure your test environment
  env.setParallelism(1)

  // values are collected in a static variable
  CollectListSink.value=null
  val testSource = new TestListSource()
  val source = env.addSource(testSource)
  // create a stream of custom elements and apply transformations
  source.map(sv => sv.getIcao24())
    .windowAll(TumblingProcessingTimeWindows.of(Time.seconds(2)))
    .aggregate(new AircraftsAggregator())
    .addSink(new CollectListSink())

  // execute
  env.execute()

  
  val myDate  = Date.from(Instant.now())
  val formatter = new SimpleDateFormat("yyyy-MM-dd-HH-mm")
  val timestamp = formatter.format(myDate)
  val air = new Aircrafts(timestamp,List("abcdef","abc123","abc456","abc789"))
  CollectListSink.value.list should equal (air.list)
  }
}

// create a testing sink
class CollectListSink extends SinkFunction[Aircrafts] {

  override def invoke(air: Aircrafts, context: SinkFunction.Context): Unit = {
    CollectListSink.value=air
  }
}

object CollectListSink {
// must be static
  var value: Aircrafts = null
}

