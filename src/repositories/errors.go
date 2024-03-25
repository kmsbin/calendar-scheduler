package repositories

import "errors"

var UserNotFounded = errors.New("user not founded")
var MeetingsRangeNotFounded = errors.New("meetings range not founded")

var TokenNotFounded = errors.New("token not founded")

var EmailDuplicatedInMeetingRange = errors.New("an email has already been registered for this event")
