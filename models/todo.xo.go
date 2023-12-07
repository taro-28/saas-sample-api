package models

// Code generated by xo. DO NOT EDIT.

import (
	"context"
)

// Todo represents a row from 'saas-sample.todos'.
type Todo struct {
	ID        int    `json:"id"`        // id
	Content   string `json:"content"`   // content
	Completed bool   `json:"completed"` // completed
	// xo fields
	_exists, _deleted bool
}

// Exists returns true when the [Todo] exists in the database.
func (t *Todo) Exists() bool {
	return t._exists
}

// Deleted returns true when the [Todo] has been marked for deletion
// from the database.
func (t *Todo) Deleted() bool {
	return t._deleted
}

// Insert inserts the [Todo] to the database.
func (t *Todo) Insert(ctx context.Context, db DB) error {
	switch {
	case t._exists: // already exists
		return logerror(&ErrInsertFailed{ErrAlreadyExists})
	case t._deleted: // deleted
		return logerror(&ErrInsertFailed{ErrMarkedForDeletion})
	}
	// insert (manual)
	const sqlstr = `INSERT INTO saas-sample.todos (` +
		`id, content, completed` +
		`) VALUES (` +
		`?, ?, ?` +
		`)`
	// run
	logf(sqlstr, t.ID, t.Content, t.Completed)
	if _, err := db.ExecContext(ctx, sqlstr, t.ID, t.Content, t.Completed); err != nil {
		return logerror(err)
	}
	// set exists
	t._exists = true
	return nil
}

// Update updates a [Todo] in the database.
func (t *Todo) Update(ctx context.Context, db DB) error {
	switch {
	case !t._exists: // doesn't exist
		return logerror(&ErrUpdateFailed{ErrDoesNotExist})
	case t._deleted: // deleted
		return logerror(&ErrUpdateFailed{ErrMarkedForDeletion})
	}
	// update with primary key
	const sqlstr = `UPDATE saas-sample.todos SET ` +
		`content = ?, completed = ? ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, t.Content, t.Completed, t.ID)
	if _, err := db.ExecContext(ctx, sqlstr, t.Content, t.Completed, t.ID); err != nil {
		return logerror(err)
	}
	return nil
}

// Save saves the [Todo] to the database.
func (t *Todo) Save(ctx context.Context, db DB) error {
	if t.Exists() {
		return t.Update(ctx, db)
	}
	return t.Insert(ctx, db)
}

// Upsert performs an upsert for [Todo].
func (t *Todo) Upsert(ctx context.Context, db DB) error {
	switch {
	case t._deleted: // deleted
		return logerror(&ErrUpsertFailed{ErrMarkedForDeletion})
	}
	// upsert
	const sqlstr = `INSERT INTO saas-sample.todos (` +
		`id, content, completed` +
		`) VALUES (` +
		`?, ?, ?` +
		`)` +
		` ON DUPLICATE KEY UPDATE ` +
		`id = VALUES(id), content = VALUES(content), completed = VALUES(completed)`
	// run
	logf(sqlstr, t.ID, t.Content, t.Completed)
	if _, err := db.ExecContext(ctx, sqlstr, t.ID, t.Content, t.Completed); err != nil {
		return logerror(err)
	}
	// set exists
	t._exists = true
	return nil
}

// Delete deletes the [Todo] from the database.
func (t *Todo) Delete(ctx context.Context, db DB) error {
	switch {
	case !t._exists: // doesn't exist
		return nil
	case t._deleted: // deleted
		return nil
	}
	// delete with single primary key
	const sqlstr = `DELETE FROM saas-sample.todos ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, t.ID)
	if _, err := db.ExecContext(ctx, sqlstr, t.ID); err != nil {
		return logerror(err)
	}
	// set deleted
	t._deleted = true
	return nil
}

// TodoByID retrieves a row from 'saas-sample.todos' as a [Todo].
//
// Generated from index 'todos_id_pkey'.
func TodoByID(ctx context.Context, db DB, id int) (*Todo, error) {
	// query
	const sqlstr = `SELECT ` +
		`id, content, completed ` +
		`FROM saas-sample.todos ` +
		`WHERE id = ?`
	// run
	logf(sqlstr, id)
	t := Todo{
		_exists: true,
	}
	if err := db.QueryRowContext(ctx, sqlstr, id).Scan(&t.ID, &t.Content, &t.Completed); err != nil {
		return nil, logerror(err)
	}
	return &t, nil
}