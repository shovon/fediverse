package simpleld

// Simple LD is this insultingly superficial implementation of JSON-LD.
//
// The point of this library is to merely marshal JSON, with some context, and
// to unmarshal JSON, while parsing the context, allowing you to make an
// informed interpretation of a particular JSON document.
//
// To put it more succinctly, this library follows a more "data is meant to be
// read by humans, and only incidenttally for computers to interpret"
// philosophy.
//
// If you want a more proper JSON-LD implementation, this library is not it.
//
// This library will not filter out fields that are not defined by the context,
// thus violating Postel's law.
