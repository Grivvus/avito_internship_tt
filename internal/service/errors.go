package service

import "errors"

var TeamAlreadyExistError error = errors.New("already exists")
var ResourceNotFoundError error = errors.New("resource not found")
var PRAlreadyExistError = errors.New("already exist")
var CantReassignOnMergedPRError = errors.New("cannot reassign on merged PR")
var ReviewerNotAssignedError = errors.New("reviewer is not assigned to this PR")
var NoCandidatesError = errors.New("no active replacement candidate in team")
