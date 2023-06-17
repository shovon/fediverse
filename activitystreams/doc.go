package activitystreams

// What we want is to be able to parse a document formatted in accordance to the
// ActivityStreams 2.0 specification. It uses JSON-LD, and frankly, after
// reading up on it thoroughly, it's meant for machines. Humans are responsible
// for properly formatting the JSON-LD document, in accordance to the JSON-LD
// spec, ensuring that machines can understand it.
