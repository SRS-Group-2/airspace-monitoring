
package it.unibo.states_source

import org.apache.flink.api.common.functions.AggregateFunction




class MinimalStateAggregator extends AggregateFunction[MinimalState,Array[Double],(Double,Double,Long)] {

  var status : MinimalState = new MinimalState("",0,0,"")

  override def createAccumulator() = Array[Double](0,0)

  override def add(min: MinimalState, acc : Array[Double]) =  {
    val distance = min.getDistanceFromLatLonInKm(status)
    this.status=min
    val co2 : Double = (1062/distance)*90
    acc(0)=acc(0)+distance
    acc(1)=acc(1)+co2
    acc
  }

  override def getResult(acc:Array[Double])= (acc(0),acc(1),(System.currentTimeMillis() / 1000L))

  override def merge(a: Array[Double], b :Array[Double])= Array[Double](a(0)+b(0),a(1)+b(1))
}
