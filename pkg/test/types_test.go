package test

import (
	"fmt"
	"testing"

	"gitlab.com/TheDonDope/gocha/pkg/util/errors"
	"gitlab.com/TheDonDope/youtube-frequenter/pkg/types"
)

func TestNextOperationStringer(t *testing.T) {
	if getChannelOverviewOperation := fmt.Sprint(types.GetChannelOverviewOperation); getChannelOverviewOperation != "GetChannelOverviewOperation" {
		t.Errorf(errors.GetFormattedFailMessage("GetChannelOverviewOperation", "GetChannelOverviewOperation", getChannelOverviewOperation))
	}
	if getVideoIDsOverviewOperation := fmt.Sprint(types.GetVideoIDsOverviewOperation); getVideoIDsOverviewOperation != "GetVideoIDsOverviewOperation" {
		t.Errorf(errors.GetFormattedFailMessage("GetVideoIDsOverviewOperation", "GetVideoIDsOverviewOperation", getVideoIDsOverviewOperation))
	}
	if getCommentsOverviewOperation := fmt.Sprint(types.GetCommentsOverviewOperation); getCommentsOverviewOperation != "GetCommentsOverviewOperation" {
		t.Errorf(errors.GetFormattedFailMessage("GetCommentsOverviewOperation", "GetCommentsOverviewOperation", getCommentsOverviewOperation))
	}
	if getObviouslyRelatedChannelsOverviewOperation := fmt.Sprint(types.GetObviouslyRelatedChannelsOverviewOperation); getObviouslyRelatedChannelsOverviewOperation != "GetObviouslyRelatedChannelsOverviewOperation" {
		t.Errorf(errors.GetFormattedFailMessage("GetObviouslyRelatedChannelsOverviewOperation", "GetObviouslyRelatedChannelsOverviewOperation", getObviouslyRelatedChannelsOverviewOperation))
	}
	if noOperation := fmt.Sprint(types.NoOperation); noOperation != "NoOperation" {
		t.Errorf(errors.GetFormattedFailMessage("NoOperation", "NoOperation", noOperation))
	}
	var unknownNegativeOperation types.NextOperation = -1
	if negativeUnknownString := fmt.Sprint(unknownNegativeOperation); negativeUnknownString != "Unknown" {
		t.Errorf(errors.GetFormattedFailMessage("UnknownNegativeOperation", "Unknow", negativeUnknownString))
	}
	var positiveUnknownOperation types.NextOperation = 5
	if positiveUnknownString := fmt.Sprint(positiveUnknownOperation); positiveUnknownString != "Unknown" {
		t.Errorf(errors.GetFormattedFailMessage("PositiveUnknownOperation", "Unknow", positiveUnknownString))
	}
}
