package store

import (
	"nyota/backend/logutil"
	"nyota/backend/model"
	"nyota/backend/model/config"
	"fmt"
	"time"
)

func (store *Store) GetTenants(s *model.SessionContext) ([]*config.Tenant, error) {

	logutil.Debugf(s, "Store Layer - Get All Tenants")
	var tenants []*config.Tenant
	err := store.DB().Select(&tenants, "Select * from CCC_Tenant")
	if err != nil {
		return nil, err
	}
	return tenants, nil
}

func (store *Store) GetTenantById(s *model.SessionContext, id string) (*config.Tenant, error) {

	logutil.Debugf(s, "Store Layer - Get Tenant By Id")
	var tenant *config.Tenant
	err := store.DB().SelectOne(&tenant, "Select * from CCC_Tenant where id=$1", id)
	if err != nil {
		return nil, err
	}
	return tenant, nil
}

func (store *Store) UpsertTenant(s *model.SessionContext, data *config.Tenant) (*config.Tenant, error) {

	logutil.Debugf(s, "Store Layer - Upsert Tenant")
	if data.ID == "" {
		data.ID = fmt.Sprintf("%d", time.Now().UnixNano())
		err := store.DB().Insert(data)
		if err != nil {
			return nil, err
		}
	} else {
		_, err := store.DB().Update(data)
		if err != nil {
			return nil, err
		}
	}
	return data, nil
}

func (store *Store) DeleteTenantById(s *model.SessionContext, id string) error {

	logutil.Debugf(s, "Store Layer - Delete Tenant By Id")
	_, err := store.DB().Exec("Delete from CCC_Tenant where id=$1", id)
	if err != nil {
		return err
	}
	return nil
}
