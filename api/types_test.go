package api

import (
	"fmt"
	"testing"
)

func TestNextOperationStringer(t *testing.T) {
	if getChannelOverviewOperation := fmt.Sprint(GetChannelOverviewOperation); getChannelOverviewOperation != "GetChannelOverviewOperation" {
		t.Errorf(GetFormattedFailMessage("GetChannelOverviewOperation", "GetChannelOverviewOperation", getChannelOverviewOperation))
	}
	if getVideoIDsOverviewOperation := fmt.Sprint(GetVideoIDsOverviewOperation); getVideoIDsOverviewOperation != "GetVideoIDsOverviewOperation" {
		t.Errorf(GetFormattedFailMessage("GetVideoIDsOverviewOperation", "GetVideoIDsOverviewOperation", getVideoIDsOverviewOperation))
	}
	if getCommentsOverviewOperation := fmt.Sprint(GetCommentsOverviewOperation); getCommentsOverviewOperation != "GetCommentsOverviewOperation" {
		t.Errorf(GetFormattedFailMessage("GetCommentsOverviewOperation", "GetCommentsOverviewOperation", getCommentsOverviewOperation))
	}
	if getObviouslyRelatedChannelsOverviewOperation := fmt.Sprint(GetObviouslyRelatedChannelsOverviewOperation); getObviouslyRelatedChannelsOverviewOperation != "GetObviouslyRelatedChannelsOverviewOperation" {
		t.Errorf(GetFormattedFailMessage("GetObviouslyRelatedChannelsOverviewOperation", "GetObviouslyRelatedChannelsOverviewOperation", getObviouslyRelatedChannelsOverviewOperation))
	}
	if noOperation := fmt.Sprint(NoOperation); noOperation != "NoOperation" {
		t.Errorf(GetFormattedFailMessage("NoOperation", "NoOperation", noOperation))
	}
	var unknownNegativeOperation NextOperation = -1
	if negativeUnknownString := fmt.Sprint(unknownNegativeOperation); negativeUnknownString != "Unknown" {
		t.Errorf(GetFormattedFailMessage("UnknownNegativeOperation", "Unknow", negativeUnknownString))
	}
	var positiveUnknownOperation NextOperation = 5
	if positiveUnknownString := fmt.Sprint(positiveUnknownOperation); positiveUnknownString != "Unknown" {
		t.Errorf(GetFormattedFailMessage("PositiveUnknownOperation", "Unknow", positiveUnknownString))
	}
}
