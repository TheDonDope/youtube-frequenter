package api

import (
	"fmt"
	"testing"
)

func TestNextOperationStringer(t *testing.T) {
	if getChannelOverviewOperation := fmt.Sprint(GetChannelOverviewOperation); getChannelOverviewOperation != "GetChannelOverviewOperation" {
		t.Errorf("GetChannelOverviewOperation stringer was incorrect, got: %s, want: %s", getChannelOverviewOperation, "GetChannelOverviewOperation")
	}
	if getVideoIDsOverviewOperation := fmt.Sprint(GetVideoIDsOverviewOperation); getVideoIDsOverviewOperation != "GetVideoIDsOverviewOperation" {
		t.Errorf("GetVideoIDsOverviewOperation stringer was incorrect, got: %s, want: %s", getVideoIDsOverviewOperation, "GetVideoIDsOverviewOperation")
	}
	if getCommentsOverviewOperation := fmt.Sprint(GetCommentsOverviewOperation); getCommentsOverviewOperation != "GetCommentsOverviewOperation" {
		t.Errorf("GetCommentsOverviewOperation stringer was incorrect, got: %s, want: %s", getCommentsOverviewOperation, "GetCommentsOverviewOperation")
	}
	if getObviouslyRelatedChannelsOverviewOperation := fmt.Sprint(GetObviouslyRelatedChannelsOverviewOperation); getObviouslyRelatedChannelsOverviewOperation != "GetObviouslyRelatedChannelsOverviewOperation" {
		t.Errorf("GetObviouslyRelatedChannelsOverviewOperation stringer was incorrect, got: %s, want: %s", getObviouslyRelatedChannelsOverviewOperation, "GetObviouslyRelatedChannelsOverviewOperation")
	}
	if noOperation := fmt.Sprint(NoOperation); noOperation != "NoOperation" {
		t.Errorf("NoOperation stringer was incorrect, got: %s, want: %s", noOperation, "NoOperation")
	}
	var unknownNegativeOperation NextOperation = -1
	if negativeUnknownString := fmt.Sprint(unknownNegativeOperation); negativeUnknownString != "Unknown" {
		t.Errorf("unknownNegativeOperation stringer was incorrect, got: %s, want: %s", negativeUnknownString, "Unknown")
	}
	var positiveUnknownOperation NextOperation = 4
	if positiveUnknownString := fmt.Sprint(positiveUnknownOperation); positiveUnknownString != "Unknown" {
		t.Errorf("positiveUnknownOperation stringer was incorrect, got: %s, want: %s", positiveUnknownString, "Unknown")
	}
}
