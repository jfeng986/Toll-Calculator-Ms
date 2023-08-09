package main

import (
	"context"

	"Toll-Calculator/types"
)

type GPRCAggregatorServer struct {
	types.UnimplementedAggregatorServer
	svc Aggregator
}

func NewAggregatorGRPCServer(svc Aggregator) *GPRCAggregatorServer {
	return &GPRCAggregatorServer{
		svc: svc,
	}
}

func (s *GPRCAggregatorServer) Aggregate(ctx context.Context, req *types.AggregateRequest) (*types.None, error) {
	distance := types.Distance{
		ObuID: int(req.ObuID),
		Value: req.Value,
		Unix:  req.Unix,
	}
	return &types.None{}, s.svc.AggregateDistance(distance)
}
