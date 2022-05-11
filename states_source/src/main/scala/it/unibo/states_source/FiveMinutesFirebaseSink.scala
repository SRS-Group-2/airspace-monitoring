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


class FiveMinutesFirebaseSink[IN] extends RichSinkFunction[(Double,Double,String)] (){
var databaseUrl : String = null


override def open(parameters : Configuration) : Unit = {
super.open(parameters)
val serviceAccount =new FileInputStream(System.getenv("GOOGLE_APPLICATION_CREDENTIALS"));
    val options = new FirebaseOptions.Builder()
              .setCredentials(GoogleCredentials.fromStream(serviceAccount))
              .setProjectId("tokyo-rain-123-f6024")
              .build()
    FirebaseApp.initializeApp(options,"oneHour");



}

override def close() : Unit ={
        super.close();
}

override def invoke(res:(Double,Double,String), context: SinkFunction.Context) : Unit = {
    
    
    
    val db = FirestoreClient.getFirestore(FirebaseApp.getInstance("oneHour"));

    val docRef : DocumentReference  = db.collection("airspace").document("1h-history").collection("5m-bucket").document(res._3)
    val docRef1 : DocumentReference  = db.collection("airspace").document("5m-history")
    val data : Map[String, Any]  = new HashMap[String,Any]();
    data.put("CO2t",res._1)
    data.put("distanceKm",res._2)
    data.put("timestamp",res._3)
    val result1 : ApiFuture[WriteResult] = docRef1.set(data)
    val result : ApiFuture[WriteResult] = docRef.set(data)

    




}


  
}


