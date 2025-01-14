// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: model_attr.sql

package db

import (
	"context"
)

const createModelAttribute = `-- name: CreateModelAttribute :exec

INSERT INTO attributes (
	entity,
	object,
	name
)
VALUES (
	(SELECT id FROM entities WHERE name = 'model'),
	(SELECT id FROM models m, makes mk WHERE m.name = ?1 AND m.make = mk.id AND mk.name = ?2),
	?3
)
`

type CreateModelAttributeParams struct {
	Model string
	Make  string
	Name  string
}

// CREATE
//
//	INSERT INTO attributes (
//		entity,
//		object,
//		name
//	)
//	VALUES (
//		(SELECT id FROM entities WHERE name = 'model'),
//		(SELECT id FROM models m, makes mk WHERE m.name = ?1 AND m.make = mk.id AND mk.name = ?2),
//		?3
//	)
func (q *Queries) CreateModelAttribute(ctx context.Context, arg CreateModelAttributeParams) error {
	_, err := q.db.ExecContext(ctx, createModelAttribute, arg.Model, arg.Make, arg.Name)
	return err
}

const readModelAttribute = `-- name: ReadModelAttribute :one
SELECT
	a.id,
	m.name AS model,
	mk.name AS make,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	models m ON a.object = m.id
JOIN
	makes mk ON m.make = mk.id
WHERE
	a.entity = (SELECT id FROM entities WHERE name = 'model')
	AND mk.name = ?1
	AND m.name = ?2
	AND a.name = ?3
`

type ReadModelAttributeParams struct {
	Make  string
	Model string
	Attr  string
}

type ReadModelAttributeRow struct {
	ID          int64
	Model       string
	Make        string
	Name        string
	Value       *string
	IsProtected int64
}

// ReadModelAttribute
//
//	SELECT
//		a.id,
//		m.name AS model,
//		mk.name AS make,
//		a.name,
//		a.value,
//		a.is_protected
//	FROM
//		attributes a
//	JOIN
//		models m ON a.object = m.id
//	JOIN
//		makes mk ON m.make = mk.id
//	WHERE
//		a.entity = (SELECT id FROM entities WHERE name = 'model')
//		AND mk.name = ?1
//		AND m.name = ?2
//		AND a.name = ?3
func (q *Queries) ReadModelAttribute(ctx context.Context, arg ReadModelAttributeParams) (ReadModelAttributeRow, error) {
	row := q.db.QueryRowContext(ctx, readModelAttribute, arg.Make, arg.Model, arg.Attr)
	var i ReadModelAttributeRow
	err := row.Scan(
		&i.ID,
		&i.Model,
		&i.Make,
		&i.Name,
		&i.Value,
		&i.IsProtected,
	)
	return i, err
}

const readModelAttributes = `-- name: ReadModelAttributes :many

SELECT
	a.id,
	m.name AS model,
	mk.name AS make,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	models m ON a.object = m.id
JOIN
	makes mk ON m.make = mk.id
WHERE
	a.entity = (SELECT id FROM entities WHERE name = 'model')
ORDER BY
        mk.name,
	m.name,
	a.name
LIMIT
	COALESCE(NULLIF(?2, 0), 100) OFFSET COALESCE(?1, 0)
`

type ReadModelAttributesParams struct {
	Offset interface{}
	Limit  interface{}
}

type ReadModelAttributesRow struct {
	ID          int64
	Model       string
	Make        string
	Name        string
	Value       *string
	IsProtected int64
}

// READ
//
//	SELECT
//		a.id,
//		m.name AS model,
//		mk.name AS make,
//		a.name,
//		a.value,
//		a.is_protected
//	FROM
//		attributes a
//	JOIN
//		models m ON a.object = m.id
//	JOIN
//		makes mk ON m.make = mk.id
//	WHERE
//		a.entity = (SELECT id FROM entities WHERE name = 'model')
//	ORDER BY
//	        mk.name,
//		m.name,
//		a.name
//	LIMIT
//		COALESCE(NULLIF(?2, 0), 100) OFFSET COALESCE(?1, 0)
func (q *Queries) ReadModelAttributes(ctx context.Context, arg ReadModelAttributesParams) ([]ReadModelAttributesRow, error) {
	rows, err := q.db.QueryContext(ctx, readModelAttributes, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadModelAttributesRow
	for rows.Next() {
		var i ReadModelAttributesRow
		if err := rows.Scan(
			&i.ID,
			&i.Model,
			&i.Make,
			&i.Name,
			&i.Value,
			&i.IsProtected,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const readModelAttributesByGlob = `-- name: ReadModelAttributesByGlob :many
SELECT
	a.id,
	m.name AS model,
	mk.name AS make,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	models m ON a.object = m.id
JOIN
	makes mk ON m.make = mk.id
WHERE
	a.entity = (SELECT id FROM entities WHERE name = 'model')
	AND mk.name = ?1
	AND m.name = ?2
	AND a.name GLOB ?3
ORDER BY
        mk.name,
	m.name,
	a.name
LIMIT
	COALESCE(NULLIF(?5, 0), 100) OFFSET COALESCE(?4, 0)
`

type ReadModelAttributesByGlobParams struct {
	Make   string
	Model  string
	Glob   string
	Offset interface{}
	Limit  interface{}
}

type ReadModelAttributesByGlobRow struct {
	ID          int64
	Model       string
	Make        string
	Name        string
	Value       *string
	IsProtected int64
}

// ReadModelAttributesByGlob
//
//	SELECT
//		a.id,
//		m.name AS model,
//		mk.name AS make,
//		a.name,
//		a.value,
//		a.is_protected
//	FROM
//		attributes a
//	JOIN
//		models m ON a.object = m.id
//	JOIN
//		makes mk ON m.make = mk.id
//	WHERE
//		a.entity = (SELECT id FROM entities WHERE name = 'model')
//		AND mk.name = ?1
//		AND m.name = ?2
//		AND a.name GLOB ?3
//	ORDER BY
//	        mk.name,
//		m.name,
//		a.name
//	LIMIT
//		COALESCE(NULLIF(?5, 0), 100) OFFSET COALESCE(?4, 0)
func (q *Queries) ReadModelAttributesByGlob(ctx context.Context, arg ReadModelAttributesByGlobParams) ([]ReadModelAttributesByGlobRow, error) {
	rows, err := q.db.QueryContext(ctx, readModelAttributesByGlob,
		arg.Make,
		arg.Model,
		arg.Glob,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadModelAttributesByGlobRow
	for rows.Next() {
		var i ReadModelAttributesByGlobRow
		if err := rows.Scan(
			&i.ID,
			&i.Model,
			&i.Make,
			&i.Name,
			&i.Value,
			&i.IsProtected,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const readModelAttributesByMake = `-- name: ReadModelAttributesByMake :many
SELECT
	a.id,
	m.name AS model,
	mk.name AS make,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	models m ON a.object = m.id
JOIN
	makes mk ON m.make = mk.id
WHERE
	a.entity = (SELECT id FROM entities WHERE name = 'model')
	AND mk.name = ?1
ORDER BY
        mk.name,
	m.name,
	a.name
LIMIT
	COALESCE(NULLIF(?3, 0), 100) OFFSET COALESCE(?2, 0)
`

type ReadModelAttributesByMakeParams struct {
	Make   string
	Offset interface{}
	Limit  interface{}
}

type ReadModelAttributesByMakeRow struct {
	ID          int64
	Model       string
	Make        string
	Name        string
	Value       *string
	IsProtected int64
}

// ReadModelAttributesByMake
//
//	SELECT
//		a.id,
//		m.name AS model,
//		mk.name AS make,
//		a.name,
//		a.value,
//		a.is_protected
//	FROM
//		attributes a
//	JOIN
//		models m ON a.object = m.id
//	JOIN
//		makes mk ON m.make = mk.id
//	WHERE
//		a.entity = (SELECT id FROM entities WHERE name = 'model')
//		AND mk.name = ?1
//	ORDER BY
//	        mk.name,
//		m.name,
//		a.name
//	LIMIT
//		COALESCE(NULLIF(?3, 0), 100) OFFSET COALESCE(?2, 0)
func (q *Queries) ReadModelAttributesByMake(ctx context.Context, arg ReadModelAttributesByMakeParams) ([]ReadModelAttributesByMakeRow, error) {
	rows, err := q.db.QueryContext(ctx, readModelAttributesByMake, arg.Make, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadModelAttributesByMakeRow
	for rows.Next() {
		var i ReadModelAttributesByMakeRow
		if err := rows.Scan(
			&i.ID,
			&i.Model,
			&i.Make,
			&i.Name,
			&i.Value,
			&i.IsProtected,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const readModelAttributesByMakeModel = `-- name: ReadModelAttributesByMakeModel :many
SELECT
	a.id,
	m.name AS model,
	mk.name AS make,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	models m ON a.object = m.id
JOIN
	makes mk ON m.make = mk.id
WHERE
	a.entity = (SELECT id FROM entities WHERE name = 'model')
	AND mk.name = ?1
	AND m.name = ?2
ORDER BY
        mk.name,
	m.name,
	a.name
LIMIT
	COALESCE(NULLIF(?4, 0), 100) OFFSET COALESCE(?3, 0)
`

type ReadModelAttributesByMakeModelParams struct {
	Make   string
	Model  string
	Offset interface{}
	Limit  interface{}
}

type ReadModelAttributesByMakeModelRow struct {
	ID          int64
	Model       string
	Make        string
	Name        string
	Value       *string
	IsProtected int64
}

// ReadModelAttributesByMakeModel
//
//	SELECT
//		a.id,
//		m.name AS model,
//		mk.name AS make,
//		a.name,
//		a.value,
//		a.is_protected
//	FROM
//		attributes a
//	JOIN
//		models m ON a.object = m.id
//	JOIN
//		makes mk ON m.make = mk.id
//	WHERE
//		a.entity = (SELECT id FROM entities WHERE name = 'model')
//		AND mk.name = ?1
//		AND m.name = ?2
//	ORDER BY
//	        mk.name,
//		m.name,
//		a.name
//	LIMIT
//		COALESCE(NULLIF(?4, 0), 100) OFFSET COALESCE(?3, 0)
func (q *Queries) ReadModelAttributesByMakeModel(ctx context.Context, arg ReadModelAttributesByMakeModelParams) ([]ReadModelAttributesByMakeModelRow, error) {
	rows, err := q.db.QueryContext(ctx, readModelAttributesByMakeModel,
		arg.Make,
		arg.Model,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadModelAttributesByMakeModelRow
	for rows.Next() {
		var i ReadModelAttributesByMakeModelRow
		if err := rows.Scan(
			&i.ID,
			&i.Model,
			&i.Make,
			&i.Name,
			&i.Value,
			&i.IsProtected,
		); err != nil {
			return nil, err
		}
		items = append(items, i)
	}
	if err := rows.Close(); err != nil {
		return nil, err
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return items, nil
}

const updateModelAttributeName = `-- name: UpdateModelAttributeName :exec

UPDATE
	attributes
SET
	name = ?1
WHERE
	entity = (SELECT id FROM entities WHERE name = 'model')
	AND object = (SELECT id FROM models m, makes mk WHERE m.name = ?2 AND m.make = mk.id AND mk.name = ?3)
	AND attributes.name = ?4
`

type UpdateModelAttributeNameParams struct {
	Name  string
	Model string
	Make  string
	Attr  string
}

// UPDATE
//
//	UPDATE
//		attributes
//	SET
//		name = ?1
//	WHERE
//		entity = (SELECT id FROM entities WHERE name = 'model')
//		AND object = (SELECT id FROM models m, makes mk WHERE m.name = ?2 AND m.make = mk.id AND mk.name = ?3)
//		AND attributes.name = ?4
func (q *Queries) UpdateModelAttributeName(ctx context.Context, arg UpdateModelAttributeNameParams) error {
	_, err := q.db.ExecContext(ctx, updateModelAttributeName,
		arg.Name,
		arg.Model,
		arg.Make,
		arg.Attr,
	)
	return err
}

const updateModelAttributeProtection = `-- name: UpdateModelAttributeProtection :exec
UPDATE
	attributes
SET
	is_protected = ?1
WHERE
	entity = (SELECT id FROM entities WHERE name = 'model')
	AND object = (SELECT id FROM models m, makes mk WHERE m.name = ?2 AND m.make = mk.id AND mk.name = ?3)
	AND attributes.name = ?4
`

type UpdateModelAttributeProtectionParams struct {
	IsProtected int64
	Model       string
	Make        string
	Attr        string
}

// UpdateModelAttributeProtection
//
//	UPDATE
//		attributes
//	SET
//		is_protected = ?1
//	WHERE
//		entity = (SELECT id FROM entities WHERE name = 'model')
//		AND object = (SELECT id FROM models m, makes mk WHERE m.name = ?2 AND m.make = mk.id AND mk.name = ?3)
//		AND attributes.name = ?4
func (q *Queries) UpdateModelAttributeProtection(ctx context.Context, arg UpdateModelAttributeProtectionParams) error {
	_, err := q.db.ExecContext(ctx, updateModelAttributeProtection,
		arg.IsProtected,
		arg.Model,
		arg.Make,
		arg.Attr,
	)
	return err
}

const updateModelAttributeValue = `-- name: UpdateModelAttributeValue :exec
UPDATE
	attributes
SET
	value = ?1
WHERE
	entity = (SELECT id FROM entities WHERE name = 'model')
	AND object = (SELECT id FROM models m, makes mk WHERE m.name = ?2 AND m.make = mk.id AND mk.name = ?3)
	AND attributes.name = ?4
`

type UpdateModelAttributeValueParams struct {
	Value *string
	Model string
	Make  string
	Attr  string
}

// UpdateModelAttributeValue
//
//	UPDATE
//		attributes
//	SET
//		value = ?1
//	WHERE
//		entity = (SELECT id FROM entities WHERE name = 'model')
//		AND object = (SELECT id FROM models m, makes mk WHERE m.name = ?2 AND m.make = mk.id AND mk.name = ?3)
//		AND attributes.name = ?4
func (q *Queries) UpdateModelAttributeValue(ctx context.Context, arg UpdateModelAttributeValueParams) error {
	_, err := q.db.ExecContext(ctx, updateModelAttributeValue,
		arg.Value,
		arg.Model,
		arg.Make,
		arg.Attr,
	)
	return err
}
