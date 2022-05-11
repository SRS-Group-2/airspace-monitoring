package it.unibo.states_source

import com.google.auth.oauth2.GoogleCredentials
import com.google.firebase.cloud.FirestoreClient
import com.google.cloud.firestore.DocumentReference
import java.util.HashMap
import java.util.Map
import java.io.FileInputStream
import com.google.firebase.FirebaseApp
import com.google.firebase.FirebaseOptions
import scala.collection.JavaConverters._

import com.google.api.core.ApiFuture
import com.google.cloud.firestore.WriteResult

import org.apache.flink.streaming.api.functions.sink.RichSinkFunction
import org.apache.flink.configuration.Configuration
import org.apache.flink.streaming.api.functions.sink.SinkFunction
import com.google.auth.oauth2.GoogleCredentials
import com.google.gson.Gson


class AircraftsFirebaseSink[IN] extends RichSinkFunction[Aircrafts] (){
var databaseUrl : String = null


override def open(parameters : Configuration) : Unit = {
super.open(parameters)

val serviceAccount =new FileInputStream(System.getenv("GOOGLE_APPLICATION_CREDENTIALS"));
    val options = new FirebaseOptions.Builder()
              .setCredentials(GoogleCredentials.fromStream(serviceAccount))
              .setProjectId("tokyo-rain-123-f6024")
              .build()
    FirebaseApp.initializeApp(options)

}

override def close() : Unit ={
        super.close();
}

override def invoke(aircraft: Aircrafts, context: SinkFunction.Context) : Unit = {
    
    
    
    val db = FirestoreClient.getFirestore()

    val docRef : DocumentReference  = db.collection("airspace").document("aircraft-list")
    val data : Map[String, Object]  = new HashMap[String, Object]();
    data.put("timestamp",aircraft.getTimestamp())
    data.put("icao",aircraft.getList().asJava)
    val result : ApiFuture[WriteResult] = docRef.set(data)



}


  
}


