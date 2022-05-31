package it.unibo.states_source

import org.apache.flink.api.common.functions.AggregateFunction
import java.util.Date
import java.text.SimpleDateFormat
import java.time.Instant



class FinalAggregator extends AggregateFunction[(Double,Double,Long),Array[Double],(Int,Int,String)] {

  override def createAccumulator() = Array[Double](0,0)

  var latestTimestamp : Date = null

  override def add(input : (Double,Double,Long), acc : Array[Double]) = {
    acc(0)=acc(0)+input._1
    acc(1)=acc(1)+input._2
    acc
  }

  override def getResult(acc:Array[Double])= {
    var myDate  = Date.from(Instant.now())
    myDate = roundDate(myDate)
    latestTimestamp=myDate
    val formatter = new SimpleDateFormat("yyyy-MM-dd-HH-mm")
    val timestamp = formatter.format(myDate)
  
    (acc(0).toInt,acc(1).toInt,timestamp)
  }

  override def merge(a: Array[Double], b :Array[Double])= Array[Double](a(0)+b(0),a(1)+b(1))

  def roundDate(myDate : Date) : Date = {
  if(latestTimestamp!=null) {
      if((myDate.getMinutes() - latestTimestamp.getMinutes())>5) {
        myDate.setMinutes(myDate.getMinutes() - ( myDate.getMinutes() % 5))
      }
      else if ((myDate.getMinutes() - latestTimestamp.getMinutes())<5) {
        myDate.setMinutes(myDate.getMinutes() + ( 5 - (myDate.getMinutes() % 5)))
      }
    }
  else {
    if((myDate.getMinutes()%5)>2) {
        myDate.setMinutes(myDate.getMinutes() + (5 - (myDate.getMinutes() % 5)))
      }
      else {
        myDate.setMinutes(myDate.getMinutes() - (myDate.getMinutes() % 5))
      }
  }

    return myDate
  }
}
