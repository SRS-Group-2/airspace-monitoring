package it.unibo.states_source

import java.io.Serializable


@SerialVersionUID(100L)
class Aircrafts (val timestamp: String, val list: List[String]) extends Serializable {

  def this() {
    this("",List[String]())
  }

  def getTimestamp() = timestamp

  def getList() = list

  def toJSONString(): String = {
    return "{ \"timestamp\": " + timestamp +
      ", \"list\": [" + list.mkString(",") +"]"
      "}"
  }

  override def toString(): String = this.toJSONString()
  
}
