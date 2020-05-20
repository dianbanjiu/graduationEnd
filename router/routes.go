package router

import (
	"context"
	"github.com/gin-gonic/gin"
	"graduationEnd/controller"
	"graduationEnd/middleware"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func StartRouter() {
	// 创建路由
	r := gin.Default()
	// 跨域解决
	r.Use(middleware.Cros())

	// 登录认证
	r.POST("/api/loginAuth", controller.LoginAuth)
	// 管理员路由组
	adminRoutes := r.Group("/api/admin", middleware.AuthMiddleWare())
	{
		// 公告管理
		// 添加公告
		adminRoutes.POST("/addBoard", controller.AddBoard)
		// 删除公告
		adminRoutes.POST("/deleteBoard", controller.DeleteBoard)

		// 获取所有学生/教师的信息
		adminRoutes.GET("/getUsersInfo", controller.GetUsersInfo)
		// 删除给定学生/教师们
		adminRoutes.POST("/deleteUsers", controller.DeleteUsers)
		// 添加学生/教师
		adminRoutes.POST("/addUser", controller.AddUser)
		// 添加学生/教师们
		adminRoutes.POST("/addUsers", controller.AddUsers)
		// 实训管理
		// 删除指定的多个实训
		adminRoutes.POST("/deleteCourses", controller.DeleteCourses)
		// 添加单个实训
		adminRoutes.POST("/addCourse", controller.AddCourse)
		// 从文件导入多个实训
		adminRoutes.POST("/addCourses", controller.AddCourses)

		// 查看所有实训申请
		adminRoutes.GET("/viewAllApplyCourse", controller.ViewAllApplyProgress)
		// 处理实训申请
		adminRoutes.POST("/handleApplyCourse",controller.HandleApplyCourse)
	}

	// 教师路由组
	teacherRoutes := r.Group("/api/teacher", middleware.AuthMiddleWare())
	{
		// 公告管理
		// 添加公告
		teacherRoutes.POST("/addBoard", controller.AddBoard)
		// 删除公告
		teacherRoutes.POST("/deleteBoard", controller.DeleteBoard)

		//查看自己管理实训的所有学生
		teacherRoutes.GET("/viewAllSelectedStudents", controller.ViewAllSelectedStudents)

		// 查看自己管理的实训的所有学生的周报
		teacherRoutes.GET("/viewAllPublications", controller.ViewAllPublication)
		// 为周报评价
		teacherRoutes.POST("/evaluationAndScore", controller.EvaluationAndScore)
		// 删除周报评价
		teacherRoutes.POST("/deleteEvaluation", controller.DeleteEvaluation)
		// 查询学生成绩的分布
		teacherRoutes.GET("/studentScoreTimes",controller.StudentScoreTimes)

		// 查看学生的实训申请
		teacherRoutes.GET("/viewAllApplyProgress", controller.ViewAllApplyProgress)
		// 处理实训申请
		teacherRoutes.POST("/handleApplyCourse", controller.HandleApplyCourse)
	}

	// 学生路由组
	studentRoutes := r.Group("/api/student", middleware.AuthMiddleWare())
	{
		// 实训管理
		// 查看已选课程
		studentRoutes.GET("/alreadySelectedCourse", controller.AlreadySelectedCourse)
		// 退选
		studentRoutes.POST("/dropCourse", controller.DropCourse)
		// 选课
		studentRoutes.POST("/selectCourse", controller.SelectCourse)

		// 周报管理
		// 添加周报
		studentRoutes.POST("/addPublication", controller.AddPublication)
		// 查看自己的所有周报
		studentRoutes.GET("/viewAllPublications", controller.ViewAllPublication)
		// 自主申请实训
		studentRoutes.POST("/applyCourse", controller.ApplyCourse)
		// 查看申请进度
		studentRoutes.GET("/viewApplyProgress", controller.ViewApplyProgress)
	}

	// 个人信息获取
	r.GET("/api/getUserInfo", middleware.AuthMiddleWare(), controller.GetUserInfo)

	// get boards info
	r.GET("/api/getBoards", middleware.AuthMiddleWare(), controller.GetBoards)

	// 修改手机及密码
	r.POST("api/changeInfo", middleware.AuthMiddleWare(), controller.ChangeInfo)

	// 获取实训信息
	r.GET("/api/getCourses", middleware.AuthMiddleWare(), controller.GetCourses)
	srv := &http.Server{Addr: ":9091", Handler: r}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("shutdown server ...")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown: ", err)
	}
	log.Println("server exiting")
}
