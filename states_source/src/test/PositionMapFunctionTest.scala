package it.unibo.states_source

import org.opensky.model.StateVector
import org.scalatest.FlatSpec
import org.scalatest.Matchers
import org.scalatest.BeforeAndAfter
import org.apache.flink.streaming.util.OneInputStreamOperatorTestHarness
import org.apache.flink.streaming.api.operators.StreamFlatMap
import org.apache.flink.streaming.runtime.streamrecord.StreamRecord

class PositionMapFunctionTest extends FlatSpec with Matchers {

    "PositionMapFunction" should "transform a StateVector in a MinimalState" in {
        // instantiate your function
        val posMapFun = new PositionMapFunction()

        val sv= new StateVector("abcdef")
        sv.setLatitude(23)
        sv.setLongitude(46)

        val ms = new MinimalState("abcdef",23,46,(System.currentTimeMillis() / 1000L).toString())

        // call the methods that you have implemented
        posMapFun.map(sv) should equal (ms)
    }
}
