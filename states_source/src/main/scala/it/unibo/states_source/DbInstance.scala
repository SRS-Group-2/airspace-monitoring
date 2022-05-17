package it.unibo.states_source

import com.google.firebase.FirebaseApp
import java.io.FileInputStream
import com.google.firebase.FirebaseOptions
import com.google.auth.oauth2.GoogleCredentials



object  DbInstance {

  var instance : FirebaseApp = null

  def getInstance() : FirebaseApp = {
    synchronized {
      if(this.instance == null) {
        

        val options = 
          new FirebaseOptions.Builder()
          .setCredentials(GoogleCredentials.getApplicationDefault())
          .setProjectId(System.getenv("GOOGLE_CLOUD_PROJECT_ID"))
          .build()
          
        this.instance = FirebaseApp.initializeApp(options)
       }
    }

    return this.instance
  }


  
}
