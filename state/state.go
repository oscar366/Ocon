package state
//global var of the pointer wich is the location of the execute
var Pointer int
//this is a List that is the "sectionname":linenumber
var SectionList = make(map[string]int) 
//var storge cuz nothing better i know
var VarStorage = make(map[string]string)
//controlls debur prints
var DebugMode = true
