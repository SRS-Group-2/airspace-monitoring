package it.unibo.states_source

import collection.JavaConverters._

import org.apache.flink.streaming.api.functions.source.SourceFunction
import org.apache.flink.streaming.api.functions.source.SourceFunction.SourceContext

import org.opensky.api.OpenSkyApi
import org.opensky.model.StateVector

class OpenSkySourceFunction(private val coordinates: Option[(Double, Double, Double, Double)]) extends SourceFunction[List[StateVector]] {
  val PERIOD_S = 10L
  val PERIOD_MS = PERIOD_S * 1000

  var running: Boolean = true

  def this() {
    this(None)
  }

  def this(a: Double, b: Double, c: Double, d: Double) {
    this(Some((a, b, c, d)))
  }
 
  override def run(ctx: SourceContext[List[StateVector]]): Unit = {
    val api = new OpenSkyApi()
    while (running) {
      val os = coordinates match {
        case Some((a, b, c, d)) => api.getStates(0, null, new OpenSkyApi.BoundingBox(a, b, c, d))
        case None => api.getStates(0, null)
      }
      ctx.collect(os.getStates().asScala.toList)
      Thread.sleep(PERIOD_MS)
    }
  }

  override def cancel(): Unit = {
      running = false
  }
}