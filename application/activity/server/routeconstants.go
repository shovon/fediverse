package server

const activityRoute string = "/activity"
const sharedInbox string = "/sharedInbox"
const userRoute string = activityRoute + "/actors/:username"
const followingRoute string = userRoute + "/following"
const followersRoute string = userRoute + "/followers"
const inboxRoute string = userRoute + "/inbox"
const outboxRoute string = userRoute + "/outbox"
const likedRoute string = userRoute + "/liked"
