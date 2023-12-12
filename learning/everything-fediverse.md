# A primer on how to work with the Fediverse

## ActivityPub

The most important component of the Fediverse (beyond the Internet and HTTP) is ActivityPub.

ActivityPub's primitive is the "actor", and in any ActivityPub-based Fediverse software, a user is associated with an actor, and it's that actor that users will interact with, via their own actors, by sending each other activities, via HTTP.

ActivityPub's primary data interchange format is JSON-LD. So an actor is encoded as a JSON-LD document, and so are the _activities_ that actors send to other actors.

## JSON-LD

JSON-LD may look like JSON, but the JSON is merely the container format to represent "linked data" (LD).

To put it simply, LD is a standard for describing nodes (subject), and their relationship (predicates) to other nodes (object).

The caveat with using LD is that fields should be defined not as human-readable names (for example, `name` or `address`) but as URLs (for example, `https://example.com/ns#name` or `https://example.com/ns#address).

So, if we have an organization that owns the domain name `example.com`, they can use the domain name to express "ownership" of those fields.

Additionally, those fields can point to more than one object. And hence a JSON-LD fields are associated with an array of nodes, rather than just a single node.

Here's what a "real" JSON-LD document would look like:

```json
{
	"@id": "https://example.com/api/people/1",
	"https://example.com/ns#name": [
		{
			"@value": "John Doe"
		}
	],
	"https://example.com/ns#address": [
		{
			"@value": "123 Peachtree Avenue"
		}
	]
}
```

Of course, that looks verbose.

Fortunately, you can define the "schema" or "vocabulary" that the JSON-LD document can be interpreted with, and you do so via the `@context` field right from within the node itself.

For example:

```json
{
	"@context": {
		"ex": "https://example.com/ns#",
		"name": "ex:name",
		"address": "ex:address"
	},
	"@id": "https://example.com/api/people/1",
	"name": "John Doe",
	"address": "123 Peachtree Avenue"
}
```

> [!Note]
> Notice the `@id` field? That is a special field that represents the **node's ID**.
>
> Remember, JSON-LD, and LD are used to define relationships between nodes. The ID helps identify the source node.

In the above, you can see that `ex` is an alias for `https://example.com/ns#`, `name` is an alias for `ex:name`, which in turn is an alias for `https://example.com/ns#name`.

To get back the original "predicates", you can throw your JSON-LD document into an expander.

The above document will expand to become:

```json
[
	{
		"@id": "https://example.com/api/people/1",
		"https://example.com/ns#name": [
			{
				"@value": "John Doe"
			}
		],
		"https://example.com/ns#address": [
			{
				"@value": "123 Peachtree Avenue"
			}
		]
	}
]
```

As you can see, the root-level document is expanded and placed in an array. This is because, again, JSON-LD is a way to represent a graph of linked data, and a single JSON-LD document (or any document representing linked data) can have multiple nodes in the graph.

For example, that above node (identified as `https://example.com/api/people/1`), can have a single field (identified as `https://example.com/ns#name`), have mulitple values.

### The context as a URL

The `@context` can also be a URL to a JSON document that actually describes the schema/vocabulary.

For example, let's say `https://example.com/ns` actually points to another JSON-LD object that contains a `@context` at its root.

```json
// This can potentially be located at https://example.com/ns
{
	"@context": {
		"ex": "https://example.com/ns#",
		"name": "ex:name",
		"address": "ex:address"
	}
}
```

Then you can simply substitute the context in your document with the URL to `https://example.com/ns`.

```json
{
	"@context": "https://example.com/ns",
	"@id": "https://example.com/api/people/1",
	"name": "John Doe",
	"address": "123 Peachtree Avenue"
}
```

Again, that above document will expand to what we saw earlier, but this time, most expanders will do an additional lookup over at `https://example.com/ns` to retrieve and interpret the context.

> [!Note]
> Depending on the expansion library that you use, you can also supply your own contexts.
>
> This is especially useful for applications where each actor in a networking application already knows what the context is, and so senders are free to omit the `@context` field, if they so choose.
>
> This is especially important with ActivityPub, since the specification states that absent the context, then interpretation (typically expansion, among others) must be done with the ActivityStreams context.

> [!Note]
> As far as JSON-LD interpreters are concerned, a context pointed to by a URL is the only thing that will only ever trigger a network reqeuest. Beyond that JSON-LD is merely a data interchange format, and any application-level inconsistencies must be handled between clients/servers.

### Multiple contexts

In JSON-LD, you have the ability to provide multiple contexts.

For example,

```json
{
	"@context": [
		{
			"ex1": "https://example.com/ns1#",
			"name": "ex1:name"
		},
		{
			"ex2": "https://example.com/ns2#",
			"address": "ex2:address"
		}
	],
	"name": "Jane Doe",
	"address": "123 Peachtree Avenue"
}
```

And, the above document with two contexts will resolve like so:

```json
[
	{
		"https://example.com/ns2#address": [
			{
				"@value": "123 Peachtree Avenue"
			}
		],
		"https://example.com/ns1#name": [
			{
				"@value": "Jane Doe"
			}
		]
	}
]
```

There are many reasons why you would want multiple contexts.

One very important use case is that two parts of an application can expect two entirely distinct set of vocabularies.

While one component can read fields that conforms to one set of vocabulary, the other can read from the other set of vocabulary.

### Practicing and debugging JSON-LD

Some may find JSON-LD to be rather confusing. Fear not, you are not alone.

Fortunately, there is a tool for you to explore the various ways to work with JSON-LD.

Head on over to [json-ld.org/playground](https://json-ld.org/playground/), and start playing around.

Bear in mind: for most applications, you're probably going to only be expanding JSON-LD, but that said, it's probably a good idea to teach yourself the motivation behind JSON-LD. As a quick hint: JSON-LD is just one way to deliver RDF, which, again, is a way to establish relationships from node to node.

As you are reading through the next sections, I highly encourage you to copy and paste the JSON-LD documents into the JSON-LD playground to get a feel for what JSON-LD truly is.

### JSON-LD and triples

Before we go any further, let me go ahead and explain what JSON-LD really is.

Earlier I mentioned that JSON-LD is a way to represent a "subject -> predicate -> object" relationship.

That relationship is what is called a "triple". And we can have multiple of these triples, all thrown into a pile. Reading them individually, tracing their paths, will eventually allow you to form of a graph.

To actually make this concept of a pile of triples be actually practicable, the same subject, can have multiple predicate -> object relationships. In other words, we can have the same subject be used across multiple triples, to represent a single entity within what the graph represents.

So, let's say we have someone named Alice, and she has a house on 123 Peachtree Avenue, and her absolute favourite colour is purple.

In a very simplified syntax, the series of triples to describe Alice will look like so:

```
https://example.com/Alice https://example.com/address "123 Peachtree Avenue"
https://example.com/Alice https://example.com/color "Purple"
```

> [!Note]
> The above two triples don't need to exclusively be a part of a single file.
>
> In fact, if you wanted to, you can throw those triples into a database of triples, that can then be used to represent a graph.
>
> In fact, there are database implementations out there that are specifically geared towards storing these triples, and the class of databases responsible for that are called "triple stores", and often resort to using a query language called SPARQL.

Repeating the subject—as represented by the ID `https://example.com/Alice`—becomes repetitive.

This is where JSON-LD comes along to aleviate that repetitiveness.

The above set of triples can be instead rewritten as so:

```json
{
	"@id": "https://example.com/Alice",
	"https://example.com/address": [
		{
			"@value": "123 Peachtree Avenue"
		}
	],
	"https://example.com/color": [
		{
			"@value": "Purple"
		}
	]
}
```

And of course, with a `@context`, we can alleviate the repetitiveness even more.

```json
{
	"@context": {
		"ex": "https://example.com/",
		"address": "ex:address",
		"color": "ex:color"
	},
	"@id": "https://example.com/Alice",
	"address": "@value": "123 Peachtree Avenue",
	"color":  "Purple"
}
```

### Blank nodes

Previously, we had our root (subject) node have an explicit ID. But sometimes, when interpreting a single document, using IDs may be overkill. Sometimes, we just want to deliver fields to those reading and interpreting the JSON-LD document. In this case we can omit the ID.

If we are to omit the ID of the node, we say that the node is a blank node.

By convention, in many triples syntax, we still need to identify a blank node, and that is done by a prefixing a label with a `_:`.

Taking a previous JSON-LD example, let's omit Alice's ID, and instead give her a name. We can have her blank node ID be `_:alice`, and the set of triples in our custom syntax will look like so.

```
_:alice https://exmaple.com/name "Alice"
_:alice https://example.com/address "123 Peachtree Avenue"
_:alice https://example.com/color "Purple"
```

In JSON-LD, however, because we don't need to explicitly identify blank nodes, we can simply omit the `_:alice` ID.

```json
{
	"https://exmaple.com/name": [
		{
			"@value": "Alice"
		}
	],
	"https://example.com/address": [
		{
			"@value": "123 Peachtree Avenue"
		}
	],
	"https://example.com/color": [
		{
			"@value": "Purple"
		}
	]
}
```

And of course, using the context, the above can be abbreviated like so:

```json
{
	"@context": {
		"ex": "https://example.com/",
		"name": "ex:name",
		"address": "ex:address",
		"color": "ex:color"
	},
	"name": "Alice",
	"address": "123 Peachtree Avenue",
	"color": "Purple"
}
```

### Literals and Nodes

So far, I've been talking about associating each "field" (predicate) to a literal (object). However, one of the powers of JSON-LD and the triples that it represents is that you can form a graph!

That means a subject can point to an object, and that object can represent a subject and then points to another object, and so on and so forth.

Going back to our triples example, we can give Alice a dog

```
_:alice https://example.com/dog _:waffles
_:waffles https://example.com/name "Waffles"
```

And an equivalent JSON-LD would look like so:

```json
{
	"https://example.com/dog": [
		{
			"https://example.com/name": [
				{
					"@value": "Waffles"
				}
			]
		}
	]
}
```

And, to clean things with the help of the context, we now introduce an additional bit of information to spcify that the `"https://example.com/dog"` field points to a non-literal.

```json
{
	"@context": {
		"ex": "https://example.com/",
		"dog": {
			"@id": "ex:dog",
			"@type": "@id"
		}
	},
	"dog": {
		"name": "Waffles"
	}
}
```

Of course, if you wanted to uniquely identify Alice's dog in a global pool of triples, you can do so via the `@id` field, just as you would with the root node.

```json
{
	"@context": {
		"ex": "https://example.com/",
		"dog": {
			"@id": "ex:dog",
			"@type": "@id"
		},
		"name": "ex:name"
	},
	"dog": {
		"@id": "https://example.com/Waffles",
		"name": "Waffles"
	}
}
```

Not much would change in the resulting triples interpretation, but, instead of using the blank node syntax to represent Waffles, we'd use an actual URL.

```
_:alice https://example.com/dog https://example.com/Waffles
https://example.com/Waffles https://example.com/name "Waffles"
```

> [!Note]
> We are not using quotes around `https://example.com/Waffles`, because that's not a literal, but something that will point to something else.

### The ID Node

In the earlier example, we gave the node identified initially as a blank node, then later re-identified to "https://example.com/Waffles". That node had a single field `https://example.com/dog` (or just `dog`).

The thing is, some applications are perfectly fine if you omit all the fields entirely.

Those applications would greatly benefit from simply downloading the data that is potentially located by the value in the `@id`.

So taking the original JSON-LD document from the last example, we can omit the `https://example.com/name` (or just `name`) field, and that application will work perfectly fine, since it was written to explicitly try to look up the value associated with the ID of the node.

So this is perfectly fine in such a hypothetical application.

```json
{
	"@context": {
		"ex": "https://example.com/",
		"dog": {
			"@id": "ex:dog",
			"@type": "@id"
		}
	},
	"dog": {
		"@id": "https://example.com/Waffles"
	}
}
```

It gets even better than that.

Above, we aliased `https://example.com/dog` as a field named `dog`, but also, that that field is in fact pointing to another non-literal node.

Because of that, we JSON-LD, we can simply move the ID to just a value.

Like so:

```json
{
	"@context": {
		"ex": "https://example.com/",
		"dog": {
			"@id": "ex:dog",
			"@type": "@id"
		}
	},
	"dog": "https://example.com/Waffles"
}
```

And the application can simply look up the document that the aliased `dog` field points to, if it so chooses.

In case you're wondering, that above document will expand to this:

```json
[
	{
		"https://example.com/dog": [
			{
				"@id": "https://example.com/Waffles"
			}
		]
	}
]
```

Also, if we wanted an entity to have multiple dogs associated with it (after all, people do own more than one dog), using that compacted syntax, we can simply have it all in an array.

```json
{
	"@context": {
		"ex": "https://example.com/",
		"dog": {
			"@id": "ex:dog",
			"@type": "@id"
		}
	},
	"dog": ["https://example.com/Waffles", "https://example.com/Milo"]
}
```

And it would expand to this:

```json
[
	{
		"https://example.com/dog": [
			{
				"@id": "https://example.com/Waffles"
			},
			{
				"@id": "https://example.com/Milo"
			}
		]
	}
]
```

### Most URLs don't need to resolve to anything

So you noticed that a lot of things in the previous JSON-LD documents are URLs. Not only are the subjects often represented as URLs, and not only are some objects represented as URLs to other nodes, but so are the predicates!

But here's the thing: URLs don't need to _actually_ resolve (with the exception of a context pointed to by a URL). That is, `https://example.com` doesn't even need to exist on the Internet!

This is why the convention behind defining a predicate uses a `#` symbol (such as `https://example.com#fieldName`); that's the part of a URL that is ignored when looking up a resource on the internet.

So, for example `https://example.com#fieldName` simply resolves to `https://example.com`, because the content including and after the `#` symbol don't matter.

In JSON-LD, the vast majority of use cases of URLs are to uniqely identify things, without actually being able to locate them on the Internet!

Whether or not a particular URL resolves in a JSON-LD document is a matter of settlement exclusively between producer and consumer of the document. If the consumer expects to be able to resolve a URL, then this is a problem between the consumer and producer, and it's up to them to figure it out.

> ![Note]
> Yes, I did say almost all URLs in a JSON-LD document need not be able to be located on the Internet.
>
> That there is only **one** exception, and that is the `@context`. If there is a context defined in a document found on the Internet, then JSON-LD interpreters **MUST** be able to download it in order to make sense of the document that has its context pointing to another another resource on the Internet.
>
> For example, the URL in the context **MUST** exist, otherwise, the document is deemed invalid.
>
> ```json
> {
> 	"@context": "https://example.com/ns"
> }
> ```
>
> If `https://example.com/ns` does not point to a JSON-LD document that represents a valid context, then the JSON-LD document is invalid!
>
> Yes, this means JSON-LD clients that don't have access to the Internet will not be able to process the document.

### URLs, URIs, and IRIs in JSON-LD

In the previous section, I mentioned how most of the URLs in a JSON-LD document does not need to actually resolve into anything. It's up to the consumer of the document to determine if the document is valid or not.

Hence why not only can subject, predicate, and objects be URLs, but they can also be URIs. JSON-LD works perfectly fine with URIs

An example of such a URI would be a "mailto" URI, such as `urn:isbn:0-486-27557-4`. That URI doesn't point to anything on the Internet.

And then you have IRIs.

IRIs are a superset of URIs.

JSON-LDs also work with IRIs.

In fact, by the JSON-LD specification, IRIs are what JSON-LD works with. And IRIs are a superset of URIs, and URIs are a superset of URLs. Hence why JSON-LD works perfectly fine with URLs.

### The `@type` field

Even though the `@type` field doesn't play _that_ major of a role in terms of interpreting an LD node, it's still there, in case you need it. This way, you don't need to describe your own custom predicate to describe the "type" of a node.

```json
{
	"@type": "https://example.com/ns#Person"
}
```

You could even alias a type. So rather than entering a whole URL, you can instead use the `@context` to define something shorter.

For example:

```json
{
	"@context": {
		"ex": "https://example.com/ns#",
		"Person": "ex:Person"
	},
	"@type": "Person"
}
```

Expanding the above should yield

```json
[
	{
		"@type": ["https://example.com/ns#Person"]
	}
]
```

Notice that `@type` is expanded into an array?

This is because a single node can represent more than one type.

```json
{
	"@context": {
		"ex": "https://example.com/ns#",
		"Person": "ex:Person",
		"Employee": "ex:Employee"
	},
	"@type": ["Person", "Employee"]
}
```

### Aliasing URIs

Everything that is represented by a URI, such as `@id`s, fields, and `@type`s, can be aliased in the from inside the `@context`.

```json
{
	"@context": {
		"ex": "https://example.com/"
	},
	"ex:ns#cool": {
		"@id": "ex:"
	}
}
```

Will expand to:

```json
[
	{
		"https://example.com/ns#cool": [
			{
				"@id": "https://example.com/"
			}
		]
	}
]
```

> [!Note]
> JSON-LD doesn't only work with URIs. Instead, it works with a superset of URIs called Internationalized Resource Identifier, or IRI for short. While URIs only support ASCII, IRIs support unicode.

### Aliasing the `@type` and `@id` field.

JSON-LD allows authors to alias the `@type` and `@id` fields to anything that they so choose. Often times, authors dislike the `@` symbol.

For that reason, many authors take advantage of JSON-LD's capability to alias those fields, and the syntax for it looks like so:

```json
{
	"@context": {
		"ex": "https://example.com/ns#",
		"id": "@id",
		"type": "@type",
		"Object": "ex:Object",
		"Nothing": "ex:Nothing"
	},

	"id": "https://example.com/api/objects/1",
	"type": ["Object", "Nothing"]
}
```

And as you can see above, no need to explicitly include the `@` symbol right befre `type` or `id`.

### Compacted vs Expanded form

When people talk about JSON-LD, they typically talk about its two forms: compacted and expanded.

As far as interpreters are concerned, both forms are perfectly valid JSON-LD. Not only can you expand a compacted JSON-LD document, but you can also expand an expanded document. Not only can you compact an expanded JSON-LD document, you can compact a compacted JSON-LD document, as long as you provide the valid contexts.

You can even even mix and match the compacted and expanded forms. As long as you follow the "rules" of JSON-LD, you are dealing with perfectly valid JSON-LD documents.

I even demonstrated this in a previous example. Here it is again:

```json
{
	"@context": {
		"ex": "https://example.com/ns#",
		"name": "ex:name",
		"address": "ex:address"
	},
	"@id": "https://example.com/api/people/1",
	"name": "John Doe",
	"address": "123 Peachtree Avenue",
	"https://example.com/ns#dogs": [
		{
			"@id": "https://example.com/api/dogs/1",
			"name": "Waffles"
		}
	]
}
```

As you can see, what could have otherwise been aliased as `dogs`, and referenced as so, instead, we didn't bother aliasing, and directly used the URI as the predicate.

This is perfectly valid JSON-LD.

Just be sure to expand, before interpreting the document.

And, just as a courtesy, be sure to alias the fields, via the `@context`, compact the document, before sending it out to an intended recipient. This way, you won't need to repeat a whole URL prefix for every field. This should help save bandwidth.

## ActivityPub and ActivityStreams Administrivia

A "field" in JSON-LD is an IRI, and not the human-readable field names that everyone is used to.

For example, the field `inbox` technically doesn't make sense in ActivityPub, because JSON-LD expanders will ignore that field—without a valid alias.

For that reason, when talking about fields in these paragraphs, rather than printing the entire field (predicate) name as an IRI, instead, I will prefix with `as`, which aliases `https://www.w3.org/ns/activitystreams#`.

For example, rather than writing out `http://www.w3.org/ns/ldp#inbox`, I will write out `as:inbox`, which aliases `https://www.w3.org/ns/activitystreams#inbox`, which in turn aliases `ldp:inbox`, and `ldp` aliases `http://www.w3.org/ns/ldp#`.

That said, in JSON form, explicit aliasing is not necessary, because JSON-LD expanders are capable of resolving the so-called "human-readable" field names perfectly fine, given the appropriate aliases in the contexts.

So while I'd write `as:inbox` in these paragraphs, in JSON, as long as I provide the appropriate context, I'd simply write `inbox`, like so:

```json
{
	"@context": "https://www.w3.org/ns/activitystreams",
	"inbox": "https://sources.example.com/actors/1/inbox"
}
```

So remember: we are assuming the context derived from `https://www.w3.org/ns/activitystreams`, which is aliased to `as`. When talking about individual "fields" (wherein people are intuitively going to expect a human-readable field), to avoid any confusion, I will be prefixing with `as:` rather than spell out the field name itself. But in JSON, just use the "human-readable" field name.

## Actor

An actor is represented by a [JSON-LD](https://json-ld.org/) document. An actor does not need to be too complicated.

Here's a barebones actor.

```json
{
	"@context": "https://www.w3.org/ns/activitystreams",
	"id": "https://source.example.com",
	"inbox": "https://sources.example.com/inbox",
	"outbox": "https://sources.example.com/outbox",
	"following": "https://sources.example.com/following",
	"followers": "https://sources.example.com/followers",
	"liked": "https://sources.example.com/liked"
}
```

As long as you know the HTTP URL to the actor, you should be able to get a JSON-LD document that looks like the above, plus some additional details.

### Actors in Mastodon

Mastodon is a bit more pickier about what an actor is.

For one, an actor _must_ be any of the [ActivityStreams Actor Types](https://www.w3.org/TR/activitystreams-vocabulary/#actor-types) (via the JSON-LD `"@type"` property, which ActivityStreams aliases to just `"type"`).

Secondly, a username must be supplied via the ActivityPub `preferredUsername` field, otherwise, Mastodon is going to ignore the actor entirely.

And finally, your actors _must_ also be associated with an RSA 2048-bit public key, that will be used to verify any activity coming in. The specification used for verification is called `draft-cavage-http-signatures-12`. Because an actor is defined with JSON-LD, we must specify that the fields come from the `https://w3id.org/security/v1` schema.

So, taking the actor from earlier, and then adding the necessary fields in order to conform to Mastodon's expectations, the actor should now look like so:

```json
{
	"@context": [
		"https://www.w3.org/ns/activitystreams",
		"https://w3id.org/security/v1"
	],
	"id": "https://source.example.com",
	"type": "Person",
	"inbox": "https://sources.example.com/inbox",
	"outbox": "https://sources.example.com/outbox",
	"following": "https://sources.example.com/following",
	"followers": "https://sources.example.com/followers",
	"liked": "https://sources.example.com/liked",
	"preferredUsername": "actor",
	"publicKey": {
		"id": "https://source.example.com#main-key",
		"owner": "https://source.example.com",
		"publicKeyPem": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArLEIhmSM4UXoUbh/UNri\nOmsruokiG4GU0jz7R/rZ3lC0kGEMEJpk7x8hLEtg0DhV9DW3jPOsPi1KvLRkTgiE\nCSEEG+ULqZ3/WTZR3VX+/Tb1huemD2rBZkv9vpL+3qSRuFTvcMumonVuJ6rtT3pG\nTbsXlYmp2n7VkbPQPz6Wy3R7YeGmdNxtRiccwrpeovc+kCCoY/t467cK1ON+FDrq\nT/xgNhG2jPfotMF3ixk5/EQuakKEz2YQP4duD6D86QciZQWjw5YMv96NxV6D24CV\nn8HxEcxM5AfWvqbNLpEvi6UBUVCnM4IzJTlboPBO4tUPSu01YDqb8jbTC0f6rOCZ\nOQIDAQAB\n-----END PUBLIC KEY-----\n"
	}
}
```

## Following someone

When following an actor, you would send a "follow activity". It usually looks like this:

```json
{
	"@context": "https://www.w3.org/ns/activitystreams",
	"id": "https://source.example.com#follow/1",
	"type": "Follow",
	"actor": "https://source.example.com",
	"object": "https://destination.example.com"
}
```

However, that follow activity is merely a request to follow.

Your follow request is therefore merely "pending".

If the followee is willing to welcome the prospective follower to become an actual follower, then the followee is responsible for responding with an "accept activity". It typically looks like so:

```json
{
	"@context": "https://www.w3.org/ns/activitystreams",
	"id": "https://destination.example.com#accepts/follows/1",
	"type": "Accept",
	"actor": "https://destination.example.com",
	"object": {
		"id": "https://source.example.com#follow/1",
		"type": "Follow",
		"actor": "https://source.example.com",
		"object": "https://destination.example.com"
	}
}
```

Ideally, the followee should store the follow activity's `@id` for record-keeping, and the follower should do likewise.

## Unfollow someone

When unfollowing an actor, the unfollower must send an "undo activity".

```json
{
	"@context": "https://www.w3.org/ns/activitystreams",
	"id": "https://destination.example.com#accepts/follows/1",
	"type": "Undo",
	"actor": "https://destination.example.com",
	"object": {
		"id": "https://source.example.com#follow/1",
		"type": "Follow",
		"actor": "https://source.example.com",
		"object": "https://destination.example.com"
	}
}
```
