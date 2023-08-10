// Package jsonld provides some basic utilities to produce JSON-LD documents,
// and maybe to parse JSON-LD documents.
//
// Most other libraries will try to follow the "spirit" of JSON-LD, which is
// to say that they will try to produce JSON-LD documents that are as close
// as possible to the JSON-LD spec, but the problem is, consumers may not
// entirely honour the JSON-LD spec. This is why it's a good idea to grant the
// producer some flexibility in how they produce JSON-LD documents, in the case
// that there are clients that don't care about the specifics of JSON-LD.
//
// This package is pretty much a living library to make working with JSON-LD
// easier. It's not going to behave "robotically". Instead, it will only aid
// in guiding us to be as compliant with JSON-LD as possible, but also grant us
// flexibility for clients that ignore JSON-LD contexts.

package jsonld
