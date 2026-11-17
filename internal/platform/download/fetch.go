package download

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
)

func fetch(ctx context.Context, url, partial string, progress Progress) (bool, error) {
	offset := int64(0)
	if info, err := os.Stat(partial); err == nil {
		offset = info.Size()
	}
	activity := make(chan struct{}, 1)
	attempt, stop := watchIdle(ctx, activity)
	defer stop()
	request, err := http.NewRequestWithContext(attempt, http.MethodGet, url, nil)
	if err != nil {
		return false, err
	}
	if offset > 0 {
		request.Header.Set("Range", fmt.Sprintf("bytes=%d-", offset))
	}
	response, err := client.Do(request)
	if err != nil {
		if ctx.Err() == nil && attempt.Err() != nil {
			return true, errIdle
		}
		return retryable(err), err
	}
	defer func() { _ = response.Body.Close() }()
	switch response.StatusCode {
	case http.StatusPartialContent:
	case http.StatusOK:
		offset = 0
	case http.StatusRequestedRangeNotSatisfiable:
		_ = os.Remove(partial)
		return true, fmt.Errorf("download: range rejected for %s", url)
	default:
		return retryStatus(response.StatusCode),
			fmt.Errorf("download status %s", response.Status)
	}
	file, err := os.OpenFile(partial, os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return false, err
	}
	if err := file.Truncate(offset); err != nil {
		_ = file.Close()
		return false, err
	}
	if _, err := file.Seek(offset, io.SeekStart); err != nil {
		_ = file.Close()
		return false, err
	}
	total := int64(0)
	if response.ContentLength >= 0 {
		total = offset + response.ContentLength
	}
	if total > maxFileSize {
		_ = file.Close()
		_ = os.Remove(partial)
		return false, fmt.Errorf("download: response exceeds the size limit")
	}
	reader := &tracked{reader: response.Body, done: offset, total: total,
		step: progressStep, report: progress, activity: activity}
	_, copyErr := io.Copy(file, io.LimitReader(reader, maxFileSize+1))
	closeErr := file.Close()
	if copyErr != nil {
		if ctx.Err() == nil && attempt.Err() != nil {
			return true, errIdle
		}
		return retryable(copyErr), copyErr
	}
	if reader.done > maxFileSize {
		_ = os.Remove(partial)
		return false, fmt.Errorf("download: response exceeds the size limit")
	}
	if closeErr != nil {
		return false, closeErr
	}
	if progress != nil {
		progress(reader.done, total)
	}
	return false, nil
}
