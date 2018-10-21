package youtube

import (
	"fmt"
	"testing"

	"gitlab.com/TheDonDope/gocha/v3/pkg/errors"
)

func TestNextOperationStringer(t *testing.T) {
	if getChannelOverviewOperation := fmt.Sprint(GetChannelOverviewOperation); getChannelOverviewOperation != "GetChannelOverviewOperation" {
		t.Errorf(errors.Fail("GetChannelOverviewOperation", "GetChannelOverviewOperation", getChannelOverviewOperation))
	}
	if getVideoIDsOverviewOperation := fmt.Sprint(GetVideoIDsOverviewOperation); getVideoIDsOverviewOperation != "GetVideoIDsOverviewOperation" {
		t.Errorf(errors.Fail("GetVideoIDsOverviewOperation", "GetVideoIDsOverviewOperation", getVideoIDsOverviewOperation))
	}
	if getCommentsOverviewOperation := fmt.Sprint(GetCommentsOverviewOperation); getCommentsOverviewOperation != "GetCommentsOverviewOperation" {
		t.Errorf(errors.Fail("GetCommentsOverviewOperation", "GetCommentsOverviewOperation", getCommentsOverviewOperation))
	}
	if getObviouslyRelatedChannelsOverviewOperation := fmt.Sprint(GetObviouslyRelatedChannelsOverviewOperation); getObviouslyRelatedChannelsOverviewOperation != "GetObviouslyRelatedChannelsOverviewOperation" {
		t.Errorf(errors.Fail("GetObviouslyRelatedChannelsOverviewOperation", "GetObviouslyRelatedChannelsOverviewOperation", getObviouslyRelatedChannelsOverviewOperation))
	}
	if noOperation := fmt.Sprint(NoOperation); noOperation != "NoOperation" {
		t.Errorf(errors.Fail("NoOperation", "NoOperation", noOperation))
	}
	var unknownNegativeOperation NextOperation = -1
	if negativeUnknownString := fmt.Sprint(unknownNegativeOperation); negativeUnknownString != "Unknown" {
		t.Errorf(errors.Fail("UnknownNegativeOperation", "Unknow", negativeUnknownString))
	}
	var positiveUnknownOperation NextOperation = 5
	if positiveUnknownString := fmt.Sprint(positiveUnknownOperation); positiveUnknownString != "Unknown" {
		t.Errorf(errors.Fail("PositiveUnknownOperation", "Unknow", positiveUnknownString))
	}
}
