package it.unibo.states_source

import org.apache.flink.api.common.functions.AggregateFunction
import java.util.Date
import java.text.SimpleDateFormat
import java.time.Instant



class FinalAggregator extends AggregateFunction[(Int,Int,Long),Array[Int],(Int,Int,String)] {

  var latestTimestamp : Date = null

  override def createAccumulator() = Array[Int](0,0)

  override def add(input : (Int,Int,Long), acc : Array[Int]) =  {
    acc(0)=acc(0)+input._1
    acc(1)=acc(1)+input._2
    acc
  }

  override def getResult(acc:Array[Int])= {
    var myDate  = Date.from(Instant.now())
    myDate = roundDate(myDate)
    val formatter = new SimpleDateFormat("yyyy-MM-dd-HH-mm")
    val timestamp = formatter.format(myDate)
  
    (acc(0),acc(1),timestamp)
  }

  override def merge(a: Array[Int], b :Array[Int])= Array[Int](a(0)+b(0),a(1)+b(1))

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
