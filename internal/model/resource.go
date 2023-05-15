package model

import "ydx-goadv-gophkeeper/internal/model/enum"

type Resource struct {
	UserId int32  `db:"user_id"`
	Data   []byte `db:"data"`
	ResourceDescription
}

type ResourceDescription struct {
	Id   int32             `db:"id"`
	Meta []byte            `db:"meta"`
	Type enum.ResourceType `db:"type"`
}
