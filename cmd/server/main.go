package main

import (
	"context"
	"fmt"
	"time"

	"le-studio-api/config"
	"le-studio-api/internal/domain"
	"le-studio-api/internal/handler"
	"le-studio-api/internal/middleware"
	"le-studio-api/internal/repository"
	pgrepo "le-studio-api/internal/repository/postgres"
	"le-studio-api/internal/service"
	"le-studio-api/pkg/cloudinary"
	"le-studio-api/pkg/database"
	"le-studio-api/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// main starts the HTTP API server.
func main() {
	cfg, err := config.Load()
	if err != nil {
		panic(err)
	}
	logger, _ := zap.NewProduction()
	defer logger.Sync()

	db, err := database.Connect(cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBSSLMode, cfg.Environment == "development")
	if err != nil {
		logger.Fatal("db connection failed", zap.Error(err))
	}
	if err := autoMigrate(db); err != nil {
		logger.Fatal("migration failed", zap.Error(err))
	}
	if err := seed(context.Background(), db); err != nil {
		logger.Fatal("seed failed", zap.Error(err))
	}

	repos := repository.Repositories{Admins: pgrepo.NewAdminRepo(db), Users: pgrepo.NewUserRepo(db), Coaches: pgrepo.NewCoachRepo(db), Training: pgrepo.NewTrainingTypeRepo(db), Templates: pgrepo.NewPackTemplateRepo(db), UserPacks: pgrepo.NewUserPackRepo(db), Schedules: pgrepo.NewWeeklyScheduleRepo(db), Slots: pgrepo.NewSlotRepo(db), Bookings: pgrepo.NewBookingRepo(db), Attendance: pgrepo.NewAttendanceRepo(db), RefreshToken: pgrepo.NewRefreshTokenRepo(db)}
	v := validator.New()
	authSvc := service.NewAuthService(repos, db, cfg.JWTSecret, cfg.AccessTokenDuration, cfg.RefreshTokenDuration)
	userSvc := service.NewUserService(repos)
	coachSvc := service.NewCoachService(repos)
	trainingSvc := service.NewTrainingTypeService(repos)
	tplSvc := service.NewPackTemplateService(repos)
	upSvc := service.NewUserPackService(repos)
	schSvc := service.NewScheduleService(repos)
	bookingSvc := service.NewBookingService(repos, db)
	attendanceSvc := service.NewAttendanceService(repos)
	adminSvc := service.NewAdminService(repos, db)

	cldClient, err := cloudinary.New(cfg.CloudinaryCloudName, cfg.CloudinaryAPIKey, cfg.CloudinaryAPISecret)
	if err != nil {
		logger.Warn("cloudinary initialization failed", zap.Error(err))
	}

	authH := handler.NewAuthHandler(authSvc, v, cfg.AdminSetupKey)
	userH := handler.NewUserHandler(userSvc, v)
	coachH := handler.NewCoachHandler(coachSvc, v)
	coachH.SetCloudinaryClient(cldClient)
	trainingH := handler.NewTrainingTypeHandler(trainingSvc, v)
	tplH := handler.NewPackTemplateHandler(tplSvc, v)
	upH := handler.NewUserPackHandler(upSvc, v)
	schH := handler.NewScheduleHandler(schSvc, v)
	bookingH := handler.NewBookingHandler(bookingSvc, v)
	attendanceH := handler.NewAttendanceHandler(attendanceSvc, v)
	adminH := handler.NewAdminHandler(adminSvc, v)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.CORS(cfg.AllowedOrigins))
	r.Use(middleware.RequestLogger(logger))
	r.Use(middleware.RateLimiter(200, time.Minute))
	r.GET("/healthz", func(c *gin.Context) { response.OK(c, gin.H{"status": "ok"}) })

	v1 := r.Group("/api/v1")
	registerRoutes(v1, cfg.JWTSecret, authH, userH, coachH, trainingH, tplH, upH, schH, bookingH, attendanceH, adminH)

	addr := fmt.Sprintf(":%s", cfg.Port)
	logger.Info("server starting", zap.String("addr", addr))
	if err := r.Run(addr); err != nil {
		logger.Fatal("server stopped", zap.Error(err))
	}
}

func registerRoutes(v1 *gin.RouterGroup, jwtSecret string, authH *handler.AuthHandler, userH *handler.UserHandler, coachH *handler.CoachHandler, trainingH *handler.TrainingTypeHandler, tplH *handler.PackTemplateHandler, upH *handler.UserPackHandler, schH *handler.ScheduleHandler, bookingH *handler.BookingHandler, attendanceH *handler.AttendanceHandler, adminH *handler.AdminHandler) {
	auth := v1.Group("/auth")
	auth.POST("/register", authH.Register)
	auth.POST("/login", authH.Login)
	auth.POST("/refresh", authH.Refresh)
	auth.POST("/logout", authH.Logout)
	auth.POST("/guest", authH.Guest)

	v1.GET("/coaches", coachH.List)
	v1.GET("/coaches/:id", coachH.Get)
	v1.GET("/pack-templates", tplH.List)
	v1.GET("/pack-templates/:id", tplH.Get)
	v1.GET("/training-types", trainingH.List)
	v1.GET("/schedules", schH.List)
	v1.GET("/schedules/:id", schH.Get)
	v1.GET("/schedules/:id/slots", schH.ListSlots)

	protected := v1.Group("")
	protected.Use(middleware.Auth(jwtSecret))
	protected.GET("/users/me", userH.GetMe)
	protected.PATCH("/users/me", userH.PatchMe)
	protected.PATCH("/users/me/password", userH.ChangePassword)
	protected.GET("/users/me/packs", upH.ListByUser)
	protected.GET("/users/me/bookings", bookingH.ListByUser)
	protected.POST("/user-packs", upH.Purchase)
	protected.POST("/bookings", bookingH.Create)
	protected.GET("/bookings/:id", bookingH.Get)
	protected.PATCH("/bookings/:id/cancel", bookingH.Cancel)

	adminAuth := v1.Group("/admin/auth")
	adminAuth.POST("/register", authH.AdminRegister)
	adminAuth.POST("/login", authH.AdminLogin)
	adminAuth.POST("/refresh", authH.Refresh)
	adminAuth.POST("/logout", authH.Logout)

	admin := v1.Group("/admin")
	admin.Use(middleware.Auth(jwtSecret), middleware.AdminOnly())
	admin.GET("/auth/me", adminH.GetMe)
	admin.PATCH("/auth/me", adminH.UpdateMe)
	admin.GET("/stats/overview", adminH.StatsOverview)
	// Allow admins to list coaches via the admin namespace as well
	admin.GET("/coaches", coachH.List)
	admin.GET("/admins", adminH.ListAdmins)
	admin.POST("/admins", adminH.CreateAdmin)
	admin.GET("/users", userH.AdminListUsers)
	admin.GET("/users/:id", userH.AdminGetUser)
	admin.POST("/users", userH.AdminCreateUser)
	admin.PATCH("/users/:id", userH.AdminUpdateUser)
	admin.DELETE("/users/:id", userH.AdminDeleteUser)
	admin.PATCH("/users/:id/promote", userH.AdminPromoteGuest)
	admin.PATCH("/users/:id/toggle-active", userH.AdminToggleActive)
	admin.POST("/coaches", coachH.AdminCreate)
	admin.PATCH("/coaches/:id", coachH.AdminUpdate)
	admin.DELETE("/coaches/:id", coachH.AdminDelete)
	admin.PATCH("/coaches/:id/toggle-active", coachH.AdminToggleActive)
	admin.POST("/coaches/:id/photo", coachH.AdminUploadPhoto)
	admin.POST("/training-types", trainingH.AdminCreate)
	admin.GET("/training-types", trainingH.List)
	admin.PATCH("/training-types/:id", trainingH.AdminUpdate)
	admin.DELETE("/training-types/:id", trainingH.AdminDelete)
	admin.POST("/pack-templates", tplH.AdminCreate)
	admin.PATCH("/pack-templates/:id", tplH.AdminUpdate)
	admin.DELETE("/pack-templates/:id", tplH.AdminDelete)
	admin.PATCH("/pack-templates/:id/reorder", tplH.AdminReorder)
	admin.GET("/user-packs", upH.AdminList)
	admin.GET("/user-packs/:id", upH.AdminGet)
	admin.PATCH("/user-packs/:id", upH.AdminUpdate)
	admin.PATCH("/user-packs/:id/mark-paid", upH.AdminMarkPaid)
	admin.PATCH("/user-packs/:id/adjust", upH.AdminAdjust)
	admin.DELETE("/user-packs/:id", upH.AdminDelete)
	admin.POST("/schedules", schH.AdminCreate)
	admin.GET("/schedules", schH.AdminList)
	admin.GET("/schedules/:id", schH.Get)
	admin.PATCH("/schedules/:id", schH.AdminUpdate)
	admin.POST("/schedules/:id/publish", schH.AdminPublish)
	admin.DELETE("/schedules/:id", schH.AdminDelete)
	admin.POST("/schedules/:id/slots", schH.AdminCreateSlot)
	admin.PATCH("/slots/:id", schH.AdminUpdateSlot)
	admin.DELETE("/slots/:id", schH.AdminCancelSlot)
	admin.GET("/bookings", bookingH.AdminList)
	admin.PATCH("/bookings/:id/cancel", bookingH.AdminCancel)
	admin.POST("/attendance", attendanceH.Mark)
	admin.GET("/attendance", attendanceH.List)
	admin.PATCH("/attendance/:id", attendanceH.Update)
	admin.DELETE("/attendance/:id", attendanceH.Delete)
}

func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(&domain.Admin{}, &domain.User{}, &domain.Coach{}, &domain.TrainingType{}, &domain.PackTemplate{}, &domain.UserPack{}, &domain.WeeklySchedule{}, &domain.Slot{}, &domain.Booking{}, &domain.Attendance{}, &domain.RefreshToken{})
}

func seed(ctx context.Context, db *gorm.DB) error {
	var count int64
	if err := db.WithContext(ctx).Model(&domain.Admin{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		h, err := bcrypt.GenerateFromPassword([]byte("admin123"), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		if err := db.WithContext(ctx).Create(&domain.Admin{Name: "Default Admin", Email: "admin@example.com", PasswordHash: string(h)}).Error; err != nil {
			return err
		}
	}
	if err := db.WithContext(ctx).Model(&domain.TrainingType{}).Count(&count).Error; err != nil {
		return err
	}
	if count == 0 {
		types := []domain.TrainingType{{Code: "A", Name: "Pilates A", Color: "#D96557", IsActive: true}, {Code: "B", Name: "Pilates B", Color: "#3A6EA5", IsActive: true}, {Code: "C", Name: "Pilates C", Color: "#6B8E23", IsActive: true}}
		if err := db.WithContext(ctx).Create(&types).Error; err != nil {
			return err
		}
		for _, t := range types {
			templates := []domain.PackTemplate{{Name: t.Code + " - 1 Session", NumberOfSessions: 1, Price: 20, TrainingTypeID: t.ID, IsActive: true, DisplayOrder: 1}, {Name: t.Code + " - 5 Sessions", NumberOfSessions: 5, Price: 90, TrainingTypeID: t.ID, IsActive: true, DisplayOrder: 2}, {Name: t.Code + " - 10 Sessions", NumberOfSessions: 10, Price: 170, TrainingTypeID: t.ID, IsActive: true, DisplayOrder: 3}, {Name: t.Code + " - 20 Sessions", NumberOfSessions: 20, Price: 320, TrainingTypeID: t.ID, IsActive: true, DisplayOrder: 4}}
			if err := db.WithContext(ctx).Create(&templates).Error; err != nil {
				return err
			}
		}
	}
	return nil
}
