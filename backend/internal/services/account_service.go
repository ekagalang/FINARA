package services

import (
	"errors"
	"finara-backend/internal/models"
	"finara-backend/internal/repository"
)

type AccountService interface {
	CreateAccount(account *models.Account) error
	GetAccountByID(id uint) (*models.Account, error)
	GetAccountByCode(companyID uint, code string) (*models.Account, error)
	GetAccountsByCompanyID(companyID uint) ([]models.Account, error)
	GetAccountsByType(companyID uint, accountType models.AccountType) ([]models.Account, error)
	GetActiveAccounts(companyID uint) ([]models.Account, error)
	UpdateAccount(id uint, updates map[string]interface{}) (*models.Account, error)
	DeleteAccount(id uint) error
	InitializeDefaultAccounts(companyID uint) error
}

type accountService struct {
	accountRepo repository.AccountRepository
}

func NewAccountService(accountRepo repository.AccountRepository) AccountService {
	return &accountService{accountRepo: accountRepo}
}

func (s *accountService) CreateAccount(account *models.Account) error {
	// Validasi: cek apakah kode akun sudah ada
	existing, _ := s.accountRepo.FindByCode(account.CompanyID, account.Code)
	if existing != nil && existing.ID > 0 {
		return errors.New("account code already exists")
	}

	// Validasi: akun header tidak boleh punya balance
	if account.IsHeader && account.Balance != 0 {
		return errors.New("header account cannot have balance")
	}

	return s.accountRepo.Create(account)
}

func (s *accountService) GetAccountByID(id uint) (*models.Account, error) {
	return s.accountRepo.FindByID(id)
}

func (s *accountService) GetAccountByCode(companyID uint, code string) (*models.Account, error) {
	return s.accountRepo.FindByCode(companyID, code)
}

func (s *accountService) GetAccountsByCompanyID(companyID uint) ([]models.Account, error) {
	return s.accountRepo.FindByCompanyID(companyID)
}

func (s *accountService) GetAccountsByType(companyID uint, accountType models.AccountType) ([]models.Account, error) {
	return s.accountRepo.FindByType(companyID, accountType)
}

func (s *accountService) GetActiveAccounts(companyID uint) ([]models.Account, error) {
	return s.accountRepo.GetActiveAccounts(companyID)
}

func (s *accountService) UpdateAccount(id uint, updates map[string]interface{}) (*models.Account, error) {
	account, err := s.accountRepo.FindByID(id)
	if err != nil {
		return nil, errors.New("account not found")
	}

	// Update fields
	if name, ok := updates["name"].(string); ok {
		account.Name = name
	}
	if description, ok := updates["description"].(string); ok {
		account.Description = description
	}
	if isActive, ok := updates["is_active"].(bool); ok {
		account.IsActive = isActive
	}

	if err := s.accountRepo.Update(account); err != nil {
		return nil, err
	}

	return account, nil
}

func (s *accountService) DeleteAccount(id uint) error {
	// TODO: Validasi apakah akun sudah digunakan di transaksi
	return s.accountRepo.Delete(id)
}

func (s *accountService) InitializeDefaultAccounts(companyID uint) error {
	defaultAccounts := []models.Account{
		// ASET
		{CompanyID: companyID, Code: "1-0000", Name: "ASET", Type: models.AccountTypeAsset, Category: models.CategoryCurrentAsset, Level: 1, IsHeader: true},
		{CompanyID: companyID, Code: "1-1000", Name: "Aset Lancar", Type: models.AccountTypeAsset, Category: models.CategoryCurrentAsset, Level: 2, IsHeader: true},
		{CompanyID: companyID, Code: "1-1100", Name: "Kas", Type: models.AccountTypeAsset, Category: models.CategoryCurrentAsset, Level: 3, IsHeader: false},
		{CompanyID: companyID, Code: "1-1200", Name: "Bank", Type: models.AccountTypeAsset, Category: models.CategoryCurrentAsset, Level: 3, IsHeader: false},
		{CompanyID: companyID, Code: "1-1300", Name: "Piutang Usaha", Type: models.AccountTypeAsset, Category: models.CategoryCurrentAsset, Level: 3, IsHeader: false},
		{CompanyID: companyID, Code: "1-1400", Name: "Persediaan Barang", Type: models.AccountTypeAsset, Category: models.CategoryCurrentAsset, Level: 3, IsHeader: false},
		
		{CompanyID: companyID, Code: "1-2000", Name: "Aset Tetap", Type: models.AccountTypeAsset, Category: models.CategoryFixedAsset, Level: 2, IsHeader: true},
		{CompanyID: companyID, Code: "1-2100", Name: "Peralatan", Type: models.AccountTypeAsset, Category: models.CategoryFixedAsset, Level: 3, IsHeader: false},
		{CompanyID: companyID, Code: "1-2200", Name: "Kendaraan", Type: models.AccountTypeAsset, Category: models.CategoryFixedAsset, Level: 3, IsHeader: false},
		{CompanyID: companyID, Code: "1-2300", Name: "Gedung", Type: models.AccountTypeAsset, Category: models.CategoryFixedAsset, Level: 3, IsHeader: false},

		// LIABILITAS
		{CompanyID: companyID, Code: "2-0000", Name: "LIABILITAS", Type: models.AccountTypeLiability, Category: models.CategoryCurrentLiability, Level: 1, IsHeader: true},
		{CompanyID: companyID, Code: "2-1000", Name: "Liabilitas Lancar", Type: models.AccountTypeLiability, Category: models.CategoryCurrentLiability, Level: 2, IsHeader: true},
		{CompanyID: companyID, Code: "2-1100", Name: "Utang Usaha", Type: models.AccountTypeLiability, Category: models.CategoryCurrentLiability, Level: 3, IsHeader: false},
		{CompanyID: companyID, Code: "2-1200", Name: "Utang Pajak", Type: models.AccountTypeLiability, Category: models.CategoryCurrentLiability, Level: 3, IsHeader: false},
		
		{CompanyID: companyID, Code: "2-2000", Name: "Liabilitas Jangka Panjang", Type: models.AccountTypeLiability, Category: models.CategoryLongTermLiability, Level: 2, IsHeader: true},
		{CompanyID: companyID, Code: "2-2100", Name: "Utang Bank Jangka Panjang", Type: models.AccountTypeLiability, Category: models.CategoryLongTermLiability, Level: 3, IsHeader: false},

		// EKUITAS
		{CompanyID: companyID, Code: "3-0000", Name: "EKUITAS", Type: models.AccountTypeEquity, Category: models.CategoryEquity, Level: 1, IsHeader: true},
		{CompanyID: companyID, Code: "3-1000", Name: "Modal", Type: models.AccountTypeEquity, Category: models.CategoryEquity, Level: 2, IsHeader: false},
		{CompanyID: companyID, Code: "3-2000", Name: "Laba Ditahan", Type: models.AccountTypeEquity, Category: models.CategoryEquity, Level: 2, IsHeader: false},

		// PENDAPATAN
		{CompanyID: companyID, Code: "4-0000", Name: "PENDAPATAN", Type: models.AccountTypeRevenue, Category: models.CategoryOperatingRevenue, Level: 1, IsHeader: true},
		{CompanyID: companyID, Code: "4-1000", Name: "Pendapatan Usaha", Type: models.AccountTypeRevenue, Category: models.CategoryOperatingRevenue, Level: 2, IsHeader: false},
		{CompanyID: companyID, Code: "4-2000", Name: "Pendapatan Lain-lain", Type: models.AccountTypeRevenue, Category: models.CategoryOtherRevenue, Level: 2, IsHeader: false},

		// BEBAN
		{CompanyID: companyID, Code: "5-0000", Name: "BEBAN", Type: models.AccountTypeExpense, Category: models.CategoryOperatingExpense, Level: 1, IsHeader: true},
		{CompanyID: companyID, Code: "5-1000", Name: "Beban Operasional", Type: models.AccountTypeExpense, Category: models.CategoryOperatingExpense, Level: 2, IsHeader: true},
		{CompanyID: companyID, Code: "5-1100", Name: "Beban Gaji", Type: models.AccountTypeExpense, Category: models.CategoryOperatingExpense, Level: 3, IsHeader: false},
		{CompanyID: companyID, Code: "5-1200", Name: "Beban Sewa", Type: models.AccountTypeExpense, Category: models.CategoryOperatingExpense, Level: 3, IsHeader: false},
		{CompanyID: companyID, Code: "5-1300", Name: "Beban Listrik & Air", Type: models.AccountTypeExpense, Category: models.CategoryOperatingExpense, Level: 3, IsHeader: false},
		{CompanyID: companyID, Code: "5-1400", Name: "Beban Telepon & Internet", Type: models.AccountTypeExpense, Category: models.CategoryOperatingExpense, Level: 3, IsHeader: false},
		
		{CompanyID: companyID, Code: "5-2000", Name: "Beban Lain-lain", Type: models.AccountTypeExpense, Category: models.CategoryOtherExpense, Level: 2, IsHeader: false},
	}

	for _, account := range defaultAccounts {
		if err := s.accountRepo.Create(&account); err != nil {
			return err
		}
	}

	return nil
}