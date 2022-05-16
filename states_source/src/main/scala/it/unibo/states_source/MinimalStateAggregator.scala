
package it.unibo.states_source

import org.apache.flink.api.common.functions.AggregateFunction




class MinimalStateAggregator extends AggregateFunction[MinimalState,Array[Int],(Int,Int,Long)] {

  var status : MinimalState = new MinimalState("",0,0,"")

  override def createAccumulator() = {
    status= new MinimalState("",0,0,"")
    Array[Int](0,0)
  }

  override def add(min: MinimalState, acc : Array[Int]) =  {
    val distance = min.getDistanceFromLatLonInKm(status)
    this.status=min
    var co2 : Int = 0
    if(distance>0){
      co2 = (distance*0.375).toInt
    }
    acc(0)=acc(0)+co2
    acc(1)=acc(1)+distance
    acc
  }

  override def getResult(acc:Array[Int])= (acc(0),acc(1),(System.currentTimeMillis() / 1000L))

  override def merge(a: Array[Int], b :Array[Int])= Array[Int](a(0)+b(0),a(1)+b(1))
}
