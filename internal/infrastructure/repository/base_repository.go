package repository

import (
	"database/sql"
	"fmt"

	"github.com/Masterminds/squirrel"
)

// BaseRepository представляет базовый репозиторий с настроенным query builder
type BaseRepository struct {
	db      *sql.DB
	builder squirrel.StatementBuilderType
}

// NewBaseRepository создает новый базовый репозиторий
func NewBaseRepository(db *sql.DB) *BaseRepository {
	return &BaseRepository{
		db: db,
		// Настраиваем Squirrel для PostgreSQL с placeholder'ами $1, $2, etc.
		builder: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar),
	}
}

// DB возвращает подключение к базе данных
func (r *BaseRepository) DB() *sql.DB {
	return r.db
}

// Builder возвращает настроенный query builder
func (r *BaseRepository) Builder() squirrel.StatementBuilderType {
	return r.builder
}

// ExecQuery выполняет запрос и возвращает результат
func (r *BaseRepository) ExecQuery(query squirrel.Sqlizer) (sql.Result, error) {
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ошибка построения SQL запроса: %w", err)
	}

	result, err := r.db.Exec(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}

	return result, nil
}

// QueryRow выполняет запрос и возвращает одну строку
func (r *BaseRepository) QueryRow(query squirrel.Sqlizer) *sql.Row {
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		// Возвращаем row с ошибкой
		return r.db.QueryRow("SELECT $1", err.Error())
	}

	return r.db.QueryRow(sqlQuery, args...)
}

// Query выполняет запрос и возвращает множество строк
func (r *BaseRepository) Query(query squirrel.Sqlizer) (*sql.Rows, error) {
	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("ошибка построения SQL запроса: %w", err)
	}

	rows, err := r.db.Query(sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("ошибка выполнения запроса: %w", err)
	}

	return rows, nil
}

// BuildSelectQuery создает базовый SELECT запрос
func (r *BaseRepository) BuildSelectQuery(table string) squirrel.SelectBuilder {
	return r.builder.Select("*").From(table)
}

// BuildInsertQuery создает базовый INSERT запрос
func (r *BaseRepository) BuildInsertQuery(table string) squirrel.InsertBuilder {
	return r.builder.Insert(table)
}

// BuildUpdateQuery создает базовый UPDATE запрос
func (r *BaseRepository) BuildUpdateQuery(table string) squirrel.UpdateBuilder {
	return r.builder.Update(table)
}

// BuildDeleteQuery создает базовый DELETE запрос
func (r *BaseRepository) BuildDeleteQuery(table string) squirrel.DeleteBuilder {
	return r.builder.Delete(table)
}

// BuildCountQuery создает COUNT запрос
func (r *BaseRepository) BuildCountQuery(table string) squirrel.SelectBuilder {
	return r.builder.Select("COUNT(*)").From(table)
}

// BuildExistsQuery создает EXISTS запрос
func (r *BaseRepository) BuildExistsQuery(table string) squirrel.SelectBuilder {
	return r.builder.Select("EXISTS").From(table)
}

// GetLastInsertID получает ID последней вставленной записи
func (r *BaseRepository) GetLastInsertID(result sql.Result) (int64, error) {
	id, err := result.LastInsertId()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения ID: %w", err)
	}
	return id, nil
}

// GetRowsAffected получает количество затронутых строк
func (r *BaseRepository) GetRowsAffected(result sql.Result) (int64, error) {
	affected, err := result.RowsAffected()
	if err != nil {
		return 0, fmt.Errorf("ошибка получения количества затронутых строк: %w", err)
	}
	return affected, nil
}

// TableNames содержит названия таблиц
var TableNames = struct {
	Products         string
	ProductTypes     string
	Materials        string
	MaterialTypes    string
	MeasurementUnits string
	ProductMaterials string
}{
	Products:         "products",
	ProductTypes:     "product_types",
	Materials:        "materials",
	MaterialTypes:    "material_types",
	MeasurementUnits: "measurement_units",
	ProductMaterials: "product_materials",
}

// ColumnNames содержит названия колонок для каждой таблицы
var ProductColumns = struct {
	ID                       string
	Article                  string
	ProductTypeID            string
	Name                     string
	Description              string
	ImagePath                string
	MinPartnerPrice          string
	PackageDimensions        string
	Weights                  string
	QualityCertificatePath   string
	StandardNumber           string
	ProductionTime           string
	CostPrice                string
	WorkshopNumber           string
	RequiredWorkers          string
	RollWidth                string
	CreatedAt                string
	UpdatedAt                string
}{
	ID:                       "id",
	Article:                  "article",
	ProductTypeID:            "product_type_id",
	Name:                     "name",
	Description:              "description",
	ImagePath:                "image_path",
	MinPartnerPrice:          "min_partner_price",
	PackageDimensions:        "package_dimensions",
	Weights:                  "weights",
	QualityCertificatePath:   "quality_certificate_path",
	StandardNumber:           "standard_number",
	ProductionTime:           "production_time",
	CostPrice:                "cost_price",
	WorkshopNumber:           "workshop_number",
	RequiredWorkers:          "required_workers",
	RollWidth:                "roll_width",
	CreatedAt:                "created_at",
	UpdatedAt:                "updated_at",
}

var MaterialColumns = struct {
	ID                   string
	Article              string
	MaterialTypeID       string
	Name                 string
	Description          string
	MeasurementUnitID    string
	PackageQuantity      string
	CostPerUnit          string
	StockQuantity        string
	MinStockQuantity     string
	ImagePath            string
	CreatedAt            string
	UpdatedAt            string
}{
	ID:                   "id",
	Article:              "article",
	MaterialTypeID:       "material_type_id",
	Name:                 "name",
	Description:          "description",
	MeasurementUnitID:    "measurement_unit_id",
	PackageQuantity:      "package_quantity",
	CostPerUnit:          "cost_per_unit",
	StockQuantity:        "stock_quantity",
	MinStockQuantity:     "min_stock_quantity",
	ImagePath:            "image_path",
	CreatedAt:            "created_at",
	UpdatedAt:            "updated_at",
}

var ProductMaterialColumns = struct {
	ID              string
	ProductID       string
	MaterialID      string
	QuantityPerUnit string
	CreatedAt       string
	UpdatedAt       string
}{
	ID:              "id",
	ProductID:       "product_id",
	MaterialID:      "material_id",
	QuantityPerUnit: "quantity_per_unit",
	CreatedAt:       "created_at",
	UpdatedAt:       "updated_at",
} 