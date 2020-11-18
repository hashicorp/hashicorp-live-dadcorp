package api

import (
	"errors"

	"github.com/hashicorp/go-memdb"
)

var (
	ErrAccessPolicyNotFound            = errors.New("access policy not found")
	ErrAccessPolicyAlreadyExists       = errors.New("access policy already exists")
	ErrConsulClusterNotFound           = errors.New("consul cluster not found")
	ErrConsulClusterAlreadyExists      = errors.New("consul cluster already exists")
	ErrNomadClusterNotFound            = errors.New("nomad cluster not found")
	ErrNomadClusterAlreadyExists       = errors.New("nomad cluster already exists")
	ErrVaultClusterNotFound            = errors.New("vault cluster not found")
	ErrVaultClusterAlreadyExists       = errors.New("vault cluster already exists")
	ErrTerraformWorkspaceNotFound      = errors.New("terraform workspace not found")
	ErrTerraformWorkspaceAlreadyExists = errors.New("terraform workspace already exists")
)

type Storer struct {
	db *memdb.MemDB
}

func NewStorer() (*Storer, error) {
	db, err := memdb.NewMemDB(&memdb.DBSchema{
		Tables: map[string]*memdb.TableSchema{
			"accessPolicy": {
				Name: "accessPolicy",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID", Lowercase: true},
					},
				},
			},
			"nomadCluster": {
				Name: "nomadCluster",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID", Lowercase: true},
					},
				},
			},
			"vaultCluster": {
				Name: "vaultCluster",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID", Lowercase: true},
					},
				},
			},
			"consulCluster": {
				Name: "consulCluster",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID", Lowercase: true},
					},
				},
			},
			"terraformWorkspace": {
				Name: "terraformWorkspace",
				Indexes: map[string]*memdb.IndexSchema{
					"id": {
						Name:    "id",
						Unique:  true,
						Indexer: &memdb.StringFieldIndex{Field: "ID", Lowercase: true},
					},
				},
			},
		},
	})
	if err != nil {
		return nil, err
	}
	return &Storer{
		db: db,
	}, nil
}

func (s *Storer) GetAccessPolicy(id string) (AccessPolicy, error) {
	txn := s.db.Txn(false)
	ap, err := txn.First("accessPolicy", "id", id)
	if err != nil {
		return AccessPolicy{}, err
	}
	if ap == nil {
		return AccessPolicy{}, ErrAccessPolicyNotFound
	}
	return *ap.(*AccessPolicy), nil
}

func (s *Storer) CreateAccessPolicy(ap AccessPolicy) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	exists, err := txn.First("accessPolicy", "id", ap.ID)
	if err != nil {
		return err
	}
	if exists != nil {
		return ErrAccessPolicyAlreadyExists
	}
	err = txn.Insert("accessPolicy", &ap)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) UpdateAccessPolicy(ap AccessPolicy) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("accessPolicy", "id", ap.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAccessPolicyNotFound
	}
	err = txn.Insert("accessPolicy", &ap)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) DeleteAccessPolicy(id string) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("accessPolicy", "id", id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrAccessPolicyNotFound
	}
	err = txn.Delete("accessPolicy", existing)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) GetConsulCluster(id string) (ConsulCluster, error) {
	txn := s.db.Txn(false)
	cluster, err := txn.First("consulCluster", "id", id)
	if err != nil {
		return ConsulCluster{}, err
	}
	if cluster == nil {
		return ConsulCluster{}, ErrConsulClusterNotFound
	}
	return *cluster.(*ConsulCluster), nil
}

func (s *Storer) CreateConsulCluster(cluster ConsulCluster) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	exists, err := txn.First("consulCluster", "id", cluster.ID)
	if err != nil {
		return err
	}
	if exists != nil {
		return ErrConsulClusterAlreadyExists
	}
	err = txn.Insert("consulCluster", &cluster)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) UpdateConsulCluster(cluster ConsulCluster) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("consulCluster", "id", cluster.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrConsulClusterNotFound
	}
	err = txn.Insert("consulCluster", &cluster)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) DeleteConsulCluster(id string) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("consulCluster", "id", id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrConsulClusterNotFound
	}
	err = txn.Delete("consulCluster", existing)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) GetVaultCluster(id string) (VaultCluster, error) {
	txn := s.db.Txn(false)
	cluster, err := txn.First("vaultCluster", "id", id)
	if err != nil {
		return VaultCluster{}, err
	}
	if cluster == nil {
		return VaultCluster{}, ErrVaultClusterNotFound
	}
	return *cluster.(*VaultCluster), nil
}

func (s *Storer) CreateVaultCluster(cluster VaultCluster) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	exists, err := txn.First("vaultCluster", "id", cluster.ID)
	if err != nil {
		return err
	}
	if exists != nil {
		return ErrVaultClusterAlreadyExists
	}
	err = txn.Insert("vaultCluster", &cluster)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) UpdateVaultCluster(cluster VaultCluster) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("vaultCluster", "id", cluster.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrVaultClusterNotFound
	}
	err = txn.Insert("vaultCluster", &cluster)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) DeleteVaultCluster(id string) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("vaultCluster", "id", id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrVaultClusterNotFound
	}
	err = txn.Delete("vaultCluster", existing)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) GetNomadCluster(id string) (NomadCluster, error) {
	txn := s.db.Txn(false)
	cluster, err := txn.First("nomadCluster", "id", id)
	if err != nil {
		return NomadCluster{}, err
	}
	if cluster == nil {
		return NomadCluster{}, ErrNomadClusterNotFound
	}
	return *cluster.(*NomadCluster), nil
}

func (s *Storer) CreateNomadCluster(cluster NomadCluster) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	exists, err := txn.First("nomadCluster", "id", cluster.ID)
	if err != nil {
		return err
	}
	if exists != nil {
		return ErrNomadClusterAlreadyExists
	}
	err = txn.Insert("nomadCluster", &cluster)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) UpdateNomadCluster(cluster NomadCluster) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("nomadCluster", "id", cluster.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNomadClusterNotFound
	}
	err = txn.Insert("nomadCluster", &cluster)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) DeleteNomadCluster(id string) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("nomadCluster", "id", id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrNomadClusterNotFound
	}
	err = txn.Delete("nomadCluster", existing)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) GetTerraformWorkspace(id string) (TerraformWorkspace, error) {
	txn := s.db.Txn(false)
	cluster, err := txn.First("terraformWorkspace", "id", id)
	if err != nil {
		return TerraformWorkspace{}, err
	}
	if cluster == nil {
		return TerraformWorkspace{}, ErrTerraformWorkspaceNotFound
	}
	return *cluster.(*TerraformWorkspace), nil
}

func (s *Storer) CreateTerraformWorkspace(workspace TerraformWorkspace) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	exists, err := txn.First("terraformWorkspace", "id", workspace.ID)
	if err != nil {
		return err
	}
	if exists != nil {
		return ErrTerraformWorkspaceAlreadyExists
	}
	err = txn.Insert("terraformWorkspace", &workspace)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) UpdateTerraformWorkspace(workspace TerraformWorkspace) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("terraformWorkspace", "id", workspace.ID)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrTerraformWorkspaceNotFound
	}
	err = txn.Insert("terraformWorkspace", &workspace)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}

func (s *Storer) DeleteTerraformWorkspace(id string) error {
	txn := s.db.Txn(true)
	defer txn.Abort()
	existing, err := txn.First("terraformWorkspace", "id", id)
	if err != nil {
		return err
	}
	if existing == nil {
		return ErrTerraformWorkspaceNotFound
	}
	err = txn.Delete("terraformWorkspace", existing)
	if err != nil {
		return err
	}
	txn.Commit()
	return nil
}
