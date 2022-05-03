package it.unibo.states_source

import org.apache.flink.api.common.serialization.SerializationSchema

class JSONAircraftsSerializer extends SerializationSchema[Aircrafts]{

    override def serialize(air: Aircrafts) :Array[Byte] = {
        return air.toString().getBytes();
    }
  
}
