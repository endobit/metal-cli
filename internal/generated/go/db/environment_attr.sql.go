// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: environment_attr.sql

package db

import (
	"context"
)

const createEnvironmentAttribute = `-- name: CreateEnvironmentAttribute :exec

INSERT INTO attributes (
	entity,
	object,
	name
)
VALUES (
	(SELECT id FROM entities WHERE name = 'environment'),
	(SELECT id FROM environments env WHERE env.name = ?1 AND zone = (SELECT id FROM zones z WHERE z.name = ?2)),
	?3
)
`

type CreateEnvironmentAttributeParams struct {
	Environment string
	Zone        string
	Name        string
}

// CREATE
//
//	INSERT INTO attributes (
//		entity,
//		object,
//		name
//	)
//	VALUES (
//		(SELECT id FROM entities WHERE name = 'environment'),
//		(SELECT id FROM environments env WHERE env.name = ?1 AND zone = (SELECT id FROM zones z WHERE z.name = ?2)),
//		?3
//	)
func (q *Queries) CreateEnvironmentAttribute(ctx context.Context, arg CreateEnvironmentAttributeParams) error {
	_, err := q.db.ExecContext(ctx, createEnvironmentAttribute, arg.Environment, arg.Zone, arg.Name)
	return err
}

const readEnvironmentAttribute = `-- name: ReadEnvironmentAttribute :one
SELECT
	a.id,
	env.name AS environment,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	environments env ON a.object = env.id
JOIN
	zones z ON env.zone = z.id
WHERE
	e.name = 'environment'
	AND z.name = ?1
	AND env.name = ?2
	AND a.name = ?3
`

type ReadEnvironmentAttributeParams struct {
	Zone        string
	Environment string
	Attr        string
}

type ReadEnvironmentAttributeRow struct {
	ID          int64
	Environment string
	Name        string
	Value       *string
	IsProtected int64
}

// ReadEnvironmentAttribute
//
//	SELECT
//		a.id,
//		env.name AS environment,
//		a.name,
//		a.value,
//		a.is_protected
//	FROM
//		attributes a
//	JOIN
//		entities e ON a.entity = e.id
//	JOIN
//		environments env ON a.object = env.id
//	JOIN
//		zones z ON env.zone = z.id
//	WHERE
//		e.name = 'environment'
//		AND z.name = ?1
//		AND env.name = ?2
//		AND a.name = ?3
func (q *Queries) ReadEnvironmentAttribute(ctx context.Context, arg ReadEnvironmentAttributeParams) (ReadEnvironmentAttributeRow, error) {
	row := q.db.QueryRowContext(ctx, readEnvironmentAttribute, arg.Zone, arg.Environment, arg.Attr)
	var i ReadEnvironmentAttributeRow
	err := row.Scan(
		&i.ID,
		&i.Environment,
		&i.Name,
		&i.Value,
		&i.IsProtected,
	)
	return i, err
}

const readEnvironmentAttributes = `-- name: ReadEnvironmentAttributes :many

SELECT
	a.id,
	env.name AS environment,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	environments env ON a.object = env.id
JOIN
	zones z ON env.zone = z.id
WHERE
	e.name = 'environment'
ORDER BY
        z.name,
	env.name,
	a.name
LIMIT
	COALESCE(NULLIF(?2, 0), 100) OFFSET COALESCE(?1, 0)
`

type ReadEnvironmentAttributesParams struct {
	Offset interface{}
	Limit  interface{}
}

type ReadEnvironmentAttributesRow struct {
	ID          int64
	Environment string
	Name        string
	Value       *string
	IsProtected int64
}

// READ
//
//	SELECT
//		a.id,
//		env.name AS environment,
//		a.name,
//		a.value,
//		a.is_protected
//	FROM
//		attributes a
//	JOIN
//		entities e ON a.entity = e.id
//	JOIN
//		environments env ON a.object = env.id
//	JOIN
//		zones z ON env.zone = z.id
//	WHERE
//		e.name = 'environment'
//	ORDER BY
//	        z.name,
//		env.name,
//		a.name
//	LIMIT
//		COALESCE(NULLIF(?2, 0), 100) OFFSET COALESCE(?1, 0)
func (q *Queries) ReadEnvironmentAttributes(ctx context.Context, arg ReadEnvironmentAttributesParams) ([]ReadEnvironmentAttributesRow, error) {
	rows, err := q.db.QueryContext(ctx, readEnvironmentAttributes, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadEnvironmentAttributesRow
	for rows.Next() {
		var i ReadEnvironmentAttributesRow
		if err := rows.Scan(
			&i.ID,
			&i.Environment,
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

const readEnvironmentAttributesByEnvironment = `-- name: ReadEnvironmentAttributesByEnvironment :many
SELECT
	a.id,
	env.name AS environment,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	environments env ON a.object = env.id
JOIN
	zones z ON env.zone = z.id
WHERE
	e.name = 'environment'
	AND z.name = ?1
	AND env.name = ?2
ORDER BY
	env.name,
	a.name
LIMIT
	COALESCE(NULLIF(?4, 0), 100) OFFSET COALESCE(?3, 0)
`

type ReadEnvironmentAttributesByEnvironmentParams struct {
	Zone        string
	Environment string
	Offset      interface{}
	Limit       interface{}
}

type ReadEnvironmentAttributesByEnvironmentRow struct {
	ID          int64
	Environment string
	Name        string
	Value       *string
	IsProtected int64
}

// ReadEnvironmentAttributesByEnvironment
//
//	SELECT
//		a.id,
//		env.name AS environment,
//		a.name,
//		a.value,
//		a.is_protected
//	FROM
//		attributes a
//	JOIN
//		entities e ON a.entity = e.id
//	JOIN
//		environments env ON a.object = env.id
//	JOIN
//		zones z ON env.zone = z.id
//	WHERE
//		e.name = 'environment'
//		AND z.name = ?1
//		AND env.name = ?2
//	ORDER BY
//		env.name,
//		a.name
//	LIMIT
//		COALESCE(NULLIF(?4, 0), 100) OFFSET COALESCE(?3, 0)
func (q *Queries) ReadEnvironmentAttributesByEnvironment(ctx context.Context, arg ReadEnvironmentAttributesByEnvironmentParams) ([]ReadEnvironmentAttributesByEnvironmentRow, error) {
	rows, err := q.db.QueryContext(ctx, readEnvironmentAttributesByEnvironment,
		arg.Zone,
		arg.Environment,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadEnvironmentAttributesByEnvironmentRow
	for rows.Next() {
		var i ReadEnvironmentAttributesByEnvironmentRow
		if err := rows.Scan(
			&i.ID,
			&i.Environment,
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

const readEnvironmentAttributesByGlob = `-- name: ReadEnvironmentAttributesByGlob :many
SELECT
	a.id,
	env.name AS environment,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	environments env ON a.object = env.id
JOIN
	zones z ON env.zone = z.id
WHERE
	e.name = 'environment'
	AND z.name = ?1
	AND env.name = ?2
	AND a.name GLOB ?3
ORDER BY
	env.name,
	a.name
LIMIT
	COALESCE(NULLIF(?5, 0), 100) OFFSET COALESCE(?4, 0)
`

type ReadEnvironmentAttributesByGlobParams struct {
	Zone        string
	Environment string
	Glob        string
	Offset      interface{}
	Limit       interface{}
}

type ReadEnvironmentAttributesByGlobRow struct {
	ID          int64
	Environment string
	Name        string
	Value       *string
	IsProtected int64
}

// ReadEnvironmentAttributesByGlob
//
//	SELECT
//		a.id,
//		env.name AS environment,
//		a.name,
//		a.value,
//		a.is_protected
//	FROM
//		attributes a
//	JOIN
//		entities e ON a.entity = e.id
//	JOIN
//		environments env ON a.object = env.id
//	JOIN
//		zones z ON env.zone = z.id
//	WHERE
//		e.name = 'environment'
//		AND z.name = ?1
//		AND env.name = ?2
//		AND a.name GLOB ?3
//	ORDER BY
//		env.name,
//		a.name
//	LIMIT
//		COALESCE(NULLIF(?5, 0), 100) OFFSET COALESCE(?4, 0)
func (q *Queries) ReadEnvironmentAttributesByGlob(ctx context.Context, arg ReadEnvironmentAttributesByGlobParams) ([]ReadEnvironmentAttributesByGlobRow, error) {
	rows, err := q.db.QueryContext(ctx, readEnvironmentAttributesByGlob,
		arg.Zone,
		arg.Environment,
		arg.Glob,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadEnvironmentAttributesByGlobRow
	for rows.Next() {
		var i ReadEnvironmentAttributesByGlobRow
		if err := rows.Scan(
			&i.ID,
			&i.Environment,
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

const readEnvironmentAttributesByZone = `-- name: ReadEnvironmentAttributesByZone :many
SELECT
	a.id,
	env.name AS environment,
	a.name,
	a.value,
	a.is_protected
FROM
	attributes a
JOIN
	entities e ON a.entity = e.id
JOIN
	environments env ON a.object = env.id
JOIN
	zones z ON env.zone = z.id
WHERE
	e.name = 'environment'
	AND z.name = ?1
ORDER BY
	z.name,
	env.name,
	a.name
LIMIT
	COALESCE(NULLIF(?3, 0), 100) OFFSET COALESCE(?2, 0)
`

type ReadEnvironmentAttributesByZoneParams struct {
	Zone   string
	Offset interface{}
	Limit  interface{}
}

type ReadEnvironmentAttributesByZoneRow struct {
	ID          int64
	Environment string
	Name        string
	Value       *string
	IsProtected int64
}

// ReadEnvironmentAttributesByZone
//
//	SELECT
//		a.id,
//		env.name AS environment,
//		a.name,
//		a.value,
//		a.is_protected
//	FROM
//		attributes a
//	JOIN
//		entities e ON a.entity = e.id
//	JOIN
//		environments env ON a.object = env.id
//	JOIN
//		zones z ON env.zone = z.id
//	WHERE
//		e.name = 'environment'
//		AND z.name = ?1
//	ORDER BY
//		z.name,
//		env.name,
//		a.name
//	LIMIT
//		COALESCE(NULLIF(?3, 0), 100) OFFSET COALESCE(?2, 0)
func (q *Queries) ReadEnvironmentAttributesByZone(ctx context.Context, arg ReadEnvironmentAttributesByZoneParams) ([]ReadEnvironmentAttributesByZoneRow, error) {
	rows, err := q.db.QueryContext(ctx, readEnvironmentAttributesByZone, arg.Zone, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadEnvironmentAttributesByZoneRow
	for rows.Next() {
		var i ReadEnvironmentAttributesByZoneRow
		if err := rows.Scan(
			&i.ID,
			&i.Environment,
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

const updateEnvironmentAttributeName = `-- name: UpdateEnvironmentAttributeName :exec

UPDATE
	attributes
SET
	name = ?1
WHERE
	entity = (SELECT id FROM entities WHERE name = 'environment')
	AND object = (
	    SELECT
		id
	    FROM
		environments env
	    WHERE
	    	env.name = ?2
	    	AND zone = (SELECT id FROM zones z WHERE z.name = ?3)
	)
	AND attributes.name = ?4
`

type UpdateEnvironmentAttributeNameParams struct {
	Name        string
	Environment string
	Zone        string
	Attr        string
}

// UPDATE
//
//	UPDATE
//		attributes
//	SET
//		name = ?1
//	WHERE
//		entity = (SELECT id FROM entities WHERE name = 'environment')
//		AND object = (
//		    SELECT
//			id
//		    FROM
//			environments env
//		    WHERE
//		    	env.name = ?2
//		    	AND zone = (SELECT id FROM zones z WHERE z.name = ?3)
//		)
//		AND attributes.name = ?4
func (q *Queries) UpdateEnvironmentAttributeName(ctx context.Context, arg UpdateEnvironmentAttributeNameParams) error {
	_, err := q.db.ExecContext(ctx, updateEnvironmentAttributeName,
		arg.Name,
		arg.Environment,
		arg.Zone,
		arg.Attr,
	)
	return err
}

const updateEnvironmentAttributeProtection = `-- name: UpdateEnvironmentAttributeProtection :exec
UPDATE
	attributes
SET
	is_protected = ?1
WHERE
	entity = (SELECT id FROM entities WHERE name = 'environment')
	AND object = (
	    SELECT
		id
	    FROM
		environments env
	    WHERE
	    	env.name = ?2
	    	AND zone = (SELECT id FROM zones z WHERE z.name = ?3)
	)
	AND attributes.name = ?4
`

type UpdateEnvironmentAttributeProtectionParams struct {
	IsProtected int64
	Environment string
	Zone        string
	Attr        string
}

// UpdateEnvironmentAttributeProtection
//
//	UPDATE
//		attributes
//	SET
//		is_protected = ?1
//	WHERE
//		entity = (SELECT id FROM entities WHERE name = 'environment')
//		AND object = (
//		    SELECT
//			id
//		    FROM
//			environments env
//		    WHERE
//		    	env.name = ?2
//		    	AND zone = (SELECT id FROM zones z WHERE z.name = ?3)
//		)
//		AND attributes.name = ?4
func (q *Queries) UpdateEnvironmentAttributeProtection(ctx context.Context, arg UpdateEnvironmentAttributeProtectionParams) error {
	_, err := q.db.ExecContext(ctx, updateEnvironmentAttributeProtection,
		arg.IsProtected,
		arg.Environment,
		arg.Zone,
		arg.Attr,
	)
	return err
}

const updateEnvironmentAttributeValue = `-- name: UpdateEnvironmentAttributeValue :exec
UPDATE
	attributes
SET
	value = ?1
WHERE
	entity = (SELECT id FROM entities WHERE name = 'environment')
	AND object = (
	    SELECT
		id
	    FROM
		environments env
	    WHERE
	    	env.name = ?2
	    	AND zone = (SELECT id FROM zones z WHERE z.name = ?3)
	)
	AND attributes.name = ?4
`

type UpdateEnvironmentAttributeValueParams struct {
	Value       *string
	Environment string
	Zone        string
	Attr        string
}

// UpdateEnvironmentAttributeValue
//
//	UPDATE
//		attributes
//	SET
//		value = ?1
//	WHERE
//		entity = (SELECT id FROM entities WHERE name = 'environment')
//		AND object = (
//		    SELECT
//			id
//		    FROM
//			environments env
//		    WHERE
//		    	env.name = ?2
//		    	AND zone = (SELECT id FROM zones z WHERE z.name = ?3)
//		)
//		AND attributes.name = ?4
func (q *Queries) UpdateEnvironmentAttributeValue(ctx context.Context, arg UpdateEnvironmentAttributeValueParams) error {
	_, err := q.db.ExecContext(ctx, updateEnvironmentAttributeValue,
		arg.Value,
		arg.Environment,
		arg.Zone,
		arg.Attr,
	)
	return err
}
