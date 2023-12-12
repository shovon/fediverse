# A primer on how to work with the Fediverse

## ActivityPub

The most important component of the Fediverse (beyond the Internet and HTTP) is ActivityPub.

ActivityPub's primitive is the "actor", and in any ActivityPub-based Fediverse software, a user is associated with an actor, and it's that actor that users will interact with, via their own actors, by sending each other activities, via HTTP.

ActivityPub's primary data interchange format is JSON-LD. So an actor is encoded as a JSON-LD document, and so are the _activities_ that actors send to other actors.

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
