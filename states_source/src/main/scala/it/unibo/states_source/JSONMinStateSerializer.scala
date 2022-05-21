package it.unibo.states_source

import org.apache.flink.api.common.serialization.SerializationSchema

class JSONMinStateSerializer extends SerializationSchema[MinimalState] {

  override def serialize(ms: MinimalState) :Array[Byte] = {
    return ms.toString().getBytes()
  }

}
