package model


var E  = &Export{unexport:unexport{name: "3333"}}

type Export struct {
	unexport
}


type unexport struct {
	name string
}

