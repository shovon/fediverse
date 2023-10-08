# Fediverse Utils

Just a bunch of stuff related to the Fediverse. Not just ActivityPub, but the individual things that make it up, among other things that the Fediverse relies on, such as WebFinger.

This is all just one big WIP.

Stay tuned.

Development (requires [Air](https://github.com/cosmtrek/air)):

```
env -S "`cat ./.example.env`" air .
```

Following (planned to be deprecated):

```
env -S "`cat ./.example.env`" go r
```

Create a post (planned to be deprecated):

```
env -S "`cat ./.example.env`" go run ./application/cli* create "$POST"