// Package jsonld provides some basic utilities to produce JSON-LD documents,
// and maybe to parse JSON-LD documents.
//
// Most other libraries will try to follow the "spirit" of JSON-LD, which is
// to say that they will try to produce JSON-LD documents that are as close
// as possible to the JSON-LD spec, but the problem is, consumers may not
// entirely honour the JSON-LD spec. This is why it's a good idea to grant the
// producer some flexibility in how they produce JSON-LD documents, in the case
// that there are clients that don't care about the specifics of JSON-LD.

package jsonld
