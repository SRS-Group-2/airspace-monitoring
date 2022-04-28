package it.unibo.states_source

class Vectors(val minstates: List[MinimalState],val timestamp: String){
    def toJSONString(): String = {
    return "{ \"minstates\": [" + minstates.mkString(",") + "]" +
           ", \"timestamp\": " + timestamp +
           "}"
  }

  override def toString(): String = this.toJSONString()
}