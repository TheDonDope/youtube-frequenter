package types

// NextOperation declares an enum type for operation names
type NextOperation int

const (
	// GetChannelOverviewOperation is a NextOperation enum value
	GetChannelOverviewOperation NextOperation = iota
	// GetVideoIDsOverviewOperation is a NextOperation enum value
	GetVideoIDsOverviewOperation
	// GetCommentsOverviewOperation is a NextOperation enum value
	GetCommentsOverviewOperation
	// GetObviouslyRelatedChannelsOverviewOperation is a NextOperation enum value
	GetObviouslyRelatedChannelsOverviewOperation
	// NoOperation is a NextOperation enum value
	NoOperation
)

// This function is the Stringer for the NextOperation enum
func (operation NextOperation) String() string {
	operationNames := [...]string{
		"GetChannelOverviewOperation",
		"GetVideoIDsOverviewOperation",
		"GetCommentsOverviewOperation",
		"GetObviouslyRelatedChannelsOverviewOperation",
		"NoOperation"}
	if operation < GetChannelOverviewOperation || operation > NoOperation {
		return "Unknown"
	}
	return operationNames[operation]
}
