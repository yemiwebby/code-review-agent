package comment

import (
	"sync"
	"time"
)

var (
	Mu         = &sync.Mutex{}
	AIComments = map[int]*AIComment{}
)

type AIComment struct {
	ID        int
	Body      string
	File      string
	Timestamp time.Time
	Line      int
	FilePath  string
	OldPatch  string
}

func StoreComment(id int, body, file string, line int, patch string) {
	Mu.Lock()
	defer Mu.Unlock()

	if _, exists := AIComments[id]; !exists {
		AIComments[id] = &AIComment{
			ID:        id,
			Body:      body,
			File:      file,
			Timestamp: time.Now(),
			Line:      line,
			FilePath:  file,
			OldPatch:  patch,
		}
	}
}

func GetComment(id int) (*AIComment, bool) {
	Mu.Lock()
	defer Mu.Unlock()

	comment, exists := AIComments[id]
	return comment, exists
}

func DeleteComment(id int) {
	Mu.Lock()
	defer Mu.Unlock()

	delete(AIComments, id)
}
