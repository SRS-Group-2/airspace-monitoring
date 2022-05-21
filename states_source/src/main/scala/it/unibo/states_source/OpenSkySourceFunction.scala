package it.unibo.states_source

import collection.JavaConverters._

import org.apache.flink.streaming.api.functions.source.SourceFunction
import org.apache.flink.streaming.api.functions.source.SourceFunction.SourceContext

import org.opensky.api.OpenSkyApi
import org.opensky.model.StateVector

import org.slf4j.Logger;
import org.slf4j.LoggerFactory;

class OpenSkySourceFunction(private val coordinates: Option[(Double, Double, Double, Double)]) extends SourceFunction[StateVector] {
  val PERIOD_S = 15L
  val PERIOD_MS = PERIOD_S * 1000

  val LOG : Logger = LoggerFactory.getLogger(classOf[OpenSkySourceFunction])

  var running: Boolean = true

  def this() {
    this(None)
  }

  def this(coor: String) {
    this(Some(
      ((coor.split(",")(0).toDouble),
      (coor.split(",")(1).toDouble),
      (coor.split(",")(2).toDouble),
      (coor.split(",")(3).toDouble))
      )
    )
  }

  def this(a: Double, b: Double, c: Double, d: Double) {
    this(Some((a, b, c, d)))
  }
 
  override def run(ctx: SourceContext[StateVector]): Unit = {
    val api = new OpenSkyApi()
    while (running) {
      val os = coordinates match {
        case Some((a, b, c, d)) => api.getStates(0, null, new OpenSkyApi.BoundingBox(a, b, c, d))
        case None => api.getStates(0, null)
      }
      os.getStates().asScala.foreach {
        ctx.collect(_)
      }
      LOG.info("Read data from OpenSky and sent")
      Thread.sleep(PERIOD_MS)
    }
  }

  override def cancel(): Unit = {
      running = false
  }
}
