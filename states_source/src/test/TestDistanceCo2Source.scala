package it.unibo.states_source


import org.apache.flink.streaming.api.functions.source.SourceFunction
import org.apache.flink.streaming.api.functions.source.SourceFunction.SourceContext
import org.opensky.model.StateVector

class TestDistanceCo2Source extends SourceFunction[StateVector]{


    override def run(ctx: SourceContext[StateVector]): Unit = {

      val sv1= new StateVector("abcdef")
      sv1.setLatitude(22)
      sv1.setLongitude(22)
      val sv2= new StateVector("abcdef")
      sv2.setLatitude(23)
      sv2.setLongitude(23)
      val sv3= new StateVector("abcde1")
      sv3.setLatitude(22)
      sv3.setLongitude(22)
      val sv4= new StateVector("abcde1")
      sv4.setLatitude(23)
      sv4.setLongitude(23)
      
      
      List(sv1,sv2,sv3,sv4).foreach {
          ctx.collect(_)
      }
      Thread.sleep(8000)
    }

    override def cancel() : Unit = {
    }  
}
