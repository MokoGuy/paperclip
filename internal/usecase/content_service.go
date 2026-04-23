package usecase

import (
	"context"
	"sync"

	apirepo "github.com/MokoGuy/paperclip/internal/repository/api"
)

type ContentResult struct {
	ID      int
	Title   string
	Content string
	Err     error
}

type ContentService struct {
	apiDoc *apirepo.DocumentRepository
}

func NewContentService(apiDoc *apirepo.DocumentRepository) *ContentService {
	return &ContentService{apiDoc: apiDoc}
}

func (s *ContentService) GetContents(ctx context.Context, ids []int) []ContentResult {
	results := make([]ContentResult, len(ids))
	sem := make(chan struct{}, 5)
	var wg sync.WaitGroup

	for i, id := range ids {
		wg.Add(1)
		go func(idx, docID int) {
			defer wg.Done()
			sem <- struct{}{}
			defer func() { <-sem }()

			detail, err := s.apiDoc.GetContent(ctx, docID)
			if err != nil {
				results[idx] = ContentResult{ID: docID, Err: err}
			} else {
				results[idx] = ContentResult{
					ID:      docID,
					Title:   detail.Title,
					Content: detail.Content,
				}
			}
		}(i, id)
	}

	wg.Wait()
	return results
}
