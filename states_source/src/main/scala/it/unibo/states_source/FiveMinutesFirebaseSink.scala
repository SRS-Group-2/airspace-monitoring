package it.unibo.states_source

import com.google.firebase.cloud.FirestoreClient
import com.google.cloud.firestore.DocumentReference
import java.util.HashMap
import java.util.Map
import scala.collection.JavaConverters._

import com.google.api.core.ApiFuture
import com.google.cloud.firestore.WriteResult

import org.apache.flink.streaming.api.functions.sink.RichSinkFunction
import org.apache.flink.configuration.Configuration
import org.apache.flink.streaming.api.functions.sink.SinkFunction
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;


@SerialVersionUID(100L)
class FiveMinutesFirebaseSink[IN] extends RichSinkFunction[(Int,Int,String)] (){

  val LOG : Logger = LoggerFactory.getLogger(classOf[FiveMinutesFirebaseSink[IN]])

  override def open(parameters : Configuration) : Unit = {
    super.open(parameters)
  }

  override def close() : Unit ={
    super.close()
  }

  override def invoke(res:(Int,Int,String), context: SinkFunction.Context) : Unit = {

    val instance = DbInstance.getInstance()
    val db = FirestoreClient.getFirestore()

    val docRef : DocumentReference  = 
      db.collection("airspace")
      .document("24h-history")
      .collection("5m-bucket")
      .document(res._3)

    val docRef1 : DocumentReference  = 
      db.collection("airspace")
      .document("5m-history")
      
    val data : Map[String, Any]  = new HashMap[String,Any]()
    data.put("CO2t",(res._1/1000).toInt)
    data.put("distanceKm",res._2)
    data.put("timestamp",res._3)
    val result1 : ApiFuture[WriteResult] = docRef1.set(data)
    val result : ApiFuture[WriteResult] = docRef.set(data)
    if(result.get()!=null && result1.get!=null) {
      LOG.info("FiveMins written on firestore")
    }
  }

}
