package tui

import (
	"fmt"
)

func TypeCastToListModel(longList DZList, monthlyList DZList) (ll *listModel, ml *listModel, err error) {
	var ok bool
	err = fmt.Errorf("Cannot type cast DZList")

	ll, ok = longList.(*listModel)
	if !ok {
		return nil, nil, err
	}

	ml, ok = monthlyList.(*listModel)
	if !ok {
		return nil, nil, err
	}

	return ll, ml, nil
}

func GenerateMonthlyText(ml *listModel) string {
	total, completed := ml.countTotalAndCompletedTasks()
	return fmt.Sprintf("Monthly tasks (%d/%d completed)", completed, total)
}

// ok
func RenderPlaceholderTextShowingRemainingTasksToRender(w int, len_monthly int) (r_tasks string) {
	if len_monthly > int(AT_LEAST_NUMBER_OF_MONTHLY_TASKS) {
		r_tasks = remainingTasks.Width((w - 2) / 2).Render(
			fmt.Sprintf("Show %d more tasks", len_monthly-int(AT_LEAST_NUMBER_OF_MONTHLY_TASKS)))
	}

	return
}
