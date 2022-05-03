package it.unibo.states_source

import org.apache.flink.api.common.serialization.SerializationSchema

class JSONVectorsSerializer extends SerializationSchema[Vectors] {

override def serialize(vec: Vectors) :Array[Byte] = {
        return vec.toString().getBytes();
    }

}