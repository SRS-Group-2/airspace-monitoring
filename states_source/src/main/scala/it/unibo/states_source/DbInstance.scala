package it.unibo.states_source

import java.io.ByteArrayInputStream
import com.google.cloud.firestore.Firestore
import com.google.cloud.firestore.FirestoreOptions
import com.google.auth.Credentials
import com.google.auth.oauth2.GoogleCredentials
import com.google.api.gax.core.FixedCredentialsProvider

object  DbInstance {
  var instance : Firestore = null

  def getInstance() : Firestore = {
    synchronized {
      if (this.instance == null) {
        val credentials: Credentials = 
          if (System.getenv("FIRESTORE_AUTHENTICATION_METHOD") == "ADC") {
            GoogleCredentials.getApplicationDefault()
          } else {
            val json = System.getenv("FIRESTORE_CREDENTIALS")
            val jsonCredentials = new ByteArrayInputStream(json.getBytes())
            GoogleCredentials.fromStream(jsonCredentials)
          }
        val options: FirestoreOptions = 
          FirestoreOptions.getDefaultInstance().toBuilder()
          // .setCredentials(GoogleCredentials.getApplicationDefault())
          .setCredentialsProvider(FixedCredentialsProvider.create(credentials))
          .setProjectId(System.getenv("GOOGLE_CLOUD_PROJECT_ID"))
          .build()
          
        this.instance = options.getService()
       }
    }
    return this.instance
  }
}
