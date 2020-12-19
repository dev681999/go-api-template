package db

import (
	"context"
	"fmt"
	"go-api-template/internal/config"

	"github.com/go-pg/pg/v10"
	"github.com/rs/zerolog"
)

// debugHook is a query hook that logs an error with a query if there are any.
// It can be installed with:
//
//   db.AddQueryHook(pgext.debugHook{})
type debugHook struct {
	// Verbose causes hook to print all queries (even those without an error).
	Verbose   bool
	EmptyLine bool
	Logger    zerolog.Logger
}

var _ pg.QueryHook = (*debugHook)(nil)

func (h debugHook) BeforeQuery(ctx context.Context, evt *pg.QueryEvent) (context.Context, error) {
	q, err := evt.FormattedQuery()
	if err != nil {
		return nil, err
	}

	if evt.Err != nil {
		h.Logger.Debug().Err(fmt.Errorf("%s executing a query:\n%s", evt.Err, q)).Msg("")
	} else if h.Verbose {
		h.Logger.Debug().Msgf("q -> %v", string(q))
	}

	return ctx, nil
}

func (debugHook) AfterQuery(context.Context, *pg.QueryEvent) error {
	return nil
}

// Connect connects to db and creates necessary tables, index
func Connect(ctx context.Context, logger zerolog.Logger, config *config.Config) (*pg.DB, error) {
	db := pg.Connect(&pg.Options{
		Addr:     config.DB.Host,
		User:     config.DB.User,
		Password: config.DB.Password,
		Database: config.DB.Database,
	})

	if err := db.Ping(ctx); err != nil {
		err = fmt.Errorf("db is down err: %v", err)
		logger.Debug().Err(err).Msg("")
		return nil, err
	}

	db.AddQueryHook(debugHook{
		Verbose: true,
		Logger:  logger,
	})

	return db, nil
}
