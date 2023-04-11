// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.17.2

package db

import (
	"context"
	"database/sql"
	"fmt"
)

type DBTX interface {
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
	PrepareContext(context.Context, string) (*sql.Stmt, error)
	QueryContext(context.Context, string, ...interface{}) (*sql.Rows, error)
	QueryRowContext(context.Context, string, ...interface{}) *sql.Row
}

func New(db DBTX) *Queries {
	return &Queries{db: db}
}

func Prepare(ctx context.Context, db DBTX) (*Queries, error) {
	q := Queries{db: db}
	var err error
	if q.countUserEmailJobsStmt, err = db.PrepareContext(ctx, countUserEmailJobs); err != nil {
		return nil, fmt.Errorf("error preparing query CountUserEmailJobs: %w", err)
	}
	if q.deleteCandidateJobInterestConditionallyStmt, err = db.PrepareContext(ctx, deleteCandidateJobInterestConditionally); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteCandidateJobInterestConditionally: %w", err)
	}
	if q.deleteUserEmailJobByEmailThreadIDStmt, err = db.PrepareContext(ctx, deleteUserEmailJobByEmailThreadID); err != nil {
		return nil, fmt.Errorf("error preparing query DeleteUserEmailJobByEmailThreadID: %w", err)
	}
	if q.getRecruiterByEmailStmt, err = db.PrepareContext(ctx, getRecruiterByEmail); err != nil {
		return nil, fmt.Errorf("error preparing query GetRecruiterByEmail: %w", err)
	}
	if q.getRecruiterOutboundMessageStmt, err = db.PrepareContext(ctx, getRecruiterOutboundMessage); err != nil {
		return nil, fmt.Errorf("error preparing query GetRecruiterOutboundMessage: %w", err)
	}
	if q.getRecruiterOutboundMessageByRecipientStmt, err = db.PrepareContext(ctx, getRecruiterOutboundMessageByRecipient); err != nil {
		return nil, fmt.Errorf("error preparing query GetRecruiterOutboundMessageByRecipient: %w", err)
	}
	if q.getRecruiterOutboundTemplateStmt, err = db.PrepareContext(ctx, getRecruiterOutboundTemplate); err != nil {
		return nil, fmt.Errorf("error preparing query GetRecruiterOutboundTemplate: %w", err)
	}
	if q.getUserEmailJobStmt, err = db.PrepareContext(ctx, getUserEmailJob); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserEmailJob: %w", err)
	}
	if q.getUserEmailJobByThreadIDStmt, err = db.PrepareContext(ctx, getUserEmailJobByThreadID); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserEmailJobByThreadID: %w", err)
	}
	if q.getUserEmailSyncHistoryStmt, err = db.PrepareContext(ctx, getUserEmailSyncHistory); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserEmailSyncHistory: %w", err)
	}
	if q.getUserOAuthTokenStmt, err = db.PrepareContext(ctx, getUserOAuthToken); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserOAuthToken: %w", err)
	}
	if q.getUserProfileByEmailStmt, err = db.PrepareContext(ctx, getUserProfileByEmail); err != nil {
		return nil, fmt.Errorf("error preparing query GetUserProfileByEmail: %w", err)
	}
	if q.incrementUserEmailStatStmt, err = db.PrepareContext(ctx, incrementUserEmailStat); err != nil {
		return nil, fmt.Errorf("error preparing query IncrementUserEmailStat: %w", err)
	}
	if q.insertRecruiterOutboundMessageStmt, err = db.PrepareContext(ctx, insertRecruiterOutboundMessage); err != nil {
		return nil, fmt.Errorf("error preparing query InsertRecruiterOutboundMessage: %w", err)
	}
	if q.insertRecruiterOutboundTemplateStmt, err = db.PrepareContext(ctx, insertRecruiterOutboundTemplate); err != nil {
		return nil, fmt.Errorf("error preparing query InsertRecruiterOutboundTemplate: %w", err)
	}
	if q.insertUserEmailJobStmt, err = db.PrepareContext(ctx, insertUserEmailJob); err != nil {
		return nil, fmt.Errorf("error preparing query InsertUserEmailJob: %w", err)
	}
	if q.listCandidateOAuthTokensStmt, err = db.PrepareContext(ctx, listCandidateOAuthTokens); err != nil {
		return nil, fmt.Errorf("error preparing query ListCandidateOAuthTokens: %w", err)
	}
	if q.listRecruiterOAuthTokensStmt, err = db.PrepareContext(ctx, listRecruiterOAuthTokens); err != nil {
		return nil, fmt.Errorf("error preparing query ListRecruiterOAuthTokens: %w", err)
	}
	if q.listSimilarRecruiterOutboundTemplatesStmt, err = db.PrepareContext(ctx, listSimilarRecruiterOutboundTemplates); err != nil {
		return nil, fmt.Errorf("error preparing query ListSimilarRecruiterOutboundTemplates: %w", err)
	}
	if q.listUserEmailJobsStmt, err = db.PrepareContext(ctx, listUserEmailJobs); err != nil {
		return nil, fmt.Errorf("error preparing query ListUserEmailJobs: %w", err)
	}
	if q.listUserOAuthTokensStmt, err = db.PrepareContext(ctx, listUserOAuthTokens); err != nil {
		return nil, fmt.Errorf("error preparing query ListUserOAuthTokens: %w", err)
	}
	if q.upsertCandidateJobInterestStmt, err = db.PrepareContext(ctx, upsertCandidateJobInterest); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertCandidateJobInterest: %w", err)
	}
	if q.upsertUserEmailSyncHistoryStmt, err = db.PrepareContext(ctx, upsertUserEmailSyncHistory); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertUserEmailSyncHistory: %w", err)
	}
	if q.upsertUserOAuthTokenStmt, err = db.PrepareContext(ctx, upsertUserOAuthToken); err != nil {
		return nil, fmt.Errorf("error preparing query UpsertUserOAuthToken: %w", err)
	}
	return &q, nil
}

func (q *Queries) Close() error {
	var err error
	if q.countUserEmailJobsStmt != nil {
		if cerr := q.countUserEmailJobsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing countUserEmailJobsStmt: %w", cerr)
		}
	}
	if q.deleteCandidateJobInterestConditionallyStmt != nil {
		if cerr := q.deleteCandidateJobInterestConditionallyStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteCandidateJobInterestConditionallyStmt: %w", cerr)
		}
	}
	if q.deleteUserEmailJobByEmailThreadIDStmt != nil {
		if cerr := q.deleteUserEmailJobByEmailThreadIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing deleteUserEmailJobByEmailThreadIDStmt: %w", cerr)
		}
	}
	if q.getRecruiterByEmailStmt != nil {
		if cerr := q.getRecruiterByEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRecruiterByEmailStmt: %w", cerr)
		}
	}
	if q.getRecruiterOutboundMessageStmt != nil {
		if cerr := q.getRecruiterOutboundMessageStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRecruiterOutboundMessageStmt: %w", cerr)
		}
	}
	if q.getRecruiterOutboundMessageByRecipientStmt != nil {
		if cerr := q.getRecruiterOutboundMessageByRecipientStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRecruiterOutboundMessageByRecipientStmt: %w", cerr)
		}
	}
	if q.getRecruiterOutboundTemplateStmt != nil {
		if cerr := q.getRecruiterOutboundTemplateStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getRecruiterOutboundTemplateStmt: %w", cerr)
		}
	}
	if q.getUserEmailJobStmt != nil {
		if cerr := q.getUserEmailJobStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserEmailJobStmt: %w", cerr)
		}
	}
	if q.getUserEmailJobByThreadIDStmt != nil {
		if cerr := q.getUserEmailJobByThreadIDStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserEmailJobByThreadIDStmt: %w", cerr)
		}
	}
	if q.getUserEmailSyncHistoryStmt != nil {
		if cerr := q.getUserEmailSyncHistoryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserEmailSyncHistoryStmt: %w", cerr)
		}
	}
	if q.getUserOAuthTokenStmt != nil {
		if cerr := q.getUserOAuthTokenStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserOAuthTokenStmt: %w", cerr)
		}
	}
	if q.getUserProfileByEmailStmt != nil {
		if cerr := q.getUserProfileByEmailStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing getUserProfileByEmailStmt: %w", cerr)
		}
	}
	if q.incrementUserEmailStatStmt != nil {
		if cerr := q.incrementUserEmailStatStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing incrementUserEmailStatStmt: %w", cerr)
		}
	}
	if q.insertRecruiterOutboundMessageStmt != nil {
		if cerr := q.insertRecruiterOutboundMessageStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertRecruiterOutboundMessageStmt: %w", cerr)
		}
	}
	if q.insertRecruiterOutboundTemplateStmt != nil {
		if cerr := q.insertRecruiterOutboundTemplateStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertRecruiterOutboundTemplateStmt: %w", cerr)
		}
	}
	if q.insertUserEmailJobStmt != nil {
		if cerr := q.insertUserEmailJobStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing insertUserEmailJobStmt: %w", cerr)
		}
	}
	if q.listCandidateOAuthTokensStmt != nil {
		if cerr := q.listCandidateOAuthTokensStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listCandidateOAuthTokensStmt: %w", cerr)
		}
	}
	if q.listRecruiterOAuthTokensStmt != nil {
		if cerr := q.listRecruiterOAuthTokensStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listRecruiterOAuthTokensStmt: %w", cerr)
		}
	}
	if q.listSimilarRecruiterOutboundTemplatesStmt != nil {
		if cerr := q.listSimilarRecruiterOutboundTemplatesStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listSimilarRecruiterOutboundTemplatesStmt: %w", cerr)
		}
	}
	if q.listUserEmailJobsStmt != nil {
		if cerr := q.listUserEmailJobsStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listUserEmailJobsStmt: %w", cerr)
		}
	}
	if q.listUserOAuthTokensStmt != nil {
		if cerr := q.listUserOAuthTokensStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing listUserOAuthTokensStmt: %w", cerr)
		}
	}
	if q.upsertCandidateJobInterestStmt != nil {
		if cerr := q.upsertCandidateJobInterestStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertCandidateJobInterestStmt: %w", cerr)
		}
	}
	if q.upsertUserEmailSyncHistoryStmt != nil {
		if cerr := q.upsertUserEmailSyncHistoryStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertUserEmailSyncHistoryStmt: %w", cerr)
		}
	}
	if q.upsertUserOAuthTokenStmt != nil {
		if cerr := q.upsertUserOAuthTokenStmt.Close(); cerr != nil {
			err = fmt.Errorf("error closing upsertUserOAuthTokenStmt: %w", cerr)
		}
	}
	return err
}

func (q *Queries) exec(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (sql.Result, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).ExecContext(ctx, args...)
	case stmt != nil:
		return stmt.ExecContext(ctx, args...)
	default:
		return q.db.ExecContext(ctx, query, args...)
	}
}

func (q *Queries) query(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) (*sql.Rows, error) {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryContext(ctx, args...)
	default:
		return q.db.QueryContext(ctx, query, args...)
	}
}

func (q *Queries) queryRow(ctx context.Context, stmt *sql.Stmt, query string, args ...interface{}) *sql.Row {
	switch {
	case stmt != nil && q.tx != nil:
		return q.tx.StmtContext(ctx, stmt).QueryRowContext(ctx, args...)
	case stmt != nil:
		return stmt.QueryRowContext(ctx, args...)
	default:
		return q.db.QueryRowContext(ctx, query, args...)
	}
}

type Queries struct {
	db                                          DBTX
	tx                                          *sql.Tx
	countUserEmailJobsStmt                      *sql.Stmt
	deleteCandidateJobInterestConditionallyStmt *sql.Stmt
	deleteUserEmailJobByEmailThreadIDStmt       *sql.Stmt
	getRecruiterByEmailStmt                     *sql.Stmt
	getRecruiterOutboundMessageStmt             *sql.Stmt
	getRecruiterOutboundMessageByRecipientStmt  *sql.Stmt
	getRecruiterOutboundTemplateStmt            *sql.Stmt
	getUserEmailJobStmt                         *sql.Stmt
	getUserEmailJobByThreadIDStmt               *sql.Stmt
	getUserEmailSyncHistoryStmt                 *sql.Stmt
	getUserOAuthTokenStmt                       *sql.Stmt
	getUserProfileByEmailStmt                   *sql.Stmt
	incrementUserEmailStatStmt                  *sql.Stmt
	insertRecruiterOutboundMessageStmt          *sql.Stmt
	insertRecruiterOutboundTemplateStmt         *sql.Stmt
	insertUserEmailJobStmt                      *sql.Stmt
	listCandidateOAuthTokensStmt                *sql.Stmt
	listRecruiterOAuthTokensStmt                *sql.Stmt
	listSimilarRecruiterOutboundTemplatesStmt   *sql.Stmt
	listUserEmailJobsStmt                       *sql.Stmt
	listUserOAuthTokensStmt                     *sql.Stmt
	upsertCandidateJobInterestStmt              *sql.Stmt
	upsertUserEmailSyncHistoryStmt              *sql.Stmt
	upsertUserOAuthTokenStmt                    *sql.Stmt
}

func (q *Queries) WithTx(tx *sql.Tx) *Queries {
	return &Queries{
		db:                     tx,
		tx:                     tx,
		countUserEmailJobsStmt: q.countUserEmailJobsStmt,
		deleteCandidateJobInterestConditionallyStmt: q.deleteCandidateJobInterestConditionallyStmt,
		deleteUserEmailJobByEmailThreadIDStmt:       q.deleteUserEmailJobByEmailThreadIDStmt,
		getRecruiterByEmailStmt:                     q.getRecruiterByEmailStmt,
		getRecruiterOutboundMessageStmt:             q.getRecruiterOutboundMessageStmt,
		getRecruiterOutboundMessageByRecipientStmt:  q.getRecruiterOutboundMessageByRecipientStmt,
		getRecruiterOutboundTemplateStmt:            q.getRecruiterOutboundTemplateStmt,
		getUserEmailJobStmt:                         q.getUserEmailJobStmt,
		getUserEmailJobByThreadIDStmt:               q.getUserEmailJobByThreadIDStmt,
		getUserEmailSyncHistoryStmt:                 q.getUserEmailSyncHistoryStmt,
		getUserOAuthTokenStmt:                       q.getUserOAuthTokenStmt,
		getUserProfileByEmailStmt:                   q.getUserProfileByEmailStmt,
		incrementUserEmailStatStmt:                  q.incrementUserEmailStatStmt,
		insertRecruiterOutboundMessageStmt:          q.insertRecruiterOutboundMessageStmt,
		insertRecruiterOutboundTemplateStmt:         q.insertRecruiterOutboundTemplateStmt,
		insertUserEmailJobStmt:                      q.insertUserEmailJobStmt,
		listCandidateOAuthTokensStmt:                q.listCandidateOAuthTokensStmt,
		listRecruiterOAuthTokensStmt:                q.listRecruiterOAuthTokensStmt,
		listSimilarRecruiterOutboundTemplatesStmt:   q.listSimilarRecruiterOutboundTemplatesStmt,
		listUserEmailJobsStmt:                       q.listUserEmailJobsStmt,
		listUserOAuthTokensStmt:                     q.listUserOAuthTokensStmt,
		upsertCandidateJobInterestStmt:              q.upsertCandidateJobInterestStmt,
		upsertUserEmailSyncHistoryStmt:              q.upsertUserEmailSyncHistoryStmt,
		upsertUserOAuthTokenStmt:                    q.upsertUserOAuthTokenStmt,
	}
}
