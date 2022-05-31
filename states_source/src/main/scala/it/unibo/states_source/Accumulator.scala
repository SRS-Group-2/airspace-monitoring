package it.unibo.states_source

class Accumulator (var min: MinimalState, var arr:Array[Double]){

    def setMin( m : MinimalState) : Unit = {
        this.min=m
    }
}
