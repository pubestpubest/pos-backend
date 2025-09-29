package seed

// Centralized seed data for easier maintenance

// Access control data
var SeedRoles = []string{
	"owner", "manager", "cashier", "waiter", "kitchen",
}

type SeedPermission struct {
	Code        string
	Description string
}

var SeedPermissions = []SeedPermission{
	{Code: "order.create", Description: "Create orders"},
	{Code: "order.update", Description: "Update orders"},
	{Code: "order.pay", Description: "Take payments"},
	{Code: "menu.manage", Description: "CRUD menu & modifiers"},
	{Code: "table.manage", Description: "CRUD tables/areas"},
	{Code: "user.manage", Description: "Manage users & roles"},
	{Code: "report.view", Description: "View reports/dashboard"},
}

var SeedRolePermissions = map[string][]string{
	"owner":   {"order.create", "order.update", "order.pay", "menu.manage", "table.manage", "user.manage", "report.view"},
	"manager": {"order.create", "order.update", "order.pay", "menu.manage", "table.manage", "report.view"},
	"cashier": {"order.pay", "report.view"},
	"waiter":  {"order.create", "order.update"},
	"kitchen": {"order.update"},
}

// Admin user
type SeedAdmin struct {
	Username     string
	Email        string
	PasswordHash string
	FullName     string
	Status       string
}

// NOTE: Replace PasswordHash at deploy time with a real bcrypt hash
var SeedAdminUser = SeedAdmin{
	Username:     "owner",
	Email:        "owner@example.com",
	PasswordHash: "$2a$12$REPLACE_ME_BCRYPT",
	FullName:     "Owner",
	Status:       "active",
}

// Venue layout
var SeedAreas = []string{"Main Hall", "Patio", "Bar"}

type SeedTable struct {
	AreaName string
	Name     string
	Seats    int
	Slug     string
}

var SeedTables = []SeedTable{
	{"Main Hall", "T1", 4, "qr_t1_x9fv"},
	{"Main Hall", "T2", 4, "qr_t2_2mrf"},
	{"Patio", "P1", 2, "qr_p1_lz3c"},
	{"Bar", "B1", 2, "qr_b1_7qk8"},
}

// Catalog
type SeedCategory struct {
	Name         string
	DisplayOrder int
}

var SeedCategories = []SeedCategory{
	{"จานเดียว", 1},
	{"เส้น", 2},
	{"กับข้าว", 3},
	{"เครื่องดื่ม", 4},
}

// รายการอาหาร/เครื่องดื่ม (หน่วย: บาท)
type SeedMenuItem struct {
	CategoryName string
	Name         string
	SKU          string
	PriceBaht    int // ใช้หน่วย "บาท" ตรง ๆ
}

var SeedMenuItems = []SeedMenuItem{
	// จานเดียว
	{"จานเดียว", "ข้าวกะเพราไก่", "RTA-GAP-KAI-01", 55},
	{"จานเดียว", "ข้าวกะเพราหมู", "RTA-GAP-MOO-01", 55},
	{"จานเดียว", "ข้าวกะเพราทะเล", "RTA-GAP-SEA-01", 75},
	{"จานเดียว", "ข้าวผัดหมู", "RTA-KAOPAD-MOO-01", 60},
	{"จานเดียว", "ข้าวผัดกุ้ง", "RTA-KAOPAD-KUNG-01", 75},
	{"จานเดียว", "ข้าวคะน้าหมูกรอบ", "RTA-KHANA-CRISPPORK-01", 70},
	{"จานเดียว", "ข้าวผัดพริกแกงไก่", "RTA-PRIKKANG-KAI-01", 60},

	// เส้น
	{"เส้น", "ผัดซีอิ๊วหมู", "RTA-SIEIW-MOO-01", 60},
	{"เส้น", "ผัดซีอิ๊วทะเล", "RTA-SIEIW-SEA-01", 80},
	{"เส้น", "ราดหน้าหมู", "RTA-RADNA-MOO-01", 65},
	{"เส้น", "ราดหน้าทะเล", "RTA-RADNA-SEA-01", 85},

	// กับข้าว (สั่งแยก/ทานร่วม)
	{"กับข้าว", "ต้มยำกุ้ง (ถ้วย)", "RTA-TOMYUM-KUNG-01", 85},
	{"กับข้าว", "ไข่เจียวหมูสับ", "RTA-KAICHIAO-MOOSAP-01", 45},
	{"กับข้าว", "คะน้าหมูกรอบ (กับ)", "RTA-KHANA-CRISPPORK-SIDE-01", 80},

	// เครื่องดื่ม
	{"เครื่องดื่ม", "ชาดำเย็น", "DRK-CHA-DAM-01", 25},
	{"เครื่องดื่ม", "ชาเย็น", "DRK-CHA-YEN-01", 30},
	{"เครื่องดื่ม", "น้ำเปล่า", "DRK-WATER-01", 10},
}

// ตัวเลือกเพิ่ม/ลด (หน่วย: บาท)
type SeedModifier struct {
	Name      string
	DeltaBaht int // หน่วย "บาท"
}

var SeedModifiers = []SeedModifier{
	{"+ไข่ดาว", 10},
	{"+ไข่เจียว", 15},
	{"+เพิ่มเนื้อ/หมู", 20},
	{"+เพิ่มข้าว", 5},
	{"-ไม่เผ็ด", 0},
	{"+เผ็ดมาก", 0},
}

// ตัวอย่างออเดอร์/จ่ายเงินสำหรับ development
type SeedSampleOrder struct {
	TableName string
	ItemSKU   string
	Quantity  int
}

var SeedDevSampleOrder = SeedSampleOrder{
	TableName: "T1",
	ItemSKU:   "RTA-GAP-MOO-01", // ข้าวกะเพราหมู
	Quantity:  2,
}
