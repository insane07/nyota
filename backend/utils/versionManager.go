package utils

import (
	"nyota/backend/model"
)

func Convert(object model.VersionEntity, destVersion int) model.VersionEntity {
	lastUpdatedVersion := object.GetVersion()
	for destVersion < lastUpdatedVersion {
		object := object.Convert()
		lastUpdatedVersion = object.GetVersion()
	}
	return object
}
