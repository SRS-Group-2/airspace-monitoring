package it.unibo.states_source


import com.google.firebase.cloud.FirestoreClient
import com.google.cloud.firestore.DocumentReference
import java.util.HashMap
import java.util.Map
import java.io.FileInputStream
import scala.collection.JavaConverters._
import com.google.api.core.ApiFuture
import com.google.cloud.firestore.WriteResult
import org.apache.flink.streaming.api.functions.sink.RichSinkFunction
import org.apache.flink.configuration.Configuration
import org.apache.flink.streaming.api.functions.sink.SinkFunction
import org.slf4j.Logger;
import org.slf4j.LoggerFactory;



@SerialVersionUID(100L)
class AircraftsFirebaseSink[IN] extends RichSinkFunction[Aircrafts] () {

  val LOG : Logger = LoggerFactory.getLogger(classOf[AircraftsFirebaseSink[IN]])


  override def open(parameters : Configuration) : Unit = {
    super.open(parameters)
  }

  override def close() : Unit ={
    super.close()
  }

  override def invoke(aircraft: Aircrafts, context: SinkFunction.Context) : Unit = {
    val instance = DbInstance.getInstance()
    val db = FirestoreClient.getFirestore()
    val docRef : DocumentReference  = db.collection("airspace").document("aircraft-list")
    val data : Map[String, Object]  = new HashMap[String, Object]()
    data.put("timestamp", aircraft.getTimestamp())
    data.put("icao24", aircraft.getList().asJava)
    val result : ApiFuture[WriteResult] = docRef.set(data)
    LOG.info("Aircrafts written on firestore")
  }

}
