package repository

import (
	"context"
	"errors"
	"fmt"
	"math"
	"user-service/internal/core/domain/entity"
	"user-service/internal/core/domain/model"
	"user-service/utils"

	"github.com/labstack/gommon/log"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error)
	CreateUserAccount(ctx context.Context, req entity.UserEntity) (int64, error)
	UpdateUserVerified(ctx context.Context, userId int64) error
	UpdatePasswordById(ctx context.Context, req entity.UserEntity) error
	GetProfileById(ctx context.Context, userId int64) (*entity.UserEntity, error)
	UpdateProfile(ctx context.Context, req entity.UserEntity) error

	// Admin customer management functions can be added here
	GetAllCustomers(ctx context.Context, query entity.QueryStringEntity) ([]entity.UserEntity, int64, int64, error)
	GetBatchCustomers(ctx context.Context, userIds []int64) ([]entity.UserEntity, error)
	GetCustomerById(ctx context.Context, customerId int64) (*entity.UserEntity, error)
	CreateCustomer(ctx context.Context, req entity.UserEntity) (int64, error)
	UpdateCustomer(ctx context.Context, req entity.UserEntity) error
	DeleteCustomer(ctx context.Context, customerId int64) error

	getDB(ctx context.Context) *gorm.DB
}

type userRepository struct {
	db     *gorm.DB
	logger *log.Logger
}

func NewUserRepository(db *gorm.DB, logger *log.Logger) UserRepositoryInterface {
	return &userRepository{db: db, logger: logger}
}

// getDB implements [UserRepositoryInterface].
func (u *userRepository) getDB(ctx context.Context) *gorm.DB {
	if tx, ok := ctx.Value(txKey{}).(*gorm.DB); ok {
		return tx
	}

	return u.db
}

// DeleteCustomer implements UserRepositoryInterface.
func (u *userRepository) DeleteCustomer(ctx context.Context, customerId int64) error {
	var (
		db        = u.getDB(ctx)
		modelUser = model.User{ID: customerId}
	)

	if err := db.WithContext(ctx).
		Model(&modelUser).
		Association("Roles").
		Clear(); err != nil {
		u.logger.Errorf("[UserRepository-1] DeleteCustomer: %v", err)
		return err
	}

	tx := db.WithContext(ctx).Delete(&modelUser)
	if tx.Error != nil {
		u.logger.Errorf("[UserRepository-2] DeleteCustomer: %v", tx.Error)
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		err := errors.New(utils.DATA_NOT_FOUND)
		u.logger.Errorf("[UserRepository-3] DeleteCustomer: %v", err)
		return err
	}

	return nil
}

// UpdateCustomers implements UserRepositoryInterface.
func (u *userRepository) UpdateCustomer(ctx context.Context, req entity.UserEntity) error {
	var (
		db            = u.getDB(ctx)
		modelRole     model.Role
		modelUser     model.User
		filteredUsers []model.User
		foundByID     *model.User
		foundByEmail  *model.User
	)

	if err := db.WithContext(ctx).
		Select("id", "email", "is_verified").
		Where("id = ? OR email = ?", req.ID, req.Email).
		Find(&filteredUsers).Error; err != nil {
		u.logger.Errorf("[UserRepository-1] UpdateCustomer: %v", err)
		return err
	}

	// Separate query result filtered by id or email
	for i := range filteredUsers {
		if filteredUsers[i].ID == req.ID {
			foundByID = &filteredUsers[i]
		}
		if filteredUsers[i].Email == req.Email {
			foundByEmail = &filteredUsers[i]
		}
	}

	if foundByID == nil {
		err := errors.New(utils.DATA_NOT_FOUND)
		u.logger.Errorf("[UserRepository-2] UpdateCustomer: %v", err)
		return err
	}

	// Other user's email
	if foundByEmail != nil && foundByEmail.ID != req.ID {
		err := errors.New(utils.EMAIL_ALREADY_EXISTS)
		u.logger.Errorf("[UserRepository-3] UpdateCustomer: %v", err)
		return err
	}

	modelRole = model.Role{
		ID: req.RoleId,
	}

	modelUser = model.User{
		ID:      req.ID,
		Name:    req.Name,
		Email:   req.Email,
		Address: req.Address,
		Lat:     req.Lat,
		Lng:     req.Lng,
		Phone:   req.Phone,
		Photo:   req.Photo,
		Roles:   []model.Role{modelRole},
	}

	if req.Password != "" {
		modelUser.Password = req.Password
	}

	if err := db.WithContext(ctx).Updates(&modelUser).Error; err != nil {
		u.logger.Errorf("[UserRepository-4] UpdateCustomer: %v", err)
		return err
	}

	return nil
}

// CreateCustomer implements UserRepositoryInterface.
func (u *userRepository) CreateCustomer(ctx context.Context, req entity.UserEntity) (int64, error) {

	var (
		db        = u.getDB(ctx)
		modelRole model.Role
		modelUser model.User
	)

	if err := db.WithContext(ctx).
		Select("is_verified", "email").
		Where("email = ?", req.Email).
		Find(&modelUser).
		Limit(1).Error; err != nil {
		u.logger.Errorf("[UserRepository-1] CreateCustomer: %v", err)
		return 0, err
	}

	if modelUser.Email != "" {
		switch modelUser.IsVerified {
		case true:
			err := errors.New(utils.EMAIL_ALREADY_EXISTS)
			u.logger.Errorf("[UserRepository-2] CreateCustomer: %v", err)
			return 0, err
		case false:
			err := errors.New(utils.EMAIL_NOT_VERIFIED)
			u.logger.Errorf("[UserRepository-3] CreateCustomer: %v", err)
			return 0, err
		}
	}

	modelRole = model.Role{
		ID: req.RoleId,
	}

	modelUser = model.User{
		Name:       req.Name,
		Email:      req.Email,
		Password:   req.Password,
		Address:    req.Address,
		Lat:        req.Lat,
		Lng:        req.Lng,
		Phone:      req.Phone,
		Photo:      req.Photo,
		Roles:      []model.Role{modelRole},
		IsVerified: true,
	}

	if err := db.WithContext(ctx).Create(&modelUser).Error; err != nil {
		u.logger.Errorf("[UserRepository-4] CreateCustomer: %v", err)
		return 0, err
	}

	return modelUser.ID, nil
}

// GetCustomerById implements UserRepositoryInterface.
func (u *userRepository) GetCustomerById(ctx context.Context, customerId int64) (*entity.UserEntity, error) {
	var (
		db        = u.getDB(ctx)
		modelUser model.User
	)

	userSelectField := `id, name, email, address, lat, lng, phone, photo, is_verified`

	if err := db.WithContext(ctx).Select(userSelectField).
		Where("id = ?", customerId).
		Preload("Roles", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name").
				Where("name = ?", "Customer")
		}).
		First(&modelUser).Error; err != nil {
		u.logger.Errorf("[UserRepository-1] GetCustomerById: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := errors.New(utils.DATA_NOT_FOUND)
			return nil, err
		}
		return nil, err
	}

	if len(modelUser.Roles) == 0 {
		err := errors.New(utils.DATA_NOT_FOUND)
		u.logger.Errorf("[UserRepository-2] GetCustomerById: %v", err)
		return nil, err
	}

	return &entity.UserEntity{
		ID:         modelUser.ID,
		Name:       modelUser.Name,
		Email:      modelUser.Email,
		RoleId:     modelUser.Roles[0].ID,
		RoleName:   modelUser.Roles[0].Name,
		Address:    modelUser.Address,
		Lat:        modelUser.Lat,
		Lng:        modelUser.Lng,
		Phone:      modelUser.Phone,
		Photo:      modelUser.Photo,
		IsVerified: modelUser.IsVerified,
	}, nil
}

// GetBatchCustomers implements UserRepositoryInterface.
func (u *userRepository) GetBatchCustomers(ctx context.Context, userIds []int64) ([]entity.UserEntity, error) {
	var (
		db         = u.getDB(ctx)
		modelUsers []model.User
	)

	chunkSize := 150

	for i := 0; i < len(userIds); i += chunkSize {
		end := min(i+chunkSize, len(userIds))

		batchCustomers := []model.User{}
		if err := db.WithContext(ctx).
			Select("id", "name", "email", "address", "phone").
			Where("id IN ?", userIds[i:end]).
			Find(&batchCustomers).Error; err != nil {
			u.logger.Errorf("[UserRepository-1] GetBatchCustomers: %v", err)
			return nil, err
		}

		modelUsers = append(modelUsers, batchCustomers...)
	}

	if len(modelUsers) == 0 {
		err := errors.New(utils.DATA_NOT_FOUND)
		u.logger.Errorf("[UserRepository-2] GetBatchCustomers: %v", err)
		return nil, err
	}

	entities := []entity.UserEntity{}
	for _, val := range modelUsers {
		entities = append(entities, entity.UserEntity{
			ID:      val.ID,
			Name:    val.Name,
			Email:   val.Email,
			Address: val.Address,
			Phone:   val.Phone,
		})
	}

	return entities, nil
}

// GetAll implements UserRepositoryInterface.
func (u *userRepository) GetAllCustomers(ctx context.Context, query entity.QueryStringEntity) ([]entity.UserEntity, int64, int64, error) {
	var (
		db         = u.getDB(ctx)
		modelUsers []model.User
	)

	var countData int64

	orderSort := fmt.Sprintf("users.%s %s", query.OrderBy, query.OrderType)
	offset := (query.Page - 1) * query.Limit

	if err := db.WithContext(ctx).
		Table("user_role").
		Where("role_id IN (?)",
			db.WithContext(ctx).Pluck("id", &model.Role{}).
				Where("name = ?", "Customer")).
		Count(&countData).Error; err != nil {
		u.logger.Errorf("[UserRepository-1] GetAllCustomers: %v", err)
		return nil, 0, 0, err
	}

	usersSelectField := `users.id AS id,
						users.name AS name,
						users.email AS email,
						users.phone AS phone,
						users.photo AS photo`

	sqlMain := db.WithContext(ctx).
		Select(usersSelectField).
		Joins("INNER JOIN user_role ON users.id = user_role.user_id").
		Where("name ILIKE ? OR email ILIKE ? OR phone ILIKE ?", "%"+query.Search+"%", "%"+query.Search+"%", "%"+query.Search+"%").
		Where("user_role.role_id IN (?)", db.WithContext(ctx).
			Select("id").
			Model(&model.Role{}).
			Where("name = ?", "Customer"))

	totalPages := int(math.Ceil(float64(countData) / float64(query.Limit)))

	if err := sqlMain.Order(orderSort).
		Limit(int(query.Limit)).
		Offset(int(offset)).
		Find(&modelUsers).Error; err != nil {
		u.logger.Errorf("[UserRepository-2] GetAllCustomers: %v", err)
		return nil, 0, 0, err
	}

	respEntities := []entity.UserEntity{}

	for _, user := range modelUsers {
		respEntities = append(respEntities, entity.UserEntity{
			ID:    user.ID,
			Name:  user.Name,
			Email: user.Email,
			Phone: user.Phone,
			Photo: user.Photo,
		})
	}

	return respEntities, countData, int64(totalPages), nil
}

// UpdateProfile implements UserRepositoryInterface.
func (u *userRepository) UpdateProfile(ctx context.Context, req entity.UserEntity) error {
	var (
		db            = u.getDB(ctx)
		modelUser     model.User
		filteredUsers []model.User
		foundByID     *model.User
		foundByEmail  *model.User
	)

	if err := db.WithContext(ctx).
		Select("id", "email", "is_verified").
		Where("id = ? OR email = ?", req.ID, req.Email).
		Find(&filteredUsers).Error; err != nil {
		u.logger.Errorf("[UserRepository-1] UpdateProfile: %v", err)
		return err
	}

	// Separate query result filtered by id or email
	for i := range filteredUsers {
		if filteredUsers[i].ID == req.ID {
			foundByID = &filteredUsers[i]
		}
		if filteredUsers[i].Email == req.Email {
			foundByEmail = &filteredUsers[i]
		}
	}

	if foundByID == nil {
		err := errors.New(utils.DATA_NOT_FOUND)
		u.logger.Errorf("[UserRepository-2] UpdateProfile: %v", err)
		return err
	}

	if foundByEmail != nil {
		// User's email
		if foundByEmail.ID == req.ID {
			if !foundByEmail.IsVerified {
				err := errors.New(utils.EMAIL_NOT_VERIFIED)
				u.logger.Errorf("[UserRepository-3] UpdateProfile: %v", err)
				return err
			}
		} else {
			// Other user's email
			err := errors.New(utils.EMAIL_ALREADY_EXISTS)
			u.logger.Errorf("[UserRepository-4] UpdateProfile: %v", err)
			return err
		}
	}

	modelUser = model.User{
		ID:      req.ID,
		Name:    req.Name,
		Email:   req.Email,
		Address: req.Address,
		Lat:     req.Lat,
		Lng:     req.Lng,
		Phone:   req.Phone,
		Photo:   req.Photo,
	}

	if err := db.WithContext(ctx).Updates(&modelUser).Error; err != nil {
		u.logger.Errorf("[UserRepository-5] UpdateProfile: %v", err)
		return err
	}

	return nil
}

// GetProfileById implements UserRepositoryInterface.
func (u *userRepository) GetProfileById(ctx context.Context, userId int64) (*entity.UserEntity, error) {
	var (
		db        = u.getDB(ctx)
		modelUser model.User
	)

	if err := db.WithContext(ctx).
		Where("id = ? AND is_verified = ?", userId, true).
		Preload("Roles").
		First(&modelUser).Error; err != nil {
		u.logger.Errorf("[UserRepository-1] GetProfileById: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := errors.New(utils.DATA_NOT_FOUND)
			return nil, err
		}
		return nil, err
	}

	return &entity.UserEntity{
		ID:       modelUser.ID,
		Name:     modelUser.Email,
		Email:    modelUser.Name,
		RoleName: modelUser.Roles[0].Name,
		Address:  modelUser.Address,
		Lat:      modelUser.Lat,
		Lng:      modelUser.Lng,
		Phone:    modelUser.Phone,
		Photo:    modelUser.Photo,
	}, nil
}

// UpdatePasswordById implements UserRepositoryInterface.
func (u *userRepository) UpdatePasswordById(ctx context.Context, req entity.UserEntity) error {
	var (
		db        = u.getDB(ctx)
		modelUser = model.User{
			ID:       req.ID,
			Password: req.Password,
		}
	)

	tx := db.WithContext(ctx).Updates(&modelUser)
	if tx.Error != nil {
		u.logger.Errorf("[UserRepository-1] UpdatePasswordById: %v", tx.Error)
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		err := errors.New(utils.DATA_NOT_FOUND)
		u.logger.Errorf("[UserRepository-2] UpdatePasswordById: %v", err)
		return err
	}

	return nil
}

// UpdateUserVerified implements UserRepositoryInterface.
func (u *userRepository) UpdateUserVerified(ctx context.Context, userId int64) error {
	var (
		db        = u.getDB(ctx)
		modelUser = model.User{
			ID:         userId,
			IsVerified: true,
		}
	)

	tx := db.WithContext(ctx).Updates(&modelUser)
	if tx.Error != nil {
		u.logger.Errorf("[UserRepository-1] UpdateUserVerified: %v", tx.Error)
		return tx.Error
	}

	if tx.RowsAffected == 0 {
		err := errors.New(utils.DATA_NOT_FOUND)
		u.logger.Errorf("[UserRepository-2] UpdateUserVerified: %v", err)
		return err
	}

	return nil
}

// CreateUserAccount implements UserRepositoryInterface.
func (u *userRepository) CreateUserAccount(ctx context.Context, req entity.UserEntity) (int64, error) {
	var (
		db        = u.getDB(ctx)
		modelUser model.User
		modelRole model.Role
	)

	if err := db.WithContext(ctx).
		Select("is_verified").
		Where("email = ?", req.Email).
		Find(&modelUser).
		Limit(1).Error; err != nil {
		u.logger.Errorf("[UserRepository-1] CreateUserAccount: %v", err)
		return 0, err
	}

	if modelUser.Email != "" {
		switch modelUser.IsVerified {
		case true:
			err := errors.New(utils.EMAIL_ALREADY_EXISTS)
			u.logger.Errorf("[UserRepository-2] CreateUserAccount: %v", err)
			return 0, err
		case false:
			err := errors.New(utils.EMAIL_NOT_VERIFIED)
			u.logger.Errorf("[UserRepository-3] CreateUserAccount: %v", err)
			return 0, err
		}
	}

	modelRole.ID = req.RoleId

	modelUser = model.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: req.Password,
		Roles:    []model.Role{modelRole},
	}

	if err := db.WithContext(ctx).Create(&modelUser).Error; err != nil {
		u.logger.Errorf("[UserRepository-4] CreateUserAccount: %v", err)
		return 0, err
	}

	return modelUser.ID, nil
}

// GetUserByEmail implements UserRepositoryInterface.
func (u *userRepository) GetUserByEmail(ctx context.Context, email string) (*entity.UserEntity, error) {
	var (
		db        = u.getDB(ctx)
		modelUser model.User
	)

	if err := db.WithContext(ctx).
		Omit("created_at", "deleted_at", "updated_at").
		Where("email = ?", email).
		Preload("Roles", func(db *gorm.DB) *gorm.DB {
			return db.Select("id", "name")
		}).
		First(&modelUser).Error; err != nil {
		u.logger.Errorf("[UserRepository-1] GetUserByEmail: %v", err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			err := errors.New(utils.DATA_NOT_FOUND)
			return nil, err
		}
		return nil, err
	}

	return &entity.UserEntity{
		ID:         modelUser.ID,
		Name:       modelUser.Name,
		Password:   modelUser.Password,
		Email:      modelUser.Email,
		RoleName:   modelUser.Roles[0].Name,
		Address:    modelUser.Address,
		Lat:        modelUser.Lat,
		Lng:        modelUser.Lng,
		Phone:      modelUser.Phone,
		Photo:      modelUser.Photo,
		IsVerified: modelUser.IsVerified,
	}, nil
}
