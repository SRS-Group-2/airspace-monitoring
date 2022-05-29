package it.unibo.states_source


import org.apache.flink.streaming.api.functions.source.SourceFunction
import org.apache.flink.streaming.api.functions.source.SourceFunction.SourceContext
import org.opensky.model.StateVector

class TestListSource extends SourceFunction[StateVector]{


    override def run(ctx: SourceContext[StateVector]): Unit = {

      val sv1= new StateVector("abcdef")
      val sv2= new StateVector("abc123")
      val sv3= new StateVector("abc456")
      val sv4= new StateVector("abc789")
      val sv5= new StateVector("abcdef")
      
      List(sv1,sv2,sv3,sv4,sv5).foreach {
          ctx.collect(_)
      }
      Thread.sleep(2000)
    }

    override def cancel() : Unit = {
    }  
}
