// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: host_interfaces.sql

package db

import (
	"context"
)

const createHostInterface = `-- name: CreateHostInterface :exec

INSERT INTO host_interfaces (
	host,
	name
)
VALUES (
	(
		SELECT id
		FROM hosts h
		WHERE h.name = ?1
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = ?2) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?3 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?2)))
		)
	),
	?4
)
`

type CreateHostInterfaceParams struct {
	Host    string
	Zone    string
	Cluster string
	Name    string
}

// CREATE
//
//	INSERT INTO host_interfaces (
//		host,
//		name
//	)
//	VALUES (
//		(
//			SELECT id
//			FROM hosts h
//			WHERE h.name = ?1
//			AND (
//				(h.zone = (SELECT id FROM zones z WHERE z.name = ?2) AND h.cluster IS NULL)
//				OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?3 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?2)))
//			)
//		),
//		?4
//	)
func (q *Queries) CreateHostInterface(ctx context.Context, arg CreateHostInterfaceParams) error {
	_, err := q.db.ExecContext(ctx, createHostInterface,
		arg.Host,
		arg.Zone,
		arg.Cluster,
		arg.Name,
	)
	return err
}

const deleteHostInterface = `-- name: DeleteHostInterface :exec

DELETE FROM
	host_interfaces
WHERE
	id = ?1
`

type DeleteHostInterfaceParams struct {
	ID int64
}

// DELETE
//
//	DELETE FROM
//		host_interfaces
//	WHERE
//		id = ?1
func (q *Queries) DeleteHostInterface(ctx context.Context, arg DeleteHostInterfaceParams) error {
	_, err := q.db.ExecContext(ctx, deleteHostInterface, arg.ID)
	return err
}

const readHostInterface = `-- name: ReadHostInterface :one
SELECT
	hi.id,
	hi.name,
	hi.ip,
	hi.mac,
	hi.netmask,
	hi.is_dhcp,
	hi.is_pxe,
	hi.is_management,
	hi.type,
	hi.bond_mode,
	mhi.name AS master_interface,
	n.name AS network,
	h.name AS host,
	z.name AS zone,
	c.name AS cluster
FROM
	host_interfaces hi
LEFT JOIN
	host_interfaces mhi ON hi.master = mhi.id
LEFT JOIN
	networks n ON hi.network = n.id
JOIN
	hosts h ON hi.host = h.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
LEFT JOIN
	clusters c ON h.cluster = c.id
WHERE
	h.name = ?1
	AND (
		(h.zone IS NOT NULL AND z.name = ?2)
		OR (h.cluster IS NOT NULL AND c.name = ?3 AND z.name = ?2)
	)
	AND hi.name = ?4
`

type ReadHostInterfaceParams struct {
	Host    string
	Zone    string
	Cluster string
	Name    string
}

type ReadHostInterfaceRow struct {
	ID              int64
	Name            string
	IP              *string
	MAC             *string
	Netmask         *string
	IsDHCP          int64
	IsPXE           int64
	IsManagement    int64
	Type            *string
	BondMode        *string
	MasterInterface string
	Network         string
	Host            string
	Zone            *string
	Cluster         *string
}

// ReadHostInterface
//
//	SELECT
//		hi.id,
//		hi.name,
//		hi.ip,
//		hi.mac,
//		hi.netmask,
//		hi.is_dhcp,
//		hi.is_pxe,
//		hi.is_management,
//		hi.type,
//		hi.bond_mode,
//		mhi.name AS master_interface,
//		n.name AS network,
//		h.name AS host,
//		z.name AS zone,
//		c.name AS cluster
//	FROM
//		host_interfaces hi
//	LEFT JOIN
//		host_interfaces mhi ON hi.master = mhi.id
//	LEFT JOIN
//		networks n ON hi.network = n.id
//	JOIN
//		hosts h ON hi.host = h.id
//	LEFT JOIN
//		zones z ON COALESCE(h.zone, c.zone) = z.id
//	LEFT JOIN
//		clusters c ON h.cluster = c.id
//	WHERE
//		h.name = ?1
//		AND (
//			(h.zone IS NOT NULL AND z.name = ?2)
//			OR (h.cluster IS NOT NULL AND c.name = ?3 AND z.name = ?2)
//		)
//		AND hi.name = ?4
func (q *Queries) ReadHostInterface(ctx context.Context, arg ReadHostInterfaceParams) (ReadHostInterfaceRow, error) {
	row := q.db.QueryRowContext(ctx, readHostInterface,
		arg.Host,
		arg.Zone,
		arg.Cluster,
		arg.Name,
	)
	var i ReadHostInterfaceRow
	err := row.Scan(
		&i.ID,
		&i.Name,
		&i.IP,
		&i.MAC,
		&i.Netmask,
		&i.IsDHCP,
		&i.IsPXE,
		&i.IsManagement,
		&i.Type,
		&i.BondMode,
		&i.MasterInterface,
		&i.Network,
		&i.Host,
		&i.Zone,
		&i.Cluster,
	)
	return i, err
}

const readHostInterfaces = `-- name: ReadHostInterfaces :many

SELECT
	hi.id,
	hi.name,
	hi.ip,
	hi.mac,
	hi.netmask,
	hi.is_dhcp,
	hi.is_pxe,
	hi.is_management,
	hi.type,
	hi.bond_mode,
	mhi.name AS master_interface,
	n.name AS network,
	h.name AS host,
	z.name AS zone,
	c.name AS cluster
FROM
	host_interfaces hi
LEFT JOIN
	host_interfaces mhi ON hi.master = mhi.id
LEFT JOIN
	networks n ON hi.network = n.id
JOIN
	hosts h ON hi.host = h.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
LEFT JOIN
	clusters c ON h.cluster = c.id
ORDER BY
	z.name,
	COALESCE(c.name, ''),
	h.name,
	hi.name
LIMIT
	COALESCE(NULLIF(?2, 0), 100) OFFSET COALESCE(?1, 0)
`

type ReadHostInterfacesParams struct {
	Offset interface{}
	Limit  interface{}
}

type ReadHostInterfacesRow struct {
	ID              int64
	Name            string
	IP              *string
	MAC             *string
	Netmask         *string
	IsDHCP          int64
	IsPXE           int64
	IsManagement    int64
	Type            *string
	BondMode        *string
	MasterInterface string
	Network         string
	Host            string
	Zone            *string
	Cluster         *string
}

// READ
//
//	SELECT
//		hi.id,
//		hi.name,
//		hi.ip,
//		hi.mac,
//		hi.netmask,
//		hi.is_dhcp,
//		hi.is_pxe,
//		hi.is_management,
//		hi.type,
//		hi.bond_mode,
//		mhi.name AS master_interface,
//		n.name AS network,
//		h.name AS host,
//		z.name AS zone,
//		c.name AS cluster
//	FROM
//		host_interfaces hi
//	LEFT JOIN
//		host_interfaces mhi ON hi.master = mhi.id
//	LEFT JOIN
//		networks n ON hi.network = n.id
//	JOIN
//		hosts h ON hi.host = h.id
//	LEFT JOIN
//		zones z ON COALESCE(h.zone, c.zone) = z.id
//	LEFT JOIN
//		clusters c ON h.cluster = c.id
//	ORDER BY
//		z.name,
//		COALESCE(c.name, ''),
//		h.name,
//		hi.name
//	LIMIT
//		COALESCE(NULLIF(?2, 0), 100) OFFSET COALESCE(?1, 0)
func (q *Queries) ReadHostInterfaces(ctx context.Context, arg ReadHostInterfacesParams) ([]ReadHostInterfacesRow, error) {
	rows, err := q.db.QueryContext(ctx, readHostInterfaces, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadHostInterfacesRow
	for rows.Next() {
		var i ReadHostInterfacesRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.IP,
			&i.MAC,
			&i.Netmask,
			&i.IsDHCP,
			&i.IsPXE,
			&i.IsManagement,
			&i.Type,
			&i.BondMode,
			&i.MasterInterface,
			&i.Network,
			&i.Host,
			&i.Zone,
			&i.Cluster,
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

const readHostInterfacesByCluster = `-- name: ReadHostInterfacesByCluster :many
SELECT
	hi.id,
	hi.name,
	hi.ip,
	hi.mac,
	hi.netmask,
	hi.is_dhcp,
	hi.is_pxe,
	hi.is_management,
	hi.type,
	hi.bond_mode,
	mhi.name AS master_interface,
	n.name AS network,
	h.name AS host,
	z.name AS zone,
	c.name AS cluster
FROM
	host_interfaces hi
LEFT JOIN
	host_interfaces mhi ON hi.master = mhi.id
LEFT JOIN
	networks n ON hi.network = n.id
JOIN
	hosts h ON hi.host = h.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
LEFT JOIN
	clusters c ON h.cluster = c.id
WHERE
	c.name = ?1
	AND
	z.name = ?2
ORDER BY
	z.name,
	c.name,
	h.name,
	hi.name
LIMIT
	COALESCE(NULLIF(?4, 0), 100) OFFSET COALESCE(?3, 0)
`

type ReadHostInterfacesByClusterParams struct {
	Cluster string
	Zone    string
	Offset  interface{}
	Limit   interface{}
}

type ReadHostInterfacesByClusterRow struct {
	ID              int64
	Name            string
	IP              *string
	MAC             *string
	Netmask         *string
	IsDHCP          int64
	IsPXE           int64
	IsManagement    int64
	Type            *string
	BondMode        *string
	MasterInterface string
	Network         string
	Host            string
	Zone            *string
	Cluster         *string
}

// ReadHostInterfacesByCluster
//
//	SELECT
//		hi.id,
//		hi.name,
//		hi.ip,
//		hi.mac,
//		hi.netmask,
//		hi.is_dhcp,
//		hi.is_pxe,
//		hi.is_management,
//		hi.type,
//		hi.bond_mode,
//		mhi.name AS master_interface,
//		n.name AS network,
//		h.name AS host,
//		z.name AS zone,
//		c.name AS cluster
//	FROM
//		host_interfaces hi
//	LEFT JOIN
//		host_interfaces mhi ON hi.master = mhi.id
//	LEFT JOIN
//		networks n ON hi.network = n.id
//	JOIN
//		hosts h ON hi.host = h.id
//	LEFT JOIN
//		zones z ON COALESCE(h.zone, c.zone) = z.id
//	LEFT JOIN
//		clusters c ON h.cluster = c.id
//	WHERE
//		c.name = ?1
//		AND
//		z.name = ?2
//	ORDER BY
//		z.name,
//		c.name,
//		h.name,
//		hi.name
//	LIMIT
//		COALESCE(NULLIF(?4, 0), 100) OFFSET COALESCE(?3, 0)
func (q *Queries) ReadHostInterfacesByCluster(ctx context.Context, arg ReadHostInterfacesByClusterParams) ([]ReadHostInterfacesByClusterRow, error) {
	rows, err := q.db.QueryContext(ctx, readHostInterfacesByCluster,
		arg.Cluster,
		arg.Zone,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadHostInterfacesByClusterRow
	for rows.Next() {
		var i ReadHostInterfacesByClusterRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.IP,
			&i.MAC,
			&i.Netmask,
			&i.IsDHCP,
			&i.IsPXE,
			&i.IsManagement,
			&i.Type,
			&i.BondMode,
			&i.MasterInterface,
			&i.Network,
			&i.Host,
			&i.Zone,
			&i.Cluster,
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

const readHostInterfacesByGlob = `-- name: ReadHostInterfacesByGlob :many
SELECT
	hi.id,
	hi.name,
	hi.ip,
	hi.mac,
	hi.netmask,
	hi.is_dhcp,
	hi.is_pxe,
	hi.is_management,
	hi.type,
	hi.bond_mode,
	mhi.name AS master_interface,
	n.name AS network,
	h.name AS host,
	z.name AS zone,
	c.name AS cluster
FROM
	host_interfaces hi
LEFT JOIN
	host_interfaces mhi ON hi.master = mhi.id
LEFT JOIN
	networks n ON hi.network = n.id
JOIN
	hosts h ON hi.host = h.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
LEFT JOIN
	clusters c ON h.cluster = c.id
WHERE
	h.name = ?1
	AND (
		(h.zone IS NOT NULL AND z.name = ?2)
		OR (h.cluster IS NOT NULL AND c.name = ?3 AND z.name = ?2)
	)
	AND hi.name GLOB ?4
ORDER BY
        z.name,
	COALESCE(c.name, ''),
	h.name,
	hi.name
LIMIT
	COALESCE(NULLIF(?6, 0), 100) OFFSET COALESCE(?5, 0)
`

type ReadHostInterfacesByGlobParams struct {
	Host    string
	Zone    string
	Cluster string
	Glob    string
	Offset  interface{}
	Limit   interface{}
}

type ReadHostInterfacesByGlobRow struct {
	ID              int64
	Name            string
	IP              *string
	MAC             *string
	Netmask         *string
	IsDHCP          int64
	IsPXE           int64
	IsManagement    int64
	Type            *string
	BondMode        *string
	MasterInterface string
	Network         string
	Host            string
	Zone            *string
	Cluster         *string
}

// ReadHostInterfacesByGlob
//
//	SELECT
//		hi.id,
//		hi.name,
//		hi.ip,
//		hi.mac,
//		hi.netmask,
//		hi.is_dhcp,
//		hi.is_pxe,
//		hi.is_management,
//		hi.type,
//		hi.bond_mode,
//		mhi.name AS master_interface,
//		n.name AS network,
//		h.name AS host,
//		z.name AS zone,
//		c.name AS cluster
//	FROM
//		host_interfaces hi
//	LEFT JOIN
//		host_interfaces mhi ON hi.master = mhi.id
//	LEFT JOIN
//		networks n ON hi.network = n.id
//	JOIN
//		hosts h ON hi.host = h.id
//	LEFT JOIN
//		zones z ON COALESCE(h.zone, c.zone) = z.id
//	LEFT JOIN
//		clusters c ON h.cluster = c.id
//	WHERE
//		h.name = ?1
//		AND (
//			(h.zone IS NOT NULL AND z.name = ?2)
//			OR (h.cluster IS NOT NULL AND c.name = ?3 AND z.name = ?2)
//		)
//		AND hi.name GLOB ?4
//	ORDER BY
//	        z.name,
//		COALESCE(c.name, ''),
//		h.name,
//		hi.name
//	LIMIT
//		COALESCE(NULLIF(?6, 0), 100) OFFSET COALESCE(?5, 0)
func (q *Queries) ReadHostInterfacesByGlob(ctx context.Context, arg ReadHostInterfacesByGlobParams) ([]ReadHostInterfacesByGlobRow, error) {
	rows, err := q.db.QueryContext(ctx, readHostInterfacesByGlob,
		arg.Host,
		arg.Zone,
		arg.Cluster,
		arg.Glob,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadHostInterfacesByGlobRow
	for rows.Next() {
		var i ReadHostInterfacesByGlobRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.IP,
			&i.MAC,
			&i.Netmask,
			&i.IsDHCP,
			&i.IsPXE,
			&i.IsManagement,
			&i.Type,
			&i.BondMode,
			&i.MasterInterface,
			&i.Network,
			&i.Host,
			&i.Zone,
			&i.Cluster,
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

const readHostInterfacesByHost = `-- name: ReadHostInterfacesByHost :many
SELECT
	hi.id,
	hi.name,
	hi.ip,
	hi.mac,
	hi.netmask,
	hi.is_dhcp,
	hi.is_pxe,
	hi.is_management,
	hi.type,
	hi.bond_mode,
	mhi.name AS master_interface,
	n.name AS network,
	h.name AS host,
	z.name AS zone,
	c.name AS cluster
FROM
	host_interfaces hi
LEFT JOIN
	host_interfaces mhi ON hi.master = mhi.id
LEFT JOIN
	networks n ON hi.network = n.id
JOIN
	hosts h ON hi.host = h.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
LEFT JOIN
	clusters c ON h.cluster = c.id
WHERE
	h.name = ?1
	AND (
		(h.zone IS NOT NULL AND z.name = ?2)
		OR (h.cluster IS NOT NULL AND c.name = ?3 AND z.name = ?2)
	)
ORDER BY
	z.name,
	COALESCE(c.name, ''),
	h.name,
	hi.name
LIMIT
	COALESCE(NULLIF(?5, 0), 100) OFFSET COALESCE(?4, 0)
`

type ReadHostInterfacesByHostParams struct {
	Host    string
	Zone    string
	Cluster string
	Offset  interface{}
	Limit   interface{}
}

type ReadHostInterfacesByHostRow struct {
	ID              int64
	Name            string
	IP              *string
	MAC             *string
	Netmask         *string
	IsDHCP          int64
	IsPXE           int64
	IsManagement    int64
	Type            *string
	BondMode        *string
	MasterInterface string
	Network         string
	Host            string
	Zone            *string
	Cluster         *string
}

// ReadHostInterfacesByHost
//
//	SELECT
//		hi.id,
//		hi.name,
//		hi.ip,
//		hi.mac,
//		hi.netmask,
//		hi.is_dhcp,
//		hi.is_pxe,
//		hi.is_management,
//		hi.type,
//		hi.bond_mode,
//		mhi.name AS master_interface,
//		n.name AS network,
//		h.name AS host,
//		z.name AS zone,
//		c.name AS cluster
//	FROM
//		host_interfaces hi
//	LEFT JOIN
//		host_interfaces mhi ON hi.master = mhi.id
//	LEFT JOIN
//		networks n ON hi.network = n.id
//	JOIN
//		hosts h ON hi.host = h.id
//	LEFT JOIN
//		zones z ON COALESCE(h.zone, c.zone) = z.id
//	LEFT JOIN
//		clusters c ON h.cluster = c.id
//	WHERE
//		h.name = ?1
//		AND (
//			(h.zone IS NOT NULL AND z.name = ?2)
//			OR (h.cluster IS NOT NULL AND c.name = ?3 AND z.name = ?2)
//		)
//	ORDER BY
//		z.name,
//		COALESCE(c.name, ''),
//		h.name,
//		hi.name
//	LIMIT
//		COALESCE(NULLIF(?5, 0), 100) OFFSET COALESCE(?4, 0)
func (q *Queries) ReadHostInterfacesByHost(ctx context.Context, arg ReadHostInterfacesByHostParams) ([]ReadHostInterfacesByHostRow, error) {
	rows, err := q.db.QueryContext(ctx, readHostInterfacesByHost,
		arg.Host,
		arg.Zone,
		arg.Cluster,
		arg.Offset,
		arg.Limit,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadHostInterfacesByHostRow
	for rows.Next() {
		var i ReadHostInterfacesByHostRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.IP,
			&i.MAC,
			&i.Netmask,
			&i.IsDHCP,
			&i.IsPXE,
			&i.IsManagement,
			&i.Type,
			&i.BondMode,
			&i.MasterInterface,
			&i.Network,
			&i.Host,
			&i.Zone,
			&i.Cluster,
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

const readHostInterfacesByZone = `-- name: ReadHostInterfacesByZone :many
SELECT
	hi.id,
	hi.name,
	hi.ip,
	hi.mac,
	hi.netmask,
	hi.is_dhcp,
	hi.is_pxe,
	hi.is_management,
	hi.type,
	hi.bond_mode,
	mhi.name AS master_interface,
	n.name AS network,
	h.name AS host,
	z.name AS zone,
	c.name AS cluster
FROM
	host_interfaces hi
LEFT JOIN
	host_interfaces mhi ON hi.master = mhi.id
LEFT JOIN
	networks n ON hi.network = n.id
JOIN
	hosts h ON hi.host = h.id
LEFT JOIN
	zones z ON COALESCE(h.zone, c.zone) = z.id
LEFT JOIN
	clusters c ON h.cluster = c.id
WHERE
	z.name = ?1
ORDER BY
	COALESCE(c.name, ''),
	h.name,
	hi.name
LIMIT
	COALESCE(NULLIF(?3, 0), 100) OFFSET COALESCE(?2, 0)
`

type ReadHostInterfacesByZoneParams struct {
	Zone   string
	Offset interface{}
	Limit  interface{}
}

type ReadHostInterfacesByZoneRow struct {
	ID              int64
	Name            string
	IP              *string
	MAC             *string
	Netmask         *string
	IsDHCP          int64
	IsPXE           int64
	IsManagement    int64
	Type            *string
	BondMode        *string
	MasterInterface string
	Network         string
	Host            string
	Zone            *string
	Cluster         *string
}

// ReadHostInterfacesByZone
//
//	SELECT
//		hi.id,
//		hi.name,
//		hi.ip,
//		hi.mac,
//		hi.netmask,
//		hi.is_dhcp,
//		hi.is_pxe,
//		hi.is_management,
//		hi.type,
//		hi.bond_mode,
//		mhi.name AS master_interface,
//		n.name AS network,
//		h.name AS host,
//		z.name AS zone,
//		c.name AS cluster
//	FROM
//		host_interfaces hi
//	LEFT JOIN
//		host_interfaces mhi ON hi.master = mhi.id
//	LEFT JOIN
//		networks n ON hi.network = n.id
//	JOIN
//		hosts h ON hi.host = h.id
//	LEFT JOIN
//		zones z ON COALESCE(h.zone, c.zone) = z.id
//	LEFT JOIN
//		clusters c ON h.cluster = c.id
//	WHERE
//		z.name = ?1
//	ORDER BY
//		COALESCE(c.name, ''),
//		h.name,
//		hi.name
//	LIMIT
//		COALESCE(NULLIF(?3, 0), 100) OFFSET COALESCE(?2, 0)
func (q *Queries) ReadHostInterfacesByZone(ctx context.Context, arg ReadHostInterfacesByZoneParams) ([]ReadHostInterfacesByZoneRow, error) {
	rows, err := q.db.QueryContext(ctx, readHostInterfacesByZone, arg.Zone, arg.Offset, arg.Limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var items []ReadHostInterfacesByZoneRow
	for rows.Next() {
		var i ReadHostInterfacesByZoneRow
		if err := rows.Scan(
			&i.ID,
			&i.Name,
			&i.IP,
			&i.MAC,
			&i.Netmask,
			&i.IsDHCP,
			&i.IsPXE,
			&i.IsManagement,
			&i.Type,
			&i.BondMode,
			&i.MasterInterface,
			&i.Network,
			&i.Host,
			&i.Zone,
			&i.Cluster,
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

const updateHostInterfaceBondMode = `-- name: UpdateHostInterfaceBondMode :exec
UPDATE
	host_interfaces
SET
	bond_mode = ?1
WHERE
	host_interfaces.name = ?2
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = ?3
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = ?4) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?4)))
		)
	)
`

type UpdateHostInterfaceBondModeParams struct {
	Bond      *string
	Interface string
	Host      string
	Zone      string
	Cluster   string
}

// UpdateHostInterfaceBondMode
//
//	UPDATE
//		host_interfaces
//	SET
//		bond_mode = ?1
//	WHERE
//		host_interfaces.name = ?2
//		AND host = (
//			SELECT id
//			FROM hosts h
//			WHERE h.name = ?3
//			AND (
//				(h.zone = (SELECT id FROM zones z WHERE z.name = ?4) AND h.cluster IS NULL)
//				OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?4)))
//			)
//		)
func (q *Queries) UpdateHostInterfaceBondMode(ctx context.Context, arg UpdateHostInterfaceBondModeParams) error {
	_, err := q.db.ExecContext(ctx, updateHostInterfaceBondMode,
		arg.Bond,
		arg.Interface,
		arg.Host,
		arg.Zone,
		arg.Cluster,
	)
	return err
}

const updateHostInterfaceDHCP = `-- name: UpdateHostInterfaceDHCP :exec
UPDATE
	host_interfaces
SET
	is_dhcp = ?1
WHERE
	host_interfaces.name = ?2
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = ?3
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = ?4) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?4)))
		)
	)
`

type UpdateHostInterfaceDHCPParams struct {
	DHCP      int64
	Interface string
	Host      string
	Zone      string
	Cluster   string
}

// UpdateHostInterfaceDHCP
//
//	UPDATE
//		host_interfaces
//	SET
//		is_dhcp = ?1
//	WHERE
//		host_interfaces.name = ?2
//		AND host = (
//			SELECT id
//			FROM hosts h
//			WHERE h.name = ?3
//			AND (
//				(h.zone = (SELECT id FROM zones z WHERE z.name = ?4) AND h.cluster IS NULL)
//				OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?4)))
//			)
//		)
func (q *Queries) UpdateHostInterfaceDHCP(ctx context.Context, arg UpdateHostInterfaceDHCPParams) error {
	_, err := q.db.ExecContext(ctx, updateHostInterfaceDHCP,
		arg.DHCP,
		arg.Interface,
		arg.Host,
		arg.Zone,
		arg.Cluster,
	)
	return err
}

const updateHostInterfaceIP = `-- name: UpdateHostInterfaceIP :exec
UPDATE
	host_interfaces
SET
	ip = ?1
WHERE
	host_interfaces.name = ?2
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = ?3
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = ?4) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?4)))
		)
	)
`

type UpdateHostInterfaceIPParams struct {
	IP        *string
	Interface string
	Host      string
	Zone      string
	Cluster   string
}

// UpdateHostInterfaceIP
//
//	UPDATE
//		host_interfaces
//	SET
//		ip = ?1
//	WHERE
//		host_interfaces.name = ?2
//		AND host = (
//			SELECT id
//			FROM hosts h
//			WHERE h.name = ?3
//			AND (
//				(h.zone = (SELECT id FROM zones z WHERE z.name = ?4) AND h.cluster IS NULL)
//				OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?4)))
//			)
//		)
func (q *Queries) UpdateHostInterfaceIP(ctx context.Context, arg UpdateHostInterfaceIPParams) error {
	_, err := q.db.ExecContext(ctx, updateHostInterfaceIP,
		arg.IP,
		arg.Interface,
		arg.Host,
		arg.Zone,
		arg.Cluster,
	)
	return err
}

const updateHostInterfaceMAC = `-- name: UpdateHostInterfaceMAC :exec
UPDATE
	host_interfaces
SET
	mac = ?1
WHERE
	host_interfaces.name = ?2
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = ?3
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = ?4) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?4)))
		)
	)
`

type UpdateHostInterfaceMACParams struct {
	MAC       *string
	Interface string
	Host      string
	Zone      string
	Cluster   string
}

// UpdateHostInterfaceMAC
//
//	UPDATE
//		host_interfaces
//	SET
//		mac = ?1
//	WHERE
//		host_interfaces.name = ?2
//		AND host = (
//			SELECT id
//			FROM hosts h
//			WHERE h.name = ?3
//			AND (
//				(h.zone = (SELECT id FROM zones z WHERE z.name = ?4) AND h.cluster IS NULL)
//				OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?4)))
//			)
//		)
func (q *Queries) UpdateHostInterfaceMAC(ctx context.Context, arg UpdateHostInterfaceMACParams) error {
	_, err := q.db.ExecContext(ctx, updateHostInterfaceMAC,
		arg.MAC,
		arg.Interface,
		arg.Host,
		arg.Zone,
		arg.Cluster,
	)
	return err
}

const updateHostInterfaceManagement = `-- name: UpdateHostInterfaceManagement :exec
UPDATE
	host_interfaces
SET
	is_management = ?1
WHERE
	host_interfaces.name = ?2
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = ?3
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = ?4) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?4)))
		)
	)
`

type UpdateHostInterfaceManagementParams struct {
	Management int64
	Interface  string
	Host       string
	Zone       string
	Cluster    string
}

// UpdateHostInterfaceManagement
//
//	UPDATE
//		host_interfaces
//	SET
//		is_management = ?1
//	WHERE
//		host_interfaces.name = ?2
//		AND host = (
//			SELECT id
//			FROM hosts h
//			WHERE h.name = ?3
//			AND (
//				(h.zone = (SELECT id FROM zones z WHERE z.name = ?4) AND h.cluster IS NULL)
//				OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?4)))
//			)
//		)
func (q *Queries) UpdateHostInterfaceManagement(ctx context.Context, arg UpdateHostInterfaceManagementParams) error {
	_, err := q.db.ExecContext(ctx, updateHostInterfaceManagement,
		arg.Management,
		arg.Interface,
		arg.Host,
		arg.Zone,
		arg.Cluster,
	)
	return err
}

const updateHostInterfaceMaster = `-- name: UpdateHostInterfaceMaster :exec
UPDATE
	host_interfaces
SET
	master = (
		SELECT id
		FROM host_interfaces hi
		WHERE hi.name = ?1
		AND host = (
			SELECT id
			FROM hosts h
			WHERE h.name = ?2
			AND (
				(h.zone = (SELECT id FROM zones z WHERE z.name = ?3) AND h.cluster IS NULL)
				OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?4 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?3)))
			)
		)
	)
WHERE
	host_interfaces.name = ?5
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = ?2
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = ?3) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?4 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?3)))
		)
	)
`

type UpdateHostInterfaceMasterParams struct {
	Master    string
	Host      string
	Zone      string
	Cluster   string
	Interface string
}

// UpdateHostInterfaceMaster
//
//	UPDATE
//		host_interfaces
//	SET
//		master = (
//			SELECT id
//			FROM host_interfaces hi
//			WHERE hi.name = ?1
//			AND host = (
//				SELECT id
//				FROM hosts h
//				WHERE h.name = ?2
//				AND (
//					(h.zone = (SELECT id FROM zones z WHERE z.name = ?3) AND h.cluster IS NULL)
//					OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?4 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?3)))
//				)
//			)
//		)
//	WHERE
//		host_interfaces.name = ?5
//		AND host = (
//			SELECT id
//			FROM hosts h
//			WHERE h.name = ?2
//			AND (
//				(h.zone = (SELECT id FROM zones z WHERE z.name = ?3) AND h.cluster IS NULL)
//				OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?4 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?3)))
//			)
//		)
func (q *Queries) UpdateHostInterfaceMaster(ctx context.Context, arg UpdateHostInterfaceMasterParams) error {
	_, err := q.db.ExecContext(ctx, updateHostInterfaceMaster,
		arg.Master,
		arg.Host,
		arg.Zone,
		arg.Cluster,
		arg.Interface,
	)
	return err
}

const updateHostInterfaceName = `-- name: UpdateHostInterfaceName :exec

UPDATE
	host_interfaces
SET
	name = ?1
WHERE
	host_interfaces.name = ?2
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = ?3
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = ?4) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?4)))
		)
	)
`

type UpdateHostInterfaceNameParams struct {
	Name      string
	Interface string
	Host      string
	Zone      string
	Cluster   string
}

// UPDATE
//
//	UPDATE
//		host_interfaces
//	SET
//		name = ?1
//	WHERE
//		host_interfaces.name = ?2
//		AND host = (
//			SELECT id
//			FROM hosts h
//			WHERE h.name = ?3
//			AND (
//				(h.zone = (SELECT id FROM zones z WHERE z.name = ?4) AND h.cluster IS NULL)
//				OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?4)))
//			)
//		)
func (q *Queries) UpdateHostInterfaceName(ctx context.Context, arg UpdateHostInterfaceNameParams) error {
	_, err := q.db.ExecContext(ctx, updateHostInterfaceName,
		arg.Name,
		arg.Interface,
		arg.Host,
		arg.Zone,
		arg.Cluster,
	)
	return err
}

const updateHostInterfaceNetmask = `-- name: UpdateHostInterfaceNetmask :exec
UPDATE
	host_interfaces
SET
	netmask = ?1
WHERE
	host_interfaces.name = ?2
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = ?3
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = ?4) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?4)))
		)
	)
`

type UpdateHostInterfaceNetmaskParams struct {
	Mask      *string
	Interface string
	Host      string
	Zone      string
	Cluster   string
}

// UpdateHostInterfaceNetmask
//
//	UPDATE
//		host_interfaces
//	SET
//		netmask = ?1
//	WHERE
//		host_interfaces.name = ?2
//		AND host = (
//			SELECT id
//			FROM hosts h
//			WHERE h.name = ?3
//			AND (
//				(h.zone = (SELECT id FROM zones z WHERE z.name = ?4) AND h.cluster IS NULL)
//				OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?4)))
//			)
//		)
func (q *Queries) UpdateHostInterfaceNetmask(ctx context.Context, arg UpdateHostInterfaceNetmaskParams) error {
	_, err := q.db.ExecContext(ctx, updateHostInterfaceNetmask,
		arg.Mask,
		arg.Interface,
		arg.Host,
		arg.Zone,
		arg.Cluster,
	)
	return err
}

const updateHostInterfaceNetwork = `-- name: UpdateHostInterfaceNetwork :exec
UPDATE
	host_interfaces
SET
	network = (
		SELECT id
		FROM networks n
		WHERE n.name = ?1
		AND n.zone = (SELECT id FROM zones z WHERE z.name = ?2)
	)
WHERE
	host_interfaces.name = ?3
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = ?4
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = ?2) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?2)))
		)
	)
`

type UpdateHostInterfaceNetworkParams struct {
	Network   string
	Zone      string
	Interface string
	Host      string
	Cluster   string
}

// UpdateHostInterfaceNetwork
//
//	UPDATE
//		host_interfaces
//	SET
//		network = (
//			SELECT id
//			FROM networks n
//			WHERE n.name = ?1
//			AND n.zone = (SELECT id FROM zones z WHERE z.name = ?2)
//		)
//	WHERE
//		host_interfaces.name = ?3
//		AND host = (
//			SELECT id
//			FROM hosts h
//			WHERE h.name = ?4
//			AND (
//				(h.zone = (SELECT id FROM zones z WHERE z.name = ?2) AND h.cluster IS NULL)
//				OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?2)))
//			)
//		)
func (q *Queries) UpdateHostInterfaceNetwork(ctx context.Context, arg UpdateHostInterfaceNetworkParams) error {
	_, err := q.db.ExecContext(ctx, updateHostInterfaceNetwork,
		arg.Network,
		arg.Zone,
		arg.Interface,
		arg.Host,
		arg.Cluster,
	)
	return err
}

const updateHostInterfacePXE = `-- name: UpdateHostInterfacePXE :exec
UPDATE
	host_interfaces
SET
	is_pxe = ?1
WHERE
	host_interfaces.name = ?2
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = ?3
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = ?4) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?4)))
		)
	)
`

type UpdateHostInterfacePXEParams struct {
	PXE       int64
	Interface string
	Host      string
	Zone      string
	Cluster   string
}

// UpdateHostInterfacePXE
//
//	UPDATE
//		host_interfaces
//	SET
//		is_pxe = ?1
//	WHERE
//		host_interfaces.name = ?2
//		AND host = (
//			SELECT id
//			FROM hosts h
//			WHERE h.name = ?3
//			AND (
//				(h.zone = (SELECT id FROM zones z WHERE z.name = ?4) AND h.cluster IS NULL)
//				OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?4)))
//			)
//		)
func (q *Queries) UpdateHostInterfacePXE(ctx context.Context, arg UpdateHostInterfacePXEParams) error {
	_, err := q.db.ExecContext(ctx, updateHostInterfacePXE,
		arg.PXE,
		arg.Interface,
		arg.Host,
		arg.Zone,
		arg.Cluster,
	)
	return err
}

const updateHostInterfaceType = `-- name: UpdateHostInterfaceType :exec
UPDATE
	host_interfaces
SET
	type = ?1
WHERE
	host_interfaces.name = ?2
	AND host = (
		SELECT id
		FROM hosts h
		WHERE h.name = ?3
		AND (
			(h.zone = (SELECT id FROM zones z WHERE z.name = ?4) AND h.cluster IS NULL)
			OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?4)))
		)
	)
`

type UpdateHostInterfaceTypeParams struct {
	Type      *string
	Interface string
	Host      string
	Zone      string
	Cluster   string
}

// UpdateHostInterfaceType
//
//	UPDATE
//		host_interfaces
//	SET
//		type = ?1
//	WHERE
//		host_interfaces.name = ?2
//		AND host = (
//			SELECT id
//			FROM hosts h
//			WHERE h.name = ?3
//			AND (
//				(h.zone = (SELECT id FROM zones z WHERE z.name = ?4) AND h.cluster IS NULL)
//				OR (h.cluster = (SELECT id FROM clusters c WHERE c.name = ?5 AND h.zone = (SELECT id FROM zones z WHERE z.name = ?4)))
//			)
//		)
func (q *Queries) UpdateHostInterfaceType(ctx context.Context, arg UpdateHostInterfaceTypeParams) error {
	_, err := q.db.ExecContext(ctx, updateHostInterfaceType,
		arg.Type,
		arg.Interface,
		arg.Host,
		arg.Zone,
		arg.Cluster,
	)
	return err
}
