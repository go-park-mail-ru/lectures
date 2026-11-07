// Package platform содержит внешние механизмы, выбранные при сборке приложения.
package platform

import (
	"lectures/cleanarchreview/internal/usecase/article"
	"sync/atomic"
	"time"
)

type Clock struct{}

func (Clock) Now() time.Time { return time.Now().UTC() }

var _ article.Clock = Clock{}

type IDs struct{ last atomic.Int64 }

func (g *IDs) Next() int64 { return g.last.Add(1) }

var _ article.IDs = (*IDs)(nil)
