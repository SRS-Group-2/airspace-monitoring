package it.unibo.states_source

import java.lang.Math
import java.io.Serializable

class MinimalState(val icao: String, val latitude: Double, val longitude: Double/*,val onGround: Boolean*/, val timestamp: String) extends Serializable {

  def getIcao() : String = {
    return icao
  }

  def getTimestamp() : String = {
    return timestamp
  }

  def toJSONString(): String = {
    return "{ \"icao24\": " +"\""+ icao +"\""+  
           ", \"lat\": " + latitude +
           ", \"lon\": " + longitude +
           ", \"timestamp\": " + timestamp +            
           "}"
  }

  override def toString(): String = this.toJSONString()


  def deg2rad(deg : Double) : Double = {
    return deg * (Math.PI/180)
  }

  def getDistanceFromLatLonInKm(other : MinimalState) : Int ={
    if(other.latitude==0) {
      return 0
    }
    val R = 6371 // Radius of the earth in km
    val dLat = deg2rad(other.latitude-latitude)  // deg2rad below
    val dLon = deg2rad(other.longitude-longitude)
    val a = 
      Math.sin(dLat/2) * Math.sin(dLat/2) +
      Math.cos(deg2rad(latitude)) * Math.cos(deg2rad(other.latitude)) * 
      Math.sin(dLon/2) * Math.sin(dLon/2)
      ;
    var c = 2 * Math.atan2(Math.sqrt(a), Math.sqrt(1-a))
    var d = R * c // Distance in km
    return d.toInt
  }

}
