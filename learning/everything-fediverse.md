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

Additionally, those fields can point to more than one object. And hence a JSON-LD fields are associated with an array of nodes, rather than an just a single node.

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

> [!Note]
> Depedning on the expansion library that you use, you can also supply your own contexts.
>
> This is especially useful for applications where each actor in a networking application already knows what the context is, and so senders are free to omit the context, if they so choose.
>
> This is especially important with ActivityPub, since the specification states that absent the context, then interpretation (typically expansion, among others) must be done with the ActivityStreams context.

In the above, you can see that `ex` is an alias for `https://example.com/ns#`, `name` is an alias for `ex:name` (which in turn is an alias for `https://example.com/ns#name`).

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

As you can see, the root-level document is expanded and placed in an array.

The `@context` can also be a URL to a JSON document that actually describes the schema/vocabulary.

For example, let's say `https://example.com/ns` actually points to another JSON-LD object that contains a `@context` at its root.

```json
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

### Actual Linked Data

LD not only links a node to other nodes, via a subject -> predicate -> object relationship, but of course, as the name implies, it also links data!

Let's take the above document, add a field that gives the person one or many dogs.

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

Here, we have the root node (subject) that points to—at least—another node, as predicated by `https://example.com/ns#dogs` (object).

Notice now you have this ugly `"https://example.com/ns#dogs"`? Let's clean that up by moving it to the context.

```json
{
	"@context": {
		"ex": "https://example.com/ns#",
		"name": "ex:name",
		"address": "ex:address",
		"dogs": {
			"@type": "@id",
			"@id": "ex:dogs"
		}
	},
	"@id": "https://example.com/api/people/1",
	"name": "John Doe",
	"address": "123 Peachtree Avenue",
	"dogs": [
		{
			"@id": "https://example.com/api/dogs/1",
			"name": "Waffles"
		}
	]
}
```

The `{"@type": "@id", "@id": "ex:dogs"}`, pretty much describes the field associated with `https://example.com/ns#` to be a full node, rather than a value node.

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
		],
		"https://example.com/ns#dogs": [
			{
				"https://example.com/ns#name": [
					{
						"@id": "https://example.com/api/dogs/1",
						"@value": "Waffles"
					}
				]
			}
		]
	}
]
```

Given that it is generally a good idea to first expand a JSON-LD document prior interpreting it, we can simply move `{"name": "Waffles"}`, out of the array, and into a single non-array value.

```json
{
	"@context": {
		"ex": "https://example.com/ns#",
		"name": "ex:name",
		"address": "ex:address",
		"dogs": {
			"@type": "@id",
			"@id": "ex:dogs"
		}
	},
	"@id": "https://example.com/api/people/1",
	"name": "John Doe",
	"address": "123 Peachtree Avenue",
	"dogs": {
		"@id": "https://example.com/api/dogs/1",
		"name": "Waffles"
	}
}
```

And it will expand the same way had that one node linked with `dogs` been inside the array in the first place.

If we wanted to, we can then move that single document in `dogs` to an entirely separate document, somewhere on the Internet, pointed to by `https://example.com/api/dogs/1`.

The dog:

```json
{
	"@context": {
		"name": "https://example.com/ns#name"
	},
	"@id": "https://example.com/api/dogs/1",
	"name": "Waffles"
}
```

Its owner:

```json
{
	"@context": {
		"ex": "https://example.com/ns#",
		"name": "ex:name",
		"address": "ex:address",
		"dogs": {
			"@type": "@id",
			"@id": "ex:dogs"
		}
	},
	"@id": "https://example.com/api/people/1",
	"name": "John Doe",
	"address": "123 Peachtree Avenue",
	"dogs": "https://example.com/api/dogs/1"
}
```

Bear in mind, unlike resolving the context, which involves making a request over the Internet, a single ID field will not yield any such requests, during the expansion.

So, for example, of the following two documents are not equivalent.

```json
{
	"@context": {
		"ex": "https://example.com/ns#",
		"name": "ex:name",
		"dogs": {
			"@type": "@id",
			"@id": "ex:dogs"
		}
	},
	"dogs": {
		"@id": "https://example.com/api/dogs/1",
		"name": "Waffles"
	}
}
```

```json
{
	"@context": {
		"dogs": {
			"@type": "@id",
			"@id": "https://example.com/ns#dogs"
		}
	},
	"dogs": "https://example.com/api/dogs/1"
}
```

The first one will expand to:

```json
[
	{
		"https://example.com/ns#dogs": [
			{
				"@id": "https://example.com/api/dogs/1",
				"https://example.com/ns#name": [
					{
						"@value": "Waffles"
					}
				]
			}
		]
	}
]
```

And the second one to:

```json
[
	{
		"https://example.com/ns#dogs": [
			{
				"@id": "https://example.com/api/dogs/1"
			}
		]
	}
]
```

It doesn't matter whether `https://example.com/api/dogs/1` points to a document that is represented by the object represented by `dogs` in the first of the two above documents, in the end of the day, the responsibility lies squarely on the interpreter of the document. If the interpreter prefers to always lookup the document associated with the `@id`, then they can do so, otherwise, they are also free to interpret `dogs`, as-is, even if it is missing the `ex:name` field.

The first one will expand to:

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

## ActivityPub and ActivityStreams Administrivia

A "field" in JSON-LD is an IRI, and not the human-readable field names that everyone is used to.

For example, the field `inbox` technically doesn't make sense in ActivityPub, because JSON-LD expanders will ignore that field—without a valid alias.

For that reason, when talking about fields in these paragraphs, rather than printing the entire field (predicate) name as an IRI, instead, I will prefix with `as`, which aliases `https://www.w3.org/ns/activitystreams#`.

For example, rather than writing out `http://www.w3.org/ns/ldp#inbox`, I will write out `as:inbox`, which aliases `https://www.w3.org/ns/activitystreams#inbox`, which in turn aliases `ldp:inbox`, and `ldp` aliases `http://www.w3.org/ns/ldp#`.

That said, in JSON form, explicit aliasing is not necessary, because JSON-LD expanders are perfectly capable of resolving the so-called "human-readable" field names perfectly fine, given the appropriate aliases in the contexts.

So while I'd write `as:inbox` in this paragraphs, in JSON, as long as I provide the appropriate context, I'd simply write `inbox`, like so:

```json
{
	"@context": "https://www.w3.org/ns/activitystreams",
	"inbox": "https://sources.example.com/actors/1/inbox"
}
```

## Actor

An actor is represented by a [JSON-LD](https://json-ld.org/) document. An actor does not need to be too complicated.

Here's a barebones actor.

```json
{
	"@context": "https://www.w3.org/ns/activitystreams",
	"id": "https://source.example.com/actors/1",
	"inbox": "https://sources.example.com/actors/1/inbox",
	"outbox": "https://sources.example.com/actors/1/outbox",
	"following": "https://sources.example.com/actors/1/following",
	"followers": "https://sources.example.com/actors/1/followers",
	"liked": "https://sources.example.com/actors/1/liked"
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
	"id": "https://source.example.com/actors/1",
	"type": "Person",
	"inbox": "https://sources.example.com/actors/1/inbox",
	"outbox": "https://sources.example.com/actors/1/outbox",
	"following": "https://sources.example.com/actors/1/following",
	"followers": "https://sources.example.com/actors/1/followers",
	"liked": "https://sources.example.com/actors/1/liked",
	"preferredUsername": "actor",
	"publicKey": {
		"id": "https://source.example.com/actors/1#main-key",
		"owner": "https://source.example.com/actors/1",
		"publicKeyPem": "-----BEGIN PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEArLEIhmSM4UXoUbh/UNri\nOmsruokiG4GU0jz7R/rZ3lC0kGEMEJpk7x8hLEtg0DhV9DW3jPOsPi1KvLRkTgiE\nCSEEG+ULqZ3/WTZR3VX+/Tb1huemD2rBZkv9vpL+3qSRuFTvcMumonVuJ6rtT3pG\nTbsXlYmp2n7VkbPQPz6Wy3R7YeGmdNxtRiccwrpeovc+kCCoY/t467cK1ON+FDrq\nT/xgNhG2jPfotMF3ixk5/EQuakKEz2YQP4duD6D86QciZQWjw5YMv96NxV6D24CV\nn8HxEcxM5AfWvqbNLpEvi6UBUVCnM4IzJTlboPBO4tUPSu01YDqb8jbTC0f6rOCZ\nOQIDAQAB\n-----END PUBLIC KEY-----\n"
	}
}
```

## Following someone

When following someone, you would send an activity in the form of a JSON-LD document, of type `https://www.w3.org/ns/activitystreams#Follow`. It usually looks like this:

```json
{
	"@context": "https://www.w3.org/ns/activitystreams",
	"id": "https://example.com#follows/follow/1",
	"type": "Follow",
	"actor": "https://source.example.com/actors/1",
	"object": "https://destination.example.com/actors/1"
}
```

But, in order to actually send a follow, we must also ensure that the
