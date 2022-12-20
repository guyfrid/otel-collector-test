package split_processor

import (
	"context"
	"fmt"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/consumer/pdata"
	"splitbatchprocessor/split"
)

var _ component.TraceProcessor = (*splitBatch)(nil)

type splitBatch struct {
	next consumer.TraceConsumer
}

func NewSplitBatch(next consumer.TraceConsumer) *splitBatch {
	return &splitBatch{next: next}
}

func (s *splitBatch) ConsumeTraces(ctx context.Context, batch pdata.Traces) error {
	for i := 0; i < batch.ResourceSpans().Len(); i++ {
		// we split the incoming batch into a collection of ResourceSpans
		rss := split.SplitByTrace(batch.ResourceSpans().At(i))
		for _, newBatch := range rss {
			trace := pdata.NewTraces()
			trace.ResourceSpans().Append(newBatch)
			fmt.Println("####### I'M CN-CONTRAST PROCESSOR ##########")
			if err := s.next.ConsumeTraces(ctx, trace); err != nil {
				// we fail fast: if we get an error from the next, we break the processing for this batch
				return err
			}
		}
	}

	return nil
}

func (s *splitBatch) GetCapabilities() component.ProcessorCapabilities {
	return component.ProcessorCapabilities{MutatesConsumedData: true}
}

func (s *splitBatch) Start(_ context.Context, host component.Host) error {
	return nil
}

func (s *splitBatch) Shutdown(context.Context) error {
	return nil
}
