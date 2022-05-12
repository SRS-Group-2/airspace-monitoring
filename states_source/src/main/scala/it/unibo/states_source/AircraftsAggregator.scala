package it.unibo.states_source

import org.apache.flink.api.common.functions.AggregateFunction
import java.util.Date
import java.text.SimpleDateFormat
import java.time.Instant



class AircraftsAggregator extends AggregateFunction[String,List[String],Aircrafts] {

  override def createAccumulator() = List[String]()

  override def add(icao: String, acc: List[String]) = acc:+icao 

  override def getResult(acc:List[String])= {
    
    val myDate  = Date.from(Instant.now())
    val formatter = new SimpleDateFormat("yyyy-MM-dd-HH-mm");
    val timestamp = formatter.format(myDate);
    
    new Aircrafts(timestamp,acc.distinct)

    }

  override def merge(a: List[String],b:List[String])= a++b
}
