package it.unibo.states_source

import java.io.FileInputStream
import com.google.cloud.firestore.Firestore
import com.google.cloud.firestore.FirestoreOptions
import com.google.auth.oauth2.GoogleCredentials
import com.google.api.gax.core.FixedCredentialsProvider

object  DbInstance {
  var instance : Firestore = null

  def getInstance() : Firestore = {
    synchronized {
      if (this.instance == null) {
        val options: FirestoreOptions = 
          FirestoreOptions.getDefaultInstance().toBuilder()
          // .setCredentials(GoogleCredentials.getApplicationDefault())
          .setCredentialsProvider(FixedCredentialsProvider.create(GoogleCredentials.getApplicationDefault()))
          .setProjectId(System.getenv("GOOGLE_CLOUD_PROJECT_ID"))
          .build()
          
        this.instance = options.getService()
       }
    }
    return this.instance
  }
}
