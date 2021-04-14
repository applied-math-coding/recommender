package services

import (
	"bufio"
	"fmt"
	"main/domain"
	"math"
	"mime/multipart"
	"sort"
	"strings"

	"github.com/pkg/errors"
)

// Data globally references current training data.
var ItemSets []domain.ItemSet

func ProcessFileData(fileHeader *multipart.FileHeader, broadcaster *Broadcaster) {
	defer close(broadcaster.ProgressChan)
	EmitProgress("Opening file.", 0.0, domain.ProgressState.Running, false, broadcaster.ProgressChan)
	file, e2 := fileHeader.Open()
	if e2 != nil {
		HandleError(errors.Wrap(e2, "fileHeader.Open fails"), nil, false)
		EmitProgress("Error when reading file.", 0.0, domain.ProgressState.Error, false, broadcaster.ProgressChan)
		return
	}
	var canceled bool
	ItemSets, canceled = ReadData(file, broadcaster)
	if canceled {
		EmitProgress("Canceled.", 1.0, domain.ProgressState.Canceled, false, broadcaster.ProgressChan)
	} else {
		EmitProgress("Finished.", 0.0, domain.ProgressState.Finished, false, broadcaster.ProgressChan)
	}
}

func ReadData(file multipart.File, broadcaster *Broadcaster) ([]domain.ItemSet, domain.Canceled) {
	scanner := bufio.NewScanner(file)
	res := make([]domain.ItemSet, 0)
	recordCount := 0
	for scanner.Scan() {
		if broadcaster.Cancel {
			return nil, true
		}
		record := strings.Split(scanner.Text(), ",")
		sort.Strings(record)
		res = append(res, record)
		recordCount++
		if math.Mod(float64(recordCount), 100000) == 0 {
			EmitProgress(
				fmt.Sprintf("Reading records: %d", recordCount),
				0.0,
				domain.ProgressState.Running,
				false,
				broadcaster.ProgressChan)
		}
	}
	return res, false
}
