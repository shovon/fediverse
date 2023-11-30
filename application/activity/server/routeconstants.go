package server

const ActivityRoute string = "/activity"
const SharedInbox string = "/sharedInbox"
const UserRoute string = ActivityRoute + "/actors/:username"
const GlobalFollowingRoute string = ActivityRoute + "/following/:following"
const FollowingRoute string = UserRoute + "/following"
const FollowersRoute string = UserRoute + "/followers"
const InboxRoute string = UserRoute + "/inbox"
const OutboxRoute string = UserRoute + "/outbox"
const LikedRoute string = UserRoute + "/liked"
