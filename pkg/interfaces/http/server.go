package http

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/tclaudel/golang_blockchain/config"
	"github.com/tclaudel/golang_blockchain/internal/entity"
	"github.com/tclaudel/golang_blockchain/pkg/repositories"
	"go.uber.org/zap"
)

type Server struct {
	*http.Server
	cfg            *config.Config
	logger         *zap.Logger
	blockchainNode *entity.BlockchainNode
	repositories.Repositories
}

func NewServer(cfg *config.Config, logger *zap.Logger, blockChainNode *entity.BlockchainNode, repositories repositories.Repositories) Server {
	server := Server{
		&http.Server{
			Addr:         cfg.Http.Address,
			ReadTimeout:  time.Duration(cfg.Http.ReadTimeout) * time.Second,
			WriteTimeout: time.Duration(cfg.Http.WriteTimeout) * time.Second,
		},
		cfg,
		logger,
		blockChainNode,
		repositories,
	}

	r := mux.NewRouter()
	r.Use(NewMiddlewareRequestID(logger))
	r.Use(NewMiddlewareRequestLogger(logger))
	r.HandleFunc("/blockchain", server.GetBlockchain).Methods(http.MethodGet)

	server.Handler = r

	return server
}

func (s *Server) Start(cfg *config.Config, logger *zap.Logger, errChan chan error) {
	if s.Addr == "" {
		errChan <- errors.New("server address is missing in configuration")
		return
	}

	if !strings.HasPrefix(cfg.Http.Address, ":") {
		s.Addr = ":" + cfg.Http.Address
	}

	logger.Info("Starting http server " + s.Addr)
	if err := s.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
		errChan <- err
		return
	}
	logger.Info("Application stopped")
}

func (s *Server) Shutdown(ctx context.Context, logger *zap.Logger) {
	if err := s.Server.Shutdown(ctx); err != nil {
		logger.Fatal(err.Error())
		return
	}

	logger.Info("Application shutdown")
}
