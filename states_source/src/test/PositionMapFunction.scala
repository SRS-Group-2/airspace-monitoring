package it.unibo.states_source

import org.opensky.model.StateVector
import org.apache.flink.api.common.functions.MapFunction

class PositionMapFunction extends MapFunction[StateVector,MinimalState] {

    override def map(sv: StateVector): MinimalState = {
      
      return new MinimalState( 
        sv.getIcao24(), 
        sv.getLatitude(), 
        sv.getLongitude(),
        (System.currentTimeMillis() / 1000L).toString()
      )
    }
}

