
package it.unibo.states_source

import org.apache.flink.api.common.functions.AggregateFunction



class MinimalStateAggregator extends AggregateFunction[MinimalState,Accumulator,(Double,Double,Long)] {


  override def createAccumulator() = {
    new Accumulator(new MinimalState("",0,0,""),Array[Double](0,0))
  }

  override def add(min: MinimalState, acc : Accumulator) =  {
    val distance = min.getDistanceFromLatLonInKm(acc.min)
    var co2 : Double = 0
    if(distance>0){
      co2 = 0.001
    }
    acc.arr(0)=acc.arr(0)+co2
    acc.arr(1)=acc.arr(1)+distance
    acc.setMin(min)
    acc
  }

  override def getResult(acc:Accumulator)= (acc.arr(0),acc.arr(1),(System.currentTimeMillis() / 1000L))

  override def merge(a: Accumulator, b :Accumulator)=new Accumulator(a.min,Array(a.arr(0)+b.arr(0),a.arr(1)+b.arr(1)))
}
