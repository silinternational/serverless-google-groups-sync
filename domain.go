package domain

// GroupDiff holds the information needed by the sync process in relation to
//   the source group name,
//   the target group name (e.g. the google group),
//   the email addresses of the members of the source group
//   the email address of the current members of the target group
//   the email address of the source group that should be added to the target group
//   the email address that should be delete from the target group, since they
//     don't appear in the source group
type GroupDiff struct {
	SourceGroup string
	TargetGroup string
	SourceMembers []string
	TargetMembers []string
	MembersToAdd []string
	MembersToDelete []string
}