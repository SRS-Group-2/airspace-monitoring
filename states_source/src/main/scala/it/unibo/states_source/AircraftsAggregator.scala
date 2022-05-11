package it.unibo.states_source

import org.apache.flink.api.common.functions.AggregateFunction



class AircraftsAggregator extends AggregateFunction[String,List[String],Aircrafts] {

  override def createAccumulator() = List[String]()

  override def add(icao: String, acc: List[String]) = acc:+icao 

  override def getResult(acc:List[String])= new Aircrafts((System.currentTimeMillis() / 1000L).toString(),acc.distinct)

  override def merge(a: List[String],b:List[String])= a++b
}
