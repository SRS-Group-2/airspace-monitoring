package it.unibo.states_source



class Aircrafts (val timestamp: String, val icao24: List[String]) {

    def toJSONString(): String = {
    return "{ \"timestamp\": " + timestamp +
           ", \"list\": [" + icao24.mkString(",") +"]"
           "}"
  }

  override def toString(): String = this.toJSONString()
  
}
