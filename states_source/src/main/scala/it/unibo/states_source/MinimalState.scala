package it.unibo.states_source

class MinimalState(val icao: String, val latitude: Double, val longitude: Double/*longitude,val onGround: Boolean, val timestamp: String*/) {
  def toJSONString(): String = {
    return "{ \"icao24\": " + icao + 
           ", \"latitude\": " + latitude +
           ", \"longitude\": " + longitude +
          /*  ", \"onGround\": " + onGround +
           ", \"timestamp\": " + timestamp +*/
           "}"
  }

  override def toString(): String = this.toJSONString()
}