package it.unibo.states_source

import org.apache.flink.api.common.functions.AggregateFunction
import java.util.Date
import java.text.SimpleDateFormat
import java.time.Instant



class FinalAggregator extends AggregateFunction[(Double,Double,Long),Array[Double],(Double,Double,String)] {
  

  override def createAccumulator() = Array[Double](0,0)

  override def add(input : (Double,Double,Long), acc : Array[Double]) =  {
    acc(0)=acc(0)+input._1
    acc(1)=acc(1)+input._2
    acc
  }

  override def getResult(acc:Array[Double])= {
  val myDate  = Date.from(Instant.now())
  val formatter = new SimpleDateFormat("yyyy-MM-dd-HH-mm");
  val timestamp = formatter.format(myDate);
  
   (acc(0),acc(1),timestamp)
  }

  override def merge(a: Array[Double], b :Array[Double])= Array[Double](a(0)+b(0),a(1)+b(1))
}
